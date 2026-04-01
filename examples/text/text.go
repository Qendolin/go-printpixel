package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "image/png"

	"github.com/Qendolin/go-printpixel/core"
	"github.com/Qendolin/go-printpixel/core/data"
	"github.com/Qendolin/go-printpixel/core/font"
	"github.com/Qendolin/go-printpixel/core/glw"
	"github.com/Qendolin/go-printpixel/core/shader"
	"github.com/Qendolin/go-printpixel/pkg/window"
	"github.com/Qendolin/go-printpixel/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const Width = 800
const Height = 1000

func loadFontMesh(text []rune, fnt *font.Font, style font.Style) {
	verts, texs := font.Mesh(text, fnt, style)
	vao := data.NewVao(nil)
	vao.Bind()
	vBuf := data.NewBuffer(nil, data.BufVertexAttribute)
	vBuf.Bind()
	vBuf.WriteStatic(verts)
	vao.MustLayout(0, 2, float32(0), false, 0, 0)
	tBuf := data.NewBuffer(nil, data.BufVertexAttribute)
	tBuf.Bind()
	tBuf.WriteStatic(texs)
	vao.MustLayout(1, 2, float32(0), false, 0, 0)
}

func main() {
	win := setup()
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	fntdef, err := os.Open("Regular.fnt")
	panicIf(err)
	fnt, err := font.Parse(fntdef)
	panicIf(err)

	fontSize := float32(68 * 96. / 72.)
	style := font.Style{
		Size:       fontSize,
		Kerning:    true,
		LineHeight: 1,
		TabSize:    4,
	}

	text := font.Layout([]rune(
		"Hello, world!\nThis is a test.\nLVA Ta gj fi\n"+
			"This line is really long, it should be broken up.\n"+
			"AndAWordThatsTooLong. Also, ellipsis is coming, and you can't see this part."), fnt, font.LayoutSpecs{
		Ellipsis: true,
		Height:   Height,
		Width:    Width,
	}, style)

	loadFontMesh(text, fnt, style)

	pagePaths := make([]string, len(fnt.Pages))
	fntDir := utils.MustResolvePath("")
	for i, file := range fnt.Pages {
		pagePaths[i] = filepath.Join(fntDir, file)
	}
	bm, err := core.NewTexture3D(core.
		InitPaths(len(pagePaths), pagePaths[0], pagePaths[1:]...), data.RGBA8)
	panicIf(err)
	bm.Bind(0)

	prog := newProgram("./text.vert", "./text.frag")
	prog.Bind()

	fontScale := fontSize / float32(fnt.Size)
	sclX, sclY := win.GlWindow.GetContentScale()
	mProj := mgl32.Scale3D(2./Width*sclX, 2./Height*sclY, 1.)
	mModel := mgl32.Translate3D(-Width/2, Height/2, 0).Mul4(mgl32.Scale3D(fontScale, fontScale, 1.))
	mMP := mProj.Mul4(mModel)
	prog.MustGetUniform("u_transform").Set(mMP)
	prog.MustGetUniform("u_type").Set(2)

	for !win.GlWindow.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(6*len(text)))
		win.GlWindow.SwapBuffers()
		glfw.PollEvents()
	}
}

func newProgram(vsPath, fsPath string) *shader.Program {
	vs, err := shader.NewShaderFromPath(vsPath, shader.TypeVertex)
	panicIf(err)

	fs, err := shader.NewShaderFromPath(fsPath, shader.TypeFragment)
	panicIf(err)
	prog, err := shader.NewProgram(vs, fs)
	panicIf(err)
	fs.Destroy()
	vs.Destroy()

	if ok, log := prog.Validate(); ok {
		fmt.Printf("Program Validation Log: \n\n%v\n\n", log)
	} else {
		panicIf(fmt.Errorf("Program Validation Log: \n\n%v\n\n", log))
	}

	return prog
}

func setup() *window.Window {
	conf := window.SimpleConfig{
		Width:  Width,
		Height: Height,
		Debug:  true,
		Title:  "Text Example",
		DebugHandler: func(err glw.DebugMessage) {
			if err.Critical {
				log.Fatalf("%v\n%v", err, err.Stack)
			}
			log.Printf("%v\n", err)
		},
	}
	win, err := window.New(conf)
	panicIf(err)

	return win
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
