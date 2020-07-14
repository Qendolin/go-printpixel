package main

import (
	"log"
	"os"

	_ "image/png"

	"github.com/Qendolin/go-printpixel/core/glcontext"
	"github.com/Qendolin/go-printpixel/core/glwindow"
	"github.com/Qendolin/go-printpixel/pkg/layout"
	"github.com/Qendolin/go-printpixel/pkg/window"
	"github.com/go-gl/gl/v3.3-core/gl"
)

func main() {
	win := setup()

	img, err := os.Open("./image.png")
	panicIf(err)
	defer img.Close()

	g := layout.NewGraphic()
	g.Texture.Bind(0)
	g.Texture.ApplyDefaults()
	err = g.Texture.AllocFile(img, 0, gl.RGBA, gl.RGBA)
	panicIf(err)
	win.Child = &layout.Aspect{
		Child: g,
		Ratio: 640 / 640,
	}

	win.GlWindow.SetSizeCallback(func(_ glwindow.Extended, _ int, _ int) {
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
	}
	win, err := window.New("Image Example", cfg)
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
