#!/bin/bash

docker run --rm \
  --user $(id -u):$(id -g) \
  -v /$(pwd)/manager/data:/latency-db keinos/sqlite3 \
  sh -c "sqlite3 /latency-db/database.db < /latency-db/dump.sql"