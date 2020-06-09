package layout

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Screen struct {
	width  int
	height int
	Child  Layoutable
}

func NewScreen(mon glfw.Monitor) Screen {
	vidMode := mon.GetVideoMode()
	return Screen{
		width:  vidMode.Width,
		height: vidMode.Height,
	}
}

func NewScreenByDimensions(width, height int) Screen {
	return Screen{
		width:  width,
		height: height,
	}
}

func (b Screen) Layout() {
	if b.Child == nil {
		return
	}

	b.Child.SetX(0)
	b.Child.SetY(0)
	b.Child.SetWidth(b.width)
	b.Child.SetHeight(b.height)

	if l, ok := b.Child.(Layouter); ok {
		l.Layout()
	}
}
