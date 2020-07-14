package renderer

import (
	"github.com/Qendolin/go-printpixel/core/data"
	"github.com/Qendolin/go-printpixel/core/shader"
	"github.com/Qendolin/go-printpixel/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type TextureCopy struct {
	source      data.Texture2D
	destination data.Texture2D
	program     shader.Program
	quad        data.Vao
	fbo         data.Fbo
}

func NewTextureCopy() *TextureCopy {

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

	quadVao := data.NewVao()
	quadVao.BindFor(func() {
		quadVbo := data.NewVbo()
		quadVbo.Bind(gl.ARRAY_BUFFER)
		quadVbo.WriteStatic(quadVertices)
		quadVbo.MustLayout(0, 2, float32(0), false, 0)

		quadVbo.Unbind(gl.ARRAY_BUFFER)
	})

	return &TextureCopy{
		program: *quadShaderProg,
		quad:    *quadVao,
		fbo:     *data.NewFbo(),
	}
}

func (copy TextureCopy) Draw() {
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
}

func (copy TextureCopy) Bind() {
	copy.quad.Bind()
	copy.program.Bind()
	copy.source.Bind(0)
	copy.fbo.Bind(data.FboTargetReadWrite)
}

func (copy TextureCopy) Unbind() {
	copy.quad.Unbind()
	copy.program.Unbind()
	copy.source.Unbind(0)
	copy.fbo.Bind(data.FboTargetReadWrite)
}

func (copy TextureCopy) BindFor(context utils.BindingClosure) {
	copy.Bind()
	context()
	copy.Unbind()
}

func (copy TextureCopy) Destroy() {
	copy.quad.Destroy()
	copy.program.Destroy()
	copy.fbo.Destroy()
}

func (copy *TextureCopy) SetSource(tex data.Texture2D) {
	copy.source = tex
}

func (copy *TextureCopy) SetDestination(tex data.Texture2D) {
	copy.destination = tex
	copy.fbo.AttachTexture(tex.GLTexture, 0, 0)
}
