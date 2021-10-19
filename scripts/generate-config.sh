#!/bin/bash

# Create manager env
ENV_FILE_MANAGER=./manager/.env
ENV_FILE_MANAGER_EXAMPLE=./manager/.env.example
if test -f "$ENV_FILE_MANAGER"; then
    echo "$ENV_FILE_MANAGER exists."
else 
    echo "$ENV_FILE_MANAGER does not exist, creating it."
    cp $ENV_FILE_MANAGER_EXAMPLE $ENV_FILE_MANAGER
fi

# Create map env
ENV_FILE_MAP=./map/.env
ENV_FILE_MAP_EXAMPLE=./map/.env.example
if test -f "$ENV_FILE_MAP"; then
    echo "$ENV_FILE_MAP exists."
else 
    echo "$ENV_FILE_MAP does not exist, creating it."
    cp $ENV_FILE_MAP_EXAMPLE $ENV_FILE_MAP
fi