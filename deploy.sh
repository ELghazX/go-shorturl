#!/bin/bash
set -e

echo "Deploying Go URL Shortener..."

echo "Pulling latest image..."
docker compose pull

echo "Stopping old containers..."
docker compose down

echo "Starting new containers..."
docker compose up -d

echo "Cleaning up..."
docker image prune -f

echo "Deployment completed!"
echo "Application is running at $BASE_URL"
