package test

import (
	"runtime"
	"testing"

	_ "image/jpeg"
	_ "image/png"

	"github.com/Qendolin/go-printpixel/core/glw"
	"github.com/Qendolin/go-printpixel/pkg/window"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type TestingWindow struct {
	glw.Window
	closeCheckCount int
	isHeadless      bool
}

func (win TestingWindow) ShouldClose() bool {
	win.closeCheckCount++
	if win.closeCheckCount >= 10 || win.isHeadless {
		return true
	}
	return win.Window.ShouldClose()
}

func WrapWindow(win glw.Window) glw.Window {
	return TestingWindow{
		Window:     win,
		isHeadless: Args.Headless,
	}
}

func NewWindow(t *testing.T) (w *window.Window, close func()) {
	runtime.LockOSThread()
	conf := glw.BasicConfig("Test Window | "+t.Name(), 1600, 900, glw.DontCare, glw.DontCare)
	conf.DebugContext = true
	conf.DebugHandler = func(err glw.DebugMessage) {
		if err.Critical {
			t.Error(err, "\n"+err.Stack)
		} else {
			t.Log(err)
		}
	}
	win, err := window.NewCustom(conf)
	if err != nil {
		t.Fatal(err)
	}
	win.GlWindow = WrapWindow(win.GlWindow)

	gl.ClearColor(0, 0, 0, 1)

	return win, win.Close
}
