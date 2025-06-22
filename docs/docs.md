## How Search Works

We use **trigram matching** to help users find content quickly and accurately.

### How Our Search Works

1. **Setup**: We use PostgreSQL's `pg_trgm` extension to enable trigram matching
2. **Similarity**: We set the similarity threshold to 0.10 (10% match required)
3. **Matching**: The system compares trigrams between the search query and content

## Pagination

We use **Keyset pagination** for efficient data retrieval. This approach is more efficient than offset pagination, especially for large datasets.

```sql
WHERE deleted_at IS NULL
ORDER BY created_at DESC 
LIMIT $1
```

**Benefits:**
- Consistent performance regardless of page depth

## Tools Reasoning

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

## Challenges Encountered

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

## Future Improvements

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


**Random Notes:**
- I intenlty exposed the grpcui for MVP purpuses

