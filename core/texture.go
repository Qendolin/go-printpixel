package core

import (
	"bytes"
	"errors"
	"image"
	"image/draw"
	"io"
	"io/ioutil"
	"math"
	"os"

	"github.com/Qendolin/go-printpixel/core/data"
	"github.com/Qendolin/go-printpixel/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type errorReader struct {
	error
}

func (er errorReader) Read(_ []byte) (int, error) {
	return 0, er.error
}

type ArrayReader struct {
	Stride int
	Array  []io.Reader
	n      int64
}

func (ar *ArrayReader) Read(b []byte) (n int, err error) {
	defer func() {
		ar.n += int64(n)
	}()

	i := int(ar.n / int64(ar.Stride))
	left := ar.Stride - int(ar.n%int64(ar.Stride))
	if i >= len(ar.Array) {
		return 0, io.EOF
	}
	r := ar.Array[i]
	n, err = r.Read(b)
	min := n
	if left < n {
		min = left
	}
	if errors.Is(err, io.EOF) {
		if n == 0 {
			return copy(b, make([]byte, left)), nil
		}
	} else if err != nil {
		return n, err
	}
	return min, nil
}

func PixelReader(img image.Image) (w, h int, r io.Reader) {
	var buf []byte
	switch i := img.(type) {
	case *image.Uniform:
		r, g, b, a := i.RGBA()
		buf = []byte{byte(r >> 8), byte(g >> 8), byte(b >> 8), byte(a >> 8)}
		w, h = 1, 1
	case *image.RGBA:
		buf = i.Pix
		w, h = i.Rect.Size().X, i.Rect.Size().Y
	default:
		rgba := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
		draw.Draw(rgba, rgba.Rect, img, image.Point{}, draw.Src)
		buf = rgba.Pix
		w, h = rgba.Rect.Size().X, rgba.Rect.Size().Y
	}
	return w, h, bytes.NewReader(buf)
}

func ImageReader(file io.Reader) (w, h int, r io.Reader) {
	img, _, err := image.Decode(file)
	if err != nil {
		return 0, 0, errorReader{err}
	}
	return PixelReader(img)
}

func FileReader(path string) (w, h int, r io.Reader) {
	path, err := utils.ResolvePath(path)
	if err != nil {
		return 0, 0, errorReader{err}
	}
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, errorReader{err}
	}
	return ImageReader(f)
}

func NewTexture1DEmpty(w int) *data.Texture1D {
	t := data.NewTexture1D(nil, data.Tex1DTarget1D)
	t.Bind(0)
	t.AllocEmpty(0, gl.RGBA, int32(w), gl.RGBA)
	t.ApplyDefaults()
	t.Unbind(0)
	return t
}

func NewTexture1D(w int, r io.Reader) (*data.Texture1D, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	t := data.NewTexture1D(nil, data.Tex1DTarget1D)
	t.Bind(0)
	t.AllocBytes(b, 0, gl.RGBA, int32(w), gl.RGBA)
	t.ApplyDefaults()
	t.Unbind(0)
	return t, nil
}

func NewTexture1DBytes(w int, b []byte) *data.Texture1D {
	t := data.NewTexture1D(nil, data.Tex1DTarget1D)
	t.Bind(0)
	t.AllocBytes(b, 0, gl.RGBA, int32(w), gl.RGBA)
	t.ApplyDefaults()
	t.Unbind(0)
	return t
}

func NewTexture2DEmpty(w, h int) *data.Texture2D {
	t := data.NewTexture2D(nil, data.Tex2DTarget2D)
	t.Bind(0)
	t.AllocEmpty(0, gl.RGBA, int32(w), int32(h), gl.RGBA)
	t.ApplyDefaults()
	t.Unbind(0)
	return t
}

func NewTexture2D(w, h int, r io.Reader) (*data.Texture2D, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	t := data.NewTexture2D(nil, data.Tex2DTarget2D)
	t.Bind(0)
	t.AllocBytes(b, 0, gl.RGBA, int32(w), int32(h), gl.RGBA)
	t.ApplyDefaults()
	t.Unbind(0)
	return t, nil
}

func NewTexture2DBytes(w, h int, b []byte) *data.Texture2D {
	t := data.NewTexture2D(nil, data.Tex2DTarget2D)
	t.Bind(0)
	t.AllocBytes(b, 0, gl.RGBA, int32(w), int32(h), gl.RGBA)
	t.ApplyDefaults()
	t.Unbind(0)
	return t
}

func LoadTexture(path string) *data.Texture2D {
	path, err := utils.ResolvePath(path)
	if err != nil {
		return newErrorTexture("not found")
	}
	f, err := os.Open(path)
	if err != nil {
		return newErrorTexture("not found")
	}
	t, err := NewTexture2D(ImageReader(f))
	if err != nil {
		return newErrorTexture("not decodeable")
	}
	return t
}

func newErrorTexture(cause string) *data.Texture2D {
	buf := make([]byte, 64*64*3)
	switch cause {
	case "not found":
		for x := 0; x < 64; x++ {
			for y := 0; y < 64; y++ {
				if ((x/8)+(y/8))%2 == 0 {
					buf[3*(x+y*64)+0] = 255
					buf[3*(x+y*64)+2] = 255
				}
			}
		}

	case "not decodeable":
		fallthrough
	default:
		for x := 0; x < 64; x++ {
			for y := 0; y < 64; y++ {
				if x <= 3 || x >= 60 || y <= 3 || y >= 60 || math.Abs(float64(y-x)) <= 3 || math.Abs(float64(y-63+x)) <= 3 {
					buf[3*(x+y*64)+0] = 255
				}
			}
		}
	}
	t := data.NewTexture2D(nil, data.Tex2DTarget2D)
	t.Bind(0)
	t.AllocBytes(buf, 0, gl.RGBA, 64, 64, gl.RGB)
	t.FilterMode(data.FilterNearest, data.FilterNearest)
	t.Unbind(0)
	return t
}

func NewCubemapEmpty(w, h, d int) *data.Texture3D {
	t := data.NewTexture3D(nil, data.Tex3DTargetCubeMap)
	t.Bind(0)
	t.AllocEmpty(0, gl.RGBA, int32(w), int32(h), int32(d), gl.RGBA)
	t.ApplyDefaults()
	t.Unbind(0)
	return t
}

/*
	rs - right (+x), left (-x), top (+y), bottom (-y), back (+z), front (-z)
*/
func NewCubemap(w, h, d int, rs [6]io.Reader) (*data.Texture3D, error) {
	t := data.NewTexture3D(nil, data.Tex3DTargetCubeMap)
	t.Bind(0)
	for i, r := range rs {
		b, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		face := data.TexTarget(int(data.Tex2DTargetCubeMapPositiveX) + i)
		t.As2D(face).AllocBytes(b, 0, gl.RGBA, int32(w), int32(h), gl.RGBA)
	}
	t.ApplyDefaults()
	t.Unbind(0)
	return t, nil
}

/*
	bufs - right (+x), left (-x), top (+y), bottom (-y), back (+z), front (-z)
*/
func NewCubemapBytes(w, h, d int, bufs [6][]byte) (*data.Texture3D, error) {
	t := data.NewTexture3D(nil, data.Tex3DTargetCubeMap)
	t.Bind(0)
	for i, b := range bufs {
		face := data.TexTarget(int(data.Tex2DTargetCubeMapPositiveX) + i)
		t.As2D(face).AllocBytes(b, 0, gl.RGBA, int32(w), int32(h), gl.RGBA)
	}
	t.ApplyDefaults()
	t.Unbind(0)
	return t, nil
}
