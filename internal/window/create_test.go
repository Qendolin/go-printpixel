// +build !headless

package window

import (
	"testing"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.M) {
	err := Init()
	if err != nil {
		panic(err)
	}
	t.Run()
	Terminate()
}

func TestCreateWindowNormal(t *testing.T) {
	hints := NewHints()
	hints.ContextVersionMajor.Value = 3
	hints.ContextVersionMinor.Value = 2
	hints.Visible.Value = false
	win, err := NewWindow(hints, "Test Window", 800, 450, nil)
	assert.NoError(t, err)
	assert.NotNil(t, win)
	w, h := win.GetSize()
	assert.Equal(t, 800, w)
	assert.Equal(t, 450, h)

	win.Destroy()
}

func TestCreateWindowMaximized(t *testing.T) {
	hints := NewHints()
	hints.ContextVersionMajor.Value = 3
	hints.ContextVersionMinor.Value = 2
	hints.Maximized.Value = true
	hints.Visible.Value = false
	win, err := NewWindow(hints, "Test Window", 800, 450, nil)
	assert.NoError(t, err)
	assert.NotNil(t, win)
	w, h := win.GetSize()
	assert.True(t, 800 == w || 450 == h)
	win.Destroy()
}

func TestCreateWindowScaledToMon(t *testing.T) {
	monitor := glfw.GetPrimaryMonitor()
	vidMode := monitor.GetVideoMode()

	hints := NewHints()
	hints.ContextVersionMajor.Value = 3
	hints.ContextVersionMinor.Value = 2
	hints.Maximized.Value = true
	hints.ScaleToMonitor.Value = true
	hints.Visible.Value = false
	win, err := NewWindow(hints, "Test Window", 800, 450, nil)
	assert.NoError(t, err)
	assert.NotNil(t, win)
	w, h := win.GetSize()
	assert.Equal(t, vidMode.Width, w)
	assert.Equal(t, vidMode.Height, h)

	win.Destroy()
}
