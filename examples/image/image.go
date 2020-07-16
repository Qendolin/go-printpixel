package main

import (
	"log"
	"os"

	_ "image/png"

	"github.com/Qendolin/go-printpixel/core"
	"github.com/Qendolin/go-printpixel/core/glcontext"
	"github.com/Qendolin/go-printpixel/core/glwindow"
	"github.com/Qendolin/go-printpixel/pkg/scene"
	"github.com/Qendolin/go-printpixel/pkg/window"
)

func main() {
	win := setup()

	img, err := os.Open("./image.png")
	panicIf(err)
	defer img.Close()
	tex, err := core.NewTexture2DFromFile(img)
	panicIf(err)

	win.Child = &scene.Aspect{
		Child: &scene.Graphic{
			Texture: tex,
		},
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
	win, err := window.New("Image Example", &cfg)
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
