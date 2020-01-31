package context

import (
	"errors"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	StatusUninitialized   = iota
	StatusGlfwInitialized = 1 << iota
	StatusGlInitialized
)

var status int

var (
	ErrGlfwNotInitialized = errors.New("GLFW has not been initialized. You have to call InitGlfw() first.")
	ErrGlfwNoContext      = errors.New("GLFW has no context. You have to call MakeContextCurrent() on a *glfw.Window first.")
	ErrGlNotInitialized   = errors.New("OpenGL has not been initialized. You have to call InitGl() first.")
)

func InitGl(cfg glConfig) (err error) {
	if status&StatusGlfwInitialized == 0 {
		err = ErrGlfwNotInitialized
		return
	}

	if glfw.GetCurrentContext() == nil {
		err = ErrGlfwNoContext
		return
	}

	if err = gl.Init(); err != nil {
		return
	}

	if err = cfg.Apply(); err != nil {
		return
	}

	return
}

func InitGlfw() (err error) {
	if status&StatusGlfwInitialized == 0 {
		err = glfw.Init()
		if err == nil {
			status |= StatusGlfwInitialized
		}
	}
	return err
}

func Terminate() {
	if status&StatusGlfwInitialized > 0 {
		glfw.Terminate()
		status &= ^StatusGlfwInitialized
	}
}

func Status() int {
	return status
}
