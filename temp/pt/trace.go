package main

import (
	"image"
	"image/color"
	_ "image/png"
	"math"
	"os"

	"github.com/dennwc/gotrace"
)

func main() {
	f, err := os.Open("./g.png")
	panicIf(err)
    src, _, err := image.Decode(f)
	dr := image.Rect(0, 0, src.Bounds().Dx()*9, src.Bounds().Dy()*9)
	img := image.NewRGBA(dr)
	scale(img, src)

	bm := gotrace.NewBitmapFromImage(img, func(x, y int, c color.Color) bool {
		r, g, b, a := c.RGBA()
		return median(float64(r)/float64(a),float64(g)/float64(a),float64(b)/float64(a)) > 0.5
	})
	paths, _ := gotrace.Trace(bm, nil)
	f, err = os.OpenFile("out.svg", os.O_CREATE, 0666)
	panicIf(err)
	gotrace.WriteSvg(f, img.Bounds(), paths, "")
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func median(r, g, b float64) float64 {
    return math.Max(math.Min(r, g), math.Min(math.Max(r, g), b))
}