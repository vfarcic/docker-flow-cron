## TODO

* Remove

```bash
go test ./... -cover -run UnitTest

go build -v -o docker-flow-cron && \
    ./docker-flow-cron
```

## Creating Jobs

```bash
docker stack deploy -c stack.yml cron

docker stack ps cron

curl -XPUT \
    -d '{
    "Name": "my-job",
    "Image": "alpine",
    "Command": "echo \"hello World\"",
    "Schedule": "@every 30s"
}' "http://localhost:8080/v1/docker-flow-cron/job"

curl -XGET \
    "http://localhost:8080/v1/docker-flow-cron/job"

curl -XGET \
    "http://localhost:8080/v1/docker-flow-cron/job/my-job"

# retrieve job execution logs (`docker service logs`)
```