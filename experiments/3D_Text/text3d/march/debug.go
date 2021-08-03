package march

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/Qendolin/go-printpixel/experiments/3D_Text/text3d/march/field"
)

var dbgImage *image.RGBA

const debug = false

func dbgReset() {
	dbgImage = nil
	os.RemoveAll("./debug/")
}

func dbgSave(typ string, nr int) {
	if dbgImage != nil {
		os.MkdirAll("./debug/", 0666)
		debug, _ := os.Create(fmt.Sprintf("./debug/%s_%03d.png", typ, nr))
		png.Encode(debug, dbgImage)
	}
}

func dbgInit(f field.ScalarField) {
	dbgImage = image.NewRGBA(image.Rect(0, 0, f.Width(), f.Height()))
	for x := 0; x < f.Width(); x++ {
		for y := 0; y < f.Height(); y++ {
			if f.Get(x, y) >= 0.5 {
				dbgImage.Set(x, y, color.RGBA{127, 127, 127, 255})
			} else {
				dbgImage.Set(x, y, color.Black)
			}
		}
	}
}

func dbgMarkStart(x, y int) {
	dbgImage.Set(x, y, color.RGBA{0, 255, 0, 255})
}

func dbgMarkVisit(x, y int) {
	r, g, b, _ := dbgImage.At(x, y).RGBA()
	dbgImage.Set(x, y, color.RGBA{uint8(r + 128), uint8(g), uint8(b), 255})
}
