package glw_test

import (
	"testing"

	"github.com/Qendolin/go-printpixel/core/glw"
	"github.com/Qendolin/go-printpixel/core/test"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	test.ParseArgs()
	m.Run()
}

func TestGlfwInit(t *testing.T) {
	assert.NoError(t, glfw.Init())
	glfw.Terminate()
}

func TestCreateWindowNormal(t *testing.T) {
	winConf := glw.BasicConfig("Test Window", 800, 450, glw.DontCare, glw.DontCare)
	winConf.Visible = false
	win, err := glw.New(winConf)
	defer win.Destroy()
	require.NoError(t, err)
	require.NotNil(t, win)
	w, h := win.GetSize()
	assert.Equal(t, 800, w)
	assert.Equal(t, 450, h)
}

func TestCreateWindowAuto(t *testing.T) {
	winConf := glw.DefaultConfig()
	winConf.Title = "Test Window"
	winConf.Visible = false
	win, err := glw.New(winConf)
	defer win.Destroy()
	require.NoError(t, err)
	require.NotNil(t, win)

	monitor := glfw.GetPrimaryMonitor()
	vidMode := monitor.GetVideoMode()

	w, h := win.GetSize()
	assert.Equal(t, vidMode.Width/2, w)
	assert.Equal(t, vidMode.Height/2, h)
	x, y := win.GetPos()
	assert.Equal(t, vidMode.Width/4, x)
	assert.Equal(t, vidMode.Height/4, y)
}

func TestCreateWindowMaximized(t *testing.T) {
	hints := glw.BasicConfig("Test Window", 99, 99, glw.DontCare, glw.DontCare)
	hints.Maximized = true
	hints.Visible = true
	win, err := glw.New(hints)
	defer win.Destroy()
	require.NoError(t, err)
	require.NotNil(t, win)

	if win.HasProblem(glw.CannotMaximize) {
		t.Skipf("%v", glw.CannotMaximize)
		return
	}

	monitor := glfw.GetPrimaryMonitor()
	vidMode := monitor.GetVideoMode()

	assert.True(t, win.GetAttrib(glfw.Maximized) == glfw.True)
	w, h := win.GetSize()
	left, top, right, bot := win.GetFrameSize()
	assert.Equal(t, vidMode.Width-left-right, w)
	assert.Equal(t, vidMode.Height-top-bot, h)
}

func TestCreateWindowFullscreen(t *testing.T) {
	hints := glw.BasicConfig("Test Window", 1920, 1080, glw.DontCare, glw.DontCare)
	hints.Fullscreen = true
	hints.Visible = false
	win, err := glw.New(hints)
	defer win.Destroy()
	require.NoError(t, err)
	require.NotNil(t, win)

	monitor := glfw.GetPrimaryMonitor()
	vidMode := monitor.GetVideoMode()

	w, h := win.GetSize()
	assert.Equal(t, vidMode.Width, w)
	assert.Equal(t, vidMode.Height, h)
}

func TestGlInit(t *testing.T) {
	conf := glw.BasicConfig("Test Window", 800, 450, glw.DontCare, glw.DontCare)
	conf.Visible = false
	conf.DebugContext = true
	conf.DebugHandler = func(err glw.DebugMessage) {
		assert.NoError(t, err)
	}
	win, err := glw.New(conf)
	defer win.Destroy()
	require.NoError(t, err)
	win.MakeContextCurrent()
	assert.NoError(t, err)
	gl.GetString(gl.VERSION)
}
