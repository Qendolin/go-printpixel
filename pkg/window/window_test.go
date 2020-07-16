package window_test

import (
	"os"
	"testing"

	_ "image/png"

	"github.com/Qendolin/go-printpixel/core"
	"github.com/Qendolin/go-printpixel/core/glcontext"
	"github.com/Qendolin/go-printpixel/core/glwindow"
	"github.com/Qendolin/go-printpixel/pkg/scene"
	"github.com/Qendolin/go-printpixel/pkg/test"
	"github.com/Qendolin/go-printpixel/pkg/window"
	"github.com/Qendolin/go-printpixel/utils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	test.ParseArgs()
	m.Run()
}

func TestWindowNormal(t *testing.T) {
	hints := glwindow.NewHints()
	cfg := glcontext.NewGlConfig(0)
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

	screenLo := scene.NewScreenByDimensions(1920, 1080)
	gridLo := scene.NewGrid([]scene.TrackDef{
		{Value: 1, Unit: scene.Percent},
	}, []scene.TrackDef{
		{Value: 0.5, Unit: scene.Percent},
		{Value: 0.5, Unit: scene.Percent},
	})
	screenLo.Child = &gridLo

	gridLo.Children[0][0] = win
	scene.Layout(screenLo)

	win.Run()

	assert.Equal(t, 0, win.X())
	assert.Equal(t, 0, win.Y())
	assert.Equal(t, 1920, win.Width())
	assert.Equal(t, 1080/2, win.Height())
}

func TestGraphic(t *testing.T) {
	win, close := test.NewWindow(t)
	defer close()

	absPath, err := utils.ResolvePath("res://assets/textures/uv.png")
	assert.NoError(t, err)
	imgFile, err := os.Open(absPath)
	assert.NoError(t, err)
	defer imgFile.Close()
	tex, err := core.NewTexture2DFromFile(imgFile)
	assert.NoError(t, err)

	win.Child = &scene.Graphic{
		Texture: tex,
	}

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
