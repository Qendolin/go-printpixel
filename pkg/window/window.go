package window

import (
	"runtime"

	"github.com/Qendolin/go-printpixel/core/glcontext"
	"github.com/Qendolin/go-printpixel/core/glwindow"
	"github.com/Qendolin/go-printpixel/pkg/scene"
	"github.com/Qendolin/go-printpixel/renderer"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Window struct {
	GlWindow     glwindow.Extended
	Child        scene.Layoutable
	BeforeUpdate func()
	AfterUpdate  func()
	Renderers    map[string]renderer.Renderer
	// Top, Right, Bottom, Left
	margins        []int
	init           bool
	drawables      []map[string][]renderer.ZDrawable
	alphaDrawables []map[string][]renderer.ZDrawable
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

	mLeft, mTop, mRight, mBot := glfwWin.GetFrameSize()

	win := Window{GlWindow: glfwWin, margins: []int{mTop, mRight, mBot, mLeft}}

	runtime.LockOSThread()
	win.GlWindow.MakeContextCurrent()
	err = glcontext.InitGl(glCfg)

	win.Renderers = map[string]renderer.Renderer{
		renderer.TextureQuad: renderer.NewTextureQuadRenderer(),
		renderer.Debug:       renderer.NewDebugRenderer(),
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

func (win *Window) InnerAspect() float32 {
	return float32(win.InnerWidth()) / float32(win.InnerHeight())
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

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	var depth int
	if dl, al := len(win.drawables), len(win.alphaDrawables); dl > al {
		depth = dl
	} else {
		depth = al
	}

	for d := range win.drawables {
		for key, drawables := range win.drawables[d] {
			renderer := win.Renderers[key]
			renderer.SetScale(2/float32(win.InnerWidth()), 2/float32(win.InnerHeight()), -2/float32(depth))
			renderer.Bind()
			renderer.Draw(drawables...)
			renderer.Unbind()
		}
	}

	gl.Enable(gl.BLEND)
	for d := range win.alphaDrawables {
		for key, drawables := range win.alphaDrawables[len(win.alphaDrawables)-d-1] {
			renderer := win.Renderers[key]
			renderer.SetScale(2/float32(win.InnerWidth()), 2/float32(win.InnerHeight()), -2/float32(depth))
			renderer.Bind()
			renderer.Draw(drawables...)
			renderer.Unbind()
		}
	}
	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.BLEND)

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

	win.drawables = make([]map[string][]renderer.ZDrawable, 0)
	win.alphaDrawables = make([]map[string][]renderer.ZDrawable, 0)

	if l, ok := win.Child.(scene.Layouter); ok {
		tree := scene.Layout(l)
		for _, node := range tree.Nodes {
			var d renderer.Drawable
			if d, ok = node.Value.(renderer.Drawable); !ok {
				continue
			}
			ds := &win.drawables
			if d.HasAlpha() {
				ds = &win.alphaDrawables
			}

			if node.Depth >= len(*ds) {
				// expand drawables
				ds1 := make([]map[string][]renderer.ZDrawable, node.Depth+1)
				copy(ds1, *ds)
				for i := range ds1 {
					if ds1[i] == nil {
						ds1[i] = make(map[string][]renderer.ZDrawable)
					}
				}
				*ds = ds1
			}

			zd := renderer.ZDrawable{Drawable: d, Z: node.Depth}
			if bucket, ok := (*ds)[node.Depth][d.GetRenderer()]; ok {
				bucket = append(bucket, zd)
				(*ds)[node.Depth][d.GetRenderer()] = bucket
			} else {
				(*ds)[node.Depth][d.GetRenderer()] = []renderer.ZDrawable{zd}
			}
		}
	} else if d, ok := win.Child.(renderer.Drawable); ok {
		if d.HasAlpha() {
			win.alphaDrawables = []map[string][]renderer.ZDrawable{{d.GetRenderer(): {{Drawable: d, Z: 0}}}}
		} else {
			win.drawables = []map[string][]renderer.ZDrawable{{d.GetRenderer(): {{Drawable: d, Z: 0}}}}
		}
	}

	return nil
}
