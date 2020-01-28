package window

import "github.com/go-gl/glfw/v3.2/glfw"

type glfwEnum int

//ClientApi
const (
	OpenGLAPI   = glfwEnum(glfw.OpenGLAPI)
	OpenGLESAPI = glfwEnum(glfw.OpenGLESAPI)
)

//ContextRobustness
const (
	NoRobustness        = glfwEnum(glfw.NoRobustness)
	NoResetNotification = glfwEnum(glfw.NoResetNotification)
	LoseContextOnReset  = glfwEnum(glfw.LoseContextOnReset)
)

//ContextReleaseBehavior
const (
	AnyReleaseBehavior   = glfwEnum(glfw.AnyReleaseBehavior)
	ReleaseBehaviorFlush = glfwEnum(glfw.ReleaseBehaviorFlush)
	ReleaseBehaviorNone  = glfwEnum(glfw.ReleaseBehaviorNone)
)

//OpenGLProfile
const (
	OpenGLAnyProfile    = glfwEnum(glfw.OpenGLAnyProfile)
	OpenGLCoreProfile   = glfwEnum(glfw.OpenGLCoreProfile)
	OpenGLCompatProfile = glfwEnum(glfw.OpenGLCompatProfile)
)

const (
	DontCare = glfwEnum(glfw.DontCare)
)

type glfwHint interface {
	intValue() int
	Code() glfw.Hint
}

type hint struct {
	target glfw.Hint
}

func (h hint) Code() glfw.Hint {
	return h.target
}

func applyHint(hint glfwHint) {
	glfw.WindowHint(hint.Code(), hint.intValue())
}

type boolHint struct {
	hint
	Value bool
}

func (bh boolHint) intValue() int {
	if bh.Value {
		return 1
	}
	return 0
}

type intHint struct {
	hint
	Value int
}

func (ih intHint) intValue() int {
	return ih.Value
}

type enumHint struct {
	hint
	Value glfwEnum
}

func (eh enumHint) intValue() int {
	return int(eh.Value)
}

type hints struct {
	windowHints
	contextHints
	framebufferHints
}

func NewHints() hints {
	return hints{
		newWindowHints(),
		newContextHints(),
		newFramebufferHints(),
	}
}

func (h hints) apply() {
	applyHint(h.ClientAPI)
	applyHint(h.ContextVersionMajor)
	applyHint(h.ContextVersionMinor)
	applyHint(h.ContextRobustness)
	applyHint(h.ContextReleaseBehavior)
	applyHint(h.OpenGLForwardCompatible)
	applyHint(h.OpenGLDebugContext)
	applyHint(h.OpenGLProfile)
	applyHint(h.SRGBCapable)
	applyHint(h.DepthBits)
	applyHint(h.Samples)
	applyHint(h.RefreshRate)
	applyHint(h.DoubleBuffer)
	applyHint(h.Focused)
	applyHint(h.Visible)
	applyHint(h.Resizable)
	applyHint(h.Decorated)
	applyHint(h.Floating)
	applyHint(h.AutoIconify)
	applyHint(h.Maximized)
}

//window related hints
//See https://www.glfw.org/docs/latest/window_guide.html
type windowHints struct {
	Focused     boolHint
	Visible     boolHint
	Resizable   boolHint
	Decorated   boolHint
	Floating    boolHint
	AutoIconify boolHint
	Maximized   boolHint
}

func newWindowHints() windowHints {
	return windowHints{
		Resizable:   boolHint{Value: true, hint: hint{target: glfw.Resizable}},
		Visible:     boolHint{Value: true, hint: hint{target: glfw.Visible}},
		Decorated:   boolHint{Value: true, hint: hint{target: glfw.Decorated}},
		Focused:     boolHint{Value: true, hint: hint{target: glfw.Focused}},
		AutoIconify: boolHint{Value: true, hint: hint{target: glfw.AutoIconify}},
		Floating:    boolHint{Value: false, hint: hint{target: glfw.Floating}},
		Maximized:   boolHint{Value: false, hint: hint{target: glfw.Maximized}},
	}
}

//context related hints
type contextHints struct {
	//See https://www.glfw.org/docs/latest/window_guide.html#GLFW_CLIENT_API_hint
	ClientAPI           enumHint
	ContextVersionMajor intHint
	ContextVersionMinor intHint
	//See https://www.glfw.org/docs/latest/window_guide.html#GLFW_CONTEXT_ROBUSTNESS_hint
	ContextRobustness enumHint
	//See https://www.glfw.org/docs/latest/window_guide.html#GLFW_CONTEXT_RELEASE_BEHAVIOR_hint
	ContextReleaseBehavior  enumHint
	OpenGLForwardCompatible boolHint
	OpenGLDebugContext      boolHint
	//See https://www.glfw.org/docs/latest/window_guide.html#GLFW_OPENGL_PROFILE_hint
	OpenGLProfile enumHint
	SRGBCapable   boolHint
}

func newContextHints() contextHints {
	return contextHints{
		ClientAPI:               enumHint{Value: OpenGLAPI, hint: hint{target: glfw.ClientAPI}},
		ContextVersionMajor:     intHint{Value: 1, hint: hint{target: glfw.ContextVersionMajor}},
		ContextVersionMinor:     intHint{Value: 0, hint: hint{target: glfw.ContextVersionMinor}},
		ContextRobustness:       enumHint{Value: NoRobustness, hint: hint{target: glfw.ContextRobustness}},
		ContextReleaseBehavior:  enumHint{Value: AnyReleaseBehavior, hint: hint{target: glfw.ContextReleaseBehavior}},
		OpenGLForwardCompatible: boolHint{Value: false, hint: hint{target: glfw.OpenGLForwardCompatible}},
		OpenGLDebugContext:      boolHint{Value: false, hint: hint{target: glfw.OpenGLDebugContext}},
		OpenGLProfile:           enumHint{Value: OpenGLAnyProfile, hint: hint{target: glfw.OpenGLProfile}},
		SRGBCapable:             boolHint{Value: false, hint: hint{target: glfw.SRGBCapable}},
	}
}

//framebuffer related hints
type framebufferHints struct {
	DepthBits    intHint
	Samples      intHint
	RefreshRate  intHint
	DoubleBuffer boolHint
}

func newFramebufferHints() framebufferHints {
	return framebufferHints{
		DepthBits:    intHint{Value: 24, hint: hint{target: glfw.DepthBits}},
		Samples:      intHint{Value: 0, hint: hint{target: glfw.Samples}},
		RefreshRate:  intHint{Value: int(DontCare), hint: hint{target: glfw.RefreshRate}},
		DoubleBuffer: boolHint{Value: true, hint: hint{target: glfw.DoubleBuffer}},
	}
}
