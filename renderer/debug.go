package renderer

import (
	"github.com/Qendolin/go-printpixel/core/shader"
	"github.com/Qendolin/go-printpixel/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
)

const (
	Debug = "DebugRenderer"
)

type DebugRenderer struct {
	program    shader.Program
	uTransform shader.Uniform
	Base
}

func NewDebugRenderer() *DebugRenderer {
	vs, err := shader.NewShaderFromPath("@lib/assets/shaders/debug.vert", shader.TypeVertex)
	if err != nil {
		panic(err)
	}
	fs, err := shader.NewShaderFromPath("@lib/assets/shaders/debug.frag", shader.TypeFragment)
	if err != nil {
		panic(err)
	}
	prog, err := shader.NewProgram(vs, fs)
	if err != nil {
		panic(err)
	}
	fs.Destroy()
	vs.Destroy()

	return &DebugRenderer{
		program:    *prog,
		uTransform: prog.MustGetUniform("u_transform"),
	}
}

func (renderer *DebugRenderer) Bind() {
	renderer.program.Bind()
}

func (renderer *DebugRenderer) Unbind() {
	renderer.program.Unbind()
}

func (renderer *DebugRenderer) BindFor(context utils.BindingClosure) {
	renderer.Bind()
	context()
	renderer.Unbind()
}

func (r *DebugRenderer) Draw(ds ...ZDrawable) {
	for i, d := range ds {
		if i == 0 {
			d.GetMesh().Bind()
		}
		d.GetTextures()[0].Bind(0)
		r.uTransform.Set(CalcModelProjectionMat(d, r.Scale))
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
		if i == len(ds)-1 {
			d.GetMesh().Unbind()
		}
	}
}

func (renderer *DebugRenderer) Destroy() {
	renderer.program.Destroy()
}
