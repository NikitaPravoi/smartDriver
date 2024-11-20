#!/bin/bash

# env.sh - Environment management script

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

# Function to check if .env file exists
check_env_file() {
    if [ ! -f .env ]; then
        echo -e "${RED}Error: .env file not found${NC}"
        echo "Creating .env file from example..."
        cp .env.example .env
        echo -e "${GREEN}Created .env file. Please edit it with your configuration.${NC}"
        exit 1
    fi
}

# Function to validate environment variables
validate_env() {
    local required_vars=(
        "DB_USER"
        "DB_PASSWORD"
        "DB_NAME"
        "CENTRIFUGO_API_KEY"
        "CENTRIFUGO_TOKEN_SECRET"
    )

    local missing_vars=()

    for var in "${required_vars[@]}"; do
        if [ -z "${!var}" ]; then
            missing_vars+=("$var")
        fi
    done

    if [ ${#missing_vars[@]} -ne 0 ]; then
        echo -e "${RED}Error: Missing required environment variables:${NC}"
        printf '%s\n' "${missing_vars[@]}"
        exit 1
    fi
}

# Main script
case "$1" in
    "check")
        check_env_file
        source .env
        validate_env
        echo -e "${GREEN}Environment configuration is valid.${NC}"
        ;;
    "init")
        if [ ! -f .env ]; then
            cp .env.example .env
            echo -e "${GREEN}Created .env file from example. Please edit it with your configuration.${NC}"
        else
            echo -e "${RED}Error: .env file already exists${NC}"
        fi
        ;;
    *)
        echo "Usage: $0 {check|init}"
        exit 1
        ;;
esac

exit 0