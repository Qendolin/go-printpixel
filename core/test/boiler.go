package test

import (
	"image"
	"runtime"
	"testing"

	_ "image/jpeg"
	_ "image/png"

	"github.com/Qendolin/go-printpixel/core/glw"
	"github.com/Qendolin/go-printpixel/core/shader"
	"github.com/go-gl/gl/v3.3-core/gl"
)

const MaxFramesHeadless = 10

type TestingWindow struct {
	glw.Window
	expectedHash    string
	t               *testing.T
	closeCheckCount int
	isHeadless      bool
}

func (win *TestingWindow) ShouldClose() bool {
	win.closeCheckCount++
	if win.closeCheckCount > MaxFramesHeadless && win.isHeadless {
		win.assertResult()
		return true
	}
	if win.Window.ShouldClose() {
		win.assertResult()
		return true
	}
	return false
}

func (win *TestingWindow) assertResult() {
	w, h := win.GetSize()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	gl.ReadPixels(0, 0, int32(w), int32(h), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	actual := dHash(img)

	if win.expectedHash == "" {
		win.t.Errorf("expected hash has not been calculated, current hash is %v\n", actual)
		win.t.FailNow()
		return
	}

	if d := distance(win.expectedHash, actual); d > 3 {
		win.t.Errorf("result has a distance greater than 3 from expected result. Distance: %d\n", d)
	}
}

func NewWindow(t *testing.T, expectedHash string) (w glw.Window, close func()) {
	runtime.LockOSThread()
	conf := glw.BasicConfig("Test Window | "+t.Name(), 800, 450, glw.DontCare, glw.DontCare)
	conf.DebugContext = true
	conf.DebugHandler = func(err glw.DebugMessage) {
		if err.Critical {
			t.Error(err, "\n"+err.Stack)
		} else {
			t.Log(err)
		}
	}
	glfwWin, err := glw.New(conf)
	if err != nil {
		t.Fatal(err)
	}

	win := &TestingWindow{
		Window:          glfwWin,
		closeCheckCount: 0,
		isHeadless:      Args.Headless,
		t:               t,
		expectedHash:    expectedHash,
	}

	left, top, _, _ := win.GetFrameSize()
	win.SetPos(left, top)

	gl.ClearColor(0, 0, 0, 1)
	gl.Viewport(0, 0, int32(800), int32(450))

	return win, func() {
		win.Destroy()
	}
}

func NewProgram(t *testing.T, vsPath, fsPath string) *shader.Program {
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

	return prog
}
