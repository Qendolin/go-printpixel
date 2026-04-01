package scene

import (
	"github.com/Qendolin/go-printpixel/core"
	"github.com/Qendolin/go-printpixel/core/data"
	"github.com/Qendolin/go-printpixel/renderer"
	"github.com/go-gl/mathgl/mgl32"
)

type Graphic struct {
	Texture *data.Texture2D
	Alpha   bool
	b       Box
	t       mgl32.Mat3
	q       *data.Vao
}

func NewGraphic() *Graphic {
	return &Graphic{
		Texture: data.NewTexture2D(nil, data.Tex2DTarget2D),
	}
}

func LoadGraphic(path0 string, paths ...string) (g *Graphic) {
	return &Graphic{
		Texture: core.MustNewTexture2D(core.InitPaths(0, path0, paths...), data.ColorFormatDefault),
	}
}

func (g *Graphic) SetWidth(w int) {
	g.b.width = w
	g.t.Set(0, 0, float32(w))
}

func (g *Graphic) Width() int {
	return g.b.width
}

func (g *Graphic) SetHeight(h int) {
	g.b.height = h
	g.t.Set(1, 1, float32(h))
}

func (g *Graphic) Height() int {
	return g.b.height
}

func (g *Graphic) SetX(x int) {
	g.b.x = x
	g.t.Set(2, 0, float32(x))
}

func (g *Graphic) X() int {
	return g.b.x
}

func (g *Graphic) SetY(y int) {
	g.b.y = y
	g.t.Set(2, 1, float32(y))
}

func (g *Graphic) Y() int {
	return g.b.y
}

func (g *Graphic) GetMesh() *data.Vao {
	if g.q == nil {
		g.q = core.Quad()
	}
	return g.q
}

func (g *Graphic) GetTextures() []data.GLTexture {
	return []data.GLTexture{g.Texture.GLTexture}
}

func (g *Graphic) GetTransform() mgl32.Mat3 {
	return g.t
}

func (g *Graphic) GetRenderer() string {
	return renderer.TextureQuad
}

func (g *Graphic) HasAlpha() bool {
	return g.Alpha
}
