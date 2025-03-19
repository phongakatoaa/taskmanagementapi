#!/usr/bin/env bash

set -o nounset
set -o errexit
set +e

composeFiles="-f ../docker-compose.yml"

echo "==> stopping and removing environment"

docker-compose ${composeFiles} --profile "delayed-start" stop
docker-compose ${composeFiles} --profile "delayed-start" rm -v -f

echo "==> removing all unattached volumes"

docker volume rm "$(docker volume ls -f dangling=true -q)"

echo "==> done"
