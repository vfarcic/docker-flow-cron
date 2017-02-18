## TODO

* Remove

```bash
go build -v -o docker-flow-cron && \
    ./docker-flow-cron
```

## Creating Jobs

```bash
docker stack deploy -c stack.yml cron

docker stack ps cron

curl -XPUT \
    -d '{
    "Image": "alpine",
    "Command": "echo \"hello World\"",
    "Schedule": "@every 15s"
}' "http://localhost:8080/v1/docker-flow-cron/job/my-job"

# Wait for 15 seconds

docker service ls

# Wait for 15 seconds

docker service ls
```

```
ID            NAME           MODE        REPLICAS  IMAGE
p9vctaewty4o  frosty_wilson  replicated  1/1       alpine:latest
```

```bash
curl -XGET \
    "http://localhost:8080/v1/docker-flow-cron/job"
```

```bash
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
      "ServiceId": "z5fcqdh8k2cagntozf8fmf0q3",
      "CreatedAt": "2017-02-18T00:52:05.677039953Z",
      "Status": {
        "Timestamp": "2017-02-18T00:52:07.584812118Z",
        "State": "complete",
        "Message": "finished",
        "ContainerStatus": {
          "ContainerID": "4b6650cf28eec355a6078700b0c893067d5599946bb7cb56f69804eda625fb02"
        },
        "PortStatus": {}
      }
    },
    {
      "ServiceId": "7idpd6461941hc4q4z4xpbatm",
      "CreatedAt": "2017-02-18T00:51:50.728345649Z",
      "Status": {
        "Timestamp": "2017-02-18T00:51:52.641787817Z",
        "State": "complete",
        "Message": "finished",
        "ContainerStatus": {
          "ContainerID": "e022110949ab571c595d5523211a27167690368ec78d66798818382644f23033"
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
      "ServiceId": "qyoqwu1r6bc74l6w8u1lt3zi9",
      "CreatedAt": "2017-02-18T00:53:25.640689104Z",
      "Status": {
        "Timestamp": "2017-02-18T00:53:27.495741267Z",
        "State": "complete",
        "Message": "finished",
        "ContainerStatus": {
          "ContainerID": "ab31969ae825c91d07dff5714cf2ed67c32531473fb9a91863184b2152af125a"
        },
        "PortStatus": {}
      }
    },
    {
      "ServiceId": "j03weoytp0s3j1e4mpqvc02wg",
      "CreatedAt": "2017-02-18T00:53:55.663722097Z",
      "Status": {
        "Timestamp": "2017-02-18T00:53:57.62880576Z",
        "State": "complete",
        "Message": "finished",
        "ContainerStatus": {
          "ContainerID": "83c821e629d057721a160b9435dd2209e781e5843217372d90fbe75703af1a2e"
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
docker service logs qyoqwu1r6bc74l6w8u1lt3zi9
```

```
quirky_keller.1.3uld8sytjzq5@moby    | hello World
```

```bash

curl -XDELETE \
    "http://localhost:8080/v1/docker-flow-cron/job/my-other-job"

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

# TODO: Add the command to stop the cron service

# TODO: Add the command to start the cron service

curl -XGET \
    "http://localhost:8080/v1/docker-flow-cron/job/my-job"








curl -XDELETE \
    "http://localhost:8080/v1/docker-flow-cron/job/my-job"

docker service ls
```