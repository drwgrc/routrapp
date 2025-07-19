#!/bin/bash

# Database Migration Script
# A wrapper for the Go migration tool

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo -e "${BLUE}🗄️  Database Migration Tool${NC}"
echo "================================"

# Check if we're in the backend directory
if [[ ! -f "$PROJECT_ROOT/go.mod" ]]; then
    echo -e "${RED}Error: This script must be run from the backend directory${NC}"
    echo "Current directory: $(pwd)"
    echo "Expected: $PROJECT_ROOT"
    exit 1
fi

# Change to project root
cd "$PROJECT_ROOT"

# If no arguments provided, show usage
if [[ $# -eq 0 ]]; then
    echo -e "${YELLOW}Usage:${NC}"
    echo "  $0 create <migration_name>    # Create new migration"
    echo "  $0 validate                   # Validate all migrations"
    echo "  $0 status                     # Show migration status"
    echo ""
    echo -e "${YELLOW}Examples:${NC}"
    echo "  $0 create \"add user sessions\""
    echo "  $0 validate"
    echo "  $0 status"
    echo ""
    exit 0
fi

# Parse command
COMMAND="$1"
shift

case "$COMMAND" in
    "create")
        if [[ $# -eq 0 ]]; then
            echo -e "${RED}Error: Migration name is required${NC}"
            echo "Usage: $0 create <migration_name>"
            exit 1
        fi
        MIGRATION_NAME="$1"
        echo -e "${GREEN}Creating migration: $MIGRATION_NAME${NC}"
        go run scripts/migrate.go -action=create -name="$MIGRATION_NAME"
        ;;
    
    "validate")
        echo -e "${GREEN}Validating migrations...${NC}"
        go run scripts/migrate.go -action=validate
        ;;
    
    "status")
        echo -e "${GREEN}Checking migration status...${NC}"
        go run scripts/migrate.go -action=status
        ;;
    
    "help"|"-h"|"--help")
        go run scripts/migrate.go
        ;;
    
    *)
        echo -e "${RED}Error: Unknown command '$COMMAND'${NC}"
        echo "Run '$0' without arguments to see usage information"
        exit 1
        ;;
esac

echo -e "${GREEN}✅ Done!${NC}" 