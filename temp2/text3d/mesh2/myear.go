package mesh2

// https://wiki.delphigl.com/index.php/Ear_Clipping_Triangulierung

type point struct {
	x, y float32
}

type triangle = [3]point

func pointInTriangle(p, tp1, tp2, tp3 point) bool {
	b0 := ((tp2.x-tp1.x)*(tp3.y-tp1.y) - (tp3.x-tp1.x)*(tp2.y-tp1.y))
	if b0 != 0 {
		b1 := (((tp2.x-p.x)*(tp3.y-p.y) - (tp3.x-p.x)*(tp2.y-p.y)) / b0)
		b2 := (((tp3.x-p.x)*(tp1.y-p.y) - (tp1.x-p.x)*(tp3.y-p.y)) / b0)
		b3 := 1 - b1 - b2

		return (b1 > 0) && (b2 > 0) && (b3 > 0)
	}

	return false
}

func getVert(list []point, i int) point {
	i = i % len(list)
	if i < 0 {
		i = len(list) + i
	}
	return list[i]
}

func Earcut(verts [][2]float32) [][2]float32 {
	poly := make([]point, len(verts))
	tris := make([]triangle, 0, len(verts)/3)

	for i, v := range verts {
		poly[i] = point{v[0], v[1]}
	}

	makeCcw(poly)

	i := -1
	lastEar := -1

	// Abbrechen, wenn nach zwei ganzen Durchläufen keine Ecke gefunden wurde, oder nur noch
	// drei Ecken übrig sind.
	for lastEar <= len(poly)*2 && len(poly) > 3 {
		lastEar++
		i++
		if i >= len(poly) {
			i = 0
		}

		// Suche drei benachbarte Punkte aus der Liste
		p1 := getVert(poly, i-1)
		p := getVert(poly, i)
		p2 := getVert(poly, i+1)

		// Berechne, ob die Ecke konvex oder konkav ist
		l := ((p1.x-p.x)*(p2.y-p.y) - (p1.y-p.y)*(p2.x-p.x))

		// Nur weitermachen, wenn die Ecke konvex ist
		if l >= 0 {
			continue
		}

		// Überprüfe ob irgendein anderer Punkt aus dem Polygon
		// das ausgewählte Dreieck schneidet
		inTriangle := false
		for j := 2; j < len(poly)-1; j++ {
			pt := getVert(poly, i+j)

			if pointInTriangle(pt, p1, p, p2) {
				inTriangle = true
				break
			}
		}

		if inTriangle {
			continue
		}

		// Ist dies nicht der Fall, so entferne die ausgwewählte Ecke und bilde
		// ein neues Dreieck
		tri := triangle{{p1.x, p1.y}, {p.x, p.y}, {p2.x, p2.y}}
		tris = append(tris, tri)

		poly = append(poly[:i], poly[i+1:]...)

		lastEar = 0
		i--
	}

	if len(poly) == 3 {
		p1 := poly[0]
		p := poly[1]
		p2 := poly[2]

		tri := triangle{{p1.x, p1.y}, {p.x, p.y}, {p2.x, p2.y}}
		tris = append(tris, tri)
	}

	result := make([][2]float32, len(tris)*3)
	for i, tri := range tris {
		result[i*3][0] = tri[0].x
		result[i*3][1] = tri[0].y
		result[i*3+1][0] = tri[1].x
		result[i*3+1][1] = tri[1].y
		result[i*3+2][0] = tri[2].x
		result[i*3+2][1] = tri[2].y
	}
	return result
}

// https://en.wikipedia.org/wiki/Curve_orientation#Orientation_of_a_simple_polygon
func makeCcw(poly []point) {
	if isCw(poly) {
		for i, j := 0, len(poly)-1; i < j; i, j = i+1, j-1 {
			poly[i], poly[j] = poly[j], poly[i]
		}
	}
}

// https://stackoverflow.com/a/10298685/7448536
func isCw(poly []point) bool {
	var signedArea float32
	for i, p1 := range poly {
		var p2 point
		if i == len(poly)-1 {
			p2 = poly[0]
		} else {
			p2 = poly[i+1]
		}

		signedArea += (p1.x*p2.y - p2.x*p1.y)
	}
	return signedArea < 0
}
