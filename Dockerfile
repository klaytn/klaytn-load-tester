ARG DOCKER_BASE_IMAGE=klaytn/build_base:1.2-go.1.18-solc0.8.13

FROM ${DOCKER_BASE_IMAGE}

ENV PKG_DIR /locust-docker-pkg
ENV SRC_DIR /go/src/github.com/klaytn/klaytn-load-tester
ENV GOPATH /go

RUN mkdir -p $PKG_DIR/bin

ADD . $SRC_DIR

RUN cd $SRC_DIR/klayslave && go build -ldflags "-linkmode external -extldflags -static"
RUN cp $SRC_DIR/klayslave/klayslave $PKG_DIR/bin
RUN pip3 install locust==1.2.3
