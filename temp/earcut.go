package main

import (
	"math"
	"sort"
)

type float = float32

type node struct {
	prev    *node
	prevZ   *node
	next    *node
	nextZ   *node
	i       int
	x       float
	y       float
	z       float
	steiner bool
}

func newNode(i int, x, y float) *node {
	return &node{
		i: i,
		x: x,
		y: y,
	}
}

func Earcut(data []float, holeIndices []int, dim int) []int {
	var (
		hasHoles    bool
		outerLength int
		outerNode   *node
		triangles   []int
	)

	hasHoles = len(holeIndices) != 0
	if hasHoles {
		outerLength = holeIndices[0] * dim
	} else {
		outerLength = len(data)
	}
	outerNode = linkedList(data, 0, outerLength, dim, true)

	if outerNode == nil || outerNode.next == outerNode.prev {
		return triangles
	}

	var (
		minX    float
		minY    float
		maxX    float
		maxY    float
		x       float
		y       float
		invSize float
	)

	if hasHoles {
		outerNode = eliminateHoles(data, holeIndices, outerNode, dim)
	}

	// if the shape is not too simple, we'll use z-order curve hash later; calculate polygon bbox
	if len(data) > 80*dim {
		minX, maxX = data[0], data[0]
		minY, maxY = data[1], data[1]

		for i := dim; i < outerLength; i += dim {
			x = data[i]
			y = data[i+1]
			if x < minX {
				minX = x
			}
			if y < minY {
				minY = y
			}
			if x > maxX {
				maxX = x
			}
			if y > maxY {
				maxY = y
			}
		}

		// minX, minY and invSize are later used to transform coords into integers for z-order calculation
		invSize = float(math.Max(float64(maxX-minX), float64(maxY-minY)))
		if invSize != 0 {
			invSize = 1 / invSize
		} else {
			invSize = 0
		}
	}

	earcutLinked(outerNode, &triangles, dim, minX, minY, invSize, 0)
	return triangles
}

func linkedList(data []float, start int, end int, dim int, clockwise bool) *node {
	var last *node

	if clockwise == (signedArea(data, start, end, dim) > 0) {
		for i := start; i < end; i += dim {
			last = insertNode(i, data[i], data[i+1], last)
		}
	} else {
		for i := end - dim; i >= start; i -= dim {
			last = insertNode(i, data[i], data[i+1], last)
		}
	}

	if last != nil && equalsNode(last, last.next) {
		removeNode(last)
		last = last.next
	}

	return last
}

// main ear slicing loop which triangulates a polygon (given as a linked list)
func earcutLinked(ear *node, triangles *[]int, dim int, minX, minY, invSize float, pass int) {
	if ear == nil {
		return
	}

	// interlink polygon nodes in z-order
	if pass == 0 && invSize != 0 {
		indexCurve(ear, minX, minY, invSize)
	}

	stop := ear
	var prev, next *node

	// iterate through ears, slicing them one by one
	for ear.prev != ear.next {
		prev = ear.prev
		next = ear.next

		var ok bool
		if invSize != 0 {
			ok = isEarHashed(ear, minX, minY, invSize)
		} else {
			ok = isEar(ear)
		}
		if ok {
			// cut off the triangle
			*triangles = append(*triangles, prev.i/dim, ear.i/dim, next.i/dim)

			removeNode(ear)

			// skipping the next vertex leads to less sliver triangles
			ear = next.next
			stop = next.next

			continue
		}

		ear = next

		// if we looped through the whole remaining polygon and can't find any more ears
		if ear == stop {
			// try filtering points and slicing again
			if pass == 0 {
				earcutLinked(filterPoints(ear, nil), triangles, dim, minX, minY, invSize, 1)

				// if this didn't work, try curing all small self-intersections locally
			} else if pass == 1 {
				ear = cureLocalIntersections(filterPoints(ear, nil), triangles, dim)
				earcutLinked(ear, triangles, dim, minX, minY, invSize, 0)

				// as a last resort, try splitting the remaining polygon into two
			} else if pass == 2 {
				splitEarcut(ear, triangles, dim, minX, minY, invSize)
			}

			break
		}
	}
	return
}

// check whether a polygon node forms a valid ear with adjacent nodes
func isEar(ear *node) bool {
	a := ear.prev
	b := ear
	c := ear.next

	if area(a, b, c) >= 0 {
		return false // reflex, can't be an ear
	}

	// now make sure we don't have other points inside the potential ear
	p := ear.next.next

	for p != ear.prev {
		if pointInTriangle(a.x, a.y, b.x, b.y, c.x, c.y, p.x, p.y) &&
			area(p.prev, p, p.next) >= 0 {
			return false
		}
		p = p.next
	}

	return true
}

func min3(a, b, c float) float {
	return float(math.Min(math.Min(float64(a), float64(b)), float64(c)))
}

func max3(a, b, c float) float {
	return float(math.Max(math.Max(float64(a), float64(b)), float64(c)))
}

func isEarHashed(ear *node, minX, minY, invSize float) bool {
	a := ear.prev
	b := ear
	c := ear.next

	if area(a, b, c) >= 0 {
		return false // reflex, can't be an ear
	}

	minTX := min3(a.x, b.x, c.x)
	minTY := min3(a.y, b.y, c.y)
	maxTX := max3(a.x, b.x, c.x)
	maxTY := max3(a.y, b.y, c.y)

	// z-order range for the current triangle bbox;
	minZ := zOrder(minTX, minTY, minX, minY, invSize)
	maxZ := zOrder(maxTX, maxTY, minX, minY, invSize)

	p := ear.prevZ
	n := ear.nextZ

	// look for points inside the triangle in both directions
	for p != nil && p.z >= minZ && n != nil && n.z <= maxZ {
		if p != ear.prev && p != ear.next &&
			pointInTriangle(a.x, a.y, b.x, b.y, c.x, c.y, p.x, p.y) &&
			area(p.prev, p, p.next) >= 0 {
			return false
		}
		p = p.prevZ

		if n != ear.prev && n != ear.next &&
			pointInTriangle(a.x, a.y, b.x, b.y, c.x, c.y, n.x, n.y) &&
			area(n.prev, n, n.next) >= 0 {
			return false
		}
		n = n.nextZ
	}

	// look for remaining points in decreasing z-order
	for p != nil && p.z >= minZ {
		if p != ear.prev && p != ear.next &&
			pointInTriangle(a.x, a.y, b.x, b.y, c.x, c.y, p.x, p.y) &&
			area(p.prev, p, p.next) >= 0 {
			return false
		}
		p = p.prevZ
	}

	// look for remaining points in increasing z-order
	for n != nil && n.z <= maxZ {
		if n != ear.prev && n != ear.next &&
			pointInTriangle(a.x, a.y, b.x, b.y, c.x, c.y, n.x, n.y) &&
			area(n.prev, n, n.next) >= 0 {
			return false
		}
		n = n.nextZ
	}

	return true
}

// go through all polygon nodes and cure small local self-intersections
func cureLocalIntersections(start *node, triangles *[]int, dim int) *node {
	p := start
	for {
		a := p.prev
		b := p.next.next

		if !equalsNode(a, b) && intersects(a, p, p.next, b) && locallyInside(a, b) && locallyInside(b, a) {

			*triangles = append(*triangles, a.i/dim, p.i/dim, b.i/dim)

			// remove two nodes involved
			removeNode(p)
			removeNode(p.next)

			start = b
			p = b
		}
		p = p.next
		if p == nil || p == start {
			break
		}
	}

	return filterPoints(p, nil)
}

// try splitting polygon into two and triangulate them independently
func splitEarcut(start *node, triangles *[]int, dim int, minX, minY, invSize float) {
	// look for a valid diagonal that divides the polygon into two
	a := start
	for {
		b := a.next.next
		for b != a.prev {
			if a.i != b.i && isValidDiagonal(a, b) {
				// split the polygon in two by the diagonal
				c := splitPolygon(a, b)

				// filter colinear points around the cuts
				a = filterPoints(a, a.next)
				c = filterPoints(c, c.next)

				// run earcut on each half
				earcutLinked(a, triangles, dim, minX, minY, invSize, 0)
				earcutLinked(c, triangles, dim, minX, minY, invSize, 0)
				return
			}
			b = b.next
		}
		a = a.next
		if a == nil || a == start {
			break
		}
	}
	return
}

// eliminate colinear or duplicate points
func filterPoints(start, end *node) *node {
	if start == nil {
		return nil
	}
	if end == nil {
		end = start
	}

	p := start
	var again bool
	for {
		again = false

		if !p.steiner && (equalsNode(p, p.next) || area(p.prev, p, p.next) == 0) {
			removeNode(p)
			end = p.prev
			p = p.prev
			if p == p.next {
				break
			}
			again = true

		} else {
			p = p.next
		}
		if p == nil || !again && p == end {
			break
		}
	}

	return end
}

func eliminateHoles(data []float, holeIndices []int, outerNode *node, dim int) *node {
	var (
		queue  = []*node{}
		length = len(holeIndices)
		start  int
		end    int
		list   *node
	)

	for i := 0; i < length; i++ {
		start = holeIndices[0] * dim
		if i < length-1 {
			end = holeIndices[i+1] * dim
		} else {
			end = len(data)
		}
		list = linkedList(data, start, end, dim, false)
		if list == list.next {
			list.steiner = true
		}
		queue = append(queue, getLeftmost(list))
	}

	sort.Slice(queue, func(i, j int) bool {
		return queue[i].x < queue[j].x
	})

	for _, list := range queue {
		eliminateHole(list, outerNode)
		outerNode = filterPoints(outerNode, outerNode.next)
	}

	return outerNode
}

func eliminateHole(hole *node, outerNode *node) {
	outerNode = findHoleBridge(hole, outerNode)
	if outerNode == nil {
		return
	}
	bridge := splitPolygon(outerNode, hole)
	filterPoints(outerNode, outerNode.next)
	filterPoints(bridge, bridge.next)
}

func findHoleBridge(hole *node, outerNode *node) *node {
	p := outerNode
	hx := hole.x
	hy := hole.y
	qx := float(-math.MaxFloat32)
	var m *node

	for {
		if hy <= p.y && hy >= p.next.y && p.next.y != p.y {
			var x = p.x + (hy-p.y)*(p.next.x-p.x)/(p.next.y-p.y)
			if x <= hx && x > qx {
				qx = x
				if x == hx {
					if hy == p.y {
						return p
					}
					if hy == p.next.y {
						return p.next
					}
				}
				if p.x < p.next.x {
					m = p
				} else {
					m = p.next
				}
			}
		}
		p = p.next
		if p == nil || p == outerNode {
			break
		}
	}

	if m == nil {
		return nil
	}

	if hx == qx {
		return m
	}

	stop := m
	mx := m.x
	my := m.y
	tanMin := float(math.MaxFloat32)
	var tan float32

	p = m

	for {
		var (
			ax float
			cx float
		)
		if hy < my {
			ax = hx
			cx = qx
		} else {
			ax = qx
			cx = hx
		}
		if hx >= p.x && p.x >= mx && hx != p.x &&
			pointInTriangle(ax, hy, mx, my, cx, hy, p.x, p.y) {

			tan = float(math.Abs(float64(hy-p.y))) / (hx - p.x) // tangential

			if locallyInside(p, hole) &&
				(tan < tanMin || (tan == tanMin && (p.x > m.x || (p.x == m.x && sectorContainsSector(m, p))))) {
				m = p
				tanMin = tan
			}
		}

		p = p.next
		if p == nil || p == stop {
			break
		}
	}

	return m
}

func pointInTriangle(ax, ay, bx, by, cx, cy, px, py float) bool {
	return (cx-px)*(ay-py)-(ax-px)*(cy-py) >= 0 &&
		(ax-px)*(by-py)-(bx-px)*(ay-py) >= 0 &&
		(bx-px)*(cy-py)-(cx-px)*(by-py) >= 0
}

func locallyInside(a, b *node) bool {
	if area(a.prev, a, a.next) < 0 {
		return area(a, b, a.next) >= 0 && area(a, a.prev, b) >= 0
	}
	return area(a, b, a.prev) < 0 || area(a, a.next, b) < 0
}

func middleInside(a, b *node) bool {
	p := a
	inside := false
	px := (a.x + b.x) / 2
	py := (a.y + b.y) / 2
	for {
		if ((p.y > py) != (p.next.y > py)) && p.next.y != p.y &&
			(px < (p.next.x-p.x)*(py-p.y)/(p.next.y-p.y)+p.x) {
			inside = !inside
		}
		p = p.next
		if p == nil || p == a {
			break
		}
	}

	return inside
}

func splitPolygon(a, b *node) *node {
	a2 := newNode(a.i, a.x, a.y)
	b2 := newNode(b.i, b.x, b.y)
	an := a.next
	bp := b.prev

	a.next = b
	b.prev = a

	a2.next = an
	an.prev = a2

	b2.next = a2
	a2.prev = b2

	bp.next = b2
	b2.prev = bp

	return b2
}

// create a node and optionally link it with previous one (in a circular doubly linked list)
func insertNode(i int, x, y float, last *node) *node {
	p := newNode(i, x, y)

	if last == nil {
		p.prev = p
		p.next = p

	} else {
		p.next = last.next
		p.prev = last
		last.next.prev = p
		last.next = p
	}
	return p
}

func removeNode(p *node) {
	p.next.prev = p.prev
	p.prev.next = p.next

	if p.prevZ != nil {
		p.prevZ.nextZ = p.nextZ
	}
	if p.nextZ != nil {
		p.nextZ.prevZ = p.prevZ
	}
}

func sectorContainsSector(m *node, p *node) bool {
	return area(m.prev, m, p.prev) < 0 && area(p.next, m, m.next) < 0
}

func indexCurve(start *node, minX, minY float, invSize float) {
	p := start
	for {
		if p.z == 0 {
			p.z = zOrder(p.x, p.y, minX, minY, invSize)
		}
		p.prevZ = p.prev
		p.nextZ = p.next
		p = p.next
		if p == start {
			break
		}
	}

	p.prevZ.nextZ = nil
	p.prevZ = nil

	sortLinked(p)
}

// Simon Tatham's linked list merge sort algorithm
// http://www.chiark.greenend.org.uk/~sgtatham/algorithms/listsort.html
func sortLinked(list *node) *node {
	var (
		p         *node
		q         *node
		tail      *node
		e         *node
		pSize     int
		qSize     int
		numMerges int
	)
	inSize := 1

	for {
		p = list
		list = nil
		tail = nil
		numMerges = 0

		for p != nil {
			numMerges++
			q = p
			pSize = 0
			for i := 0; i < inSize; i++ {
				pSize++
				q = q.nextZ
				if q == nil {
					break
				}
			}
			qSize = inSize

			for pSize > 0 || (qSize > 0 && q != nil) {

				if pSize != 0 && (qSize == 0 || q == nil || p.z <= q.z) {
					e = p
					p = p.nextZ
					pSize--
				} else {
					e = q
					q = q.nextZ
					qSize--
				}

				if tail != nil {
					tail.nextZ = e
				} else {
					list = e
				}

				e.prevZ = tail
				tail = e
			}

			p = q
		}

		tail.nextZ = nil
		inSize *= 2

		if numMerges <= 1 {
			break
		}
	}

	return list
}

// z-order of a point given coords and inverse of the longer side of data bbox
func zOrder(x, y, minX, minY, invSize float) float {
	// coords are transformed into non-negative 15-bit integer range
	nx := uint(32767 * (x - minX) * invSize)
	ny := uint(32767 * (y - minY) * invSize)

	nx = (nx | (nx << 8)) & 0x00FF00FF
	nx = (nx | (nx << 4)) & 0x0F0F0F0F
	nx = (nx | (nx << 2)) & 0x33333333
	nx = (nx | (nx << 1)) & 0x55555555

	ny = (ny | (ny << 8)) & 0x00FF00FF
	ny = (ny | (ny << 4)) & 0x0F0F0F0F
	ny = (ny | (ny << 2)) & 0x33333333
	ny = (ny | (ny << 1)) & 0x55555555

	return float(nx | (ny << 1))
}

// check if a diagonal between two polygon nodes is valid (lies in polygon interior)
func isValidDiagonal(a *node, b *node) bool {
	return a.next.i != b.i && a.prev.i != b.i && !intersectsPolygon(a, b) && // dones't intersect other edges
		(locallyInside(a, b) && locallyInside(b, a) && middleInside(a, b) && // locally visible
			(area(a.prev, a, b.prev) != 0 || area(a, b.prev, b) != 0) || // does not create opposite-facing sectors
			equalsNode(a, b) && area(a.prev, a, a.next) > 0 && area(b.prev, b, b.next) > 0) // special zero-length case
}

func area(p *node, q *node, r *node) float {
	return (q.y-p.y)*(r.x-q.x) - (q.x-p.x)*(r.y-q.y)
}

func equalsNode(p1, p2 *node) bool {
	return p1.x == p2.x && p1.y == p2.y
}

// check if two segments intersect
func intersects(p1 *node, q1 *node, p2 *node, q2 *node) bool {
	o1 := sign(area(p1, q1, p2))
	o2 := sign(area(p1, q1, q2))
	o3 := sign(area(p2, q2, p1))
	o4 := sign(area(p2, q2, q1))

	if o1 != o2 && o3 != o4 {
		return true
	} // general case

	if o1 == 0 && onSegment(p1, p2, q1) {
		return true
	} // p1, q1 and p2 are collinear and p2 lies on p1q1
	if o2 == 0 && onSegment(p1, q2, q1) {
		return true
	} // p1, q1 and q2 are collinear and q2 lies on p1q1
	if o3 == 0 && onSegment(p2, p1, q2) {
		return true
	} // p2, q2 and p1 are collinear and p1 lies on p2q2
	if o4 == 0 && onSegment(p2, q1, q2) {
		return true
	} // p2, q2 and q1 are collinear and q1 lies on p2q2

	return false
}

// for collinear points p, q, r, check if point q lies on segment pr
func onSegment(p *node, q *node, r *node) bool {
	return float64(q.x) <= math.Max(float64(p.x), float64(r.x)) &&
		float64(q.x) >= math.Min(float64(p.x), float64(r.x)) &&
		float64(q.y) <= math.Max(float64(p.y), float64(r.y)) &&
		float64(q.y) >= math.Min(float64(p.y), float64(r.y))
}

func sign(num float) float {
	if num > 0 {
		return 1
	} else if num < 0 {
		return -1
	}
	return 0
}

// check if a polygon diagonal intersects any polygon segments
func intersectsPolygon(a *node, b *node) bool {
	p := a
	for {
		if p.i != a.i && p.next.i != a.i && p.i != b.i && p.next.i != b.i &&
			intersects(p, p.next, a, b) {
			return true
		}
		p = p.next
		if p == nil || p == a {
			break
		}
	}

	return false
}

func getLeftmost(start *node) *node {
	p := start
	leftmost := start
	for {
		if p.x < leftmost.x || (p.x == leftmost.x && p.y < leftmost.y) {
			leftmost = p
		}
		p = p.next
		if p == nil || p == start {
			break
		}
	}
	return leftmost
}

func signedArea(data []float, start int, end int, dim int) float {
	var sum float
	j := end - dim
	for i := start; i < end; i += dim {
		sum += (data[j] - data[i]) * (data[i+1] + data[j+1])
		j = i
	}
	return sum
}
