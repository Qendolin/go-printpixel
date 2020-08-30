package data_test

import (
	"image"
	"image/draw"
	_ "image/png"
	"math/rand"
	"os"
	"testing"

	"github.com/Qendolin/go-printpixel/core"
	"github.com/Qendolin/go-printpixel/core/data"
	"github.com/Qendolin/go-printpixel/core/test"
	"github.com/Qendolin/go-printpixel/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	test.ParseArgs()
	m.Run()
}

func TestFileTexture(t *testing.T) {
	win, close := test.NewWindow(t)
	defer close()

	absPath, err := utils.ResolvePath("res://assets/textures/uv.png")
	require.NoError(t, err)
	imgFile, err := os.Open(absPath)
	require.NoError(t, err)
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	require.NoError(t, err)
	rgba := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	draw.Draw(rgba, rgba.Rect, img, image.Point{}, draw.Src)
	buf := rgba.Pix
	w, h := rgba.Rect.Size().X, rgba.Rect.Size().Y

	tex := data.NewTexture2D(nil, data.Tex2DTarget2D)
	tex.Bind(0)
	tex.ApplyDefaults()

	assert.NoError(t, tex.AllocBytes(buf, 0, gl.RGBA, int32(w), int32(h), gl.RGBA))

	test.NewProgram(t, "res://assets/shaders/quad_tex.vert", "res://assets/shaders/quad_tex.frag").Bind()
	core.Quad().Bind()

	for !win.ShouldClose() {
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
		win.SwapBuffers()
		glfw.PollEvents()
	}

}

func TestGeneratedTexture(t *testing.T) {
	win, close := test.NewWindow(t)
	defer close()

	tex := data.NewTexture2D(nil, data.Tex2DTarget2D)
	tex.Bind(0)
	tex.ApplyDefaults()

	data := make([]byte, 256*256*3)

	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			data[(x+y*256)*3+0] = byte(x)
			data[(x+y*256)*3+1] = byte(y)
			data[(x+y*256)*3+2] = 0
		}
	}

	assert.NoError(t, tex.AllocBytes(data, 0, gl.RGB, 256, 256, gl.RGB))

	test.NewProgram(t, "res://assets/shaders/quad_tex.vert", "res://assets/shaders/quad_tex.frag").Bind()
	core.Quad().Bind()

	for !win.ShouldClose() {
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
		win.SwapBuffers()
		glfw.PollEvents()
	}
}

// Texture is not initialized and may contain garbage data
func TestUpdatingTexture(t *testing.T) {
	win, close := test.NewWindow(t)
	defer close()

	tex := data.NewTexture2D(nil, data.Tex2DTarget2D)
	tex.Bind(0)
	tex.FilterMode(data.FilterLinear, data.FilterLinear)
	tex.WrapMode(data.WrapClampToEdge, data.WrapClampToEdge)

	tex.AllocEmpty(0, gl.RGB, 128, 128, gl.RGB)

	test.NewProgram(t, "res://assets/shaders/quad_tex.vert", "res://assets/shaders/quad_tex.frag").Bind()
	core.Quad().Bind()

	for !win.ShouldClose() {
		assert.NoError(t, tex.WriteBytes(0, int32(rand.Intn(128)), int32(rand.Intn(128)), 1, 1, gl.RGB, []byte{byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255))}))
		assert.NoError(t, tex.WriteBytes(0, int32(rand.Intn(128)), int32(rand.Intn(128)), 1, 1, gl.RGB, []byte{byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255))}))
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
		win.SwapBuffers()
		glfw.PollEvents()
	}
}
