package core_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	_ "image/png"

	"github.com/Qendolin/go-printpixel/core"
	"github.com/Qendolin/go-printpixel/core/data"
	"github.com/Qendolin/go-printpixel/core/test"
	"github.com/Qendolin/go-printpixel/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	test.ParseArgs()
	m.Run()
}

func TestQuad(t *testing.T) {
	win, close := test.NewWindow(t, "40080100200400801000000")
	defer close()

	prog := test.NewProgram(t, "@lib/assets/shaders/quad_uv.vert", "@lib/assets/shaders/quad_uv.frag")
	prog.Bind()

	core.Quad().Bind()

	for !win.ShouldClose() {
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func TestTextureQuad(t *testing.T) {
	win, close := test.NewWindow(t, "240c83705206c1682d000000")
	defer close()

	prog := test.NewProgram(t, "@lib/assets/shaders/quad_tex.vert", "@lib/assets/shaders/quad_tex.frag")
	prog.Bind()

	absPath, err := utils.ResolvePath("@lib/assets/textures/uv.png")
	require.NoError(t, err)
	img, err := os.Open(absPath)
	require.NoError(t, err)
	defer img.Close()

	tex, err := core.NewTexture2D(core.InitFiles(0, img), data.RGBA8)
	require.NoError(t, err)
	tex.Bind(0)

	core.Quad().Bind()

	for !win.ShouldClose() {
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func TestTextureNotFound(t *testing.T) {
	win, close := test.NewWindow(t, "241b8090020dc04837000000")
	defer close()

	prog := test.NewProgram(t, "@lib/assets/shaders/quad_tex.vert", "@lib/assets/shaders/quad_tex.frag")
	prog.Bind()

	tex := core.MustNewTexture2D(core.InitPaths(0, "/file/that/does/not.exist"), data.RGBA8)
	tex.Bind(0)

	core.Quad().Bind()

	for !win.ShouldClose() {
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func TestTextureUndecodeable(t *testing.T) {
	win, close := test.NewWindow(t, "1c0e81502e0941a813000000")
	defer close()

	prog := test.NewProgram(t, "@lib/assets/shaders/quad_tex.vert", "@lib/assets/shaders/quad_tex.frag")
	prog.Bind()

	tex := core.MustNewTexture2D(core.InitFiles(0, bytes.NewBuffer([]byte{1, 2, 3, 4, 5})), data.RGBA8)
	tex.Bind(0)

	core.Quad().Bind()

	for !win.ShouldClose() {
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func TestTextureError(t *testing.T) {
	win, close := test.NewWindow(t, "41881104a04418801000000")
	defer close()

	prog := test.NewProgram(t, "@lib/assets/shaders/quad_tex.vert", "@lib/assets/shaders/quad_tex.frag")
	prog.Bind()

	tex := core.MustNewTexture2D(core.TextureInitializer{Levels: []interface{}{fmt.Errorf("some error")}}, data.RGBA8)
	tex.Bind(0)

	core.Quad().Bind()

	for !win.ShouldClose() {
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
		win.SwapBuffers()
		glfw.PollEvents()
	}
}
