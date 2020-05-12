package data_test

import (
	_ "image/png"
	"os"
	"testing"

	"github.com/Qendolin/go-printpixel/internal/canvas"
	"github.com/Qendolin/go-printpixel/internal/data"
	"github.com/Qendolin/go-printpixel/internal/test"
	"github.com/Qendolin/go-printpixel/internal/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

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

	tex := data.NewTexture(data.Texture2D)
	tex.Bind(0)
	tex.FilterMode(data.FilterLinear, data.FilterLinear)
	tex.WrapMode(data.WrapClampToEdge, data.WrapClampToEdge, data.WrapClampToEdge)
	err = tex.WriteFromFile2D(imgFile, 0, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE)
	if err != nil {
		t.Fatal(err)
	}

	prog := test.NewProgram(t, "assets/shaders/quad_tex.vert", "assets/shaders/quad_tex.frag")
	cnv := canvas.NewCanvasWithProgram(prog)

	for !win.ShouldClose() {
		cnv.BindFor(func() []func() {
			cnv.Draw()
			return nil
		})
		win.SwapBuffers()
		glfw.PollEvents()
	}

}

func TestGeneratedTexture(t *testing.T) {
	win, close := test.NewWindow(t)
	defer close()

	tex := data.NewTexture(data.Texture2D)
	tex.Bind(0)
	tex.FilterMode(data.FilterLinear, data.FilterLinear)
	tex.WrapMode(data.WrapClampToEdge, data.WrapClampToEdge, data.WrapClampToEdge)

	data := make([]byte, 100*100*3)

	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			data[(x+y*100)*3+0] = byte(x)
			data[(x+y*100)*3+1] = byte(y)
			data[(x+y*100)*3+2] = 0
		}
	}

	tex.WriteFromBytes(data, 100, 100, 0, gl.RGB, gl.RGB)

	prog := test.NewProgram(t, "assets/shaders/quad_tex.vert", "assets/shaders/quad_tex.frag")
	cnv := canvas.NewCanvasWithProgram(prog)

	for !win.ShouldClose() {
		cnv.BindFor(func() []func() {
			cnv.Draw()
			return nil
		})
		win.SwapBuffers()
		glfw.PollEvents()
	}
}
