# 🔧 nginx Configuration Setup Guide

> **Purpose:** Configure nginx to handle CORS for GreenRide APIs  
> **Location:** Remote server (cloudshell)  
> **Files:** `nginx-config-complete.conf` (provided in this repo)

---

## 📋 **What This Configuration Does**

1. ✅ **Handles CORS** - Adds CORS headers to all API responses
2. ✅ **Proxies Requests** - Routes requests to backend services (ports 8610 & 8611)
3. ✅ **Handles OPTIONS** - Properly handles preflight requests
4. ✅ **Production Ready** - Includes timeouts, buffering, and error handling
5. ✅ **Dual API Support** - Configures both Mobile API and Admin API

---

## 🚀 **Quick Setup (Copy/Paste)**

### **Step 1: Backup Existing nginx Config**

```bash
# On cloudshell/server
sudo cp /etc/nginx/sites-available/default /etc/nginx/sites-available/default.backup
# OR if you have a custom config:
sudo cp /etc/nginx/sites-available/greenride /etc/nginx/sites-available/greenride.backup
```

### **Step 2: Copy Configuration File**

**Option A: Upload the file from your local machine**

```bash
# From your local machine (Windows PowerShell)
scp nginx-config-complete.conf user@18.143.118.157:/tmp/nginx-config-complete.conf
```

**Option B: Create file directly on server**

```bash
# On cloudshell/server
sudo nano /etc/nginx/sites-available/greenride-api
# Then copy/paste the entire content from nginx-config-complete.conf
```

### **Step 3: Enable the Configuration**

```bash
# Create symlink to enable the site
sudo ln -s /etc/nginx/sites-available/greenride-api /etc/nginx/sites-enabled/

# OR if replacing default:
sudo rm /etc/nginx/sites-enabled/default
sudo ln -s /etc/nginx/sites-available/greenride-api /etc/nginx/sites-enabled/greenride-api
```

### **Step 4: Test Configuration**

```bash
# Test nginx configuration syntax
sudo nginx -t

# Expected output:
# nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
# nginx: configuration file /etc/nginx/nginx.conf test is successful
```

### **Step 5: Reload nginx**

```bash
# Reload nginx (graceful restart)
sudo systemctl reload nginx

# OR restart nginx
sudo systemctl restart nginx

# Check status
sudo systemctl status nginx
```

---

## 📝 **Configuration Details**

### **Mobile API (Port 8610)**
- **Server Name:** `api.greenrideafrica.com` or `18.143.118.157`
- **Backend:** `localhost:8610`
- **Access Log:** `/var/log/nginx/greenride-mobile-api-access.log`
- **Error Log:** `/var/log/nginx/greenride-mobile-api-error.log`

### **Admin API (Port 8611)**
- **Server Name:** `admin-api.greenrideafrica.com`
- **Backend:** `localhost:8611`
- **Access Log:** `/var/log/nginx/greenride-admin-api-access.log`
- **Error Log:** `/var/log/nginx/greenride-admin-api-error.log`

### **CORS Headers Applied:**
```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS, PATCH
Access-Control-Allow-Headers: Content-Type, Authorization, Accept, Origin, X-Requested-With, Cache-Control, X-CSRF-Token
Access-Control-Allow-Credentials: true
Access-Control-Max-Age: 3600
```

---

## 🔍 **Verification**

### **Test CORS Headers**

```bash
# Test Mobile API
curl -I -X OPTIONS http://18.143.118.157/health \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: GET"

# Should return:
# HTTP/1.1 204 No Content
# Access-Control-Allow-Origin: *
# Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS, PATCH
# ...

# Test Admin API
curl -I -X OPTIONS http://18.143.118.157:8611/health \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: GET"
```

### **Test from Browser**

1. Open browser console (F12)
2. Navigate to `http://localhost:3000`
3. Check Network tab - API calls should succeed
4. No CORS errors should appear

---

## ⚙️ **Customization Options**

### **1. Restrict CORS Origins (More Secure)**

Replace `*` with specific origins:

```nginx
# Instead of:
add_header Access-Control-Allow-Origin "*" always;

# Use:
set $cors_origin "";
if ($http_origin ~* "^https?://(localhost|127\.0\.0\.1|admin\.greenrideafrica\.com)") {
    set $cors_origin $http_origin;
}
add_header Access-Control-Allow-Origin $cors_origin always;
```

### **2. Enable Caching for Static Endpoints**

Already included in the config for `/dashboard` and `/analytics` routes.

### **3. Add Rate Limiting**

```nginx
# Add to http block in nginx.conf
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;

# Then in server block:
limit_req zone=api_limit burst=20 nodelay;
```

### **4. Enable Compression**

```nginx
gzip on;
gzip_types text/plain text/css application/json application/javascript text/xml application/xml;
gzip_min_length 1000;
```

---

## 🐛 **Troubleshooting**

### **Problem: nginx won't start**

```bash
# Check syntax
sudo nginx -t

# Check logs
sudo tail -f /var/log/nginx/error.log

# Common issues:
# - Port already in use
# - Invalid syntax
# - Missing upstream backend
```

### **Problem: CORS still not working**

1. **Check headers are being sent:**
   ```bash
   curl -I http://18.143.118.157/health
   ```

2. **Check nginx is actually proxying:**
   ```bash
   sudo tail -f /var/log/nginx/greenride-admin-api-access.log
   ```

3. **Verify backend is running:**
   ```bash
   curl http://localhost:8611/health
   ```

### **Problem: 502 Bad Gateway**

- Backend service is not running
- Check backend logs
- Verify upstream server in nginx config matches actual backend port

### **Problem: 504 Gateway Timeout**

- Increase timeout values in nginx config:
  ```nginx
  proxy_connect_timeout 120s;
  proxy_send_timeout 120s;
  proxy_read_timeout 120s;
  ```

---

## 📊 **Current vs. Recommended Setup**

| Aspect | Current | With nginx |
|--------|---------|-----------|
| **CORS** | ❌ Not handled | ✅ Handled by nginx |
| **Direct Access** | ✅ `18.143.118.157:8611` | ✅ `18.143.118.157` (port 80) |
| **Security** | ⚠️ Direct backend access | ✅ nginx as reverse proxy |
| **Caching** | ❌ No caching | ✅ Optional caching |
| **SSL/HTTPS** | ❌ HTTP only | ✅ Can add SSL easily |

---

## 🔐 **Security Recommendations**

1. **Restrict CORS Origins** - Don't use `*` in production
2. **Enable Rate Limiting** - Prevent abuse
3. **Add SSL/HTTPS** - Encrypt traffic
4. **Firewall Rules** - Only allow nginx port (80/443) from outside
5. **Backend Access** - Only allow localhost access to backend ports (8610, 8611)

---

## ✅ **After Setup**

Once nginx is configured:

1. **Update Frontend `.env.local`:**
   ```bash
   # Change from:
   NEXT_PUBLIC_API_URL=http://18.143.118.157:8611
   
   # To (if using domain):
   NEXT_PUBLIC_API_URL=http://admin-api.greenrideafrica.com
   
   # OR (if still using IP):
   NEXT_PUBLIC_API_URL=http://18.143.118.157
   ```

2. **Restart Next.js:**
   ```bash
   npm run dev
   ```

3. **Test:** CORS errors should be gone! ✅

---

## 📝 **Summary**

**What to do:**
1. Copy `nginx-config-complete.conf` to server
2. Place in `/etc/nginx/sites-available/greenride-api`
3. Enable with symlink
4. Test: `sudo nginx -t`
5. Reload: `sudo systemctl reload nginx`
6. Verify CORS headers are present

**Result:**
- ✅ CORS errors fixed
- ✅ Frontend can connect to remote backend
- ✅ Production-ready configuration

---

**Need help?** Check nginx logs: `/var/log/nginx/error.log`
