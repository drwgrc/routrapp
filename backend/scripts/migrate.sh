#!/bin/bash

# Set environment variables if not already set
: ${DB_HOST:="localhost"}
: ${DB_PORT:="5432"}
: ${DB_USER:="postgres"}
: ${DB_PASSWORD:="postgres"}
: ${DB_NAME:="routrapp"}
: ${DB_SSL_MODE:="disable"}

# Create migrations directory if it doesn't exist
MIGRATIONS_DIR="$(dirname "$0")/../internal/repositories/postgres/migrations"

# Check if migrations directory exists
if [ ! -d "$MIGRATIONS_DIR" ]; then
    echo "Migrations directory not found: $MIGRATIONS_DIR"
    exit 1
fi

# Check if migrate command is available
if ! command -v migrate &> /dev/null; then
    echo "migrate command not found. Installing golang-migrate..."
    if [ "$(uname)" == "Darwin" ]; then
        # macOS
        brew install golang-migrate
    elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
        # Linux
        curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
        sudo mv migrate /usr/local/bin/
    else
        echo "Unsupported OS. Please install golang-migrate manually: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"
        exit 1
    fi
fi

# Build connection string
DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}"

# Parse command (up, down, create, etc.)
CMD=$1
shift

# Execute migration command
case "$CMD" in
    "up")
        echo "Running migrations up..."
        migrate -path "$MIGRATIONS_DIR" -database "$DB_URL" up "$@"
        ;;
    "down")
        echo "Running migrations down..."
        migrate -path "$MIGRATIONS_DIR" -database "$DB_URL" down "$@"
        ;;
    "create")
        echo "Creating migration..."
        migrate create -ext sql -dir "$MIGRATIONS_DIR" -seq "$@"
        ;;
    "version")
        echo "Current migration version:"
        migrate -path "$MIGRATIONS_DIR" -database "$DB_URL" version
        ;;
    "force")
        echo "Forcing migration version..."
        migrate -path "$MIGRATIONS_DIR" -database "$DB_URL" force "$@"
        ;;
    *)
        echo "Usage: $0 <up|down|create|version|force> [args]"
        echo "Examples:"
        echo "  $0 up             # Apply all migrations"
        echo "  $0 up 1           # Apply 1 migration"
        echo "  $0 down           # Revert last migration"
        echo "  $0 down 2         # Revert last 2 migrations"
        echo "  $0 create my_migration # Create a new migration"
        echo "  $0 version        # Show current migration version"
        echo "  $0 force 1        # Force migration version"
        exit 1
        ;;
esac 