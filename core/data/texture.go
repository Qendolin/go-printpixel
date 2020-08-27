package data

import (
	"encoding/binary"
	"fmt"
	"unsafe"

	"github.com/Qendolin/go-printpixel/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
)

var MipMapDefaultBaseLevel = 0
var MipMapDefaultMaxLevel = 0

type TexWrapMode int

//Texture Wrap Modes
const (
	WrapDefault           = WrapClampToEdge
	WrapClampToEdge       = TexWrapMode(gl.CLAMP_TO_EDGE)
	WrapClampToBorder     = TexWrapMode(gl.CLAMP_TO_BORDER)
	WrapMirroredRepeat    = TexWrapMode(gl.MIRRORED_REPEAT)
	WrapRepeat            = TexWrapMode(gl.REPEAT)
	WrapMirrorClampToEdge = TexWrapMode(gl.MIRROR_CLAMP_TO_EDGE)
)

type TexFilterMode int

var FilterMagDefault = FilterLinear

//Texture magnification filters
const (
	FilterNearest = TexFilterMode(gl.NEAREST)
	FilterLinear  = TexFilterMode(gl.LINEAR)
)

var FilterMinDefault = FilterLinearMipMapLinear

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
	Tex1DTarget1D      = TexTarget(gl.TEXTURE_1D)
	Tex1DTargetProxy1D = TexTarget(gl.PROXY_TEXTURE_1D)
	Tex1DTargetBuffer  = TexTarget(gl.TEXTURE_BUFFER)
)

//2D Texture targets
const (
	Tex2DTarget2D               = TexTarget(gl.TEXTURE_2D)
	Tex2DTargetProxy2D          = TexTarget(gl.PROXY_TEXTURE_2D)
	Tex2DTarget1DArray          = TexTarget(gl.TEXTURE_1D_ARRAY)
	Tex2DTargetProxy1DArray     = TexTarget(gl.PROXY_TEXTURE_1D_ARRAY)
	Tex2DTargetRectangle        = TexTarget(gl.TEXTURE_RECTANGLE)
	Tex2DTargetProxyRectangle   = TexTarget(gl.PROXY_TEXTURE_RECTANGLE)
	Tex2DTargetCubeMapPositiveX = TexTarget(gl.TEXTURE_CUBE_MAP_POSITIVE_X)
	Tex2DTargetCubeMapPositiveY = TexTarget(gl.TEXTURE_CUBE_MAP_POSITIVE_Y)
	Tex2DTargetCubeMapPositiveZ = TexTarget(gl.TEXTURE_CUBE_MAP_POSITIVE_Z)
	Tex2DTargetCubeMapNegativeX = TexTarget(gl.TEXTURE_CUBE_MAP_NEGATIVE_X)
	Tex2DTargetCubeMapNegativeY = TexTarget(gl.TEXTURE_CUBE_MAP_NEGATIVE_Y)
	Tex2DTargetCubeMapNegativeZ = TexTarget(gl.TEXTURE_CUBE_MAP_NEGATIVE_Z)
	Tex2DTargetProxyCubeMap     = TexTarget(gl.PROXY_TEXTURE_CUBE_MAP)
)

//3D Texture targets
const (
	Tex3DTarget3D           = TexTarget(gl.TEXTURE_3D)
	Tex3DTargetCubeMap      = TexTarget(gl.TEXTURE_CUBE_MAP)
	Tex3DTargetProxy3D      = TexTarget(gl.PROXY_TEXTURE_3D)
	Tex3DTarget2DArray      = TexTarget(gl.TEXTURE_2D_ARRAY)
	Tex3DTargetProxy2DArray = TexTarget(gl.PROXY_TEXTURE_2D_ARRAY)
)

func (tt TexTarget) IsArray() bool {
	switch tt {
	case Tex2DTarget1DArray, Tex2DTargetProxy1DArray, Tex3DTarget2DArray, Tex3DTargetProxy2DArray, Tex3DTargetCubeMap:
		return true
	default:
		return false
	}
}

func (tt TexTarget) Dimensions() int {
	switch tt {
	case Tex1DTarget1D, Tex1DTargetProxy1D, Tex1DTargetBuffer:
		return 1
	case Tex3DTarget3D, Tex3DTargetProxy3D, Tex3DTarget2DArray, Tex3DTargetProxy2DArray, Tex3DTargetCubeMap:
		return 3
	case Tex2DTarget1DArray, Tex2DTarget2D, Tex2DTargetCubeMapNegativeX, Tex2DTargetCubeMapNegativeY, Tex2DTargetCubeMapNegativeZ,
		Tex2DTargetCubeMapPositiveX, Tex2DTargetCubeMapPositiveY, Tex2DTargetCubeMapPositiveZ, Tex2DTargetProxy1DArray, Tex2DTargetProxy2D, Tex2DTargetProxyCubeMap,
		Tex2DTargetProxyRectangle, Tex2DTargetRectangle:
		return 2
	default:
		return 0
	}
}

type GLTexture struct {
	*uint32
	Target TexTarget
}

func NewTexture(id *uint32, texType TexTarget) *GLTexture {
	return &GLTexture{uint32: id, Target: texType}
}

func NewTexture1D(id *uint32, texType TexTarget) *Texture1D {
	return &Texture1D{
		GLTexture{uint32: id, Target: texType},
	}
}

func NewTexture2D(id *uint32, texType TexTarget) *Texture2D {
	return &Texture2D{
		GLTexture{uint32: id, Target: texType},
	}
}

func NewTexture3D(id *uint32, texType TexTarget) *Texture3D {
	return &Texture3D{
		GLTexture{uint32: id, Target: texType},
	}
}

func (tex *GLTexture) Id() *uint32 {
	if tex.uint32 == nil {
		tex.uint32 = new(uint32)
		gl.GenTextures(1, tex.uint32)
	}
	return tex.uint32
}

func (tex *GLTexture) Destroy() {
	gl.DeleteTextures(1, tex.uint32)
	*tex.uint32 = 0
}

func (tex *GLTexture) As(target TexTarget) *GLTexture {
	return &GLTexture{
		uint32: tex.uint32,
		Target: target,
	}
}

func (tex *GLTexture) As1D(target TexTarget) *Texture1D {
	if target == 0 {
		if tex.Target.Dimensions() == 1 {
			target = tex.Target
		} else {
			target = Tex1DTarget1D
		}
	}
	return &Texture1D{
		GLTexture: *tex.As(TexTarget(target)),
	}
}

func (tex *GLTexture) As2D(target TexTarget) *Texture2D {
	if target == 0 {
		if tex.Target.Dimensions() == 2 {
			target = tex.Target
		} else {
			target = Tex2DTarget2D
		}
	}
	return &Texture2D{
		GLTexture: *tex.As(TexTarget(target)),
	}
}

func (tex *GLTexture) As3D(target TexTarget) *Texture3D {
	if target == 0 {
		if tex.Target.Dimensions() == 3 {
			target = tex.Target
		} else {
			target = Tex3DTargetCubeMap
		}
	}
	return &Texture3D{
		GLTexture: *tex.As(TexTarget(target)),
	}
}

func (tex *GLTexture) WrapMode(sMode, tMode, rMode TexWrapMode) {
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

func (tex *GLTexture) Bind(unit int) {
	gl.ActiveTexture(uint32(gl.TEXTURE0 + unit))
	gl.BindTexture(uint32(tex.Target), *tex.Id())
}

func (tex *GLTexture) Unbind(unit int) {
	gl.ActiveTexture(uint32(gl.TEXTURE0 + unit))
	gl.BindTexture(uint32(tex.Target), 0)
}

func (tex *GLTexture) BindFor(unit int, context utils.BindingClosure) {
	tex.Bind(unit)
	context()
	tex.Unbind(unit)
}

func (tex *GLTexture) FilterMode(minMode, magMode TexFilterMode) {
	if minMode != 0 {
		gl.TexParameteri(uint32(tex.Target), gl.TEXTURE_MIN_FILTER, int32(minMode))
	}
	if magMode != 0 {
		gl.TexParameteri(uint32(tex.Target), gl.TEXTURE_MAG_FILTER, int32(magMode))
	}
}

func (tex *GLTexture) GenerateMipmap() {
	gl.GenerateMipmap(uint32(tex.Target))
}

func (tex *GLTexture) ApplyDefaults() {
	tex.WrapMode(WrapDefault, WrapDefault, WrapDefault)
	tex.FilterMode(FilterMinDefault, FilterMagDefault)
	tex.MipMapLevels(MipMapDefaultBaseLevel, MipMapDefaultMaxLevel)
}

func (tex *GLTexture) MipMapLevels(base, max int) {
	gl.TexParameteri(uint32(tex.Target), gl.TEXTURE_BASE_LEVEL, int32(base))
	gl.TexParameteri(uint32(tex.Target), gl.TEXTURE_MAX_LEVEL, int32(max))
}

func dataPtr(data interface{}, glType uint32, len int) (unsafe.Pointer, error) {
	if data == nil {
		return unsafe.Pointer(nil), nil
	}

	minSize := getGlTypeSize(glType) * len
	size := binary.Size(data)
	if size == -1 {
		return nil, fmt.Errorf("data is not a fixed-size value or a slice of fixed-size values, or a pointer to such data. Actual type: %T", data)
	}
	if size < minSize {
		return nil, NotEnoughError{ActualSize: int64(size), RequiredSize: int64(minSize)}
	}
	return gl.Ptr(data), nil
}

func (tex *GLTexture) AllocType(level int32, internalFormat uint32, width, height, depth int32, format, dataType uint32, data interface{}) error {
	switch tex.Target.Dimensions() {
	case 1:
		addr, err := dataPtr(data, dataType, int(width))
		if err != nil {
			return err
		}
		gl.TexImage1D(uint32(tex.Target), level, int32(internalFormat), width, 0, format, dataType, addr)
	case 2:
		addr, err := dataPtr(data, dataType, int(width*height))
		if err != nil {
			return err
		}
		gl.TexImage2D(uint32(tex.Target), level, int32(internalFormat), width, height, 0, format, dataType, addr)
	case 3:
		addr, err := dataPtr(data, dataType, int(width*height*depth))
		if err != nil {
			return err
		}
		gl.TexImage3D(uint32(tex.Target), level, int32(internalFormat), width, height, depth, 0, format, dataType, addr)
	default:
		return fmt.Errorf("texture has unsupported dimensions of %d", tex.Target.Dimensions())
	}
	return nil
}

func (tex *GLTexture) Alloc(level int32, internalFormat uint32, width, height, depth int32, format uint32, data interface{}) error {
	typ, _, err := getGlType(data)
	if err != nil {
		return err
	}
	return tex.AllocType(level, internalFormat, width, height, depth, format, typ, data)
}

func (tex *GLTexture) AllocBytes(bytes []byte, level int32, internalFormat uint32, width, height, depth int32, format uint32) error {
	return tex.AllocType(level, internalFormat, width, height, depth, format, gl.UNSIGNED_BYTE, bytes)
}

func (tex *GLTexture) AllocEmpty(level int32, internalFormat uint32, width, height, depth int32, format uint32) error {
	return tex.AllocType(level, internalFormat, width, height, depth, format, gl.UNSIGNED_BYTE, nil)
}

func (tex *GLTexture) WriteType(level, x, y, z, width, height, depth int32, format, dataType uint32, data interface{}) error {
	switch tex.Target.Dimensions() {
	case 1:
		addr, err := dataPtr(data, dataType, int(width))
		if err != nil {
			return err
		}
		gl.TexSubImage1D(uint32(tex.Target), level, x, width, format, dataType, addr)
	case 2:
		addr, err := dataPtr(data, dataType, int(width*height))
		if err != nil {
			return err
		}
		gl.TexSubImage2D(uint32(tex.Target), level, x, y, width, height, format, dataType, addr)
	case 3:
		addr, err := dataPtr(data, dataType, int(width*height*depth))
		if err != nil {
			return err
		}
		gl.TexSubImage3D(uint32(tex.Target), level, x, y, z, width, height, depth, format, dataType, addr)
	default:
		return fmt.Errorf("texture has unsupported dimensions of %d", tex.Target.Dimensions())
	}
	return nil
}

func (tex *GLTexture) Write(level, x, y, z, width, height, depth int32, format uint32, data interface{}) error {
	typ, _, err := getGlType(data)
	if err != nil {
		return err
	}
	return tex.WriteType(level, x, y, z, width, height, depth, format, typ, data)
}

func (tex *GLTexture) WriteBytes(bytes []byte, level, x, y, z, width, height, depth int32, format uint32) error {
	return tex.WriteType(level, x, y, z, width, height, depth, format, gl.UNSIGNED_BYTE, bytes)
}

type Texture1D struct {
	GLTexture
}

func (tex *Texture1D) WrapMode(sMode TexWrapMode) {
	tex.GLTexture.WrapMode(sMode, 0, 0)
}

func (tex *Texture1D) AllocType(level int32, internalFormat uint32, width int32, format, dataType uint32, data interface{}) error {
	return tex.GLTexture.AllocType(level, internalFormat, width, 0, 0, format, dataType, data)
}

func (tex *Texture1D) Alloc(level int32, internalFormat uint32, width int32, format uint32, data interface{}) error {
	return tex.GLTexture.Alloc(level, internalFormat, width, 0, 0, format, data)
}

func (tex *Texture1D) AllocBytes(level int32, internalFormat uint32, width int32, format uint32, bytes []byte) error {
	return tex.AllocType(level, internalFormat, width, format, gl.UNSIGNED_BYTE, bytes)
}

func (tex *Texture1D) AllocEmpty(level int32, internalFormat uint32, width int32, format uint32) error {
	return tex.AllocType(level, internalFormat, width, format, gl.UNSIGNED_BYTE, nil)
}

func (tex *Texture1D) WriteType(level, x, width int32, format, dataType uint32, data interface{}) error {
	return tex.GLTexture.WriteType(level, x, 0, 0, width, 0, 0, format, dataType, data)
}

func (tex *Texture1D) Write(level, x, width int32, format uint32, data interface{}) error {
	return tex.GLTexture.Write(level, x, 0, 0, width, 0, 0, format, data)
}

func (tex *Texture1D) WriteBytes(level, x, width int32, format uint32, bytes []byte) error {
	return tex.WriteType(level, x, width, format, gl.UNSIGNED_BYTE, bytes)
}

type Texture2D struct {
	GLTexture
}

func (tex *Texture2D) WrapMode(sMode, tMode TexWrapMode) {
	tex.GLTexture.WrapMode(sMode, tMode, 0)
}

func (tex *Texture2D) AllocType(level int32, internalFormat uint32, width, height int32, format, dataType uint32, data interface{}) error {
	return tex.GLTexture.AllocType(level, internalFormat, width, height, 0, format, dataType, data)
}

func (tex *Texture2D) Alloc(level int32, internalFormat uint32, width, height int32, format uint32, data interface{}) error {
	return tex.GLTexture.Alloc(level, internalFormat, width, height, 0, format, data)
}

func (tex *Texture2D) AllocBytes(bytes []byte, level int32, internalFormat uint32, width, height int32, format uint32) error {
	return tex.AllocType(level, internalFormat, width, height, format, gl.UNSIGNED_BYTE, bytes)
}

func (tex *Texture2D) AllocEmpty(level int32, internalFormat uint32, width, height int32, format uint32) error {
	return tex.AllocType(level, internalFormat, width, height, format, gl.UNSIGNED_BYTE, nil)
}

func (tex *Texture2D) WriteType(level, x, y, width, height int32, format, dataType uint32, data interface{}) error {
	return tex.GLTexture.WriteType(level, x, y, 0, width, height, 0, format, dataType, data)
}

func (tex *Texture2D) Write(level, x, y, width, height int32, format uint32, data interface{}) error {
	return tex.GLTexture.Write(level, x, y, 0, width, height, 0, format, data)
}

func (tex *Texture2D) WriteBytes(level, x, y, width, height int32, format uint32, bytes []byte) error {
	return tex.WriteType(level, x, y, width, height, format, gl.UNSIGNED_BYTE, bytes)
}

type Texture3D struct {
	GLTexture
}

/*
	Order: right (+x), left (-x), top (+y), bottom (-y), back (+z), front (-z)
*/
func (tex *Texture3D) As2DSides() []Texture2D {
	return []Texture2D{
		*tex.GLTexture.As2D(Tex2DTargetCubeMapPositiveX + 0),
		*tex.GLTexture.As2D(Tex2DTargetCubeMapPositiveX + 1),
		*tex.GLTexture.As2D(Tex2DTargetCubeMapPositiveX + 2),
		*tex.GLTexture.As2D(Tex2DTargetCubeMapPositiveX + 3),
		*tex.GLTexture.As2D(Tex2DTargetCubeMapPositiveX + 4),
		*tex.GLTexture.As2D(Tex2DTargetCubeMapPositiveX + 5),
	}
}

func (tex *Texture3D) WrapMode(sMode, tMode, rMode TexWrapMode) {
	tex.GLTexture.WrapMode(sMode, tMode, rMode)
}

func (tex *Texture3D) AllocType(level int32, internalFormat uint32, width, height, depth int32, format, dataType uint32, data interface{}) error {
	return tex.GLTexture.AllocType(level, internalFormat, width, height, depth, format, dataType, data)
}

func (tex *Texture3D) Alloc(level int32, internalFormat uint32, width, height, depth int32, format uint32, data interface{}) error {
	return tex.GLTexture.Alloc(level, internalFormat, width, height, depth, format, data)
}

func (tex *Texture3D) AllocBytes(level int32, internalFormat uint32, width, height, depth int32, format uint32, bytes []byte) error {
	return tex.AllocType(level, internalFormat, width, height, depth, format, gl.UNSIGNED_BYTE, bytes)
}

func (tex *Texture3D) AllocEmpty(level int32, internalFormat uint32, width, height, depth int32, format uint32) error {
	return tex.AllocType(level, internalFormat, width, height, depth, format, gl.UNSIGNED_BYTE, nil)
}

func (tex *Texture3D) WriteType(level, x, y, z, width, height, depth int32, format, dataType uint32, data interface{}) error {
	return tex.GLTexture.WriteType(level, x, y, z, width, height, depth, format, dataType, data)
}

func (tex *Texture3D) Write(level, x, y, z, width, height, depth int32, format uint32, data interface{}) error {
	return tex.GLTexture.Write(level, x, y, z, width, height, depth, format, data)
}

func (tex *Texture3D) WriteBytes(level, x, y, z, width, height, depth int32, format uint32, bytes []byte) error {
	return tex.WriteType(level, x, y, z, width, height, depth, format, gl.UNSIGNED_BYTE, bytes)
}
