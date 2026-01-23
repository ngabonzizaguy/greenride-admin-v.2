#!/bin/bash

# 🔍 Feedback Endpoints 404 Debug Script
# Run this on the server: bash debug_feedback_endpoints.sh

echo "🔍 Checking Backend Status..."
echo ""

# 1. Check if backend is running on port 8611
echo "1️⃣ Checking if backend is running on port 8611..."
if netstat -tlnp 2>/dev/null | grep -q ":8611" || ss -tlnp 2>/dev/null | grep -q ":8611"; then
    echo "✅ Backend is running on port 8611"
    netstat -tlnp 2>/dev/null | grep ":8611" || ss -tlnp 2>/dev/null | grep ":8611"
else
    echo "❌ Backend is NOT running on port 8611"
fi

echo ""
echo "2️⃣ Testing backend /health endpoint directly..."
curl -s -o /dev/null -w "HTTP Status: %{http_code}\n" http://localhost:8611/health || echo "❌ Cannot connect to backend"

echo ""
echo "3️⃣ Testing feedback/stats endpoint directly..."
STATS_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" http://localhost:8611/feedback/stats 2>&1)
HTTP_STATUS=$(echo "$STATS_RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
BODY=$(echo "$STATS_RESPONSE" | sed '/HTTP_STATUS/d')
echo "HTTP Status: $HTTP_STATUS"
echo "Response: $BODY" | head -c 200
echo ""

echo ""
echo "4️⃣ Testing through nginx..."
NGINX_STATS=$(curl -s -w "\nHTTP_STATUS:%{http_code}" http://localhost/admin/api/feedback/stats 2>&1)
NGINX_HTTP_STATUS=$(echo "$NGINX_STATS" | grep "HTTP_STATUS" | cut -d: -f2)
NGINX_BODY=$(echo "$NGINX_STATS" | sed '/HTTP_STATUS/d')
echo "HTTP Status: $NGINX_HTTP_STATUS"
echo "Response: $NGINX_BODY" | head -c 200
echo ""

echo ""
echo "5️⃣ Checking for Docker containers..."
if command -v docker &> /dev/null; then
    echo "Docker containers:"
    docker ps | grep -E "greenride|backend|api" || echo "No GreenRide containers found"
    
    echo ""
    echo "Checking if backend is in Docker..."
    BACKEND_CONTAINER=$(docker ps | grep -i "greenride.*api\|backend.*admin" | awk '{print $1}' | head -1)
    if [ -n "$BACKEND_CONTAINER" ]; then
        echo "✅ Found backend container: $BACKEND_CONTAINER"
        echo "Container logs (last 20 lines with 'feedback' or 'GET\|POST'):"
        docker logs --tail 20 "$BACKEND_CONTAINER" 2>&1 | grep -i "feedback\|GET.*feedback\|POST.*feedback" || echo "No feedback-related logs found"
    else
        echo "ℹ️  Backend not running in Docker (might be systemd or manual)"
    fi
fi

echo ""
echo "6️⃣ Checking systemd services..."
if command -v systemctl &> /dev/null; then
    systemctl list-units --type=service --state=running | grep -i "greenride\|backend" || echo "No GreenRide systemd services found"
fi

echo ""
echo "7️⃣ Checking backend process..."
ps aux | grep -E "greenride|backend.*8611|main.*8611" | grep -v grep || echo "No backend process found"

echo ""
echo "✅ Debug complete!"
echo ""
echo "📋 Next Steps:"
if [ "$HTTP_STATUS" != "200" ] && [ "$HTTP_STATUS" != "401" ]; then
    echo "⚠️  Backend endpoint returns $HTTP_STATUS - backend might need restart"
    echo "   To restart:"
    echo "   - If Docker: docker restart <container-name>"
    echo "   - If systemd: sudo systemctl restart greenride-admin"
    echo "   - If manual: kill process and restart"
fi
