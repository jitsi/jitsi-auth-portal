FROM jap:latest
MAINTAINER Sam Whited <swhited@atlassian.com>

RUN apt-get update
RUN apt-get -y --force-yes install awscli

COPY secrets-entrypoint /secrets-entrypoint
RUN chmod +x /secrets-entrypoint

ENTRYPOINT ["/secrets-entrypoint"]
CMD ["dumb-init", "jap", "-http=:8080"]
