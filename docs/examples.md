## TODO

* Change to containers

```bash
go test ./... -cover -run UnitTest

go build -v -o docker-flow-cron

./docker-flow-cron

curl -XPUT \
    -d '{"Name": "my-job", "Image": "alpine", "Command": "hello World", "Schedule": "@every 10s"}' \
    "http://localhost:8080/v1/docker-flow-cron/job"

# Stop docker-flow-cron

./docker-flow-cron
```