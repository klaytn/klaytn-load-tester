FROM golang:1.18-buster as builder

RUN apt update && apt install -y make

ENV SRC_DIR /go/src/github.com/klaytn/klaytn-load-tester
ENV GOPATH /go

WORKDIR $SRC_DIR
ADD . .

RUN make && cp $SRC_DIR/build/bin/klayslave /bin/

FROM python:3.7-buster

RUN pip3 install locust==1.2.3

COPY --from=builder /bin/klayslave /bin/klayslave
