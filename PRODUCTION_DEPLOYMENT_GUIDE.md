# 🚀 Production Deployment Guide

> **Goal:** Deploy Admin Dashboard to Production and Connect to Production Backend  
> **From:** Dev (`http://18.143.118.157/admin/api`)  
> **To:** Prod (Production API URL)

---

## 📋 **What Needs to Change for Production**

### **1. Environment Variables**

#### **For Development (Current):**
```bash
# .env.local (on your local machine)
NEXT_PUBLIC_API_URL=http://18.143.118.157/admin/api
NEXT_PUBLIC_DEMO_MODE=false
NEXT_PUBLIC_GOOGLE_MAPS_KEY=AIzaSyBXaxfZrk9-qdmnuYY-3YMhNcnxyH_lj8Q
```

#### **For Production:**
```bash
# .env.production (when building for production)
NEXT_PUBLIC_API_URL=https://admin-api.greenrideafrica.com/admin/api
# OR if using IP:
# NEXT_PUBLIC_API_URL=https://18.143.118.157/admin/api
NEXT_PUBLIC_DEMO_MODE=false
NEXT_PUBLIC_GOOGLE_MAPS_KEY=AIzaSyBXaxfZrk9-qdmnuYY-3YMhNcnxyH_lj8Q
```

**Key Changes:**
- ✅ Change `http://` → `https://` (production should use HTTPS)
- ✅ Use production domain (e.g., `admin-api.greenrideafrica.com`) or production IP
- ✅ Keep `NEXT_PUBLIC_DEMO_MODE=false` (always use real API)

---

## 🔧 **Step-by-Step: Deploy to Production**

### **Step 1: Create `.env.production` File**

In your project root, create `.env.production`:

```bash
# Production Backend API
NEXT_PUBLIC_API_URL=https://admin-api.greenrideafrica.com/admin/api
NEXT_PUBLIC_DEMO_MODE=false

# Google Maps (production key if different)
NEXT_PUBLIC_GOOGLE_MAPS_KEY=YOUR_PROD_GOOGLE_MAPS_KEY
```

**Note:** Replace `https://admin-api.greenrideafrica.com` with your actual production API URL.

---

### **Step 2: Build Next.js for Production**

```bash
# Build the production version
npm run build

# This creates:
# - .next/ folder (production build)
# - Static files
# - Server files
```

---

### **Step 3: Deploy Admin Dashboard to Server**

**Option A: Deploy to Same Server (where nginx is)**

```bash
# On production server
# 1. Upload build files
scp -r .next ubuntu@18.143.118.157:/opt/admin-dashboard/

# 2. Upload package.json and node_modules (or install on server)
scp package.json ubuntu@18.143.118.157:/opt/admin-dashboard/

# 3. Upload .env.production
scp .env.production ubuntu@18.143.118.157:/opt/admin-dashboard/.env.production
```

**Option B: Deploy via Docker**

```dockerfile
# Dockerfile for Admin Dashboard
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM node:20-alpine AS runner
WORKDIR /app
ENV NODE_ENV production
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/public ./public
RUN npm ci --only=production
EXPOSE 3000
CMD ["npm", "start"]
```

Then run:
```bash
docker build -t greenride-admin-dashboard:latest .
docker run -d -p 3600:3000 \
  -e NEXT_PUBLIC_API_URL=https://admin-api.greenrideafrica.com/admin/api \
  -e NEXT_PUBLIC_DEMO_MODE=false \
  greenride-admin-dashboard:latest
```

**Option C: Deploy via PM2**

```bash
# On production server
cd /opt/admin-dashboard
npm install --production
npm run build
pm2 start npm --name "admin-dashboard" -- start
pm2 save
```

---

### **Step 4: Update nginx Configuration**

**On Production Server:** Update `/opt/nginx/conf/greenride/nginx.conf` to handle production domain:

```nginx
# Admin Domain - admin.greenrideafrica.com (Production)
server {
    listen 80;
    server_name admin.greenrideafrica.com;
    
    # Redirect HTTP to HTTPS (if SSL is configured)
    # return 301 https://$server_name$request_uri;
    
    # OR serve directly if no SSL yet
    # ... (same config as dev)
}

# If SSL is configured:
server {
    listen 443 ssl http2;
    server_name admin.greenrideafrica.com;
    
    ssl_certificate /etc/nginx/certs/admin.crt;
    ssl_certificate_key /etc/nginx/certs/admin.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    
    # ... (same proxy config as dev)
}
```

---

### **Step 5: Verify Production Backend**

**Ensure Production Backend is Running:**

```bash
# Test production backend
curl https://admin-api.greenrideafrica.com/admin/api/health

# Should return:
# {"service":"admin","status":"ok"}
```

---

### **Step 6: Update CORS Origins in nginx**

**In production nginx config**, update CORS allowlist:

```nginx
# CORS allowlist (production domains)
map $http_origin $cors_allow_origin {
    default "";
    "https://admin.greenrideafrica.com" $http_origin;
    "https://admin-dev.greenrideafrica.com" $http_origin;
    # Remove localhost origins for production
    # "http://localhost:3000" $http_origin;
}
```

---

## 📊 **Configuration Comparison**

| Setting | Development | Production |
|---------|-------------|------------|
| **API URL** | `http://18.143.118.157/admin/api` | `https://admin-api.greenrideafrica.com/admin/api` |
| **Protocol** | HTTP | HTTPS ✅ |
| **Domain** | IP address | Domain name ✅ |
| **Demo Mode** | `false` | `false` ✅ |
| **CORS Origins** | `localhost:3000` | `admin.greenrideafrica.com` ✅ |
| **Build Mode** | `npm run dev` | `npm run build` → `npm start` ✅ |

---

## ✅ **Verification Checklist**

After deploying to production:

- [ ] `.env.production` is created with production API URL
- [ ] `NEXT_PUBLIC_DEMO_MODE=false` is set
- [ ] Admin dashboard is built (`npm run build`)
- [ ] Admin dashboard is running on server (port 3600 or 3001)
- [ ] nginx config updated for production domain
- [ ] CORS origins updated (removed localhost)
- [ ] SSL certificates configured (if using HTTPS)
- [ ] Production backend is accessible
- [ ] Can login to admin dashboard
- [ ] Data loads correctly (not mock data)

---

## 🔍 **How to Verify It's Using Production Data**

### **1. Check Browser Console**

Open `https://admin.greenrideafrica.com` → Console (F12) → Look for:
```
[API Client] API_BASE_URL: https://admin-api.greenrideafrica.com/admin/api
[API Client] DEMO_MODE: false
```

### **2. Check Network Tab**

F12 → Network → Look for API calls to:
- `https://admin-api.greenrideafrica.com/admin/api/*`
- **NOT** `http://18.143.118.157/*` (that's dev)

### **3. Verify Data Source**

The data should match your **production database**, not dev database:
- User counts, ride counts, revenue should match production
- Test with a known production user/ride

---

## 🚨 **Important Notes**

### **1. Environment Variables**

**Next.js Environment Variable Priority:**
1. `.env.production` (used when `NODE_ENV=production`)
2. `.env.local` (always loaded, but `.env.production` takes precedence)
3. `.env` (fallback)

**For Production Build:**
- Create `.env.production` file
- Don't commit `.env.production` to git (add to `.gitignore`)
- Set environment variables on server or in deployment config

### **2. HTTPS/SSL**

**Production MUST use HTTPS:**
- Update API URL to `https://` (not `http://`)
- Configure SSL certificates in nginx
- Update CORS origins to `https://admin.greenrideafrica.com`

### **3. Build vs Dev**

**Development:**
```bash
npm run dev  # Uses .env.local, hot reload
```

**Production:**
```bash
npm run build  # Builds for production, uses .env.production
npm start      # Runs production server
```

---

## 🐛 **Troubleshooting**

### **Problem: Still seeing dev data**

**Check:**
1. Is `.env.production` correct?
2. Did you rebuild after changing `.env.production`? (`npm run build`)
3. Check browser console for `API_BASE_URL`
4. Clear browser cache

**Fix:**
```bash
# Rebuild with production env
rm -rf .next
npm run build
```

### **Problem: CORS errors in production**

**Check:**
1. Is nginx CORS config updated for production domain?
2. Is API URL using correct domain?
3. Check nginx logs: `docker logs nginx`

**Fix:**
Update nginx CORS origins to include production domain.

### **Problem: API calls failing**

**Check:**
1. Is production backend running?
2. Is production backend accessible?
3. Check backend logs

**Test:**
```bash
curl https://admin-api.greenrideafrica.com/admin/api/health
```

---

## 📝 **Quick Reference**

### **Files to Update for Production:**

1. **`.env.production`** (create this file)
   ```
   NEXT_PUBLIC_API_URL=https://admin-api.greenrideafrica.com/admin/api
   NEXT_PUBLIC_DEMO_MODE=false
   ```

2. **`/opt/nginx/conf/greenride/nginx.conf`** (on server)
   - Update CORS origins
   - Add SSL config (if needed)

3. **Build command:**
   ```bash
   npm run build
   npm start
   ```

---

## ✅ **Summary**

**To switch from dev to prod:**

1. ✅ Create `.env.production` with production API URL
2. ✅ Build: `npm run build`
3. ✅ Deploy to server
4. ✅ Update nginx CORS origins
5. ✅ Configure SSL (if needed)
6. ✅ Verify: Check browser console for correct API URL

**That's it!** The admin dashboard will now use production data from the production backend. 🚀
