FROM golang:1.7.5

MAINTAINER 	Viktor Farcic <viktor@farcic.com>

RUN apt-get update && \
    apt-get install -y apt-transport-https ca-certificates curl software-properties-common expect && \
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - && \
    add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/debian $(lsb_release -cs) stable" && \
    apt-get update && \
    apt-get -y install docker-ce

RUN go get github.com/docker/docker/api/types && \
    go get github.com/docker/docker/api/types/filters && \
    go get github.com/docker/docker/api/types/swarm && \
    go get github.com/docker/docker/client && \
    go get gopkg.in/robfig/cron.v2 && \
    go get golang.org/x/net/context && \
    go get github.com/gorilla/mux && \
    go get github.com/stretchr/testify/suite

COPY . /src
WORKDIR /src