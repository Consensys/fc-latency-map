#!/bin/bash

DATABASE_FILE=./manager/data/database.db
if test -f "$DATABASE_FILE"; then
  echo "Remove $DATABASE_FILE."
  rm $DATABASE_FILE
fi

docker run --rm \
  --user $(id -u):$(id -g) \
  -v /$(pwd)/manager/data:/latency-db keinos/sqlite3 \
  sh -c "sqlite3 /latency-db/database.db < /latency-db/dump.sql"