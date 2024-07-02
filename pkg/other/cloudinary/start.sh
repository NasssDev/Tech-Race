#!/bin/bash
if [ "$1" == "start-build" ]; then
  # If the argument is "start", bring the containers up (use it for first time)
  docker-compose up -d --build
fi

if [ "$1" == "start-run" ]; then
  # If the argument is "start", bring the containers up (use it for first time)
  docker-compose up -d --build
  docker exec -it cloudinarace sh -c "go run main/main.go"
fi

if [ "$1" == "cloudinarace" ]; then
  docker exec -it cloudinarace sh -c "go run main/main.go"
fi