package shader_test

import (
	"fmt"
	"testing"

	"github.com/Qendolin/go-printpixel/internal/canvas"
	"github.com/Qendolin/go-printpixel/internal/context"
	"github.com/Qendolin/go-printpixel/internal/shader"
	"github.com/Qendolin/go-printpixel/internal/window"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func TestUniformColor(t *testing.T) {
	err := context.InitGlfw()
	if err != nil {
		t.Error(err)
	}
	defer context.Terminate()

	hints := window.NewHints()
	hints.ContextVersionMajor.Value = 3
	hints.ContextVersionMinor.Value = 2
	win, err := window.New(hints, "Test Window", 800, 450, nil)
	defer win.Destroy()
	if err != nil {
		t.Error(err)
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
		t.Error(err)
	}
	gl.ClearColor(1, 0, 0, 1)

	vs, err := shader.NewShaderFromModulePath("assets/shaders/quad_uniform.vert", shader.TypeVertex)
	if err != nil {
		t.Error(err)
	}

	fs, err := shader.NewShaderFromModulePath("assets/shaders/quad_uniform.frag", shader.TypeFragment)
	if err != nil {
		t.Error(err)
	}
	prog, err := shader.NewProgram(vs, fs)
	if err != nil {
		t.Error(err)
	}

	uColor, err := shader.NewUniform(*prog, "u_color")
	if err != nil {
		t.Error(err)
	}

	prog.Bind()
	uColor.Set(mgl32.Vec3{0, 1, 0})

	cnv := canvas.NewCanvasWithProgram(*prog)
	for !win.ShouldClose() {
		cnv.BindFor(func() []func() {
			cnv.Draw()
			return nil
		})
		win.SwapBuffers()
		glfw.PollEvents()
	}
}
