package glwindow

import "github.com/go-gl/glfw/v3.3/glfw"

type GlfwEnum int

//ClientApi
const (
	OpenGLAPI   = GlfwEnum(glfw.OpenGLAPI)
	OpenGLESAPI = GlfwEnum(glfw.OpenGLESAPI)
)

//ContextRobustness
const (
	NoRobustness        = GlfwEnum(glfw.NoRobustness)
	NoResetNotification = GlfwEnum(glfw.NoResetNotification)
	LoseContextOnReset  = GlfwEnum(glfw.LoseContextOnReset)
)

//ContextReleaseBehavior
const (
	AnyReleaseBehavior   = GlfwEnum(glfw.AnyReleaseBehavior)
	ReleaseBehaviorFlush = GlfwEnum(glfw.ReleaseBehaviorFlush)
	ReleaseBehaviorNone  = GlfwEnum(glfw.ReleaseBehaviorNone)
)

//OpenGLProfile
const (
	OpenGLAnyProfile    = GlfwEnum(glfw.OpenGLAnyProfile)
	OpenGLCoreProfile   = GlfwEnum(glfw.OpenGLCoreProfile)
	OpenGLCompatProfile = GlfwEnum(glfw.OpenGLCompatProfile)
)

const (
	DontCare = GlfwEnum(glfw.DontCare)
)
