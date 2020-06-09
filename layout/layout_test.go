package layout_test

import (
	"testing"

	"github.com/Qendolin/go-printpixel/layout"
	"github.com/stretchr/testify/assert"
)

func TestGrid(t *testing.T) {

	var (
		width  = 100
		height = 300
	)

	grid := layout.NewGrid([]layout.TrackDef{
		{.25, layout.Percent},
		{.75, layout.Percent},
	}, []layout.TrackDef{
		{100, layout.Pixel},
		{200, layout.Pixel},
	})
	grid.SetWidth(width)
	grid.SetHeight(height)
	grid.Children[0][0] = &layout.SimpleBox{}
	grid.Children[0][1] = &layout.SimpleBox{}
	grid.Children[1][0] = &layout.SimpleBox{}
	grid.Children[1][1] = &layout.SimpleBox{}

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
