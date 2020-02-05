@ECHO OFF
echo %%~dp0 is "%~dp0"
pushd %~dp0
echo %%~dp0 is "%~dp0"

set dest=..\internal\profiler\

rd /s /q %dest%gl-proxy\
go-gen-proxy github.com/go-gl/gl/v3.2-core/gl %dest%gl-proxy/
for %%f in (%dest%gl-proxy\*) do (
	echo %%f
	(echo // +build debug) >%%f.tmp
	type %%f >>%%f.tmp
	move /y %%f.tmp %%f
)

rd /s /q %dest%gl-proxy-noop\
go-gen-proxy github.com/go-gl/gl/v3.2-core/gl %dest%gl-proxy-noop/ noop
for %%f in (%dest%gl-proxy-noop\*) do (
	(echo // +build !debug) >%%f.tmp
	type %%f >>%%f.tmp
	move /y %%f.tmp %%f
)