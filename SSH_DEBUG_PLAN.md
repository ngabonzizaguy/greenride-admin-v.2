# 🔧 SSH Debug Plan - Feedback 404 Issue

## 📋 SSH Connection Details
- **Host:** `ec2-18-143-118-157.ap-southeast-1.compute.amazonaws.com` (or `18.143.118.157`)
- **Username:** `ubuntu`
- **Key File:** `dev.pem`

## 🎯 What We'll Do Next

### Step 1: Connect via SSH ✅
First, we need to locate the `dev.pem` key file or connect via SSH.

**Option A: If `dev.pem` is available locally**
```bash
# Set correct permissions (required for SSH key)
icacls dev.pem /inheritance:r
icacls dev.pem /grant:r "%username%:R"

# Test SSH connection
ssh -i dev.pem ubuntu@18.143.118.157 "echo 'Connection successful!'"
```

**Option B: If key is in different location**
Provide the full path to `dev.pem`, or we can use AWS Session Manager/EC2 Instance Connect if configured.

### Step 2: Check Backend Status 🔍
```bash
# Check if backend is running on port 8611
ssh -i dev.pem ubuntu@18.143.118.157 "sudo netstat -tlnp | grep 8611 || sudo ss -tlnp | grep 8611"

# Check backend health
ssh -i dev.pem ubuntu@18.143.118.157 "curl -s http://localhost:8611/health"
```

### Step 3: Test Feedback Endpoints Directly 🧪
```bash
# Test /feedback/stats directly (bypass nginx)
ssh -i dev.pem ubuntu@18.143.118.157 "curl -s -w '\nHTTP: %{http_code}\n' http://localhost:8611/feedback/stats"

# Test /feedback/search directly
ssh -i dev.pem ubuntu@18.143.118.157 "curl -s -X POST -w '\nHTTP: %{http_code}\n' -H 'Content-Type: application/json' -d '{\"page\":1,\"limit\":10}' http://localhost:8611/feedback/search"
```

### Step 4: Test Through Nginx 🌐
```bash
# Test through nginx (same as frontend)
ssh -i dev.pem ubuntu@18.143.118.157 "curl -s -w '\nHTTP: %{http_code}\n' http://localhost/admin/api/feedback/stats"

# Compare with working endpoint
ssh -i dev.pem ubuntu@18.143.118.157 "curl -s -w '\nHTTP: %{http_code}\n' http://localhost/admin/api/dashboard/stats"
```

### Step 5: Check Backend Process/Container 🐳
```bash
# Check Docker containers
ssh -i dev.pem ubuntu@18.143.118.157 "docker ps | grep -E 'greenride|backend|api'"

# Check systemd services
ssh -i dev.pem ubuntu@18.143.118.157 "systemctl list-units --type=service --state=running | grep -i greenride"

# Check running processes
ssh -i dev.pem ubuntu@18.143.118.157 "ps aux | grep -E 'greenride|backend.*8611|main.*8611' | grep -v grep"
```

### Step 6: Check Backend Logs 📋
```bash
# If running in Docker
ssh -i dev.pem ubuntu@18.143.118.157 "docker logs --tail 50 <container-name> | grep -i feedback"

# If running as systemd service
ssh -i dev.pem ubuntu@18.143.118.157 "journalctl -u greenride-admin --tail 50 | grep -i feedback"

# Check for route registration
ssh -i dev.pem ubuntu@18.143.118.157 "docker logs --tail 100 <container-name> | grep -E 'GET.*feedback|POST.*feedback|feedback/stats|feedback/search'"
```

### Step 7: Restart Backend (If Needed) 🔄
Based on how backend runs:

**If Docker:**
```bash
ssh -i dev.pem ubuntu@18.143.118.157 "docker restart <container-name>"
```

**If systemd:**
```bash
ssh -i dev.pem ubuntu@18.143.118.157 "sudo systemctl restart greenride-admin"
```

**If manual process:**
```bash
ssh -i dev.pem ubuntu@18.143.118.157 "pkill -f 'greenride.*8611' && cd /path/to/backend && ./start.sh"
```

### Step 8: Verify Fix ✅
```bash
# Test endpoints again after restart
ssh -i dev.pem ubuntu@18.143.118.157 "curl -s http://localhost:8611/feedback/stats | head -c 200"
ssh -i dev.pem ubuntu@18.143.118.157 "curl -s http://localhost/admin/api/feedback/stats | head -c 200"
```

## 🚀 Next Action

**Please provide one of the following:**

1. **Path to `dev.pem` key file** (if different from current directory)
   - Example: `C:\Users\YourName\Downloads\dev.pem`
   - Or: `D:\keys\dev.pem`

2. **Or confirm you want to use AWS Session Manager/EC2 Instance Connect** instead

Once we have the key file location, I'll:
1. ✅ Test SSH connection
2. ✅ Run all debug commands
3. ✅ Identify the issue
4. ✅ Fix it (restart backend)
5. ✅ Verify the fix works

---

## 📝 Expected Results

After running the commands, we should see:
- ✅ **Backend is running** on port 8611
- ✅ **Endpoints return 200 or 401** (not 404)
- ✅ **Feedback routes are registered** in logs
- ✅ **Nginx rewrite works** correctly

If we see 404, we'll restart the backend and it should fix the issue.
