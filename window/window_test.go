package window_test

import (
	"os"
	"testing"

	_ "image/png"

	"github.com/Qendolin/go-printpixel/internal/utils"
	"github.com/Qendolin/go-printpixel/layout"
	"github.com/Qendolin/go-printpixel/test"
	"github.com/Qendolin/go-printpixel/window"
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
	win.Window = test.WrapWindow(win.Window)
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
	screenLo.Layout()

	win.Run()

	assert.Equal(t, 0, win.X())
	assert.Equal(t, 0, win.Y())
	assert.Equal(t, 1920, win.Width())
	assert.Equal(t, 1080/2, win.Height())
}

func TestViewer(t *testing.T) {
	win, close := test.NewWindow(t)
	defer close()

	absPath, err := utils.ResolvePath("assets/textures/uv.png")
	assert.NoError(t, err)
	imgFile, err := os.Open(absPath)
	assert.NoError(t, err)
	defer imgFile.Close()
	v := layout.NewViewer()
	v.Target.Texture.Bind(0)
	err = v.Target.Texture.AllocFile(imgFile, 0, gl.RGBA, gl.RGBA)
	assert.NoError(t, err)
	v.Target.Texture.ApplyDefaults()
	win.Child = v

	prevW := win.Width()
	prevH := win.Height()
	win.AfterUpdate = func() {
		if prevW != win.Width() || prevH != win.Height() {
			prevW = win.Width()
			prevH = win.Height()
			win.Layout()
			v.Draw()
			win.Update()
			v.Draw()
		}
	}

	v.Draw()
	win.Update()
	v.Draw()
	win.Run()
}
