// +build !headless

package canvas_test

import (
	"testing"

	"github.com/Qendolin/go-printpixel/internal/canvas"
	"github.com/Qendolin/go-printpixel/internal/test"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func TestCanvasQuad(t *testing.T) {
	win, close := test.NewWindow(t)
	defer close()

	prog := test.NewProgram(t, "assets/shaders/quad_uv.vert", "assets/shaders/quad_uv.frag")

	cnv := canvas.NewCanvasWithProgram(prog)
	for !win.ShouldClose() {
		cnv.BindFor(func() []func() {
			cnv.Draw()
			return nil
		})
		win.SwapBuffers()
		glfw.PollEvents()
	}
}
