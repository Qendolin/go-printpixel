package core

import (
	"github.com/Qendolin/go-printpixel/core/data"
)

var (
	quadVertices = []float32{0.5, -0.5, 0.5, 0.5, -0.5, -0.5, -0.5, 0.5}
	quadInstance = &data.ContextInstance{}
)

func Quad() *data.Vao {
	id := quadInstance.Create(func() *uint32 {
		quadVao := data.NewVao(nil)
		quadVao.Bind()
		quadVbo := data.Buffer{Target: data.BufVertexAttribute}
		quadVbo.Bind()
		quadVbo.WriteStatic(quadVertices)
		quadVao.MustLayout(0, 2, float32(0), false, 0, 0)

		quadVao.Unbind()
		quadVbo.Unbind()

		return quadVao.Id()
	})
	return data.NewVao(id)
}
