package window

import (
	"github.com/Qendolin/go-printpixel/internal/context"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func New(hints Hints, title string, width, height int, monitor *glfw.Monitor) (win *glfw.Window, err error) {

	if context.Status()&context.StatusGlfwInitialized == 0 {
		err = context.ErrGlfwNotInitialized
		return
	}

	glfw.DefaultWindowHints()
	hints.apply()

	if monitor == nil && (hints.Fullscreen.Value) {
		monitor = glfw.GetPrimaryMonitor()
	}

	win, err = glfw.CreateWindow(width, height, title, monitor, nil)
	win.MakeContextCurrent()

	if hints.Vsync.Value {
		glfw.SwapInterval(1)
	}

	return
}
