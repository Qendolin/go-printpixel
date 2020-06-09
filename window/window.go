package window

import (
	"runtime"

	"github.com/Qendolin/go-printpixel/internal/context"
	iWin "github.com/Qendolin/go-printpixel/internal/window"
	"github.com/Qendolin/go-printpixel/layout"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Updater interface {
	Update()
}

type Window struct {
	handle *glfw.Window
	//Top, Right, Bottom, Left
	margins []int
	Child   layout.Layoutable
}

type Hints iWin.Hints

func NewHints() Hints {
	return Hints(iWin.NewHints())
}

type GlConfig context.GlConfig

func NewGlConfig(errorChanSize int) GlConfig {
	return GlConfig(context.NewGlConfig(errorChanSize))
}

func New(hints Hints, title string, width, height int, monitor *glfw.Monitor) (Window, error) {
	context.InitGlfw()
	glfwWin, err := iWin.New(iWin.Hints(hints), title, width, height, monitor)
	if err != nil {
		return Window{}, err
	}

	mLeft, mTop, mRight, mBot := glfwWin.GetFrameSize()

	win := Window{handle: glfwWin, margins: []int{mTop, mRight, mBot, mLeft}}
	return win, nil
}

func (win Window) SetX(x int) {
	_, y := win.handle.GetPos()
	win.handle.SetPos(x+win.margins[3], y)
}

func (win Window) SetY(y int) {
	x, _ := win.handle.GetPos()
	win.handle.SetPos(x, y+win.margins[0])
}

func (win Window) X() int {
	x, _ := win.handle.GetPos()
	return x - win.margins[3]
}

func (win Window) Y() int {
	_, y := win.handle.GetPos()
	return y - win.margins[0]
}

func (win Window) SetWidth(width int) {
	_, h := win.handle.GetSize()
	win.handle.SetSize(width-win.margins[1]-win.margins[3], h)
}

func (win Window) SetHeight(height int) {
	w, _ := win.handle.GetSize()
	win.handle.SetSize(w, height-win.margins[0]-win.margins[2])
}

func (win Window) Width() int {
	w, _ := win.handle.GetSize()
	return w + win.margins[1] + win.margins[3]
}

func (win Window) Height() int {
	_, h := win.handle.GetSize()
	return h + win.margins[0] + win.margins[2]
}

func (win Window) Run(cfg GlConfig) {
	win.handle.MakeContextCurrent()
	context.InitGl(context.GlConfig(cfg))
	for !win.handle.ShouldClose() {
		win.Update()
	}
}

func (win Window) Close() {
	runtime.LockOSThread()
	win.handle.Destroy()
	win.handle = nil
	runtime.UnlockOSThread()
}

func (win Window) Update() {
	win.handle.SwapBuffers()
	glfw.PollEvents()
}

func (win Window) Layout() {
	if win.Child == nil {
		return
	}
	win.Child.SetX(0)
	win.Child.SetY(0)
	win.Child.SetWidth(win.Width())
	win.Child.SetHeight(win.Height())
}
