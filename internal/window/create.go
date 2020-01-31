package window

import "github.com/go-gl/glfw/v3.3/glfw"

import "github.com/Qendolin/go-printpixel/internal/context"

func New(hints hints, title string, width, height int, monitor *glfw.Monitor) (win *glfw.Window, err error) {

	if context.Status()&context.StatusGlfwInitialized == 0 {
		err = context.ErrGlfwNotInitialized
		return
	}

	glfw.DefaultWindowHints()
	hints.apply()

	if monitor == nil && (hints.Maximized.Value || hints.ScaleToMonitor.Value) {
		monitor = glfw.GetPrimaryMonitor()
	}
	if hints.ScaleToMonitor.Value {
		vidMode := monitor.GetVideoMode()
		width = vidMode.Width
		height = vidMode.Height
	}

	win, err = glfw.CreateWindow(width, height, title, monitor, nil)
	return
}
