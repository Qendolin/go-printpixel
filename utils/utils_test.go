package utils_test

import (
	"os"
	"testing"

	"github.com/Qendolin/go-printpixel/core/test"
	"github.com/Qendolin/go-printpixel/utils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	test.ParseArgs()
	m.Run()
}

func TestResolveModulePath(t *testing.T) {
	_, err := utils.ResolveModulePath("")
	assert.NoError(t, err)
}

func TestResolvePath(t *testing.T) {
	modPath := utils.MustResolveModulePath("")
	path, err := utils.ResolvePath("res://")
	assert.NoError(t, err)
	assert.Equal(t, modPath, path)

	absPath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	path, err = utils.ResolvePath(absPath)
	assert.NoError(t, err)
	assert.Equal(t, absPath, path)
}

func TestNullTerminatedString(t *testing.T) {
	nullStr := "Test\x00"
	str := "Test"
	assert.Equal(t, nullStr, utils.NullTerm(str))
	assert.Equal(t, nullStr, utils.NullTerm(nullStr))
}
