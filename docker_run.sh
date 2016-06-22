#!/bin/bash

docker run -d \
    -p 8080:8080 \
    --name dev-dummy-game-service \
    -e "SERVICE_TAGS=dev" \
    registry.wargaming.net/freya/dummy-game-service
