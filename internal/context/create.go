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

var ErrGlfwNotInitialized error = errors.New("GLFW has not been initialized. You have to call InitGlfw() first.")
var ErrGlNotInitialized error = errors.New("OpenGL has not been initialized. You have to call InitGl() first.")

func InitGl(cfg glConfig) (err error) {
	if status&StatusGlfwInitialized == 0 {
		err = ErrGlfwNotInitialized
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
