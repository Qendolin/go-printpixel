package data

import (
	"image"
	"image/draw"
	"io"

	"github.com/Qendolin/go-printpixel/internal/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type TexWrapMode int

//Texture Wrap Modes
const (
	WrapClampToEdge       = TexWrapMode(gl.CLAMP_TO_EDGE)
	WrapClampToBorder     = TexWrapMode(gl.CLAMP_TO_BORDER)
	WrapMirroredRepeat    = TexWrapMode(gl.MIRRORED_REPEAT)
	WrapRepeat            = TexWrapMode(gl.REPEAT)
	WrapMirrorClampToEdge = TexWrapMode(gl.MIRROR_CLAMP_TO_EDGE)
)

type TexFilterMode int

//Common texture filters
const (
	FilterNearest = TexFilterMode(gl.NEAREST)
	FilterLinear  = TexFilterMode(gl.LINEAR)
)

//Texture minification filters
const (
	FilterNearestMipMapNearest = TexFilterMode(gl.NEAREST_MIPMAP_NEAREST)
	FilterLinearMipMapNearest  = TexFilterMode(gl.LINEAR_MIPMAP_NEAREST)
	FilterNearestMipMapLinear  = TexFilterMode(gl.NEAREST_MIPMAP_LINEAR)
	FilterLinearMipMapLinear   = TexFilterMode(gl.LINEAR_MIPMAP_LINEAR)
)

type TexTarget int

//1D Texture targets
const (
	Texture1D      = TexTarget(gl.TEXTURE_1D)
	TextureProxy1D = TexTarget(gl.PROXY_TEXTURE_1D)
)

//2D Texture targets
const (
	Texture2D               = TexTarget(gl.TEXTURE_2D)
	TextureProxy2D          = TexTarget(gl.PROXY_TEXTURE_2D)
	Texture1DArray          = TexTarget(gl.TEXTURE_1D_ARRAY)
	TextureProxy1DArray     = TexTarget(gl.PROXY_TEXTURE_1D_ARRAY)
	TextureRectangle        = TexTarget(gl.TEXTURE_RECTANGLE)
	TextureProxyRectangle   = TexTarget(gl.PROXY_TEXTURE_RECTANGLE)
	TextureCubeMapPositiveX = TexTarget(gl.TEXTURE_CUBE_MAP_POSITIVE_X)
	TextureCubeMapPositiveY = TexTarget(gl.TEXTURE_CUBE_MAP_POSITIVE_Y)
	TextureCubeMapPositiveZ = TexTarget(gl.TEXTURE_CUBE_MAP_POSITIVE_Z)
	TextureCubeMapNegativeX = TexTarget(gl.TEXTURE_CUBE_MAP_NEGATIVE_X)
	TextureCubeMapNegativeY = TexTarget(gl.TEXTURE_CUBE_MAP_NEGATIVE_Y)
	TextureCubeMapNegativeZ = TexTarget(gl.TEXTURE_CUBE_MAP_NEGATIVE_Z)
	TextureProxyCubeMap     = TexTarget(gl.PROXY_TEXTURE_CUBE_MAP)
)

//3D Texture targets
const (
	Texture3D           = TexTarget(gl.TEXTURE_3D)
	TextureProxy3D      = TexTarget(gl.PROXY_TEXTURE_3D)
	Texture2DArray      = TexTarget(gl.TEXTURE_2D_ARRAY)
	TextureProxy2DArray = TexTarget(gl.PROXY_TEXTURE_2D_ARRAY)
)

type Texture struct {
	*uint32
	Type TexTarget
}

func NewTexture(texType TexTarget) *Texture {
	id := new(uint32)
	gl.GenTextures(1, id)
	return &Texture{uint32: id, Type: texType}
}

func (tex *Texture) Id() uint32 {
	return *tex.uint32
}

func (tex *Texture) BindAs(target TexTarget, unit uint32) {
	gl.ActiveTexture(unit)
	gl.BindTexture(uint32(target), tex.Id())
}

func (tex *Texture) Bind(unit uint32) {
	tex.BindAs(tex.Type, unit)
}

func (tex *Texture) UnbindAs(target TexTarget, unit uint32) {
	gl.ActiveTexture(unit)
	gl.BindTexture(uint32(target), 0)
}

func (tex *Texture) Unbind(unit uint32) {
	tex.UnbindAs(tex.Type, unit)
}

func (tex *Texture) BindForAs(target TexTarget, unit uint32, context utils.BindingClosure) {
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

func (tex *Texture) WrapModeAs(target TexTarget, sMode, tMode, rMode TexWrapMode) {
	if sMode != 0 {
		gl.TexParameteri(uint32(target), gl.TEXTURE_WRAP_S, int32(sMode))
	}
	if tMode != 0 {
		gl.TexParameteri(uint32(target), gl.TEXTURE_WRAP_T, int32(tMode))
	}
	if tMode != 0 {
		gl.TexParameteri(uint32(target), gl.TEXTURE_WRAP_R, int32(rMode))
	}
}

func (tex *Texture) WrapMode(sMode, tMode, rMode TexWrapMode) {
	tex.WrapModeAs(tex.Type, sMode, tMode, rMode)
}

func (tex *Texture) FilterModeAs(target TexTarget, minMode, magMode TexFilterMode) {
	if minMode != 0 {
		gl.TexParameteri(uint32(target), gl.TEXTURE_MIN_FILTER, int32(minMode))
	}
	if magMode != 0 {
		gl.TexParameteri(uint32(target), gl.TEXTURE_MAG_FILTER, int32(magMode))
	}
}

func (tex *Texture) FilterMode(minMode, magMode TexFilterMode) {
	tex.FilterModeAs(tex.Type, minMode, magMode)
}

func (tex *Texture) WriteAs(target TexTarget, level, internalFormat, width, height, depth int32, format, dataType uint32, data interface{}) {
	dataPtr := gl.Ptr(data)
	if target == Texture1D || target == TextureProxy1D {
		gl.TexImage1D(uint32(target), level, internalFormat, width, 0, format, dataType, dataPtr)
	} else if target == Texture3D || target == TextureProxy3D || target == Texture2DArray || target == TextureProxy2DArray {
		gl.TexImage3D(uint32(target), level, internalFormat, width, height, depth, 0, format, dataType, dataPtr)
	} else {
		gl.TexImage2D(uint32(target), level, internalFormat, width, height, 0, format, dataType, dataPtr)
	}
}

func (tex *Texture) Write(level, internalFormat, width, height, depth int32, format, dataType uint32, data interface{}) {
	tex.WriteAs(tex.Type, level, internalFormat, width, height, depth, format, dataType, data)
}

func (tex *Texture) WriteFromFile2DAs(file io.Reader, target TexTarget, level, internalFormat int32, format, dataType uint32) error {
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

func (tex *Texture) WriteFromFile1DAs(file io.Reader, target TexTarget, level, internalFormat int32, format, dataType uint32) error {
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
		err := tex.WriteFromFile2DAs(file, TexTarget(int(TextureCubeMapPositiveX)+i), level, internalFormat, format, dataType)
		if err != nil {
			return err
		}
	}
	return nil
}

func (tex *Texture) WriteFromImageAs(img image.Image, target TexTarget, level, internalFormat int32, format, dataType uint32) {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	size := img.Bounds().Size()
	tex.WriteAs(target, level, internalFormat, int32(size.X), int32(size.Y), 0, format, dataType, rgba.Pix)
}

func (tex *Texture) WriteFromImage(img image.Image, level, internalFormat int32, format, dataType uint32) {
	tex.WriteFromImageAs(img, tex.Type, level, internalFormat, format, dataType)
}

func (tex *Texture) WriteFromBytes(bytes []byte, width, height int32, level, internalFormat int32, format uint32) {
	tex.WriteFromBytesAs(bytes, width, height, tex.Type, level, internalFormat, format)
}

func (tex *Texture) WriteFromBytesAs(bytes []byte, width, height int32, target TexTarget, level, internalFormat int32, format uint32) {
	tex.WriteAs(target, level, internalFormat, width, height, 0, format, gl.BYTE, bytes)
}
