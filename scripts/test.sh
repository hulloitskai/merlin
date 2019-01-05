#!/usr/bin/env bash

## Bring image online.
export VERSION="$(make version)"
docker-compose up -d && sleep 15

## Set environment variables for healthcheck.
export HEALTH_ENDPOINT=http://localhost:3000

## Return success only if both healthcheck and image unload was successful.
if ! (./scripts/healthcheck.sh && docker-compose down); then exit -1; fi
