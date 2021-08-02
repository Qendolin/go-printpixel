package data

import (
	"github.com/Qendolin/go-printpixel/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Vao struct {
	*uint32
}

func NewVao(id *uint32) *Vao {
	return &Vao{uint32: id}
}

func (vao *Vao) Id() *uint32 {
	if vao.uint32 == nil {
		vao.uint32 = new(uint32)
		gl.GenVertexArrays(1, vao.uint32)
	}
	return vao.uint32
}

func (vao *Vao) Bind() {
	gl.BindVertexArray(*vao.Id())
}

func (vao *Vao) Unbind() {
	gl.BindVertexArray(0)
}

func (vao *Vao) BindFor(context utils.BindingClosure) {
	vao.Bind()
	context()
	vao.Unbind()
}

//TODO: Destroy VBOs, note that VAOs can share VBOs
func (vao *Vao) Destroy() {
	gl.DeleteVertexArrays(1, vao.uint32)
	*vao.uint32 = 0
}

// Specifies the layout of attribute {index} of the currently bound GL_ARRAY_BUFFER
func (vao *Vao) Layout(index int, size int, dataType interface{}, normalized bool, stride int, offset int) (err error) {
	var glType uint32
	var isFloat bool
	glType, isFloat, err = getGlType(dataType)
	if err != nil {
		return
	}
	gl.EnableVertexAttribArray(uint32(index))
	if isFloat {
		// gl.VertexAttribPointerWithOffset(uint32(index), int32(size), glType, normalized, int32(stride), uintptr(offset))
		gl.VertexAttribPointer(uint32(index), int32(size), glType, normalized, int32(stride), gl.PtrOffset(offset))
	} else {
		gl.VertexAttribIPointer(uint32(index), int32(size), glType, int32(stride), gl.PtrOffset(offset))
	}
	return
}

// Like Layout but panics if there is an error
func (vao *Vao) MustLayout(index int, size int, dataType interface{}, normalized bool, stride int, offset int) {
	if err := vao.Layout(index, size, dataType, normalized, stride, offset); err != nil {
		panic(err)
	}
}
