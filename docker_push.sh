#!/bin/bash

docker tag  registry.wargaming.net/freya/dummy-game-service wgplatform/dummy-game-service
docker push wgplatform/dummy-game-service 

