package shader_test

import (
	"testing"

	"github.com/Qendolin/go-printpixel/internal/canvas"
	"github.com/Qendolin/go-printpixel/internal/shader"
	"github.com/Qendolin/go-printpixel/internal/test"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func TestMain(m *testing.M) {
	test.ParseArgs()
	m.Run()
}

func TestUniformColor(t *testing.T) {
	win, close := test.NewWindow(t)
	defer close()

	prog := test.NewProgram(t, "assets/shaders/quad_uniform.vert", "assets/shaders/quad_uniform.frag")

	uColor, err := shader.NewUniform(prog, "u_color")
	if err != nil {
		t.Fatal(err)
	}

	prog.Bind()
	uColor.Set(mgl32.Vec3{0, 1, 0})

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
