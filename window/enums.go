package window

import (
	iWin "github.com/Qendolin/go-printpixel/internal/window"
)

//ClientApi
const (
	OpenGLAPI   = iWin.OpenGLAPI
	OpenGLESAPI = iWin.OpenGLESAPI
)

//ContextRobustness
const (
	NoRobustness        = iWin.NoRobustness
	NoResetNotification = iWin.NoResetNotification
	LoseContextOnReset  = iWin.LoseContextOnReset
)

//ContextReleaseBehavior
const (
	AnyReleaseBehavior   = iWin.AnyReleaseBehavior
	ReleaseBehaviorFlush = iWin.ReleaseBehaviorFlush
	ReleaseBehaviorNone  = iWin.ReleaseBehaviorNone
)

//OpenGLProfile
const (
	OpenGLAnyProfile    = iWin.OpenGLAnyProfile
	OpenGLCoreProfile   = iWin.OpenGLCoreProfile
	OpenGLCompatProfile = iWin.OpenGLCompatProfile
)

const (
	DontCare = iWin.DontCare
)
