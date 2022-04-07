package window

import (
	"github.com/Qendolin/go-printpixel/core/glw"
)

type SimpleConfig struct {
	Title         string
	Width         int
	Height        int
	NoVsync       bool
	FixedSize     bool
	Hidden        bool
	Maximized     bool
	Debug         bool
	Multisampling int
	DebugHandler  glw.DebugHandler
}

func (simple *SimpleConfig) ToFullConfig() glw.Config {
	conf := glw.DefaultConfig()
	conf.Title = simple.Title
	conf.Visible = !simple.Hidden
	conf.Vsync = !simple.NoVsync
	conf.Resizable = !simple.FixedSize
	conf.Maximized = simple.Maximized
	conf.Samples = simple.Multisampling
	conf.DebugContext = simple.Debug
	conf.DebugHandler = simple.DebugHandler
	conf.Width = simple.Width
	conf.Height = simple.Height
	return conf
}
