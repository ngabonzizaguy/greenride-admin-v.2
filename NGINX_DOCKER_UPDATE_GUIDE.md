# 🐳 nginx Docker Update Guide

> **Setup:** nginx running in Docker container  
> **Config File:** `~/nginx.conf` (mounted into container)  
> **Docker Compose:** `docker-compose-greenride-nginx.yml`

---

## ✅ **Yes, Update the File in `~`**

Since nginx is running in Docker, the `nginx.conf` in your home directory is mounted into the container. **Update that file directly.**

---

## 📋 **Step-by-Step Instructions**

### **Step 1: Backup Current Config**

```bash
# On cloudshell
cd ~
cp nginx.conf nginx.conf.backup
```

### **Step 2: Update nginx.conf**

```bash
# Edit the file
nano nginx.conf

# OR use vi
vi nginx.conf
```

**Then:**
1. Copy the **entire content** from `nginx-config-fixed.conf`
2. Paste it into `nginx.conf` (replace everything)
3. Save and exit (`Ctrl+X`, then `Y`, then `Enter` for nano)

### **Step 3: Verify Docker Compose Setup**

Check how nginx.conf is mounted:

```bash
# View docker-compose file
cat docker-compose-greenride-nginx.yml
```

**Expected volume mount:**
```yaml
volumes:
  - ./nginx.conf:/etc/nginx/nginx.conf
```

### **Step 4: Test nginx Configuration**

**Option A: Test inside container (Recommended)**

```bash
# Find nginx container name
docker ps | grep nginx

# Test config inside container
docker exec <nginx-container-name> nginx -t

# Example:
docker exec greenride-nginx nginx -t
```

**Expected output:**
```
nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
nginx: configuration file /etc/nginx/nginx.conf test is successful
```

**Option B: Test locally (if nginx is installed)**

```bash
nginx -t -c ~/nginx.conf
```

### **Step 5: Reload/Restart nginx Container**

**Option A: Reload (No Downtime - Recommended)**

```bash
# Reload nginx inside container
docker exec <nginx-container-name> nginx -s reload

# Example:
docker exec greenride-nginx nginx -s reload
```

**Option B: Restart Container**

```bash
# Restart using docker-compose
docker-compose -f docker-compose-greenride-nginx.yml restart nginx

# OR restart specific container
docker restart <nginx-container-name>
```

**Option C: Recreate Container (if needed)**

```bash
# Stop, remove, and recreate
docker-compose -f docker-compose-greenride-nginx.yml up -d --force-recreate nginx
```

---

## 🔍 **Verify the Update**

### **1. Check nginx is Running**

```bash
docker ps | grep nginx
```

### **2. Test API Endpoints**

```bash
# Test Mobile API
curl -I http://localhost/api/health
# OR from outside:
curl -I http://18.143.118.157/api/health

# Test Admin API
curl -I http://localhost/admin/api/health
# OR from outside:
curl -I http://18.143.118.157/admin/api/health
```

### **3. Check nginx Logs**

```bash
# View nginx access logs
docker logs <nginx-container-name>

# OR follow logs in real-time
docker logs -f <nginx-container-name>
```

---

## 🐛 **Troubleshooting**

### **Problem: nginx -t fails**

**Check syntax errors:**
```bash
docker exec <nginx-container-name> nginx -t
```

**Common issues:**
- Missing semicolons
- Incorrect brackets
- Invalid directives

**Fix:** Review error message and correct the config file.

### **Problem: Container won't start**

**Check container logs:**
```bash
docker logs <nginx-container-name>
```

**Check if port is already in use:**
```bash
sudo netstat -tulpn | grep :80
```

### **Problem: Changes not taking effect**

**Make sure you:**
1. ✅ Saved the file (`Ctrl+X`, `Y`, `Enter`)
2. ✅ Reloaded nginx (`docker exec ... nginx -s reload`)
3. ✅ Checked container is using the mounted file

**Verify mount:**
```bash
# Check if file is mounted correctly
docker exec <nginx-container-name> cat /etc/nginx/nginx.conf | head -20
```

**Compare with your local file:**
```bash
head -20 ~/nginx.conf
```

---

## 📝 **Quick Reference Commands**

```bash
# 1. Backup
cp ~/nginx.conf ~/nginx.conf.backup

# 2. Edit
nano ~/nginx.conf

# 3. Test
docker exec <nginx-container-name> nginx -t

# 4. Reload
docker exec <nginx-container-name> nginx -s reload

# 5. Check logs
docker logs -f <nginx-container-name>

# 6. Restart container (if needed)
docker-compose -f docker-compose-greenride-nginx.yml restart nginx
```

---

## ✅ **Summary**

1. ✅ **Update `~/nginx.conf`** (the file in your home directory)
2. ✅ **Test:** `docker exec <container> nginx -t`
3. ✅ **Reload:** `docker exec <container> nginx -s reload`
4. ✅ **Verify:** Test API endpoints

**That's it!** The Docker container will automatically use the updated config file because it's mounted as a volume.

---

## 🔗 **Docker Compose File Reference**

If you need to check or modify the docker-compose file:

```bash
# View the compose file
cat docker-compose-greenride-nginx.yml

# Expected structure:
# services:
#   nginx:
#     image: nginx:alpine
#     volumes:
#       - ./nginx.conf:/etc/nginx/nginx.conf
#     ports:
#       - "80:80"
#       - "8080:8080"
```

---

**Ready to update!** 🚀
