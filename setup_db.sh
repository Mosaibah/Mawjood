#!/bin/bash

# Database setup script for Mawjood
echo "Setting up Mawjood database..."

# Check if container is running
if ! docker ps | grep -q roach1; then
    echo "Starting CockroachDB container..."
    docker run -d --name=roach1 --hostname=roach1 \
        -p 26257:26257 -p 8080:8080 \
        -v "roach1:/cockroach/cockroach-data" \
        cockroachdb/cockroach:v25.2.1 start-single-node --insecure
    
    # Wait for container to be ready
    sleep 5
fi

# Run the database setup
echo "Running database setup..."
docker exec -i roach1 cockroach sql --insecure < db_setup.sql

echo "Database setup complete!"
echo "Web UI available at: http://localhost:8080"
echo "Connection: localhost:26257 (database: mawjood)" 