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

// Reducer performs a polygon/polyline reduction algorithm.
type Reducer interface {
	// Reduces vertices and returns a new slice or nil if len < 3
	Reduce(polyline []vertex) []vertex
}

// DetailReducer performs a custom algorithm that is similar to Visvalingam–Whyatt
//
// Detail gives good results btween 0.1 and 0.01
type DetailReducer struct {
	Detail float32
}

func (r DetailReducer) Reduce(polyline []vertex) (result []vertex) {
	if r.Detail < 1e-6 || len(polyline) == 0 {
		return append([]vertex{}, polyline...)
	}

	result = make([]vertex, 0, len(polyline))
	result = append(result, polyline[0])
	var (
		a = polyline[0]
		b = polyline[1]
		c vertex
	)

	for i := 2; i <= len(polyline); i++ {
		if i == len(polyline) {
			c = polyline[0]
		} else {
			c = polyline[i]
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

// CountReducer performs the same algorithm as DetailReducer but
// removes the least important vertices until at most 'Count' are left
type CountReducer struct {
	Count int
}

func (r CountReducer) Reduce(polyline []vertex) (result []vertex) {
	return reduceToCount(polyline, r.Count)
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

func reduceToCount(polyline []vertex, count int) (result []vertex) {
	if len(polyline) <= count || len(polyline) == 0 {
		return append([]vertex{}, polyline...)
	}

	q := make(triHeightQueue, 0, len(polyline))

	var (
		a = &linkedTri{B: polyline[0], index: 0}
		b = &linkedTri{A: a, B: polyline[1], index: 1}
		c *linkedTri
	)
	first := a
	a.C = b

	for i := 2; i <= len(polyline); i++ {
		if i == len(polyline) {
			first.A = c
			c = first
		} else {
			c = &linkedTri{B: polyline[i], index: i}
		}

		b.C = c
		c.A = b

		heap.Push(&q, b)

		a = b
		b = c
	}

	heap.Init(&q)

	over := len(polyline) - count
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

// PercentReducer performs the same algorithm as DetailReducer but
// removes the least important vertices until at most ('Percent'*100)% are left
type PercentReducer struct {
	Percent float32
}

func (r PercentReducer) Reduce(polyline []vertex) (result []vertex) {
	return reduceToCount(polyline, int(r.Percent*float32(len(polyline))))
}

// RDPReducer performs the Ramer–Douglas–Peucker algorithm
type RDPReducer struct {
	Epsilon float32
}

func (r RDPReducer) Reduce(polyline []vertex) (result []vertex) {
	if len(polyline) <= 3 {
		return append([]vertex{}, polyline...)
	}

	return RDP(polyline, r.Epsilon)
}
