package main

import (
	"fmt"
	"log"

	"github.com/Qendolin/go-printpixel/core/glcontext"
	"github.com/Qendolin/go-printpixel/pkg/layout"
	"github.com/Qendolin/go-printpixel/pkg/window"
	"github.com/go-gl/gl/v3.3-core/gl"
)

const (
	Width  = 1024
	Height = 1024
)

func main() {
	win := setup()

	g := layout.NewGraphic()
	g.Texture.Bind(0)
	g.Texture.ApplyDefaults()
	g.Texture.AllocEmpty(0, gl.RGB, Width, Height, gl.RGB)
	win.Child = g

	life := NewLife(Width, Height)
	win.BeforeUpdate = func() {
		life.Step()
		tex := life.Texture()
		g.Texture.Bind(0)
		g.Texture.WriteBytes(tex, 0, 0, 0, Width, Height, gl.RGB)
		fmt.Println(win.GlWindow.Delta())
	}

	win.Layout()
	win.Run()
	win.Close()
}

func setup() *window.Window {
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

func handleErrors(errs <-chan glcontext.GlError) {
	for err := range errs {
		if err.Fatal {
			log.Fatalf("%v\n%v", err, err.Stack)
		}
		log.Printf("%v\n", err)
	}
}
