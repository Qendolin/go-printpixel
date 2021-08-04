package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "image/png"

	"github.com/Qendolin/go-printpixel/core"
	"github.com/Qendolin/go-printpixel/core/data"
	"github.com/Qendolin/go-printpixel/core/glw"
	"github.com/Qendolin/go-printpixel/pkg/scene"
	"github.com/Qendolin/go-printpixel/pkg/window"
	"github.com/go-gl/gl/v3.3-core/gl"
)

const Stars = 400

func main() {
	win := setup()
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	tex := core.MustNewTexture2D(core.
		InitPaths(0, "./star.png").
		WithFilters(data.FilterLinearMipMapLinear, data.FilterLinear).
		WithRequiredLevels().
		WithGeneratedMipMap(), 0)
	tex.Bind(0)

	starStack := &scene.Stack{
		Children: make([]scene.Layoutable, Stars),
	}
	win.Child = scene.Centered(starStack)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < Stars; i++ {
		size := 20 + 10*rand.Float32()
		starStack.Children[i] = &scene.Absolute{
			Unit: scene.Percent,
			DX:   rand.Float32(),
			DY:   rand.Float32(),
			W:    1,
			H:    1,
			Child: &scene.Absolute{
				W: size,
				H: size,
				Child: &scene.Graphic{
					Texture: tex,
					Alpha:   true,
				},
			},
		}
	}

	win.GlWindow.SetSizeCallback(func(_ glw.Window, _ int, _ int) {
		win.Layout()
	})

	win.AfterUpdate = func() {
		fmt.Printf("%10s, %.1ffps\n", win.GlWindow.Delta(), 1/win.GlWindow.Delta().Seconds())
	}

	win.Layout()
	win.Run()
	win.Close()
}

func setup() *window.Window {
	cfg := window.SimpleConfig{
		Width:  1600,
		Height: 900,
		Debug:  true,
		Title:  "Stars Example",
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
