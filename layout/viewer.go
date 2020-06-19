package layout

import "github.com/Qendolin/go-printpixel/internal/canvas"

type Viewer struct {
	Canvas canvas.Canvas
	SimpleBox
}

func (v Viewer) Draw() {
	v.Canvas.Bind()
	v.Canvas.Draw()
	v.Canvas.Unbind()
}
