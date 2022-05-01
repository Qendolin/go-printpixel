package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type BindingClosure func()

type aliasValue struct {
	raw   string
	cache string
}

var aliases map[string]aliasValue = map[string]aliasValue{
	"@": {},
}

func init() {
	libroot := ""
	if _, currentFile, _, ok := runtime.Caller(0); ok {
		if p, ok := FindModuleRoot(currentFile); ok {
			libroot = p
		}
	}

	if libroot == "" {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		libroot = filepath.Clean(wd)
	}
	SetAlias("lib", libroot)

	modroot := ""
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if p, ok := FindModuleRoot(wd); ok {
		modroot = p
	}

	if libroot == "" {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		modroot = filepath.Clean(wd)
	}
	SetAlias("mod", modroot)
}

func FindModuleRoot(start string) (string, bool) {
	if start == "" {
		panic("start not set")
	}
	path := filepath.Clean(start)

	// Look for enclosing go.mod.
	for {
		if fi, err := os.Stat(filepath.Join(path, "go.mod")); err == nil && !fi.IsDir() {
			return path, true
		}

		// Get the parent directory
		parent := filepath.Dir(path)
		if parent == path {
			// Reached the outer most directory
			break
		}
		path = parent
	}
	return "", false
}

// alias must not contain '/'. value can contain another alias
func SetAlias(alias, value string) {
	alias = normalizeAliasName(alias)

	if alias == "@" {
		panic("the specified alias name is empty")
	}

	aliases[alias] = aliasValue{
		raw: value,
	}

	err := buildAliasCache()
	if err != nil {
		panic(fmt.Sprintf("invalid alias path %q: %v", value, err))
	}
}

func RemoveAlias(alias string) {
	alias = normalizeAliasName(alias)
	if _, ok := aliases[alias]; !ok {
		return
	}
	delete(aliases, alias)
	buildAliasCache()
}

// Checks for an alias and returns it, including the @ symbol
func MatchAlias(path string) (string, bool) {
	if !strings.HasPrefix(path, "@") {
		return "", false
	}

	if strings.HasPrefix(path, "@@") {
		return "@", true
	}

	return strings.SplitN(path, "/", 2)[0], true
}

// Ensures single leading '@', lowercase and no '/'
func normalizeAliasName(alias string) string {
	alias = strings.TrimLeftFunc(alias, func(r rune) bool { return r == '@' })
	alias = strings.ToLower(alias)

	if strings.Contains(alias, "/") {
		panic(fmt.Sprintf("the specified alias name %q must not contain '/'", alias))
	}

	return "@" + alias
}

func buildAliasCache() error {
	for key, value := range aliases {
		absPath, err := filepath.Abs(resolveAlias(value.raw, []string{}))
		if err != nil {
			return err
		}
		value.cache = absPath
		aliases[key] = value
	}
	return nil
}

// recursivly resolves aliases. panics on cyclic references
func resolveAlias(path string, visited []string) string {
	alias, ok := MatchAlias(path)
	if !ok {
		return path
	}

	for _, a := range visited {
		if a == alias {
			panic(fmt.Sprintf("cyclic aliases: %v, %v", strings.Join(visited, ", "), alias))
		}
	}

	value, ok := aliases[alias]
	if ok {
		path = value.raw + strings.TrimPrefix(path, alias)
	} else {
		return path
	}

	visited = append(visited, alias)
	return resolveAlias(path, visited)
}

// Resolves a path relative to the current module or the working directory if the module is not avaliable
// Won't return a path thats outside of the root folder, e.g.: "./../" will be changed to just "./"
func resolveAliasPath(path string) (string, error) {
	alias, ok := MatchAlias(path)
	if !ok {
		return path, nil
	}

	value, ok := aliases[alias]
	if !ok {
		return "", fmt.Errorf("could not resovle path %q, alias %q does not exist", path, alias)
	}

	// recursive aliases must be resolved
	if missing, ok := MatchAlias(value.cache); ok {
		return "", fmt.Errorf("could not resovle alias %q (resolved to %q), alias %q is missing", alias, value.cache, missing)
	}

	path = value.cache + strings.TrimPrefix(path, alias)

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	// path must not go outside the aliased directory
	if !strings.HasPrefix(absPath, value.cache) {
		return "", fmt.Errorf("path %q ends outside the aliased path %q", absPath, value.cache)
	}

	return path, nil
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
// Using the prefix @@ will prevent alias resolution and the first @ will be removed.
func ResolvePath(path string) (abs string, err error) {
	path, err = resolveAliasPath(path)
	if err != nil {
		return "", err
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
