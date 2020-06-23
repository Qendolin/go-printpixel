package test

import (
	"runtime"
	"testing"

	"github.com/Qendolin/go-printpixel/window"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type TestingWindow struct {
	window.Extended
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

func WrapWindow(win window.Extended) window.Extended {
	return TestingWindow{
		Extended:   win,
		isHeadless: Args.Headless,
	}
}

func NewWindow(t *testing.T) (w window.Layout, close func()) {
	runtime.LockOSThread()
	hints := window.NewHints()
	cfg := window.NewGlConfig(0)
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

	win, err := window.NewCustom("Test Window", 1600, 900, hints, nil, cfg)
	if err != nil {
		t.Fatal(err)
	}
	win.Window = WrapWindow(win.Window)

	gl.ClearColor(1, 0, 0, 1)

	return win, win.Close
}
