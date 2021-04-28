package main

import (
	"image"
	"image/color"
	"log"
	"math"
)

type Point [2]float32

func (p Point) Dot(q Point) float32 {
	return p[0] * q[0] + p[1] * q[1]
}

func (p Point) SquareLength() float32 {
	return p[0] * p[0] + p[1] * p[1]
}

func (p Point) Length() float32 {
	return float32(math.Sqrt(float64(p[0] * p[0] + p[1] * p[1])))
}

func (p Point) Sub(q Point) Point {
	return Point{p[0] - q[0], p[1] - q[1]}
}

type side int
const (
	none = side(0)
	up = side(1<<iota)
	right
	down
	left
)

type Marcher struct {
	Img   *image.RGBA
	Quality int // Quality is used for binary search interations
	Discriminator func(color.Color) float64
}

func (m Marcher) Process() [][]Point {
	loops := make([][]Point, 0)
	x, y := m.findStart()
	loops = append(loops, m.walk(x, y))
	return loops
}

func (m Marcher) findStart() (x, y int) {
	b := m.Img.Bounds()
	for y := b.Max.Y-1; y >= b.Min.Y; y-- {
		for x := b.Max.X-1; x >= b.Min.X; x-- {
			if m.Discriminator(m.Img.At(x, y)) >= 0.5 {
				return x, y
			}
		}
	}
	return 0, 0
}

func (m Marcher) lookup(value int, c8, c4, c2, c1 color.Color, previousStep side) (colors [][2]color.Color, orients []float64, sides []side, nextStep side) {
	var dir side
	switch value {
	default:
		return [][2]color.Color{}, []float64{}, []side{}, none
	case 1:
		return [][2]color.Color{{c1,c2}, {c1,c8}}, []float64{1,-1}, []side{down, left}, left
	case 2:
		return [][2]color.Color{{c2,c1}, {c2,c4}}, []float64{-1,-1}, []side{down,right}, down
	case 3:
		return [][2]color.Color{{c1,c8}, {c2,c4}}, []float64{-1,-1}, []side{left,right}, left
	case 4:
		return [][2]color.Color{{c4,c2}, {c4,c8}}, []float64{1,-1}, []side{right,up}, right
	case 5:
		if previousStep == up {
			dir = left
		} else {
			dir = right
		}
		return [][2]color.Color{{c1,c2}, {c1,c8}, {c4,c2}, {c4,c8}}, []float64{1,-1,1,-1}, []side{down,left,right,up}, dir
	case 6:
		return [][2]color.Color{{c2,c1}, {c4,c8}}, []float64{-1,-1}, []side{down,up}, down
	case 7:
		return [][2]color.Color{{c1,c8}, {c4,c8}}, []float64{-1,-1}, []side{left,up}, left
	case 8:
		return [][2]color.Color{{c8,c1}, {c8,c4}}, []float64{1,1}, []side{left,up}, up
	case 9:
		return [][2]color.Color{{c1,c2}, {c8,c4}}, []float64{1,1}, []side{down,up}, up
	case 10:
		if previousStep == right {
			dir = up
		} else {
			dir = down
		}
		return [][2]color.Color{{c2,c1}, {c2,c4}, {c8,c1}, {c8,c4}}, []float64{-1,-1,1,1}, []side{down,right,left,up}, dir
	case 11:
		return [][2]color.Color{{c2,c4}, {c8,c4}}, []float64{-1,1}, []side{right,up}, up
	case 12:
		return [][2]color.Color{{c4,c2}, {c8,c1}}, []float64{1,1}, []side{right,left}, right
	case 13:
		return [][2]color.Color{{c1,c2}, {c4,c2}}, []float64{1,1}, []side{down,right}, right
	case 14:
		return [][2]color.Color{{c2,c1}, {c8,c1}}, []float64{-1,1}, []side{down,left}, down
	}
}

func (m Marcher) step(loop *[]Point, previousStep side, x, y int) side {
	value := 0
	c8,c4,c2,c1 := m.Img.At(x, y), m.Img.At(x+1, y), m.Img.At(x+1, y+1), m.Img.At(x, y+1)
	if m.Discriminator(c8) >= 0.5 {
		value |= 8
	}
	if m.Discriminator(c4) >= 0.5{
		value |= 4
	}
	if m.Discriminator(c2) >= 0.5 {
		value |= 2
	}
	if m.Discriminator(c1) >= 0.5 {
		value |= 1
	}

	colors, orients, sides, nextStep := m.lookup(value, c8,c4,c2,c1, previousStep)

	for i := 0; i < len(colors); i++ {
		if i % 2 == 1 {
			continue
		}
		pair := colors[i]
		side := sides[i]
		value := float32(m.findIntersection(pair[0], pair[1], orients[i]))

		var p Point
		switch side {
		case up:
			p = Point{value, 0}
		case right:
			p = Point{1, value}
		case down:
			p = Point{value, 1}
		case left:
			p = Point{0, value}
		case none:
			p = Point{0.5, 0.5}
		}
		p[0] += float32(x)
		p[1] += float32(y)
		*loop = append(*loop, p)
	}

	return nextStep
}

func (m Marcher) walk(startX, startY int) []Point {
	loop := make([]Point, 0)
	step := down
	x, y := startX, startY
	i := 0
	for {
		i++
		nextStep := m.step(&loop, step, x, y)
		switch nextStep {
		case up:
			y--
		case left:
			x--
		case down:
			y++
		case right:
			x++
		case none:
			log.Println("NextStep is none, cannot continue")
			return loop
		}
		if x == startX && y == startY {
			break
		}
	}
	log.Printf("%v iters", i)
	return loop
}

func (m Marcher) findIntersection(s, e color.Color, orientation float64) float64 {
	t := 0.5
	epsilon := 1 / float64(int64(8 << m.Quality))

	for i := 1; i <= m.Quality; i++ {
		step := 1 / float64(int64(2 << i))
		v := m.Discriminator(clerp(s, e, t))
		d := math.Abs(0.5 - v)
		if d < epsilon {
			break
		} else if (v > 0.5) {
			t += step
		} else {
			t -= step
		}
	}
	if t < 0 || t > 1 {
		log.Panicf("t out of bounds %v", t)
	}

	if orientation < 0 {
		t = 1-t
	}

	return t
}