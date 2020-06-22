package data_test

import (
	_ "image/png"
	"math/rand"
	"os"
	"testing"

	"github.com/Qendolin/go-printpixel/internal/data"
	"github.com/Qendolin/go-printpixel/internal/renderer"
	"github.com/Qendolin/go-printpixel/internal/test"
	"github.com/Qendolin/go-printpixel/internal/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	test.ParseArgs()
	m.Run()
}

func TestFileTexture(t *testing.T) {
	win, close := test.NewWindow(t)
	defer close()

	absPath, err := utils.ResolvePath("assets/textures/uv.png")
	assert.NoError(t, err)
	imgFile, err := os.Open(absPath)
	assert.NoError(t, err)
	defer imgFile.Close()

	tex := data.NewTexture(data.Tex2DTarget2D).As2D(0)
	tex.Bind(0)
	tex.ApplyDefaults()
	err = tex.AllocFile(imgFile, 0, gl.RGBA, gl.RGBA)
	assert.NoError(t, err)

	prog := test.NewProgram(t, "assets/shaders/quad_tex.vert", "assets/shaders/quad_tex.frag")
	cnv := renderer.NewTextureQuadWithProgram(prog)
	//Set Texture to nil in order to use our texture
	cnv.Texture = nil

	for !win.ShouldClose() {
		cnv.BindFor(func() {
			cnv.Draw()
		})
		win.SwapBuffers()
		glfw.PollEvents()
	}

}

func TestGeneratedTexture(t *testing.T) {
	win, close := test.NewWindow(t)
	defer close()

	tex := data.NewTexture(data.Tex2DTarget2D).As2D(0)
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

	tex.AllocBytes(data, 0, gl.RGB, 256, 256, gl.RGB)

	prog := test.NewProgram(t, "assets/shaders/quad_tex.vert", "assets/shaders/quad_tex.frag")
	cnv := renderer.NewTextureQuadWithProgram(prog)
	//Set Texture to nil in order to use our texture
	cnv.Texture = nil

	for !win.ShouldClose() {
		cnv.BindFor(func() {
			cnv.Draw()
		})
		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func TestUpdatingTexture(t *testing.T) {
	win, close := test.NewWindow(t)
	defer close()

	tex := data.NewTexture(data.Tex2DTarget2D).As2D(0)
	tex.Bind(0)
	tex.FilterMode(data.FilterLinear, data.FilterLinear)
	tex.WrapMode(data.WrapClampToEdge, data.WrapClampToEdge)

	tex.AllocEmpty(0, gl.RGB, 100, 100, gl.RGB)

	prog := test.NewProgram(t, "assets/shaders/quad_tex.vert", "assets/shaders/quad_tex.frag")
	cnv := renderer.NewTextureQuadWithProgram(prog)
	//Set Texture to nil in order to use our texture
	cnv.Texture = nil

	cnv.Bind()
	for !win.ShouldClose() {
		cnv.BindFor(func() {
			tex.WriteBytes([]byte{255, 255, 255}, 0, int32(rand.Intn(100)), int32(rand.Intn(100)), 1, 1, gl.RGB)
			cnv.Draw()
		})
		win.SwapBuffers()
		glfw.PollEvents()
	}
}
