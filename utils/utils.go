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

// Resolves a path relative to the current module or the working directory if the module is not avaliable
// Won't return a path thats outside of the root folder, e.g.: "./../" will be changed to just "./"
func ResolveModulePath(path string) (string, error) {
	if rootPath == "" {
		if _, p, _, ok := runtime.Caller(0); ok {
			pkgs, err := packages.Load(&packages.Config{
				Mode: packages.NeedName | packages.NeedModule,
			}, p)
			if err == nil && len(pkgs) > 0 && pkgs[0].Module != nil && pkgs[0].Module.Dir != "" {
				rootPath = filepath.Clean(pkgs[0].Module.Dir)
				return ResolveModulePath(path)
			}
		}
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		rootPath = filepath.Clean(wd)
	}

	path = filepath.Join(rootPath, path)
	if !strings.HasPrefix(path, rootPath) {
		path = rootPath
	}
	return path, nil
}

func MustResolveModulePath(path string) string {
	absPath, err := ResolveModulePath(path)
	if err != nil {
		panic(err)
	}
	return absPath
}

func MustResolvePath(path string) string {
	absPath, err := ResolvePath(path)
	if err != nil {
		panic(err)
	}
	return absPath
}

// ResolvePath acts like filepath.Abs with support for aliases.
//
// A path starting with @mod wil resolve to the module root or the wd of the binary.
// An extra @ can be used for escaping.
func ResolvePath(path string) (string, error) {
	if strings.HasPrefix(path, "@mod/") {
		return ResolveModulePath(path[4:])
	}
	if path == "@mod" {
		return ResolveModulePath("")
	}
	if strings.HasPrefix(path, "@@") {
		path = path[1:]
	}
	return filepath.Abs(path)
}

func NullTerm(str string) string {
	if str[len(str)-1] == '\x00' {
		return str
	} else {
		return str + "\x00"
	}
}
