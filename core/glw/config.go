package glw

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

func btoi(b bool) int {
	if b {
		return glfw.True
	}
	return glfw.False
}

type Config struct {
	WindowHints
	ContextHints
	FramebufferHints
	MonitorHints
	ExtendedConfig
}

func DefaultConfig() Config {
	return Config{
		DefaultWindowHints(),
		DefaultContextHints(),
		DefaultFramebufferHints(),
		DefaultMonitorHints(),
		DefaultExtendedConfig(),
	}
}

func BasicConfig(title string, width, height int, x, y int) Config {
	conf := DefaultConfig()
	conf.Title = title
	conf.Width = width
	conf.Height = height
	conf.X = x
	conf.Y = y
	return conf
}

func (c *Config) load() {
	c.ContextHints.load()
	c.WindowHints.load()
	c.FramebufferHints.load()
	c.MonitorHints.load()
	c.ExtendedConfig.load(c)
}

// Window related hints
// See https://www.glfw.org/docs/latest/window_guide.html
type WindowHints struct {
	Focused                bool
	Visible                bool
	Resizable              bool
	Decorated              bool
	Floating               bool
	AutoIconify            bool
	Maximized              bool
	ScaleToMonitor         bool
	TransparentFramebuffer bool
}

func DefaultWindowHints() WindowHints {
	return WindowHints{
		Focused:                true,
		Visible:                true,
		Resizable:              true,
		Decorated:              true,
		Floating:               false,
		AutoIconify:            true,
		Maximized:              false,
		ScaleToMonitor:         false,
		TransparentFramebuffer: false,
	}
}

func (h WindowHints) load() {
	glfw.WindowHint(glfw.Resizable, btoi(h.Resizable))
	glfw.WindowHint(glfw.Visible, btoi(h.Visible))
	glfw.WindowHint(glfw.Decorated, btoi(h.Decorated))
	glfw.WindowHint(glfw.Focused, btoi(h.Focused))
	glfw.WindowHint(glfw.AutoIconify, btoi(h.AutoIconify))
	glfw.WindowHint(glfw.Floating, btoi(h.Floating))
	glfw.WindowHint(glfw.Maximized, btoi(h.Maximized))
	glfw.WindowHint(glfw.ScaleToMonitor, btoi(h.ScaleToMonitor))
	glfw.WindowHint(glfw.TransparentFramebuffer, btoi(h.TransparentFramebuffer))
}

// Context related hints
type ContextHints struct {
	// See https://www.glfw.org/docs/latest/window_guide.html#GLFW_CLIENT_API_hint
	ClientAPI           GlfwEnum
	ContextCreationAPI  GlfwEnum
	ContextVersionMajor int
	ContextVersionMinor int
	// See https://www.glfw.org/docs/latest/window_guide.html#GLFW_CONTEXT_ROBUSTNESS_hint
	ContextRobustness GlfwEnum
	// See https://www.glfw.org/docs/latest/window_guide.html#GLFW_CONTEXT_RELEASE_BEHAVIOR_hint
	ContextReleaseBehavior  GlfwEnum
	OpenGLForwardCompatible bool
	// Sets glDebugMessageCallback and enables GL_DEBUG_OUTPUT and GL_DEBUG_OUTPUT_SYNCHRONOUS
	DebugContext bool
	// See https://www.glfw.org/docs/latest/window_guide.html#GLFW_OPENGL_PROFILE_hint
	OpenGLProfile GlfwEnum
	SRGBCapable   bool
}

func DefaultContextHints() ContextHints {
	return ContextHints{
		ClientAPI:               OpenGLAPI,
		ContextCreationAPI:      NativeContextAPI,
		ContextVersionMajor:     3,
		ContextVersionMinor:     3,
		ContextRobustness:       NoRobustness,
		ContextReleaseBehavior:  AnyReleaseBehavior,
		OpenGLForwardCompatible: false,
		DebugContext:            false,
		OpenGLProfile:           OpenGLAnyProfile,
		SRGBCapable:             false,
	}
}

func (h ContextHints) load() {
	glfw.WindowHint(glfw.ClientAPI, int(h.ClientAPI))
	glfw.WindowHint(glfw.ContextCreationAPI, int(h.ContextCreationAPI))
	glfw.WindowHint(glfw.ContextVersionMajor, h.ContextVersionMajor)
	glfw.WindowHint(glfw.ContextVersionMinor, h.ContextVersionMinor)
	glfw.WindowHint(glfw.ContextRobustness, int(h.ContextRobustness))
	glfw.WindowHint(glfw.ContextReleaseBehavior, int(h.ContextReleaseBehavior))
	glfw.WindowHint(glfw.OpenGLForwardCompatible, btoi(h.OpenGLForwardCompatible))
	glfw.WindowHint(glfw.OpenGLDebugContext, btoi(h.DebugContext))
	glfw.WindowHint(glfw.OpenGLProfile, int(h.OpenGLProfile))
	glfw.WindowHint(glfw.SRGBCapable, btoi(h.SRGBCapable))
}

// Framebuffer related hints
type FramebufferHints struct {
	DepthBits, StencilBits       int
	RedBits, GreenBits, BlueBits int
	Samples                      int
	DoubleBuffer                 bool
}

func DefaultFramebufferHints() FramebufferHints {
	return FramebufferHints{
		DepthBits:    24,
		StencilBits:  8,
		RedBits:      8,
		GreenBits:    8,
		BlueBits:     8,
		Samples:      0,
		DoubleBuffer: true,
	}
}

func (h FramebufferHints) load() {
	glfw.WindowHint(glfw.DepthBits, h.DepthBits)
	glfw.WindowHint(glfw.StencilBits, h.StencilBits)
	glfw.WindowHint(glfw.RedBits, h.RedBits)
	glfw.WindowHint(glfw.GreenBits, h.GreenBits)
	glfw.WindowHint(glfw.BlueBits, h.BlueBits)
	glfw.WindowHint(glfw.Samples, h.Samples)
	glfw.WindowHint(glfw.DoubleBuffer, btoi(h.DoubleBuffer))
}

type MonitorHints struct {
	RefreshRate int
}

func DefaultMonitorHints() MonitorHints {
	return MonitorHints{
		RefreshRate: DontCare,
	}
}

func (h MonitorHints) load() {
	glfw.WindowHint(glfw.RefreshRate, h.RefreshRate)
}

type ExtendedConfig struct {
	Fullscreen          bool
	Vsync               bool
	Width, Height       int
	MinWidth, MinHeight int
	MaxWidth, MaxHeight int
	AspectRatio         float64
	X, Y                int
	Title               string
	Monitor             *glfw.Monitor
	SharedContext       *glfw.Window
	DebugHandler        DebugHandler
}

func DefaultExtendedConfig() ExtendedConfig {
	return ExtendedConfig{
		Fullscreen:    false,
		Vsync:         true,
		Width:         DontCare,
		Height:        DontCare,
		MinWidth:      DontCare,
		MinHeight:     DontCare,
		MaxWidth:      DontCare,
		MaxHeight:     DontCare,
		AspectRatio:   0,
		X:             DontCare,
		Y:             DontCare,
		Title:         "",
		Monitor:       nil,
		SharedContext: nil,
		DebugHandler:  DefaultDebugHandler,
	}
}

func (c *ExtendedConfig) load(conf *Config) {
	if c.Monitor == nil && c.Fullscreen {
		c.Monitor = glfw.GetPrimaryMonitor()
	}

	vidMode := glfw.GetPrimaryMonitor().GetVideoMode()
	if vidMode == nil || vidMode.Width == 0 || vidMode.Height == 0 {
		vidMode = &glfw.VidMode{
			Width:  1920,
			Height: 1080,
		}
	}

	if !conf.Maximized {
		if c.Width == DontCare {
			c.Width = vidMode.Width / 2
		}
		if c.Height == DontCare {
			c.Height = vidMode.Height / 2
		}
		if c.X == DontCare {
			c.X = (vidMode.Width - c.Width) / 2
		}
		if c.Y == DontCare {
			c.Y = (vidMode.Height - c.Height) / 2
		}
	}
}
