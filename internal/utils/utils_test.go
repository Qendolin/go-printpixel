package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveModulePath(t *testing.T) {
	_, err := ResolveModulePath("")
	assert.NoError(t, err)
}

func TestResolvePath(t *testing.T) {
	modPath := MustResolveModulePath("")
	path, err := ResolvePath("")
	assert.NoError(t, err)
	assert.Equal(t, modPath, path)

	absPath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	path, err = ResolvePath(absPath)
	assert.NoError(t, err)
	assert.Equal(t, absPath, path)
}

func TestNullTerminatedString(t *testing.T) {
	nullStr := "Test\x00"
	str := "Test"
	assert.Equal(t, nullStr, NullTerm(str))
	assert.Equal(t, nullStr, NullTerm(nullStr))
}
