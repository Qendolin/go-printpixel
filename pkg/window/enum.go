package window

import (
	"github.com/Qendolin/go-printpixel/core/glwindow"
)

//ClientApi
const (
	OpenGLAPI   = glwindow.OpenGLAPI
	OpenGLESAPI = glwindow.OpenGLESAPI
)

//ContextRobustness
const (
	NoRobustness        = glwindow.NoRobustness
	NoResetNotification = glwindow.NoResetNotification
	LoseContextOnReset  = glwindow.LoseContextOnReset
)

//ContextReleaseBehavior
const (
	AnyReleaseBehavior   = glwindow.AnyReleaseBehavior
	ReleaseBehaviorFlush = glwindow.ReleaseBehaviorFlush
	ReleaseBehaviorNone  = glwindow.ReleaseBehaviorNone
)

//OpenGLProfile
const (
	OpenGLAnyProfile    = glwindow.OpenGLAnyProfile
	OpenGLCoreProfile   = glwindow.OpenGLCoreProfile
	OpenGLCompatProfile = glwindow.OpenGLCompatProfile
)

const (
	DontCare = glwindow.DontCare
)
