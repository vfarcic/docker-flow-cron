## TODO

* Remove

```bash
go test ./... -cover -run UnitTest -p=1

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

curl -XPUT \
    -d '{
    "Image": "alpine",
    "Command": "echo \"hello World\"",
    "Schedule": "@every 15s"
}' "http://localhost:8080/v1/docker-flow-cron/job/my-other-job"

# Wait for 15 seconds

curl -XGET \
    "http://localhost:8080/v1/docker-flow-cron/job/my-job"

docker service logs blbfbgamad7yg3ll8egjhxcrw

# NOTE: Requires Docker 1.13+ with experimental features enabled

curl -XDELETE \
    "http://localhost:8080/v1/docker-flow-cron/job/my-other-job"

curl -XGET \
    "http://localhost:8080/v1/docker-flow-cron/job"







curl -XDELETE \
    "http://localhost:8080/v1/docker-flow-cron/job/my-job"

docker service ls
```