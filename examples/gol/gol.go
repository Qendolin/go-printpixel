package main

import (
	"fmt"
	"log"

	"github.com/Qendolin/go-printpixel/layout"
	"github.com/Qendolin/go-printpixel/window"
	"github.com/go-gl/gl/v3.3-core/gl"
)

const (
	Width  = 1024
	Height = 1024
)

func main() {
	win := setup()

	view := layout.NewViewer()
	view.Target.Texture.Bind(0)
	view.Target.Texture.ApplyDefaults()
	view.Target.Texture.AllocEmpty(0, gl.RGB, Width, Height, gl.RGB)
	win.Child = view

	life := NewLife(Width, Height)
	win.BeforeUpdate = func() {
		life.Step()
		tex := life.Texture()
		view.Target.Texture.Bind(0)
		view.Target.Texture.WriteBytes(tex, 0, 0, 0, Width, Height, gl.RGB)
		view.Draw()
		fmt.Println(win.Window.Delta())
	}

	win.Run()
	win.Close()
}

func setup() window.Layout {
	cfg := window.SimpleConfig{
		Width:        Width,
		Height:       Height,
		Unresizeable: true,
		NoVsync:      true,
		Debug:        true,
	}
	win, err := window.New("Game of Life Example", cfg)
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
