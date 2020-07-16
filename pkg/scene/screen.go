package scene

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

func (s Screen) Layout() []Layoutable {
	if s.Child == nil {
		return nil
	}

	s.Child.SetX(0)
	s.Child.SetY(0)
	s.Child.SetWidth(s.width)
	s.Child.SetHeight(s.height)

	return []Layoutable{s.Child}
}
