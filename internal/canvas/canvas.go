package canvas

import (
	"github.com/Qendolin/go-printpixel/internal/data"
	"github.com/Qendolin/go-printpixel/internal/shader"
	"github.com/Qendolin/go-printpixel/internal/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
)

var (
	quadVertices = []float32{1, -1, 1, 1, -1, -1, -1, 1}
)

type Canvas struct {
	Program shader.Program
	Texture *data.Texture2D
	quad    data.Vao
}

func NewCanvas() *Canvas {
	vs, err := shader.NewShaderFromPath("assets/shaders/quad_tex.vert", shader.TypeVertex)
	if err != nil {
		panic(err)
	}

	fs, err := shader.NewShaderFromPath("assets/shaders/quad_tex.frag", shader.TypeFragment)
	if err != nil {
		panic(err)
	}

	quadShaderProg, err := shader.NewProgram(vs, fs)
	if err != nil {
		panic(err)
	}

	fs.Destroy()
	vs.Destroy()
	return NewCanvasWithProgram(*quadShaderProg)
}

func NewCanvasWithProgram(prog shader.Program) *Canvas {
	quadVao := data.NewVao()
	quadVao.BindFor(func() {
		quadVbo := data.NewVbo()
		quadVbo.Bind(gl.ARRAY_BUFFER)
		quadVbo.WriteStatic(quadVertices)
		quadVbo.MustLayout(0, 2, float32(0), false, 0)

		quadVbo.Unbind(gl.ARRAY_BUFFER)
	})
	return &Canvas{Program: prog, quad: *quadVao}
}

func (canvas *Canvas) Bind() {
	canvas.quad.Bind()
	canvas.Program.Bind()
	if canvas.Texture != nil {
		canvas.Texture.Bind(0)
	}
}

func (canvas *Canvas) Unbind() {
	canvas.quad.Unbind()
	canvas.Program.Unbind()
	if canvas.Texture != nil {
		canvas.Texture.Unbind(0)
	}
}

func (canvas *Canvas) BindFor(context utils.BindingClosure) {
	canvas.Bind()
	context()
	canvas.Unbind()
}

func (canvas *Canvas) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
}

func (canvas *Canvas) Destroy() {
	canvas.Program.Destroy()
	canvas.quad.Destroy()
	if canvas.Texture != nil {
		canvas.Texture.Destroy()
	}
}
