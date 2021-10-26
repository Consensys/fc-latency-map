#!/bin/bash

docker build -f manager/Dockerfile \
  --build-arg FC_UID="$(id -u)" \
  --build-arg FC_GID="$(id -g)" \
  -t fc-latency-manager .
