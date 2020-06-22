// An implementation of Conway's Game of Life.
// https://golang.org/doc/play/life.go
// Modified
package main

import (
	"math/rand"
	"sync"
)

const TileSize = 64

// Field represents a two-dimensional field of cells.
type Field struct {
	s    [][]bool
	w, h int
}

// NewField returns an empty field of the specified width and height.
func NewField(w, h int) *Field {
	s := make([][]bool, h)
	for i := range s {
		s[i] = make([]bool, w)
	}
	return &Field{s: s, w: w, h: h}
}

// Set sets the state of the specified cell to the given value.
func (f *Field) Set(x, y int, b bool) {
	f.s[y][x] = b
}

// Alive reports whether the specified cell is alive.
// If the x or y coordinates are outside the field boundaries they are wrapped
// toroidally. For instance, an x value of -1 is treated as width-1.
func (f *Field) Alive(x, y int) bool {
	x += f.w
	x %= f.w
	y += f.h
	y %= f.h
	return f.s[y][x]
}

// Next returns the state of the specified cell at the next time step.
func (f *Field) Next(x, y int) bool {
	// Count the adjacent cells that are alive.
	alive := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (j != 0 || i != 0) && f.Alive(x+i, y+j) {
				alive++
			}
		}
	}
	// Return next state according to the game rules:
	//   exactly 3 neighbors: on,
	//   exactly 2 neighbors: maintain current state,
	//   otherwise: off.
	return alive == 3 || alive == 2 && f.Alive(x, y)
}

func Work(wg *sync.WaitGroup, trigger chan bool, l *Life, sx, sy int) {
	for range trigger {
		// Update the state of the next field (b) from the current field (a).
		for y := sx; y < sx+TileSize; y++ {
			for x := sy; x < sy+TileSize; x++ {
				l.b.Set(x, y, l.a.Next(x, y))
			}
		}
		wg.Done()
	}
}

// Life stores the state of a round of Conway's Game of Life.
type Life struct {
	a, b *Field
	w, h int
	p    []chan bool
	wg   *sync.WaitGroup
	tb   []byte
}

// NewLife returns a new Life game state with a random initial state.
func NewLife(w, h int) *Life {
	a := NewField(w, h)
	for i := 0; i < (w * h / 4); i++ {
		a.Set(rand.Intn(w), rand.Intn(h), true)
	}

	l := new(Life)

	var wg sync.WaitGroup
	pool := make([]chan bool, 0)
	for x := 0; x < w; x += TileSize {
		for y := 0; y < h; y += TileSize {
			ch := make(chan bool, 1)
			go Work(&wg, ch, l, x, y)
			pool = append(pool, ch)
		}
	}

	*l = Life{
		a: a, b: NewField(w, h),
		w: w, h: h,
		p:  pool,
		wg: &wg,
	}
	return l
}

// Step advances the game by one instant, recomputing and updating all cells.
func (l *Life) Step() {
	l.wg.Add(len(l.p))
	for _, worker := range l.p {
		worker <- true
	}
	l.wg.Wait()
	// Swap fields a and b.
	l.a, l.b = l.b, l.a
}

func (l *Life) Texture() []byte {
	texBuf := l.tb
	if texBuf == nil {
		texBuf = make([]byte, l.w*l.h*3)
		l.tb = texBuf
	}
	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			if l.a.s[x][y] {
				texBuf[(x+y*l.w)*3+0] = 255
				texBuf[(x+y*l.w)*3+1] = 255
				texBuf[(x+y*l.w)*3+2] = 255
			} else {
				texBuf[(x+y*l.w)*3+0] = 0
				texBuf[(x+y*l.w)*3+1] = 0
				texBuf[(x+y*l.w)*3+2] = 0
			}
		}
	}
	return texBuf
}
