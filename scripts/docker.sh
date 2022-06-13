#!/bin/bash

#################################################
source $(dirname "$0")/env.sh
cd $ROOT
dependencies
#################################################

docker_up() {
  docker-compose -f deployments/docker-compose.yaml -p dev_solocms up
}

docker_down() {
  docker-compose -f deployments/docker-compose.yaml -p dev_solocms down
}

case $1 in
docker_up)
  docker_down
  docker_up
  ;;
docker_down)
  docker_down
  ;;
*)
  echo "docker_up or docker_down"
  ;;
esac