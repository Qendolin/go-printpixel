package renderer

import (
	"github.com/Qendolin/go-printpixel/internal/data"
	"github.com/Qendolin/go-printpixel/internal/shader"
	"github.com/Qendolin/go-printpixel/internal/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type TextureCopy struct {
	Source      data.Texture2D
	Destination data.Texture2D
	program     shader.Program
	quad        data.Vao
	fbo         *uint32
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

	fbo := new(uint32)
	gl.GenFramebuffers(1, fbo)

	return &TextureCopy{
		program: *quadShaderProg,
		quad:    *quadVao,
		fbo:     fbo,
	}
}

func (copy TextureCopy) Draw() {
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
}

func (copy TextureCopy) Bind() {
	copy.quad.Bind()
	copy.program.Bind()
	copy.Source.Bind(0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, *copy.fbo)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, uint32(copy.Destination.Target), copy.Destination.Id(), 0)
}

func (copy TextureCopy) Unbind() {
	copy.quad.Unbind()
	copy.program.Unbind()
	copy.Source.Unbind(0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (copy TextureCopy) BindFor(context utils.BindingClosure) {
	copy.Bind()
	context()
	copy.Unbind()
}

func (copy TextureCopy) Destroy() {
	copy.quad.Destroy()
	copy.program.Destroy()
	gl.DeleteFramebuffers(1, copy.fbo)
}
