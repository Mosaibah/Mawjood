#!/bin/bash
set -e

echo "🔨 Building and Pushing Mawjood Images..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if .env exists
if [ ! -f .env ]; then
    echo -e "${RED}❌ .env file not found!${NC}"
    echo "Create .env with DOCKER_USERNAME=your_dockerhub_username"
    exit 1
fi

# Load environment variables
source .env

# Check if DOCKER_USERNAME is set
if [ -z "$DOCKER_USERNAME" ]; then
    echo -e "${RED}❌ DOCKER_USERNAME not set in .env file!${NC}"
    exit 1
fi

echo -e "${YELLOW}📦 Building Docker images...${NC}"

# Build images
docker build -f Dockerfile.cms -t ${DOCKER_USERNAME}/mawjood-cms:latest .
docker build -f Dockerfile.discovery -t ${DOCKER_USERNAME}/mawjood-discovery:latest .

echo -e "${GREEN}✅ Images built successfully${NC}"

echo -e "${YELLOW}📤 Pushing to Docker Hub...${NC}"

# Push images
docker push ${DOCKER_USERNAME}/mawjood-cms:latest
docker push ${DOCKER_USERNAME}/mawjood-discovery:latest

echo -e "${GREEN}✅ Images pushed successfully${NC}"
echo ""
echo -e "${YELLOW}🚀 Ready for deployment! Run on server:${NC}"
echo -e "   ./deploy-prod.sh" 