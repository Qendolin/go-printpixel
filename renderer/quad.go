package renderer

import (
	"github.com/Qendolin/go-printpixel/core/shader"
	"github.com/Qendolin/go-printpixel/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
)

const (
	TextureQuad = "TextureQuadRenderer"
)

type TextureQuadRenderer struct {
	program    shader.Program
	uTransform shader.Uniform
	Base
}

func NewTextureQuadRenderer() *TextureQuadRenderer {
	vs, err := shader.NewShaderFromPath("@mod/assets/shaders/quad_tex_transform.vert", shader.TypeVertex)
	if err != nil {
		panic(err)
	}
	fs, err := shader.NewShaderFromPath("@mod/assets/shaders/quad_tex_transform.frag", shader.TypeFragment)
	if err != nil {
		panic(err)
	}
	prog, err := shader.NewProgram(vs, fs)
	if err != nil {
		panic(err)
	}
	fs.Destroy()
	vs.Destroy()

	return &TextureQuadRenderer{
		program:    *prog,
		uTransform: prog.MustGetUniform("u_transform"),
	}
}

func (renderer *TextureQuadRenderer) Bind() {
	renderer.program.Bind()
}

func (renderer *TextureQuadRenderer) Unbind() {
	renderer.program.Unbind()
}

func (renderer *TextureQuadRenderer) BindFor(context utils.BindingClosure) {
	renderer.Bind()
	context()
	renderer.Unbind()
}

func (r *TextureQuadRenderer) Draw(ds ...ZDrawable) {
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

func (renderer *TextureQuadRenderer) Destroy() {
	renderer.program.Destroy()
}
