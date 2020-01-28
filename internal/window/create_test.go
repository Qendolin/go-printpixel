package window

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.M) {
	Init()
	t.Run()
	Terminate()
}

func TestCreateWindow(t *testing.T) {
	hints := NewHints()
	hints.ContextVersionMajor.Value = 3
	hints.ContextVersionMinor.Value = 2
	win, err := NewWindow(hints, "Test Window", 800, 450, nil)
	assert.NoError(t, err)
	assert.NotNil(t, win)
}
