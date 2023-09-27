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
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o payment gitlab.com/mefit/mefit-server/cmd/payment/main.go
#RUN go build -v  -o cron1 youtab/cron/expWorker/main.go
#RUN go build -v  -o cron2 youtab/cron/mailWorker/main.go

FROM alpine
WORKDIR /
ENV GOPATH /go
RUN apk update && \
     apk add libc6-compat && \
     apk add ca-certificates
COPY --from=builder /go/src/payment .
COPY static static
COPY templates templates
EXPOSE 80
CMD ["./payment"]
