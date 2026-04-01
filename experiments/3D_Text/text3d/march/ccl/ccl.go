// package ccl is a port of https://github.com/spwhitt/cclabel
package ccl

import "github.com/Qendolin/go-printpixel/experiments/3D_Text/text3d/march/field"

type LabelField struct {
	Width, Height int
	Labels        []int
}

func (lf *LabelField) Get(x, y int) int {
	return lf.Labels[x+y*lf.Width]
}

func (lf *LabelField) Set(x, y int, l int) {
	lf.Labels[x+y*lf.Width] = l
}

// LabelConnectedComponents implements 8-connectivity connected component labeling
//
// Algorithm obtained from "Optimizing Two-Pass Connected-Component Labeling"
// by Kesheng Wu, Ekow Otoo, and Kenji Suzuki
func LabelConnectedComponents(field field.ScalarField) (labels LabelField, components map[int][][2]int) {
	labels, uf := cclFirstPass(field)
	uf.Flatten()

	// Second pass
	components = map[int][][2]int{}
	for y := 0; y < labels.Height; y++ {
		for x := 0; x < labels.Width; x++ {
			label := uf.Find(labels.Get(x, y))
			labels.Set(x, y, label)
			component, ok := components[label]
			if !ok {
				component = [][2]int{{x, y}}
			} else {
				component = append(component, [2]int{x, y})
			}
			components[label] = component
		}
	}
	return labels, components
}

// Pixel names were chosen as shown:
//   -------------
//   | a | b | c |
//   -------------
//   | d | e |   |
//   -------------
//   |   |   |   |
//   -------------
//
// The current pixel is e
// a, b, c, and d are its neighbors of interest
//
// If a point lies outside the bounds of the image, it is ignored
func cclFirstPass(field field.ScalarField) (LabelField, UnionFind) {
	w, h := field.Width(), field.Height()
	labels := LabelField{
		Width:  w,
		Height: h,
		Labels: make([]int, w*h),
	}
	uf := UnionFind{
		NextLabel: 1,
		labels:    []int{0},
	}

	for y := 0; y < field.Height(); y++ {
		for x := 0; x < field.Width(); x++ {
			// If the current point is not of interest then continue
			if field.Get(x, y) < 0.5 {
				continue
			}

			// If pixel b is in the image and above threshold:
			//    a, d, and c are its neighbors, so they are all part of the same component
			//    Therefore, there is no reason to check their labels
			//    so simply assign b's label to e
			if y > 0 && field.Get(x, y-1) >= 0.5 {
				labels.Set(x, y, labels.Get(x, y-1))
				continue
			}

			// If pixel c is in the image and above threshold:
			//    b is its neighbor, but a and d are not
			//    Therefore, we must check a and d's labels
			if x+1 < w && y > 0 && field.Get(x+1, y-1) >= 0.5 {
				c := labels.Get(x+1, y-1)
				labels.Set(x, y, c)

				// If pixel a is in the image and above threshold:
				//    Then a and c are connected through e
				//    Therefore, we must union their sets
				if x > 0 && field.Get(x-1, y-1) >= 0.5 {
					a := labels.Get(x-1, y-1)
					uf.Union(c, a)
					continue
				}

				// If pixel d is in the image and above threshold:
				//    Then d and c are connected through e
				//    Therefore we must union their sets
				if x > 0 && field.Get(x-1, y) >= 0.5 {
					d := labels.Get(x-1, y)
					uf.Union(c, d)
					continue
				}
				continue
			}

			// If pixel a is in the image and above threshold:
			//    We already know b and c are below threshold
			//    d is a's neighbor, so they already have the same label
			//    So simply assign a's label to e
			if x > 0 && y > 0 && field.Get(x-1, y-1) >= 0.5 {
				labels.Set(x, y, labels.Get(x-1, y-1))
				continue
			}

			// If pixel d is in the image and above threshold
			//    We already know a, b, and c are below threshold
			//    so simpy assign d's label to e
			if x > 0 && field.Get(x-1, y) >= 0.5 {
				labels.Set(x, y, labels.Get(x-1, y))
				continue
			}

			// All the neighboring pixels are below threshold,
			// Therefore the current pixel is a new component
			labels.Set(x, y, uf.MakeLabel())

		}
	}

	return labels, uf
}
