package window

import (
	"github.com/Qendolin/go-printpixel/core/glcontext"
	"github.com/Qendolin/go-printpixel/core/glwindow"
)

type SimpleConfig struct {
	Width        int
	Height       int
	NoVsync      bool
	Unresizeable bool
	Invisible    bool
	Maximized    bool
	Debug        bool
	errorsIn     <-chan glcontext.GlError
	errorsOut    chan glcontext.GlError
}

func (cfg *SimpleConfig) ToHints() glwindow.Hints {
	h := NewHints()
	h.Visible.Value = !cfg.Invisible
	h.Vsync.Value = !cfg.NoVsync
	h.Resizable.Value = !cfg.Unresizeable
	h.Maximized.Value = cfg.Maximized
	return h
}

func (cfg *SimpleConfig) Errors() <-chan glcontext.GlError {
	if cfg.errorsOut != nil {
		return cfg.errorsOut
	}

	cfg.errorsOut = make(chan glcontext.GlError, 0)
	if cfg.errorsIn != nil {
		go cfg.pipeErrors()
	}
	return cfg.errorsOut
}

func (cfg *SimpleConfig) ToGlConfig() glcontext.GlConfig {
	size := 1
	if cfg.Debug {
		size = 0
	}
	glCfg := NewGlConfig(size)
	glCfg.Debug = cfg.Debug
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
