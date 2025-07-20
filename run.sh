#!/bin/bash

# Simple script to run either development or production mode

MODE=${1:-dev}

case $MODE in
  dev)
    echo "Starting development environment..."
    docker-compose up --build
    ;;
  prod)
    echo "Starting production environment..."
    docker-compose -f docker-compose.prod.yml up --build
    ;;
  down)
    echo "Stopping containers..."
    docker-compose down
    ;;
  prod-down)
    echo "Stopping production containers..."
    docker-compose -f docker-compose.prod.yml down
    ;;
  *)
    echo "Usage: $0 {dev|prod|down|prod-down}"
    exit 1
    ;;
esac 