#!/bin/bash
###############################################################################
# GreenRide — Manual Deployment Script (run on EC2)
#
# Deploys backend API + admin frontend from the latest main branch.
# Equivalent to what the GitHub Actions CI/CD pipeline does.
#
# Usage:
#   chmod +x deploy.sh
#   ./deploy.sh              # deploy both
#   ./deploy.sh backend      # backend only
#   ./deploy.sh frontend     # frontend only
#
# Paths (matches CI/CD):
#   Repo clone:   /home/ubuntu/greenride-admin-v.2
#   Backend dir:  /home/ubuntu/greenride-api
###############################################################################

set -euo pipefail

# ─── Config ──────────────────────────────────────────────────────────────────
REPO_DIR="/home/ubuntu/greenride-admin-v.2"
BACKEND_DIR="/home/ubuntu/greenride-api"
GIT_BRANCH="main"

# Container names (must match existing setup)
BACKEND_CONTAINER="greenride-api"
FRONTEND_CONTAINER="greenride-admin-v2"

# Ports
API_PORT=8610
ADMIN_PORT=8611
FRONTEND_PORT=3601

# Go version for building backend
GO_VERSION="1.21"

# ─── Colors ──────────────────────────────────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

log()  { echo -e "${CYAN}[$(date '+%H:%M:%S')]${NC} $*"; }
ok()   { echo -e "${GREEN}[OK]${NC}    $*"; }
warn() { echo -e "${YELLOW}[WARN]${NC}  $*"; }
err()  { echo -e "${RED}[ERROR]${NC} $*"; }
step() { echo -e "\n${BOLD}── $* ──${NC}"; }

# ─── Parse args ──────────────────────────────────────────────────────────────
DEPLOY_TARGET="${1:-all}"

if [[ "$DEPLOY_TARGET" != "all" && "$DEPLOY_TARGET" != "backend" && "$DEPLOY_TARGET" != "frontend" ]]; then
  echo "Usage: $0 [all|backend|frontend]"
  exit 1
fi

echo -e "${BOLD}"
echo "╔═══════════════════════════════════════════════╗"
echo "║        GreenRide Deployment Script            ║"
echo "║        Target: ${DEPLOY_TARGET}                           "
echo "║        $(date)     ║"
echo "╚═══════════════════════════════════════════════╝"
echo -e "${NC}"

# ─── Pre-flight ──────────────────────────────────────────────────────────────
step "Pre-flight checks"

if [ ! -d "$REPO_DIR" ]; then
  err "Repo directory not found: $REPO_DIR"
  exit 1
fi
ok "Repo directory exists"

if ! command -v docker &>/dev/null; then
  err "Docker is not installed"
  exit 1
fi
ok "Docker available"

if ! command -v git &>/dev/null; then
  err "Git is not installed"
  exit 1
fi
ok "Git available"

# Check disk space (warn if < 2GB free)
FREE_MB=$(df -m / | tail -1 | awk '{print $4}')
if [ "$FREE_MB" -lt 2048 ]; then
  warn "Low disk space: ${FREE_MB}MB free. Consider pruning Docker images."
else
  ok "Disk space: ${FREE_MB}MB free"
fi

# ─── Pull latest code ───────────────────────────────────────────────────────
step "Pulling latest code from ${GIT_BRANCH}"

cd "$REPO_DIR"
BEFORE_SHA=$(git rev-parse HEAD 2>/dev/null || echo "unknown")
git fetch origin "$GIT_BRANCH"
git reset --hard "origin/$GIT_BRANCH"
AFTER_SHA=$(git rev-parse HEAD 2>/dev/null || echo "unknown")

if [ "$BEFORE_SHA" = "$AFTER_SHA" ]; then
  log "Already up to date ($AFTER_SHA)"
else
  ok "Updated: $BEFORE_SHA → $AFTER_SHA"
  log "Recent commits:"
  git log --oneline -5
fi
echo ""

# ═════════════════════════════════════════════════════════════════════════════
#  BACKEND DEPLOYMENT
# ═════════════════════════════════════════════════════════════════════════════
deploy_backend() {
  step "Building backend Go binary"

  mkdir -p "$BACKEND_DIR"
  cd "$REPO_DIR/backend/greenride-api-clean"

  # Build strategy: native Go > Docker-based build
  # CGO_ENABLED=0 is REQUIRED — the binary runs on Alpine (musl), not glibc
  if command -v go &>/dev/null; then
    GO_INSTALLED=$(go version | grep -oP 'go\d+\.\d+' || echo "")
    log "Using native Go ($GO_INSTALLED) with static linking (CGO_ENABLED=0)"
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o greenride-api-linux ./main
  else
    log "Go not installed — building inside Docker"
    docker run --rm \
      -v "$(pwd)":/src \
      -w /src \
      -e CGO_ENABLED=0 \
      -e GOOS=linux \
      -e GOARCH=amd64 \
      golang:${GO_VERSION}-alpine \
      go build -o greenride-api-linux ./main
  fi

  if [ ! -f "greenride-api-linux" ]; then
    err "Binary not found after build. Build failed."
    return 1
  fi
  BINARY_SIZE=$(du -h greenride-api-linux | cut -f1)
  ok "Binary built: greenride-api-linux ($BINARY_SIZE)"

  step "Deploying backend"

  # Copy binary to deployment dir
  cp greenride-api-linux "$BACKEND_DIR/"
  chmod +x "$BACKEND_DIR/greenride-api-linux"

  # Copy config and locale files
  cp -r internal/locales "$BACKEND_DIR/internal/" 2>/dev/null || {
    mkdir -p "$BACKEND_DIR/internal/locales"
    cp -r internal/locales/* "$BACKEND_DIR/internal/locales/"
  }
  cp config.yaml "$BACKEND_DIR/" 2>/dev/null || true

  # prod.yaml: only copy if missing (never overwrite server secrets)
  if [ -d "$BACKEND_DIR/prod.yaml" ]; then
    rm -rf "$BACKEND_DIR/prod.yaml"
  fi
  if [ ! -f "$BACKEND_DIR/prod.yaml" ]; then
    cp prod.yaml "$BACKEND_DIR/" 2>/dev/null || warn "No prod.yaml in repo"
  else
    log "prod.yaml already exists on server — not overwriting"
  fi

  # Copy Dockerfile
  cp Dockerfile "$BACKEND_DIR/"

  step "Restarting backend container"

  cd "$BACKEND_DIR"

  # Stop existing container
  docker stop "$BACKEND_CONTAINER" 2>/dev/null && log "Stopped $BACKEND_CONTAINER" || true
  docker rm "$BACKEND_CONTAINER" 2>/dev/null && log "Removed $BACKEND_CONTAINER" || true

  # Build new image
  docker build -t "$BACKEND_CONTAINER" .
  ok "Docker image built"

  # Find the Docker network (if exists)
  NETWORK_FLAG=""
  if docker network inspect greenride-network &>/dev/null; then
    NETWORK_FLAG="--network greenride-network"
    log "Joining greenride-network"
  fi

  # Start container
  docker run -d \
    --name "$BACKEND_CONTAINER" \
    --restart unless-stopped \
    $NETWORK_FLAG \
    -p ${API_PORT}:${API_PORT} \
    -p ${ADMIN_PORT}:${ADMIN_PORT} \
    -p 8612:8612 \
    -v "$BACKEND_DIR/config.yaml:/app/config.yaml" \
    -v "$BACKEND_DIR/prod.yaml:/app/prod.yaml" \
    -v "$BACKEND_DIR/logs:/app/logs" \
    "$BACKEND_CONTAINER"

  ok "Backend container started"

  # Health check with retry
  step "Backend health check"
  for i in 1 2 3 4 5 6; do
    sleep 2
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:${API_PORT}/health" 2>/dev/null || echo "000")
    if [ "$HTTP_CODE" = "200" ]; then
      ok "API (port $API_PORT) is healthy"
      break
    fi
    if [ "$i" -eq 6 ]; then
      err "API health check failed after 12s (HTTP $HTTP_CODE)"
      log "Check logs: docker logs $BACKEND_CONTAINER --tail 50"
      return 1
    fi
    log "Waiting... (attempt $i/6, HTTP $HTTP_CODE)"
  done

  ADMIN_HTTP=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:${ADMIN_PORT}/health" 2>/dev/null || echo "000")
  if [ "$ADMIN_HTTP" = "200" ]; then
    ok "Admin API (port $ADMIN_PORT) is healthy"
  else
    warn "Admin API health check: HTTP $ADMIN_HTTP"
  fi

  # Clean up build artifact from repo dir
  rm -f "$REPO_DIR/backend/greenride-api-clean/greenride-api-linux"
}

# ═════════════════════════════════════════════════════════════════════════════
#  FRONTEND DEPLOYMENT
# ═════════════════════════════════════════════════════════════════════════════
deploy_frontend() {
  step "Building frontend"

  cd "$REPO_DIR"

  # Stop existing container
  docker stop "$FRONTEND_CONTAINER" 2>/dev/null && log "Stopped $FRONTEND_CONTAINER" || true
  docker rm "$FRONTEND_CONTAINER" 2>/dev/null && log "Removed $FRONTEND_CONTAINER" || true

  # Build new image
  docker build -t "$FRONTEND_CONTAINER" .
  ok "Frontend Docker image built"

  step "Starting frontend container"

  docker run -d \
    --name "$FRONTEND_CONTAINER" \
    --restart unless-stopped \
    -p ${FRONTEND_PORT}:3000 \
    "$FRONTEND_CONTAINER"

  ok "Frontend container started on port $FRONTEND_PORT"

  # Health check
  step "Frontend health check"
  for i in 1 2 3 4 5; do
    sleep 3
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:${FRONTEND_PORT}" 2>/dev/null || echo "000")
    if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "302" ] || [ "$HTTP_CODE" = "307" ]; then
      ok "Frontend is healthy (HTTP $HTTP_CODE)"
      break
    fi
    if [ "$i" -eq 5 ]; then
      warn "Frontend health check: HTTP $HTTP_CODE (may still be starting)"
      log "Check logs: docker logs $FRONTEND_CONTAINER --tail 30"
    fi
    log "Waiting... (attempt $i/5, HTTP $HTTP_CODE)"
  done
}

# ═════════════════════════════════════════════════════════════════════════════
#  EXECUTE
# ═════════════════════════════════════════════════════════════════════════════

if [[ "$DEPLOY_TARGET" == "all" || "$DEPLOY_TARGET" == "backend" ]]; then
  deploy_backend
fi

if [[ "$DEPLOY_TARGET" == "all" || "$DEPLOY_TARGET" == "frontend" ]]; then
  deploy_frontend
fi

# ─── Cleanup ─────────────────────────────────────────────────────────────────
step "Cleanup"
PRUNED=$(docker image prune -f 2>/dev/null | tail -1)
log "$PRUNED"

# ─── Summary ─────────────────────────────────────────────────────────────────
echo ""
echo -e "${BOLD}═══════════════════════════════════════════════${NC}"
echo -e "${GREEN}${BOLD}  Deployment complete!${NC}"
echo -e "${BOLD}═══════════════════════════════════════════════${NC}"
echo ""
echo "  Commit:    $(git -C "$REPO_DIR" log --oneline -1)"
echo "  Deployed:  $DEPLOY_TARGET"
echo "  Time:      $(date)"
echo ""

if [[ "$DEPLOY_TARGET" == "all" || "$DEPLOY_TARGET" == "backend" ]]; then
  echo "  Backend:   http://localhost:$API_PORT/health"
  echo "  Admin API: http://localhost:$ADMIN_PORT/health"
fi
if [[ "$DEPLOY_TARGET" == "all" || "$DEPLOY_TARGET" == "frontend" ]]; then
  echo "  Frontend:  http://localhost:$FRONTEND_PORT"
fi

echo ""
echo "  Running containers:"
docker ps --format "  {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep -i greenride || true
echo ""
