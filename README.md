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

### Demo Screenshots

*CMS Service gRPC UI Interface:*
<!-- Add CMS gRPC UI screenshot here -->

*Discovery Service gRPC UI Interface:*
<!-- Add Discovery gRPC UI screenshot here -->

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

## 🧪 Testing

We implement comprehensive unit tests for all core services. Here are the test results:

```bash
🧪 Running all tests...
task: [test] go test ./packages/cms/...
?       github.com/mosaibah/Mawjood/packages/cms/mock   [no test files]
?       github.com/mosaibah/Mawjood/packages/cms/server [no test files]
ok      github.com/mosaibah/Mawjood/packages/cms/store  0.556s
ok      github.com/mosaibah/Mawjood/packages/cms/v1     0.790s
task: [test] go test ./packages/discovery/...
?       github.com/mosaibah/Mawjood/packages/discovery/mock     [no test files]
?       github.com/mosaibah/Mawjood/packages/discovery/server   [no test files]
ok      github.com/mosaibah/Mawjood/packages/discovery/store    0.303s
ok      github.com/mosaibah/Mawjood/packages/discovery/v1       0.543s
✅ All tests completed!
```

## 🔍 How Search Works

We use **trigram matching** to help users find content quickly and accurately.

### How Our Search Works

1. **Setup**: We use PostgreSQL's `pg_trgm` extension to enable trigram matching
2. **Similarity**: We set the similarity threshold to 0.10 (10% match required)
3. **Matching**: The system compares trigrams between the search query and content

## 📄 Pagination

We use **Keyset pagination** for efficient data retrieval. This approach is more efficient than offset pagination, especially for large datasets.

```sql
WHERE deleted_at IS NULL
ORDER BY created_at DESC 
LIMIT $1
```

**Benefits:**
- Consistent performance regardless of page depth

## 🛠️ Tools Reasoning

**gRPC:**
- Transfers data in binary format, much faster than JSON (7x) !
- Helps serve millions of users efficiently

**CockroachDB:**
- Distributed SQL database for easy scaling
- Stores data across multiple nodes for reliability

**Bazel:**
- Build tool that manages dependencies and builds
- Ensures consistent builds across environments

**ClickHouse (not implemented):**
- Ideal for the Discovery service's read-heavy workload
- Couldn't implement due to CDC licensing requirements

## 🚧 Challenges Encountered

**Bazel Configuration:**
- Took 70% of implementation time due to steep learning curve
- Major changes from previous Bazel versions required relearning

**Server Deployment:**
- Deployed project with Docker containers on production server
- Configured Nginx reverse proxy for gRPC services
- Set up gRPC-UI for easier API testing and debugging

**Protocol Buffer Generation:**
- Multiple issues with Bazel and proto file compilation
- Had to make compromises in build configuration

## 🚀 Future Improvements

**Providers Service:**
- Connect to external content providers (YouTube, Spotify, etc.)
- Users share links, system automatically extracts metadata
- Store title, description, tags, and other content details

**ClickHouse Integration:**
- Use ClickHouse for Discovery service read operations
- Separate read/write workloads for better performance
- Couldn't implement due to CDC licensing requirements

**Architecture Enhancement:**
- Discovery service queries ClickHouse only
- Remove CockroachDB dependency for read operations
- Better separation between operational and analytical workloads

**Infrastructure Organization:**
- Create dedicated `infra` folder for deployment configs
- Centralize Docker, Kubernetes, and server configurations
- Currently avoided due to extensive path updates required
- k8s :)
- CI/CD pipeline (very simple and straight forward since I'm already using docker registry)

**Random Notes:**
- I intenlty exposed the grpcui for MVP purpuses

## 🏗️ Folder Structure & Module Boundaries


### Service Architecture

- **Discovery Service**: User-facing content search and discovery (read-only)
- **CMS Service**: Admin content management operations (full CRUD)
- **Providers Service**: External platform integrations (future)

**Service Interactions**:
- Users → Discovery Service only
- Admins → CMS Service only  
- Services are independent with no direct communication
- Both services currently share the same database

### Module Structure

#### `packages/proto/` - API Contracts
Defines all gRPC service contracts and shared messages. Generated Go code is used by all services.

**Key APIs**:
```proto
service DiscoveryService {
  rpc SearchContents(SearchContentsRequest) returns (SearchContentsResponse);
  rpc ListContents(ListContentsRequest) returns (ListContentsResponse);
  rpc GetContent(GetContentRequest) returns (Content);
}

service CMSService {
  rpc CreateContent(CreateContentRequest) returns (Content);
  rpc UpdateContent(UpdateContentRequest) returns (Content);
  rpc DeleteContent(DeleteContentRequest) returns (google.protobuf.Empty);
  rpc ListContents(ListContentsRequest) returns (ListContentsResponse);
}
```

#### `packages/discovery/` - Discovery Service
**Purpose**: User-facing content search and discovery
- **Database Access**: Read-only operations
- **Constraints**: No write operations, no external API calls
- **Structure**: `v1/` (business logic), `store/` (database layer), `server/` (gRPC setup), `mock/` (testing)

#### `packages/cms/` - CMS Service  
**Purpose**: Admin content management operations
- **Database Access**: Full read/write operations
- **Capabilities**: CRUD operations, external API integrations
- **Structure**: `v1/` (business logic), `store/` (database layer), `server/` (gRPC setup), `mock/` (testing)

#### `packages/providers/` - Future Service
**Purpose**: External platform integrations (YouTube, Spotify, etc.)
- **Database Access**: None (stateless processing)
- **Planned Structure**: `v1/` (integration logic), `youtube/` (YouTube client), `podcast/` (podcast clients)

### Communication Rules

1. **Client → Discovery Service**: Direct gRPC calls
2. **Admin → CMS Service**: Direct gRPC calls
3. **CMS Service → External APIs**: HTTP/REST calls
4. **All Services → Database**: Direct connections
5. **All Services → Proto Code**: Import generated code

### Benefits

- **Clear Separation**: Distinct service responsibilities
- **Independent Scaling**: Services scale based on different load patterns
- **Security Isolation**: Admin operations separate from user operations
- **Future-Proof**: Easy to add new services
- **Testing**: Independent module testing with clear interfaces

## 📊 Sample Data
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