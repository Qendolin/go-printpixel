package main

import (
	"fmt"
	"runtime"

	"github.com/Qendolin/go-printpixel/core/glw"
	"github.com/Qendolin/go-printpixel/experiments/Input/input"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	runtime.LockOSThread()
	glwConf := glw.BasicConfig("Input Test", 1600/2, 900/2, glw.DontCare, glw.DontCare)
	glwConf.DebugContext = true
	win, err := glw.New(glwConf)
	panicIf(err)

	input.Default.Bind(win.GetGLFWWindow())
	err = input.Default.AddTrigger(input.Combo("ModControl", "KeySpace"), "print_hello")
	panicIf(err)
	err = input.Default.AddTrigger(input.Combo("ModControl", "KeyK"), "switch_mode")
	panicIf(err)

	altMode := map[input.Trigger]input.Action{
		input.Combo("ModControl", "KeyX"): "win_exit",
		input.Combo("ModControl", "KeyF"): "win_fullscreen",
	}

	var (
		isFullscreen = false
	)

	input.Default.On("print_hello", func(ae input.ActionEvent) {
		fmt.Println("Hello World!")
	})
	input.Default.On("switch_mode", func(ae input.ActionEvent) {
		input.Default.SetOverride(altMode, false)
	})
	input.Default.On("win_exit", func(ae input.ActionEvent) {
		win.SetShouldClose(true)
	})
	input.Default.On("win_fullscreen", func(ae input.ActionEvent) {
		if isFullscreen {
			win.SetMonitor(nil, 100, 100, 1600/2, 900/2, 0)
		} else {
			mon := glfw.GetPrimaryMonitor()
			vidMode := mon.GetVideoMode()
			win.SetMonitor(mon, 0, 0, vidMode.Width, vidMode.Height, 0)
		}
		isFullscreen = !isFullscreen
	})

	for !win.ShouldClose() {
		input.Default.Update()

		win.SwapBuffers()
	}
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
