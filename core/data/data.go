package data

import (
	"sync"

	"github.com/go-gl/glfw/v3.3/glfw"
)

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
