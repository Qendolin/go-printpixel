package data

import (
	"encoding/binary"

	"github.com/Qendolin/go-printpixel/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type BufferTarget int

const (
	// Vertex attributes
	BufVertexAttribute = BufferTarget(gl.ARRAY_BUFFER)
	// Vertex array indices
	BufElementIndex = BufferTarget(gl.ELEMENT_ARRAY_BUFFER)
)

type Buffer struct {
	*uint32
	Target BufferTarget
}

func NewBuffer(id *uint32, bufType BufferTarget) *Buffer {
	return &Buffer{uint32: id, Target: bufType}
}

func (buf *Buffer) Id() *uint32 {
	if buf.uint32 == nil {
		buf.uint32 = new(uint32)
		gl.GenBuffers(1, buf.uint32)
	}
	return buf.uint32
}

func (buf *Buffer) Bind() {
	gl.BindBuffer(uint32(buf.Target), *buf.Id())
}

func (buf *Buffer) Unbind() {
	gl.BindBuffer(uint32(buf.Target), 0)
}

func (buf *Buffer) BindFor(context utils.BindingClosure) {
	buf.Bind()
	context()
	buf.Unbind()
}

func (buf *Buffer) Destroy() {
	gl.DeleteBuffers(1, buf.uint32)
	*buf.uint32 = 0
}

func (buf *Buffer) WriteStatic(data interface{}) {
	buf.Write(gl.STATIC_DRAW, data)
}

// data - a silce of some type
func (buf *Buffer) Write(mode uint32, data interface{}) {
	size := binary.Size(data)
	if size == -1 {
		// Ignore, gl will throw error anyway
	}
	gl.BufferData(uint32(buf.Target), size, gl.Ptr(data), mode)
}
