package core

import (
	"image"
	"io"
	"math"
	"os"

	"github.com/Qendolin/go-printpixel/core/data"
	"github.com/Qendolin/go-printpixel/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
)

func NewTexture1DEmpty(w int) *data.Texture1D {
	t := data.NewTexture1D(nil, data.Tex1DTarget1D)
	t.Bind(0)
	t.AllocEmpty(0, gl.RGBA, 2, gl.RGBA)
	t.ApplyDefaults()
	t.Unbind(0)
	return t
}

func NewTexture1DFromBytes(b []byte, w int) *data.Texture1D {
	t := data.NewTexture1D(nil, data.Tex1DTarget1D)
	t.Bind(0)
	t.AllocBytes(b, 0, gl.RGBA, 2, gl.RGBA)
	t.ApplyDefaults()
	t.Unbind(0)
	return t
}

func NewTexture1DFromFile(r io.Reader) (*data.Texture1D, error) {
	t := data.NewTexture1D(nil, data.Tex1DTarget1D)
	t.Bind(0)
	err := t.AllocFile(r, 0, gl.RGBA)
	if err != nil {
		return nil, err
	}
	t.ApplyDefaults()
	t.Unbind(0)
	return t, nil
}

func NewTexture1DFromImage(img image.Image) *data.Texture1D {
	t := data.NewTexture1D(nil, data.Tex1DTarget1D)
	t.Bind(0)
	t.AllocImage(img, 0, gl.RGBA)
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

func NewTexture2DFromBytes(b []byte, w, h int) *data.Texture2D {
	t := data.NewTexture2D(nil, data.Tex2DTarget2D)
	t.Bind(0)
	t.AllocBytes(b, 0, gl.RGBA, int32(w), int32(h), gl.RGBA)
	t.ApplyDefaults()
	t.Unbind(0)
	return t
}

func NewTexture2DFromFile(r io.Reader) (*data.Texture2D, error) {
	t := data.NewTexture2D(nil, data.Tex2DTarget2D)
	t.Bind(0)
	err := t.AllocFile(r, 0, gl.RGBA)
	if err != nil {
		return nil, err
	}
	t.ApplyDefaults()
	t.Unbind(0)
	return t, nil
}

func NewTexture2DFromImage(img image.Image) *data.Texture2D {
	t := data.NewTexture2D(nil, data.Tex2DTarget2D)
	t.Bind(0)
	t.AllocImage(img, 0, gl.RGBA)
	t.ApplyDefaults()
	t.Unbind(0)
	return t
}

func LoadTexture(path string) (t *data.Texture2D) {
	path, err := utils.ResolvePath(path)
	if err != nil {
		return newTextureError("not found")
	}
	f, err := os.Open(path)
	if err != nil {
		return newTextureError("not found")
	}
	t, err = NewTexture2DFromFile(f)
	if err != nil {
		return newTextureError("not decodeable")
	}
	return t
}

func newTextureError(cause string) *data.Texture2D {
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
	bufs - right (+x), left (-x), top (+y), bottom (-y), back (+z), front (-z)
*/
func NewCubemapFromBytes(bufs [6][]byte, w, h, d int) *data.Texture3D {
	t := data.NewTexture3D(nil, data.Tex3DTargetCubeMap)
	t.Bind(0)
	for i, b := range bufs {
		face := data.TexTarget(int(data.Tex2DTargetCubeMapPositiveX) + i)
		t.As2D(face).AllocBytes(b, 0, gl.RGBA, int32(w), int32(h), gl.RGBA)
	}
	t.ApplyDefaults()
	t.Unbind(0)
	return t
}

/*
	files - right (+x), left (-x), top (+y), bottom (-y), back (+z), front (-z)
*/
func NewCubemapFromFiles(files [6]io.Reader) (*data.Texture3D, error) {
	t := data.NewTexture3D(nil, data.Tex3DTargetCubeMap)
	t.Bind(0)
	for i, r := range files {
		face := data.TexTarget(int(data.Tex2DTargetCubeMapPositiveX) + i)
		if err := t.As2D(face).AllocFile(r, 0, gl.RGBA); err != nil {
			return nil, err
		}
	}
	t.ApplyDefaults()
	t.Unbind(0)
	return t, nil
}

/*
	images - right (+x), left (-x), top (+y), bottom (-y), back (+z), front (-z)
*/
func NewCubemapFromImages(images [6]image.Image) *data.Texture3D {
	t := data.NewTexture3D(nil, data.Tex3DTargetCubeMap)
	t.Bind(0)
	for i, img := range images {
		face := data.TexTarget(int(data.Tex2DTargetCubeMapPositiveX) + i)
		t.As2D(face).AllocImage(img, 0, gl.RGBA)
	}
	t.ApplyDefaults()
	t.Unbind(0)
	return t
}
