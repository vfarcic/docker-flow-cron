# Docker Flow Cron

**THIS PROJECT IS IN DESIGN PHASE, DO NOT USE IT IN PRODUCTION YET!**

Docker Swarm services are designed to the long lasting processes that, potentially, live forever. Docker does not have a mechanism to schedule jobs based on time interval or, to put it in other words, it does not have the ability to use cron-like syntax for time-based scheduling.

*Docker Flow Cron* is designed to overcome some of the limitations behind Docker Swarm services and provide cron-like time-based scheduling while maintaining fault tolerance features available in Docker Swarm.

Please visit the **[project documentation](http://cron.dockerflow.com)** for more info or join the #df-cron Slack channel in [DevOps20](http://slack.devops20toolkit.com/) if you have any questions, suggestions, or problems.

# Usage

There is currently two different ways of using Docker Flow Cron

- Using the [Docker Flow Cron API](#docker-flow-cron-api) directly to manage (add, list, delete) scheduled jobs

- Using the [*Docker Flow Swarm Listener*](#docker-flow-swarm-listener-support) support to manage jobs by creating/deleting regular Docker Services.


## Docker Flow Cron API
#### Put Job

> Adds a job to docker-flow-cron

The following body parameters can be used to send a *create job* `PUT` request to *Docker Flow Cron*. They should be added to the base address **[CRON_IP]:[CRON_PORT]/v1/docker-flow-cron/job/[jobName]**.

|param           |Description                                                        |Mandatory|Example  |
|----------------|-------------------------------------------------------------------|---------|---------|
|image           |Docker image.                                                      |yes      |alpine   |
|serviceName     |Docker service name                                                |no       |my-cronjob  |
|command         |The command that will be executed when a job is created.           |no       |echo "hello World"|
|schedule        |The schedule that defines the frequency of the job execution.Check the [scheduling section](#scheduling) for more info. |yes|@every 15s|
|args            |The list of arguments that can be used with the `docker service create` command.<br><br>`--restart-condition` cannot be set to `any`. If not specified, it will be set to `none`.<br>`--name` argument is not allowed. Use serviceName param instead<br><br>Any other argument supported by `docker service create` is allowed.|no|TODO|

TODO: Example

#### Get All Jobs

> Gets all scheduled jobs

The following `GET` request **[CRON_IP]:[CRON_PORT]/v1/docker-flow-cron/job**. can be used to get all scheduled jobs from Docker Flow Cron.


#### Get Job

> Gets a job from docker-flow-cron

The following `GET` request **[CRON_IP]:[CRON_PORT]/v1/docker-flow-cron/[jobName]**. can be used to get a job from Docker Flow Cron.

#### Delete Job

> Deletes a job from docker-flow-cron

The following `DELETE` request **[CRON_IP]:[CRON_PORT]/v1/docker-flow-cron/[jobName]**. can be used to delete a job from Docker Flow Cron.


## *Docker Flow Swarm Listener* support

Using the *Docker Flow Swarm Listener* support, Docker Services can schedule jobs.
Docker Flow Swarm Listener listens to Docker Swarm events and sends requests to Docker Flow Cron when changes occurs, 
every time a service is created or deleted Docker Flow Cron gets notified and manages job scheduling.


A Docker Service is created with the following syntax:

```docker service create [OPTIONS] IMAGE [COMMAND] [ARG...]```

> Check the [documentation](https://docs.docker.com/engine/reference/commandline/service_create/) for more information.

The same syntax is used to schedule a job in Docker Flow Cron:

- ```IMAGE``` specifies the Docker Image
- ```[COMMAND]``` specifies the command to schedule.
- ```[OPTIONS]``` specifies some necessary options making scheduling Docker Services possible.


The following Docker Service args ```[OPTIONS]``` should be used for scheduled Docker Services:

|arg                 |Description                                                        |Mandatory|Example   |
|--------------------|-------------------------------------------------------------------|---------|----------|
|--replicas          | Set 0 to prevent the service from running immedietely. Set to 1 to run command on service creation.  |yes      |0 or 1    |
|--restart-condition | Set to ```none``` to prevent the service from using Docker Swarm's ability to autorestart exited services.                        |yes      |none      |



The following Docker Service labels ```[OPTIONS]``` needs to be used for scheduled Docker Services:

|label           |Description                                                        |Prefix|Mandatory|Example   |
|----------------|-------------------------------------------------------------------|------|---------|----------|
|cron            |Enable scheduling                                                  |com.df|yes      |true      |
|image           |Docker image.                                                      |com.df.cron|yes      |alpine    |
|name            |Cronjob name.                                                      |com.df.cron|yes      |my-cronjob|
|schedule        |The schedule that defines the frequency of the job execution. Check the [scheduling section](#scheduling) for more info.|com.df.cron|yes|@every 15s|
|command         |The command that is scheduled, only used for Docker Flow Cron registration. Use the same command you set for your docker service to run.|com.df.cron|No   |echo Hello World|

**All labels needs to be prefixed**

> Examples:
- ```--labels "com.df.cron=true"```
- ```--labels "com.df.cron.name=my-job"```


## Tutorial

#### Docker Flow Cron using Docker Stacks

> Create the overlay network

```docker network create -d overlay cron```

> Create the cron stack 

docker-compose.yml
```
version: "3"

services:
  cron:
    image: vfarcic/docker-flow-cron
    networks:
      - cron
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    ports:
      - 8080:8080
    deploy:
      placement:
        constraints: [node.role == manager]
  
  swarm-listener:
    image: vfarcic/docker-flow-swarm-listener
    networks:
      - cron
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      - DF_NOTIFY_CREATE_SERVICE_URL=http://cron:8080/v1/docker-flow-cron/job/create
      - DF_NOTIFY_REMOVE_SERVICE_URL=http://cron:8080/v1/docker-flow-cron/job/remove
    deploy:
      placement:
        constraints: [node.role == manager]
networks:
  cron:
    external: true
```

Deploy the stack above by executing ```docker stack deploy -c docker-compose.yml cron```



#### Docker Flow Cron using Docker Services

> Create Docker Flow Swarm Listener

```
docker service create --name swarm-listener \
    --network cron \
    --mount "type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock" \
    -e DF_NOTIFY_CREATE_SERVICE_URL=http://cron:8080/v1/docker-flow-cron/create \
    -e DF_NOTIFY_REMOVE_SERVICE_URL=http://cron:8080/v1/docker-flow-cron/remove \
    --constraint 'node.role==manager' \
    vfarcic/docker-flow-swarm-listener
```

> Create Docker Flow Cron

```
docker service create --name cron \
    --network cron \
    --mount "type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock" \
    -p 8080:8080 \
    --constraint 'node.role==manager' \
    vfarcic/docker-flow-cron
```

#### Example: Scheduling a Docker Service

Now that the cron stack is up and running we can add our first scheduled Docker Service.
The example service below will run ```echo "Hello World"``` every 10 seconds using an Alpine image.

> Example: Scheduled job using Docker Services

```
docker service create --name cronjob \
    --network cron --replicas 0 \
    --restart-condition=none \
    -l "com.df.notify=true" \
    -l "com.df.cron=true" \
    -l "com.df.cron.name=cronjob" \
    -l "com.df.cron.image=alpine" \
    -l "com.df.cron.command=echo Hello World" \
    -l "com.df.cron.schedule=@every 10s" \
    alpine \
    echo Hello world
```

The example below is a more realistic use case for *Docker Flow Cron*.
It will remove unusued docker data by running ```docker system prune``` every day.

> Example: Scheduling a docker command to run on the host
```
docker service create --name docker_cleanup \
    --network cron --replicas 0 \
    --restart-condition=none \
    --mount "type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock" \
    -l "com.df.notify=true" \
    -l "com.df.cron=true" \
    -l "com.df.cron.name=docker_cleanup" \
    -l "com.df.cron.image=docker" \
    -l "com.df.cron.command=docker system prune -f" \
    -l "com.df.cron.schedule=@daily" \
    docker \
    docker system prune -f
```

## Scheduling
Docker Flow Cron uses the library [robfig/cron](https://godoc.org/github.com/robfig/cron) to provide a simple cron syntax for scheduling.

> Examples
```
0 30 * * * *        Every hour on the half hour
@hourly             Every hour
@every 1h30m        Every hour thirty
```

#### Predefined schedules
You may use one of several pre-defined schedules in place of a cron expression.
```
Entry                  | Description                                | Equivalent To
-----                  | -----------                                | -------------
@yearly (or @annually) | Run once a year, midnight, Jan. 1st        | 0 0 0 1 1 *
@monthly               | Run once a month, midnight, first of month | 0 0 0 1 * *
@weekly                | Run once a week, midnight on Sunday        | 0 0 0 * * 0
@daily (or @midnight)  | Run once a day, midnight                   | 0 0 0 * * *
@hourly                | Run once an hour, beginning of hour        | 0 0 * * * *
```

#### Intervals 
You may also schedule a job to execute at fixed intervals

```
@every <duration>
@every 2h30m15s
```

Check the library [documentation](https://godoc.org/github.com/robfig/cron) for more information.
