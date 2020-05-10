package data

import (
	"image"
	"image/draw"
	"io"

	"github.com/Qendolin/go-printpixel/internal/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
)

//Texture Wrap Modes
const (
	WrapClampToEdge       = gl.CLAMP_TO_EDGE
	WrapClampToBorder     = gl.CLAMP_TO_BORDER
	WrapMirroredRepeat    = gl.MIRRORED_REPEAT
	WrapRepeat            = gl.REPEAT
	WrapMirrorClampToEdge = gl.MIRROR_CLAMP_TO_EDGE
)

//Common texture filters
const (
	FilterNearest = gl.NEAREST
	FilterLinear  = gl.LINEAR
)

//Texture minification filters
const (
	FilterNearestMipMapNearest = gl.NEAREST_MIPMAP_NEAREST
	FilterLinearMipMapNearest  = gl.LINEAR_MIPMAP_NEAREST
	FilterNearestMipMapLinear  = gl.NEAREST_MIPMAP_LINEAR
	FilterLinearMipMapLinear   = gl.LINEAR_MIPMAP_LINEAR
)

//1D Texture targets
const (
	Texture1D      = gl.TEXTURE_1D
	TextureProxy1D = gl.PROXY_TEXTURE_1D
)

//2D Texture targets
const (
	Texture2D               = gl.TEXTURE_2D
	TextureProxy2D          = gl.PROXY_TEXTURE_2D
	Texture1DArray          = gl.TEXTURE_1D_ARRAY
	TextureProxy1DArray     = gl.PROXY_TEXTURE_1D_ARRAY
	TextureRectangle        = gl.TEXTURE_RECTANGLE
	TextureProxyRectangle   = gl.PROXY_TEXTURE_RECTANGLE
	TextureCubeMapPositiveX = gl.TEXTURE_CUBE_MAP_POSITIVE_X
	TextureCubeMapPositiveY = gl.TEXTURE_CUBE_MAP_POSITIVE_Y
	TextureCubeMapPositiveZ = gl.TEXTURE_CUBE_MAP_POSITIVE_Z
	TextureCubeMapNegativeX = gl.TEXTURE_CUBE_MAP_NEGATIVE_X
	TextureCubeMapNegativeY = gl.TEXTURE_CUBE_MAP_NEGATIVE_Y
	TextureCubeMapNegativeZ = gl.TEXTURE_CUBE_MAP_NEGATIVE_Z
	TextureProxyCubeMap     = gl.PROXY_TEXTURE_CUBE_MAP
)

//3D Texture targets
const (
	Texture3D           = gl.TEXTURE_3D
	TextureProxy3D      = gl.PROXY_TEXTURE_3D
	Texture2DArray      = gl.TEXTURE_2D_ARRAY
	TextureProxy2DArray = gl.PROXY_TEXTURE_2D_ARRAY
)

type Texture struct {
	*uint32
	Type uint32
}

func NewTexture(texType uint32) *Texture {
	id := new(uint32)
	gl.GenTextures(1, id)
	return &Texture{uint32: id, Type: texType}
}

func (tex *Texture) Id() uint32 {
	return *tex.uint32
}

func (tex *Texture) BindAs(target, unit uint32) {
	gl.ActiveTexture(unit)
	gl.BindTexture(target, tex.Id())
}

func (tex *Texture) Bind(unit uint32) {
	tex.BindAs(tex.Type, unit)
}

func (tex *Texture) UnbindAs(target, unit uint32) {
	gl.ActiveTexture(unit)
	gl.BindTexture(target, 0)
}

func (tex *Texture) Unbind(unit uint32) {
	tex.UnbindAs(tex.Type, unit)
}

func (tex *Texture) BindForAs(target, unit uint32, context utils.BindingClosure) {
	tex.BindAs(target, unit)
	defered := context()
	tex.UnbindAs(target, unit)
	for _, deferedFunc := range defered {
		deferedFunc()
	}
}

func (tex *Texture) BindFor(unit uint32, context utils.BindingClosure) {
	tex.BindForAs(tex.Type, unit, context)
}

func (tex *Texture) WrapModeAs(target uint32, sMode, tMode, rMode int32) {
	if sMode != 0 {
		gl.TexParameteri(target, gl.TEXTURE_WRAP_S, sMode)
	}
	if tMode != 0 {
		gl.TexParameteri(target, gl.TEXTURE_WRAP_T, tMode)
	}
	if tMode != 0 {
		gl.TexParameteri(target, gl.TEXTURE_WRAP_R, rMode)
	}
}

func (tex *Texture) WrapMode(sMode, tMode, rMode int32) {
	tex.WrapModeAs(tex.Type, sMode, tMode, rMode)
}

func (tex *Texture) FilterModeAs(target uint32, minMode, magMode int32) {
	if minMode != 0 {
		gl.TexParameteri(target, gl.TEXTURE_MIN_FILTER, minMode)
	}
	if magMode != 0 {
		gl.TexParameteri(target, gl.TEXTURE_MAG_FILTER, magMode)
	}
}

func (tex *Texture) FilterMode(minMode, magMode int32) {
	tex.FilterModeAs(tex.Type, minMode, magMode)
}

func (tex *Texture) WriteAs(target uint32, level, internalFormat, width, height, depth int32, format, dataType uint32, data interface{}) {
	dataPtr := gl.Ptr(data)
	if target == gl.TEXTURE_1D || target == gl.PROXY_TEXTURE_1D {
		gl.TexImage1D(target, level, internalFormat, width, 0, format, dataType, dataPtr)
	} else if target == gl.TEXTURE_3D || target == gl.PROXY_TEXTURE_3D || target == gl.TEXTURE_2D_ARRAY || target == gl.PROXY_TEXTURE_2D_ARRAY {
		gl.TexImage3D(target, level, internalFormat, width, height, depth, 0, format, dataType, dataPtr)
	} else {
		gl.TexImage2D(target, level, internalFormat, width, height, 0, format, dataType, dataPtr)
	}
}

func (tex *Texture) Write(level, internalFormat, width, height, depth int32, format, dataType uint32, data interface{}) {
	tex.WriteAs(tex.Type, level, internalFormat, width, height, depth, format, dataType, data)
}

func (tex *Texture) WriteFromFile2DAs(file io.Reader, target uint32, level, internalFormat int32, format, dataType uint32) error {
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	size := img.Bounds().Size()
	tex.WriteAs(target, level, internalFormat, int32(size.X), int32(size.Y), 0, format, dataType, rgba.Pix)
	return nil
}

func (tex *Texture) WriteFromFile2D(file io.Reader, level, internalFormat int32, format, dataType uint32) error {
	return tex.WriteFromFile2DAs(file, tex.Type, level, internalFormat, format, dataType)
}

func (tex *Texture) WriteFromFile1DAs(file io.Reader, target uint32, level, internalFormat int32, format, dataType uint32) error {
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	size := img.Bounds().Size()
	tex.WriteAs(target, level, internalFormat, int32(size.X), 0, 0, format, dataType, rgba.Pix)
	return nil
}

func (tex *Texture) WriteFromFile1D(file io.Reader, level, internalFormat int32, format, dataType uint32) error {
	return tex.WriteFromFile1DAs(file, tex.Type, level, internalFormat, format, dataType)
}

/*
	files - right (+x), left (-x), top (+y), bottom (-y), back (+z), front (-z)
*/
func (tex *Texture) WriteFromFile3D(files [6]io.Reader, level, internalFormat int32, format, dataType uint32) error {
	for i, file := range files {
		err := tex.WriteFromFile2DAs(file, TextureCubeMapPositiveX+uint32(i), level, internalFormat, format, dataType)
		if err != nil {
			return err
		}
	}
	return nil
}
