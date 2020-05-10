package data_test

import (
	_ "image/png"
	"os"
	"testing"

	"github.com/Qendolin/go-printpixel/internal/canvas"
	"github.com/Qendolin/go-printpixel/internal/context"
	"github.com/Qendolin/go-printpixel/internal/data"
	"github.com/Qendolin/go-printpixel/internal/utils"
	"github.com/Qendolin/go-printpixel/internal/window"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func TestTextures(t *testing.T) {
	err := context.InitGlfw()
	if err != nil {
		t.Fatal(err)
	}
	defer context.Terminate()

	hints := window.NewHints()
	win, err := window.New(hints, "Test Window", 800, 450, nil)
	defer win.Destroy()
	if err != nil {
		t.Fatal(err)
	}
	win.MakeContextCurrent()

	cfg := context.NewGlConfig(0)
	cfg.Debug = true
	go func() {
		for err := range cfg.Errors {
			if err.Fatal {
				t.Error(err)
			}
			t.Log(err)
		}
	}()
	err = context.InitGl(cfg)
	if err != nil {
		t.Fatal(err)
	}
	gl.ClearColor(1, 0, 0, 1)

	imgFile, err := os.Open(utils.MustResolveModulePath("assets/textures/uv.png"))
	if err != nil {
		t.Fatal(err)
	}
	defer imgFile.Close()

	tex := data.NewTexture(data.Texture2D)
	tex.Bind(gl.TEXTURE0)
	tex.FilterMode(data.FilterLinear, data.FilterLinear)
	tex.WrapMode(data.WrapClampToEdge, data.WrapClampToEdge, data.WrapClampToEdge)
	err = tex.WriteFromFile2D(imgFile, 0, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE)
	if err != nil {
		t.Fatal(err)
	}

	cnv := canvas.NewCanvas()

	for !win.ShouldClose() {
		cnv.BindFor(func() []func() {
			cnv.Draw()
			return nil
		})
		win.SwapBuffers()
		glfw.PollEvents()
	}

}
