package window_test

import (
	"os"
	"testing"

	_ "image/png"

	"github.com/Qendolin/go-printpixel/pkg/layout"
	"github.com/Qendolin/go-printpixel/pkg/test"
	"github.com/Qendolin/go-printpixel/pkg/window"
	"github.com/Qendolin/go-printpixel/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	test.ParseArgs()
	m.Run()
}

func TestWindowNormal(t *testing.T) {
	hints := window.NewHints()
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

	win, err := window.NewCustom("Test Window", 1600, 900, hints, nil, cfg)
	assert.NoError(t, err)
	win.GlWindow = test.WrapWindow(win.GlWindow)
	win.Run()
}

func TestScreenLayout(t *testing.T) {
	win, close := test.NewWindow(t)
	defer close()

	screenLo := layout.NewScreenByDimensions(1920, 1080)
	gridLo := layout.NewGrid([]layout.TrackDef{
		{Value: 1, Unit: layout.Percent},
	}, []layout.TrackDef{
		{Value: 0.5, Unit: layout.Percent},
		{Value: 0.5, Unit: layout.Percent},
	})
	screenLo.Child = &gridLo

	gridLo.Children[0][0] = win
	layout.Layout(screenLo)

	win.Run()

	assert.Equal(t, 0, win.X())
	assert.Equal(t, 0, win.Y())
	assert.Equal(t, 1920, win.Width())
	assert.Equal(t, 1080/2, win.Height())
}

func TestGraphic(t *testing.T) {
	win, close := test.NewWindow(t)
	defer close()

	absPath, err := utils.ResolvePath("assets/textures/uv.png")
	assert.NoError(t, err)
	imgFile, err := os.Open(absPath)
	assert.NoError(t, err)
	defer imgFile.Close()

	v := layout.NewGraphic()
	v.Texture.Bind(0)
	v.Texture.ApplyDefaults()
	err = v.Texture.AllocFile(imgFile, 0, gl.RGBA, gl.RGBA)
	assert.NoError(t, err)
	win.Child = v

	prevW := win.Width()
	prevH := win.Height()
	win.BeforeUpdate = func() {
		if prevW != win.Width() || prevH != win.Height() {
			prevW = win.Width()
			prevH = win.Height()
			win.Layout()
		}
	}

	win.Layout()
	win.Run()
}
