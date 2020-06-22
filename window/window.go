package window

import (
	"github.com/Qendolin/go-printpixel/internal/context"
	iWin "github.com/Qendolin/go-printpixel/internal/window"
	"github.com/Qendolin/go-printpixel/layout"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Extended = iWin.Extended

type Layout struct {
	Window iWin.Extended
	//Top, Right, Bottom, Left
	margins      []int
	Child        layout.Layoutable
	BeforeUpdate func()
	AfterUpdate  func()
	init         bool
}

type Hints = iWin.Hints

func NewHints() Hints {
	return Hints(iWin.NewHints())
}

type GlConfig context.GlConfig

func NewGlConfig(errorChanSize int) GlConfig {
	return GlConfig(context.NewGlConfig(errorChanSize))
}

func New(hints Hints, title string, width, height int, monitor *glfw.Monitor) (Layout, error) {
	err := context.InitGlfw()
	if err != nil {
		return Layout{}, err
	}
	glfwWin, err := iWin.New(iWin.Hints(hints), title, width, height, monitor)
	if err != nil {
		return Layout{}, err
	}

	mLeft, mTop, mRight, mBot := glfwWin.GetVisibleFrameSize()

	win := Layout{Window: glfwWin, margins: []int{mTop, mRight, mBot, mLeft}}
	return win, nil
}

func (win Layout) SetX(x int) {
	_, y := win.Window.GetPos()
	win.Window.SetPos(x+win.margins[3], y)
}

func (win Layout) SetY(y int) {
	x, _ := win.Window.GetPos()
	win.Window.SetPos(x, y+win.margins[0])
}

func (win Layout) X() int {
	x, _ := win.Window.GetPos()
	return x - win.margins[3]
}

func (win Layout) Y() int {
	_, y := win.Window.GetPos()
	return y - win.margins[0]
}

func (win Layout) SetWidth(width int) {
	_, h := win.Window.GetSize()
	win.Window.SetSize(width-win.margins[1]-win.margins[3], h)
}

func (win Layout) SetHeight(height int) {
	w, _ := win.Window.GetSize()
	win.Window.SetSize(w, height-win.margins[0]-win.margins[2])
}

func (win Layout) Width() int {
	w, _ := win.Window.GetSize()
	return w + win.margins[1] + win.margins[3]
}

func (win Layout) Height() int {
	_, h := win.Window.GetSize()
	return h + win.margins[0] + win.margins[2]
}

func (win *Layout) Init(cfg GlConfig) (err error) {
	if !win.init {
		win.Window.MakeContextCurrent()
		err = context.InitGl(context.GlConfig(cfg))
	}
	win.init = true
	return
}

func (win Layout) Run() {
	for !win.Window.ShouldClose() {
		win.Update()
	}
}

func (win Layout) Close() {
	win.Window.Destroy()
	win.Window = nil
}

func (win Layout) Update() {
	if win.BeforeUpdate != nil {
		win.BeforeUpdate()
	}
	win.Window.SwapBuffers()
	glfw.PollEvents()
	if win.AfterUpdate != nil {
		win.AfterUpdate()
	}
}

func (win Layout) Layout() {
	if win.Child == nil {
		return
	}
	win.Child.SetX(0)
	win.Child.SetY(0)
	win.Child.SetWidth(win.Width())
	win.Child.SetHeight(win.Height())

	if l, ok := win.Child.(layout.Layouter); ok {
		l.Layout()
	}
}
