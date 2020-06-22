package layout

import (
	"github.com/Qendolin/go-printpixel/internal/data"
	"github.com/Qendolin/go-printpixel/internal/renderer"
	"github.com/go-gl/gl/v3.3-core/gl"
)

var _viewCopyer *renderer.TextureCopy

type Viewer struct {
	Target renderer.TextureQuad
	SimpleBox
}

func NewViewer() *Viewer {
	return &Viewer{
		Target: *renderer.NewTextureQuad(),
	}
}

func viewCopyer() *renderer.TextureCopy {
	if _viewCopyer == nil {
		_viewCopyer = renderer.NewTextureCopy()
	}
	return _viewCopyer
}

func (v *Viewer) Layout() {
	vc := viewCopyer()

	vc.SetSource(*v.Target.Texture)

	dst := data.NewTexture2D(data.Tex2DTarget2D)
	dst.Bind(0)
	dst.AllocEmpty(0, gl.RGB, int32(v.width), int32(v.height), gl.RGB)
	dst.ApplyDefaults()
	vc.SetDestination(*dst)
	vc.Bind()
	vc.Draw()
	vc.Unbind()
}

func (v Viewer) Draw() {
	v.Target.Bind()
	v.Target.Draw()
	v.Target.Unbind()
}
