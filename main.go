package main

import (
	_ "image/png"

	"github.com/Qendolin/go-printpixel/pkg/scene"
	"github.com/Qendolin/go-printpixel/pkg/window"
)

func main() {
	win, err := window.New("Hello World!", &window.SimpleConfig{Height: 900, Width: 1600, Unresizeable: true})
	if err != nil {
		panic(err)
	}
	win.Child = &scene.Aspect{
		Child: scene.LoadGraphic("@mod/assets/textures/hello.png"),
		Ratio: 640 / 640,
	}
	win.Layout()

	win.Run()
}
