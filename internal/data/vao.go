package data

import (
	"github.com/Qendolin/go-printpixel/internal/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Vao struct {
	*uint32
}

func NewVao() *Vao {
	id := new(uint32)
	gl.GenVertexArrays(1, id)
	return &Vao{id}
}

func (vao *Vao) Id() uint32 {
	return *vao.uint32
}

func (vao *Vao) Bind() {
	gl.BindVertexArray(vao.Id())
}

func (vao *Vao) Unbind() {
	gl.BindVertexArray(0)
}

func (vao *Vao) BindFor(context utils.BindingClosure) {
	vao.Bind()
	defered := context()
	vao.Unbind()
	for _, deferedFunc := range defered {
		deferedFunc()
	}
}

//TODO: Destroy VBOs, note that VAOs can share VBOs
func (vao *Vao) Destroy() {
	gl.DeleteVertexArrays(1, vao.uint32)
	vao.uint32 = nil
}
