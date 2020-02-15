package data

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/Qendolin/go-printpixel/internal/utils"
	"github.com/go-gl/gl/v3.2-core/gl"
)

type TypeErr struct {
	Type reflect.Type
}

func (terr TypeErr) Error() string {
	return fmt.Sprintf("Invalid dataType %v, doesn't correspond to any gl type", terr.Type)
}

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

func (vbo *Vbo) Unbind(target uint32) {
	gl.BindBuffer(target, 0)
}

func (vbo *Vbo) BindFor(target uint32, context utils.BindingClosure) {
	vbo.Bind(target)
	defered := context()
	vbo.Unbind(target)
	for _, deferedFunc := range defered {
		deferedFunc()
	}
}

/*
	data - a silce of some type
*/
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
	glType, isFloat, err = getGlType(reflect.TypeOf(dataType))
	if err != nil {
		return
	}
	if isFloat {
		gl.VertexAttribPointer(uint32(index), int32(size), glType, normalized, int32(stride), nil)
	} else {
		gl.VertexAttribIPointer(uint32(index), int32(size), glType, int32(stride), nil)
	}
	gl.EnableVertexAttribArray(uint32(index))
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

func getGlType(dataType reflect.Type) (glType uint32, float bool, err error) {
	name := dataType.Name()
	switch name {
	case "byte":
		fallthrough
	case "uint8":
		glType = gl.UNSIGNED_BYTE
	case "int8":
		glType = gl.BYTE
	case "int16":
		glType = gl.SHORT
	case "uint16":
		glType = gl.UNSIGNED_SHORT
	case "int":
		fallthrough
	case "rune":
		fallthrough
	case "int32":
		glType = gl.INT
	case "uint":
		fallthrough
	case "uint32":
		glType = gl.UNSIGNED_INT
	case "float32":
		glType = gl.FLOAT
		float = true
	case "float16":
		glType = gl.HALF_FLOAT
		float = true
	case "float64":
		glType = gl.DOUBLE
		float = true
	}
	if glType == 0 {
		err = TypeErr{
			Type: dataType,
		}
	}
	return
}

func (vbo *Vbo) Destroy() {
	gl.DeleteBuffers(1, vbo.uint32)
	vbo.uint32 = nil
}
