package context_test

import (
	"testing"

	"github.com/Qendolin/go-printpixel/internal/context"
	"github.com/Qendolin/go-printpixel/internal/test"
	"github.com/Qendolin/go-printpixel/internal/window"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	test.ParseArgs()
	m.Run()
}

func TestGlfwInit(t *testing.T) {
	err := context.InitGlfw()
	assert.NoError(t, err)
	defer context.Terminate()
}

func TestCreateWindowNormal(t *testing.T) {
	err := context.InitGlfw()
	assert.NoError(t, err)
	defer context.Terminate()

	hints := window.NewHints()
	hints.Visible.Value = false
	win, err := window.New(hints, "Test Window", 800, 450, nil)
	defer win.Destroy()
	assert.NoError(t, err)
	assert.NotNil(t, win)
	w, h := win.GetSize()
	assert.Equal(t, 800, w)
	assert.Equal(t, 450, h)
}

func TestCreateWindowMaximized(t *testing.T) {
	err := context.InitGlfw()
	assert.NoError(t, err)
	defer context.Terminate()

	hints := window.NewHints()
	hints.Maximized.Value = true
	hints.Visible.Value = false
	win, err := window.New(hints, "Test Window", 1920, 1080, nil)
	defer win.Destroy()
	assert.NoError(t, err)
	assert.NotNil(t, win)
	w, h := win.GetSize()
	assert.True(t, 1920 == w || 1080 == h)
}

func TestCreateWindowScaledToMon(t *testing.T) {
	err := context.InitGlfw()
	assert.NoError(t, err)
	defer context.Terminate()

	monitor := glfw.GetPrimaryMonitor()
	vidMode := monitor.GetVideoMode()

	hints := window.NewHints()
	hints.Maximized.Value = true
	hints.ScaleToMonitor.Value = true
	hints.Visible.Value = false
	win, err := window.New(hints, "Test Window", 800, 450, nil)
	defer win.Destroy()
	assert.NoError(t, err)
	assert.NotNil(t, win)
	w, h := win.GetSize()
	assert.Equal(t, vidMode.Width, w)
	assert.Equal(t, vidMode.Height, h)
}

func TestGlInit(t *testing.T) {
	err := context.InitGlfw()
	assert.NoError(t, err)
	defer context.Terminate()

	hints := window.NewHints()
	hints.Visible.Value = false
	win, err := window.New(hints, "Test Window", 800, 450, nil)
	defer win.Destroy()
	assert.NoError(t, err)
	win.MakeContextCurrent()
	cfg := context.NewGlConfig(0)
	cfg.Debug = true
	go func() {
		for err := range cfg.Errors {
			assert.NoError(t, err)
		}
	}()
	err = context.InitGl(cfg)
	assert.NoError(t, err)
	gl.GetString(gl.VERSION)
}
