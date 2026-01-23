# 🔧 How to Replace nginx.conf on Cloudshell

> **Goal:** Replace entire `~/nginx.conf` with contents of `nginx-config-fixed.conf`

---

## ✅ **Method 1: Using nano (Recommended - Visual)**

### **Step 1: Open file in nano**
```bash
nano ~/nginx.conf
```

### **Step 2: Select all and delete**
1. Press `Ctrl + A` (select all)
2. Press `Delete` or `Backspace` (delete everything)

### **Step 3: Paste new content**
1. **If using SSH client (PuTTY, Windows Terminal, etc.):**
   - Right-click in the terminal → Paste
   - OR `Shift + Insert` (Windows)
   - OR `Ctrl + Shift + V` (Linux/Mac)

2. **If using cloudshell web interface:**
   - Right-click → Paste
   - OR `Ctrl + V`

### **Step 4: Save and exit**
1. Press `Ctrl + O` (write out/save)
2. Press `Enter` (confirm filename)
3. Press `Ctrl + X` (exit)

---

## ✅ **Method 2: Using cat with heredoc (Fastest)**

**Copy this entire block and paste into cloudshell:**

```bash
cat > ~/nginx.conf << 'NGINX_EOF'
worker_processes auto;
events {
    worker_connections 1024;
    use epoll;
}

http {
    # GreenRide 后端API服务 (Mobile API - Port 8610)
    upstream greenride_api_backend {
        server host.docker.internal:8610;
        keepalive 32;
        keepalive_requests 1000;
        keepalive_timeout 60s;
    }
    
    # GreenRide 管理后台API服务 (Admin API - Port 8611)
    upstream greenride_admin_api_backend {
        server host.docker.internal:8611;
        keepalive 32;
        keepalive_requests 1000;
        keepalive_timeout 60s;
    }
    
    # GreenRide WebSocket服务
    upstream greenride_websocket_backend {
        server host.docker.internal:8612;
        keepalive 32;
        keepalive_requests 1000;
        keepalive_timeout 60s;
    }
    
    # GreenRide 前端应用 (Mobile App Frontend)
    upstream greenride_frontend_backend {
        server host.docker.internal:3000;
        keepalive 32;
        keepalive_requests 1000;
        keepalive_timeout 60s;
    }
    
    # GreenRide 管理后台前端 (Admin Dashboard)
    upstream greenride_admin_frontend_backend {
        server host.docker.internal:3001;
        keepalive 32;
        keepalive_requests 1000;
        keepalive_timeout 60s;
    }
    
    # 日志格式
    log_format detailed '$remote_addr - $remote_user [$time_local] "$request" '
                        '$status $body_bytes_sent "$http_referer" '
                        '"$http_user_agent" "$http_x_forwarded_for" '
                        'upstream_addr=$upstream_addr '
                        'upstream_status=$upstream_status '
                        'upstream_response_time=$upstream_response_time '
                        'request_time=$request_time';

    # MIME types
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    # CORS allowlist (admin + local dev)
    map $http_origin $cors_allow_origin {
        default "";
        "http://localhost:3000" $http_origin;
        "http://localhost:3600" $http_origin;
        "https://admin.greenrideafrica.com" $http_origin;
    }
    
    # 文件传输优化
    sendfile        on;
    tcp_nopush      on;
    tcp_nodelay     on;
    keepalive_timeout 65;
    types_hash_max_size 2048;
    
    # 压缩配置
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_proxied any;
    gzip_comp_level 6;
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/json
        application/javascript
        application/xml+rss
        application/atom+xml
        image/svg+xml;

    # ============================================================================
    # GreenRide 主站 - 用户端 (Port 80)
    # ============================================================================
    server {
        listen 80;
        server_name _;
        
        # 日志配置
        access_log /var/log/nginx/greenride_main_access.log detailed;
        error_log /var/log/nginx/greenride_main_error.log info;
        
        # Client settings
        client_max_body_size 10M;
        client_body_timeout 60s;
        client_header_timeout 60s;
        
        # ============================================================================
        # Mobile API Routes - /api/* → Backend Port 8610 (root path)
        # ============================================================================
        # IMPORTANT: Backend uses root path (/), so we strip /api prefix
        location /api/ {
            # Handle OPTIONS preflight requests
            if ($request_method = 'OPTIONS') {
                add_header 'Access-Control-Allow-Origin' $cors_allow_origin always;
                add_header 'Vary' 'Origin' always;
                add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS' always;
                add_header 'Access-Control-Allow-Headers' 'Accept, Accept-Language, Content-Language, Content-Type, Authorization, X-Requested-With, X-Real-IP, X-Forwarded-For, X-Forwarded-Proto, Cache-Control, DNT, User-Agent, If-Modified-Since, Range' always;
                add_header 'Access-Control-Allow-Credentials' 'true' always;
                add_header 'Access-Control-Max-Age' '86400' always;
                add_header 'Content-Type' 'text/plain; charset=utf-8' always;
                add_header 'Content-Length' '0' always;
                return 204;
            }
            
            # Strip /api prefix and proxy to backend root path
            rewrite ^/api/(.*)$ /$1 break;
            proxy_pass http://greenride_api_backend;
            
            # Proxy headers
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_set_header X-Forwarded-Host $host;
            proxy_set_header X-Forwarded-Port $server_port;
            proxy_redirect off;
            port_in_redirect off;
            
            # Timeout configuration
            proxy_connect_timeout 60s;
            proxy_send_timeout 60s;
            proxy_read_timeout 120s;
            proxy_buffering on;
            proxy_buffer_size 4k;
            proxy_buffers 8 4k;
            proxy_busy_buffers_size 8k;
            
            # HTTP/1.1 and connection keepalive
            proxy_http_version 1.1;
            proxy_set_header Connection "";
            
            # CORS headers - Add to all responses
            add_header 'Access-Control-Allow-Origin' $cors_allow_origin always;
            add_header 'Vary' 'Origin' always;
            add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS' always;
            add_header 'Access-Control-Allow-Headers' 'Accept, Accept-Language, Content-Language, Content-Type, Authorization, X-Requested-With, X-Real-IP, X-Forwarded-For, X-Forwarded-Proto, Cache-Control, DNT, User-Agent, If-Modified-Since, Range' always;
            add_header 'Access-Control-Allow-Credentials' 'true' always;
            add_header 'Access-Control-Expose-Headers' 'Content-Length, Content-Range, Content-Type, Date, Server, Transfer-Encoding' always;
        }

        # ============================================================================
        # Admin API Routes - /admin/api/* → Backend Port 8611 (root path)
        # ============================================================================
        # IMPORTANT: Backend uses root path (/), so we strip /admin/api prefix
        location /admin/api/ {
            # Handle OPTIONS preflight requests
            if ($request_method = 'OPTIONS') {
                add_header 'Access-Control-Allow-Origin' $cors_allow_origin always;
                add_header 'Vary' 'Origin' always;
                add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS' always;
                add_header 'Access-Control-Allow-Headers' 'Accept, Accept-Language, Content-Language, Content-Type, Authorization, X-Requested-With, X-Real-IP, X-Forwarded-For, X-Forwarded-Proto, Cache-Control, DNT, User-Agent, If-Modified-Since, Range' always;
                add_header 'Access-Control-Allow-Credentials' 'true' always;
                add_header 'Access-Control-Max-Age' '86400' always;
                add_header 'Content-Type' 'text/plain; charset=utf-8' always;
                add_header 'Content-Length' '0' always;
                return 204;
            }
            
            # Strip /admin/api prefix and proxy to backend root path
            rewrite ^/admin/api/(.*)$ /$1 break;
            proxy_pass http://greenride_admin_api_backend;
            
            # Proxy headers
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_set_header X-Forwarded-Host $host;
            proxy_set_header X-Forwarded-Port $server_port;
            proxy_redirect off;
            port_in_redirect off;
            
            # Timeout configuration
            proxy_connect_timeout 60s;
            proxy_send_timeout 60s;
            proxy_read_timeout 120s;
            proxy_buffering on;
            proxy_buffer_size 4k;
            proxy_buffers 8 4k;
            proxy_busy_buffers_size 8k;
            
            # HTTP/1.1 and connection keepalive
            proxy_http_version 1.1;
            proxy_set_header Connection "";
            
            # CORS headers - Add to all responses
            add_header 'Access-Control-Allow-Origin' $cors_allow_origin always;
            add_header 'Vary' 'Origin' always;
            add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS' always;
            add_header 'Access-Control-Allow-Headers' 'Accept, Accept-Language, Content-Language, Content-Type, Authorization, X-Requested-With, X-Real-IP, X-Forwarded-For, X-Forwarded-Proto, Cache-Control, DNT, User-Agent, If-Modified-Since, Range' always;
            add_header 'Access-Control-Allow-Credentials' 'true' always;
            add_header 'Access-Control-Expose-Headers' 'Content-Length, Content-Range, Content-Type, Date, Server, Transfer-Encoding' always;
        }

        # ============================================================================
        # Admin Frontend Static Files - /admin/static/*
        # ============================================================================
        location ~* ^/admin/static/(.+)$ {
            proxy_pass http://greenride_admin_frontend_backend/static/$1;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            expires 1y;
            add_header Cache-Control "public, immutable";
            add_header X-Served-By "nginx-admin";
        }
        
        # ============================================================================
        # Static Files Cache - Main Frontend
        # ============================================================================
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
            proxy_pass http://greenride_frontend_backend;
            expires 1y;
            add_header Cache-Control "public, immutable";
            add_header X-Served-By "nginx";
        }
        
        # ============================================================================
        # Admin Frontend Routes - /admin → Admin Frontend (Port 3001)
        # ============================================================================
        location /admin {
            proxy_pass http://greenride_admin_frontend_backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_redirect off;
            port_in_redirect off;
            
            # HTTP/1.1 and connection keepalive
            proxy_http_version 1.1;
            proxy_set_header Connection "";
        }
        
        # ============================================================================
        # WebSocket Connection - /ws
        # ============================================================================
        location /ws {
            proxy_pass http://greenride_websocket_backend;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "Upgrade";
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_cache_bypass $http_upgrade;
            proxy_read_timeout 86400;
        }
        
        # ============================================================================
        # Main Frontend Application - / → Mobile Frontend (Port 3000)
        # ============================================================================
        location / {
            proxy_pass http://greenride_frontend_backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_redirect off;
            port_in_redirect off;
            
            # HTTP/1.1 and connection keepalive
            proxy_http_version 1.1;
            proxy_set_header Connection "";
        }
    }

    # ============================================================================
    # GreenRide 管理后台专用端口 (Port 8080)
    # ============================================================================
    server {
        listen 8080;
        server_name _;
        
        # 日志配置
        access_log /var/log/nginx/greenride_admin_access.log detailed;
        error_log /var/log/nginx/greenride_admin_error.log info;
        
        # Client settings
        client_max_body_size 10M;
        client_body_timeout 60s;
        client_header_timeout 60s;
        
        # ============================================================================
        # Admin API Routes - /admin/api/* → Backend Port 8611 (root path)
        # ============================================================================
        location /admin/api/ {
            # Handle OPTIONS preflight requests
            if ($request_method = 'OPTIONS') {
                add_header 'Access-Control-Allow-Origin' $cors_allow_origin always;
                add_header 'Vary' 'Origin' always;
                add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS' always;
                add_header 'Access-Control-Allow-Headers' 'Accept, Accept-Language, Content-Language, Content-Type, Authorization, X-Requested-With, X-Real-IP, X-Forwarded-For, X-Forwarded-Proto, Cache-Control, DNT, User-Agent, If-Modified-Since, Range' always;
                add_header 'Access-Control-Allow-Credentials' 'true' always;
                add_header 'Access-Control-Max-Age' '86400' always;
                add_header 'Content-Type' 'text/plain; charset=utf-8' always;
                add_header 'Content-Length' '0' always;
                return 204;
            }
            
            # Strip /admin/api prefix and proxy to backend root path
            rewrite ^/admin/api/(.*)$ /$1 break;
            proxy_pass http://greenride_admin_api_backend;
            
            # Proxy headers
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_set_header X-Forwarded-Host $host;
            proxy_set_header X-Forwarded-Port $server_port;
            proxy_redirect off;
            port_in_redirect off;
            
            # Timeout configuration
            proxy_connect_timeout 60s;
            proxy_send_timeout 60s;
            proxy_read_timeout 120s;
            proxy_buffering on;
            proxy_buffer_size 4k;
            proxy_buffers 8 4k;
            proxy_busy_buffers_size 8k;
            
            # HTTP/1.1 and connection keepalive
            proxy_http_version 1.1;
            proxy_set_header Connection "";
            
            # CORS headers - Add to all responses
            add_header 'Access-Control-Allow-Origin' $cors_allow_origin always;
            add_header 'Vary' 'Origin' always;
            add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS' always;
            add_header 'Access-Control-Allow-Headers' 'Accept, Accept-Language, Content-Language, Content-Type, Authorization, X-Requested-With, X-Real-IP, X-Forwarded-For, X-Forwarded-Proto, Cache-Control, DNT, User-Agent, If-Modified-Since, Range' always;
            add_header 'Access-Control-Allow-Credentials' 'true' always;
            add_header 'Access-Control-Expose-Headers' 'Content-Length, Content-Range, Content-Type, Date, Server, Transfer-Encoding' always;
        }

        # ============================================================================
        # Admin API Routes (Alternative) - /admin/* → Backend Port 8611 (root path)
        # ============================================================================
        # This handles routes like /admin/dashboard/stats (without /api in path)
        location /admin/ {
            # Handle OPTIONS preflight requests
            if ($request_method = 'OPTIONS') {
                add_header 'Access-Control-Allow-Origin' $cors_allow_origin always;
                add_header 'Vary' 'Origin' always;
                add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS' always;
                add_header 'Access-Control-Allow-Headers' 'Accept, Accept-Language, Content-Language, Content-Type, Authorization, X-Requested-With, X-Real-IP, X-Forwarded-For, X-Forwarded-Proto, Cache-Control, DNT, User-Agent, If-Modified-Since, Range' always;
                add_header 'Access-Control-Allow-Credentials' 'true' always;
                add_header 'Access-Control-Max-Age' '86400' always;
                add_header 'Content-Type' 'text/plain; charset=utf-8' always;
                add_header 'Content-Length' '0' always;
                return 204;
            }
            
            # Strip /admin prefix and proxy to backend root path
            rewrite ^/admin/(.*)$ /$1 break;
            proxy_pass http://greenride_admin_api_backend;
            
            # Proxy headers
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_set_header X-Forwarded-Host $host;
            proxy_set_header X-Forwarded-Port $server_port;
            proxy_redirect off;
            port_in_redirect off;
            
            # Timeout configuration
            proxy_connect_timeout 60s;
            proxy_send_timeout 60s;
            proxy_read_timeout 120s;
            proxy_buffering on;
            proxy_buffer_size 4k;
            proxy_buffers 8 4k;
            proxy_busy_buffers_size 8k;
            
            # HTTP/1.1 and connection keepalive
            proxy_http_version 1.1;
            proxy_set_header Connection "";
            
            # CORS headers - Add to all responses
            add_header 'Access-Control-Allow-Origin' $cors_allow_origin always;
            add_header 'Vary' 'Origin' always;
            add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS' always;
            add_header 'Access-Control-Allow-Headers' 'Accept, Accept-Language, Content-Language, Content-Type, Authorization, X-Requested-With, X-Real-IP, X-Forwarded-For, X-Forwarded-Proto, Cache-Control, DNT, User-Agent, If-Modified-Since, Range' always;
            add_header 'Access-Control-Allow-Credentials' 'true' always;
            add_header 'Access-Control-Expose-Headers' 'Content-Length, Content-Range, Content-Type, Date, Server, Transfer-Encoding' always;
        }
        
        # ============================================================================
        # Static Files Cache - Admin Frontend
        # ============================================================================
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
            proxy_pass http://greenride_admin_frontend_backend;
            expires 1y;
            add_header Cache-Control "public, immutable";
            add_header X-Served-By "nginx";
        }
        
        # ============================================================================
        # Admin Frontend Application - / → Admin Frontend (Port 3001)
        # ============================================================================
        location / {
            proxy_pass http://greenride_admin_frontend_backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_redirect off;
            port_in_redirect off;
            
            # HTTP/1.1 and connection keepalive
            proxy_http_version 1.1;
            proxy_set_header Connection "";
        }
    }
}
NGINX_EOF
```

**Then verify:**
```bash
cat ~/nginx.conf | head -20
```

---

## ✅ **Method 3: Using scp (from your local machine)**

**On your Windows machine (PowerShell):**

```powershell
# Navigate to project directory
cd D:\greenride-admin-v.2

# Copy file to cloudshell
scp nginx-config-fixed.conf ubuntu@18.143.118.157:~/nginx.conf
```

**Then on cloudshell, verify:**
```bash
cat ~/nginx.conf | head -20
```

---

## ✅ **Method 4: Using vi/vim (if you prefer)**

```bash
vi ~/nginx.conf
```

**Commands:**
1. Press `gg` (go to top)
2. Press `dG` (delete to end of file)
3. Press `i` (insert mode)
4. Paste content (right-click or `Shift+Insert`)
5. Press `Esc` (exit insert mode)
6. Type `:wq` (write and quit)
7. Press `Enter`

---

## 🎯 **After Replacing File - Test & Reload**

### **1. Test nginx config syntax**
```bash
# Find nginx container
docker ps | grep nginx

# Test config (replace <container-name>)
docker exec <container-name> nginx -t
```

**Expected output:**
```
nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
nginx: configuration file /etc/nginx/nginx.conf test is successful
```

### **2. Reload nginx**
```bash
# Reload (no downtime)
docker exec <container-name> nginx -s reload

# OR restart container
docker restart <container-name>
```

### **3. Verify it's working**
```bash
# Test Admin API endpoint
curl -I http://localhost/admin/api/health

# Check logs
docker logs <container-name> | tail -20
```

---

## ⚠️ **Troubleshooting**

### **Problem: Paste doesn't work in nano**
- Try `Ctrl + Shift + V` instead of `Ctrl + V`
- Or use Method 2 (cat with heredoc) - it's more reliable

### **Problem: File is read-only**
```bash
# Make sure you have write permissions
chmod 644 ~/nginx.conf
```

### **Problem: nginx test fails**
```bash
# Check error message
docker exec <container-name> nginx -t

# Common issues:
# - Missing semicolons
# - Incorrect brackets
# - Invalid directives
```

---

## 📝 **Quick Reference**

| Method | Speed | Reliability | Best For |
|--------|-------|-------------|----------|
| **nano** | Medium | High | Visual editing |
| **cat heredoc** | Fast | Very High | Quick replacement |
| **scp** | Fast | Very High | From local machine |
| **vi/vim** | Medium | High | If you know vim |

**Recommendation:** Use **Method 2 (cat heredoc)** - fastest and most reliable! 🚀
