```bash
# TODO: Setup a cluster

docker container run -it --rm \
    -v /var/run/docker.sock:/var/run/docker.sock \
    docker docker image prune -f

docker service create --name cron-prune-images \
    --mount "type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock" \
    docker docker image prune -f

docker service ps cron-prune-images

docker service logs -t cron-prune-images

docker service rm cron-prune-images

docker service create --name cron-prune-images \
    --mount "type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock" \
    --restart-delay 10s \
    docker docker image prune -f

docker service ps cron-prune-images

docker service logs -t cron-prune-images

docker service rm cron-prune-images

docker service create --name cron-prune-images \
    --mount "type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock" \
    --restart-delay 10s \
    --mode global \
    docker docker image prune -f

docker service ps cron-prune-images

docker service logs -t cron-prune-images

docker service rm cron-prune-images

# NOTE: No flexibility expressing the desired time.

# NOTE: No feature to restart imediatelly on failure.

# NOTE: What is missing?
```