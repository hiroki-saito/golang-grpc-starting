FROM golang:1.18

RUN apt-get update \
    && apt-get install git autoconf automake libtool curl make g++ unzip -y \
    && cd /usr/local \
    && git clone https://github.com/google/protobuf.git \
    && cd protobuf \
    && git submodule update --init --recursive \
    && ./autogen.sh \
    && ./configure \
    && make \
    && make check \
    && make install \
    && ldconfig \
    && apt-get clean

ADD . /go/src/app
RUN cd /go/src/app \
    && go mod tidy
ENV PATH $PATH:$GOPATH/bin
