FROM alpine:latest

RUN wget https://dl.google.com/go/go1.13.3.linux-amd64.tar.gz
RUN tar -xvf go1.13.3.linux-amd64.tar.gz
RUN apk add sudo
RUN sudo mv go /usr/local

#RUN apk --allow-untrusted -X https://apkproxy.herokuapp.com/sgerrand/alpine-pkg-glibc add glibc glibc-bin

RUN apk --no-cache add ca-certificates wget
RUN wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub
RUN wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.28-r0/glibc-2.28-r0.apk
RUN apk add glibc-2.28-r0.apk

ENV GOROOT=/usr/local/go
ENV GOPATH=$HOME/go
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin
ENV GOGCCFLAGS=-std=c99
ENV GCCGO=gccgo
ENV CC=gcc
ENV CXX=g++
ENV CGO_ENABLED=1
ENV CGO_CPPFLAGS=
ENV CGO_CFLAGS="-g -O2 --std=gnu11"
ENV CGO_CXXFLAGS="-g -O2 --std=gnu11"
ENV CGO_FFLAGS="-g -O2"
ENV CGO_LDFLAGS="-g -O2"

RUN ls $GOROOT
RUN ls $GOROOT/bin
RUN echo Go Config:
RUN go version
RUN go env

RUN sudo apk update
RUN sudo apk add gcc
RUN sudo apk add mesa-dev
RUN sudo apk add libx11-dev
RUN sudo apk add libc-dev
RUN sudo apk add libx11-dev
RUN sudo apk add libxcursor-dev
RUN sudo apk add libxi-dev
RUN sudo apk add libxinerama-dev
RUN sudo apk add libxrandr-dev
RUN sudo apk add xorg-server
RUN sudo apk add xvfb
RUN sudo apk add coreutils
RUN sudo apk add xvfb-run --update-cache --repository http://dl-3.alpinelinux.org/alpine/edge/main/ --allow-untrusted

COPY ./build/package/entryfile.sh /entryfile.sh
ADD ./ /root/src

ENTRYPOINT ["/entryfile.sh"]