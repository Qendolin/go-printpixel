package window_test

import (
	"math"
	"os"
	"testing"

	_ "image/png"

	"github.com/Qendolin/go-printpixel/core"
	"github.com/Qendolin/go-printpixel/core/data"
	"github.com/Qendolin/go-printpixel/core/glcontext"
	"github.com/Qendolin/go-printpixel/core/glwindow"
	"github.com/Qendolin/go-printpixel/pkg/scene"
	"github.com/Qendolin/go-printpixel/pkg/test"
	"github.com/Qendolin/go-printpixel/pkg/window"
	"github.com/Qendolin/go-printpixel/renderer"
	"github.com/Qendolin/go-printpixel/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
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

type SpecialGraphic struct {
	scene.Graphic
}

func (s *SpecialGraphic) GetRenderer() string {
	return renderer.TextureQuad + "Special"
}

func TestDepthSorting(t *testing.T) {
	win, close := test.NewWindow(t)
	defer close()

	win.Renderers[renderer.TextureQuad+"Special"] = renderer.NewTextureQuadRenderer()

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)

	win.Child = &scene.Aspect{
		Ratio: 1.,
		Child: scene.Centered(&scene.Stack{
			Children: []scene.Layoutable{
				//Left (back) to right (front)
				&scene.Absolute{
					Child: &scene.Layer{Child: &scene.Layer{Child: &scene.Graphic{
						Texture: core.NewTexture2DFromBytes([]byte{255, 0, 0, 127}, 1, 1),
					}}},
					Unit: scene.Percent,
					DX:   0.125, DY: 0.75,
					W: 0.5, H: 0.4,
				},
				&scene.Absolute{
					Child: &scene.Layer{Child: &SpecialGraphic{
						Graphic: scene.Graphic{Texture: core.NewTexture2DFromBytes([]byte{0, 255, 0, 127}, 1, 1)},
					}},
					Unit: scene.Percent,
					DX:   0.5, DY: 0.75,
					W: 0.5, H: 0.4,
				},
				&scene.Absolute{
					Child: &scene.Graphic{
						Texture: core.NewTexture2DFromBytes([]byte{0, 0, 255, 127}, 1, 1),
					},
					Unit: scene.Percent,
					DX:   0.875, DY: 0.75,
					W: 0.5, H: 0.4,
				},
				//Left (front) to right (back)
				&scene.Absolute{
					Child: &scene.Graphic{
						Texture: core.NewTexture2DFromBytes([]byte{255, 0, 0, 127}, 1, 1),
					},
					Unit: scene.Percent,
					DX:   0.125, DY: 0.25,
					W: 0.5, H: 0.4,
				},
				&scene.Absolute{
					Child: &scene.Layer{Child: &SpecialGraphic{
						Graphic: scene.Graphic{Texture: core.NewTexture2DFromBytes([]byte{0, 255, 0, 127}, 1, 1)},
					}},
					Unit: scene.Percent,
					DX:   0.5, DY: 0.25,
					W: 0.5, H: 0.4,
				},
				&scene.Absolute{
					Child: &scene.Layer{Child: &scene.Layer{Child: &scene.Graphic{
						Texture: core.NewTexture2DFromBytes([]byte{0, 0, 255, 127}, 1, 1),
					}}},
					Unit: scene.Percent,
					DX:   0.875, DY: 0.25,
					W: 0.5, H: 0.4,
				},
			},
		}),
	}

	win.Layout()
	win.Run()
}

func TestWindowPosition(t *testing.T) {
	win, close := test.NewWindow(t)
	defer close()

	crossTexData := make([]byte, 16*16)
	for i := range crossTexData {
		if (i%16+1)/2 == 4 || i >= 7*16 && i < 9*16 {
			crossTexData[i] = 255
		}
	}
	corssTex := data.NewTexture2D(nil, data.Tex2DTarget2D)
	corssTex.Bind(0)
	corssTex.AllocBytes(crossTexData, 0, gl.RGB, 16, 16, gl.RED)
	corssTex.ApplyDefaults()

	abs := &scene.Absolute{
		Child: &scene.Absolute{
			Child: &scene.Graphic{
				Texture: corssTex,
			},
			W: 16,
			H: 16,
		},
		Unit: scene.Percent,
		W:    1,
		H:    1,
	}
	win.Child = abs

	win.GlWindow.SetSizeCallback(func(_ glwindow.Extended, _ int, _ int) {
		win.Layout()
	})

	dir := 0.0
	win.BeforeUpdate = func() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		dir += win.GlWindow.Delta().Seconds() / 4 * math.Pi
		abs.DX = .5 * float32(math.Cos(dir))
		abs.DY = .5 * float32(math.Sin(dir))
		scene.Layout(abs)
	}

	win.Layout()
	win.Run()
}
