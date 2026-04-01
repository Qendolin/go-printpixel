// +build !windows

package glw

func (win extWindow) GetFrameSize() (left, top, right, bottom int) {
	return win.Window.GetFrameSize()
}
