# 🔴 Fix Redirect Issue

## **Problem**

You're getting **307 redirects** to the wrong domain:
- `http://localhost/admin/api/health` → `https://admin-dev.greenrideafrica.com/admin/api/health` ❌
- `http://18.143.118.157/admin/api/health` → `https://admin-dev.greenrideafrica.com/admin/api/health` ❌

**This means:**
1. ❌ Current nginx config has redirect rules (not updated yet)
2. ❌ Redirecting to wrong domain (`admin-dev` instead of `admin`)
3. ❌ Forcing HTTPS when we want HTTP to work

---

## ✅ **Solution: Check Current nginx Config**

**On cloudshell, check what's causing the redirect:**

```bash
# View current nginx config
cat ~/nginx.conf | grep -i "redirect\|return\|location.*admin"

# OR search for admin-dev
cat ~/nginx.conf | grep -i "admin-dev"

# OR check for HTTPS redirects
cat ~/nginx.conf | grep -i "https\|ssl"
```

**Look for lines like:**
- `return 301` or `return 307`
- `rewrite.*permanent`
- `location.*admin-dev`
- `server_name.*admin-dev`

---

## 🔧 **Fix Steps**

### **Step 1: Backup Current Config**

```bash
cp ~/nginx.conf ~/nginx.conf.old
```

### **Step 2: Replace with Fixed Config**

Use the updated `nginx-config-fixed.conf` (which I just updated to handle `admin.greenrideafrica.com` correctly).

### **Step 3: Test Config**

```bash
docker exec nginx nginx -t
```

### **Step 4: Reload nginx**

```bash
docker exec nginx nginx -s reload
```

### **Step 5: Test Again**

```bash
# Should return 200 OK (no redirect)
curl -I http://localhost/admin/api/health

# Should return 200 OK (no redirect)
curl -I http://18.143.118.157/admin/api/health
```

---

## ⚠️ **If Redirects Still Happen After Fix**

If you still see redirects after updating nginx, check:

### **1. Cloudflare Settings**

The redirect from `admin.greenrideafrica.com` HTTP → HTTPS is coming from **Cloudflare**, not nginx.

**To fix:**
1. Go to Cloudflare Dashboard
2. Select domain: `greenrideafrica.com`
3. Go to **SSL/TLS** → **Overview**
4. Set encryption mode to **"Flexible"** (allows HTTP) OR **"Full"** (requires HTTPS on backend)
5. Go to **Rules** → **Page Rules**
6. Check if there's a rule forcing HTTPS for `admin.greenrideafrica.com/*`
7. Remove or modify the rule

### **2. Check for Multiple nginx Configs**

```bash
# Check if there are other nginx config files
docker exec nginx find /etc/nginx -name "*.conf" -type f

# Check main config includes
docker exec nginx cat /etc/nginx/nginx.conf | grep include
```

---

## 📋 **Quick Diagnostic Commands**

```bash
# 1. Check current nginx config for redirects
cat ~/nginx.conf | grep -E "return|rewrite|admin-dev" -i

# 2. Check what nginx container is actually using
docker exec nginx cat /etc/nginx/nginx.conf | head -50

# 3. Check nginx error logs
docker logs nginx 2>&1 | tail -20

# 4. Test if backend is reachable directly
curl -I http://localhost:8611/health
```

---

## 🎯 **Expected Result After Fix**

**Before (BAD):**
```
HTTP/1.1 307 Temporary Redirect
Location: https://admin-dev.greenrideafrica.com/admin/api/health
```

**After (GOOD):**
```
HTTP/1.1 200 OK
Content-Type: application/json
{"status":"ok","service":"admin"}
```

---

## 🚨 **Action Required**

**First, check what's in your current nginx config:**

```bash
cat ~/nginx.conf | grep -i "admin-dev\|redirect\|return 30"
```

**Share the output** and I'll help you identify the exact redirect rule that needs to be removed!
