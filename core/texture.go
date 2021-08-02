package core

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"io"
	"math"
	"math/rand"
	"os"

	"github.com/Qendolin/go-printpixel/core/data"
	"github.com/Qendolin/go-printpixel/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
)

func init() {
	// magenta and black checkers
	pix := make([]byte, 64*64*4)
	for x := 0; x < 64; x++ {
		for y := 0; y < 64; y++ {
			if ((x/8)+(y/8))%2 == 0 {
				pix[4*(x+y*64)+0] = 255
				pix[4*(x+y*64)+2] = 255
				pix[4*(x+y*64)+3] = 255
			}
		}
	}
	errorImageData[ErrNotFound] = pix

	// random colors
	pix = make([]byte, 64*64*4)
	for i := range pix {
		pix[i] = byte(rand.Float32() * 256)
	}
	errorImageData[ErrDecode] = pix

	// red square with an X
	pix = make([]byte, 64*64*4)
	for x := 0; x < 64; x++ {
		for y := 0; y < 64; y++ {
			if x <= 3 || x >= 60 || y <= 3 || y >= 60 || math.Abs(float64(y-x)) <= 3 || math.Abs(float64(y-63+x)) <= 3 {
				pix[4*(x+y*64)+0] = 255
				pix[4*(x+y*64)+3] = 255
			}
		}
	}
	errorImageData[nil] = pix
}

type ImageError error

var (
	ErrNotFound ImageError = fmt.Errorf("the image was not found or could not be opened")
	ErrDecode              = fmt.Errorf("the image could not be decoded")
)

var errorImageData = map[ImageError][]byte{}

func newErrorTexture(err error) *data.Texture2D {
	pix, ok := errorImageData[err]
	if !ok {
		pix = errorImageData[nil]
	}
	tex := data.NewTexture2D(nil, data.Tex2DTarget2D)
	tex.Bind(0)
	err = tex.AllocBytes(pix, 0, gl.RGBA8, 64, 64, gl.RGBA)
	if err != nil {
		panic(err)
	}
	tex.FilterMode(data.FilterLinear, data.FilterNearest)
	tex.WrapMode(data.WrapClampToEdge, data.WrapClampToEdge)
	tex.Unbind(0)
	return tex
}

func InitImages(layers int, img0 image.Image, imgs ...image.Image) TextureInitializer {
	imgs = append([]image.Image{img0}, imgs...)
	tini := TextureInitializer{
		Format: data.RGBA,
	}

	if layers > 0 {
		tini.Levels = make([]interface{}, len(imgs)/layers)
		tini.Target = data.Tex3DTarget2DArray
		tini.Depth = layers
	} else {
		tini.Levels = make([]interface{}, len(imgs))
		tini.Target = data.Tex2DTarget2D
	}

	if layers == 0 {
		layers = 1
	}

	var lvl []byte
	var size int
	for i, img := range imgs {
		var buf []byte
		var w, h int
		switch ti := img.(type) {
		case *image.Uniform:
			r, g, b, a := ti.RGBA()
			buf = []byte{byte(r >> 8), byte(g >> 8), byte(b >> 8), byte(a >> 8)}
			w, h = 1, 1
		case *image.RGBA:
			buf = ti.Pix
			w, h = ti.Rect.Size().X, ti.Rect.Size().Y
		default:
			rgba := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
			draw.Draw(rgba, rgba.Rect, img, image.Point{}, draw.Src)
			buf = rgba.Pix
			w, h = rgba.Rect.Size().X, rgba.Rect.Size().Y
		}

		if i == 0 {
			tini.Width, tini.Height = w, h
		}

		l, j := i/layers, i%layers
		if l >= len(tini.Levels) {
			break
		}
		if j == 0 {
			size = len(buf)
			lvl = make([]byte, size*layers)
			tini.Levels[l] = lvl
		}

		copy(lvl[j*size:(j+1)*size], buf)
	}

	return tini
}

func InitFiles(layers int, file0 io.Reader, files ...io.Reader) TextureInitializer {
	files = append([]io.Reader{file0}, files...)
	imgs := make([]image.Image, len(files))
	for i, file := range files {
		img, _, err := image.Decode(file)
		if err != nil {
			return TextureInitializer{
				Levels: []interface{}{ErrDecode},
			}
		}
		imgs[i] = img
	}
	return InitImages(layers, imgs[0], imgs[1:]...)
}

func InitPaths(layers int, path0 string, paths ...string) TextureInitializer {
	paths = append([]string{path0}, paths...)
	files := make([]io.Reader, len(paths))
	for i, path := range paths {
		path, err := utils.ResolvePath(path)
		if err != nil {
			return TextureInitializer{
				Levels: []interface{}{ErrNotFound},
			}
		}
		file, err := os.Open(path)
		if err != nil {
			return TextureInitializer{
				Levels: []interface{}{ErrNotFound},
			}
		}
		defer file.Close()
		files[i] = file
	}
	return InitFiles(layers, files[0], files[1:]...)
}

func InitBytes(w, h, d int, format data.ColorFormat, layers int, buf0 []byte, bufs ...[]byte) TextureInitializer {
	bufs = append([][]byte{buf0}, bufs...)
	if layers == 0 {
		layers = 1
	}
	lvls := make([]interface{}, len(bufs)/layers)

	for l := 0; l < len(lvls); l++ {
		if l*layers >= len(bufs) {
			break
		}

		if bufs[l*layers] == nil {
			continue
		}

		size := len(bufs[l*layers])
		lvl := make([]byte, size*layers)
		lvls[l] = lvl

		for i := 0; i < layers; i++ {
			if l*layers+i >= len(bufs) {
				break
			}
			copy(lvl[i*size:(i+1)*size], bufs[l*layers+i])
		}
	}

	var typ data.TexTarget
	if d != 0 {
		typ = data.Tex3DTarget3D
	} else if h != 0 {
		typ = data.Tex2DTarget2D
	} else {
		typ = data.Tex1DTarget1D
	}
	return TextureInitializer{
		Width: w, Height: h, Depth: d,
		Format: format,
		Levels: lvls,
		Target: typ,
	}
}

func InitEmpty(w, h, d int) TextureInitializer {
	return InitBytes(w, h, d, gl.RGBA, 0, nil)
}

type TextureInitializer struct {
	Width, Height, Depth int
	Format               data.ColorFormat
	Target               data.TexTarget
	Levels               []interface{}
	GenerateMipMap       bool
	MagFilter            data.TexFilterMode
	MinFilter            data.TexFilterMode
}

func (tini TextureInitializer) As(target data.TexTarget) TextureInitializer {
	tini.Target = target
	return tini
}

func (tini TextureInitializer) WithFilters(minFilter, magFilter data.TexFilterMode) TextureInitializer {
	tini.MinFilter = minFilter
	tini.MagFilter = magFilter
	return tini
}

func (tini TextureInitializer) WithGeneratedMipMap() TextureInitializer {
	tini.GenerateMipMap = true
	return tini
}

func (tini TextureInitializer) WithLevels(n int) TextureInitializer {
	dst := make([]interface{}, n)
	copy(dst, tini.Levels)
	tini.Levels = dst
	return tini
}

func (tini TextureInitializer) WithRequiredLevels() TextureInitializer {
	max := math.Max(float64(tini.Width), math.Max(float64(tini.Height), float64(tini.Depth)))
	return tini.WithLevels(int(math.Ceil(math.Log2(max))))
}

func (tini *TextureInitializer) Level(index int) (interface{}, error) {
	if !tini.HasLevel(index) {
		return nil, nil
	}
	lvl := tini.Levels[index]
	if err, ok := lvl.(error); ok {
		return nil, err
	}
	return lvl, nil
}

func (tini *TextureInitializer) HasLevel(index int) bool {
	return index < len(tini.Levels)
}

func calcLevelSize(bw, bh, bd int, typ data.TexTarget, lvl int) (w, h, d int) {
	w, h, d = 2*bw/(2<<lvl), 2*bh/(2<<lvl), 2*bd/(2<<lvl)
	if w == 0 {
		w = 1
	}
	if h == 0 {
		h = 1
	}
	if d == 0 {
		d = 1
	}
	if typ.IsArray() {
		switch typ.Dimensions() {
		case 2:
			h = bh
		case 3:
			d = bd
		}
	}
	return
}

func NewTexture(tini TextureInitializer, format data.ColorFormat) (*data.GLTexture, error) {
	t := data.NewTexture(nil, tini.Target)
	t.Bind(0)

	lvl := 0
	for ; tini.HasLevel(lvl) || lvl == 0; lvl++ {
		w, h, d := calcLevelSize(tini.Width, tini.Height, tini.Depth, tini.Target, lvl)
		px, err := tini.Level(lvl)
		if err != nil {
			return nil, err
		}
		err = t.Alloc(int32(lvl), format.InternalFormatEnum(), int32(w), int32(h), int32(d), tini.Format.PixelFormatEnum(), px)
		if err != nil {
			if ne := new(data.NotEnoughError); errors.As(err, ne) {
				return nil, fmt.Errorf("level %d data has %d bytes but expected %d (%dx%dx%d)", lvl, ne.ActualSize, ne.RequiredSize, w, h, d)
			}
			return nil, err
		}
	}

	t.WrapMode(data.WrapDefault, data.WrapDefault, data.WrapDefault)
	if tini.MinFilter == 0 {
		tini.MinFilter = data.FilterMinDefault
	}
	if tini.MagFilter == 0 {
		tini.MagFilter = data.FilterMagDefault
	}
	t.FilterMode(tini.MinFilter, tini.MagFilter)
	t.MipMapLevels(0, lvl-1)
	if tini.GenerateMipMap {
		t.GenerateMipmap()
	}
	t.Unbind(0)
	return t, nil
}

func NewTexture1D(tini TextureInitializer, format data.ColorFormat) (*data.Texture1D, error) {
	if tini.Target.Dimensions() != 1 {
		tini.Target = data.Tex1DTarget1D
	}

	t, err := NewTexture(tini, format)
	if err != nil {
		return nil, err
	}
	return t.As1D(tini.Target), nil
}

func NewTexture2D(tini TextureInitializer, format data.ColorFormat) (*data.Texture2D, error) {
	if tini.Target.Dimensions() != 2 {
		tini.Target = data.Tex2DTarget2D
	}

	t, err := NewTexture(tini, format)
	if err != nil {
		return nil, err
	}
	return t.As2D(tini.Target), nil
}

func MustNewTexture2D(tini TextureInitializer, format data.ColorFormat) *data.Texture2D {
	t, err := NewTexture2D(tini, format)
	if err != nil {
		return newErrorTexture(err)
	}
	return t
}

func NewTexture3D(tini TextureInitializer, format data.ColorFormat) (*data.Texture3D, error) {
	if tini.Target.Dimensions() != 3 {
		tini.Target = data.Tex3DTarget3D
	}

	t, err := NewTexture(tini, format)
	if err != nil {
		return nil, err
	}
	return t.As3D(tini.Target), nil
}

// func NewCubemapEmpty(w, h, d int) *data.Texture3D {
// 	t := data.NewTexture3D(nil, data.Tex3DTargetCubeMap)
// 	t.Bind(0)
// 	t.AllocEmpty(0, gl.RGBA, int32(w), int32(h), int32(d), gl.RGBA)
// 	t.ApplyDefaults()
// 	t.Unbind(0)
// 	return t
// }

/*
	bufs - right (+x), left (-x), top (+y), bottom (-y), back (+z), front (-z)
*/
// func NewCubemap(w, h, d int, bufs [6][]byte) *data.Texture3D {
// 	t := data.NewTexture3D(nil, data.Tex3DTargetCubeMap)
// 	t.Bind(0)
// 	for i, b := range bufs {
// 		face := data.TexTarget(int(data.Tex2DTargetCubeMapPositiveX) + i)
// 		t.As2D(face).AllocBytes(b, 0, gl.RGBA, int32(w), int32(h), gl.RGBA)
// 	}
// 	t.ApplyDefaults()
// 	t.Unbind(0)
// 	return t
// }

// func NewVolumeEmpty(w, h, d int) *data.Texture3D {
// 	t := data.NewTexture3D(nil, data.Tex3DTarget3D)
// 	t.Bind(0)
// 	t.AllocEmpty(0, gl.RGBA, int32(w), int32(h), int32(d), gl.RGBA)
// 	t.ApplyDefaults()
// 	t.Unbind(0)
// 	return t
// }

// func NewVolume(w, h, d int, b []byte) *data.Texture3D {
// 	t := data.NewTexture3D(nil, data.Tex3DTarget3D)
// 	t.Bind(0)
// 	t.AllocBytes(b, 0, gl.RGBA, int32(w), int32(h), int32(d), gl.RGBA)
// 	t.ApplyDefaults()
// 	t.Unbind(0)
// 	return t
// }
