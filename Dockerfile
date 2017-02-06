FROM alpine:3.5
MAINTAINER 	Viktor Farcic <viktor@farcic.com>

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

EXPOSE 8080

CMD ["docker-flow-cron"]

COPY docker-flow-cron /usr/local/bin/docker-flow-cron
RUN chmod +x /usr/local/bin/docker-flow-cron
