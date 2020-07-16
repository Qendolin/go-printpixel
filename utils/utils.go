package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/tools/go/packages"
)

type BindingClosure func()

var rootPath string

func ResolveModulePath(path string) (string, error) {
	if rootPath == "" {
		if _, p, _, ok := runtime.Caller(0); ok {
			pkgs, err := packages.Load(&packages.Config{
				Mode: packages.NeedName | packages.NeedModule,
			}, p)
			if err == nil && len(pkgs) > 0 && pkgs[0].Module != nil && pkgs[0].Module.Dir != "" {
				rootPath = pkgs[0].Module.Dir
				return ResolveModulePath(path)
			}
		}
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		rootPath = wd
	}
	return filepath.Join(rootPath, path), nil
}

func MustResolveModulePath(path string) string {
	absPath, err := ResolveModulePath(path)
	if err != nil {
		panic(err)
	}
	return absPath
}

func MustResolvePath(path string) string {
	absPath, err := ResolveModulePath(path)
	if err != nil {
		panic(err)
	}
	return absPath
}

func ResolvePath(path string) (string, error) {
	if !strings.HasPrefix(path, "res://") {
		return filepath.Abs(path)
	}
	return ResolveModulePath(path[6:])
}

func NullTerm(str string) string {
	if str[len(str)-1] == '\x00' {
		return str
	} else {
		return str + "\x00"
	}
}
