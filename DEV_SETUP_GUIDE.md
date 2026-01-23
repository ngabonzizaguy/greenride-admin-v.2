# 🛠️ Development Setup Guide

> **Purpose:** Configure Admin Dashboard to connect to Remote Dev Backend  
> **Target:** Development environment (not production)  
> **Backend:** `18.143.118.157:8611` (Remote Dev Backend)

---

## 🎯 **Goal**

Connect your Admin Dashboard to the **Remote Dev Backend** instead of localhost, so you:
- ✅ Don't need to run local backend
- ✅ Use real dev data from remote server
- ✅ Avoid slow database queries (remote backend is closer to database)
- ✅ Faster development workflow

---

## 📋 **Solution: Configure Admin Dashboard for Remote Dev Backend**

### **Step 1: Update `.env.local` File**

Edit `.env.local` in the root directory:

```bash
# Remote Dev Backend (18.143.118.157:8611)
NEXT_PUBLIC_API_URL=http://18.143.118.157:8611

# Use real API (not mock data)
NEXT_PUBLIC_DEMO_MODE=false

# Google Maps API Key (if needed)
NEXT_PUBLIC_GOOGLE_MAPS_KEY=AIzaSyDif39v3Gx4YXonS3-A8pINUMi3hxRfC3U
```

### **Step 2: Restart Next.js Dev Server**

After changing `.env.local`, **restart your Next.js dev server**:

```bash
# Stop the current dev server (Ctrl+C)
# Then restart:
npm run dev
```

**Important:** Next.js only loads `.env.local` on startup, so changes require a restart.

---

## ✅ **Verification**

### **1. Check Environment Variables**

Open browser console (F12) and look for:
```
[API Client] API_BASE_URL: http://18.143.118.157:8611
[API Client] DEMO_MODE: false
```

### **2. Test Connection**

1. Open Admin Dashboard: `http://localhost:3000`
2. Try to login or load dashboard
3. Check Network tab (F12 → Network) - should see requests to `18.143.118.157:8611`

### **3. Expected Behavior**

- ✅ Admin dashboard connects to remote dev backend
- ✅ No need to run local backend
- ✅ Uses real data from dev database
- ✅ Faster than local backend (remote backend is closer to database)

---

## 🔄 **Switching Between Local and Remote**

### **Option A: Use Local Backend (for local development)**

```bash
# .env.local
NEXT_PUBLIC_API_URL=http://localhost:8611
NEXT_PUBLIC_DEMO_MODE=false
```

**Requires:**
- Local backend must be running (`go run main/main.go`)
- Local backend connects to remote database

### **Option B: Use Remote Dev Backend (recommended for dev)**

```bash
# .env.local
NEXT_PUBLIC_API_URL=http://18.143.118.157:8611
NEXT_PUBLIC_DEMO_MODE=false
```

**Benefits:**
- ✅ No need to run local backend
- ✅ Faster (remote backend is closer to database)
- ✅ Uses real dev data
- ✅ Always available (remote backend runs 24/7)

---

## 📊 **Configuration Comparison**

| Configuration | API URL | Requires Local Backend? | Database | Speed |
|---------------|---------|-------------------------|----------|-------|
| **Local Dev** | `http://localhost:8611` | ✅ Yes | Remote (via local backend) | ⚠️ Slow (remote DB) |
| **Remote Dev** | `http://18.143.118.157:8611` | ❌ No | Remote (direct) | ✅ Fast (close to DB) |
| **Production** | `https://api.greenrideafrica.com:8611` | ❌ No | Production | ✅ Fast |

---

## 🚀 **Recommended Setup for Development**

### **For Daily Development:**

1. **Use Remote Dev Backend** (`18.143.118.157:8611`)
   - Faster workflow
   - No need to run local backend
   - Real dev data

2. **Use Local Backend Only When:**
   - Testing backend changes
   - Debugging backend code
   - Adding new backend features

---

## 🔧 **Troubleshooting**

### **Problem: Admin dashboard can't connect to remote backend**

**Check:**
1. Is `.env.local` correct? (`NEXT_PUBLIC_API_URL=http://18.143.118.157:8611`)
2. Did you restart Next.js dev server?
3. Is remote backend accessible? (`curl http://18.143.118.157:8611/health`)
4. Check browser console for errors

**Solution:**
```bash
# Test remote backend connectivity
curl http://18.143.118.157:8611/health

# If not accessible, check:
# - Internet connection
# - Firewall settings
# - VPN (if required)
```

### **Problem: CORS errors**

**Cause:** Remote backend might not allow requests from `localhost:3000`

**Solution:** Backend needs to allow CORS from your origin. Check backend CORS configuration.

### **Problem: Slow responses**

**Cause:** Network latency to remote backend

**Solution:**
- Use local backend if latency is too high
- Or optimize backend queries
- Add caching (Redis) on backend

---

## 📝 **Quick Reference**

### **Current Setup (after this guide):**

```bash
# .env.local
NEXT_PUBLIC_API_URL=http://18.143.118.157:8611  # Remote Dev Backend
NEXT_PUBLIC_DEMO_MODE=false                      # Use real API
```

### **To Switch Back to Local:**

```bash
# .env.local
NEXT_PUBLIC_API_URL=http://localhost:8611        # Local Backend
NEXT_PUBLIC_DEMO_MODE=false                      # Use real API
```

**Remember:** Always restart Next.js dev server after changing `.env.local`!

---

## ✅ **Summary**

1. **Edit `.env.local`**: Change `NEXT_PUBLIC_API_URL` to `http://18.143.118.157:8611`
2. **Restart Next.js**: Stop and restart `npm run dev`
3. **Verify**: Check browser console for correct API URL
4. **Done!**: Admin dashboard now uses remote dev backend

**Benefits:**
- ✅ No need to run local backend
- ✅ Faster development workflow
- ✅ Uses real dev data
- ✅ Always available

---

**Ready to proceed?** Follow the steps above! 🚀
