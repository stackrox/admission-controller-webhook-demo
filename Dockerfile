FROM golang:1.13 as builder
WORKDIR /go/src/app
ADD . .
RUN make build-linux-amd64

FROM scratch
COPY --from=builder /go/src/app/bin/linux/runasuser-admission-controller ./
ENTRYPOINT ["./runasuser-admission-controller"]
EXPOSE 8443
