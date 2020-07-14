package window

import (
	"runtime"

	"github.com/Qendolin/go-printpixel/core/glcontext"
	"github.com/Qendolin/go-printpixel/core/glwindow"
	"github.com/Qendolin/go-printpixel/pkg/layout"
	"github.com/Qendolin/go-printpixel/renderer"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Window struct {
	GlWindow     glwindow.Extended
	Child        layout.Layoutable
	BeforeUpdate func()
	AfterUpdate  func()
	//Top, Right, Bottom, Left
	margins         []int
	init            bool
	drawables       []layout.Drawable
	texQuadRenderer renderer.TextureQuadRenderer
}

func NewHints() glwindow.Hints {
	return glwindow.NewHints()
}

func NewGlConfig(errorChanSize int) glcontext.GlConfig {
	return glcontext.NewGlConfig(errorChanSize)
}

func New(title string, cfg SimpleConfig) (*Window, error) {
	w := cfg.Width
	h := cfg.Height
	if w == 0 {
		w = 800
	}
	if h == 0 {
		h = 450
	}
	return NewCustom(title, w, h, cfg.ToHints(), nil, cfg.ToGlConfig())
}

func NewCustom(title string, width, height int, hints glwindow.Hints, monitor *glfw.Monitor, glCfg glcontext.GlConfig) (*Window, error) {
	err := glcontext.InitGlfw()
	if err != nil {
		return nil, err
	}
	glfwWin, err := glwindow.New(glwindow.Hints(hints), title, width, height, monitor)
	if err != nil {
		return nil, err
	}

	mLeft, mTop, mRight, mBot := glfwWin.GetVisibleFrameSize()

	win := Window{GlWindow: glfwWin, margins: []int{mTop, mRight, mBot, mLeft}}

	runtime.LockOSThread()
	win.GlWindow.MakeContextCurrent()
	err = glcontext.InitGl(glCfg)

	win.texQuadRenderer = *renderer.NewTextureQuadRenderer()

	return &win, err
}

func (win *Window) SetX(x int) {
	_, y := win.GlWindow.GetPos()
	win.GlWindow.SetPos(x+win.margins[3], y)
}

func (win *Window) SetY(y int) {
	x, _ := win.GlWindow.GetPos()
	win.GlWindow.SetPos(x, y+win.margins[0])
}

func (win *Window) X() int {
	x, _ := win.GlWindow.GetPos()
	return x - win.margins[3]
}

func (win *Window) Y() int {
	_, y := win.GlWindow.GetPos()
	return y - win.margins[0]
}

func (win *Window) SetWidth(width int) {
	_, h := win.GlWindow.GetSize()
	win.GlWindow.SetSize(width-win.margins[1]-win.margins[3], h)
}

func (win *Window) SetHeight(height int) {
	w, _ := win.GlWindow.GetSize()
	win.GlWindow.SetSize(w, height-win.margins[0]-win.margins[2])
}

func (win *Window) Width() int {
	w, _ := win.GlWindow.GetSize()
	return w + win.margins[1] + win.margins[3]
}

func (win *Window) Height() int {
	_, h := win.GlWindow.GetSize()
	return h + win.margins[0] + win.margins[2]
}

func (win *Window) Run() {
	for !win.GlWindow.ShouldClose() {
		win.Update()
	}
}

func (win *Window) Close() {
	win.GlWindow.Destroy()
	win.GlWindow = nil
}

func (win *Window) Update() {
	if win.BeforeUpdate != nil {
		win.BeforeUpdate()
	}
	win.texQuadRenderer.Bind()
	tq := make([]renderer.TextureQuad, len(win.drawables))
	for i, d := range win.drawables {
		tq[i] = d.TextureQuad()
	}
	win.texQuadRenderer.Draw(1/float32(win.Width()), 1/float32(win.Height()), tq...)
	win.GlWindow.SwapBuffers()
	glfw.PollEvents()
	if win.AfterUpdate != nil {
		win.AfterUpdate()
	}
}

func (win *Window) Layout() []layout.Layoutable {
	if win.Child == nil {
		return nil
	}
	win.Child.SetX(0)
	win.Child.SetY(0)
	win.Child.SetWidth(win.Width())
	win.Child.SetHeight(win.Height())

	if l, ok := win.Child.(layout.Layouter); ok {
		graph := layout.Layout(l)
		win.drawables = make([]layout.Drawable, 0)
		for node := range graph {
			if d, ok := node.(layout.Drawable); ok {
				win.drawables = append(win.drawables, d)
			}
		}
	} else if d, ok := win.Child.(layout.Drawable); ok {
		win.drawables = []layout.Drawable{d}
	}

	return nil
}
