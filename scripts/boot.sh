#!/usr/bin/env bash

set -o errexit

composeFiles="-f ../docker-compose.yml"

echo '==> refresh images'
docker-compose ${composeFiles} --profile "delayed-start" pull

echo '==> building environment'
docker-compose ${composeFiles} --profile "delayed-start" build

echo '==> starting environment'
docker-compose ${composeFiles} up -d

echo '==> done'