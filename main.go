package main

import (
	_ "image/png"

	"github.com/Qendolin/go-printpixel/pkg/scene"
	"github.com/Qendolin/go-printpixel/pkg/window"
)

func main() {
	win, err := window.New(window.SimpleConfig{Title: "Hello World!", Height: 900, Width: 1600, FixedSize: true})
	if err != nil {
		panic(err)
	}
	win.Child = &scene.Aspect{
		Child: scene.LoadGraphic("@lib/assets/textures/hello.png"),
		//lint:ignore SA4000 So what?
		Ratio: 640 / 640,
	}
	win.Layout()

	win.Run()
}
