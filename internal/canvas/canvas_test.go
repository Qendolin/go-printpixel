// +build !headless

package canvas_test

import (
	"fmt"
	"testing"

	"github.com/Qendolin/go-printpixel/internal/canvas"
	"github.com/Qendolin/go-printpixel/internal/context"
	"github.com/Qendolin/go-printpixel/internal/window"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func TestCanvasQuad(t *testing.T) {
	err := context.InitGlfw()
	if err != nil {
		panic(err)
	}
	defer context.Terminate()

	hints := window.NewHints()
	hints.ContextVersionMajor.Value = 3
	hints.ContextVersionMinor.Value = 2
	win, err := window.New(hints, "Test Window", 800, 450, nil)
	defer win.Destroy()
	if err != nil {
		panic(err)
	}
	win.MakeContextCurrent()

	cfg := context.NewGlConfig(0)
	cfg.Debug = true
	go func() {
		for err := range cfg.Errors {
			if err.Fatal {
				panic(err.Error())
			}
			fmt.Printf("%v\n", err)
		}
	}()
	err = context.InitGl(cfg)
	if err != nil {
		panic(err)
	}
	gl.ClearColor(1, 0, 0, 1)

	cnv := canvas.NewCanvas()
	for !win.ShouldClose() {
		cnv.BindFor(func() []func() {
			cnv.Draw()
			return nil
		})
		win.SwapBuffers()
		glfw.PollEvents()
	}
}
