FROM alpine:edge

RUN apk update
RUN apk add git musl-dev go=1.13.10-r0

ENV GOROOT=/usr/lib/go
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

RUN echo Go Config: &&\
	go version &&\
	go env

# Dependencies for GLFW
RUN	apk add \
	build-base \
	libx11-dev \ 
	libxcursor-dev \
	libxrandr-dev \
	libxinerama-dev \
	libxi-dev \
	mesa-dev

# Required to run xvfb-run
RUN apk add mesa-dri-gallium xvfb-run

# virtualgl includes glxinfo
RUN apk add virtualgl --update-cache --repository http://dl-3.alpinelinux.org/alpine/edge/testing/ --allow-untrusted

COPY ./build/package/entryfile.sh /root/entryfile.sh
ADD ./ /root/src/

ENTRYPOINT ["sh", "/root/entryfile.sh"]