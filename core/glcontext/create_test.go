package glcontext_test

import (
	"testing"

	"github.com/Qendolin/go-printpixel/core/glcontext"
	"github.com/Qendolin/go-printpixel/core/glwindow"
	"github.com/Qendolin/go-printpixel/core/test"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	test.ParseArgs()
	m.Run()
}

func TestGlfwInit(t *testing.T) {
	err := glcontext.InitGlfw()
	assert.NoError(t, err)
	defer glcontext.Terminate()
}

func TestCreateWindowNormal(t *testing.T) {
	err := glcontext.InitGlfw()
	assert.NoError(t, err)
	defer glcontext.Terminate()

	hints := glwindow.NewHints()
	hints.Visible.Value = false
	win, err := glwindow.New(hints, "Test Window", 800, 450, nil)
	defer win.Destroy()
	assert.NoError(t, err)
	assert.NotNil(t, win)
	w, h := win.GetSize()
	assert.Equal(t, 800, w)
	assert.Equal(t, 450, h)
}

func TestCreateWindowMaximized(t *testing.T) {
	err := glcontext.InitGlfw()
	assert.NoError(t, err)
	defer glcontext.Terminate()

	monitor := glfw.GetPrimaryMonitor()
	vidMode := monitor.GetVideoMode()

	hints := glwindow.NewHints()
	hints.Maximized.Value = true
	hints.Visible.Value = false
	win, err := glwindow.New(hints, "Test Window", 1, 1, nil)
	defer win.Destroy()
	assert.NoError(t, err)
	assert.NotNil(t, win)
	w, h := win.GetSize()
	assert.True(t, w == vidMode.Width && h == vidMode.Height)
}

func TestCreateWindowFullscreen(t *testing.T) {
	err := glcontext.InitGlfw()
	assert.NoError(t, err)
	defer glcontext.Terminate()

	monitor := glfw.GetPrimaryMonitor()
	vidMode := monitor.GetVideoMode()

	hints := glwindow.NewHints()
	hints.Fullscreen.Value = true
	hints.Visible.Value = false
	win, err := glwindow.New(hints, "Test Window", 1920, 1080, nil)
	defer win.Destroy()
	assert.NoError(t, err)
	assert.NotNil(t, win)
	w, h := win.GetSize()
	assert.Equal(t, vidMode.Width, w)
	assert.Equal(t, vidMode.Height, h)
}

func TestGlInit(t *testing.T) {
	err := glcontext.InitGlfw()
	assert.NoError(t, err)
	defer glcontext.Terminate()

	hints := glwindow.NewHints()
	hints.Visible.Value = false
	win, err := glwindow.New(hints, "Test Window", 800, 450, nil)
	defer win.Destroy()
	assert.NoError(t, err)
	win.MakeContextCurrent()
	cfg := glcontext.NewGlConfig(0)
	cfg.Debug = true
	go func() {
		for err := range cfg.Errors {
			assert.NoError(t, err)
		}
	}()
	err = glcontext.InitGl(cfg)
	assert.NoError(t, err)
	gl.GetString(gl.VERSION)
}
