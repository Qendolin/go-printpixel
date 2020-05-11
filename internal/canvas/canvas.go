package canvas

import (
	"github.com/Qendolin/go-printpixel/internal/data"
	"github.com/Qendolin/go-printpixel/internal/shader"
	"github.com/Qendolin/go-printpixel/internal/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
)

var (
	quadVertices   = []float32{1, -1, 1, 1, -1, -1, -1, 1}
	quadVao        *data.Vao
	quadShaderProg *shader.Program
	isInit         = false
)

func _init() {
	if isInit {
		return
	}
	isInit = true

	quadVao = data.NewVao()
	quadVao.BindFor(func() (defered []func()) {
		quadVbo := data.NewVbo()
		quadVbo.Bind(gl.ARRAY_BUFFER)
		quadVbo.WriteStatic(quadVertices)
		quadVbo.MustLayout(0, 2, float32(0), false, 0)

		defered = append(defered, func() {
			quadVbo.Unbind(gl.ARRAY_BUFFER)
		})
		return
	})

	quadVertShader, err := shader.NewShaderFromPath("assets/shaders/quad_tex.vert", shader.TypeVertex)
	if err != nil {
		panic(err)
	}

	quadFragShader, err := shader.NewShaderFromPath("assets/shaders/quad_tex.frag", shader.TypeFragment)
	if err != nil {
		panic(err)
	}

	quadShaderProg, err = shader.NewProgram(quadVertShader, quadFragShader)
	if err != nil {
		panic(err)
	}

}

type Canvas struct {
	Program shader.Program
}

func NewCanvas() *Canvas {
	if !isInit {
		_init()
	}
	return &Canvas{*quadShaderProg}
}

func NewCanvasWithProgram(prog shader.Program) *Canvas {
	if !isInit {
		_init()
	}
	return &Canvas{prog}
}

func (canvas *Canvas) Bind() {
	quadVao.Bind()
	canvas.Program.Bind()
}

func (canvas *Canvas) Unbind() {
	quadVao.Unbind()
	canvas.Program.Unbind()
}

func (canvas *Canvas) BindFor(context utils.BindingClosure) {
	canvas.Bind()
	defered := context()
	canvas.Unbind()
	for _, deferedFunc := range defered {
		deferedFunc()
	}
}

func (canvas *Canvas) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
}
