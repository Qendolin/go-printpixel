package window

import (
	"testing"

	"github.com/go-gl/glfw/v3.2/glfw"
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
	win, err := NewWindow(hints, "Test Window", 800, 450, nil)
	assert.NoError(t, err)
	assert.NotNil(t, win)

	win.Destroy()
}

func TestCreateWindowMaximized(t *testing.T) {
	hints := NewHints()
	hints.ContextVersionMajor.Value = 3
	hints.ContextVersionMinor.Value = 2
	hints.Maximized.Value = true
	win, err := NewWindow(hints, "Test Window", 800, 450, nil)
	assert.NoError(t, err)
	assert.NotNil(t, win)

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
	win, err := NewWindow(hints, "Test Window", 800, 450, nil)
	assert.NoError(t, err)
	assert.NotNil(t, win)
	w, h := win.GetSize()
	assert.Equal(t, vidMode.Width, w, h)

	win.Destroy()
}
