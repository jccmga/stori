#!/bin/bash

CURR_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
docker compose -f "$CURR_DIR/build/docker-compose.yml" up db -d --wait --wait-timeout 30
FILEPATH=$1 EMAIL=$2 docker compose -f "$CURR_DIR/build/docker-compose.yml" up main --build
exit $(docker inspect main --format='{{.State.ExitCode}}')
