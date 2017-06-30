FROM golang:1.6 AS build
ADD . /src
WORKDIR /src
RUN go get -d -v -t
RUN go test --cover ./... --run UnitTest
RUN go build -v -o docker-flow-cron

FROM alpine:3.5
MAINTAINER 	Viktor Farcic <viktor@farcic.com>

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

EXPOSE 8080

CMD ["docker-flow-cron"]

ENV DOCKER_VERSION 1.13.1
RUN set -x \
    && apk add --no-cache curl \
	&& curl -fSL "https://get.docker.com/builds/Linux/x86_64/docker-${DOCKER_VERSION}.tgz" -o docker.tgz \
	&& tar -xzvf docker.tgz \
	&& mv docker/* /usr/local/bin/ \
	&& rmdir docker \
	&& rm docker.tgz \
	&& apk del curl

COPY --from=build /src/docker-flow-cron /usr/local/bin/docker-flow-cron
RUN chmod +x /usr/local/bin/docker-flow-cron
