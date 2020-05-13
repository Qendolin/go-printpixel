#!/bin/sh
cd /root/src/
echo Starting Xvfb
export DISPLAY=:10
Xvfb :10 -screen 0 1024x768x24 +extension GLX +render -noreset -ac &
echo Installing go packages
go get -v -t -d ./...
echo Starting Test
go test -timeout 120s ./... -headless