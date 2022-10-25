FROM golang:1.18-buster as builder

ENV SRC_DIR /go/src/github.com/klaytn/klaytn-load-tester
ENV GOPATH /go

WORKDIR $SRC_DIR
ADD . .

RUN (cd klayslave && \
        go build -ldflags "-linkmode external -extldflags -static" && \
        cp klayslave /bin/)

FROM python:3.7-buster

RUN pip3 install locust==1.2.3

COPY --from=builder /bin/klayslave /bin/klayslave
