#!/bin/bash

docker run --user $(id -u):$(id -g) -v /$(pwd)/data:/latency-db keinos/sqlite3 sh -c "sqlite3 /latency-db/database.db < /latency-db/dump.sql"