package canvas

import "github.com/go-gl/gl/v3.2-core/gl"

import "encoding/binary"

type Vbo struct {
	*uint32
}

func NewVbo() *Vbo {
	id := new(uint32)
	gl.GenBuffers(1, id)
	return &Vbo{id}
}

func (vbo *Vbo) Id() uint32 {
	return *vbo.uint32
}

func (vbo *Vbo) Bind(target uint32) {
	gl.BindBuffer(target, vbo.Id())
}

/*
	data - a silce of some type
*/
func (vbo *Vbo) WriteStatic(data interface{}, length int, dataType interface{}) {
	vbo.Bind(gl.ARRAY_BUFFER)
	size := binary.Size(data)
	if size == -1 {
		//TODO
		panic("invalid data type")
	}
	gl.BufferData(gl.ARRAY_BUFFER, size, gl.Ptr(data), gl.STATIC_DRAW)
}

func (vbo *Vbo) Destroy() {
	gl.DeleteBuffers(1, vbo.uint32)
	vbo.uint32 = nil
}
