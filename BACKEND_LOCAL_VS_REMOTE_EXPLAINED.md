# 🔍 Backend: Local vs Remote - Complete Explanation

> **For:** Your Understanding  
> **Purpose:** Clarify local vs remote backend, database connections, and which services connect to what

---

## 🚨 **Issue 1: Why Backend Keeps Failing**

### **The Error:**
```
listen tcp :8610: bind: Only one usage of each socket address (protocol/network address/port) is normally permitted.
listen tcp :8611: bind: Only one usage of each socket address (protocol/network address/port) is normally permitted.
```

### **Root Cause:**
**Ports 8610 and 8611 are already in use** by another backend process (PID: 18844 - `main.exe`)

### **Solution:**
Kill the existing process before starting a new one:

```bash
# Option 1: Kill by PID
taskkill /F /PID 18844

# Option 2: Kill by process name
taskkill /F /IM main.exe

# Option 3: Kill all Go processes
taskkill /F /IM go.exe
taskkill /F /IM main.exe
```

**Then restart:**
```bash
cd backend/greenride-api-clean
go run main/main.go
```

### **Why This Happens:**
- Previous backend instance didn't shut down cleanly
- Multiple terminal windows running backend simultaneously
- Process crashed but port wasn't released
- Windows sometimes holds ports for a few seconds after process ends

---

## 🌐 **Local vs Remote Backend - Are They Connected?**

### **Short Answer: NO - They are COMPLETELY SEPARATE**

### **Detailed Explanation:**

#### **1. Local Backend (On Your Computer)**
- **Location:** `D:\greenride-admin-v.2\backend\greenride-api-clean\`
- **Runs on:** `localhost` or `127.0.0.1`
- **Ports:** 
  - `8610` - Mobile API
  - `8611` - Admin API
- **Database:** Connects to **REMOTE database** at `18.143.118.157:3306`
- **Status:** Only runs when you start it manually
- **Purpose:** Local development and testing

#### **2. Remote Backend (On AWS Server)**
- **Location:** AWS EC2 server (`18.143.118.157`)
- **Runs on:** Cloud server (always running)
- **Ports:** 
  - `8610` - Mobile API
  - `8611` - Admin API
- **Database:** Connects to **SAME remote database** at `18.143.118.157:3306`
- **Status:** Running 24/7 (production)
- **Purpose:** Serves mobile app and admin dashboard in production

### **Key Point: They Share the SAME Database!**

```
┌─────────────────┐         ┌──────────────────┐
│  Local Backend  │         │  Remote Backend  │
│  (Your PC)      │         │  (AWS Server)    │
│  Port 8610/8611 │         │  Port 8610/8611  │
└────────┬────────┘         └────────┬─────────┘
         │                            │
         │                            │
         └────────────┬───────────────┘
                      │
                      ▼
         ┌─────────────────────────┐
         │   MySQL Database         │
         │   18.143.118.157:3306   │
         │   Database: greenride    │
         └─────────────────────────┘
```

**Both backends connect to the SAME database**, so:
- ✅ Data created by local backend appears in remote backend
- ✅ Data created by remote backend appears in local backend
- ⚠️ **They can conflict** if both modify the same data simultaneously
- ⚠️ **They are NOT synced** - they're just reading/writing to the same database

---

## 📱 **Which Backend Does the Mobile App Talk To?**

### **Mobile App (Flutter) Configuration:**

The mobile app **does NOT** connect to your local backend. It connects to:

#### **Development/Testing:**
- **URL:** `http://18.143.118.157:8610/` (Remote backend on AWS)
- **Port:** 8610 (Mobile API)

#### **Production:**
- **URL:** `https://api.greenrideafrica.com/` (via nginx on AWS)
- **Port:** 8610 (Mobile API)

### **Why Not Local?**
- Mobile app runs on a **physical device** or **emulator**
- `localhost` on your computer is **not accessible** from mobile device
- Mobile app needs a **network-accessible IP address**
- Your local backend is only accessible from your computer

### **To Test Mobile App with Local Backend:**
You would need to:
1. Find your computer's local IP (e.g., `192.168.1.100`)
2. Configure mobile app to use `http://192.168.1.100:8610`
3. Ensure mobile device and computer are on same network
4. Configure firewall to allow port 8610

**This is NOT recommended** - use remote backend for mobile app testing.

---

## 💻 **Which Backend Does the Admin Dashboard Talk To?**

### **Admin Dashboard (Next.js) Configuration:**

The admin dashboard **CAN** connect to either local or remote backend, depending on `.env.local`:

#### **Current Configuration (`.env.local`):**
```bash
NEXT_PUBLIC_API_URL=http://localhost:8611
NEXT_PUBLIC_DEMO_MODE=false
```

**This means:**
- ✅ Admin dashboard connects to **LOCAL backend** (port 8611)
- ✅ Only works when local backend is running
- ✅ Uses real data from database (not mock data)

#### **To Connect to Remote Backend:**
Change `.env.local`:
```bash
NEXT_PUBLIC_API_URL=http://18.143.118.157:8611
NEXT_PUBLIC_DEMO_MODE=false
```

**This would:**
- ✅ Admin dashboard connects to **REMOTE backend** (AWS server)
- ✅ Works even if local backend is not running
- ✅ Uses production data

---

## 🔄 **Data Flow Diagram**

### **Scenario 1: Admin Dashboard → Local Backend**
```
Admin Dashboard (Next.js)
    ↓ (http://localhost:8611)
Local Backend (Your PC)
    ↓ (MySQL connection)
Remote Database (18.143.118.157:3306)
    ↑ (MySQL connection)
Remote Backend (AWS Server)
    ↑ (http://18.143.118.157:8610)
Mobile App (Flutter)
```

**Result:**
- Admin dashboard sees data from database
- Mobile app sees data from database
- Both see the **same data** (same database)
- But they're talking to **different backend instances**

### **Scenario 2: Admin Dashboard → Remote Backend**
```
Admin Dashboard (Next.js)
    ↓ (http://18.143.118.157:8611)
Remote Backend (AWS Server)
    ↓ (MySQL connection)
Remote Database (18.143.118.157:3306)
    ↑ (MySQL connection)
Remote Backend (AWS Server)
    ↑ (http://18.143.118.157:8610)
Mobile App (Flutter)
```

**Result:**
- Both admin dashboard and mobile app use **same backend instance**
- Both see the **same data** (same database)
- More consistent, but uses production backend

---

## 🗄️ **Database Connection Details**

### **Both Backends Use Same Database:**

**Local Backend:**
```yaml
# config.yaml or dev.yaml
database:
  dsn: "greenride:GreenRide2024!@tcp(18.143.118.157:3306)/greenride"
```

**Remote Backend:**
```yaml
# Same config on server
database:
  dsn: "greenride:GreenRide2024!@tcp(18.143.118.157:3306)/greenride"
```

**Key Points:**
- ✅ **Same database host:** `18.143.118.157:3306`
- ✅ **Same database name:** `greenride`
- ✅ **Same credentials:** `greenride:GreenRide2024!`
- ✅ **Same tables:** Both use `t_users`, `t_orders`, `t_feedbacks`, etc.

**This means:**
- Data created by local backend → visible to remote backend
- Data created by remote backend → visible to local backend
- They're **reading/writing to the same database**
- But they're **separate processes** running on different machines

---

## ⚠️ **Important Warnings**

### **1. Data Conflicts**
If both backends are running simultaneously:
- ⚠️ They can modify the same data at the same time
- ⚠️ Race conditions possible
- ⚠️ Last write wins (data can be overwritten)

**Recommendation:** Only run one backend at a time for development.

### **2. Performance Impact**
- Local backend connects to **remote database** (AWS Singapore)
- Network latency can be **very slow** (as you've experienced)
- Database queries take longer from your location

**Recommendation:** Use remote backend for admin dashboard if local is too slow.

### **3. Port Conflicts**
- Local backend uses ports 8610 and 8611
- Remote backend uses same ports (but on different machine)
- **No conflict** between local and remote (different machines)
- **Conflict** if you try to run two local instances

---

## 🎯 **Quick Reference**

### **Which Backend Should I Use?**

| Use Case | Backend | Why |
|----------|---------|-----|
| **Local development** | Local (`localhost:8611`) | Faster iteration, easier debugging |
| **Testing with real data** | Remote (`18.143.118.157:8611`) | Uses production data, no local setup needed |
| **Mobile app** | Remote (`18.143.118.157:8610`) | Mobile device can't access localhost |
| **Production admin** | Remote (`18.143.118.157:8611`) | Always available, production-ready |

### **Current Setup:**

**Admin Dashboard:**
- ✅ Connects to: **Local backend** (`localhost:8611`)
- ✅ Requires: Local backend to be running
- ✅ Database: Remote database (shared with remote backend)

**Mobile App:**
- ✅ Connects to: **Remote backend** (`18.143.118.157:8610`)
- ✅ Always works: Remote backend runs 24/7
- ✅ Database: Remote database (shared with local backend)

---

## 🔧 **Troubleshooting**

### **Problem: Backend won't start (port in use)**

**Solution:**
```bash
# Find process using ports
netstat -ano | findstr ":8610 :8611"

# Kill the process (replace PID with actual process ID)
taskkill /F /PID <PID>

# Or kill all Go processes
taskkill /F /IM go.exe
taskkill /F /IM main.exe
```

### **Problem: Admin dashboard can't connect**

**Check:**
1. Is local backend running? (`netstat -ano | findstr ":8611"`)
2. Is `.env.local` correct? (`NEXT_PUBLIC_API_URL=http://localhost:8611`)
3. Is `NEXT_PUBLIC_DEMO_MODE=false`?
4. Restart Next.js dev server after changing `.env.local`

### **Problem: Slow database queries**

**Cause:** Local backend connecting to remote database (AWS Singapore) from your location

**Solutions:**
1. Use remote backend instead (change `.env.local` to `18.143.118.157:8611`)
2. Add caching (Redis) for slow queries
3. Optimize database queries (add indexes)
4. Use nginx caching (already configured)

---

## 📊 **Summary Table**

| Aspect | Local Backend | Remote Backend |
|--------|---------------|---------------|
| **Location** | Your computer | AWS EC2 (`18.143.118.157`) |
| **Ports** | 8610, 8611 | 8610, 8611 |
| **Database** | Remote (`18.143.118.157:3306`) | Remote (`18.143.118.157:3306`) |
| **Status** | Manual start | Always running |
| **Admin Dashboard** | Can connect (if running) | Can connect (always) |
| **Mobile App** | Cannot connect | Always connects |
| **Data** | Same database (shared) | Same database (shared) |
| **Synced?** | NO - separate processes | NO - separate processes |

---

## ✅ **Key Takeaways**

1. **Local and remote backends are SEPARATE** - they don't sync or communicate
2. **They share the SAME database** - both read/write to `18.143.118.157:3306`
3. **Mobile app ALWAYS uses remote backend** - can't access localhost
4. **Admin dashboard can use either** - depends on `.env.local` configuration
5. **Port conflicts happen** - kill old processes before starting new ones
6. **Slow queries are normal** - local backend connecting to remote database

---

**This should clarify everything!** If you have more questions, let me know! 🙏
