package canvas_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Qendolin/go-printpixel/internal/canvas"
	"github.com/Qendolin/go-printpixel/internal/context"
	"github.com/Qendolin/go-printpixel/internal/window"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func TestCanvasQuad(t *testing.T) {

	os.Chdir("../..")

	err := context.InitGlfw()
	if err != nil {
		panic(err)
	}
	defer context.Terminate()

	hints := window.NewHints()
	hints.ContextVersionMajor.Value = 3
	hints.ContextVersionMinor.Value = 2
	hints.Visible.Value = false
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
			fmt.Printf("%v\n", err)
		}
	}()
	err = context.InitGl(cfg)
	if err != nil {
		panic(err)
	}

	canvas := canvas.NewCanvas(1600, 800)
	for !win.ShouldClose() {
		canvas.Draw()
		glfw.PollEvents()
	}
}
