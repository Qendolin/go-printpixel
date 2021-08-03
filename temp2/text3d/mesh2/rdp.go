package mesh2

import (
	"math"
)

type vertex = [2]float32

func findPerpendicularDistance(p, p1, p2 vertex) float32 {
	if p1[0] == p2[0] {
		dist := math.Abs(float64(p[0] - p1[0]))
		return float32(dist)
	}
	slope := float64(p2[1]-p1[1]) / float64(p2[0]-p1[0])
	intercept := float64(p1[1]) - (slope * float64(p1[0]))
	dist := math.Abs(slope*float64(p[0])-float64(p[0])+intercept) / math.Sqrt(slope*slope+1)
	return float32(dist)
}

// A simple Ramer–Douglas–Peucker implementation
//
// Adapted from https://github.com/zx9597446/rdp
func RDP(points []vertex, epsilon float32) []vertex {
	if len(points) < 3 {
		return points
	}
	var (
		firstPoint = points[0]
		lastPoint  = points[len(points)-1]
		index      = -1
		dist       float32
	)
	for i := 1; i < len(points)-1; i++ {
		cDist := findPerpendicularDistance(points[i], firstPoint, lastPoint)
		if cDist > dist {
			dist = cDist
			index = i
		}
	}
	if dist > epsilon {
		l1 := points[0 : index+1]
		l2 := points[index:]
		r1 := RDP(l1, epsilon)
		r2 := RDP(l2, epsilon)
		rs := append(r1[0:len(r1)-1], r2...)
		return rs
	} else {
		ret := make([]vertex, 0, 2)
		ret = append(ret, firstPoint, lastPoint)
		return ret
	}
}
