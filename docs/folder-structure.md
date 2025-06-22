# Folder Structure & Module Boundaries

This document explains how our project folders are organized and defines clear boundaries between modules. The structure keeps our code clean, easy to understand, and scalable.

## Service Architecture Overview

Our system has two main services with distinct responsibilities:

### Service Responsibilities

- **Discovery Service**: Handles user-facing content search and discovery
- **CMS Service**: Manages content creation, updates, and admin operations  
- **Providers Service**: Integrates with external platforms (future implementation)

### Service Interactions

- **Users** → Discovery Service only
- **Admins** → CMS Service only
- Discovery and CMS services are independent (no direct communication)
- Both services share the same database (temporary - will be separated later)
- CMS service integrates with external APIs for content import

## Module Boundaries

### API Contract Module (`packages/proto`)

**Purpose**: Defines all gRPC service contracts and shared messages

**Boundaries**:
- **Input**: Protocol buffer definitions
- **Output**: Generated Go code for all services
- **Dependencies**: None (pure protocol definitions)
- **Consumers**: All services import generated code

**Key Files**:
- `discovery.proto`: User-facing search and content retrieval
- `cms.proto`: Admin content management operations
- `messages.proto`: Shared data structures and request/response types

**Interface**:
```proto
// Discovery Service API
service DiscoveryService {
  rpc SearchContents(SearchContentsRequest) returns (SearchContentsResponse);
  rpc ListContents(ListContentsRequest) returns (ListContentsResponse);
  rpc GetContent(GetContentRequest) returns (Content);
}

// CMS Service API  
service CMSService {
  rpc CreateContent(CreateContentRequest) returns (Content);
  rpc UpdateContent(UpdateContentRequest) returns (Content);
  rpc DeleteContent(DeleteContentRequest) returns (google.protobuf.Empty);
  rpc ListContents(ListContentsRequest) returns (ListContentsResponse);
  rpc ImportFromExternal(ImportRequest) returns (ImportResponse);
}
```

### Discovery Service Module (`packages/discovery`)

**Purpose**: Provides user-facing content search and discovery

**Boundaries**:
- **Input**: gRPC requests from users/frontend
- **Output**: Search results, content details, paginated lists
- **Database Access**: Read-only operations only
- **External Dependencies**: None
- **Internal Dependencies**: Proto definitions only

**Module Structure**:
- `v1/service.go`: Business logic and gRPC handler implementation
- `store/store.go`: Database read operations interface and implementation
- `server/server.go`: gRPC server setup and configuration
- `mock/mock.go`: Mock implementations for testing

**Key Constraints**:
- **No write operations** to database
- **No direct service-to-service calls**
- **No external API integrations**
- Must handle high read traffic efficiently

**Interface Contract**:
```go
type DiscoveryService struct {
    store store.ContentStoreInterface
}

type ContentStoreInterface interface {
    SearchContents(ctx context.Context, query string, limit int32, offset string) ([]*Content, string, error)
    ListContents(ctx context.Context, platformID int32, limit int32, offset string) ([]*Content, string, error)
    GetContent(ctx context.Context, contentID int32) (*Content, error)
}
```

### CMS Service Module (`packages/cms`)

**Purpose**: Manages all content operations for administrators

**Boundaries**:
- **Input**: gRPC requests from admin interfaces
- **Output**: Content CRUD responses, import results
- **Database Access**: Full read/write operations
- **External Dependencies**: External content provider APIs
- **Internal Dependencies**: Proto definitions only

**Module Structure**:
- `v1/service.go`: Business logic and gRPC handler implementation
- `store/store.go`: Database CRUD operations interface and implementation  
- `server/server.go`: gRPC server setup and configuration
- `mock/mock.go`: Mock implementations for testing

**Key Constraints**:
- **Full database access** for content management
- **No direct service-to-service calls**
- **Authorized for external API integrations**
- Must validate all input data thoroughly

**Interface Contract**:
```go
type CMSService struct {
    store store.ContentStoreInterface
}

type ContentStoreInterface interface {
    CreateContent(ctx context.Context, req *CreateContentRequest) (*Content, error)
    UpdateContent(ctx context.Context, req *UpdateContentRequest) (*Content, error)
    DeleteContent(ctx context.Context, contentID int32) error
    ListContents(ctx context.Context, limit int32, offset string) ([]*Content, string, error)
}
```

### Providers Service Module (`packages/providers`) - Future

**Purpose**: Integrates with external content platforms

**Boundaries**:
- **Input**: Import requests from CMS service
- **Output**: Standardized content data
- **Database Access**: None (stateless processing)
- **External Dependencies**: YouTube API, podcast platforms, etc.
- **Internal Dependencies**: Proto definitions, shared utilities

**Module Structure** (Planned):
- `v1/service.go`: Provider integration logic
- `youtube/client.go`: YouTube API integration
- `podcast/client.go`: Podcast platform integrations
- `server/server.go`: gRPC server setup

## Module Communication Rules

### Allowed Communication Patterns

1. **Client → Discovery Service**: Direct gRPC calls
2. **Admin → CMS Service**: Direct gRPC calls  
3. **CMS Service → External APIs**: HTTP/REST calls
4. **All Services → Database**: Direct database connections
5. **All Services → Proto Generated Code**: Import and use


## Directory Structure Details

### `packages/proto/v1/`
- **BUILD.bazel**: Bazel configuration for proto compilation
- **discovery.proto**: Discovery service API definition  
- **cms.proto**: CMS service API definition
- **messages.proto**: Shared message types

### `packages/discovery/`
- **v1/**: Core business logic
- **store/**: Database access layer (read-only)
- **server/**: gRPC server setup
- **mock/**: Testing utilities
- **BUILD.bazel**: Bazel build configuration

### `packages/cms/`
- **v1/**: Core business logic  
- **store/**: Database access layer (read/write)
- **server/**: gRPC server setup
- **mock/**: Testing utilities
- **BUILD.bazel**: Bazel build configuration

## Benefits of This Structure

1. **Clear Separation**: Each service has distinct responsibilities
2. **Independent Scaling**: Services can scale based on different load patterns  
3. **Security Isolation**: Admin operations are completely separate from user operations
4. **Future-Proof**: New services can be added without affecting existing ones
5. **Database Flexibility**: Services can eventually use separate databases
6. **Testing**: Each module can be tested independently with clear interfaces
