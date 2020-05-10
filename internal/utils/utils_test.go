package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveModulePath(t *testing.T) {
	_, err := ResolveModulePath("")
	assert.NoError(t, err)
}

func TestNullTerminatedString(t *testing.T) {
	nullStr := "Test\x00"
	str := "Test"
	assert.Equal(t, nullStr, NullTerm(str))
	assert.Equal(t, nullStr, NullTerm(nullStr))
}
