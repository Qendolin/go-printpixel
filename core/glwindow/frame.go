// +build !windows

package glwindow

func (win extWindow) GetVisibleFrameSize() (left, top, right, bottom int) {
	return win.GetFrameSize()
}
