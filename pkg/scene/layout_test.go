package scene_test

import (
	"testing"

	"github.com/Qendolin/go-printpixel/pkg/scene"
	"github.com/Qendolin/go-printpixel/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	test.ParseArgs()
	m.Run()
}

func TestGrid(t *testing.T) {

	var (
		width  = 100
		height = 300
	)

	grid := scene.Grid{
		Cols: []scene.TrackDef{
			{.25, scene.Percent},
			{.75, scene.Percent},
		},
		Rows: []scene.TrackDef{
			{100, scene.Pixel},
			{200, scene.Pixel},
		},
	}
	grid.Init()
	grid.SetWidth(width)
	grid.SetHeight(height)
	grid.Children[0][0] = &scene.Box{}
	grid.Children[0][1] = &scene.Box{}
	grid.Children[1][0] = &scene.Box{}
	grid.Children[1][1] = &scene.Box{}

	grid.Layout()

	epsilon := float64(1.5)

	assert.Equal(t, 0, grid.Children[0][0].X())
	assert.Equal(t, 0, grid.Children[0][0].Y())
	assert.InEpsilon(t, float32(width)*0.25, grid.Children[0][0].Width(), epsilon)
	assert.Equal(t, 100, grid.Children[0][0].Height())

	assert.Equal(t, 0, grid.Children[0][1].X())
	assert.Equal(t, 100, grid.Children[0][1].Y())
	assert.InEpsilon(t, float32(width)*0.25, grid.Children[0][1].Width(), epsilon)
	assert.Equal(t, 200, grid.Children[0][1].Height())

	assert.InEpsilon(t, float32(width)*0.25, grid.Children[1][0].X(), epsilon)
	assert.Equal(t, 0, grid.Children[1][0].Y())
	assert.InEpsilon(t, float32(width)*0.75, grid.Children[1][0].Width(), epsilon)
	assert.Equal(t, 100, grid.Children[1][0].Height())

	assert.InEpsilon(t, float32(width)*0.25, grid.Children[1][1].X(), epsilon)
	assert.Equal(t, 100, grid.Children[1][1].Y())
	assert.InEpsilon(t, float32(width)*0.75, grid.Children[1][1].Width(), epsilon)
	assert.Equal(t, 200, grid.Children[1][1].Height())
}

func TestAspect(t *testing.T) {
	a := scene.Aspect{
		Child: &scene.Box{},
		Ratio: 1,
		Mode:  scene.Contain,
	}

	a.SetWidth(100)
	a.SetHeight(100)
	a.Layout()
	assert.Equal(t, 100, a.Child.Width())
	assert.Equal(t, 100, a.Child.Height())

	a.SetWidth(200)
	a.SetHeight(100)
	a.Layout()
	assert.Equal(t, 100, a.Child.Width())
	assert.Equal(t, 100, a.Child.Height())

	a.SetWidth(100)
	a.SetHeight(200)
	a.Layout()
	assert.Equal(t, 100, a.Child.Width())
	assert.Equal(t, 100, a.Child.Height())

	a.Mode = scene.Cover

	a.SetWidth(100)
	a.SetHeight(100)
	a.Layout()
	assert.Equal(t, 100, a.Child.Width())
	assert.Equal(t, 100, a.Child.Height())

	a.SetWidth(200)
	a.SetHeight(100)
	a.Layout()
	assert.Equal(t, 200, a.Child.Width())
	assert.Equal(t, 200, a.Child.Height())

	a.SetWidth(100)
	a.SetHeight(200)
	a.Layout()
	assert.Equal(t, 200, a.Child.Width())
	assert.Equal(t, 200, a.Child.Height())

	a.Ratio = 2 / 1
	a.Mode = scene.FitHieght

	a.SetWidth(100)
	a.SetHeight(200)
	a.Layout()
	assert.Equal(t, 400, a.Child.Width())
	assert.Equal(t, 200, a.Child.Height())

	a.SetWidth(200)
	a.SetHeight(100)
	a.Layout()
	assert.Equal(t, 200, a.Child.Width())
	assert.Equal(t, 100, a.Child.Height())

	a.SetWidth(100)
	a.SetHeight(100)
	a.Layout()
	assert.Equal(t, 200, a.Child.Width())
	assert.Equal(t, 100, a.Child.Height())

	a.Ratio = 1. / 2.
	a.Mode = scene.FitWidth

	a.SetWidth(200)
	a.SetHeight(100)
	a.Layout()
	assert.Equal(t, 200, a.Child.Width())
	assert.Equal(t, 400, a.Child.Height())

	a.Ratio = 1
	a.Mode = scene.FitWidth

	a.SetWidth(200)
	a.SetHeight(100)
	a.Layout()
	assert.Equal(t, 200, a.Child.Width())
	assert.Equal(t, 200, a.Child.Height())
}
