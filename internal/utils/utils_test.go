package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveModulePath(t *testing.T) {
	_, err := ResolveModulePath("")
	assert.NoError(t, err)
}
