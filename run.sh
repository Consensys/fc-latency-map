#!/bin/bash

export FC_UID=$(id -u)
export FC_GID=$(id -g)

ENV_FILE_MANAGER=./manager/.env
if test -f "$ENV_FILE_MANAGER"; then
    echo "$ENV_FILE_MANAGER exists."
else 
    echo "Error: $ENV_FILE_MANAGER does not exist."
    exit 1
fi

ENV_FILE_MAP=./map/.env
if test -f "$ENV_FILE_MAP"; then
    echo "$ENV_FILE_MAP exists."
else 
    echo "Error: $ENV_FILE_MAP does not exist."
    exit 1
fi

mkdir -p exports
docker-compose up