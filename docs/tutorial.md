# Examples of Running Docker Flow Cron In a Swarm Cluster

The examples that follow assume that you already have a Swarm cluster and that you are logged into one of the managers.

## Creating Jobs

We'll start by downloading a stack that fill deploy the `docker-flow-cron` service.

```bash
curl -o cron.yml \
    https://raw.githubusercontent.com/vfarcic/docker-flow-cron/master/stack.yml
```

The definition of the stack is as follows.

```
version: "3"

services:

  main:
    image: vfarcic/docker-flow-cron
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    ports:
      - ${PORT:-8080}:8080
    deploy:
      placement:
        constraints: [node.role == manager]
```

As you can see, it is a very simple stack. It contains a single service. It mounts `/var/run/docker.sock` as a volume. The `cron` will use it for communication with Docker Engine. The internal port `8080` will be exposed as `8080` on the host unless the environment variable `PORT` is specified. Finally, we're using a constraint that will limit the `cron` to one of the manager nodes.

Let us deploy the stack.

```bash
docker stack deploy -c cron.yml cron
```

A few moments later, the service will be up and running. We can confirm that with the `stack ps` command.

```bash
docker stack ps cron
```

The output is as follows.

```
ID            NAME         IMAGE                            NODE  DESIRED STATE  CURRENT STATE          ERROR  PORTS
auy9ajs8mgyn  cron_main.1  vfarcic/docker-flow-cron:latest  moby  Running        Running 4 seconds ago
```

Now that the service is running, we can schedule the first job.

TODO: Continue

TODO: Add at least two args

```bash
curl -XPUT \
    -d '{
    "Image": "alpine",
    "Command": "echo \"hello World\"",
    "Schedule": "@every 15s"
}' "http://localhost:8080/v1/docker-flow-cron/job/my-job"
```

```
{
  "Status": "OK",
  "Message": "Job my-job has been scheduled",
  "Job": {
    "Name": "my-job",
    "Image": "alpine",
    "Command": "echo \"hello World\"",
    "Schedule": "@every 15s",
    "Args": null
  },
  "Executions": null
}
```

```bash
# Wait for 15 seconds

docker service ls
```

```
ID            NAME            MODE        REPLICAS  IMAGE
7lfroifdmw00  thirsty_bartik  replicated  0/1       alpine:latest
vp1bcimbwsj5  cron_main       replicated  1/1       vfarcic/docker-flow-cron:latest
```

```bash
# Wait for 15 seconds

docker service ls
```

```
ID            NAME            MODE        REPLICAS  IMAGE
7lfroifdmw00  thirsty_bartik  replicated  0/1       alpine:latest
lz0ehsvy3fiu  pensive_liskov  replicated  0/1       alpine:latest
vp1bcimbwsj5  cron_main       replicated  1/1       vfarcic/docker-flow-cron:latest
```

```bash
curl -XGET \
    "http://localhost:8080/v1/docker-flow-cron/job"
```

```
{
  "Status": "OK",
  "Message": "",
  "Jobs": {
    "my-job": {
      "Name": "my-job",
      "Image": "alpine:latest@sha256:dfbd4a3a8ebca874ebd2474f044a0b33600d4523d03b0df76e5c5986cb02d7e8",
      "Command": "docker service create --restart-condition none alpine echo \"hello World\"",
      "Schedule": "@every 15s",
      "Args": null
    }
  }
}
```

```bash
curl -XGET \
    "http://localhost:8080/v1/docker-flow-cron/job/my-job"
```

```
{
  "Status": "OK",
  "Message": "",
  "Job": {
    "Name": "my-job",
    "Image": "alpine:latest@sha256:dfbd4a3a8ebca874ebd2474f044a0b33600d4523d03b0df76e5c5986cb02d7e8",
    "Command": "docker service create --restart-condition none alpine echo \"hello World\"",
    "Schedule": "@every 15s",
    "Args": null
  },
  "Executions": [
    {
      "ServiceId": "85iyzvg00utszweglyle59i1y",
      "CreatedAt": "2017-02-28T22:22:24.786950411Z",
      "Status": {
        "Timestamp": "2017-02-28T22:22:26.556923596Z",
        "State": "complete",
        "Message": "finished",
        "ContainerStatus": {
          "ContainerID": "b415243a6e98204fb55fd6ce8ab378174cd58a550fa212514020ccb34bc61677"
        },
        "PortStatus": {}
      }
    },
    {
      "ServiceId": "j7s37k0jy12u6203pne4v50xz",
      "CreatedAt": "2017-02-28T22:23:25.33680838Z",
      "Status": {
        "Timestamp": "2017-02-28T22:23:27.118679889Z",
        "State": "complete",
        "Message": "finished",
        "ContainerStatus": {
          "ContainerID": "d608c4cb5f9594020f5067406c1cb50ccfe691d9b9ec220d4b12ffa05e59a777"
        },
        "PortStatus": {}
      }
    },
    ...
  ]
}
```

```bash
curl -XPUT \
    -d '{
    "Image": "alpine",
    "Command": "echo \"hello World\"",
    "Schedule": "@every 15s"
}' "http://localhost:8080/v1/docker-flow-cron/job/my-other-job"
```

```
{
  "Status": "OK",
  "Message": "Job my-other-job has been scheduled",
  "Job": {
    "Name": "my-other-job",
    "Image": "alpine",
    "Command": "echo \"hello World\"",
    "Schedule": "@every 15s",
    "Args": null
  },
  "Executions": null
}
```

```bash
# Wait for 15 seconds

curl -XGET \
    "http://localhost:8080/v1/docker-flow-cron/job/my-other-job"
```

```
{
  "Status": "OK",
  "Message": "",
  "Job": {
    "Name": "my-other-job",
    "Image": "alpine:latest@sha256:dfbd4a3a8ebca874ebd2474f044a0b33600d4523d03b0df76e5c5986cb02d7e8",
    "Command": "docker service create --restart-condition none alpine echo \"hello World\"",
    "Schedule": "@every 15s",
    "Args": null
  },
  "Executions": [
    {
      "ServiceId": "k92tbts7itviky14ucaxnozds",
      "CreatedAt": "2017-02-28T22:25:26.518009461Z",
      "Status": {
        "Timestamp": "2017-02-28T22:25:38.064998905Z",
        "State": "complete",
        "Message": "finished",
        "ContainerStatus": {
          "ContainerID": "72517cd1d328ee24c630fd0ec230473f4d2be49baf76aa4d1b4f36b7d9186d1b"
        },
        "PortStatus": {}
      }
    },
    ...
  ]
}
```

```bash
# NOTE: Requires Docker 1.13+ with experimental features enabled
docker service logs k92tbts7itviky14ucaxnozds
```

```
tender_bhabha.1.6yozcuzf9abf@moby    | hello World
```

```bash
curl -XGET \
    "http://localhost:8080/v1/docker-flow-cron/job"
```

```
{
  "Status": "OK",
  "Message": "",
  "Jobs": {
    "my-job": {
      "Name": "my-job",
      "Image": "alpine:latest@sha256:dfbd4a3a8ebca874ebd2474f044a0b33600d4523d03b0df76e5c5986cb02d7e8",
      "Command": "docker service create --restart-condition none alpine echo \"hello World\"",
      "Schedule": "@every 15s",
      "Args": null
    },
    "my-other-job": {
      "Name": "my-other-job",
      "Image": "alpine:latest@sha256:dfbd4a3a8ebca874ebd2474f044a0b33600d4523d03b0df76e5c5986cb02d7e8",
      "Command": "docker service create --restart-condition none alpine echo \"hello World\"",
      "Schedule": "@every 15s",
      "Args": null
    }
  }
}
```

```bash
curl -XDELETE \
    "http://localhost:8080/v1/docker-flow-cron/job/my-other-job"
```

```
{
  "Status": "OK",
  "Message": "my-other-job was deleted",
  "Job": {
    "Name": "",
    "Image": "",
    "Command": "",
    "Schedule": "",
    "Args": null
  },
  "Executions": null
}
```

```bash
curl -XGET \
    "http://localhost:8080/v1/docker-flow-cron/job"
```

```
{
  "Status": "OK",
  "Message": "",
  "Jobs": {
    "my-job": {
      "Name": "my-job",
      "Image": "alpine:latest@sha256:dfbd4a3a8ebca874ebd2474f044a0b33600d4523d03b0df76e5c5986cb02d7e8",
      "Command": "docker service create --restart-condition none alpine echo \"hello World\"",
      "Schedule": "@every 15s",
      "Args": null
    }
  }
}
```

```bash
docker stack rm cron

docker stack deploy -c cron.yml cron

curl -XGET \
    "http://localhost:8080/v1/docker-flow-cron/job"
```

```
{
  "Status": "OK",
  "Message": "",
  "Jobs": {
    "my-job": {
      "Name": "my-job",
      "Image": "alpine:latest@sha256:dfbd4a3a8ebca874ebd2474f044a0b33600d4523d03b0df76e5c5986cb02d7e8",
      "Command": "docker service create --restart-condition none alpine echo \"hello World\"",
      "Schedule": "@every 15s",
      "Args": null
    }
  }
}
```

```bash
curl -XDELETE \
    "http://localhost:8080/v1/docker-flow-cron/job/my-job"

curl -XGET \
    "http://localhost:8080/v1/docker-flow-cron/job"
```

```
{
  "Status": "OK",
  "Message": "",
  "Jobs": {}
}
```

```bash
docker service ls
```

```
ID            NAME       MODE        REPLICAS  IMAGE
nvmq69qthhqz  cron_main  replicated  1/1       vfarcic/docker-flow-cron:latest
```
