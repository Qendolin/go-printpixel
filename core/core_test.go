package core_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	_ "image/png"

	"github.com/Qendolin/go-printpixel/core"
	"github.com/Qendolin/go-printpixel/core/test"
	"github.com/Qendolin/go-printpixel/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	test.ParseArgs()
	m.Run()
}

func TestQuad(t *testing.T) {
	win, close := test.NewWindow(t)
	defer close()

	prog := test.NewProgram(t, "res://assets/shaders/quad_uv.vert", "res://assets/shaders/quad_uv.frag")
	prog.Bind()

	core.Quad().Bind()

	for !win.ShouldClose() {
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func TestArrayReader(t *testing.T) {
	ar := &core.ArrayReader{
		Stride: 5,
		Array:  []io.Reader{bytes.NewReader([]byte{1, 2, 3}), bytes.NewReader([]byte{1, 2, 3, 4, 5, 6, 7}), bytes.NewReader([]byte{1, 2, 3, 4, 5})},
	}

	buf := make([]byte, 15)
	n, err := io.ReadFull(ar, buf)
	assert.NoError(t, err)
	assert.Equal(t, 15, n)
	assert.EqualValues(t, []byte{1, 2, 3, 0, 0, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5}, buf)

	ar = &core.ArrayReader{
		Stride: 5,
		Array:  []io.Reader{bytes.NewReader([]byte{1, 2, 3}), bytes.NewReader([]byte{1, 2, 3, 4, 5, 6, 7}), bytes.NewReader([]byte{1, 2, 3, 4, 5})},
	}

	buf = make([]byte, 15)
	for i := 0; i < 15; i++ {
		n, err := ar.Read(buf[i : i+1])
		assert.NotEqual(t, 0, n)
		if err == io.EOF {
			break
		}
	}
	assert.NoError(t, err)
	assert.Equal(t, 15, n)
	assert.EqualValues(t, []byte{1, 2, 3, 0, 0, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5}, buf)
}

func TestTextureQuad(t *testing.T) {
	win, close := test.NewWindow(t)
	defer close()

	prog := test.NewProgram(t, "res://assets/shaders/quad_tex.vert", "res://assets/shaders/quad_tex.frag")
	prog.Bind()

	absPath, err := utils.ResolvePath("res://assets/textures/uv.png")
	assert.NoError(t, err)
	img, err := os.Open(absPath)
	assert.NoError(t, err)
	defer img.Close()

	tex, err := core.NewTexture2D(core.ImageReader(img))
	assert.NoError(t, err)
	tex.Bind(0)

	core.Quad().Bind()

	for !win.ShouldClose() {
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
		win.SwapBuffers()
		glfw.PollEvents()
	}
}
