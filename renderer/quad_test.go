package renderer_test

import (
	"os"
	"testing"

	_ "image/png"

	"github.com/Qendolin/go-printpixel/core/test"
	"github.com/Qendolin/go-printpixel/renderer"
	"github.com/Qendolin/go-printpixel/utils"
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
	prog.Bind()

	tq := renderer.NewTextureQuad()
	tq.Bind(0)

	renderer.NewQuad().Bind()

	for !win.ShouldClose() {
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func TestTextureQuadTexture(t *testing.T) {
	win, close := test.NewWindow(t)
	defer close()

	tq := renderer.NewTextureQuad()

	absPath, err := utils.ResolvePath("assets/textures/uv.png")
	assert.NoError(t, err)
	imgFile, err := os.Open(absPath)
	assert.NoError(t, err)
	defer imgFile.Close()

	tq.Bind(0)
	tq.Texture.ApplyDefaults()
	err = tq.Texture.AllocFile(imgFile, 0, gl.RGBA, gl.RGBA)
	assert.NoError(t, err)

	render := renderer.NewTextureQuadRenderer()
	render.Bind()

	for !win.ShouldClose() {
		render.Draw(1, 1, *tq)
		win.SwapBuffers()
		glfw.PollEvents()
	}
}
