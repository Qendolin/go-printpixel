package window

import (
	"github.com/Qendolin/go-printpixel/core/glcontext"
	"github.com/Qendolin/go-printpixel/core/glwindow"
)

type SimpleConfig struct {
	Width         int
	Height        int
	NoVsync       bool
	Unresizeable  bool
	Invisible     bool
	Maximized     bool
	Debug         bool
	Multisampling int
	errorsIn      <-chan glcontext.Error
	errorsOut     chan glcontext.Error
}

func (cfg *SimpleConfig) ToHints() glwindow.Hints {
	h := glwindow.NewHints()
	h.Visible.Value = !cfg.Invisible
	h.Vsync.Value = !cfg.NoVsync
	h.Resizable.Value = !cfg.Unresizeable
	h.Maximized.Value = cfg.Maximized
	h.Samples.Value = cfg.Multisampling
	return h
}

func (cfg *SimpleConfig) Errors() <-chan glcontext.Error {
	if cfg.errorsOut != nil {
		return cfg.errorsOut
	}

	cfg.errorsOut = make(chan glcontext.Error)
	if cfg.errorsIn != nil {
		go cfg.pipeErrors()
	}
	return cfg.errorsOut
}

func (cfg *SimpleConfig) ToGlConfig() glcontext.Config {
	size := 1
	if cfg.Debug {
		size = 0
	}
	glCfg := glcontext.NewGlConfig(size)
	glCfg.Debug = cfg.Debug
	glCfg.Multisampling = cfg.Multisampling > 1
	cfg.errorsIn = glCfg.Errors
	if cfg.errorsOut != nil {
		go cfg.pipeErrors()
	}
	return glCfg
}

func (cfg *SimpleConfig) pipeErrors() {
	for err := range cfg.errorsIn {
		cfg.errorsOut <- err
	}
}
