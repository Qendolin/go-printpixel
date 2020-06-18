package test

import (
	"github.com/Qendolin/go-printpixel/internal/window"
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
