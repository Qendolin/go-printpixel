package data_test

import (
	_ "image/png"
	"math/rand"
	"os"
	"testing"

	"github.com/Qendolin/go-printpixel/internal/canvas"
	"github.com/Qendolin/go-printpixel/internal/data"
	"github.com/Qendolin/go-printpixel/internal/test"
	"github.com/Qendolin/go-printpixel/internal/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func TestMain(m *testing.M) {
	test.ParseArgs()
	m.Run()
}

func TestFileTexture(t *testing.T) {
	win, close := test.NewWindow(t)
	defer close()

	absPath, err := utils.ResolvePath("assets/textures/uv.png")
	if err != nil {
		t.Fatal(err)
	}
	imgFile, err := os.Open(absPath)
	if err != nil {
		t.Fatal(err)
	}
	defer imgFile.Close()

	tex := data.NewTexture(data.Tex2DTarget2D).As2D(0)
	tex.Bind(0)
	tex.FilterMode(data.FilterLinear, data.FilterLinear)
	tex.WrapMode(data.WrapClampToEdge, data.WrapClampToEdge)
	err = tex.AllocFile(imgFile, 0, gl.RGBA, gl.RGBA)
	if err != nil {
		t.Fatal(err)
	}

	prog := test.NewProgram(t, "assets/shaders/quad_tex.vert", "assets/shaders/quad_tex.frag")
	cnv := canvas.NewCanvasWithProgram(prog)

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
	tex.FilterMode(data.FilterLinear, data.FilterLinear)
	tex.WrapMode(data.WrapClampToEdge, data.WrapClampToEdge)

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
	cnv := canvas.NewCanvasWithProgram(prog)

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
	cnv := canvas.NewCanvasWithProgram(prog)

	for !win.ShouldClose() {
		tex.WriteBytes([]byte{255, 255, 255}, 0, int32(rand.Intn(100)), int32(rand.Intn(100)), 1, 1, gl.RGB)
		cnv.BindFor(func() {
			cnv.Draw()
		})
		win.SwapBuffers()
		glfw.PollEvents()
	}
}
