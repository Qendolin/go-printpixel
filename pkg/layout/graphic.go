package layout

import (
	"github.com/Qendolin/go-printpixel/core/data"
	"github.com/Qendolin/go-printpixel/renderer"
)

type Drawable interface {
	TextureQuad() renderer.TextureQuad
}

type Graphic struct {
	Texture data.Texture2D
	q       renderer.TextureQuad
	b       SimpleBox
}

func NewGraphic() *Graphic {
	return &Graphic{
		Texture: *data.NewTexture2D(data.Tex2DTarget2D),
	}
}

func (g *Graphic) TextureQuad() renderer.TextureQuad {
	g.q.Texture = &g.Texture
	return g.q
}

func (g *Graphic) SetWidth(w int) {
	g.b.width = w
	g.q.Transform.Set(0, 0, float32(w))
}

func (g *Graphic) Width() int {
	return g.b.width
}

func (g *Graphic) SetHeight(h int) {
	g.b.height = h
	g.q.Transform.Set(1, 1, float32(h))
}

func (g *Graphic) Height() int {
	return g.b.height
}

func (g *Graphic) SetX(x int) {
	g.b.x = x
	g.q.Transform.Set(2, 0, float32(x))
}

func (g *Graphic) X() int {
	return g.b.x
}

func (g *Graphic) SetY(y int) {
	g.b.y = y
	g.q.Transform.Set(2, 1, float32(y))
}

func (g *Graphic) Y() int {
	return g.b.y
}
