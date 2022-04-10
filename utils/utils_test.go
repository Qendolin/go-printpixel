package utils_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qendolin/go-printpixel/core/test"
	"github.com/Qendolin/go-printpixel/utils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	test.ParseArgs()
	m.Run()
}

func TesResolvePath(t *testing.T) {
	_, err := utils.ResolvePath("")
	assert.NoError(t, err)
}

func TestResolveModulePath(t *testing.T) {
	path, err := utils.ResolvePath("@mod/testfile.txt")
	assert.NoError(t, err)
	assert.Equal(t, "testfile.txt", filepath.Base(path))

	absPath, err := filepath.Abs("/testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, absPath, path)

	_, err = os.Stat(path)
	assert.NoError(t, err)

	absPath, err = os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	path, err = utils.ResolvePath(absPath)
	assert.NoError(t, err)
	assert.Equal(t, absPath, path)
}

func TestRecurisveAlias(t *testing.T) {
	utils.RemoveAlias("testfile")
	utils.SetAlias("testfile", "@mod/testfile.txt")

	path, err := utils.ResolvePath("@testfile")
	assert.NoError(t, err)
	_, err = os.Stat(path)
	assert.NoError(t, err)

	utils.RemoveAlias("testalias1")
	utils.SetAlias("testalias1", "@mod/test1")
	utils.RemoveAlias("testalias3")
	utils.SetAlias("testalias3", "@testalias2/test3")
	utils.RemoveAlias("testalias2")
	utils.SetAlias("testalias2", "@testalias1/test2")

	path, err = utils.ResolvePath("@testalias3")
	assert.NoError(t, err)

	assert.True(t, strings.HasSuffix(path, filepath.Clean("test1/test2/test3")))
}

func TestPathOutsideAlias(t *testing.T) {
	_, err := utils.ResolvePath("@mod/..")
	assert.Error(t, err)

	_, err = utils.ResolvePath("@mod/../whatever")
	assert.Error(t, err)

	_, err = utils.ResolvePath("@mod/test/../../abc")
	assert.Error(t, err)
}

func TestNullTerminatedString(t *testing.T) {
	nullStr := "Test\x00"
	str := "Test"
	assert.Equal(t, nullStr, utils.NullTerm(str))
	assert.Equal(t, nullStr, utils.NullTerm(nullStr))
}
