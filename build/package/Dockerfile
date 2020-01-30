FROM alpine:edge

RUN wget https://dl.google.com/go/go1.13.3.linux-amd64.tar.gz
RUN tar -xvf go1.13.3.linux-amd64.tar.gz
RUN mv go /usr/local

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

#RUN ls $GOROOT
#RUN ls $GOROOT/bin
#RUN echo Go Config:
#RUN go version
#RUN go env

RUN apk update && \
	apk add \
	gcc \
	mesa-dev \
	libx11-dev \
	libc-dev \
	libx11-dev \
	libxcursor-dev \
	libxi-dev \
	libxinerama-dev \
	libxrandr-dev \
	xorg-server \
	xvfb \
	coreutils \
	#mesa \
	#mesa-gl \
	#mesa-demos \
	#mesa-osmesa \
	xf86-video-dummy \
	xvfb-run

RUN apk add xf86-input-libinput
RUN apk add virtualgl --update-cache --repository http://dl-3.alpinelinux.org/alpine/edge/testing/ --allow-untrusted

COPY ./build/package/xorg.conf /root/xorg.conf
COPY ./build/package/entryfile.sh /root/entryfile.sh
ADD ./ /root/src

ENTRYPOINT ["sh", "/root/entryfile.sh"]