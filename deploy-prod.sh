#!/bin/bash
set -e

echo "🚀 Deploying Mawjood to Production..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if .env.production exists
if [ ! -f .env.production ]; then
    echo -e "${RED}❌ .env.production file not found!${NC}"
    echo "Create .env.production with:"
    echo "DOCKER_USERNAME=your_dockerhub_username"
    echo "DB_PASSWORD=your_secure_password"
    echo "DOMAIN=mawjood.mosaibah.com"
    echo "EMAIL=your@email.com"
    exit 1
fi

# Load environment variables
source .env.production

echo -e "${YELLOW}📦 Building and pushing Docker images...${NC}"

# Build and push images
docker build -f Dockerfile.cms -t ${DOCKER_USERNAME}/mawjood-cms:latest .
docker build -f Dockerfile.discovery -t ${DOCKER_USERNAME}/mawjood-discovery:latest .

docker push ${DOCKER_USERNAME}/mawjood-cms:latest
docker push ${DOCKER_USERNAME}/mawjood-discovery:latest

echo -e "${GREEN}✅ Docker images pushed successfully${NC}"

echo -e "${YELLOW}🐳 Starting Docker services...${NC}"

# Start Docker services
docker-compose -f docker-compose.prod.yml up -d

echo -e "${GREEN}✅ Docker services started${NC}"

echo -e "${YELLOW}🌐 Setting up nginx configuration...${NC}"

# Copy nginx configuration
sudo cp nginx-mawjood.conf /etc/nginx/sites-available/mawjood

# Enable the site
sudo ln -sf /etc/nginx/sites-available/mawjood /etc/nginx/sites-enabled/

# Test nginx configuration
if sudo nginx -t; then
    echo -e "${GREEN}✅ Nginx configuration is valid${NC}"
    sudo systemctl reload nginx
else
    echo -e "${RED}❌ Nginx configuration error!${NC}"
    exit 1
fi

echo -e "${YELLOW}🔒 Setting up SSL certificates...${NC}"
echo -e "${YELLOW}Run this command to get SSL certificates:${NC}"
echo -e "${GREEN}sudo certbot --nginx -d ${DOMAIN} -d cms.${DOMAIN}${NC}"
echo ""
echo -e "${YELLOW}After SSL setup, replace nginx config with:${NC}"
echo -e "${GREEN}sudo cp nginx-mawjood-post-ssl.conf /etc/nginx/sites-available/mawjood${NC}"
echo -e "${GREEN}sudo systemctl reload nginx${NC}"

echo ""
echo -e "${GREEN}🎉 Deployment completed!${NC}"
echo ""
echo -e "${YELLOW}📍 Your services will be available at:${NC}"
echo -e "• Main site: https://${DOMAIN}"
echo -e "• Admin UI: https://${DOMAIN}/admin/"
echo -e "• Health check: https://${DOMAIN}/health"
echo -e "• CMS gRPC: cms.${DOMAIN}:443"
echo -e "• Discovery gRPC: ${DOMAIN}:444"
echo ""
echo -e "${YELLOW}⚡ Next steps:${NC}"
echo -e "1. Run: ${GREEN}sudo certbot --nginx -d ${DOMAIN} -d cms.${DOMAIN}${NC}"
echo -e "2. Run: ${GREEN}sudo cp nginx-mawjood-post-ssl.conf /etc/nginx/sites-available/mawjood${NC}"
echo -e "3. Run: ${GREEN}sudo systemctl reload nginx${NC}"
echo -e "4. Test your endpoints!" 