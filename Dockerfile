FROM golang:1.3.3-onbuild
MAINTAINER 3onyc <3onyc@x3tech.com>

ENV AMQP_EXCHANGE docker.events
ENV AMQP_USER guest
ENV AMQP_PASSWORD guest

ADD contrib/docker-run.sh /run.sh
CMD /run.sh
