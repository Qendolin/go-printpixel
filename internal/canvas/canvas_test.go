package canvas_test

import (
	"os"
	"testing"

	_ "image/png"

	"github.com/Qendolin/go-printpixel/internal/canvas"
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

func TestCanvasQuad(t *testing.T) {
	win, close := test.NewWindow(t)
	defer close()

	prog := test.NewProgram(t, "assets/shaders/quad_uv.vert", "assets/shaders/quad_uv.frag")

	cnv := canvas.NewCanvasWithProgram(prog)
	for !win.ShouldClose() {
		cnv.BindFor(func() {
			cnv.Draw()
		})
		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func TestCanvasTexture(t *testing.T) {
	win, close := test.NewWindow(t)
	defer close()

	cnv := canvas.NewCanvas()

	absPath, err := utils.ResolvePath("assets/textures/uv.png")
	assert.NoError(t, err)
	imgFile, err := os.Open(absPath)
	assert.NoError(t, err)
	defer imgFile.Close()
	cnv.Texture.Bind(0)
	cnv.Texture.DefaultModes()
	err = cnv.Texture.AllocFile(imgFile, 0, gl.RGBA, gl.RGBA)
	assert.NoError(t, err)

	for !win.ShouldClose() {
		cnv.BindFor(func() {
			cnv.Draw()
		})
		win.SwapBuffers()
		glfw.PollEvents()
	}
}
