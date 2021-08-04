package main

import (
	"log"

	_ "image/png"

	"github.com/Qendolin/go-printpixel/core/glw"
	"github.com/Qendolin/go-printpixel/pkg/scene"
	"github.com/Qendolin/go-printpixel/pkg/window"
)

func main() {
	win := setup()

	win.Child = &scene.Aspect{
		Child: scene.LoadGraphic("./image.png"),
		//lint:ignore SA4000 So what?
		Ratio: 640 / 640,
	}

	win.GlWindow.SetSizeCallback(func(_ glw.Window, _ int, _ int) {
		win.Layout()
	})

	win.Layout()
	win.Run()
	win.Close()
}

func setup() *window.Window {
	cfg := window.SimpleConfig{
		Width:  1600,
		Height: 900,
		Debug:  true,
		Title:  "Image Example",
		DebugHandler: func(err glw.DebugMessage) {
			if err.Critical {
				log.Fatalf("%v\n%v", err, err.Stack)
			}
			log.Printf("%v\n", err)
		},
	}
	win, err := window.New(cfg)
	panicIf(err)

	return win
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
