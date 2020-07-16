package renderer

import (
	"github.com/Qendolin/go-printpixel/core/data"
	"github.com/Qendolin/go-printpixel/utils"
	"github.com/go-gl/mathgl/mgl32"
)

type Drawable interface {
	GetMesh() *data.Vao
	GetTextures() []data.GLTexture
	GetTransform() mgl32.Mat3
	GetRenderer() string
}

type Renderer interface {
	SetScale(x, y float32)
	Bind()
	Unbind()
	BindFor(utils.BindingClosure)
	Draw(...Drawable)
}

type Base struct {
	scale mgl32.Vec3
}

func (b *Base) SetScale(x, y float32) {
	b.scale = mgl32.Vec3{x, y, 1}
}
