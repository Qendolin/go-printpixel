package test

import (
	"runtime"
	"testing"

	"github.com/Qendolin/go-printpixel/core/glcontext"
	"github.com/Qendolin/go-printpixel/core/glwindow"
	"github.com/Qendolin/go-printpixel/pkg/window"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type TestingWindow struct {
	glwindow.Extended
	closeCheckCount int
	isHeadless      bool
}

func (win TestingWindow) ShouldClose() bool {
	win.closeCheckCount++
	if win.closeCheckCount >= 10 || win.isHeadless {
		return true
	}
	return win.Extended.ShouldClose()
}

func WrapWindow(win glwindow.Extended) glwindow.Extended {
	return TestingWindow{
		Extended:   win,
		isHeadless: Args.Headless,
	}
}

func NewWindow(t *testing.T) (w *window.Window, close func()) {
	runtime.LockOSThread()
	hints := glwindow.NewHints()
	cfg := glcontext.NewGlConfig(0)
	cfg.Debug = true
	go func() {
		for err := range cfg.Errors {
			if err.Fatal {
				t.Error(err, "\n"+err.Stack)
			} else {
				t.Log(err)
			}
		}
	}()

	win, err := window.NewCustom("Test Window | "+t.Name(), 1600, 900, hints, nil, cfg)
	if err != nil {
		t.Fatal(err)
	}
	win.GlWindow = WrapWindow(win.GlWindow)

	gl.ClearColor(1, 0, 0, 1)

	return win, win.Close
}
