package renderer_test

import (
	"os"
	"testing"

	_ "image/png"

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

func TestTextureQuad(t *testing.T) {
	win, close := test.NewWindow(t)
	defer close()

	prog := test.NewProgram(t, "assets/shaders/quad_uv.vert", "assets/shaders/quad_uv.frag")

	cnv := renderer.NewTextureQuadWithProgram(prog)
	for !win.ShouldClose() {
		cnv.BindFor(func() {
			cnv.Draw()
		})
		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func TestTextureQuadTexture(t *testing.T) {
	win, close := test.NewWindow(t)
	defer close()

	cnv := renderer.NewTextureQuad()

	absPath, err := utils.ResolvePath("assets/textures/uv.png")
	assert.NoError(t, err)
	imgFile, err := os.Open(absPath)
	assert.NoError(t, err)
	defer imgFile.Close()
	cnv.Texture.Bind(0)
	cnv.Texture.ApplyDefaults()
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
