package renderer

import (
	"github.com/Qendolin/go-printpixel/core/data"
	"github.com/Qendolin/go-printpixel/core/shader"
	"github.com/Qendolin/go-printpixel/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	quadVertices = []float32{1, -1, 1, 1, -1, -1, -1, 1}
)

type TextureQuad struct {
	Texture   *data.Texture2D
	Transform mgl32.Mat3
}

func (tq *TextureQuad) Bind(texUnit int) {
	tq.Texture.Bind(texUnit)
}

func (tq *TextureQuad) Unbind(texUnit int) {
	tq.Texture.Unbind(texUnit)
}

func (tq *TextureQuad) BindFor(texUnit int, context utils.BindingClosure) {
	tq.Bind(texUnit)
	context()
	tq.Unbind(texUnit)
}

func NewTextureQuad() *TextureQuad {
	return &TextureQuad{
		Texture:   data.NewTexture2D(data.Tex2DTarget2D),
		Transform: mgl32.Ident3(),
	}
}

func NewQuad() *data.Vao {
	quadVao := data.NewVao()
	quadVao.Bind()
	quadVbo := data.NewVbo()
	quadVbo.Bind(gl.ARRAY_BUFFER)
	quadVbo.WriteStatic(quadVertices)
	quadVbo.MustLayout(0, 2, float32(0), false, 0)

	quadVbo.Unbind(gl.ARRAY_BUFFER)
	quadVao.Unbind()
	return quadVao
}

type TextureQuadRenderer struct {
	program    shader.Program
	quad       data.Vao
	uTransform shader.Uniform
}

func NewTextureQuadRenderer() *TextureQuadRenderer {
	vs, err := shader.NewShaderFromPath("assets/shaders/quad_tex_transform.vert", shader.TypeVertex)
	if err != nil {
		panic(err)
	}
	fs, err := shader.NewShaderFromPath("assets/shaders/quad_tex_transform.frag", shader.TypeFragment)
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
		quad:       *NewQuad(),
		uTransform: prog.MustGetUniform("u_transform"),
	}
}

func (renderer *TextureQuadRenderer) Bind() {
	renderer.program.Bind()
	renderer.quad.Bind()
}

func (renderer *TextureQuadRenderer) Unbind() {
	renderer.program.Unbind()
	renderer.quad.Unbind()
}

func (renderer *TextureQuadRenderer) BindFor(context utils.BindingClosure) {
	renderer.Bind()
	context()
	renderer.Unbind()
}

func (renderer *TextureQuadRenderer) Draw(scaleX, scaleY float32, quads ...TextureQuad) {
	for _, tq := range quads {
		tq.Bind(0)
		renderer.uTransform.Set(tq.Transform.Mul3(mgl32.Diag3(mgl32.Vec3{scaleX, scaleY, 1})))
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
		tq.Unbind(0)
	}
}

func (renderer *TextureQuadRenderer) Destroy() {
	renderer.program.Destroy()
	renderer.quad.Destroy()
}
