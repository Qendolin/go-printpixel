#!/bin/sh
cd /root/src/
go get -v -t -d ./...
echo Starting Test
#xvfb-run -e /dev/stderr glxinfo
xvfb-run -e /dev/stderr go test -timeout 120s ./... -headless
