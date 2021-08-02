package field

import (
	"image"
	"image/draw"
	"math"
)

type ValueFieldI struct {
	Width, Height int
	ComponentSize int
	Values        []uint32
}

func NewValueFieldI(vf *ValueField) *ValueFieldI {
	vfi := &ValueFieldI{
		Width:         vf.Width,
		Height:        vf.Height,
		ComponentSize: vf.ComponentSize,
		Values:        make([]uint32, len(vf.Values)),
	}
	for i, f := range vf.Values {
		vfi.Values[i] = uint32(float64(f) * float64(math.MaxUint32))
	}
	return vfi
}

type ValueFieldI8 struct {
	Width, Height int
	ComponentSize int
	Values        []uint8
}

func NewValueFieldI8(vf *ValueField) *ValueFieldI8 {
	vfi := &ValueFieldI8{
		Width:         vf.Width,
		Height:        vf.Height,
		ComponentSize: vf.ComponentSize,
		Values:        make([]uint8, len(vf.Values)),
	}
	for i, f := range vf.Values {
		vfi.Values[i] = uint8(float64(f) * float64(math.MaxUint8))
	}
	return vfi
}

func (vf *ValueFieldI) Get(x, y, dim int) uint32 {
	componentIdx := x + y*vf.Width
	return vf.Values[componentIdx*vf.ComponentSize+dim]
}

func (vf *ValueFieldI) GetComponent(x, y int) []uint32 {
	componentIdx := x + y*vf.Width
	return vf.Values[componentIdx*vf.ComponentSize : componentIdx*vf.ComponentSize+vf.ComponentSize]
}

func (vf *ValueFieldI) Set(x, y, dim int, v uint32) {
	componentIdx := x + y*vf.Width
	vf.Values[componentIdx*vf.ComponentSize+dim] = v
}

func (vf *ValueFieldI) SetComponent(x, y int, c []uint32) {
	copy(vf.GetComponent(x, y), c)
}

type ValueField struct {
	Width, Height int
	ComponentSize int
	Values        []float32
}

func NewValueField(w, h, c int) *ValueField {
	return &ValueField{
		Width:         w,
		Height:        h,
		ComponentSize: c,
		Values:        make([]float32, w*h*c),
	}
}

func RGBValueField(img image.Image) *ValueField {
	br := img.Bounds()
	rgba, ok := img.(*image.RGBA)
	if !ok {
		rgba = image.NewRGBA(br)
		draw.Draw(rgba, br, img, br.Min, draw.Src)
	}

	pix := rgba.Pix

	vf := ValueField{
		Width:         br.Dx(),
		Height:        br.Dy(),
		ComponentSize: 3,
		Values:        make([]float32, len(rgba.Pix)/4*3),
	}

	for y := br.Min.Y; y < br.Max.Y; y++ {
		for x := br.Min.X; x < br.Max.X; x++ {
			rgb8 := pix[y*rgba.Stride+x*4 : y*rgba.Stride+x*4+3]

			vf.SetComponent(x+br.Min.X, y+br.Min.Y, []float32{
				float32(rgb8[0]) / float32(math.MaxUint8),
				float32(rgb8[1]) / float32(math.MaxUint8),
				float32(rgb8[2]) / float32(math.MaxUint8),
			})
		}
	}

	return &vf
}

func RGBAValueField(img image.Image) *ValueField {
	br := img.Bounds()
	rgba, ok := img.(*image.RGBA)
	if !ok {
		rgba = image.NewRGBA(br)
		draw.Draw(rgba, br, img, br.Min, draw.Src)
	}

	pix := rgba.Pix

	vf := ValueField{
		Width:         br.Dx(),
		Height:        br.Dy(),
		ComponentSize: 4,
		Values:        make([]float32, len(rgba.Pix)),
	}

	for y := br.Min.Y; y < br.Max.Y; y++ {
		for x := br.Min.X; x < br.Max.X; x++ {
			rgb8 := pix[y*rgba.Stride+x*4 : y*rgba.Stride+x*4+4]

			vf.SetComponent(x+br.Min.X, y+br.Min.Y, []float32{
				float32(rgb8[0]) / float32(math.MaxUint8),
				float32(rgb8[1]) / float32(math.MaxUint8),
				float32(rgb8[2]) / float32(math.MaxUint8),
				float32(rgb8[3]) / float32(math.MaxUint8),
			})
		}
	}

	return &vf
}

func (vf *ValueField) Get(x, y, dim int) float32 {
	componentIdx := x + y*vf.Width
	return vf.Values[componentIdx*vf.ComponentSize+dim]
}

func (vf *ValueField) GetComponent(x, y int) []float32 {
	componentIdx := x + y*vf.Width
	return vf.Values[componentIdx*vf.ComponentSize : componentIdx*vf.ComponentSize+vf.ComponentSize]
}

func (vf *ValueField) Set(x, y, dim int, v float32) {
	componentIdx := x + y*vf.Width
	vf.Values[componentIdx*vf.ComponentSize+dim] = v
}

func (vf *ValueField) SetComponent(x, y int, c []float32) {
	copy(vf.GetComponent(x, y), c)
}

type MSDF struct {
	width, height int
	Values        []float32
}

func NewMSDF(vf *ValueField) *MSDF {
	msdf := MSDF{
		width:  vf.Width,
		height: vf.Height,
		Values: make([]float32, len(vf.Values)/vf.ComponentSize),
	}
	for y := 0; y < vf.Height; y++ {
		for x := 0; x < vf.Width; x++ {
			comp := vf.GetComponent(x, y)
			msdf.Set(x, y, median(comp[0], comp[1], comp[2]))
		}
	}
	return &msdf
}

func (msdf *MSDF) Get(x, y int) float32 {
	return msdf.Values[x+y*msdf.width]
}

func (msdf *MSDF) Set(x, y int, v float32) {
	msdf.Values[x+y*msdf.width] = v
}

func (msdf *MSDF) Width() int {
	return msdf.width
}

func (msdf *MSDF) Height() int {
	return msdf.height
}

func (msdf *MSDF) Raw() []float32 {
	return msdf.Values
}

// func median(a, b, c float32) float32 {
// 	x := a - b
// 	y := b - c
// 	z := a - c
// 	if x*y > 0 {
// 		return b
// 	}
// 	if x*z > 0 {
// 		return c
// 	}
// 	return a
// }

// func median(a, b, c float32) float32 {
// 	if a > b {
// 		a, b = b, a
// 	}
// 	if b > c {
// 		b = c
// 	}
// 	if a > b {
// 		b = a
// 	}
// 	return b
// }

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func median(a, b, c float32) float32 {
	return max(min(a, b), min(max(a, b), c))
}

type ScalarField interface {
	Width() int
	Height() int
	Get(x, y int) float32
	Raw() []float32
}
