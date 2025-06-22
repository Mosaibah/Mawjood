# 🚀 Mawjood Production Deployment Guide

Simple deployment using Docker + System Nginx + Certbot

## 📋 Prerequisites

- Digital Ocean droplet (4GB+ recommended)
- Docker & Docker Compose installed
- Nginx installed
- Git installed
- Domain pointing to your server: `mawjood.mosaibah.com` & `cms.mawjood.mosaibah.com`

## 🔧 Setup

### 1. Build and Push Images (Local Machine)

```bash
# Create .env with your Docker Hub username
cat > .env << EOF
DOCKER_USERNAME=your_dockerhub_username
DB_PASSWORD=$(openssl rand -base64 32)
DOMAIN=mawjood.mosaibah.com
EMAIL=your@email.com
EOF

# Build and push images to Docker Hub
./build-and-push.sh
```

### 2. Deploy on Server

```bash
# Clone to /var/www (on server)
sudo git clone https://github.com/yourusername/Mawjood.git /var/www/mawjood

# Set proper ownership
sudo chown -R $USER:$USER /var/www/mawjood
cd /var/www/mawjood

# Create same .env file on server
cat > .env << EOF
DOCKER_USERNAME=your_dockerhub_username
DB_PASSWORD=your_secure_password_here
DOMAIN=mawjood.mosaibah.com
EMAIL=your@email.com
EOF

# Deploy (pulls images and configures everything)
chmod +x deploy-prod.sh
./deploy-prod.sh
```

### 3. Setup SSL Certificates

```bash
# Get SSL certificates
sudo certbot --nginx -d mawjood.mosaibah.com -d cms.mawjood.mosaibah.com
```

### 4. Update Nginx for gRPC

```bash
# Replace with SSL-enabled config
sudo cp nginx-mawjood-post-ssl.conf /etc/nginx/sites-available/mawjood
sudo systemctl reload nginx
```

## 🎯 Architecture

```
Internet → System Nginx (SSL) → Docker Services
```

**Files:**
- `docker-compose.prod.yml` - Docker services (no nginx)
- `nginx-mawjood.conf` - Initial HTTP config
- `nginx-mawjood-post-ssl.conf` - Final HTTPS + gRPC config

## 📍 Endpoints

- **Main Site:** `https://mawjood.mosaibah.com`
- **Admin UI:** `https://mawjood.mosaibah.com/admin/`
- **Health Check:** `https://mawjood.mosaibah.com/health`
- **CMS gRPC:** `cms.mawjood.mosaibah.com:443`
- **Discovery gRPC:** `mawjood.mosaibah.com:444`

## 🔄 Updates

```bash
# Update services
docker-compose -f docker-compose.prod.yml pull
docker-compose -f docker-compose.prod.yml up -d
```

## 🐛 Troubleshooting

```bash
# Check service status
docker-compose -f docker-compose.prod.yml ps

# View logs
docker-compose -f docker-compose.prod.yml logs -f

# Test nginx
sudo nginx -t

# Reload nginx
sudo systemctl reload nginx
```

## ✅ Benefits

- ✅ **Simple**: No Docker nginx complexity
- ✅ **Fast**: Direct service connections
- ✅ **Secure**: SSL handled by system nginx
- ✅ **Scalable**: Easy to add more services
- ✅ **Maintainable**: Clean separation of concerns 