#!/bin/bash

dest=../internal/profiler

rm -rf $dest/gl-proxy/
go-gen-proxy github.com/go-gl/gl/v3.2-core/gl $dest/gl-proxy/
for file in $dest/gl-proxy/*.go; do
    sed -i '1s/^/\/\/ +build debug\n/' $file
done

rm -rf $dest/gl-proxy-noop/
go-gen-proxy github.com/go-gl/gl/v3.2-core/gl $dest/gl-proxy-noop/ noop
for file in $dest/gl-proxy-noop/*.go; do
    sed -i '1s/^/\/\/ +build !debug\n/' $file
done