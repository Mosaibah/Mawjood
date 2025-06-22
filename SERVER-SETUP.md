# 🖥️ Server Setup Guide

Quick setup guide for deploying Mawjood to `/var/www` on your Digital Ocean server.

## 🚀 **Step 1: Initial Server Setup**

```bash
# SSH into your server
ssh root@your-server-ip

# Update system
apt update && apt upgrade -y

# Install essential packages
apt install -y docker.io docker-compose nginx git ufw htop

# Start and enable Docker
systemctl start docker
systemctl enable docker

# Add user to docker group (if not root)
usermod -aG docker $USER
``

## 🔧 **Step 2: Configure Firewall**

```bash
# Configure UFW firewall
ufw allow ssh
ufw allow 80/tcp
ufw allow 443/tcp
ufw allow 444/tcp  # For Discovery gRPC
ufw --force enable

# Check status
ufw status
```

## 📁 **Step 3: Prepare /var/www Directory**

```bash
# Ensure /var/www exists and has proper permissions
mkdir -p /var/www
chown -R $USER:$USER /var/www
chmod 755 /var/www
```

## 📦 **Step 4: Clone and Deploy Mawjood**

```bash
# Clone project
git clone https://github.com/yourusername/Mawjood.git /var/www/mawjood
cd /var/www/mawjood

# Set ownership
chown -R $USER:$USER /var/www/mawjood

# Create environment file
cat > .env.production << EOF
DOCKER_USERNAME=your_dockerhub_username
DB_PASSWORD=$(openssl rand -base64 32)
DOMAIN=mawjood.mosaibah.com
EMAIL=your@email.com
EOF

# Deploy
chmod +x deploy-prod.sh
./deploy-prod.sh
```

## 🔒 **Step 5: Setup SSL with Certbot**

```bash
# Install Certbot
apt install -y certbot python3-certbot-nginx

# Get SSL certificates
certbot --nginx -d mawjood.mosaibah.com -d cms.mawjood.mosaibah.com

# Update nginx config for gRPC
cp nginx-mawjood-post-ssl.conf /etc/nginx/sites-available/mawjood
systemctl reload nginx
```

## ✅ **Step 6: Verify Deployment**

```bash
# Check Docker services
docker-compose -f docker-compose.prod.yml ps

# Test endpoints
curl https://mawjood.mosaibah.com/health
curl https://mawjood.mosaibah.com/

# Test gRPC (if grpcurl installed)
grpcurl mawjood.mosaibah.com:444 list
grpcurl cms.mawjood.mosaibah.com:443 list
```

## 🔄 **Future Updates**

```bash
# Navigate to project directory
cd /var/www/mawjood

# Pull latest changes
git pull origin main

# Rebuild and restart services
docker-compose -f docker-compose.prod.yml down
./deploy-prod.sh
```

## 📍 **File Locations**

- **Project**: `/var/www/mawjood/`
- **Nginx Config**: `/etc/nginx/sites-available/mawjood`
- **SSL Certs**: `/etc/letsencrypt/live/mawjood.mosaibah.com/`
- **Docker Data**: Docker volumes (managed by Docker)

## 🛠️ **Useful Commands**

```bash
# Check nginx status
systemctl status nginx
nginx -t

# Check SSL certificates
certbot certificates

# Monitor Docker logs
cd /var/www/mawjood
docker-compose -f docker-compose.prod.yml logs -f

# Check system resources
htop
df -h
docker system df
```

This setup ensures your Mawjood deployment is properly organized in `/var/www` with correct permissions and easy maintenance! 🎉 