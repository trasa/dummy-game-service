# Dummy Game Service

Simple web service for ... well, nothing, actually.

This is put together to handle incoming webhook requests from other
services. It basically just receives the request and responds
200 OK. Or maybe an error, if that's what you configured it to do.

Other than logging everything to the console, that's about it.

## Building and Deploying

To create the docker container for fake service, see ```docker_build.sh``` which basically just runs

```
docker build -t registry.wargaming.net/freya/dummy-game-service .
```

To run the container, see ```docker_run.sh```
