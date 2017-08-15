# Dummy Game Service

Simple web service for ... well, nothing, actually.

This is put together to handle incoming webhook requests from other
services. It basically just receives the request and responds
200 OK. Or maybe an error, if that's what you configured it to do.

Listeners can register via <http://dummy-game.wgseattle.com/static/#/> to 
see the webhook traffic as they arrive. There's also a way to verify
the hash given a secret key and a (complete) webhook body.

Other than logging everything to the console, that's about it.

## Building and Deploying to registry.wargaming.net

To create the docker container for fake service, see ```docker_build.sh``` which basically just runs

```
docker build -t registry.wargaming.net/freya/dummy-game-service .
```

To run the container locally, see ```docker_run.sh```

## Deploying to AWS

First, you'll need to build & deploy the docker image to create 
"registry.wargaming.net/freya/dummy-game-service:latest"

```
./docker_build.sh
```

Next, we need to tag & push this docker image up to dockerhub - this
creates [docker.io/wgplatform/dummy-game-service](https://hub.docker.com/r/wgplatform/dummy-game-service/)

```
./docker_push.sh
```

Now that our image is up on dockerhub, we need to tell AWS
ECS to use the new image. 

1. Go to [ECS Task Definitions page](https://us-west-2.console.aws.amazon.com/ecs/home?region=us-west-2#/taskDefinitions)
1. Select check-box for task definition "dummy-game-server" and click "Create new revision"
1. Leave everything as-is and click "create"
1. On the new revision, click the "Actions v" button / dropdown.
1. Select "Update Service"
1. On the "Update Service" page click the button labeled 
   "Update Service" to update the service.
1. Wait, patiently, for AWS to "do its thing".   
   
**Important Note** 
If your service is configured to use only 1 EC2 instance, you _must_
set the service accordingly:

* Number of tasks: 1
* Minimum Healthy Percent: 0
* Maximum Percent: 200

Otherwise, AWS will never kill the old task (making healthy percent = 0%)
and will therefore never start the new task. Note that 
[dummy-game-service](https://us-west-2.console.aws.amazon.com/ecs/home?region=us-west-2#/clusters/platform/services/dummy-game-server-service/update)
is already configured this way, but if you create a new ECS service
this is one of those "gotchas" that you'll be swearing at for a while.