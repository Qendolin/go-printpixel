package march

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

func Scale(dst draw.Image, src image.Image) {
	sr := src.Bounds()
	dr := dst.Bounds()
	mx := float64(sr.Dx()-1) / float64(dr.Dx())
	my := float64(sr.Dy()-1) / float64(dr.Dy())
	for x := dr.Min.X; x < dr.Max.X; x++ {
		for y := dr.Min.Y; y < dr.Max.Y; y++ {
			gx, tx := math.Modf(float64(x) * mx)
			gy, ty := math.Modf(float64(y) * my)
			srcX, srcY := int(gx), int(gy)
			r00, g00, b00, a00 := src.At(srcX, srcY).RGBA()
			r10, g10, b10, a10 := src.At(srcX+1, srcY).RGBA()
			r01, g01, b01, a01 := src.At(srcX, srcY+1).RGBA()
			r11, g11, b11, a11 := src.At(srcX+1, srcY+1).RGBA()
			result := color.RGBA64{
				R: blerp(r00, r10, r01, r11, tx, ty),
				G: blerp(g00, g10, g01, g11, tx, ty),
				B: blerp(b00, b10, b01, b11, tx, ty),
				A: blerp(a00, a10, a01, a11, tx, ty),
			}
			dst.Set(x, y, result)
		}
	}
}

func clerp(s, e color.Color, t float64) color.Color {
	sr, sg, sb, sa := s.RGBA()
	er, eg, eb, ea := e.RGBA()
	return color.RGBA64{
		R: uint16(lerp(float64(sr), float64(er), t)),
		G: uint16(lerp(float64(sg), float64(eg), t)),
		B: uint16(lerp(float64(sb), float64(eb), t)),
		A: uint16(lerp(float64(sa), float64(ea), t)),
	}
}

func lerp(s, e, t float64) float64 { return s + (e-s)*t }
func blerp(c00, c10, c01, c11 uint32, tx, ty float64) uint16 {
	return uint16(lerp(
		lerp(float64(c00), float64(c10), tx),
		lerp(float64(c01), float64(c11), tx),
		ty,
	))
}
