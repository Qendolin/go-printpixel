package context

import (
	"log"
	"os"
	"testing"

	"github.com/Qendolin/go-printpixel/internal/window"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	window.Init()
	hints := window.NewHints()
	hints.ContextVersionMajor.Value = 3
	hints.ContextVersionMinor.Value = 2
	win, err := window.NewWindow(hints, "Test Window", 800, 450, nil)
	assert.NoError(t, err)
	win.MakeContextCurrent()
	cfg := NewGlConfig()
	cfg.Debug = true
	cfg.Logger = log.New(os.Stdout, "Test ", log.LstdFlags)
	err = Init(cfg)
	assert.NoError(t, err)
	win.Destroy()
	window.Terminate()
}
