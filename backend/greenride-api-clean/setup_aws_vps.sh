#!/bin/bash

# GreenRide VPS Setup Script (Ubuntu)
# This script installs Docker and Docker Compose

echo "🛠️ Starting VPS Setup..."

# 1. Update system
sudo apt-get update
sudo apt-get upgrade -y

# 2. Install Docker
echo "🐳 Installing Docker..."
sudo apt-get install -y ca-certificates curl gnupg lsb-release
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# 3. Enable Docker without sudo (optional but convenient)
sudo usermod -aG docker $USER

# 4. Success message
echo "✅ VPS Setup Complete!"
echo "Please log out and log back in to use docker without sudo."
docker --version
