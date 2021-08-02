package march

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/Qendolin/go-printpixel/temp2/text3d/march/field"
)

// Polyline is a polygonal chain
type Polyline = []vertex

// Polygon is a simple polygon with optional holes
type Polygon struct {
	ExteriorBound  Polyline
	InteriorBounds []Polyline
}

// 	Edges  │ Corners │ Indices
// ────────┼─────────┼─────────
// 	┌─4─┐  │  8───4  │  3───2
// 	3   2  │  │   │  │  │   │
// 	└─1─┘  │  2───1  │  1───0
var (
	edge1 = [2]float32{0.5, 1}
	edge2 = [2]float32{1, 0.5}
	edge3 = [2]float32{0, 0.5}
	edge4 = [2]float32{0.5, 0}
)

// Ordering the edges from CCW to CW makes deduplication easier
// since they will from a loop that always goes CW.
// For a better explaination on ordering see stepTable.
// Every second edge in this table is actually unnecessary.
var contourTable = [][][2]float32{
	{},                           // 0
	{edge2, edge1},               // 1
	{edge1, edge3},               // 2
	{edge2, edge3},               // 3
	{edge4, edge2},               // 4
	{edge4, edge1},               // 5
	{edge1, edge3, edge4, edge2}, // 6 up -> first pair, down -> second pair
	{edge4, edge3},               // 7
	{edge3, edge4},               // 8
	{edge3, edge4, edge2, edge1}, // 9 right -> first pair, left -> second pair
	{edge1, edge4},               // 10
	{edge2, edge4},               // 11
	{edge3, edge2},               // 12
	{edge3, edge1},               // 13
	{edge1, edge2},               // 14
	{},                           // 15
}

// Corner indices for each edge
// TODO: Figure out pattern
var cornerIndexTable = [][]int{
	{}, // 0
	// {0, 2, 0, 1},             // 1
	{2, 0, 1, 0}, // 1
	// {1, 0, 1, 3}, // 2
	{1, 0, 3, 1}, // 2
	// {0, 2, 1, 3}, // 3
	{2, 0, 3, 1}, // 3
	// {2, 3, 2, 0}, // 4
	{3, 2, 2, 0}, // 4
	// {2, 3, 0, 1},             // 5
	{3, 2, 1, 0},             // 5
	{1, 0, 1, 3, 2, 3, 2, 0}, // 6 down -> first pairs, up -> second pairs
	// {2, 3, 1, 3},             // 7
	{3, 2, 3, 1},             // 7
	{3, 1, 3, 2},             // 8
	{3, 1, 3, 2, 0, 2, 0, 1}, // 9 left -> first pair, right -> second pair
	{1, 0, 3, 2},             // 10
	// {0, 2, 3, 2}, // 11 ???
	{2, 0, 3, 2}, // 11
	{3, 1, 2, 0}, // 12
	// {3, 1, 0, 1}, // 13
	{3, 1, 1, 0}, // 13
	{1, 0, 2, 0}, // 14
	{},           // 15
}

type step [2]int

var (
	none  = step{0, 0}
	up    = step{0, -1}
	down  = step{0, +1}
	left  = step{-1, 0}
	right = step{+1, 0}
)

// Consider all corners that are inside the shape
// Go clockwise and take the last
// Then the next (cw) crossed edge will be the next direction
// ········│ Example case 14:
// ·8───4··│  Corners 2, 4 and 8 are inside. The shape crosses the edges 2 and 1.
// ·│·┌─┼──┘  Starting at corner 2 and going clockwise we find that 4 is the last corner.
// ·2─┼─1     The next edge after 4 is 2 and pointing to the right side.
// ───┘       So the next direction is for case 14 is 'right'.
// As case 6 and 9 are ambiguous the direction can't be precomputed
// and are decided at runtime
var stepTable = []step{
	none,  // 0
	down,  // 1
	left,  // 2
	left,  // 3
	right, // 4
	down,  // 5
	none,  // 6 up -> right, down -> left
	left,  // 7
	up,    // 8
	none,  // 9 right -> down, left -> up
	up,    // 10
	up,    // 11
	right, // 12
	down,  // 13
	right, // 14
	none,  // 15
}

var _debugImage *image.RGBA

const debug = false

type point [2]int
type vertex = [2]float32

// FIXME: When a border is only 1 px wide, │ E.g.: ··
// sometimes holes will be missed.         │       · ·
//                                         │        ··

func March(f field.ScalarField, flipY bool) []Polygon {
	_debugImage = nil
	if debug {
		resetDebugImages()
	}

	width, height := f.Width(), f.Height()
	values := f.Raw()

	polygons := []Polygon{}
	boundId := 1
	polygonIndices := map[int]int{}

	// Could use a map here but is much slower
	boundIds := make([]int, width*height)
	markVisit := func(x, y, _ int) {
		boundIds[y*width+x] = boundId
	}

	for y := 0; y < height; y++ {
		last := false
		outerBoundId := -1
		for x := 0; x < width; x++ {
			i := y*width + x
			curr := values[i] >= 0.5
			if last || !curr {
				last = curr
				continue
			}
			last = curr

			if boundId := boundIds[i]; boundId != 0 {
				if _, isOuter := polygonIndices[boundId]; isOuter {
					outerBoundId = boundId
				}
				continue
			}

			if debug {
				createDebugImage(f, []point{{x, y}})
			}

			bound, isInnerBound := trace(values, width, height, [2]int{x, y}, flipY, markVisit)
			if isInnerBound {
				poly := &polygons[polygonIndices[outerBoundId]]
				poly.InteriorBounds = append(poly.InteriorBounds, bound)
			} else {
				polygonIndices[boundId] = len(polygons)
				polygons = append(polygons, Polygon{ExteriorBound: bound})
			}

			if debug {
				if isInnerBound {
					saveDebugImage("hole", boundId)
				} else {
					saveDebugImage("shape", boundId)
				}
			}
			boundId++
		}
	}

	return polygons
}

func resetDebugImages() {
	os.RemoveAll("./debug/")
}

func saveDebugImage(typ string, nr int) {
	if _debugImage != nil {
		os.MkdirAll("./debug/", 0666)
		debug, _ := os.Create(fmt.Sprintf("./debug/%s_%03d.png", typ, nr))
		png.Encode(debug, _debugImage)
	}
}

func createDebugImage(f field.ScalarField, starts []point) {
	_debugImage = image.NewRGBA(image.Rect(0, 0, f.Width(), f.Height()))
	for x := 0; x < f.Width(); x++ {
		for y := 0; y < f.Height(); y++ {
			if f.Get(x, y) >= 0.5 {
				_debugImage.Set(x, y, color.RGBA{127, 127, 127, 255})
			} else {
				_debugImage.Set(x, y, color.Black)
			}
		}
	}

	for _, p := range starts {
		_debugImage.Set(p[0], p[1], color.RGBA{0, 255, 0, 255})
	}
}

func trace(values []float32, width, height int, start point, flipY bool, onVisit func(x, y, value int)) (verts []vertex, isHole bool) {
	x, y := start[0], start[1]
	prevStep := none
	ySign := float32(1)
	if flipY {
		ySign = -1
	}
	if onVisit == nil {
		onVisit = func(x, y, value int) {}
	}

	outTurns := 0
	inTurns := 0

	// Failsafe to prevent endless loop
	// Cannot loop more times than number of pixels
	failsafeCount := 0
	failsafeMax := width * height

	for failsafeCount <= failsafeMax {
		value, next, values := calcNextStep(values, width, height, x, y, prevStep)
		_ = values
		onVisit(x, y, value)
		if next == none {
			// Should never happen
			// But a panic is better than an endless loop
			panic("Next step is 'none', how did this happen?")
		}

		if value == 1 || value == 2 || value == 4 || value == 8 {
			outTurns++
		} else if value == 7 || value == 11 || value == 13 || value == 14 {
			inTurns++
		} else if value == 6 || value == 9 {
			// TODO: This could work since a start is always at the top left
			// but it might not. In a simple loops the outer order might be
			// 6 (down), 9 (right), 6 (up), 9 (left), so ccw
			// and the inner order might be
			// 6 (down), 9 (left), 6 (up), 9 (right), so cw.
			// This may not be the case in general (!) just something that I noticed.
			inTurns++
		}

		if value == 6 || value == 9 {
			if prevStep == up || prevStep == right || (prevStep == none && value == 6) {
				vert := contourTable[value][0]
				corners := cornerIndexTable[value][4:]
				vx, vy := calcFracVert(vert, corners, values)
				verts = append(verts, [2]float32{float32(x) + vx, ySign * (float32(y) + vy)})
			} else if prevStep == down || prevStep == left || (prevStep == none && value == 9) {
				vert := contourTable[value][2]
				corners := cornerIndexTable[value][:4]
				vx, vy := calcFracVert(vert, corners, values)
				verts = append(verts, [2]float32{float32(x) + vx, ySign * (float32(y) + vy)})
			} else {
				panic("Invalid march state")
			}
		} else {
			// The second edge can be skipped as it would just be a duplicate
			vert := contourTable[value][0]
			corners := cornerIndexTable[value]
			vx, vy := calcFracVert(vert, corners, values)
			verts = append(verts, [2]float32{float32(x) + vx, ySign * (float32(y) + vy)})
		}

		x += next[0]
		y += next[1]

		if x == start[0] && y == start[1] {
			break
		}

		if debug {
			r, g, b, _ := _debugImage.At(x, y).RGBA()
			_debugImage.Set(x, y, color.RGBA{uint8(r + 128), uint8(g), uint8(b), 255})
		}
		prevStep = next

		failsafeCount++
	}

	if failsafeCount == failsafeMax {
		// Should never happen
		panic("Endless loop detected, maximum iterations exceeded.")
	}

	return verts, inTurns > outTurns
}

func calcFracVert(vert vertex, corners []int, values [4]float32) (vx, vy float32) {
	fract := lerp(values[corners[0]], values[corners[1]])

	vx, vy = vert[0], vert[1]
	if vx == 0.5 {
		vx = fract
	} else {
		vy = fract
	}
	return vx, vy
}

// Corners:
// 8───4
// │   │
// 2───1
// Where 1 is the current position
func calcNextStep(f []float32, w, h, x, y int, prevStep step) (value int, nextStep step, values [4]float32) {
	i := y*w + x
	if x < w && y < h {
		v := f[i]
		values[0] = v
		if v >= 0.5 {
			value |= 1
		}
	}

	if y < h && x-1 >= 0 {
		v := f[i-1]
		values[1] = v
		if v >= 0.5 {
			value |= 2
		}
	}

	if x < w && y-1 >= 0 {
		v := f[i-w]
		values[2] = v
		if v >= 0.5 {
			value |= 4
		}
	}

	if x-1 >= 0 && y-1 >= 0 {
		v := f[i-1-w]
		values[3] = v
		if v >= 0.5 {
			value |= 8
		}
	}

	if value == 6 {
		if prevStep == up {
			nextStep = right
		} else if prevStep == down {
			nextStep = left
		} else {
			// Bad march state, probably becase start is a 6
			// Just assume that prevStep is down
			nextStep = left
		}
	} else if value == 9 {
		if prevStep == right {
			nextStep = down
		} else if prevStep == left {
			nextStep = up
		} else {
			// Bad march state, probably becase start is a 9
			// Just assume that prevStep is right
			nextStep = down
		}
	} else {
		nextStep = stepTable[value]
	}

	return
}

func lerp(v1, v2 float32) float32 {
	// if the values are the same, it doesn't matter, so just return 0
	if math.Abs(float64(v1-v2)) < 0.00001 {
		return 0
	}

	// the delta interpolation value is equal to the difference between the threshold from first
	// value over the difference of value 2 from value 1
	return (0.5 - v1) / (v2 - v1)
}
