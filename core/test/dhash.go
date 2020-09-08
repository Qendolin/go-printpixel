package test

import (
	"image"
	"math/big"

	"github.com/nfnt/resize"
)

func rgbaToGray(img *image.RGBA) *image.Gray {
	var (
		bounds = img.Bounds()
		gray   = image.NewGray(bounds)
	)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			var rgba = img.At(x, y)
			gray.Set(x, y, rgba)
		}
	}
	return gray
}

func dHash(rgba *image.RGBA) string {
	gray := rgbaToGray(rgba)
	gray = resize.Resize(12, 11, gray, resize.NearestNeighbor).(*image.Gray)
	var bits string

	for y := 0; y < gray.Rect.Dy(); y++ {
		for x := 0; x < gray.Rect.Dx()-1; x++ {
			i := (x + y*gray.Rect.Dx())
			if gray.Pix[i] > gray.Pix[i+1] {
				bits += "1"
			} else {
				bits += "0"
			}
		}
	}

	hash := big.NewInt(0)
	hash.SetString(bits, 2)
	return hash.Text(16)
}

func distance(a, b string) int {
	ha := big.NewInt(0)
	ha.SetString(a, 16)
	hb := big.NewInt(0)
	hb.SetString(b, 16)

	d := 0
	t := ha.Xor(hb, ha).Text(2)
	for _, b := range t {
		if b == '1' {
			d++
		}
	}
	return d
}
