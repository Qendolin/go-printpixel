package glwindow

import "github.com/go-gl/glfw/v3.3/glfw"

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
		return glfw.True
	}
	return glfw.False
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
	Value GlfwEnum
}

func (eh enumHint) intValue() int {
	return int(eh.Value)
}

type Hints struct {
	windowHints
	contextHints
	framebufferHints
	customHints
}

func NewHints() Hints {
	return Hints{
		newWindowHints(),
		newContextHints(),
		newFramebufferHints(),
		newCustomHints(),
	}
}

func (h Hints) apply() {
	h.contextHints.apply()
	h.windowHints.apply()
	h.framebufferHints.apply()
}

// Window related hints
// See https://www.glfw.org/docs/latest/window_guide.html
type windowHints struct {
	Focused        boolHint
	Visible        boolHint
	Resizable      boolHint
	Decorated      boolHint
	Floating       boolHint
	AutoIconify    boolHint
	Maximized      boolHint
	ScaleToMonitor boolHint
}

func newWindowHints() windowHints {
	return windowHints{
		Resizable:      boolHint{Value: true, hint: hint{target: glfw.Resizable}},
		Visible:        boolHint{Value: true, hint: hint{target: glfw.Visible}},
		Decorated:      boolHint{Value: true, hint: hint{target: glfw.Decorated}},
		Focused:        boolHint{Value: true, hint: hint{target: glfw.Focused}},
		AutoIconify:    boolHint{Value: true, hint: hint{target: glfw.AutoIconify}},
		Floating:       boolHint{Value: false, hint: hint{target: glfw.Floating}},
		Maximized:      boolHint{Value: false, hint: hint{target: glfw.Maximized}},
		ScaleToMonitor: boolHint{Value: false, hint: hint{target: glfw.ScaleToMonitor}},
	}
}

func (h windowHints) apply() {
	applyHint(h.Focused)
	applyHint(h.Visible)
	applyHint(h.Resizable)
	applyHint(h.Decorated)
	applyHint(h.Floating)
	applyHint(h.AutoIconify)
	applyHint(h.Maximized)
	applyHint(h.ScaleToMonitor)
}

// Context related hints
type contextHints struct {
	// See https://www.glfw.org/docs/latest/window_guide.html#GLFW_CLIENT_API_hint
	ClientAPI           enumHint
	ContextVersionMajor intHint
	ContextVersionMinor intHint
	// See https://www.glfw.org/docs/latest/window_guide.html#GLFW_CONTEXT_ROBUSTNESS_hint
	ContextRobustness enumHint
	// See https://www.glfw.org/docs/latest/window_guide.html#GLFW_CONTEXT_RELEASE_BEHAVIOR_hint
	ContextReleaseBehavior  enumHint
	OpenGLForwardCompatible boolHint
	OpenGLDebugContext      boolHint
	// See https://www.glfw.org/docs/latest/window_guide.html#GLFW_OPENGL_PROFILE_hint
	OpenGLProfile enumHint
	SRGBCapable   boolHint
}

func newContextHints() contextHints {
	return contextHints{
		ClientAPI:               enumHint{Value: OpenGLAPI, hint: hint{target: glfw.ClientAPI}},
		ContextVersionMajor:     intHint{Value: 3, hint: hint{target: glfw.ContextVersionMajor}},
		ContextVersionMinor:     intHint{Value: 3, hint: hint{target: glfw.ContextVersionMinor}},
		ContextRobustness:       enumHint{Value: NoRobustness, hint: hint{target: glfw.ContextRobustness}},
		ContextReleaseBehavior:  enumHint{Value: AnyReleaseBehavior, hint: hint{target: glfw.ContextReleaseBehavior}},
		OpenGLForwardCompatible: boolHint{Value: false, hint: hint{target: glfw.OpenGLForwardCompatible}},
		OpenGLDebugContext:      boolHint{Value: false, hint: hint{target: glfw.OpenGLDebugContext}},
		OpenGLProfile:           enumHint{Value: OpenGLAnyProfile, hint: hint{target: glfw.OpenGLProfile}},
		SRGBCapable:             boolHint{Value: false, hint: hint{target: glfw.SRGBCapable}},
	}
}

func (h contextHints) apply() {
	applyHint(h.ClientAPI)
	applyHint(h.ContextVersionMajor)
	applyHint(h.ContextVersionMinor)
	applyHint(h.ContextRobustness)
	applyHint(h.ContextReleaseBehavior)
	applyHint(h.OpenGLForwardCompatible)
	applyHint(h.OpenGLDebugContext)
	applyHint(h.OpenGLProfile)
	applyHint(h.SRGBCapable)
}

// Framebuffer related hints
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

func (h framebufferHints) apply() {
	applyHint(h.DepthBits)
	applyHint(h.Samples)
	applyHint(h.RefreshRate)
	applyHint(h.DoubleBuffer)
}

type customHints struct {
	Fullscreen boolHint
	Vsync      boolHint
}

func newCustomHints() customHints {
	return customHints{
		Fullscreen: boolHint{Value: false},
		Vsync:      boolHint{Value: true},
	}
}
