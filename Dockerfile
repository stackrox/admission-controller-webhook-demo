FROM alpine:latest

RUN apk update && \
    apk upgrade && \
    apk --no-cache add ca-certificates && \
    rm -rf /var/cache/apk/*
WORKDIR /code
USER 1001
COPY bin/linux/runasuser-admission-controller .
ENTRYPOINT ["/code/runasuser-admission-controller"]
EXPOSE 8443
