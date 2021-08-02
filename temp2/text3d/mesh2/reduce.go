package mesh2

import (
	"container/heap"
	"math"
)

func triangleDeterminant(a, b, c vertex) float64 {
	return float64(a[0]*b[1]) + float64(b[0]*c[1]) + float64(c[0]*a[1]) - float64(a[1]*b[0]) - float64(b[1]*c[0]) - float64(c[1]*a[0])
}

func triangleAltitude(doubleArea float64, a, c vertex) float64 {
	base := math.Sqrt(float64((a[0]-c[0])*(a[0]-c[0]) + (a[1]-c[1])*(a[1]-c[1])))
	return doubleArea / base
}

type Reducer interface {
	// Reduces vertices and returns a new slice or nil if len < 3
	Reduce(polygon []vertex) []vertex
}

// detail best between 0.1 to 0.01
type DetailReducer struct {
	Detail float32
}

func (r DetailReducer) Reduce(polygon []vertex) (result []vertex) {
	if r.Detail < 1e-6 || len(polygon) == 0 {
		return append([]vertex{}, polygon...)
	}

	result = make([]vertex, 0, len(polygon))
	result = append(result, polygon[0])
	var (
		a = polygon[0]
		b = polygon[1]
		c vertex
	)

	for i := 2; i <= len(polygon); i++ {
		if i == len(polygon) {
			c = polygon[0]
		} else {
			c = polygon[i]
		}

		alt := 0.0
		det := triangleDeterminant(a, b, c)
		doubleArea := math.Abs(det)

		if doubleArea > 1e-6 {
			alt = triangleAltitude(doubleArea, a, c)
		}

		if float32(alt) >= r.Detail {
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

type linkedTri struct {
	A, C      *linkedTri
	B         vertex
	index     int
	memoValue float64
	unchanged bool
}

func (v *linkedTri) Value() float64 {
	if v.unchanged {
		return v.memoValue
	}

	det := triangleDeterminant(v.A.B, v.B, v.C.B)
	doubleArea := math.Abs(det)
	alt := 0.0
	if doubleArea > 1e-6 {
		alt = triangleAltitude(doubleArea, v.A.B, v.C.B)
	}
	v.memoValue = alt
	v.unchanged = true
	return alt
}

// triHeightQueue implements a priority queue by triangle height
// https://pkg.go.dev/container/heap@go1.16.6#example-package-PriorityQueue
type triHeightQueue []*linkedTri

func (q triHeightQueue) Len() int           { return len(q) }
func (q triHeightQueue) Less(i, j int) bool { return q[i].Value() < q[j].Value() }
func (q triHeightQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *triHeightQueue) Push(x interface{}) {
	n := len(*q)
	item := x.(*linkedTri)
	item.index = n
	*q = append(*q, item)
}

func (q *triHeightQueue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*q = old[0 : n-1]
	return item
}

// CountReducer removes the least important vertices until at most 'Count' are left
type CountReducer struct {
	Count int
}

func (r CountReducer) Reduce(polygon []vertex) (result []vertex) {
	if len(polygon) <= r.Count || len(polygon) == 0 {
		return append([]vertex{}, polygon...)
	}

	q := make(triHeightQueue, 0, len(polygon))

	var (
		a = &linkedTri{B: polygon[0], index: 0}
		b = &linkedTri{A: a, B: polygon[1], index: 1}
		c *linkedTri
	)
	first := a
	a.C = b

	for i := 2; i <= len(polygon); i++ {
		if i == len(polygon) {
			first.A = c
			c = first
		} else {
			c = &linkedTri{B: polygon[i], index: i}
		}

		b.C = c
		c.A = b

		heap.Push(&q, b)

		a = b
		b = c
	}

	heap.Init(&q)

	over := len(polygon) - r.Count
	for i := 0; i < over; i++ {
		tri := heap.Pop(&q).(*linkedTri)
		tri.A.unchanged = false
		tri.A.C = tri.C
		tri.C.unchanged = false
		tri.C.A = tri.A
		heap.Fix(&q, tri.A.index)
		heap.Fix(&q, tri.C.index)
	}

	result = make([][2]float32, 0, len(q))
	result = append(result, q[0].B)

	for v := q[0].C; v != q[0]; v = v.C {
		result = append(result, v.B)
	}

	if len(result) < 3 {
		return nil
	}

	return result
}
