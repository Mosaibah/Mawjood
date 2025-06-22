# 🎯 Mawjood - Content Discovery Platform

A microservices-based content management and discovery platform built with Go, gRPC, and CockroachDB.

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Discovery     │    │      CMS        │    │   CockroachDB   │
│   Service       │    │    Service      │    │    Database     │
│   (Port 9002)   │    │  (Port 9001)    │    │   (Port 26257)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │      Nginx      │
                    │   (SSL + gRPC)  │
                    │   Port 443      │
                    └─────────────────┘
```

## 🚀 Services

- **Discovery Service**: User-facing content search and retrieval
- **CMS Service**: Admin content management operations  
- **gRPC-UI**: Web interfaces for both gRPC services
- **CockroachDB**: Distributed SQL database with full-text search

## 📍 Endpoints

- **Main API**: `https://mawjood.mosaibah.com`
- **CMS gRPC UI**: `https://mawjood.mosaibah.com/cms/`
- **Discovery gRPC UI**: `https://mawjood.mosaibah.com/api/`
- **Health Check**: `https://mawjood.mosaibah.com/health`

## 🛠️ Tech Stack

- **Backend**: Go 1.24+ with gRPC
- **Database**: CockroachDB v25.2.1
- **Deployment**: Docker + Docker Compose
- **Proxy**: Nginx with SSL termination
- **Build**: Bazel for Go code generation

## 🔧 Local Development

### Prerequisites
- Go 1.24+
- Docker & Docker Compose
- [Task](https://taskfile.dev/) (optional but recommended)


**LIVE Endpoints:**
- Main API: `https://mawjood.mosaibah.com`
- CMS UI: `https://mawjood.mosaibah.com/cms/`
- Discovery UI: `https://mawjood.mosaibah.com/api/`


### Local Development

```bash
# Clone and start everything
git clone https://github.com/yourusername/Mawjood.git
cd Mawjood
task run  # or: docker-compose up -d
```

### Local Endpoints
- **Database Admin**: `http://localhost:8080`
- **CMS gRPC UI**: `http://localhost:8081` 
- **Discovery gRPC UI**: `http://localhost:8082`

### Common Tasks

```bash
task run      # Start everything
task logs     # View logs
task test     # Run tests
task restart  # Restart services
task stop     # Stop everything
task clean    # Clean up (removes data!)
```

### Manual Commands (if not using Task)

```bash
# Start
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down

# Run tests
go test ./packages/cms/... ./packages/discovery/...
```
## 📚 Documentation

- **[Deployment Guide](DEPLOYMENT-GUIDE.md)** - Complete deployment instructions
- **[Server Setup](SERVER-SETUP.md)** - Server preparation and configuration
- **[Project Structure](notes.md)** - Architecture and folder organization
