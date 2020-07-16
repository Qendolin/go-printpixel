package window

import (
	"runtime"

	"github.com/Qendolin/go-printpixel/core/glcontext"
	"github.com/Qendolin/go-printpixel/core/glwindow"
	"github.com/Qendolin/go-printpixel/pkg/scene"
	"github.com/Qendolin/go-printpixel/renderer"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Window struct {
	GlWindow     glwindow.Extended
	Child        scene.Layoutable
	BeforeUpdate func()
	AfterUpdate  func()
	Renderers    map[string]renderer.Renderer
	//Top, Right, Bottom, Left
	margins   []int
	init      bool
	drawables map[string][]renderer.Drawable
}

func New(title string, cfg *SimpleConfig) (*Window, error) {
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

func NewCustom(title string, width, height int, hints glwindow.Hints, monitor *glfw.Monitor, glCfg glcontext.Config) (*Window, error) {
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

	win.Renderers = map[string]renderer.Renderer{
		renderer.TextureQuad: renderer.NewTextureQuadRenderer(),
	}

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

func (win *Window) SetWidth(w int) {
	_, h := win.GlWindow.GetSize()
	win.GlWindow.SetSize(w-win.margins[1]-win.margins[3], h)
}

func (win *Window) SetHeight(h int) {
	w, _ := win.GlWindow.GetSize()
	win.GlWindow.SetSize(w, h-win.margins[0]-win.margins[2])
}

func (win *Window) Width() int {
	w, _ := win.GlWindow.GetSize()
	return w + win.margins[1] + win.margins[3]
}

func (win *Window) Height() int {
	_, h := win.GlWindow.GetSize()
	return h + win.margins[0] + win.margins[2]
}

func (win *Window) SetInnerX(x int) {
	_, y := win.GlWindow.GetPos()
	win.GlWindow.SetPos(x, y)
}

func (win *Window) SetInnerY(y int) {
	x, _ := win.GlWindow.GetPos()
	win.GlWindow.SetPos(x, y)
}

func (win *Window) InnerX() int {
	x, _ := win.GlWindow.GetPos()
	return x
}

func (win *Window) InnerY() int {
	_, y := win.GlWindow.GetPos()
	return y
}

func (win *Window) SetInnerWidth(w int) {
	_, h := win.GlWindow.GetSize()
	win.GlWindow.SetSize(w, h)
}

func (win *Window) SetInnerHeight(h int) {
	w, _ := win.GlWindow.GetSize()
	win.GlWindow.SetSize(w, h)
}

func (win *Window) InnerWidth() int {
	w, _ := win.GlWindow.GetSize()
	return w
}

func (win *Window) InnerHeight() int {
	_, h := win.GlWindow.GetSize()
	return h
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

	for key, drawables := range win.drawables {
		renderer := win.Renderers[key]
		renderer.SetScale(2/float32(win.InnerWidth()), 2/float32(win.InnerHeight()))
		renderer.Bind()
		renderer.Draw(drawables...)
		renderer.Unbind()
	}

	win.GlWindow.SwapBuffers()
	glfw.PollEvents()
	if win.AfterUpdate != nil {
		win.AfterUpdate()
	}
}

func (win *Window) Layout() []scene.Layoutable {
	if win.Child == nil {
		return nil
	}
	win.Child.SetX(0)
	win.Child.SetY(0)
	win.Child.SetWidth(win.InnerWidth())
	win.Child.SetHeight(win.InnerHeight())

	win.drawables = map[string][]renderer.Drawable{}

	if l, ok := win.Child.(scene.Layouter); ok {
		tree := scene.Layout(l)
		for _, node := range tree.Nodes {
			if d, ok := node.Value.(renderer.Drawable); ok {
				if bucket, ok := win.drawables[d.GetRenderer()]; ok {
					bucket = append(bucket, d)
					win.drawables[d.GetRenderer()] = bucket
				} else {
					win.drawables[d.GetRenderer()] = []renderer.Drawable{d}
				}
			}
		}
	} else if d, ok := win.Child.(renderer.Drawable); ok {
		win.drawables[d.GetRenderer()] = []renderer.Drawable{d}
	}

	return nil
}
