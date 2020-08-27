package data

import (
	"strconv"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type ColorFormat uint32

const (
	Red              ColorFormat = 0b0001 << 8
	Green            ColorFormat = 0b0010 << 8
	Blue             ColorFormat = 0b0100 << 8
	Alpha            ColorFormat = 0b1000 << 8
	Standard         ColorFormat = 0b0001 << 12
	Compressed       ColorFormat = 0b0010 << 12
	Depth            ColorFormat = 0b0100 << 12
	Stencil          ColorFormat = 0b1000 << 12
	Float            ColorFormat = 0b0001 << 16
	Integer          ColorFormat = 0b0010 << 16
	UnsignedInteger  ColorFormat = 0b0100 << 16
	SignedNormalized ColorFormat = 0b1000 << 16
	Special          ColorFormat = 0b0001 << 20
	Reversed         ColorFormat = 0b0010 << 20
)

const (
	RG               = Red | Green
	RGB              = Red | Green | Blue
	RGBA             = Red | Green | Blue | Alpha
	R8               = Red | 8
	RG8              = RG | 8
	RGB8             = RGB | 8
	RGBA8            = RGBA | 8
	R8I              = R8 | Integer
	RG8I             = RG8 | Integer
	RGB8I            = RGB8 | Integer
	RGBA8I           = RGBA8 | Integer
	R8UI             = R8 | UnsignedInteger
	RG8UI            = RG8 | UnsignedInteger
	RGB8UI           = RGB8 | UnsignedInteger
	RGBA8UI          = RGBA8 | UnsignedInteger
	R8SN             = R8 | SignedNormalized
	RG8SN            = RG8 | SignedNormalized
	RGB8SN           = RGB8 | SignedNormalized
	RGBA8SN          = RGBA8 | SignedNormalized
	SRGB8            = RGB8 | Standard
	SRGBA8           = RGB8 | Standard
	R16              = Red | 16
	RG16             = RG | 16
	RGB16            = RGB | 16
	RGBA16           = RGBA | 16
	R16I             = R16 | Integer
	RG16I            = RG16 | Integer
	RGB16I           = RGB16 | Integer
	RGBA16I          = RGBA16 | Integer
	R16UI            = R16 | UnsignedInteger
	RG16UI           = RG16 | UnsignedInteger
	RGB16UI          = RGB16 | UnsignedInteger
	RGBA16UI         = RGBA16 | UnsignedInteger
	R16SN            = R16 | SignedNormalized
	RG16SN           = RG16 | SignedNormalized
	RGB16SN          = RGB16 | SignedNormalized
	RGBA16SN         = RGBA16 | SignedNormalized
	R16F             = R16 | Float
	RG16F            = RG16 | Float
	RGB16F           = RGB16 | Float
	RGBA16F          = RGBA16 | Float
	R32I             = Red | 32 | Integer
	RG32I            = RG | 32 | Integer
	RGB32I           = RGB | 32 | Integer
	RGBA32I          = RGBA | 32 | Integer
	R32UI            = Red | 32 | UnsignedInteger
	RG32UI           = RG | 32 | UnsignedInteger
	RGB32UI          = RGB | 32 | UnsignedInteger
	RGBA32UI         = RGBA | 32 | UnsignedInteger
	R32F             = Red | 32 | Float
	RG32F            = RG | 32 | Float
	RGB32F           = RGB | 32 | Float
	RGBA32F          = RGBA | 32 | Float
	Depth16          = Depth | 16
	Depth24          = Depth | 24
	Depth32          = Depth | 32
	Depth32F         = Depth | Float | 32
	DepthStencil     = Depth | Stencil
	Depth24Stencil8  = Depth | Stencil | 24 | 8
	Depth32FStencil8 = Depth | Stencil | Float | 32 | 8
	Stencil8         = Stencil | 8
)

func (cf ColorFormat) String() string {
	var s string
	if cf&Special != 0 {
		s += "Special "
	}
	if cf&Depth != 0 {
		s += "Depth"
	}
	if cf&Stencil != 0 {
		s += "Stencil"
	}
	if cf&Standard != 0 {
		s += "s"
	}
	var col string
	if cf&Red != 0 {
		col += "R"
	}
	if cf&Green != 0 {
		col += "G"
	}
	if cf&Blue != 0 {
		col += "B"
	}
	if cf&Reversed != 0 {
		var tmp string
		for _, v := range col {
			tmp = string(v) + tmp
		}
		col = tmp
	}
	if cf&Alpha != 0 {
		col += "A"
	}
	s += col

	s += " " + strconv.FormatUint(uint64(cf&0x3F), 10)

	if cf&SignedNormalized != 0 {
		s += "SNORM"
	}
	if cf&Integer != 0 {
		s += "I"
	}
	if cf&UnsignedInteger != 0 {
		s += "UI"
	}
	if cf&Float != 0 {
		s += "F"
	}

	return s
}

var ColorFormatDefault = RGBA8

func (cf ColorFormat) InternalFormatEnum() uint32 {
	if cf == 0 {
		if ColorFormatDefault == 0 {
			return 0
		}
		return ColorFormatDefault.InternalFormatEnum()
	}

	if cf&(DepthStencil) != 0 {
		cf &= 0x000FC03F
	} else {
		cf &= 0x001F3f3F
	}

	if glenum := getGlColorFormatEnum(cf); glenum != 0 {
		return glenum
	}

	if cf&(DepthStencil) != 0 {
		if cf&Float != 0 {
			return gl.DEPTH32F_STENCIL8
		} else if cf&0xFF == 0 {
			return gl.DEPTH_STENCIL
		}
		return gl.DEPTH24_STENCIL8
	}
	if cf&Depth != 0 {
		if cf&Float != 0 {
			return gl.DEPTH_COMPONENT32F
		} else if bits := cf & 0xFF; bits == 0 {
			return gl.DEPTH_COMPONENT
		} else if bits < 24 {
			return gl.DEPTH_COMPONENT16
		} else if bits < 32 {
			return gl.DEPTH_COMPONENT24
		}
		return gl.DEPTH_COMPONENT32
	}
	if cf&Stencil != 0 {
		if bits := cf & 0xFF; bits == 0 {
			return gl.STENCIL_INDEX8
		} else if bits < 4 {
			return gl.STENCIL_INDEX1
		} else if bits < 8 {
			return gl.STENCIL_INDEX4
		} else if bits < 16 {
			return gl.STENCIL_INDEX8
		}
		return gl.STENCIL_INDEX16
	}

	// cf is not a depth or stencil format

	comps := cf & 0x00000F00
	if comps == 0 {
		comps = RGBA
	} else if cf&Alpha != 0 {
		comps = RGBA
	} else if cf&Blue != 0 {
		comps = RGB
	} else if cf&Green != 0 {
		comps = RG
	}
	cf &= 0xFFFFF0FF
	cf |= comps

	typ := cf & 0x000F0000
	if typ := cf & 0x000F0000; typ != 0 && typ != Integer && typ != Float && typ != UnsignedInteger && typ != SignedNormalized {
		cf &= 0xFFF0FFFF
		typ = 0
	}

	if glenum := getGlColorFormatEnum(cf); glenum != 0 {
		return glenum
	}

	bits := cf & 0x000000FF
	switch typ {
	case 0, SignedNormalized:
		if bits < 16 {
			bits = 8
		} else {
			bits = 16
		}
	case Integer, UnsignedInteger:
		if bits < 16 {
			bits = 8
		} else if bits < 32 {
			bits = 16
		} else {
			bits = 32
		}
	case Float:
		if bits < 32 {
			bits = 16
		} else {
			bits = 32
		}
	}
	cf &= 0xFFFFFF00
	cf |= bits

	if glenum := getGlColorFormatEnum(cf); glenum != 0 {
		return glenum
	}

	cf &= ^(Compressed | Special)
	if glenum := getGlColorFormatEnum(cf); glenum != 0 {
		return glenum
	}

	return gl.RGBA8
}

func getGlColorFormatEnum(cf ColorFormat) uint32 {
	switch cf {
	case Red:
		return gl.RED
	case RG:
		return gl.RG
	case RGB:
		return gl.RGB
	case RGBA:
		return gl.RGBA
	case RGBA | 2:
		return gl.RGBA2
	case RGB | 4:
		return gl.RGB4
	case RGBA | 4:
		return gl.RGBA4
	case RGB | 5:
		return gl.RGB5
	case RGB | 10:
		return gl.RGB10
	case RGB | 12:
		return gl.RGB12
	case RGBA | 12:
		return gl.RGBA12
	case R8:
		return gl.R8
	case RG8:
		return gl.RG8
	case RGB8:
		return gl.RGB8
	case RGBA8:
		return gl.RGBA8
	case R16:
		return gl.R16
	case RG16:
		return gl.RG16
	case RGB16:
		return gl.RGB16
	case RGBA16:
		return gl.RGBA16
	case R8SN:
		return gl.R8_SNORM
	case RG8SN:
		return gl.RG8_SNORM
	case RGB8SN:
		return gl.RGB8_SNORM
	case RGBA8SN:
		return gl.RGBA8_SNORM
	case R16SN:
		return gl.R16_SNORM
	case RG16SN:
		return gl.RG16_SNORM
	case RGB16SN:
		return gl.RGB16_SNORM
	case RGBA16SN:
		return gl.RGBA16_SNORM
	case R8UI:
		return gl.R8UI
	case RG8UI:
		return gl.RG8UI
	case RGB8UI:
		return gl.RGB8UI
	case RGBA8UI:
		return gl.RGBA8UI
	case R16UI:
		return gl.R16UI
	case RG16UI:
		return gl.RG16UI
	case RGB16UI:
		return gl.RGB16UI
	case RGBA16UI:
		return gl.RGBA16UI
	case R32UI:
		return gl.R32UI
	case RG32UI:
		return gl.RG32UI
	case RGB32UI:
		return gl.RGB32UI
	case RGBA32UI:
		return gl.RGBA32UI
	case R8I:
		return gl.R8I
	case RG8I:
		return gl.RG8I
	case RGB8I:
		return gl.RGB8I
	case RGBA8I:
		return gl.RGBA8I
	case R16I:
		return gl.R16I
	case RG16I:
		return gl.RG16I
	case RGB16I:
		return gl.RGB16I
	case RGBA16I:
		return gl.RGBA16I
	case R32I:
		return gl.R32I
	case RG32I:
		return gl.RG32I
	case RGB32I:
		return gl.RGB32I
	case RGBA32I:
		return gl.RGBA32I
	case R16F:
		return gl.R16F
	case RG16F:
		return gl.RG16F
	case RGB16F:
		return gl.RGB16F
	case RGBA16F:
		return gl.RGBA16F
	case R32F:
		return gl.R32F
	case RG32F:
		return gl.RG32F
	case RGB32F:
		return gl.RGB32F
	case RGBA32F:
		return gl.RGBA32F
	case Red | Compressed:
		return gl.COMPRESSED_RED
	case RG | Compressed:
		return gl.COMPRESSED_RG
	case RGB | Compressed:
		return gl.COMPRESSED_RGB
	case RGBA | Compressed:
		return gl.COMPRESSED_RGBA
	case RGB | Standard | Compressed:
		return gl.COMPRESSED_SRGB
	case RGBA | Standard | Compressed:
		return gl.COMPRESSED_SRGB_ALPHA
	case Depth:
		return gl.DEPTH_COMPONENT
	case Depth16:
		return gl.DEPTH_COMPONENT16
	case Depth24:
		return gl.DEPTH_COMPONENT24
	case Depth32:
		return gl.DEPTH_COMPONENT32
	case Depth32F:
		return gl.DEPTH_COMPONENT32F
	case DepthStencil:
		return gl.DEPTH_STENCIL
	case Depth24Stencil8:
		return gl.DEPTH24_STENCIL8
	case Depth32FStencil8:
		return gl.DEPTH32F_STENCIL8
	case Stencil | 1:
		return gl.STENCIL_INDEX1
	case Stencil | 4:
		return gl.STENCIL_INDEX4
	case Stencil | 8:
		return gl.STENCIL_INDEX8
	case Stencil | 16:
		return gl.STENCIL_INDEX16
	case RGBA | Special | 12:
		return gl.RGB10_A2
	case RGBA | Special | UnsignedInteger | 12:
		return gl.RGB10_A2UI
	case RGB | Special | Float | 32:
		return gl.R11F_G11F_B10F
	case RGB | Special | 8:
		return gl.R3_G3_B2
	case RGBA | Special | 6:
		return gl.RGB5_A1
	case RGB | Special | Float | 14:
		return gl.RGB9_E5
	case RGB | Special | 16:
		return gl.RGB565
	}
	return 0
}

func (cf ColorFormat) PixelFormatEnum() uint32 {
	if cf == 0 {
		if ColorFormatDefault == 0 {
			return 0
		}
		return ColorFormatDefault.PixelFormatEnum()
	}

	if cf&Compressed != 0 {
		return 0
	}

	switch ds := cf & (Depth | Stencil); ds {
	case Depth | Stencil:
		return gl.DEPTH_STENCIL
	case Depth:
		return gl.DEPTH_COMPONENT
	case Stencil:
		return gl.STENCIL_INDEX
	}

	rev := cf&Reversed != 0

	colBits := cf & RGBA
	if cf&(Integer|UnsignedInteger) == 0 {
		switch {
		case (colBits <= Red):
			return gl.RED
		case colBits <= RG:
			return gl.RG
		case colBits <= RGB && rev:
			return gl.BGR
		case colBits <= RGB:
			return gl.RGB
		case colBits <= RGBA && rev:
			return gl.BGRA
		case colBits <= RGBA:
			return gl.RGBA
		}
	} else {
		switch {
		case (colBits <= Red):
			return gl.RED_INTEGER
		case colBits <= RG:
			return gl.RG_INTEGER
		case colBits <= RGB && rev:
			return gl.BGR_INTEGER
		case colBits <= RGB:
			return gl.RGB_INTEGER
		case colBits <= RGBA && rev:
			return gl.BGRA_INTEGER
		case colBits <= RGBA:
			return gl.RGBA_INTEGER
		}
	}

	return gl.RGBA
}
