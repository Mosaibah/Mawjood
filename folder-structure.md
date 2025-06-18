# Folder Structure Explanation

This document explains how our project's folders are organized. This structure is designed to keep our code clean, easy to understand, and scalable, based on the services defined in `system-design.md`.

## Service Architecture

Our system is composed of three main services, each with distinct responsibilities:

```mermaid
graph TD
    A[External Search API] -.-> B[CMS]
    C[User] --> D[Discovery]
    E[Admin] --> B
    B --> F[(Database)]
    D --> F
    
    style A fill:#e1f5fe
    style B fill:#f3e5f5
    style C fill:#e8f5e8
    style D fill:#fff3e0
    style E fill:#fce4ec
    style F fill:#f1f8e9
```

### Service Responsibilities

- **Discovery Service**: Handles user-facing content search and discovery
- **CMS Service**: Manages content creation, updates, and external API integrations
- **Providers Service**: Integrates with external platforms (YouTube, Podcast platforms) - future implementation

### Service Interactions

- Users interact only with the **Discovery** service
- Admins interact only with the **CMS** service  
- **Discovery** and **CMS** services are independent and don't communicate directly
- Both services share the same database (temporary design - will be separated later)
- **CMS** service can integrate with external APIs for content ingestion

## Main `packages` Directory

The `packages` directory holds all the major parts of our application. We have three main services: `discovery`, `cms`, and `providers`. We also have a special folder called `proto` that helps our services talk to each other.

### `packages/proto`

This folder acts as the single source of truth for our API contracts. It contains the `.proto` files that define our gRPC services and messages.

```proto
// packages/proto/mawjood/v1/discovery.proto

syntax = "proto3";

package mawjood.v1;

import "google/protobuf/timestamp.proto";

// Service for user-facing content discovery and search
service DiscoveryService {
  // Searches for content based on a query
  rpc SearchContents(SearchContentsRequest) returns (SearchContentsResponse);
  
  // Lists all content with pagination
  rpc ListContents(ListContentsRequest) returns (ListContentsResponse);
  
  // Gets a specific piece of content by ID
  rpc GetContent(GetContentRequest) returns (Content);
}
```

```proto
// packages/proto/mawjood/v1/cms.proto

syntax = "proto3";

package mawjood.v1;

import "google/protobuf/timestamp.proto";
import "google/protobuf/empty.proto";

// Service for content management (admin operations)
service CMSService {
  // Creates a new piece of content
  rpc CreateContent(CreateContentRequest) returns (Content);
  
  // Updates an existing piece of content
  rpc UpdateContent(UpdateContentRequest) returns (Content);
  
  // Deletes a piece of content
  rpc DeleteContent(DeleteContentRequest) returns (google.protobuf.Empty);
  
  // Lists all content for admin management
  rpc ListContents(ListContentsRequest) returns (ListContentsResponse);
  
  // Imports content from external sources
  rpc ImportFromExternal(ImportRequest) returns (ImportResponse);
}
```

-   **`mawjood/v1/discovery.proto`**: Defines the `DiscoveryService` for user-facing operations
-   **`mawjood/v1/cms.proto`**: Defines the `CMSService` for admin content management
-   **`mawjood/v1/providers.proto`**: Will define the `ProvidersService` for external integrations (future)
-   **`gen/go/`**: Bazel generates Go code from our `.proto` files

### `packages/discovery` (Discovery Service)

This service handles all user-facing content discovery and search operations.

-   **`v1/`**: Core business logic for discovery operations
    -   `service.go`: Implements the `DiscoveryServiceServer` interface

    ```go
    // packages/discovery/v1/service.go
    import discoverypbv1 "mawjood/gen/go/mawjood/v1"

    type DiscoveryService struct{
        discoverypbv1.UnimplementedDiscoveryServiceServer
        store store.ContentStoreInterface
    }

    func (s *DiscoveryService) SearchContents(ctx context.Context, req *discoverypbv1.SearchContentsRequest) (*discoverypbv1.SearchContentsResponse, error){
        // Search logic using CockroachDB's full-text search
        results, err := s.store.SearchContents(ctx, req.Query, req.PageSize, req.PageToken)
        // ...
    }
    ```
-   **`store/`**: Database operations for read-only content access
    -   `store.go`: Interface and implementation for content retrieval operations
-   **`server/`**: gRPC server setup for the discovery service
-   **`mock/`**: Mock implementations for testing
-   **`Dockerfile`**: Container configuration
-   **`BUILD.bazel`**: Bazel build configuration

### `packages/cms` (CMS Service)

This service handles all content management operations for administrators.

-   **`v1/`**: Core business logic for content management
    -   `service.go`: Implements the `CMSServiceServer` interface

    ```go
    // packages/cms/v1/service.go
    import cmspbv1 "mawjood/gen/go/mawjood/v1"

    type CMSService struct{
        cmspbv1.UnimplementedCMSServiceServer
        store store.ContentStoreInterface
        externalClient external.ExternalAPIClient
    }

    func (s *CMSService) CreateContent(ctx context.Context, req *cmspbv1.CreateContentRequest) (*cmspbv1.Content, error){
        // Content creation logic with validation
        newContent, err := s.store.CreateContent(ctx, req)
        // ...
    }
    ```
-   **`store/`**: Database operations for content CRUD operations
    -   `store.go`: Interface and implementation for content management
-   **`external/`**: External API integrations (YouTube, Podcast platforms)
    -   `client.go`: Client for external content providers
-   **`server/`**: gRPC server setup for the CMS service
-   **`mock/`**: Mock implementations for testing
-   **`Dockerfile`**: Container configuration
-   **`BUILD.bazel`**: Bazel build configuration

### `packages/providers` (Providers Service)

This service will handle integrations with external content platforms (future implementation).

-   **`v1/`**: Core business logic for provider integrations
-   **`youtube/`**: YouTube API integration
-   **`podcast/`**: Podcast platform integrations
-   **`server/`**: gRPC server setup
-   **`mock/`**: Mock implementations for testing
-   **`Dockerfile`**: Container configuration
-   **`BUILD.bazel`**: Bazel build configuration

## Database Schema

Both Discovery and CMS services share the same CockroachDB database (temporary design):

### Tables

#### `platforms`

This table stores information about content platforms (YouTube, Spotify, Apple Podcasts, etc.).

```sql
CREATE TABLE platforms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    api_endpoint VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### `contents`

This table stores the core content metadata with a reference to its platform.

```sql
CREATE TABLE contents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    language VARCHAR(50),
    duration_seconds INT,
    published_at TIMESTAMPTZ,
    content_type VARCHAR(20) NOT NULL,
    platform_id UUID REFERENCES platforms(id),
    external_id VARCHAR(255), -- ID from the external platform (e.g., YouTube video ID)
    external_url VARCHAR(500), -- Direct URL to the content on the platform
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### `tags`

This table stores unique tags.

```sql
CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE
);
```

#### `content_tags`

This is a join table to associate content with tags.

```sql
CREATE TABLE content_tags (
    content_id UUID NOT NULL REFERENCES contents(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (content_id, tag_id)
);
```

### Indexes

To support efficient querying for the search functionality:

```sql
-- Index for searching content by title
CREATE INDEX ON contents (title);

-- For full-text search, CockroachDB supports trigram indexes.
-- This can be used to implement the SearchContents RPC.
CREATE INVERTED INDEX ON contents (title gin_trgm_ops);
CREATE INVERTED INDEX ON contents (description gin_trgm_ops);

-- Index for finding tags by name
CREATE INDEX ON tags (name);

-- Index for finding platforms by name
CREATE INDEX ON platforms (name);

-- Indexes on the join table for efficient lookups in both directions
CREATE INDEX ON content_tags (tag_id);

-- Index for content by platform
CREATE INDEX ON contents (platform_id);

-- Index for external ID lookups (useful for avoiding duplicates)
CREATE INDEX ON contents (platform_id, external_id);
```

### Sample Platform Data

```sql
-- Insert some common platforms
INSERT INTO platforms (name, description, api_endpoint) VALUES
('YouTube', 'YouTube video platform', 'https://www.googleapis.com/youtube/v3'),
('Spotify', 'Spotify podcast platform', 'https://api.spotify.com/v1'),
('Apple Podcasts', 'Apple Podcasts platform', 'https://itunes.apple.com'),
('SoundCloud', 'SoundCloud audio platform', 'https://api.soundcloud.com'),
('Internal', 'Internally created content', NULL);
```

## Why this structure?

1.  **Clear Separation**: Discovery and CMS services have distinct responsibilities and don't interfere with each other
2.  **Scalability**: Each service can be scaled independently based on load
3.  **Security**: Admin operations are isolated in the CMS service
4.  **Future-Proof**: Providers service can be added later without affecting existing services
5.  **Database Independence**: Services can eventually use separate databases for better isolation
