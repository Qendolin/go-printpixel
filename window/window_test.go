package window_test

import (
	"testing"

	"github.com/Qendolin/go-printpixel/layout"
	"github.com/Qendolin/go-printpixel/test"
	"github.com/Qendolin/go-printpixel/window"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	test.ParseArgs()
	m.Run()
}

func TestWindowNormal(t *testing.T) {
	hints := window.NewHints()

	win, err := window.New(hints, "Test Window", 1600, 900, nil)
	assert.NoError(t, err)
	win.Window = test.WrapWindow(win.Window)

	cfg := window.NewGlConfig(0)
	cfg.Debug = true
	go func() {
		for err := range cfg.Errors {
			if err.Fatal {
				t.Error(err)
			}
			t.Log(err)
		}
	}()

	win.Run(cfg)
	win.Close()
}

func TestScreenLayout(t *testing.T) {
	hints := window.NewHints()

	screenLo := layout.NewScreenByDimensions(1920, 1080)
	gridLo := layout.NewGrid([]layout.TrackDef{
		{Value: 1, Unit: layout.Percent},
	}, []layout.TrackDef{
		{Value: 0.5, Unit: layout.Percent},
		{Value: 0.5, Unit: layout.Percent},
	})
	screenLo.Child = &gridLo

	win, err := window.New(hints, "Test Window", 1600, 900, nil)
	assert.NoError(t, err)
	win.Window = test.WrapWindow(win.Window)

	gridLo.Children[0][0] = win
	screenLo.Layout()

	cfg := window.NewGlConfig(0)
	cfg.Debug = true
	go func() {
		for err := range cfg.Errors {
			if err.Fatal {
				t.Error(err)
			}
			t.Log(err)
		}
	}()

	win.Run(cfg)

	assert.Equal(t, 0, win.X())
	assert.Equal(t, 0, win.Y())
	assert.Equal(t, 1920, win.Width())
	assert.Equal(t, 1080/2, win.Height())

	win.Close()
}
