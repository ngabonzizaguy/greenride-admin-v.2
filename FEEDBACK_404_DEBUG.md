# 🔍 Feedback Endpoints 404 Debug Guide

## 📋 Issue
Console shows 404 errors for feedback endpoints:
- `POST /admin/api/feedback/search` → 404 Not Found
- `GET /admin/api/feedback/stats` → 404 Not Found

## ✅ What We Know
1. **Frontend calls are correct:**
   - `POST /feedback/search` (line 1056 in `api-client.ts`)
   - `GET /feedback/stats` (line 1168 in `api-client.ts`)

2. **Backend endpoints are registered:**
   - `POST /feedback/search` → `SearchFeedback` (line 148 in `admin.go`)
   - `GET /feedback/stats` → `GetFeedbackStats` (line 151 in `admin.go`)

3. **Nginx rewrite should work:**
   - `/admin/api/feedback/stats` → rewrites to → `/feedback/stats`
   - `/admin/api/feedback/search` → rewrites to → `/feedback/search`

4. **Other endpoints work:**
   - `/admin/api/dashboard/stats` → ✅ Works (200 OK)

## 🔧 Debug Steps (Run on Server)

### **Step 1: Check if Backend is Running**

```bash
# Check if backend process is running on port 8611
sudo netstat -tlnp | grep 8611

# OR
sudo ss -tlnp | grep 8611

# Expected output: Something like:
# tcp  0  0  0.0.0.0:8611  0.0.0.0:*  LISTEN  12345/backend
```

### **Step 2: Test Backend Directly (Bypass Nginx)**

```bash
# Test feedback/stats endpoint directly
curl -X GET http://localhost:8611/feedback/stats \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -v

# Test feedback/search endpoint directly
curl -X POST http://localhost:8611/feedback/search \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"page":1,"limit":10}' \
  -v
```

**Expected results:**
- ✅ **200 OK** → Backend is working, issue is with nginx
- ❌ **404 Not Found** → Backend endpoints aren't registered
- ❌ **Connection refused** → Backend isn't running

### **Step 3: Test Through Nginx**

```bash
# Test through nginx (same as frontend)
curl -X GET http://localhost/admin/api/feedback/stats \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -v

# Test feedback/search through nginx
curl -X POST http://localhost/admin/api/feedback/search \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"page":1,"limit":10}' \
  -v
```

**Compare with dashboard/stats (which works):**
```bash
# This should work
curl -X GET http://localhost/admin/api/dashboard/stats \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -v
```

### **Step 4: Check Nginx Rewrite Logs**

```bash
# Check nginx error logs
docker exec nginx tail -n 50 /var/log/nginx/error.log | grep feedback

# Check nginx access logs
docker exec nginx tail -n 50 /var/log/nginx/access.log | grep feedback
```

### **Step 5: Verify Backend Endpoint Registration**

```bash
# Check backend logs for route registration
# Look for lines like:
# [GIN-debug] POST   /feedback/search
# [GIN-debug] GET    /feedback/stats

# If you have access to backend logs:
tail -n 100 /path/to/backend.log | grep feedback
```

### **Step 6: Test Health Endpoint**

```bash
# Verify backend is responding at all
curl http://localhost:8611/health

# Expected: {"service":"admin","status":"ok"}
```

## 🎯 Most Likely Causes

### **Cause 1: Backend Not Running**
**Solution:** Start/restart backend
```bash
# Find backend process
ps aux | grep greenride

# Kill old process if needed
kill -9 <PID>

# Start backend (adjust path as needed)
cd /path/to/backend
./backend-admin
```

### **Cause 2: Backend Routes Not Registered**
**Solution:** Check if backend code was recompiled after adding feedback endpoints
```bash
# Restart backend to re-register routes
systemctl restart greenride-admin
# OR
docker restart greenride-admin-backend
```

### **Cause 3: Nginx Config Not Applied**
**Solution:** Reload nginx
```bash
# Test config first
docker exec nginx nginx -t

# If OK, reload
docker exec nginx nginx -s reload
```

## 📝 Quick Fix Checklist

1. [ ] Backend is running on port 8611
2. [ ] Backend `/health` endpoint responds
3. [ ] Backend `/feedback/stats` responds when called directly
4. [ ] Nginx config is correct and applied
5. [ ] Nginx rewrite rule works for `/admin/api/*` → `/*`
6. [ ] Token is valid and included in Authorization header

## 🔍 What to Report Back

After running the debug steps, report:
1. **Backend direct test results** (Step 2)
2. **Nginx test results** (Step 3)
3. **Backend logs** showing route registration
4. **Nginx error/access logs** (if any errors)

This will help identify the exact issue.
