package main

import (
	"log"
	"math"
	"time"

	"github.com/Qendolin/go-printpixel/layout"
	"github.com/Qendolin/go-printpixel/window"
	"github.com/go-gl/gl/v3.3-core/gl"
)

func main() {
	win := setup()

	view := layout.NewViewer()
	view.Target.Texture.Bind(0)
	view.Target.Texture.ApplyDefaults()
	view.Target.Texture.AllocEmpty(0, gl.RGB, 1600, 900, gl.RGB)
	win.Child = view

	start := time.Now()
	win.BeforeUpdate = func() {
		time := time.Since(start).Seconds()
		x := time
		y := math.Sin(time*math.Pi)*.5 + .5
		view.Target.Texture.Bind(0)
		view.Target.Texture.WriteBytes([]byte{255, 255, 255}, 0, 100+int32(x*50)%1400, 100+int32(y*500), 1, 1, gl.RGB)
		view.Draw()
	}

	win.Run()
	win.Close()
}

func setup() window.Layout {
	cfg := window.SimpleConfig{
		Width:        1600,
		Height:       900,
		Unresizeable: true,
		Debug:        true,
	}
	win, err := window.New("Graph Example", cfg)
	panicIf(err)

	go handleErrors(cfg.Errors())
	return win
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func handleErrors(errs <-chan window.GlError) {
	for err := range errs {
		if err.Fatal {
			log.Fatalf("%v\n%v", err, err.Stack)
		}
		log.Printf("%v\n", err)
	}
}
