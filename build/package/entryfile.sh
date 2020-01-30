#!/bin/sh
cd /root/src/
##xvfb-run -e /dev/stderr --server-args=':99 -screen 0 640x480x8 +extension GLX +render -noreset -ac' glxinfo
#Xorg -noreset +extension GLX +extension RANDR +extension RENDER -logfile /dev/stdout -config /root/xorg.conf :10 &
#sleep 2
#echo ================
#export DISPLAY=:10
#vglrun -d :10 glxinfo
#echo ================
#vglrun go test --tags=headless ./...
go test --tags=headless ./...