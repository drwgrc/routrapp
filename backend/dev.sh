#!/bin/sh

# Create tmp directory if it doesn't exist
mkdir -p /app/tmp

echo "Starting development server with hot reload..."

# Loop to rebuild and restart on changes
while true; do
  echo "Building..."
  
  # Build and run the application
  go build -o /app/tmp/main /app/cmd/api && /app/tmp/main &
  PID=$!
  
  echo "Server started with PID: $PID"
  
  # Watch for changes to Go files
  inotifyd - /app:c | while read EVENTS; do
    if echo "$EVENTS" | grep -q "\.go$"; then
      echo "Change detected, rebuilding..."
      kill $PID
      wait $PID
      break
    fi
  done
done 