package canvas

import (
	"io/ioutil"

	"github.com/Qendolin/go-printpixel/internal"
	"github.com/Qendolin/go-printpixel/internal/shader"
	"github.com/go-gl/gl/v3.2-core/gl"
)

var (
	quadVertices   = []float32{1, -1, 1, 1, -1, -1, -1, 1}
	quadVao        *Vao
	quadShaderProg *shader.Program
	isInit         = false
)

func _init() {
	if isInit {
		return
	}

	quadVao = NewVao()
	quadVao.BindFor(func() (defered []func()) {
		quadVbo := NewVbo()
		quadVbo.Bind(gl.ARRAY_BUFFER)
		quadVbo.WriteStatic(quadVertices)
		quadVbo.MustLayout(0, 3, float32(0), false, 0)

		defered = append(defered, func() {
			quadVbo.Unbind(gl.ARRAY_BUFFER)
		})
		return
	})

	qvsSource, err := ioutil.ReadFile("./assets/shaders/quad.vert")
	if err != nil {
		panic(err)
	}
	quadVertShader, err := shader.NewVertexShader(string(qvsSource))
	if err != nil {
		panic(err)
	}

	qfsSource, err := ioutil.ReadFile("./assets/shaders/quad.frag")
	if err != nil {
		panic(err)
	}
	quadFragShader, err := shader.NewFragmentShader(string(qfsSource))
	if err != nil {
		panic(err)
	}

	quadShaderProg, err = shader.NewProgram(quadVertShader, quadFragShader)
	if err != nil {
		panic(err)
	}

}

type Canvas struct {
	Width, Height int
}

func NewCanvas(width, height int) *Canvas {
	if !isInit {
		_init()
	}
	return &Canvas{
		Width:  width,
		Height: height,
	}
}

func (canvas *Canvas) Bind() {
	quadVao.Bind()
	quadShaderProg.Bind()
}

func (canvas *Canvas) Unbind() {
	quadVao.Unbind()
	quadShaderProg.Unbind()
}

func (canvas *Canvas) BindFor(context internal.BindingClosure) {
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
