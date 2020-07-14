// +build !windows

package glwindow

func (win glfwWindow) GetVisibleFrameSize() (left, top, right, bottom int) {
	return win.GetFrameSize()
}
