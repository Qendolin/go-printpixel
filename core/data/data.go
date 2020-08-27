package data

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type NotEnoughError struct {
	RequiredSize int64
	ActualSize   int64
}

func (nee NotEnoughError) Error() string {
	return fmt.Sprintf("data has %d less bytes than %d, which are required", (nee.RequiredSize - nee.ActualSize), nee.RequiredSize)
}

type ContextInstance struct {
	ids map[uint64]*uint32
	m   sync.Mutex
}

func (ci *ContextInstance) Create(f func() *uint32) *uint32 {
	ci.m.Lock()
	defer ci.m.Unlock()
	if ci.ids == nil {
		ci.ids = make(map[uint64]*uint32, 1)
	}

	wId := (*uint64)(glfw.GetCurrentContext().GetUserPointer())
	if id, ok := ci.ids[*wId]; ok {
		return id
	}

	id := f()
	ci.ids[*wId] = id
	return id
}

type TypeError struct {
	Type reflect.Type
}

func (terr TypeError) Error() string {
	return fmt.Sprintf("Invalid dataType %v, doesn't correspond to any gl type", terr.Type)
}

func getGlType(data interface{}) (glType uint32, float bool, err error) {
	switch data.(type) {
	case byte, []byte:
		return gl.UNSIGNED_BYTE, false, nil
	case int8, []int8:
		return gl.BYTE, false, nil
	case int16, []int16:
		return gl.SHORT, false, nil
	case uint16, []uint16:
		return gl.UNSIGNED_SHORT, false, nil
	case int32, []int32:
		return gl.INT, false, nil
	case uint32, []uint32:
		return gl.UNSIGNED_INT, false, nil
	case float32, []float32:
		return gl.FLOAT, true, nil
	case float64, []float64:
		return gl.DOUBLE, true, nil
	}
	return 0, false, TypeError{
		Type: reflect.TypeOf(data),
	}
}

func getGlTypeSize(typ uint32) int {
	switch typ {
	case gl.BYTE, gl.UNSIGNED_BYTE:
		return 1
	case gl.SHORT, gl.UNSIGNED_SHORT, gl.HALF_FLOAT:
		return 2
	case gl.INT, gl.UNSIGNED_INT, gl.FIXED, gl.FLOAT:
		return 4
	case gl.DOUBLE:
		return 8
	}
	return 0
}
