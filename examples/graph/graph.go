package main

import (
	"log"
	"math"
	"time"

	"github.com/Qendolin/go-printpixel/core"
	"github.com/Qendolin/go-printpixel/core/glw"
	"github.com/Qendolin/go-printpixel/pkg/scene"
	"github.com/Qendolin/go-printpixel/pkg/window"
	"github.com/go-gl/gl/v3.3-core/gl"
)

func main() {
	win := setup()

	g := &scene.Graphic{
		Texture: core.MustNewTexture2D(core.InitEmpty(1600, 900, 0), 0),
	}
	win.Child = g

	start := time.Now()
	win.BeforeUpdate = func() {
		time := time.Since(start).Seconds()
		x := time
		y := math.Sin(time*math.Pi)*.5 + .5
		g.Texture.Bind(0)
		err := g.Texture.WriteBytes(0, 100+int32(x*50)%1400, 100+int32(y*500), 1, 1, gl.RGB, []byte{255, 255, 255})
		if err != nil {
			panic(err)
		}
	}

	win.Layout()
	win.Run()
	win.Close()
}

func setup() *window.Window {
	conf := window.SimpleConfig{
		Width:     1600,
		Height:    900,
		FixedSize: true,
		Debug:     true,
		Title:     "Graph Example",
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
