package shader_test

import (
	"testing"

	"github.com/Qendolin/go-printpixel/core"
	"github.com/Qendolin/go-printpixel/core/shader"
	"github.com/Qendolin/go-printpixel/core/test"
	"github.com/go-gl/gl/v3.3-core/gl"
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

	uColor, err := shader.NewUniform(*prog, "u_color")
	if err != nil {
		t.Fatal(err)
	}

	prog.Bind()
	uColor.Set(mgl32.Vec3{0, 1, 0})

	core.Quad().Bind()
	for !win.ShouldClose() {
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
		win.SwapBuffers()
		glfw.PollEvents()
	}
}
