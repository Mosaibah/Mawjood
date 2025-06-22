# 🎯 Mawjood - Content Discovery Platform

A microservices-based content management and discovery platform built with Go, gRPC, and CockroachDB.

## 🚀 Services

- **Discovery Service**: User-facing content search and retrieval
- **CMS Service**: Admin content management operations  
- **gRPC-UI**: Web interfaces for both gRPC services
- **CockroachDB**: Distributed SQL database with full-text search

## 📍 Live Demo Endpoints

- **CMS gRPC UI**: `https://mawjood.mosaibah.com/cms/`
- **Discovery gRPC UI**: `https://mawjood.mosaibah.com/api/`
- **Health Check**: `https://mawjood.mosaibah.com/health`

## 🛠️ Tech Stack

- **GoLang** 
- **CockroachDB** 
- **Bazel** 
- **gRPC** 
- **Nginx** 
- **DigitalOcean** 

## 🔧 Local Development

### Prerequisites
- Go 1.24+
- Docker & Docker Compose
- [Task](https://taskfile.dev/) (optional but recommended)


### Local Development

```bash
git clone https://github.com/mosaibah/mawjood.git
cd mawjood
task run  
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

## 📚 Documentation

- **[General Documentation](docs/docs.md)** - Complete deployment instructions
- **[Project Structure](docs/folder-structure)** - Architecture and folder organization

## Sample Data
Here is a sample request body for creating a content:
```json
{
  "tags": [],
  "title": "test",
  "description": "tewt",
  "language": "ar",
  "durationSeconds": 444,
  "publishedAt": "2006-01-02T15:04:05Z",
  "contentType": "CONTENT_TYPE_DOCUMENTARY",
  "url": "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG",
  "platformName": "Youtube"
}
```