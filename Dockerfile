# This file is a template, and might need editing before it works on your project.
FROM golang:1.10 AS builder
WORKDIR $GOPATH/src/gitlab.com/mefit/mefit-server
# This will download deps in docker file ignored for faster build(copy vendor folder)
# ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
# RUN chmod +x /usr/bin/dep
# COPY Gopkg.toml Gopkg.lock ./
# RUN https_proxy=https://proxy.sapph.ir:8000 dep ensure --vendor-only
COPY . ./
WORKDIR $GOPATH/src
RUN go get -u github.com/sinabakh/go-zarinpal-checkout
RUN go get -u github.com/davecgh/go-spew/spew
# RUN go get -u golang.org/x/net/trace
RUN CGO_ENABLED=0 GOOS=linux go build -o server gitlab.com/mefit/mefit-server/cmd/grpc/main.go
#RUN go build -v  -o cron1 youtab/cron/expWorker/main.go
#RUN go build -v  -o cron2 youtab/cron/mailWorker/main.go

FROM alpine
WORKDIR /
RUN apk update && \
     apk add libc6-compat && \
     apk add ca-certificates
COPY --from=builder /go/src/server .
COPY static static
COPY certs   certs
# COPY --from=builder /go/src/cron1 .
# COPY --from=builder /go/src/cron2 .
EXPOSE 8443
CMD ["./server"]
