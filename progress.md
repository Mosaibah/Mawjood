# Project Progress Tracker

This document tracks the development progress of the Mawjood project, breaking it down into small, manageable steps.

## Phase 1: Project Foundation & API Definition
- [X] **Status:** `Completed`
- **Goal:** Set up the basic monorepo structure, define the core API contracts, and configure the build system.

- [X] Set up the monorepo with the `packages` directory.
- [X] Initialize Bazel as the build system within the workspace.
- [X] Define the `DiscoveryService` API contract in `packages/proto/mawjood/v1/discovery.proto`.
- [X] Define the `CMSService` API contract in `packages/proto/mawjood/v1/cms.proto`.
- [ ] Configure Bazel to generate Go gRPC code from the `.proto` files into `gen/go`.
- [ ] Test the proto generation to ensure the build system is working correctly.

## Phase 2: Database Setup
- [ ] **Status:** `Not Started`
- **Goal:** Establish the database schema and populate it with initial data.

- [ ] Set up a local CockroachDB instance using Docker.
- [ ] Write the SQL script to create the `platforms`, `contents`, `tags`, and `content_tags` tables.
- [ ] Apply the database schema.
- [ ] Write the SQL script to create all necessary indexes for searching and relationships.
- [ ] Apply the indexes.
- [ ] Write a seeding script to insert the initial platform data (YouTube, Spotify, etc.).
- [ ] Run the seed script.

## Phase 3: CMS Service Implementation
- [ ] **Status:** `Not Started`
- **Goal:** Build the service responsible for content management.

- [ ] Create the basic directory structure for the `packages/cms` service.
- [ ] Implement the `store` package for database CRUD operations on content.
- [ ] Implement the `CreateContent` gRPC endpoint.
- [ ] Implement the `UpdateContent` gRPC endpoint.
- [ ] Implement the `DeleteContent` gRPC endpoint.
- [ ] Implement the `ListContents` gRPC endpoint for the admin view.
- [ ] Set up the main gRPC server for the CMS service.
- [ ] Create a `Dockerfile` to containerize the CMS service.
- [ ] Write unit tests for the `store` package.
- [ ] Write integration tests for the gRPC endpoints.
- [ ] **(Stretch Goal)** Implement the client for an external API (e.g., YouTube).
- [ ] **(Stretch Goal)** Implement the `ImportFromExternal` gRPC endpoint.

## Phase 4: Discovery Service Implementation
- [ ] **Status:** `Not Started`
- **Goal:** Build the user-facing service for searching and discovering content.

- [ ] Create the basic directory structure for the `packages/discovery` service.
- [ ] Implement the `store` package for read-only database access.
- [ ] Implement the `GetContent` gRPC endpoint.
- [ ] Implement the `ListContents` gRPC endpoint with pagination.
- [ ] Implement the `SearchContents` gRPC endpoint using full-text search capabilities.
- [ ] Set up the main gRPC server for the Discovery service.
- [ ] Create a `Dockerfile` to containerize the Discovery service.
- [ ] Write unit tests for the `store` package.
- [ ] Write integration tests for the gRPC endpoints.

## Phase 5: Providers Service (Future)
- [ ] **Status:** `Not Started`
- **Goal:** Design and implement the service for integrating with external content providers. This will be tackled after the core services are complete.

- [ ] Define the `ProvidersService` in `packages/proto/mawjood/v1/providers.proto`.
- [ ] Plan and implement the service logic.
