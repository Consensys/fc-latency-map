#!/bin/bash

# Export host UID and GID
export FC_UID=$(id -u)
export FC_GID=$(id -g)

# Check if url and api key are provided
ENV_FILE_MANAGER=./manager/.env
ERROR=false
if ! grep -q "FILECOIN_NODE_URL" $ENV_FILE_MANAGER; then
    echo "Invalid .env file in $ENV_FILE_MANAGER, missing FILECOIN_NODE_URL"
    ERROR=true
elif grep -Fxq "FILECOIN_NODE_URL=changeme" $ENV_FILE_MANAGER; then
    echo "Error: please add a valid FILECOIN_NODE_URL in $ENV_FILE_MANAGER"
    ERROR=true
fi

if ! grep -q "RIPE_API_KEY" $ENV_FILE_MANAGER; then
    echo "Invalid .env file in $ENV_FILE_MANAGER, missing RIPE_API_KEY"
    ERROR=true
elif grep -Fxq "RIPE_API_KEY=changeme" $ENV_FILE_MANAGER; then
    echo "Error: please add a valid RIPE_API_KEY in $ENV_FILE_MANAGER"
    ERROR=true
fi

if [ "$ERROR" = true ]; then
    echo "Exit"
    exit 1
fi

# Create exports folder
mkdir -p exports

# Start services
docker-compose up