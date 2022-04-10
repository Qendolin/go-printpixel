package field_test

import (
	"image"
	"image/png"
	"math"
	"math/rand"
	"os"
	"testing"

	"github.com/Qendolin/go-printpixel/core/test"
	"github.com/Qendolin/go-printpixel/experiments/3D_Text/text3d/march/field"
)

func TestMain(m *testing.M) {
	test.ParseArgs()
	m.Run()
}

func BenchmarkBlerpF(b *testing.B) {
	src, err := loadRGBValueField()
	if err != nil {
		b.Fatal(err)
	}

	dst := field.NewValueField(src.Width*8, src.Height*8, 3)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		field.ScaleBlerp(src, dst)
	}
}

func BenchmarkBlerpI(b *testing.B) {
	src, err := loadRGBValueField()
	if err != nil {
		b.Fatal(err)
	}
	srcI := field.NewValueFieldI(src)

	dst := field.NewValueField(src.Width*8, src.Height*8, 3)
	dstI := field.NewValueFieldI(dst)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		field.ScaleBlerpIptr(srcI, dstI)
	}
}

func BenchmarkBlerpIRand(b *testing.B) {
	src := &field.ValueFieldI{
		Width:         37,
		Height:        37,
		ComponentSize: 3,
		Values:        make([]uint32, 37*37*3),
	}

	for i := range src.Values {
		src.Values[i] = rand.Uint32()
	}

	dst := &field.ValueFieldI{
		Width:         37 * 8,
		Height:        37 * 8,
		ComponentSize: 3,
		Values:        make([]uint32, 37*8*37*8*3),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		field.ScaleBlerpIptr(src, dst)
	}
}

func BenchmarkBlerpC(b *testing.B) {
	src, err := loadRGBValueField()
	if err != nil {
		b.Fatal(err)
	}
	srcI := field.NewValueFieldI(src)

	dst := field.NewValueField(src.Width*8, src.Height*8, 3)
	dstI := field.NewValueFieldI(dst)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		field.ScaleBlerpIC(srcI, dstI)
	}
}

func BenchmarkBlerpCFSimd(b *testing.B) {
	src, err := loadRGBValueField()
	if err != nil {
		b.Fatal(err)
	}

	dst := field.NewValueField(src.Width*8, src.Height*8, 3)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		field.ScaleBlerpFSimd(src, dst)
	}
}

func BenchmarkBlerpCFSimd2(b *testing.B) {
	src, err := loadRGBAValueField()
	if err != nil {
		b.Fatal(err)
	}

	dst := field.NewValueField(src.Width*8, src.Height*8, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		field.ScaleBlerpFSimd2(src, dst)
	}
}

// FIXME: segfault
// func BenchmarkBlerpCSimd(b *testing.B) {
// 	src, err := loadRGBAValueField()
// 	if err != nil {
// 		b.Fatal(err)
// 	}
// 	srcI := field.NewValueFieldI(src)

// 	dst := field.NewValueField(src.Width*8, src.Height*8, 4)
// 	dstI := field.NewValueFieldI(dst)
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		field.ScaleBlerpSimd(srcI, dstI)
// 	}
// }

// FIXME: segfault
// func BenchmarkBlerpCSimd2(b *testing.B) {
// 	src, err := loadRGBAValueField()
// 	if err != nil {
// 		b.Fatal(err)
// 	}
// 	srcI := field.NewValueFieldI(src)

// 	dst := field.NewValueField(src.Width*8, src.Height*8, 4)
// 	dstI := field.NewValueFieldI(dst)
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		field.ScaleBlerpSimd2(srcI, dstI)
// 	}
// }

func BenchmarkBlerpCSimdRand(b *testing.B) {
	srcI := &field.ValueFieldI{
		Width:         37,
		Height:        37,
		ComponentSize: 4,
		Values:        make([]uint32, 37*37*4),
	}

	for i := range srcI.Values {
		srcI.Values[i] = rand.Uint32()
	}

	dstI := &field.ValueFieldI{
		Width:         37 * 8,
		Height:        37 * 8,
		ComponentSize: 4,
		Values:        make([]uint32, 37*8*37*8*4),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		field.ScaleBlerpSimd(srcI, dstI)
	}
}

func BenchmarkBlerpCSimd2Rand(b *testing.B) {
	srcI := &field.ValueFieldI{
		Width:         37,
		Height:        37,
		ComponentSize: 4,
		Values:        make([]uint32, 37*37*4),
	}

	for i := range srcI.Values {
		srcI.Values[i] = rand.Uint32()
	}

	dstI := &field.ValueFieldI{
		Width:         37 * 8,
		Height:        37 * 8,
		ComponentSize: 4,
		Values:        make([]uint32, 37*8*37*8*4),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		field.ScaleBlerpSimd2(srcI, dstI)
	}
}

func BenchmarkBlerpFull(b *testing.B) {
	src, err := loadRGBAValueField()
	if err != nil {
		b.Fatal(err)
	}
	dst := field.NewValueField(src.Width*8, src.Height*8, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		field.ScaleBlerpFull(src, dst)
	}
}

func TestBoth(t *testing.T) {
	src, err := loadRGBValueField()
	if err != nil {
		t.Fatal(err)
	}
	// srcI := field.NewValueFieldI(src)

	dst := field.NewValueField(src.Width*8, src.Height*8, 3)
	dstI := field.NewValueFieldI(dst)
	// field.ScaleBlerpIptr(srcI, dstI)
	// field.ScaleBlerpIC(srcI, dstI)
	field.ScaleBlerp(src, dst)

	err = saveRGBValueFiled(dst, "./float32.png")
	if err != nil {
		t.Fatal(err)
	}

	for i, v := range dstI.Values {
		dst.Values[i] = float32(v) / float32(math.MaxUint32)
	}
	err = saveRGBValueFiled(dst, "./uint32.png")
	if err != nil {
		t.Fatal(err)
	}
}

// FIXME: segfault
// func TestFSimd(t *testing.T) {
// 	src, err := loadRGBAValueField()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	dst := field.NewValueField(src.Width*8, src.Height*8, 4)
// 	field.ScaleBlerpFSimd(src, dst)

// 	err = saveRGBAValueFiled(dst, "./fsimd.png")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

func TestFSimd2(t *testing.T) {
	src, err := loadRGBAValueField()
	if err != nil {
		t.Fatal(err)
	}
	dst := field.NewValueField(src.Width*8, src.Height*8, 4)
	field.ScaleBlerpFSimd2(src, dst)

	err = saveRGBAValueFiled(dst, "./fsimd2.png")
	if err != nil {
		t.Fatal(err)
	}
}

func TestFSimd2Range(t *testing.T) {
	src, err := loadRGBAValueField()
	if err != nil {
		t.Fatal(err)
	}

	for x := 0; x < 200; x++ {
		for y := 0; y < 200; y++ {
			dst := field.NewValueField(src.Width+x, src.Height+y, 4)
			field.ScaleBlerpFSimd2(src, dst)
		}
	}

	dst := field.NewValueField(1000, 1000, 4)
	field.ScaleBlerpFSimd2(src, dst)
	dst = field.NewValueField(10000, 10000, 4)
	field.ScaleBlerpFSimd2(src, dst)
}

func TestFull(t *testing.T) {
	src, err := loadRGBAValueField()
	if err != nil {
		t.Fatal(err)
	}

	dst := field.NewValueField(src.Width*8, src.Height*1, 4)
	field.ScaleBlerpFull(src, dst)
	err = saveRGBAValueFiled(dst, "./full_x.png")
	if err != nil {
		t.Fatal(err)
	}
	dst = field.NewValueField(src.Width*1, src.Height*8, 4)
	field.ScaleBlerpFull(src, dst)
	err = saveRGBAValueFiled(dst, "./full_y.png")
	if err != nil {
		t.Fatal(err)
	}

	for sx := 1; sx < 25; sx++ {
		for sy := 1; sy < 25; sy++ {
			t.Logf("Scale X: %3d Scale Y: %3d\n", sx, sy)
			dst := field.NewValueField(src.Width*sx, src.Height*sy, 4)
			field.ScaleBlerpFull(src, dst)
		}
	}

	for i := 0; i < 100; i++ {
		sx := 1 + 100*rand.Float32()
		sy := 1 + 100*rand.Float32()
		t.Logf("Scale X: %3.2f Scale Y: %3.2f\n", sx, sy)
		dst := field.NewValueField(int(float32(src.Width)*sx), int(float32(src.Height)*sy), 4)
		field.ScaleBlerpFull(src, dst)
	}
}

// FIXME: segfault
// func TestISimd(t *testing.T) {
// 	src, err := loadRGBAValueField()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	srcI := field.NewValueFieldI(src)

// 	dst := field.NewValueField(src.Width*8, src.Height*8, 4)
// 	dstI := field.NewValueFieldI(dst)
// 	field.ScaleBlerpSimd(srcI, dstI)

// 	for i, v := range dstI.Values {
// 		dst.Values[i] = float32(v) / float32(math.MaxUint32)
// 	}
// 	err = saveRGBAValueFiled(dst, "./isimd.png")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// FIXME: segfault
// func TestISimd2(t *testing.T) {
// 	src, err := loadRGBAValueField()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	srcI := field.NewValueFieldI(src)

// 	dst := field.NewValueField(src.Width*8, src.Height*8, 4)
// 	dstI := field.NewValueFieldI(dst)
// 	field.ScaleBlerpSimd2(srcI, dstI)

// 	for i, v := range dstI.Values {
// 		dst.Values[i] = float32(v) / float32(math.MaxUint32)
// 	}
// 	err = saveRGBAValueFiled(dst, "./isimd2.png")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// FIXME: segfault
// func TestISimd3(t *testing.T) {
// 	src, err := loadRGBAValueField()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	srcI := field.NewValueFieldI8(src)

// 	dst := field.NewValueField(src.Width*8, src.Height*8, 4)
// 	dstI := field.NewValueFieldI(dst)
// 	field.ScaleBlerpSimd3(srcI, dstI)

// 	for i, v := range dstI.Values {
// 		dst.Values[i] = float32(v) / float32(math.MaxUint8)
// 	}
// 	err = saveRGBAValueFiled(dst, "./isimd3.png")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

func loadRGBValueField() (vf *field.ValueField, err error) {
	file, err := os.Open("../../../assets/g-msdf.png")
	if err != nil {
		return
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return
	}
	return field.RGBValueField(img), nil
}

func loadRGBAValueField() (vf *field.ValueField, err error) {
	file, err := os.Open("../../../assets/blerp-test.png")
	if err != nil {
		return
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return
	}
	return field.RGBAValueField(img), nil
}

func saveRGBValueFiled(vf *field.ValueField, path string) error {
	img := image.NewRGBA(image.Rect(0, 0, vf.Width, vf.Height))
	for i := 0; i < vf.Width*vf.Height; i++ {
		img.Pix[i*4+0] = uint8(vf.Values[i*3+0] * 0xff)
		img.Pix[i*4+1] = uint8(vf.Values[i*3+1] * 0xff)
		img.Pix[i*4+2] = uint8(vf.Values[i*3+2] * 0xff)
		img.Pix[i*4+3] = 0xff
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	return png.Encode(file, img)
}

func saveRGBAValueFiled(vf *field.ValueField, path string) error {
	img := image.NewRGBA(image.Rect(0, 0, vf.Width, vf.Height))
	for i := 0; i < vf.Width*vf.Height; i++ {
		img.Pix[i*4+0] = uint8(vf.Values[i*4+0] * 0xff)
		img.Pix[i*4+1] = uint8(vf.Values[i*4+1] * 0xff)
		img.Pix[i*4+2] = uint8(vf.Values[i*4+2] * 0xff)
		img.Pix[i*4+3] = uint8(0xff)
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	return png.Encode(file, img)
}
