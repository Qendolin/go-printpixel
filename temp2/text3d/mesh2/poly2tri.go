package mesh2

// #define NDEBUG
// #include <stdlib.h>
// #include <string.h>
//
// #define MPE_POLY2TRI_IMPLEMENTATION
// #include "poly2tri.h"
import "C"
import (
	"reflect"
	"unsafe"

	"github.com/Qendolin/go-printpixel/temp2/text3d/march"
)

const mpePolyPointSize = C.size_t(unsafe.Sizeof(C.MPEPolyPoint{}))

// Must not have duplicate vertices or self intersecions
// Returns nil if len(verts) < 3, removes holes with len(hole) < 3
func Poly2Tri(poly march.Polygon, indicesOffset uint32) (indices []uint32, vertices []vertex, ok bool) {
	totalVertCount := len(poly.ExteriorBound)
	if totalVertCount < 3 {
		return nil, nil, false
	}

	for _, hole := range poly.InteriorBounds {
		if len(hole) < 3 {
			continue
		}
		totalVertCount += len(hole)
	}

	vertices = make([]vertex, 0, totalVertCount)

	requiredMemory := C.MPE_PolyMemoryRequired(C.uint(totalVertCount))
	memory := C.calloc(requiredMemory, 1)
	defer C.free(memory)

	polyCtx := &C.MPEPolyContext{}
	C.MPE_PolyInitContext(polyCtx, memory, C.uint(totalVertCount))

	pointIndex := 0
	outerPoints := mpePushPointArray(polyCtx, len(poly.ExteriorBound))
	for i, v := range poly.ExteriorBound {
		outerPoints[i] = C.MPEPolyPoint{
			Index: C.uint(pointIndex),
			X:     C.float(v[0]),
			Y:     C.float(v[1]),
		}
		pointIndex++
	}
	C.MPE_PolyAddEdge(polyCtx)
	vertices = append(vertices, poly.ExteriorBound...)

	for _, bound := range poly.InteriorBounds {
		if len(bound) < 3 {
			continue
		}

		innerPoints := mpePushPointArray(polyCtx, len(bound))
		for i, v := range bound {
			innerPoints[i] = C.MPEPolyPoint{
				Index: C.uint(pointIndex),
				X:     C.float(v[0]),
				Y:     C.float(v[1]),
			}
			pointIndex++
		}
		C.MPE_PolyAddHole(polyCtx)
		vertices = append(vertices, bound...)
	}

	C.MPE_PolyTriangulate(polyCtx)

	mpeTriangles := mpePolyTriangleToSlice(polyCtx.Triangles, int(polyCtx.TriangleCount))
	indices = make([]uint32, polyCtx.TriangleCount*3)
	for i := 0; i < int(polyCtx.TriangleCount); i++ {
		tri := mpeTriangles[i]
		p0 := tri.Points[0]
		p1 := tri.Points[1]
		p2 := tri.Points[2]

		indices[i*3+0] = uint32(p0.Index) + indicesOffset
		indices[i*3+1] = uint32(p1.Index) + indicesOffset
		indices[i*3+2] = uint32(p2.Index) + indicesOffset
	}

	return indices, vertices, true
}

func mpePushPointArray(polyCtx *C.MPEPolyContext, count int) []C.MPEPolyPoint {
	arr := C.MPE_PolyPushPointArray(polyCtx, C.uint(count))
	var slice []C.MPEPolyPoint
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	header.Cap = count
	header.Len = count
	header.Data = uintptr(unsafe.Pointer(arr))

	return slice
}

func mpePolyTriangleToSlice(arr **C.MPEPolyTriangle, length int) []*C.MPEPolyTriangle {
	var slice []*C.MPEPolyTriangle
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	header.Cap = length
	header.Len = length
	header.Data = uintptr(unsafe.Pointer(arr))

	return slice
}
