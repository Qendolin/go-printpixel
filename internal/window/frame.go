// +build !windows

package window

func (win glfwWindow) GetVisibleFrameSize() (left, top, right, bottom int) {
	return win.GetFrameSize()
}
