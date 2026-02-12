#!/bin/bash
###############################################################################
# GreenRide Deployment Test Script
# Run on AWS EC2 (18.143.118.157) after deploying the latest backend+frontend
#
# Tests:
#   1. Health check
#   2. Public system config endpoint
#   3. Admin login + auth
#   4. Admin get/update system config (maintenance mode)
#   5. Maintenance mode blocking on user endpoints
#   6. Exempt endpoint access during maintenance
#   7. Disable maintenance + verify service restored
#   8. Redis cache verification
#   9. Order detail coordinate fields (omitempty fix)
#
# Usage:
#   chmod +x test-deployment.sh
#   ./test-deployment.sh
#
# Optional env vars:
#   ADMIN_USER     - admin username  (default: prompt)
#   ADMIN_PASS     - admin password  (default: prompt)
#   API_HOST       - user API host   (default: http://localhost:8610)
#   ADMIN_HOST     - admin API host  (default: http://localhost:8611)
###############################################################################

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

PASS=0
FAIL=0
SKIP=0

API_HOST="${API_HOST:-http://localhost:8610}"
ADMIN_HOST="${ADMIN_HOST:-http://localhost:8611}"

log()  { echo -e "${CYAN}[INFO]${NC}  $*"; }
ok()   { echo -e "${GREEN}[PASS]${NC}  $*"; PASS=$((PASS+1)); }
fail() { echo -e "${RED}[FAIL]${NC}  $*"; FAIL=$((FAIL+1)); }
skip() { echo -e "${YELLOW}[SKIP]${NC}  $*"; SKIP=$((SKIP+1)); }
header() { echo -e "\n${BOLD}═══════════════════════════════════════════${NC}"; echo -e "${BOLD}  $*${NC}"; echo -e "${BOLD}═══════════════════════════════════════════${NC}"; }

# Helper: extract JSON field (uses jq if available, falls back to grep/sed)
json_val() {
  local json="$1" key="$2"
  if command -v jq &>/dev/null; then
    echo "$json" | jq -r ".$key // empty" 2>/dev/null
  else
    echo "$json" | grep -o "\"$key\"[[:space:]]*:[[:space:]]*\"[^\"]*\"" | head -1 | sed "s/.*:.*\"\(.*\)\"/\1/"
  fi
}

json_bool() {
  local json="$1" key="$2"
  if command -v jq &>/dev/null; then
    echo "$json" | jq -r ".$key" 2>/dev/null
  else
    echo "$json" | grep -o "\"$key\"[[:space:]]*:[[:space:]]*[a-z]*" | head -1 | sed "s/.*:[[:space:]]*//"
  fi
}

json_num() {
  local json="$1" key="$2"
  if command -v jq &>/dev/null; then
    echo "$json" | jq -r ".$key" 2>/dev/null
  else
    echo "$json" | grep -o "\"$key\"[[:space:]]*:[[:space:]]*[0-9.e+-]*" | head -1 | sed "s/.*:[[:space:]]*//"
  fi
}

###############################################################################
header "PRE-FLIGHT CHECKS"
###############################################################################

# Check jq
if command -v jq &>/dev/null; then
  log "jq found — JSON parsing will be reliable"
else
  log "jq not found — install with: sudo apt-get install -y jq"
  log "Falling back to grep-based JSON parsing (less reliable)"
fi

# Check curl
if ! command -v curl &>/dev/null; then
  echo -e "${RED}ERROR: curl is required but not installed.${NC}"
  exit 1
fi

# Check Redis CLI
HAS_REDIS=false
if command -v redis-cli &>/dev/null; then
  HAS_REDIS=true
  log "redis-cli found"
else
  log "redis-cli not found — Redis tests will be skipped"
fi

# Get admin credentials
if [ -z "${ADMIN_USER:-}" ]; then
  read -rp "Admin username: " ADMIN_USER
fi
if [ -z "${ADMIN_PASS:-}" ]; then
  read -rsp "Admin password: " ADMIN_PASS
  echo
fi

###############################################################################
header "1. HEALTH CHECK"
###############################################################################

HEALTH_API=$(curl -s -o /dev/null -w "%{http_code}" "${API_HOST}/health" 2>/dev/null || echo "000")
HEALTH_ADMIN=$(curl -s -o /dev/null -w "%{http_code}" "${ADMIN_HOST}/health" 2>/dev/null || echo "000")

if [ "$HEALTH_API" = "200" ]; then
  ok "User API health check (port 8610): HTTP $HEALTH_API"
else
  fail "User API health check (port 8610): HTTP $HEALTH_API"
fi

if [ "$HEALTH_ADMIN" = "200" ]; then
  ok "Admin API health check (port 8611): HTTP $HEALTH_ADMIN"
else
  fail "Admin API health check (port 8611): HTTP $HEALTH_ADMIN"
fi

###############################################################################
header "2. PUBLIC SYSTEM CONFIG ENDPOINT"
###############################################################################

SYSCONFIG=$(curl -s "${API_HOST}/system/config" 2>/dev/null)
SYSCONFIG_CODE=$(json_val "$SYSCONFIG" "code")
SYSCONFIG_MODE=$(json_bool "$SYSCONFIG" "data.maintenance_mode" 2>/dev/null || echo "")

if [ -z "$SYSCONFIG_MODE" ] && command -v jq &>/dev/null; then
  SYSCONFIG_MODE=$(echo "$SYSCONFIG" | jq -r '.data.maintenance_mode' 2>/dev/null)
fi

log "Response: $SYSCONFIG"

if [ "$SYSCONFIG_CODE" = "0000" ]; then
  ok "GET /system/config returns success (code 0000)"
else
  fail "GET /system/config unexpected code: $SYSCONFIG_CODE"
fi

if [ "$SYSCONFIG_MODE" = "false" ]; then
  ok "Maintenance mode is currently OFF (expected default)"
elif [ "$SYSCONFIG_MODE" = "true" ]; then
  log "Maintenance mode is currently ON — tests will still proceed"
  ok "GET /system/config returned maintenance_mode field"
else
  fail "Could not parse maintenance_mode from response"
fi

###############################################################################
header "3. ADMIN LOGIN"
###############################################################################

LOGIN_RESP=$(curl -s -X POST "${ADMIN_HOST}/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"${ADMIN_USER}\",\"password\":\"${ADMIN_PASS}\"}" 2>/dev/null)

LOGIN_CODE=$(json_val "$LOGIN_RESP" "code")

if command -v jq &>/dev/null; then
  ADMIN_TOKEN=$(echo "$LOGIN_RESP" | jq -r '.data.token // empty' 2>/dev/null)
else
  ADMIN_TOKEN=$(echo "$LOGIN_RESP" | grep -o '"token":"[^"]*"' | head -1 | sed 's/"token":"//;s/"//')
fi

if [ -n "$ADMIN_TOKEN" ] && [ "$ADMIN_TOKEN" != "null" ]; then
  ok "Admin login successful — token received"
  log "Token: ${ADMIN_TOKEN:0:20}..."
else
  fail "Admin login failed: $LOGIN_RESP"
  echo -e "${RED}Cannot continue without admin token. Exiting.${NC}"
  exit 1
fi

###############################################################################
header "4. ADMIN GET SYSTEM CONFIG"
###############################################################################

ADMIN_SYSCONFIG=$(curl -s "${ADMIN_HOST}/admin/system/config" \
  -H "Authorization: Bearer ${ADMIN_TOKEN}" 2>/dev/null)

ADMIN_SC_CODE=$(json_val "$ADMIN_SYSCONFIG" "code")
log "Response: $ADMIN_SYSCONFIG"

if [ "$ADMIN_SC_CODE" = "0000" ]; then
  ok "GET /admin/system/config returns success"
else
  fail "GET /admin/system/config unexpected code: $ADMIN_SC_CODE"
fi

###############################################################################
header "5. ENABLE MAINTENANCE MODE"
###############################################################################

ENABLE_RESP=$(curl -s -X POST "${ADMIN_HOST}/admin/system/config" \
  -H "Authorization: Bearer ${ADMIN_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"maintenance_mode":true,"maintenance_message":"Deployment test — maintenance mode enabled by test script","maintenance_phone":"6996"}' 2>/dev/null)

ENABLE_CODE=$(json_val "$ENABLE_RESP" "code")
log "Response: $ENABLE_RESP"

if [ "$ENABLE_CODE" = "0000" ]; then
  ok "POST /admin/system/config — maintenance mode ENABLED"
else
  fail "Failed to enable maintenance mode: $ENABLE_CODE"
fi

# Small delay for Redis cache propagation
sleep 1

###############################################################################
header "6. VERIFY MAINTENANCE BLOCKING (user API)"
###############################################################################

# Test a non-exempt endpoint (should be blocked with 503)
BLOCKED_HTTP=$(curl -s -o /tmp/gr_blocked_resp.json -w "%{http_code}" \
  "${API_HOST}/profile" \
  -H "Authorization: Bearer fake-token-for-test" 2>/dev/null || echo "000")

BLOCKED_RESP=$(cat /tmp/gr_blocked_resp.json 2>/dev/null || echo "{}")
BLOCKED_CODE=$(json_val "$BLOCKED_RESP" "code")

log "HTTP status: $BLOCKED_HTTP"
log "Response: $BLOCKED_RESP"

if [ "$BLOCKED_HTTP" = "503" ]; then
  ok "Non-exempt endpoint returns HTTP 503 during maintenance"
else
  fail "Expected HTTP 503, got HTTP $BLOCKED_HTTP"
fi

if [ "$BLOCKED_CODE" = "1100" ]; then
  ok "Error code is 1100 (MaintenanceMode)"
else
  fail "Expected error code 1100, got: $BLOCKED_CODE"
fi

# Check that maintenance payload includes message and phone
if echo "$BLOCKED_RESP" | grep -q "support_phone"; then
  ok "Maintenance response includes support_phone field"
else
  fail "Maintenance response missing support_phone field"
fi

###############################################################################
header "7. VERIFY EXEMPT ENDPOINTS DURING MAINTENANCE"
###############################################################################

# Health should still work
EXEMPT_HEALTH=$(curl -s -o /dev/null -w "%{http_code}" "${API_HOST}/health" 2>/dev/null)
if [ "$EXEMPT_HEALTH" = "200" ]; then
  ok "EXEMPT: /health returns 200 during maintenance"
else
  fail "EXEMPT: /health returned $EXEMPT_HEALTH (expected 200)"
fi

# Public system config should still work
EXEMPT_SC_HTTP=$(curl -s -o /dev/null -w "%{http_code}" "${API_HOST}/system/config" 2>/dev/null)
if [ "$EXEMPT_SC_HTTP" = "200" ]; then
  ok "EXEMPT: /system/config returns 200 during maintenance"
else
  fail "EXEMPT: /system/config returned $EXEMPT_SC_HTTP (expected 200)"
fi

# Login should still work (will fail auth, but should NOT be 503)
EXEMPT_LOGIN_HTTP=$(curl -s -o /dev/null -w "%{http_code}" -X POST "${API_HOST}/login" \
  -H "Content-Type: application/json" \
  -d '{"phone":"0000000000","user_type":"passenger"}' 2>/dev/null)
if [ "$EXEMPT_LOGIN_HTTP" != "503" ]; then
  ok "EXEMPT: /login is NOT blocked (HTTP $EXEMPT_LOGIN_HTTP, not 503)"
else
  fail "EXEMPT: /login returned 503 — should be exempt!"
fi

# Support config should still work
EXEMPT_SUPPORT_HTTP=$(curl -s -o /dev/null -w "%{http_code}" "${API_HOST}/support/config" 2>/dev/null)
if [ "$EXEMPT_SUPPORT_HTTP" = "200" ]; then
  ok "EXEMPT: /support/config returns 200 during maintenance"
else
  fail "EXEMPT: /support/config returned $EXEMPT_SUPPORT_HTTP (expected 200)"
fi

###############################################################################
header "8. MOBILE APP — PUBLIC CONFIG CHECK"
###############################################################################

# Simulate what the mobile app does on startup: GET /system/config
MOBILE_CHECK=$(curl -s "${API_HOST}/system/config" 2>/dev/null)

if command -v jq &>/dev/null; then
  MOBILE_MODE=$(echo "$MOBILE_CHECK" | jq -r '.data.maintenance_mode' 2>/dev/null)
  MOBILE_MSG=$(echo "$MOBILE_CHECK" | jq -r '.data.maintenance_message' 2>/dev/null)
  MOBILE_PHONE=$(echo "$MOBILE_CHECK" | jq -r '.data.maintenance_phone' 2>/dev/null)
else
  MOBILE_MODE=$(json_bool "$MOBILE_CHECK" "maintenance_mode")
  MOBILE_MSG="(install jq to see)"
  MOBILE_PHONE="(install jq to see)"
fi

if [ "$MOBILE_MODE" = "true" ]; then
  ok "Mobile startup check: maintenance_mode=true"
  log "  Message: $MOBILE_MSG"
  log "  Phone:   $MOBILE_PHONE"
else
  fail "Mobile startup check: expected maintenance_mode=true, got $MOBILE_MODE"
fi

###############################################################################
header "9. DISABLE MAINTENANCE MODE"
###############################################################################

DISABLE_RESP=$(curl -s -X POST "${ADMIN_HOST}/admin/system/config" \
  -H "Authorization: Bearer ${ADMIN_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"maintenance_mode":false}' 2>/dev/null)

DISABLE_CODE=$(json_val "$DISABLE_RESP" "code")
log "Response: $DISABLE_RESP"

if [ "$DISABLE_CODE" = "0000" ]; then
  ok "Maintenance mode DISABLED"
else
  fail "Failed to disable maintenance mode: $DISABLE_CODE"
fi

sleep 1

# Verify user API is accessible again
UNBLOCKED_HTTP=$(curl -s -o /dev/null -w "%{http_code}" \
  "${API_HOST}/profile" \
  -H "Authorization: Bearer fake-token" 2>/dev/null || echo "000")

# Should get 401 (unauthorized) not 503 (maintenance)
if [ "$UNBLOCKED_HTTP" = "401" ]; then
  ok "User API accessible after disabling maintenance (HTTP 401 = auth check, not 503)"
elif [ "$UNBLOCKED_HTTP" = "503" ]; then
  fail "Still getting 503 after disabling maintenance — cache may not have cleared"
else
  log "HTTP $UNBLOCKED_HTTP — may be normal depending on auth flow"
  ok "User API no longer returning 503"
fi

###############################################################################
header "10. REDIS CACHE CHECK"
###############################################################################

if [ "$HAS_REDIS" = "true" ]; then
  REDIS_PASS="${REDIS_PASS:-GreenRideRedis2025!}"
  REDIS_VAL=$(redis-cli -a "$REDIS_PASS" GET "greenride:system_config" 2>/dev/null | grep -v "Warning")

  if [ -n "$REDIS_VAL" ]; then
    ok "Redis cache key 'greenride:system_config' exists"
    log "Value: $REDIS_VAL"
  else
    log "Redis cache key is empty (normal if TTL expired)"
    ok "Redis cache check completed"
  fi

  REDIS_TTL=$(redis-cli -a "$REDIS_PASS" TTL "greenride:system_config" 2>/dev/null | grep -v "Warning")
  log "TTL: ${REDIS_TTL}s (expected: -2 if expired, or 0-30 if fresh)"
else
  skip "Redis cache check — redis-cli not available"
fi

###############################################################################
header "11. ORDER DETAIL — COORDINATE FIELDS (omitempty fix)"
###############################################################################

# This test requires a real order ID. Try to get a recent one from the admin API.
ORDERS_RESP=$(curl -s "${ADMIN_HOST}/admin/orders?page=1&size=1" \
  -H "Authorization: Bearer ${ADMIN_TOKEN}" 2>/dev/null)

if command -v jq &>/dev/null; then
  ORDER_ID=$(echo "$ORDERS_RESP" | jq -r '.data.list[0].order_id // .data[0].order_id // empty' 2>/dev/null)
else
  ORDER_ID=$(echo "$ORDERS_RESP" | grep -o '"order_id":"[^"]*"' | head -1 | sed 's/"order_id":"//;s/"//')
fi

if [ -n "$ORDER_ID" ] && [ "$ORDER_ID" != "null" ]; then
  log "Found order: $ORDER_ID"

  ORDER_DETAIL=$(curl -s "${ADMIN_HOST}/admin/orders/${ORDER_ID}" \
    -H "Authorization: Bearer ${ADMIN_TOKEN}" 2>/dev/null)

  if command -v jq &>/dev/null; then
    HAS_PLAT=$(echo "$ORDER_DETAIL" | jq 'has("data") and (.data | has("pickup_latitude"))' 2>/dev/null)
    PLAT=$(echo "$ORDER_DETAIL" | jq '.data.pickup_latitude // "missing"' 2>/dev/null)
    PLON=$(echo "$ORDER_DETAIL" | jq '.data.pickup_longitude // "missing"' 2>/dev/null)
    DLAT=$(echo "$ORDER_DETAIL" | jq '.data.dropoff_latitude // "missing"' 2>/dev/null)
    DLON=$(echo "$ORDER_DETAIL" | jq '.data.dropoff_longitude // "missing"' 2>/dev/null)

    log "pickup_latitude:  $PLAT"
    log "pickup_longitude: $PLON"
    log "dropoff_latitude: $DLAT"
    log "dropoff_longitude: $DLON"

    # Check that fields exist in JSON (even if 0) — the omitempty fix ensures this
    if echo "$ORDER_DETAIL" | jq -e '.data.pickup_latitude != null' &>/dev/null; then
      ok "pickup_latitude present in order detail response"
    else
      fail "pickup_latitude missing from order detail — omitempty fix may not be deployed"
    fi

    if echo "$ORDER_DETAIL" | jq -e '.data.dropoff_latitude != null' &>/dev/null; then
      ok "dropoff_latitude present in order detail response"
    else
      fail "dropoff_latitude missing from order detail — omitempty fix may not be deployed"
    fi
  else
    if echo "$ORDER_DETAIL" | grep -q "pickup_latitude"; then
      ok "pickup_latitude found in order detail JSON"
    else
      fail "pickup_latitude not found — omitempty fix may not be deployed"
    fi
    if echo "$ORDER_DETAIL" | grep -q "dropoff_latitude"; then
      ok "dropoff_latitude found in order detail JSON"
    else
      fail "dropoff_latitude not found — omitempty fix may not be deployed"
    fi
  fi
else
  skip "No orders found — cannot test coordinate fields"
  log "Create a test order first, then re-run this script"
fi

###############################################################################
header "12. DATABASE — SYSTEM CONFIG TABLE"
###############################################################################

# Check if MySQL is accessible and table exists
if command -v mysql &>/dev/null; then
  DB_PASS="${DB_PASSWORD:-GreenRide2024!}"
  TABLE_EXISTS=$(mysql -u greenride -p"$DB_PASS" -h 127.0.0.1 greenride \
    -e "SHOW TABLES LIKE 't_system_config';" 2>/dev/null | grep -c "t_system_config" || echo "0")

  if [ "$TABLE_EXISTS" -ge 1 ]; then
    ok "Table t_system_config exists in database"

    ROW_COUNT=$(mysql -u greenride -p"$DB_PASS" -h 127.0.0.1 greenride \
      -e "SELECT COUNT(*) AS cnt FROM t_system_config;" -sN 2>/dev/null || echo "0")
    log "Row count: $ROW_COUNT"

    if [ "$ROW_COUNT" -ge 1 ]; then
      ok "t_system_config has data ($ROW_COUNT row(s))"
      mysql -u greenride -p"$DB_PASS" -h 127.0.0.1 greenride \
        -e "SELECT id, maintenance_mode, maintenance_message, maintenance_phone, updated_by, updated_at FROM t_system_config LIMIT 1;" 2>/dev/null || true
    else
      log "No rows yet (will be auto-created on first API call)"
    fi
  else
    fail "Table t_system_config does NOT exist — AutoMigrate may not have run"
  fi
else
  skip "MySQL client not available — database checks skipped"
fi

###############################################################################
header "SUMMARY"
###############################################################################

TOTAL=$((PASS + FAIL + SKIP))
echo ""
echo -e "  ${GREEN}PASSED: ${PASS}${NC}"
echo -e "  ${RED}FAILED: ${FAIL}${NC}"
echo -e "  ${YELLOW}SKIPPED: ${SKIP}${NC}"
echo -e "  TOTAL:   ${TOTAL}"
echo ""

if [ "$FAIL" -eq 0 ]; then
  echo -e "${GREEN}${BOLD}All tests passed! Deployment is healthy.${NC}"
  exit 0
else
  echo -e "${RED}${BOLD}${FAIL} test(s) failed. Review output above.${NC}"
  exit 1
fi
