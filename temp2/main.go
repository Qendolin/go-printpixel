package main

import (
	"image"
	"log"
	"os"
	"time"

	"github.com/Qendolin/go-printpixel/core/glcontext"
	"github.com/Qendolin/go-printpixel/core/glwindow"
	"github.com/Qendolin/go-printpixel/core/shader"
	"github.com/Qendolin/go-printpixel/pkg/window"
	"github.com/Qendolin/go-printpixel/temp2/text3d"
	"github.com/Qendolin/go-printpixel/temp2/text3d/mesh2"
	"github.com/go-gl/gl/v3.3-compatibility/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var shittyKeyState map[glfw.Key]bool = make(map[glfw.Key]bool)

func shittyIsKeyPress(win *window.Window, key glfw.Key) bool {
	last, ok := shittyKeyState[key]
	state := win.GlWindow.GetKey(key) == glfw.Press
	shittyKeyState[key] = state
	if !ok {
		return state
	}

	if !last && state {
		return true
	}
	return false
}

func main() {
	win := setup()
	loadStart := time.Now()
	// imgPath := "./assets/a-msdf.png"
	// imgPath := "./assets/b-msdf.png"
	// imgPath := "./assets/g-msdf.png"
	// imgPath := "./assets/star-msdf.png"
	// imgPath := "./assets/Playball-W-msdf.png"
	// imgPath := "./assets/PressStart2P-W-msdf.png"
	imgPath := "./assets/Cook-T-msdf.png"
	// imgPath := "./assets/Regular-MSDF.png"
	// imgPath := "./assets/march-test.png"
	imgFile, err := os.Open(imgPath)
	panicIf(err)
	img, _, err := image.Decode(imgFile)
	panicIf(err)
	mesh, bounds, err := text3d.Load3d(img, 1, mesh2.DetailReducer{Detail: 0.06}, [3]float32{0.5, 0.5, 0.}, 10)
	// mesh, bounds, err := text3d.Load2d(imgPath, 2, mesh2.DetailReducer{Detail: 0.06}, [3]float32{0.5, 0.5, 0.5}, [3]float32{0, 0, 1})
	panicIf(err)
	log.Printf("Loaded in: %v\n", time.Since(loadStart))
	log.Printf("%-5d vertices, %-5d indices total\n", mesh.VertexCount, mesh.IndexCount)

	win.GlWindow.MakeContextCurrent()
	panicIf(gl.Init())

	gl.ClearColor(0.05, 0.05, 0.05, 1.)
	gl.PointSize(4)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	sh, err := shader.NewProgramFromPaths(
		"@mod/assets/shaders/transform.vert",
		"@mod/assets/shaders/flat.frag",
	)
	panicIf(err)
	if ok, info := sh.Validate(); ok {
		log.Printf("Program Validation Log: \n\n%v\n\n", info)
	} else {
		log.Fatalf("Program Validation Log: \n\n%v\n\n", info)
	}
	shDbg, err := shader.NewProgramFromPaths(
		"@mod/assets/shaders/transform.vert",
		"@mod/assets/shaders/color.frag",
	)
	panicIf(err)
	if ok, info := shDbg.Validate(); ok {
		log.Printf("Program Validation Log: \n\n%v\n\n", info)
	} else {
		log.Fatalf("Program Validation Log: \n\n%v\n\n", info)
	}

	sh.Bind()
	uTransform := sh.MustGetUniform("u_transform_mat")
	uModelMat := sh.MustGetUniform("u_model_mat")
	uColor := shDbg.MustGetUniform("u_color")

	// proj := mgl32.Ortho(-float32(bounds.Dx()+5)/2, float32(bounds.Dx()+5)/2, -float32(bounds.Dy()+5)/2, float32(bounds.Dy()+5)/2, -1, 1)
	proj := mgl32.Perspective(mgl32.DegToRad(90), win.InnerAspect(), 0.1, 10000)
	pos := mgl32.Vec3{float32(bounds.Dx()) / 2, -float32(bounds.Dy()) / 2, float32(bounds.Dx()) * 0.6}
	viewMat := mgl32.Translate3D(-pos[0], -pos[1], -pos[2])
	scale := float32(1)
	angle := float32(0)
	scaleMat := mgl32.Scale3D(scale, scale, scale)
	rotMat := mgl32.HomogRotate3DY(angle)
	speed := float32(bounds.Dx()) * 0.75

	win.GlWindow.SetFramebufferSizeCallback(func(w glwindow.Extended, width, height int) {
		proj = mgl32.Perspective(mgl32.DegToRad(90), float32(width)/float32(height), 0.1, 10000)
		glwindow.ResizeGlViewport(w, width, height)
	})

	var (
		wireframe = false
		points    = false
		culling   = true
		depth     = true
	)

	for !win.GlWindow.ShouldClose() {
		if shittyIsKeyPress(win, glfw.KeyP) {
			points = !points
		}
		if shittyIsKeyPress(win, glfw.KeyC) {
			culling = !culling
			if culling {
				gl.Enable(gl.CULL_FACE)
			} else {
				gl.Disable(gl.CULL_FACE)
			}
		}
		if shittyIsKeyPress(win, glfw.KeyO) {
			wireframe = !wireframe
		}
		if shittyIsKeyPress(win, glfw.KeyY) {
			depth = !depth
			if depth {
				gl.Enable(gl.DEPTH_TEST)
			} else {
				gl.Disable(gl.DEPTH_TEST)

			}
		}
		if glfw.GetCurrentContext().GetKey(glfw.KeyW) == glfw.Press {
			pos[2] -= speed / scale * float32(win.GlWindow.Delta().Seconds())
			viewMat = mgl32.Translate3D(-pos[0], -pos[1], -pos[2])
		}
		if glfw.GetCurrentContext().GetKey(glfw.KeyS) == glfw.Press {
			pos[2] += speed / scale * float32(win.GlWindow.Delta().Seconds())
			viewMat = mgl32.Translate3D(-pos[0], -pos[1], -pos[2])
		}
		if glfw.GetCurrentContext().GetKey(glfw.KeyA) == glfw.Press {
			pos[0] -= speed / scale * float32(win.GlWindow.Delta().Seconds())
			viewMat = mgl32.Translate3D(-pos[0], -pos[1], -pos[2])
		}
		if glfw.GetCurrentContext().GetKey(glfw.KeyD) == glfw.Press {
			pos[0] += speed / scale * float32(win.GlWindow.Delta().Seconds())
			viewMat = mgl32.Translate3D(-pos[0], -pos[1], -pos[2])
		}
		if glfw.GetCurrentContext().GetKey(glfw.KeySpace) == glfw.Press {
			pos[1] += speed / scale * float32(win.GlWindow.Delta().Seconds())
			viewMat = mgl32.Translate3D(-pos[0], -pos[1], -pos[2])
		}
		if glfw.GetCurrentContext().GetKey(glfw.KeyLeftControl) == glfw.Press {
			pos[1] -= speed / scale * float32(win.GlWindow.Delta().Seconds())
			viewMat = mgl32.Translate3D(-pos[0], -pos[1], -pos[2])
		}

		if glfw.GetCurrentContext().GetKey(glfw.KeyQ) == glfw.Press {
			angle += float32(win.GlWindow.Delta().Seconds())
			rotMat = mgl32.HomogRotate3DY(angle)
		}
		if glfw.GetCurrentContext().GetKey(glfw.KeyE) == glfw.Press {
			angle -= float32(win.GlWindow.Delta().Seconds())
			rotMat = mgl32.HomogRotate3DY(angle)
		}

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		sh.Bind()

		modelMat := rotMat.Mul4(scaleMat)
		uModelMat.Set(modelMat)
		uTransform.Set(proj.Mul4(viewMat).Mul4(modelMat))

		mesh.Vao.Bind()
		gl.DrawElements(gl.TRIANGLES, int32(mesh.IndexCount), gl.UNSIGNED_INT, nil)
		if points {
			shDbg.Bind()
			uModelMat.Set(modelMat)
			uTransform.Set(proj.Mul4(viewMat).Mul4(modelMat))
			uColor.Set(mgl32.Vec4{1, 0, 0, 1})
			gl.DepthFunc(gl.ALWAYS)
			gl.DrawElements(gl.POINTS, int32(mesh.IndexCount), gl.UNSIGNED_INT, nil)
			gl.DepthFunc(gl.LESS)
		}
		if wireframe {
			shDbg.Bind()
			uModelMat.Set(modelMat)
			uTransform.Set(proj.Mul4(viewMat).Mul4(modelMat))
			uColor.Set(mgl32.Vec4{0, 0, 1, 1})
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
			gl.DepthFunc(gl.ALWAYS)
			gl.DrawElements(gl.TRIANGLES, int32(mesh.IndexCount), gl.UNSIGNED_INT, nil)
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
			gl.DepthFunc(gl.LESS)
		}

		win.GlWindow.SwapBuffers()
		glfw.PollEvents()
	}
}

func setup() *window.Window {
	cfg := glcontext.NewGlConfig(1)
	cfg.Debug = true
	cfg.Multisampling = false
	hints := glwindow.NewHints()
	hints.ContextVersionMajor.Value = 3
	hints.ContextVersionMinor.Value = 3
	hints.OpenGLForwardCompatible.Value = true
	hints.OpenGLProfile.Value = glwindow.OpenGLCoreProfile
	hints.OpenGLDebugContext.Value = true
	win, err := window.NewCustom("MS Example", 900, 900, hints, nil, cfg)
	panicIf(err)
	win.GlWindow.MakeContextCurrent()

	panicIf(gl.Init())

	go handleErrors(cfg.Errors)
	return win
}

func panicIf(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func handleErrors(errs <-chan glcontext.Error) {
	for err := range errs {
		if err.Fatal {
			log.Fatalf("%v\n%v", err, err.Stack)
		}
		log.Printf("%v\n", err)
	}
}
