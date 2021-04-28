package main

import (
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"os"
	"time"

	"github.com/Qendolin/go-printpixel/core/glcontext"
	"github.com/Qendolin/go-printpixel/core/glwindow"
	"github.com/Qendolin/go-printpixel/pkg/window"
	"github.com/go-gl/gl/v3.3-compatibility/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/zx9597446/rdp"
)

func main() {
	f, err := os.Open("./g.png")
	panicIf(err)
	start := time.Now()
	src, _, err := image.Decode(f)
	dr := image.Rect(0, 0, src.Bounds().Dx()*6, src.Bounds().Dy()*6)
	img := image.NewRGBA(dr)
	scale(img, src)
	// img := src
	log.Printf("Scale %v\n", time.Since(start))

	verts := make([]Point, 0)

	m := Marcher{Quality: 8, Img: img, Discriminator: func(c color.Color) float64 {
		r, g, b, a := c.RGBA()
		return median(float64(r)/float64(a), float64(g)/float64(a), float64(b)/float64(a))
	}}
	verts = m.Process()[0]

	log.Printf("March %v\n", time.Since(start))

	// for x := 0; x < img.Bounds().Max.X-1.; x++ {
	// 	for y := 0; y < img.Bounds().Max.Y-1; y++ {
	// 		value := 0
	// 		if sample(img, x, y) > 0.5 {
	// 			value |= 8
	// 		}
	// 		if sample(img, x+1, y) > 0.5 {
	// 			value |= 4
	// 		}
	// 		if sample(img, x, y+1) > 0.5 {
	// 			value |= 1
	// 		}
	// 		if sample(img, x+1, y+1) > 0.5 {
	// 			value |= 2
	// 		}

	// 		verts = append(verts, translate(x, y, lookup(value)...)...)
	// 	}
	// }

	temp := make([]image.Point, len(verts))
	for i := 0; i < len(verts); i++ {
		temp[i] = image.Pt(int(verts[i][0]*100), int(verts[i][1]*100))
	}
	// temp := marchingsquare.Process(img, func(r, g, b, a uint32) bool {
	// 	return  median(float64(r)/float64(a),float64(g)/float64(a),float64(b)/float64(a))*2-1 > 0
	// })
	temp = rdp.Process(temp, 30)
	verts = make([]Point, len(temp))
	for i, p := range temp {
		verts[i] = Point{float32(p.X) / 100, float32(p.Y) / 100}
	}

	verts = clean(verts, 5)

	// verts = decimate(verts, 0.002)

	log.Printf("Verts: %d\n", len(verts))
	log.Printf("End %v\n", time.Since(start))

	win := setup()
	win.GlWindow.MakeContextCurrent()
	panicIf(gl.Init())
	gl.ClearColor(0., 0., 0., 1.)
	gl.LineWidth(5.)
	gl.PointSize(5.)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(-5, float64(img.Bounds().Dx()+5), float64(img.Bounds().Dy()+5), -15, 1.0, -1.0)

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Translatef(-1, -1, 0)

	for !win.GlWindow.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.Color3f(1, 1, 1)
		gl.Begin(gl.LINE_LOOP)
		for i := 0; i < len(verts); i++ {
			gl.Vertex2f(verts[i][0], verts[i][1])
		}
		gl.End()
		//*
		gl.Begin(gl.POINTS)
		for i := 0; i < len(verts); i++ {
			gl.Color3f(calcColor(i, len(verts)))
			gl.Vertex2f(verts[i][0], verts[i][1])
		}
		gl.End()
		//*/
		win.GlWindow.SwapBuffers()
		glfw.PollEvents()
	}
}

func calcColor(i, l int) (r, g, b float32) {
	// a := float32(i)/float32(l) * 0.75 + 0.25
	a := float32(0.85)
	switch i % 3 {
	case 0:
		return a, 0, 0
	case 1:
		return 0, a, 0
	case 2:
		return 0, 0, a
	}
	return 0, 0, 0
}

func sample(img *image.RGBA, x, y int) float64 {
	r := img.Pix[(x+y*img.Bounds().Dx())*4]
	g := img.Pix[(x+y*img.Bounds().Dx())*4+1]
	b := img.Pix[(x+y*img.Bounds().Dx())*4+2]
	// return (float64(r+g+b)/(3*255))
	return median(float64(r)/255, float64(g)/255, float64(b)/255)
}

func median(r, g, b float64) float64 {
	return math.Max(math.Min(r, g), math.Min(math.Max(r, g), b))
}

func translate(x, y int, values ...float32) []float32 {
	for i, v := range values {
		if i%2 == 0 {
			values[i] = v + float32(x)
		} else {
			values[i] = v + float32(y)
		}
	}
	return values
}

func clean(verts []Point, tolerance float32) []Point {
	verts = append([]Point{verts[len(verts)-2], verts[len(verts)-1]}, verts...)
	j := 3
	p1, p2, p3 := verts[0], verts[1], verts[2]
	for i := 3; i < len(verts); i++ {
		p4 := verts[i]

		d := p2.Sub(p3).Length()
		if d < tolerance {
			// https://www.geeksforgeeks.org/program-for-point-of-intersection-of-two-lines/
			a1, b1 := p2[1]-p1[1], p1[0]-p2[0]
			c1 := a1*(p1[0]) + b1*(p1[1])

			// Line CD represented as a2x + b2y = c2
			a2, b2 := p3[1]-p4[1], p4[0]-p3[0]
			c2 := a2*(p4[0]) + b2*(p4[1])

			determinant := a1*b2 - a2*b1

			if determinant != 0 {
				x := (b2*c1 - b1*c2) / determinant
				y := (a1*c2 - a2*c1) / determinant

				verts[j-2] = Point{x, y}
				verts[j-1] = p4
			} else {
				verts[j-2] = p2
				verts[j-1] = p4
			}
		} else {
			verts[j] = p4
			j++
		}

		p1 = p2
		p2 = p3
		p3 = p4
	}
	return verts[:j-2]
}

func decimate(verts []Point, tolerance float32) []Point {
	j := 0
	p1 := verts[0]
	p2 := verts[1]
	for i := 2; i < len(verts); i++ {
		p3 := verts[i]

		v1 := p1.Sub(p2)
		v2 := p2.Sub(p3)

		s := v1.Dot(v2) / (v1.Length() * v2.Length())
		if s < 1-tolerance {
			verts[j] = p3
			j++
		}

		p1 = p2
		p2 = p3
	}

	return verts[:j]
}

// func lookup(v int) []float32 {
// 	switch v {
// 	default:
// 		return []float32{}
// 	case 1, 14:
// 		return []float32{0, 0.5, 0.5, 1}
// 	case 2, 13:
// 		return []float32{0.5, 1, 1, 0.5}
// 	case 3, 12:
// 		return []float32{0, 0.5, 1, 0.5}
// 	case 4, 11:
// 		return []float32{0.5, 0, 1, 0.5}
// 	case 5:
// 		return []float32{0, 0.5, 0.5, 0, 0.5, 1, 1, 0.5}
// 	case 6, 9:
// 		return []float32{0.5, 0, 0.5, 1}
// 	case 7, 8:
// 		return []float32{0, 0.5, 0.5, 0}
// 	case 10:
// 		return []float32{0.5, 0., 1, 0.5, 0, 0.5, 0.5, 1}
// 	}
// }

func setup() *window.Window {
	cfg := glcontext.NewGlConfig(1)
	cfg.Debug = true
	cfg.Multisampling = false
	hints := glwindow.NewHints()
	hints.ContextVersionMajor.Value = 2
	hints.ContextVersionMinor.Value = 1
	// hints.OpenGLForwardCompatible.Value = true
	// hints.OpenGLProfile.Value = glwindow.OpenGLCompatProfile
	hints.OpenGLDebugContext.Value = true
	win, err := window.NewCustom("MS Example", 900, 900, hints, nil, cfg)
	panicIf(err)

	go handleErrors(cfg.Errors)
	return win
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func handleErrors(errs <-chan glcontext.Error) {
	for err := range errs {
		if err.Fatal {
			log.Fatalf("%v\n%v", err, err.Stack)
		}
		log.Printf("%v\n", err)
	}
}
