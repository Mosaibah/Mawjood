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

echo -e "${YELLOW}📦 Building multi-platform Docker images...${NC}"
echo -e "${YELLOW}💡 Building for both ARM64 (Mac) and AMD64 (Linux servers)${NC}"

# Create and use buildx builder for multi-platform builds
docker buildx create --name mawjood-builder --use 2>/dev/null || docker buildx use mawjood-builder

# Build and push multi-platform images directly
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -f Dockerfile.cms \
  -t ${DOCKER_USERNAME}/mawjood-cms:latest \
  --push .

docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -f Dockerfile.discovery \
  -t ${DOCKER_USERNAME}/mawjood-discovery:latest \
  --push .

echo -e "${GREEN}✅ Multi-platform images built and pushed successfully${NC}"

echo -e "${GREEN}✅ Images pushed successfully${NC}"
echo ""
echo -e "${YELLOW}🚀 Ready for deployment! Run on server:${NC}"
echo -e "   ./deploy-prod.sh" 