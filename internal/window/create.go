package window

import "github.com/go-gl/glfw/v3.3/glfw"

import "errors"

var glfwInit bool

var ErrGLFWNotInitialized error = errors.New("GLFW has not been initialized. You have to call Init() first.")

func Init() (err error) {
	if !glfwInit {
		err = glfw.Init()
		if err == nil {
			glfwInit = true
		}
	}
	return err
}

func Terminate() {
	if glfwInit {
		glfw.Terminate()
	}
}

func NewWindow(hints hints, title string, width, height int, monitor *glfw.Monitor) (win *glfw.Window, err error) {

	if !glfwInit {
		err = ErrGLFWNotInitialized
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
