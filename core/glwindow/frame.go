// +build !windows

package glwindow

func (win extWindow) GetFrameSize() (left, top, right, bottom int) {
	return win.Window.GetFrameSize()
}
