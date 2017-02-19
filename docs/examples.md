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

curl -XGET \
    "http://localhost:8080/v1/docker-flow-cron/job"

curl -XGET \
    "http://localhost:8080/v1/docker-flow-cron/job/my-job"

curl -XPUT \
    -d '{
    "Image": "alpine",
    "Command": "echo \"hello World\"",
    "Schedule": "@every 15s"
}' "http://localhost:8080/v1/docker-flow-cron/job/my-other-job"

# Wait for 15 seconds

curl -XGET \
    "http://localhost:8080/v1/docker-flow-cron/job/my-other-job"

# NOTE: Requires Docker 1.13+ with experimental features enabled
docker service logs qyoqwu1r6bc74l6w8u1lt3zi9

curl -XDELETE \
    "http://localhost:8080/v1/docker-flow-cron/job/my-other-job"

curl -XGET \
    "http://localhost:8080/v1/docker-flow-cron/job"

curl -XGET \
    "http://localhost:8080/v1/docker-flow-cron/job/my-job"

# TODO: Add the command to stop the cron service

# TODO: Add the command to start the cron service

curl -XGET \
    "http://localhost:8080/v1/docker-flow-cron/job/my-job"








curl -XDELETE \
    "http://localhost:8080/v1/docker-flow-cron/job/my-job"

curl -XGET \
    "http://localhost:8080/v1/docker-flow-cron/job/my-job"

docker service ls
```