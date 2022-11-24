ARG SRC_DIR=/go/src/github.com/klaytn/klaytn-load-tester

FROM golang:1.18-buster as builder
ARG SRC_DIR

RUN apt update && apt install -y make

ENV GOPATH /go

WORKDIR $SRC_DIR
ADD . .

RUN make && cp $SRC_DIR/build/bin/klayslave /bin/

FROM python:3.7-buster
ARG SRC_DIR

RUN pip3 install locust==1.2.3
RUN mkdir -p /locust-docker-pkg/bin

COPY --from=builder /bin/klayslave /locust-docker-pkg/bin/klayslave
COPY --from=builder $SRC_DIR/dist/locustfile.py /locust-docker-pkg/locustfile.py 
RUN ln -s /locust-docker-pkg/bin/klayslave /bin/klayslave
