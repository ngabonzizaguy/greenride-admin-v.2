# 🔧 CORS Error Fix Guide

> **Issue:** CORS errors when connecting to remote backend (`18.143.118.157:8611`)  
> **Error:** `Access to fetch at 'http://18.143.118.157:8611/dashboard/stats' from origin 'http://localhost:3000' has been blocked by CORS policy`

---

## 🔍 **Root Cause**

The remote backend at `18.143.118.157:8611` is running in **production mode** (`env: prod`), which means:

1. ✅ CORS middleware is **disabled** (only enabled in dev mode)
2. ✅ Backend expects **nginx to handle CORS** (not the backend itself)
3. ✅ You're connecting **directly** to the backend (not through nginx)
4. ❌ Result: **CORS errors** because backend doesn't send CORS headers

**Code Evidence:**
```go
// In admin.go (line 56-67)
// Local Development CORS
if cfg.Env != protocol.EnvProduction {
    router.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        // ... CORS headers
    })
}
```

**Problem:** Remote backend has `env: prod`, so CORS middleware is **not active**.

---

## ✅ **Solutions**

### **Solution 1: Use Local Backend (Recommended for Development)**

**Why:** No CORS issues, faster development workflow

**Steps:**

1. **Update `.env.local`:**
   ```bash
   # Change from remote to local
   NEXT_PUBLIC_API_URL=http://localhost:8611
   NEXT_PUBLIC_DEMO_MODE=false
   ```

2. **Start Local Backend:**
   ```bash
   cd backend/greenride-api-clean
   go run main/main.go
   ```

3. **Restart Next.js:**
   ```bash
   # Stop current server (Ctrl+C)
   npm run dev
   ```

**Benefits:**
- ✅ No CORS issues (same origin: `localhost:3000` → `localhost:8611`)
- ✅ Faster development (no network latency)
- ✅ Better debugging (see backend logs directly)
- ✅ CORS enabled in dev mode automatically

---

### **Solution 2: Enable CORS on Remote Backend (If You Need Remote)**

**Requires:** Backend configuration change on remote server

**Option A: Change Remote Backend to Dev Mode**

**On Remote Server:**
1. Edit `config.yaml` or `dev.yaml`:
   ```yaml
   env: dev  # Change from 'prod' to 'dev'
   ```

2. Restart backend:
   ```bash
   # On remote server
   sudo systemctl restart greenride-api
   # OR
   # Restart backend process
   ```

**Option B: Enable CORS for Production (Not Recommended)**

**Modify `backend/greenride-api-clean/internal/handlers/admin.go`:**

```go
// Change line 56 from:
if cfg.Env != protocol.EnvProduction {
    router.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })
}

// To (always enable CORS):
router.Use(func(c *gin.Context) {
    c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
    if c.Request.Method == "OPTIONS" {
        c.AbortWithStatus(204)
        return
    }
    c.Next()
})
```

**⚠️ Warning:** This allows all origins (`*`) - **not secure for production**!

**Better Option:** Allow specific origins:
```go
allowedOrigins := []string{
    "http://localhost:3000",
    "http://localhost:3001",
    "https://admin.greenrideafrica.com",
}

origin := c.Request.Header.Get("Origin")
if contains(allowedOrigins, origin) {
    c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
}
```

---

### **Solution 3: Use nginx Proxy (Production Setup)**

**Requires:** nginx configuration on remote server

**nginx Configuration:**
```nginx
server {
    listen 80;
    server_name admin-api.greenrideafrica.com;

    location / {
        proxy_pass http://localhost:8611;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        
        # CORS headers
        add_header Access-Control-Allow-Origin "*" always;
        add_header Access-Control-Allow-Methods "POST, GET, OPTIONS, PUT, DELETE" always;
        add_header Access-Control-Allow-Headers "Content-Type, Authorization, Accept, Origin" always;
        
        if ($request_method = OPTIONS) {
            return 204;
        }
    }
}
```

**Then connect to:**
```bash
NEXT_PUBLIC_API_URL=http://admin-api.greenrideafrica.com
```

---

## 🎯 **Recommended Solution for Development**

### **Use Local Backend (Solution 1)**

**Why:**
- ✅ No CORS issues (same origin)
- ✅ Faster development
- ✅ Better debugging
- ✅ No server configuration needed

**Quick Steps:**
1. Change `.env.local`: `NEXT_PUBLIC_API_URL=http://localhost:8611`
2. Start local backend: `cd backend/greenride-api-clean && go run main/main.go`
3. Restart Next.js: `npm run dev`

---

## 📊 **Current Setup Analysis**

| Aspect | Current | Recommended |
|--------|---------|-------------|
| **Frontend** | `localhost:3000` | `localhost:3000` |
| **Backend** | `18.143.118.157:8611` (remote) | `localhost:8611` (local) |
| **Environment** | Production (CORS disabled) | Dev (CORS enabled) |
| **CORS Status** | ❌ Blocked | ✅ Allowed |
| **For Development** | ❌ Not ideal | ✅ Perfect |

---

## ⚠️ **Why CORS Errors Happen**

**CORS (Cross-Origin Resource Sharing) Policy:**
- Browser enforces CORS for security
- `localhost:3000` (frontend) → `18.143.118.157:8611` (backend) = **Different origins**
- Backend must send `Access-Control-Allow-Origin` header
- Remote backend (production mode) doesn't send this header
- Browser blocks the request ❌

**Same Origin (No CORS):**
- `localhost:3000` (frontend) → `localhost:8611` (backend) = **Same origin**
- Browser doesn't enforce CORS ✅
- Requests work directly ✅

---

## 🔧 **Quick Fix (Recommended)**

### **Step 1: Update `.env.local`**

```bash
# Change from:
NEXT_PUBLIC_API_URL=http://18.143.118.157:8611

# To:
NEXT_PUBLIC_API_URL=http://localhost:8611
```

### **Step 2: Start Local Backend**

```bash
cd backend/greenride-api-clean

# Kill any existing backend process first
taskkill /F /IM main.exe 2>$null

# Start backend
go run main/main.go
```

### **Step 3: Restart Next.js**

```bash
# Stop current server (Ctrl+C)
npm run dev
```

### **Step 4: Verify**

1. Open browser: `http://localhost:3000`
2. Check console (F12) - should see:
   ```
   [API Client] API_BASE_URL: http://localhost:8611
   ```
3. Dashboard should load without CORS errors ✅

---

## ✅ **Expected Result After Fix**

**Before (CORS Error):**
```
❌ Access to fetch at 'http://18.143.118.157:8611/dashboard/stats' 
   from origin 'http://localhost:3000' has been blocked by CORS policy
❌ Failed to fetch dashboard stats
❌ Failed to fetch sidebar counts
❌ Dashboard shows cached/empty data
```

**After (Working):**
```
✅ API calls succeed
✅ Dashboard loads real data
✅ Sidebar counts load correctly
✅ No CORS errors in console
```

---

## 📝 **Summary**

**Problem:** Remote backend in production mode doesn't send CORS headers  
**Solution:** Use local backend for development (recommended)  
**Alternative:** Enable CORS on remote backend (requires server config)

**Quick Fix:**
1. Change `.env.local` to `localhost:8611`
2. Start local backend
3. Restart Next.js
4. Done! ✅

---

**This should fix your CORS errors!** 🚀
