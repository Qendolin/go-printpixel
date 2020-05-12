package test

import (
	"runtime"
	"testing"

	"github.com/Qendolin/go-printpixel/internal/context"
	"github.com/Qendolin/go-printpixel/internal/shader"
	"github.com/Qendolin/go-printpixel/internal/window"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func NewWindow(t *testing.T) (w glfw.Window, close func()) {
	runtime.LockOSThread()
	err := context.InitGlfw()
	if err != nil {
		t.Fatal(err)
	}

	hints := window.NewHints()
	win, err := window.New(hints, "Test Window", 800, 450, nil)
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

	return *win, func() {
		win.Destroy()
		context.Terminate()
	}
}

func NewProgram(t *testing.T, vsPath, fsPath string) shader.Program {
	vs, err := shader.NewShaderFromPath(vsPath, shader.TypeVertex)
	if err != nil {
		t.Fatal(err)
	}

	fs, err := shader.NewShaderFromPath(fsPath, shader.TypeFragment)
	if err != nil {
		t.Fatal(err)
	}
	prog, err := shader.NewProgram(vs, fs)
	if err != nil {
		t.Fatal(err)
	}
	fs.Destroy()
	vs.Destroy()

	if ok, log := prog.Validate(); ok {
		t.Logf("Program Validation Log: \n\n%v\n\n", log)
	} else {
		t.Fatalf("Program Validation Log: \n\n%v\n\n", log)
	}

	return *prog
}
