package window

import (
	"github.com/Qendolin/go-printpixel/core/glw"
)

// ClientApi
const (
	OpenGLAPI   = glw.OpenGLAPI
	OpenGLESAPI = glw.OpenGLESAPI
)

// ContextRobustness
const (
	NoRobustness        = glw.NoRobustness
	NoResetNotification = glw.NoResetNotification
	LoseContextOnReset  = glw.LoseContextOnReset
)

// ContextReleaseBehavior
const (
	AnyReleaseBehavior   = glw.AnyReleaseBehavior
	ReleaseBehaviorFlush = glw.ReleaseBehaviorFlush
	ReleaseBehaviorNone  = glw.ReleaseBehaviorNone
)

// OpenGLProfile
const (
	OpenGLAnyProfile    = glw.OpenGLAnyProfile
	OpenGLCoreProfile   = glw.OpenGLCoreProfile
	OpenGLCompatProfile = glw.OpenGLCompatProfile
)

const (
	DontCare = glw.DontCare
)
