package march_test

import (
	"image"
	_ "image/png"
	"os"
	"testing"

	"github.com/Qendolin/go-printpixel/core/test"
	"github.com/Qendolin/go-printpixel/experiments/3D_Text/text3d/march"
	"github.com/Qendolin/go-printpixel/experiments/3D_Text/text3d/march/field"
)

func TestMain(m *testing.M) {
	test.ParseArgs()
	m.Run()
}

func BenchmarkMarcher(b *testing.B) {
	msdf := load("../../assets/Regular-msdf.png", 4)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		march.March(msdf, false)
	}
}

func load(path string, scale float32) *field.MSDF {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	src := field.RGBAValueField(img)
	dst := field.NewValueField(int(float32(src.Width)*scale), int(float32(src.Height)*scale), 4)
	field.ScaleBlerpFull(src, dst)
	return field.NewMSDF(dst)
}
