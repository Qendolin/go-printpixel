package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "image/png"

	"github.com/Qendolin/go-printpixel/core"
	"github.com/Qendolin/go-printpixel/core/data"
	"github.com/Qendolin/go-printpixel/core/glcontext"
	"github.com/Qendolin/go-printpixel/core/glwindow"
	"github.com/Qendolin/go-printpixel/pkg/scene"
	"github.com/Qendolin/go-printpixel/pkg/window"
	"github.com/go-gl/gl/v3.3-core/gl"
)

const Stars = 400

func main() {
	win := setup()
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	tex := core.LoadTexture("./star.png")
	tex.Bind(0)
	tex.FilterMode(data.FilterLinearMipMapLinear, data.FilterLinear)
	tex.GenerateMipmap()

	starStack := &scene.Stack{
		Children: make([]scene.Layoutable, Stars),
	}
	win.Child = &scene.Center{
		Child: starStack,
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < Stars; i++ {
		size := 13 + 10*rand.Float32()
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
				},
			},
		}
	}

	win.GlWindow.SetSizeCallback(func(_ glwindow.Extended, _ int, _ int) {
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
		Width:        1600,
		Height:       900,
		Debug:        true,
		Unresizeable: false,
	}
	win, err := window.New("Stars Example", &cfg)
	panicIf(err)

	go handleErrors(cfg.Errors())
	return win
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func handleErrors(errs <-chan glcontext.Error) {
	for err := range errs {
		if err.Fatal {
			log.Fatalf("%v\n%v", err, err.Stack)
		}
		log.Printf("%v\n", err)
	}
}
