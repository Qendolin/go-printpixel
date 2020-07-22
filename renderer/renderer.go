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

type ZDrawable struct {
	Drawable
	Z int
}

func CalcModelProjectionMat(zd ZDrawable, s mgl32.Vec3) mgl32.Mat4 {
	t := zd.GetTransform()
	// This is the premultiplied result of an ortho. projection matrix from s, assuming near=0:
	// [s_x 0   0   0
	//  0   s_y 0   0
	//  0   0   s_z -1
	//  0   0   0   1 ]
	// And an 2d transformation matrix t that has been converted to 2d with z
	// [a b c
	//  d e f
	//  g h 1]
	// [a b 0 c
	//  d e 0 f
	//  0 0 0 -z
	//  g h 0 1 ]
	return mgl32.Mat4{
		t[0] * s[0], t[1] * s[0], 0, t[2] * s[0],
		t[3] * s[1], t[4] * s[1], 0, t[5] * s[1],
		-t[6], -t[7], 0, float32(-zd.Z)*s[2] - 1,
		t[6], t[7], 0, 1,
	}
}

type Renderer interface {
	SetScale(x, y, z float32)
	Bind()
	Unbind()
	BindFor(utils.BindingClosure)
	Draw(...ZDrawable)
}

type Base struct {
	Scale mgl32.Vec3
}

func (b *Base) SetScale(x, y, z float32) {
	b.Scale = mgl32.Vec3{x, y, z}
}
