package main

import (
	"fmt"
	"log"

	"github.com/Qendolin/go-printpixel/core"
	"github.com/Qendolin/go-printpixel/core/glw"
	"github.com/Qendolin/go-printpixel/pkg/scene"
	"github.com/Qendolin/go-printpixel/pkg/window"
	"github.com/go-gl/gl/v3.3-core/gl"
)

const (
	Width  = 1024
	Height = 1024
)

func main() {
	win := setup()

	g := &scene.Graphic{
		Texture: core.MustNewTexture2D(core.InitEmpty(Width, Height, 0), 0),
	}
	win.Child = g

	life := NewLife(Width, Height)
	win.BeforeUpdate = func() {
		life.Step()
		tex := life.Texture()
		g.Texture.Bind(0)
		if err := g.Texture.WriteBytes(0, 0, 0, Width, Height, gl.RGB, tex); err != nil {
			panic(err)
		}
		fmt.Printf("%10s, %.1ffps\n", win.GlWindow.Delta(), 1/win.GlWindow.Delta().Seconds())
	}

	win.Layout()
	win.Run()
	win.Close()
}

func setup() *window.Window {
	conf := window.SimpleConfig{
		Width:     Width,
		Height:    Height,
		FixedSize: true,
		NoVsync:   true,
		Debug:     true,
		Title:     "Game of Life Example",
		DebugHandler: func(err glw.DebugMessage) {
			if err.Critical {
				log.Fatalf("%v\n%v", err, err.Stack)
			}
			log.Printf("%v\n", err)
		},
	}
	win, err := window.New(conf)
	panicIf(err)

	return win
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
