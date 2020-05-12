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
	Target TexTarget
}

func NewTexture(texType TexTarget) *Texture {
	id := new(uint32)
	gl.GenTextures(1, id)
	return &Texture{uint32: id, Target: texType}
}

func (tex *Texture) Id() uint32 {
	return *tex.uint32
}

func (tex *Texture) As(target TexTarget) *Texture {
	return &Texture{
		uint32: tex.uint32,
		Target: target,
	}
}

func (tex *Texture) Bind(unit int) {
	gl.ActiveTexture(uint32(gl.TEXTURE0 + unit))
	gl.BindTexture(uint32(tex.Target), tex.Id())
}

func (tex *Texture) Unbind(unit int) {
	gl.ActiveTexture(uint32(gl.TEXTURE0 + unit))
	gl.BindTexture(uint32(tex.Target), 0)
}

func (tex *Texture) BindFor(unit int, context utils.BindingClosure) {
	tex.Bind(unit)
	defered := context()
	tex.Unbind(unit)
	for _, deferedFunc := range defered {
		deferedFunc()
	}
}

func (tex *Texture) WrapMode(sMode, tMode, rMode TexWrapMode) {
	if sMode != 0 {
		gl.TexParameteri(uint32(tex.Target), gl.TEXTURE_WRAP_S, int32(sMode))
	}
	if tMode != 0 {
		gl.TexParameteri(uint32(tex.Target), gl.TEXTURE_WRAP_T, int32(tMode))
	}
	if rMode != 0 {
		gl.TexParameteri(uint32(tex.Target), gl.TEXTURE_WRAP_R, int32(rMode))
	}
}

func (tex *Texture) FilterMode(minMode, magMode TexFilterMode) {
	if minMode != 0 {
		gl.TexParameteri(uint32(tex.Target), gl.TEXTURE_MIN_FILTER, int32(minMode))
	}
	if magMode != 0 {
		gl.TexParameteri(uint32(tex.Target), gl.TEXTURE_MAG_FILTER, int32(magMode))
	}
}

func (tex *Texture) GenerateMipmap() {
	gl.GenerateMipmap(uint32(tex.Target))
}

func (tex *Texture) Alloc(level, internalFormat, width, height, depth int32, format, dataType uint32, data interface{}) {
	dataPtr := gl.Ptr(data)
	if tex.Target == Texture1D || tex.Target == TextureProxy1D {
		gl.TexImage1D(uint32(tex.Target), level, internalFormat, width, 0, format, dataType, dataPtr)
	} else if tex.Target == Texture3D || tex.Target == TextureProxy3D || tex.Target == Texture2DArray || tex.Target == TextureProxy2DArray {
		gl.TexImage3D(uint32(tex.Target), level, internalFormat, width, height, depth, 0, format, dataType, dataPtr)
	} else {
		gl.TexImage2D(uint32(tex.Target), level, internalFormat, width, height, 0, format, dataType, dataPtr)
	}
}

func (tex *Texture) AllocWithFile2D(file io.Reader, level, internalFormat int32, format, dataType uint32) error {
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	size := img.Bounds().Size()
	tex.Alloc(level, internalFormat, int32(size.X), int32(size.Y), 0, format, dataType, rgba.Pix)
	return nil
}

func (tex *Texture) AllocWithFile1D(file io.Reader, level, internalFormat int32, format, dataType uint32) error {
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	size := img.Bounds().Size()
	tex.Alloc(level, internalFormat, int32(size.X), 0, 0, format, dataType, rgba.Pix)
	return nil
}

/*
	files - right (+x), left (-x), top (+y), bottom (-y), back (+z), front (-z)
*/
func (tex *Texture) AllocWithFile3D(files [6]io.Reader, level, internalFormat int32, format, dataType uint32) error {
	for i, file := range files {
		err := tex.As(TexTarget(int(TextureCubeMapPositiveX)+i)).AllocWithFile2D(file, level, internalFormat, format, dataType)
		if err != nil {
			return err
		}
	}
	return nil
}

func (tex *Texture) AllocWithImage(img image.Image, level, internalFormat int32, format, dataType uint32) {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	size := img.Bounds().Size()
	tex.Alloc(level, internalFormat, int32(size.X), int32(size.Y), 0, format, dataType, rgba.Pix)
}

func (tex *Texture) AllocWithBytes(bytes []byte, width, height int32, level, internalFormat int32, format uint32) {
	tex.Alloc(level, internalFormat, width, height, 0, format, gl.BYTE, bytes)
}
