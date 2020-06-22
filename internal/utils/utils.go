package utils

import (
	"go/build"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

type BindingClosure func()

var rootPath string

func ResolveModulePath(path string) (string, error) {
	if rootPath == "" {
		pkgs, err := packages.Load(&packages.Config{
			Mode: packages.NeedName | packages.NeedModule,
		}, "github.com/Qendolin/go-printpixel")
		if err == nil && len(pkgs) > 0 && pkgs[0].Module != nil && pkgs[0].Module.Dir != "" {
			rootPath = pkgs[0].Module.Dir
			return ResolveModulePath(path)
		}

		pkg, err := build.Import("github.com/Qendolin/go-printpixel", "/", build.FindOnly)
		if err == nil && pkg.Dir != "" {
			rootPath = pkg.Dir
			return ResolveModulePath(path)
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
	if filepath.IsAbs(path) {
		return path, nil
	}
	return ResolveModulePath(path)
}

func NullTerm(str string) string {
	if str[len(str)-1] == '\x00' {
		return str
	} else {
		return str + "\x00"
	}
}
