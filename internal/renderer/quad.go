package renderer

import (
	"github.com/Qendolin/go-printpixel/internal/data"
	"github.com/Qendolin/go-printpixel/internal/shader"
	"github.com/Qendolin/go-printpixel/internal/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
)

var (
	quadVertices = []float32{1, -1, 1, 1, -1, -1, -1, 1}
)

type TextureQuad struct {
	Program shader.Program
	Texture *data.Texture2D
	quad    data.Vao
}

func NewTextureQuad() *TextureQuad {
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
	return NewTextureQuadWithProgram(*quadShaderProg)
}

func NewTextureQuadWithProgram(prog shader.Program) *TextureQuad {
	quadVao := data.NewVao()
	quadVao.BindFor(func() {
		quadVbo := data.NewVbo()
		quadVbo.Bind(gl.ARRAY_BUFFER)
		quadVbo.WriteStatic(quadVertices)
		quadVbo.MustLayout(0, 2, float32(0), false, 0)

		quadVbo.Unbind(gl.ARRAY_BUFFER)
	})
	return &TextureQuad{Program: prog, quad: *quadVao, Texture: data.NewTexture2D(data.Tex2DTarget2D)}
}

func (renderer *TextureQuad) Bind() {
	renderer.quad.Bind()
	renderer.Program.Bind()
	if renderer.Texture != nil {
		renderer.Texture.Bind(0)
	}
}

func (renderer *TextureQuad) Unbind() {
	renderer.quad.Unbind()
	renderer.Program.Unbind()
	if renderer.Texture != nil {
		renderer.Texture.Unbind(0)
	}
}

func (renderer *TextureQuad) BindFor(context utils.BindingClosure) {
	renderer.Bind()
	context()
	renderer.Unbind()
}

func (renderer *TextureQuad) Draw() {
	//gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
}

func (renderer *TextureQuad) Destroy() {
	renderer.Program.Destroy()
	renderer.quad.Destroy()
	if renderer.Texture != nil {
		renderer.Texture.Destroy()
	}
}
