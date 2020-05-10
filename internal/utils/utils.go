package utils

import (
	"os"
	"path/filepath"
)

type BindingClosure func() []func()

var rootPath string

func ResolveModulePath(path string) (absPath string, err error) {
	if rootPath == "" {
		var wd string
		wd, err = os.Getwd()
		if err != nil {
			panic(err)
		}
		rootPath = filepath.Join(wd, "../..")
	}
	if err != nil {
		return
	}
	absPath = filepath.Join(rootPath, path)
	return
}

func MustResolveModulePath(path string) string {
	absPath, err := ResolveModulePath(path)
	if err != nil {
		panic(err)
	}
	return absPath
}

func NullTerm(str string) string {
	if str[len(str)-1] == '\x00' {
		return str
	} else {
		return str + "\x00"
	}
}
