package core

import (
	"github.com/Qendolin/go-printpixel/core/data"
	"github.com/go-gl/gl/v3.3-core/gl"
)

var (
	quadVertices = []float32{0.5, -0.5, 0.5, 0.5, -0.5, -0.5, -0.5, 0.5}
	quadInstance = &data.ContextInstance{}
)

func Quad() *data.Vao {
	id := quadInstance.Create(func() *uint32 {
		quadVao := data.NewVao(nil)
		quadVao.Bind()
		quadVbo := data.Vbo{}
		quadVbo.Bind(gl.ARRAY_BUFFER)
		quadVbo.WriteStatic(quadVertices)
		quadVbo.MustLayout(0, 2, float32(0), false, 0)

		quadVbo.Unbind(gl.ARRAY_BUFFER)
		quadVao.Unbind()

		return quadVao.Id()
	})
	return data.NewVao(id)
}
