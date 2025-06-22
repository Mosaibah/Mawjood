#!/bin/bash
set -e

echo "🚀 Deploying Mawjood to Production..."
echo "📍 Working directory: $(pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if .env exists
if [ ! -f .env ]; then
    echo -e "${RED}❌ .env file not found!${NC}"
    echo "Create .env with:"
    echo "DOCKER_USERNAME=your_dockerhub_username"
    echo "DB_PASSWORD=\$(openssl rand -base64 32)"
    echo "DOMAIN=mawjood.mosaibah.com"
    echo "EMAIL=your@email.com"
    exit 1
fi

# Load environment variables
source .env

# Export variables for Docker Compose
export DOCKER_USERNAME
export DB_PASSWORD
export DOMAIN
export EMAIL

echo -e "${YELLOW}📦 Pulling Docker images from registry...${NC}"

# Pull pre-built images from registry (much faster than building on server)
docker pull ${DOCKER_USERNAME}/mawjood-cms:latest
docker pull ${DOCKER_USERNAME}/mawjood-discovery:latest

echo -e "${GREEN}✅ Docker images pulled successfully${NC}"
echo -e "${YELLOW}💡 To build and push new images, run locally:${NC}"
echo -e "   docker build -f Dockerfile.cms -t ${DOCKER_USERNAME}/mawjood-cms:latest ."
echo -e "   docker build -f Dockerfile.discovery -t ${DOCKER_USERNAME}/mawjood-discovery:latest ."
echo -e "   docker push ${DOCKER_USERNAME}/mawjood-cms:latest"
echo -e "   docker push ${DOCKER_USERNAME}/mawjood-discovery:latest"

echo -e "${YELLOW}🐳 Starting Docker services...${NC}"

# Start Docker services
docker-compose -f docker-compose.prod.yml up -d

echo -e "${GREEN}✅ Docker services started${NC}"

echo -e "${YELLOW}🔐 Setting up database authentication...${NC}"

# Wait for CockroachDB to be ready
sleep 15

# Set root password in database to match .env file
echo -e "${YELLOW}Setting root password in CockroachDB...${NC}"
docker-compose -f docker-compose.prod.yml exec -T cockroachdb cockroach sql --certs-dir=/cockroach/certs --host=localhost:26257 --execute="ALTER USER root WITH PASSWORD '$DB_PASSWORD';" || echo -e "${YELLOW}⚠️  Password may already be set${NC}"

# Restart services to ensure they connect with the password
echo -e "${YELLOW}Restarting application services...${NC}"
docker-compose -f docker-compose.prod.yml restart cms-service discovery-service

echo -e "${GREEN}✅ Database authentication configured${NC}"

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