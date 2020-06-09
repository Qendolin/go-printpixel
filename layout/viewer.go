package layout

import "github.com/Qendolin/go-printpixel/internal/canvas"

type Viewer struct {
	cnv canvas.Canvas
	SimpleBox
}

func (v Viewer) Update() {
	v.cnv.Bind()
	v.cnv.Draw()
	v.cnv.Unbind()
}
