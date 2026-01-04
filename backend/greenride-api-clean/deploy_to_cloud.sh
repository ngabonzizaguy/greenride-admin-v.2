#!/bin/bash

# GreenRide Cloud Deployment Script

echo "🚀 Starting GreenRide Deployment..."

# 1. Pull latest changes (assuming git is set up)
# git pull origin main

# 2. Build and start containers
echo "📦 Building and starting containers..."
docker compose up -d --build

# 3. Wait for DB to be ready
echo "⏳ Waiting for services to stabilize..."
sleep 10

echo "✅ Deployment Complete!"
echo "User API: http://YOUR_SERVER_IP:8610"
echo "Admin API: http://YOUR_SERVER_IP:8611"
