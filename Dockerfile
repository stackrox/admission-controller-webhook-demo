FROM alpine:latest

RUN apk add --no-cache ca-certificates
WORKDIR /code
USER 1001
COPY bin/linux/runasuser-admission-controller .
ENTRYPOINT ["/code/runasuser-admission-controller"]
EXPOSE 8443
