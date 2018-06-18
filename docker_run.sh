#!/bin/bash

docker run -d \
    -p 8080:8080 \
    --name dev-dummy-game-service \
    -e "SERVICE_TAGS=dev" \
    yourregistryhere.com/yourprojecthere/dummy-game-service
