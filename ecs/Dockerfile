FROM jap:latest
MAINTAINER Sam Whited <swhited@atlassian.com>

RUN apt-get update
RUN apt-get -y --force-yes install awscli

WORKDIR /go/src/github.com/jitsi/jap/cmd/jap

COPY secrets-entrypoint secrets-entrypoint
RUN chmod +x secrets-entrypoint

CMD ["./secrets-entrypoint", "dumb-init", "jap", "-http=:8080"]
