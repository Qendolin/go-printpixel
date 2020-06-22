package test

import (
	"runtime"
	"testing"

	"github.com/Qendolin/go-printpixel/window"
	"github.com/go-gl/gl/v2.1/gl"
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
	win, err := window.New(hints, "Test Window", 1600, 900, nil)
	if err != nil {
		t.Fatal(err)
	}
	win.Window = WrapWindow(win.Window)

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
	err = win.Init(cfg)
	if err != nil {
		t.Fatal(err)
	}

	gl.ClearColor(1, 0, 0, 1)

	return win, win.Close
}
