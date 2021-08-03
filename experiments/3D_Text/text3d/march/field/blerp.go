package field

import (
	"math"
	"unsafe"
)

// #cgo CFLAGS: -O3 -mfpmath=sse -mavx -mavx2 -msse -msse2 -msse -msse3 -mssse3 -msse4 -msse4.1 -msse4.2 -march=native
// #include "blerp.h"
import "C"

func ScaleBlerpIC(src, dst *ValueFieldI) {
	csrc := C.struct_ValueFieldIC{
		Width:  C.uint32_t(src.Width),
		Height: C.uint32_t(src.Height),
		Len:    C.uint32_t(len(src.Values)),
		Data:   (*C.uint32_t)(unsafe.Pointer(&src.Values[0])),
	}

	cdst := C.struct_ValueFieldIC{
		Width:  C.uint32_t(dst.Width),
		Height: C.uint32_t(dst.Height),
		Len:    C.uint32_t(len(dst.Values)),
		Data:   (*C.uint32_t)(unsafe.Pointer(&dst.Values[0])),
	}

	C.ScaleBlerpC(csrc, cdst)
}

func ScaleBlerpFSimd(src, dst *ValueField) {
	csrc := C.struct_ValueFieldFC{
		Width:  C.uint32_t(src.Width),
		Height: C.uint32_t(src.Height),
		Len:    C.uint32_t(len(src.Values)),
		Data:   (*C.float)(unsafe.Pointer(&src.Values[0])),
	}

	cdst := C.struct_ValueFieldFC{
		Width:  C.uint32_t(dst.Width),
		Height: C.uint32_t(dst.Height),
		Len:    C.uint32_t(len(dst.Values)),
		Data:   (*C.float)(unsafe.Pointer(&dst.Values[0])),
	}

	C.ScaleBlerpCFSimd(csrc, cdst)
}

func ScaleBlerpFSimd2(src, dst *ValueField) {
	csrc := C.struct_ValueFieldFC{
		Width:  C.uint32_t(src.Width),
		Height: C.uint32_t(src.Height),
		Len:    C.uint32_t(len(src.Values)),
		Data:   (*C.float)(unsafe.Pointer(&src.Values[0])),
	}

	cdst := C.struct_ValueFieldFC{
		Width:  C.uint32_t(dst.Width),
		Height: C.uint32_t(dst.Height),
		Len:    C.uint32_t(len(dst.Values)),
		Data:   (*C.float)(unsafe.Pointer(&dst.Values[0])),
	}

	C.ScaleBlerpCFSimd2(csrc, cdst)
}

func ScaleBlerpFull(src, dst *ValueField) {
	if src.ComponentSize != 4 || dst.ComponentSize != 4 {
		panic("src and dst component size must be 4")
	}

	if src.Width == dst.Width && src.Height == dst.Height {
		copy(dst.Values[:len(src.Values)], src.Values)
		return
	}

	csrc := C.struct_ValueFieldFC{
		Width:  C.uint32_t(src.Width),
		Height: C.uint32_t(src.Height),
		Len:    C.uint32_t(len(src.Values)),
		Data:   (*C.float)(unsafe.Pointer(&src.Values[0])),
	}

	cdst := C.struct_ValueFieldFC{
		Width:  C.uint32_t(dst.Width),
		Height: C.uint32_t(dst.Height),
		Len:    C.uint32_t(len(dst.Values)),
		Data:   (*C.float)(unsafe.Pointer(&dst.Values[0])),
	}

	C.ScaleBlerpCFull(csrc, cdst)
}

func ScaleBlerpSimd(src, dst *ValueFieldI) {
	csrc := C.struct_ValueFieldIC{
		Width:  C.uint32_t(src.Width),
		Height: C.uint32_t(src.Height),
		Len:    C.uint32_t(len(src.Values)),
		Data:   (*C.uint32_t)(unsafe.Pointer(&src.Values[0])),
	}

	cdst := C.struct_ValueFieldIC{
		Width:  C.uint32_t(dst.Width),
		Height: C.uint32_t(dst.Height),
		Len:    C.uint32_t(len(dst.Values)),
		Data:   (*C.uint32_t)(unsafe.Pointer(&dst.Values[0])),
	}

	C.ScaleBlerpCSimd(csrc, cdst)
}

func ScaleBlerpSimd2(src, dst *ValueFieldI) {
	csrc := C.struct_ValueFieldIC{
		Width:  C.uint32_t(src.Width),
		Height: C.uint32_t(src.Height),
		Len:    C.uint32_t(len(src.Values)),
		Data:   (*C.uint32_t)(unsafe.Pointer(&src.Values[0])),
	}

	cdst := C.struct_ValueFieldIC{
		Width:  C.uint32_t(dst.Width),
		Height: C.uint32_t(dst.Height),
		Len:    C.uint32_t(len(dst.Values)),
		Data:   (*C.uint32_t)(unsafe.Pointer(&dst.Values[0])),
	}

	C.ScaleBlerpCSimd2(csrc, cdst)
}

func ScaleBlerpSimd3(src *ValueFieldI8, dst *ValueFieldI) {
	csrc := C.struct_ValueFieldI8C{
		Width:  C.uint32_t(src.Width),
		Height: C.uint32_t(src.Height),
		Len:    C.uint32_t(len(src.Values)),
		Data:   (*C.uint8_t)(unsafe.Pointer(&src.Values[0])),
	}

	cdst := C.struct_ValueFieldIC{
		Width:  C.uint32_t(dst.Width),
		Height: C.uint32_t(dst.Height),
		Len:    C.uint32_t(len(dst.Values)),
		Data:   (*C.uint32_t)(unsafe.Pointer(&dst.Values[0])),
	}

	C.ScaleBlerpCSimd3(csrc, cdst)
}

func ScaleBlerp(src, dst *ValueField) {
	mx := float64(src.Width-1) / float64(dst.Width)
	my := float64(src.Height-1) / float64(dst.Height)
	for y := 0; y < dst.Height; y++ {
		for x := 0; x < dst.Width; x++ {
			gx, tx := math.Modf(float64(x) * mx)
			gy, ty := math.Modf(float64(y) * my)
			srcX, srcY := int(gx), int(gy)
			rgba00 := src.GetComponent(srcX, srcY)
			rgba10 := src.GetComponent(srcX+1, srcY)
			rgba01 := src.GetComponent(srcX, srcY+1)
			rgba11 := src.GetComponent(srcX+1, srcY+1)
			result := []float32{
				blerp(rgba00[0], rgba10[0], rgba01[0], rgba11[0], float32(tx), float32(ty)),
				blerp(rgba00[1], rgba10[1], rgba01[1], rgba11[1], float32(tx), float32(ty)),
				blerp(rgba00[2], rgba10[2], rgba01[2], rgba11[2], float32(tx), float32(ty)),
				// blerp(rgba00[3], rgba10[3], rgba01[3], rgba11[3], float32(tx), float32(ty)),
			}
			dst.SetComponent(x, y, result)
		}
	}
}

func lerp(s, e float32, t float32) float32 { return s + (e-s)*t }
func blerp(c00, c10, c01, c11 float32, tx, ty float32) float32 {
	return lerp(
		lerp(c00, c10, tx),
		lerp(c01, c11, tx),
		ty,
	)
}

func ScaleBlerpIptr(src, dst *ValueFieldI) {
	const size uintptr = 4

	sp := unsafe.Pointer(&src.Values[0])
	dp := unsafe.Pointer(&dst.Values[0])

	send := uintptr(len(src.Values)) * size
	srowstride := uintptr(src.Width) * 3 * size

	drowstride := uintptr(dst.Width) * 3 * size
	dTxF := uint64(math.MaxUint32 / dst.Width)
	dTyF := uint64(math.MaxUint32 / dst.Height)
	sTxF := uint64(math.MaxUint32 / (src.Width - 1))
	sTyF := uint64(math.MaxUint32 / (src.Height - 1))

	// s2dxF := ((uint64(dst.Width) * math.MaxUint32) / uint64(src.Width-1)) >> 32
	// s2dyF := ((uint64(dst.Height) * math.MaxUint32) / uint64(src.Height-1)) >> 32
	d2sxF := uint32((uint64(src.Width-1) * math.MaxUint32) / uint64(dst.Width))
	d2syF := uint32((uint64(src.Height-1) * math.MaxUint32) / uint64(dst.Height))

	var dyF, dyFend uint64
	var dy uintptr
	var fy uint32

	for sy := uintptr(0); sy < send-srowstride; sy += srowstride {
		syp := unsafe.Pointer(uintptr(sp) + sy)

		dyFend += sTyF
		for ; dyF < dyFend; dyF += dTyF {
			dyp := unsafe.Pointer(uintptr(dp) + dy)

			var dxF, dxFend uint64
			var dx uintptr
			var fx uint32

			for sx := uintptr(0); sx < srowstride-3*size; sx += 3 * size {
				sxp := unsafe.Pointer(uintptr(syp) + sx)

				var (
					r00 = uint32(*(*uint32)(unsafe.Pointer(uintptr(sxp) + 0*size)))
					g00 = uint32(*(*uint32)(unsafe.Pointer(uintptr(sxp) + 1*size)))
					b00 = uint32(*(*uint32)(unsafe.Pointer(uintptr(sxp) + 2*size)))
					r10 = uint32(*(*uint32)(unsafe.Pointer(uintptr(sxp) + 3*size)))
					g10 = uint32(*(*uint32)(unsafe.Pointer(uintptr(sxp) + 4*size)))
					b10 = uint32(*(*uint32)(unsafe.Pointer(uintptr(sxp) + 5*size)))
					r01 = uint32(*(*uint32)(unsafe.Pointer(uintptr(sxp) + 0*size + srowstride)))
					g01 = uint32(*(*uint32)(unsafe.Pointer(uintptr(sxp) + 1*size + srowstride)))
					b01 = uint32(*(*uint32)(unsafe.Pointer(uintptr(sxp) + 2*size + srowstride)))
					r11 = uint32(*(*uint32)(unsafe.Pointer(uintptr(sxp) + 3*size + srowstride)))
					g11 = uint32(*(*uint32)(unsafe.Pointer(uintptr(sxp) + 4*size + srowstride)))
					b11 = uint32(*(*uint32)(unsafe.Pointer(uintptr(sxp) + 5*size + srowstride)))
				)

				dxFend += sTxF
				for ; dxF < dxFend; dxF += dTxF {
					*(*uint32)(unsafe.Pointer(uintptr(dyp) + dx + 0*size)) = blerpIptr(r00, r10, r01, r11, uint32(fx), uint32(fy))
					*(*uint32)(unsafe.Pointer(uintptr(dyp) + dx + 1*size)) = blerpIptr(g00, g10, g01, g11, uint32(fx), uint32(fy))
					*(*uint32)(unsafe.Pointer(uintptr(dyp) + dx + 2*size)) = blerpIptr(b00, b10, b01, b11, uint32(fx), uint32(fy))

					dx += 3 * size
					// fx = (fx + d2sxF) & math.MaxUint32
					fx += d2sxF
				}
			}

			dy += drowstride
			// fy = (fy + d2syF) & math.MaxUint32
			fy += d2syF
		}
	}
}

// func ScaleBlerpI(src, dst *ValueFieldI) {
// 	mx := uint64((src.Width - 1) * math.MaxUint32 / dst.Width)
// 	// mx := uint64((src.Width-1)*math.MaxUint32/dst.Width) >> 8
// 	// mxi := uint64((math.MaxUint32 / dst.Width) * (src.Width))
// 	// mxi := uint64(math.MaxUint32/src.Width) * uint64(dst.Width)
// 	// mxi := uint64(math.MaxUint32/(src.Width-1)) * uint64(dst.Width)
// 	// mxi := uint64(math.MaxUint32/(src.Width-1)) * uint64(dst.Width) >> 8

// 	// mx := uint64((math.MaxUint32 / dst.Width) * (src.Width - 1))
// 	// mx := uint64((dst.Width) * math.MaxUint32 / src.Width)
// 	my := uint64((src.Height - 1) * math.MaxUint32 / dst.Height)

// 	sstride := uint64(src.Width * 3)
// 	dstride := uint64(dst.Width * 3)

// 	var srow0 []uint32
// 	srow1 := src.Values[:sstride]
// 	// 1 to pass check on first iter
// 	var srcY uint64 = 1

// 	for y := uint64(0); y < uint64(dst.Height); y++ {
// 		gy := (y * my) >> 32
// 		ty := (y * my) & math.MaxUint32
// 		coty := math.MaxUint32 - ty

// 		if gy != srcY {
// 			srcY = gy
// 			srow0 = srow1
// 			srow1 = src.Values[(srcY+1)*sstride : (srcY+2)*sstride]
// 		}

// 		drow := dst.Values[y*dstride : (y+1)*dstride]

// 		var dx, sx, ax, ex uint64
// 		for sx < sstride-3 {
// 			rgba00 := srow0[sx : sx+6]
// 			rgba01 := srow1[sx : sx+6]
// 			var (
// 				c10b = uint64(rgba00[5])
// 				c10g = uint64(rgba00[4])
// 				c10r = uint64(rgba00[3])
// 				c00b = uint64(rgba00[2])
// 				c00g = uint64(rgba00[1])
// 				c00r = uint64(rgba00[0])
// 				c11b = uint64(rgba01[5])
// 				c11g = uint64(rgba01[4])
// 				c11r = uint64(rgba01[3])
// 				c01b = uint64(rgba01[2])
// 				c01g = uint64(rgba01[1])
// 				c01r = uint64(rgba01[0])
// 			)

// 			ex += 0xffffffff

// 			// this loop condition is slow
// 			for ax < ex && dx+2 < uint64(len(drow)) {
// 				tx := ax & 0xffffffff
// 				cotx := math.MaxUint32 - tx

// 				drow[dx+0] = blerpI(c00r, c10r, c01r, c11r, tx, ty, cotx, coty)
// 				drow[dx+1] = blerpI(c00g, c10g, c01g, c11g, tx, ty, cotx, coty)
// 				drow[dx+2] = blerpI(c00b, c10b, c01b, c11b, tx, ty, cotx, coty)

// 				ax += mx
// 				dx += 3
// 			}

// 			sx += 3
// 		}

// 		// var ax uint64
// 		// for sx := uint64(0); sx < uint64(src.Width)-1; sx++ {
// 		// 	// gx := sx * mxi
// 		// 	gx := (sx + 1) * mxi
// 		// 	rgba00 := srow0[sx*3 : sx*3+6]
// 		// 	rgba01 := srow1[sx*3 : sx*3+6]

// 		// 	for ax <= gx {
// 		// 		dx := ax >> 24
// 		// 		tx := ((ax * mx) >> 16) & 0xffffffff

// 		// 		// tx := (px / mxi) * 0xffffffff
// 		// 		drow[dx*3+0] = blerpI(rgba00[0], rgba00[3], rgba01[0], rgba01[3], tx, ty)
// 		// 		drow[dx*3+1] = blerpI(rgba00[1], rgba00[4], rgba01[1], rgba01[4], tx, ty)
// 		// 		drow[dx*3+2] = blerpI(rgba00[2], rgba00[5], rgba01[2], rgba01[5], tx, ty)
// 		// 		// ax += mx
// 		// 		ax += 0xffffff
// 		// 	}
// 		// }

// 		// var tx uint64
// 		// for sx := uint64(0); sx < uint64(src.Width)-1; sx++ {
// 		// 	// gx := sx * mxi
// 		// 	gx := sx * mxi
// 		// 	rgba00 := srow0[sx*3 : sx*3+6]
// 		// 	rgba01 := srow1[sx*3 : sx*3+6]

// 		// 	for tx <= math.MaxUint32 {
// 		// 		// dx := (gx + px) >> 32
// 		// 		dx := (gx + (tx*mxi)>>32) >> 24
// 		// 		// tx := (px / mxi) * 0xffffffff
// 		// 		drow[dx*3+0] = blerpI(rgba00[0], rgba00[3], rgba01[0], rgba01[3], tx, ty)
// 		// 		drow[dx*3+1] = blerpI(rgba00[1], rgba00[4], rgba01[1], rgba01[4], tx, ty)
// 		// 		drow[dx*3+2] = blerpI(rgba00[2], rgba00[5], rgba01[2], rgba01[5], tx, ty)
// 		// 		tx += mx
// 		// 	}

// 		// 	tx &= math.MaxUint32

// 		// 	// for tx := uint64(0); tx <= math.MaxUint32; tx += mx {
// 		// 	// 	// dx := gx + (mx*tx)>>32
// 		// 	// 	dx := ((gx) * uint64(dst.Width)) >> 32
// 		// 	// 	drow[dx*3+0] = blerpI(rgba00[0], rgba00[3], rgba01[0], rgba01[3], tx, ty)
// 		// 	// 	drow[dx*3+1] = blerpI(rgba00[1], rgba00[4], rgba01[1], rgba01[4], tx, ty)
// 		// 	// 	drow[dx*3+2] = blerpI(rgba00[2], rgba00[5], rgba01[2], rgba01[5], tx, ty)
// 		// 	// }
// 		// }

// 		// for x := uint64(0); x < uint64(dst.Width); x++ {
// 		// 	gx := (x * mx) >> 32            // eq. / math.MaxUint32
// 		// 	tx := (x * mx) & math.MaxUint32 // eq. % (math.MaxUint32 + 1) or % 2^32

// 		// 	srcX := gx * 3
// 		// 	rgba00 := srow0[srcX+0 : srcX+6]
// 		// 	rgba01 := srow1[srcX+0 : srcX+6]

// 		// 	// rgba00 := srow0[srcX+0 : srcX+3]
// 		// 	// rgba10 := srow0[srcX+3 : srcX+6]
// 		// 	// rgba01 := srow1[srcX+0 : srcX+3]
// 		// 	// rgba11 := srow1[srcX+3 : srcX+6]

// 		// 	// result := []uint32{
// 		// 	// 	blerpI(rgba00[0], rgba10[0], rgba01[0], rgba11[0], tx, ty),
// 		// 	// 	blerpI(rgba00[1], rgba10[1], rgba01[1], rgba11[1], tx, ty),
// 		// 	// 	blerpI(rgba00[2], rgba10[2], rgba01[2], rgba11[2], tx, ty),
// 		// 	// }

// 		// 	// copy(drow[x*3:x*3+3], result)
// 		// 	drow[x*3+0] = blerpI(rgba00[0], rgba00[3], rgba01[0], rgba01[3], tx, ty)
// 		// 	drow[x*3+1] = blerpI(rgba00[1], rgba00[4], rgba01[1], rgba01[4], tx, ty)
// 		// 	drow[x*3+2] = blerpI(rgba00[2], rgba00[5], rgba01[2], rgba01[5], tx, ty)
// 		// 	// dst.SetComponent(int(x), int(y), result)
// 		// }
// 	}
// }

// func lerpI64(s, e, f uint64) uint64 {
// 	// basically s + f*(e-s)
// 	return (s + (f*(e-s))>>32)
// }

// func lerpI64(s, e, f uint64) uint64 {
// 	// basically s + f*(e-s)
// 	return (s*(math.MaxUint32-f) + e*f) >> 32
// }

// func lerpI32(s, e uint32, f uint64) uint32 {
// 	// basically s + f*(e-s)
// 	return s + uint32((f*(uint64(e)-uint64(s)))>>32)
// }

// func lerpI32(s, e uint32, f uint32) uint32 {
// 	// basically s + f*(e-s)
// 	return s + uint32((uint64(f)*(uint64(e)-uint64(s)))>>32)
// }

// func lerpI(s, e, f, cof uint64) uint64 {
// 	// basically s * (1 - f) + b * f
// 	return (s*cof + e*f) >> 32
// }

// func lerpI(s, e uint32, f uint64) uint32 {
// 	// basically s * (1 - f) + b * f
// 	return uint32(
// 		(uint64(s)*(math.MaxUint32-f) + uint64(e)*f) /
// 			math.MaxUint32)
// }

// func lerpI(s, e uint64, f uint64) uint64 {
// 	// basically s * (1 - f) + b * f
// 	return uint64(
// 		(uint64(s)*(math.MaxUint32-f) + uint64(e)*f) /
// 			math.MaxUint32)
// }

// func blerpI(c00, c10, c01, c11 uint32, tx, ty uint64) uint32 {
// 	return uint32(lerpI(
// 		lerpI(uint64(c00), uint64(c10), tx),
// 		lerpI(uint64(c01), uint64(c11), tx),
// 		ty,
// 	))
// }

// func blerpI(c00, c10, c01, c11 uint64, tx, ty uint64) uint32 {
// 	return uint32(lerpI32(
// 		uint32(lerpI64(c00, c10, tx)),
// 		uint32(lerpI64(c01, c11, tx)),
// 		ty,
// 	))
// }

// func blerpI(c00, c10, c01, c11 uint32, tx, ty uint64) uint32 {
// 	return lerpI(
// 		lerpI(c00, c10, tx),
// 		lerpI(c01, c11, tx),
// 		ty,
// 	)
// }

// func blerpI(c00, c10, c01, c11 uint32, tx, ty uint64) uint32 {
// 	return lerpI32(
// 		lerpI32(c00, c10, tx),
// 		lerpI32(c01, c11, tx),
// 		ty,
// 	)
// }

// func blerpI(c00, c10, c01, c11 uint64, tx, ty, cotx, coty uint64) uint32 {
// 	return uint32(lerpI(
// 		lerpI(c00, c10, tx, cotx),
// 		lerpI(c01, c11, tx, cotx),
// 		ty, coty,
// 	))
// }

func lerpIptr(s, e, f uint32) uint32 {
	return (s + (e-s)*f)
}

// func blerpIptr(c00, c10, c01, c11 uint32, tx, ty uint32) uint32 {
// 	return uint32(lerpIptr(
// 		lerpIptr(uint64(c00), uint64(c10), uint64(tx)),
// 		lerpIptr(uint64(c01), uint64(c11), uint64(tx)),
// 		uint64(ty),
// 	))
// }

func blerpIptr(c00, c10, c01, c11 uint32, tx, ty uint32) uint32 {
	return uint32(lerpIptr(
		lerpIptr(c00, c10, tx),
		lerpIptr(c01, c11, tx),
		ty,
	))
}
