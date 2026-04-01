package text3d

import (
	"image"
	_ "image/png"
	"math"

	"github.com/Qendolin/go-printpixel/core/data"
	"github.com/Qendolin/go-printpixel/experiments/3D_Text/text3d/march"
	"github.com/Qendolin/go-printpixel/experiments/3D_Text/text3d/march/field"
	"github.com/Qendolin/go-printpixel/experiments/3D_Text/text3d/mesh2"
)

type vec2 = [2]float32
type vec3 = [3]float32

const Epsilon = 1e-6

func Load3d(img image.Image, scale float32, reducer mesh2.Reducer, origin vec3, depth float32) (mesh data.TriMesh, bounds image.Rectangle, err error) {
	msdf, err := loadMsdf(img, scale)
	if err != nil {
		return
	}

	origin = vec3{origin[0] * float32(msdf.Width()), origin[1] * float32(msdf.Height()), origin[2] * depth}
	polygons := march.March(msdf, true)
	totalVertCount := reducePolygons(polygons, reducer)
	sideStride := uint32(totalVertCount)
	totalVertCount *= 4

	// the layout is [front back front back]. (duplicated to enable hard edges)
	vertexBuffer := make([]vec3, totalVertCount)
	normalBuffer := make([]vec3, totalVertCount/2, totalVertCount)
	indexBuffer := make([]uint32, 0, totalVertCount*3)
	restartIndices := make([]uint32, 0, len(polygons))
	var restartIndex uint32

	for _, poly := range polygons {
		if poly.ExteriorBound == nil {
			continue
		}
		indices, vertices, ok := mesh2.Poly2Tri(poly, restartIndex)
		if !ok {
			continue
		}

		for i, v := range vertices {
			var (
				iFront = int(restartIndex) + i
				iBack  = iFront + int(sideStride)
			)
			vertexBuffer[iFront] = vec3{v[0] - origin[0], v[1] - origin[1], 0 + origin[2]}
			vertexBuffer[iBack] = vec3{v[0] - origin[0], v[1] - origin[1], -depth + origin[2]}
			normalBuffer[iFront] = vec3{0, 0, 1}
			normalBuffer[iBack] = vec3{0, 0, -1}
		}

		restartIndex += uint32(len(poly.ExteriorBound))
		restartIndices = append(restartIndices, restartIndex)
		for _, bound := range poly.InteriorBounds {
			restartIndex += uint32(len(bound))
			restartIndices = append(restartIndices, restartIndex)
		}
		indexBuffer = append(indexBuffer, indices...)
	}

	// Duplicate front and back face for hard edges
	copy(vertexBuffer[sideStride*2:sideStride*4], vertexBuffer[sideStride*0:sideStride*2])

	faceIndicesCount := len(indexBuffer)
	// Duplicate the indices for the back face but offset them and reverse the winding order
	for i := 0; i < faceIndicesCount; i += 3 {
		indexBuffer = append(indexBuffer,
			sideStride+indexBuffer[i+2],
			sideStride+indexBuffer[i+1],
			sideStride+indexBuffer[i+0])
	}

	lastFaceRestart := sideStride * 2
	flipNormals := false
	if depth < 0 {
		flipNormals = true
	}
	var (
		v0 = vertexBuffer[restartIndices[0]-1]
		v1 = vertexBuffer[0]
		v2 vec3
	)

	// Connect edge loops and calculate normals
	for i := uint32(1); i <= sideStride; i++ {
		idx := i + sideStride*2
		if i == restartIndices[0] {
			restartIndices = restartIndices[1:]

			var (
				firstFront = lastFaceRestart
				lastFront  = idx - 1
				firstBack  = lastFaceRestart + sideStride
				lastBack   = idx + sideStride - 1
			)

			indexBuffer = append(indexBuffer, firstBack, firstFront, lastFront)
			indexBuffer = append(indexBuffer, lastFront, lastBack, firstBack)

			v2 = vertexBuffer[lastFaceRestart]
			normalBuffer = append(normalBuffer, calcVertexNormalXY(v0, v1, v2, flipNormals))
			if len(restartIndices) != 0 {
				v0, v1, v2 = vertexBuffer[restartIndices[0]-1], vertexBuffer[i], vec3{}
			}

			lastFaceRestart = idx
			continue
		}

		indexBuffer = append(indexBuffer, idx, idx-1, idx+sideStride-1)
		indexBuffer = append(indexBuffer, idx+sideStride-1, idx+sideStride, idx)

		v2 = vertexBuffer[i]
		normalBuffer = append(normalBuffer, calcVertexNormalXY(v0, v1, v2, flipNormals))
		v0, v1 = v1, v2
	}

	// Duplicate normals for back face
	normalBuffer = normalBuffer[:cap(normalBuffer)]
	copy(normalBuffer[sideStride*3:sideStride*4], normalBuffer[sideStride*2:sideStride*3])

	mesh = makeTriMesh(vertexBuffer, normalBuffer, indexBuffer)
	return mesh, image.Rect(0, 0, msdf.Width(), msdf.Height()), nil
}

func Load2d(img image.Image, scale float32, reducer mesh2.Reducer, origin vec3, normal vec3) (mesh data.TriMesh, bounds image.Rectangle, err error) {
	msdf, err := loadMsdf(img, scale)
	if err != nil {
		return
	}

	normalLenInv := lenInv(normal)
	normal[0], normal[1], normal[2] = normal[0]*normalLenInv, normal[1]*normalLenInv, normal[2]*normalLenInv

	origin = vec3{origin[0] * float32(msdf.Width()), origin[1] * float32(msdf.Height()), 0}
	polygons := march.March(msdf, true)
	totalVertCount := reducePolygons(polygons, reducer)

	vertexBuffer := make([]vec3, totalVertCount)
	indexBuffer := make([]uint32, 0, totalVertCount*3)
	var normalBuffer []vec3
	if normal != [3]float32{} {
		normalBuffer = make([]vec3, totalVertCount)
		for i := range normalBuffer {
			normalBuffer[i] = normal
		}
	}
	var indexOffset uint32

	for _, poly := range polygons {
		if poly.ExteriorBound == nil {
			continue
		}
		indices, vertices, ok := mesh2.Poly2Tri(poly, indexOffset)
		if !ok {
			continue
		}

		for i, v := range vertices {
			vertexBuffer[int(indexOffset)+i] = vec3{v[0] - origin[0], v[1] - origin[1], 0}
		}
		indexBuffer = append(indexBuffer, indices...)
		indexOffset += uint32(len(vertices))
	}

	mesh = makeTriMesh(vertexBuffer, normalBuffer, indexBuffer)

	return mesh, image.Rect(0, 0, msdf.Width(), msdf.Height()), nil
}

func makeTriMesh(vertices, normals []vec3, indices []uint32) data.TriMesh {
	var (
		vbo = data.Buffer{Target: data.BufVertexAttribute}
		ibo = data.Buffer{Target: data.BufElementIndex}
		vao data.Vao
	)

	vao.Bind()
	vbo.Bind()

	vertexData := make([]vec3, len(vertices), len(vertices)+len(normals))
	copy(vertexData, vertices)
	vao.MustLayout(0, 3, float32(0), false, 0, 0)
	if normals != nil {
		vertexData = append(vertexData, normals...)
		vao.MustLayout(1, 3, float32(0), false, 0, 3*4*len(vertices))
	}
	vbo.WriteStatic(vertexData)

	ibo.Bind()
	ibo.WriteStatic(indices)

	vao.Unbind()
	ibo.Unbind()

	return data.TriMesh{
		Vao:         vao,
		Vbo:         vbo,
		Ibo:         ibo,
		IndexCount:  len(indices),
		VertexCount: len(vertices),
	}
}

// calculate normal of v1 in the xy plane
func calcVertexNormalXY(v0, v1, v2 vec3, flip bool) vec3 {
	n01 := vec2{v1[1] - v0[1], v0[0] - v1[0]}
	n02 := vec2{v2[1] - v1[1], v1[0] - v2[0]}
	n := vec3{n01[0] + n02[0], n01[1] + n02[1]}

	l := lenInv(n)
	if flip {
		l *= -1
	}
	n[0], n[1], n[2] = n[0]*l, n[1]*l, n[2]*l

	return n
}

func lenInv(v vec3) float32 {
	return 1 / float32(math.Sqrt(float64(v[0]*v[0]+v[1]*v[1]+v[2]*v[2])))
}

func reducePolygons(polygons []march.Polygon, reducer mesh2.Reducer) int {
	totalVertCount := 0
	for i := range polygons {
		poly := &polygons[i]
		bound := poly.ExteriorBound
		bound = RemoveCloseDuplicates(bound, Epsilon)
		bound = reducer.Reduce(bound)
		poly.ExteriorBound = bound
		totalVertCount += len(bound)

		for i, bound := range poly.InteriorBounds {
			bound = RemoveCloseDuplicates(bound, Epsilon)
			bound = reducer.Reduce(bound)
			poly.InteriorBounds[i] = bound
			totalVertCount += len(bound)
		}
	}
	return totalVertCount
}

func loadMsdf(img image.Image, scale float32) (*field.MSDF, error) {
	src := field.RGBAValueField(img)
	if scale <= 1 {
		return field.NewMSDF(src), nil
	}
	dst := field.NewValueField(int(float32(src.Width)*scale), int(float32(src.Height)*scale), 4)
	field.ScaleBlerpFull(src, dst)
	return field.NewMSDF(dst), nil
}

// Removes vertices that are closer than epsilon in-place
func RemoveCloseDuplicates(shape []vec2, epsilon float64) (result []vec2) {
	var (
		vA      = shape[0]
		vB      vec2
		removed = 0
	)

	for i := 1; i < len(shape); i++ {
		vB = shape[i]

		distX := math.Abs(float64(vA[0] - vB[0]))
		distY := math.Abs(float64(vA[1] - vB[1]))
		if distX < epsilon && distY < epsilon {
			removed++
		} else {
			shape[i-removed] = vB
			vA = vB
		}

	}

	vB = shape[0]
	// Check last and first
	distX := math.Abs(float64(vA[0] - vB[0]))
	distY := math.Abs(float64(vA[1] - vB[1]))
	if distX < epsilon && distY < epsilon {
		removed++
	}

	return shape[:len(shape)-removed]
}

func Reduce2(shape []vec2) (result []vec2) {
	return mesh2.RDP(shape, 3)
}

func Reduce(shape []vec2, th float32) (result []vec2) {
	result = make([]vec2, 0, len(shape))
	var (
		a = shape[0]
		b = shape[1]
		c vec2
	)

	for i := 2; i <= len(shape)+1; i++ {
		if i >= len(shape) {
			c = shape[i-len(shape)]
		} else {
			c = shape[i]
		}

		area := math.Abs(float64(a[0]*b[1]+b[0]*c[1]+c[0]*a[1]-a[1]*b[0]-b[1]*c[0]-c[1]*a[0])) / 2

		if float32(area) >= th {
			result = append(result, b)
			a = b
		}

		b = c
	}
	return result
}

// can result in self intersection if detail is too high
func Reduce3(shape []vec2, detail float32) (result []vec2) {
	if detail < Epsilon {
		detail = 0
	}

	result = make([]vec2, 0, len(shape))
	result = append(result, shape[0])
	var (
		a = shape[0]
		b = shape[1]
		c vec2
	)

	for i := 2; i <= len(shape); i++ {
		if i == len(shape) {
			c = shape[0]
		} else {
			c = shape[i]
		}

		distX := math.Abs(float64(a[0] - b[0]))
		distY := math.Abs(float64(a[1] - b[1]))
		if distX < Epsilon && distY < Epsilon {
			b = c
			continue
		}

		if detail == 0 {
			result = append(result, b)
			a = b
			b = c
			continue
		}

		alt := 0.0
		det := float64(a[0]*b[1]) + float64(b[0]*c[1]) + float64(c[0]*a[1]) - float64(a[1]*b[0]) - float64(b[1]*c[0]) - float64(c[1]*a[0])
		doubleArea := math.Abs(det)

		if math.Abs(det) > 1e-6 {
			base := math.Sqrt(float64((a[0]-c[0])*(a[0]-c[0]) + (a[1]-c[1])*(a[1]-c[1])))
			alt = doubleArea / base
		}

		if float32(alt) >= detail {
			result = append(result, b)
			a = b
		}

		b = c
	}

	if len(result) < 3 {
		return nil
	}

	return result
}
