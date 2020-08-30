package data

import (
	"encoding/binary"

	"github.com/Qendolin/go-printpixel/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Vbo struct {
	*uint32
}

func NewVbo(id *uint32) *Vbo {
	return &Vbo{uint32: id}
}

func (vbo *Vbo) Id() *uint32 {
	if vbo.uint32 == nil {
		vbo.uint32 = new(uint32)
		gl.GenBuffers(1, vbo.uint32)
	}
	return vbo.uint32
}

func (vbo *Vbo) Bind(target uint32) {
	gl.BindBuffer(target, *vbo.Id())
}

func (vbo *Vbo) Unbind(target uint32) {
	gl.BindBuffer(target, 0)
}

func (vbo *Vbo) BindFor(target uint32, context utils.BindingClosure) {
	vbo.Bind(target)
	context()
	vbo.Unbind(target)
}

func (vbo *Vbo) Destroy() {
	gl.DeleteBuffers(1, vbo.uint32)
	*vbo.uint32 = 0
}

func (vbo *Vbo) WriteStatic(data interface{}) {
	vbo.Write(gl.STATIC_DRAW, data)
}

/*
	data - a silce of some type
*/
func (vbo *Vbo) Write(mode uint32, data interface{}) {
	size := binary.Size(data)
	if size == -1 {
		//Ignore, gl will throw error anyway
	}
	gl.BufferData(gl.ARRAY_BUFFER, size, gl.Ptr(data), mode)
}

func (vbo *Vbo) Layout(index int, size int, dataType interface{}, normalized bool, stride int) (err error) {
	var glType uint32
	var isFloat bool
	glType, isFloat, err = getGlType(dataType)
	if err != nil {
		return
	}
	gl.EnableVertexAttribArray(uint32(index))
	if isFloat {
		gl.VertexAttribPointer(uint32(index), int32(size), glType, normalized, int32(stride), nil)
	} else {
		gl.VertexAttribIPointer(uint32(index), int32(size), glType, int32(stride), nil)
	}
	return
}

/*
	Like Layout but panics if there is an error
*/
func (vbo *Vbo) MustLayout(index int, size int, dataType interface{}, normalized bool, stride int) {
	if err := vbo.Layout(index, size, dataType, normalized, stride); err != nil {
		panic(err)
	}
}