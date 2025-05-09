#!/bin/bash

# Wait for Redis nodes to be ready
echo "Waiting for Redis nodes to be ready..."
sleep 5

# Create the cluster
echo "Creating Redis cluster..."
docker exec -it redis-6379 redis-cli --cluster create \
  redis-6379:6379 \
  redis-6380:6379 \
  redis-6381:6379 \
  redis-6382:6379 \
  redis-6383:6379 \
  redis-6384:6379 \
  --cluster-replicas 1

echo "Redis cluster initialization completed!" 