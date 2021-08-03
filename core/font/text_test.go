package font_test

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	_ "image/png"

	"github.com/Qendolin/go-printpixel/core"
	"github.com/Qendolin/go-printpixel/core/data"
	"github.com/Qendolin/go-printpixel/core/font"
	"github.com/Qendolin/go-printpixel/core/test"
	"github.com/Qendolin/go-printpixel/utils"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	test.ParseArgs()
	m.Run()
}

const pt float32 = 96. / 72.
const (
	Width  = 800
	Height = 1000
)

func TestFontWrap(t *testing.T) {
	win, close := test.NewWindow(t, "70747db82b76c2cdae3d56af16c338")
	win.SetSize(Width, Height)
	defer close()
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	fntdef, err := os.Open(utils.MustResolvePath("@mod/assets/fonts/Go/Regular.fnt"))
	require.NoError(t, err)
	fnt, err := font.Parse(fntdef)
	require.NoError(t, err)

	fontSize := 68 * pt
	style := font.Style{
		Size:       fontSize,
		Kerning:    true,
		LineHeight: 1,
		TabSize:    4,
	}

	text := font.Layout([]rune(
		"Hello, world!\nThis is a test.\nLVA Ta gj\n"+
			"This line is really long, it should be broken up.\n"+
			"AndAWordThatsTooLong. Also, ellipsis is coming, can't see this."), fnt, font.LayoutSpecs{
		Ellipsis: true,
		Height:   Height,
		Width:    Width,
	}, style)

	loadFontMesh(text, fnt, style)

	pagePaths := make([]string, len(fnt.Pages))
	fntDir := utils.MustResolvePath("@mod/assets/fonts/Go/")
	for i, file := range fnt.Pages {
		pagePaths[i] = filepath.Join(fntDir, file)
	}
	bm, err := core.NewTexture3D(core.
		InitPaths(len(pagePaths), pagePaths[0], pagePaths[1:]...), data.RGBA8)
	require.NoError(t, err)
	bm.Bind(0)

	prog := test.NewProgram(t, "@mod/assets/shaders/text_test.vert", "@mod/assets/shaders/text_test.frag")
	prog.Bind()

	fontScale := fontSize / float32(fnt.Size)
	sclX, sclY := win.GetContentScale()
	mProj := mgl32.Scale3D(2./Width*sclX, 2./Height*sclY, 1.)
	mModel := mgl32.Translate3D(-Width/2, Height/2, 0).Mul4(mgl32.Scale3D(fontScale, fontScale, 1.))
	mMP := mProj.Mul4(mModel)
	prog.MustGetUniform("u_transform").Set(mMP)

	for !win.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(6*len(text)))
		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func TestTabs(t *testing.T) {
	win, close := test.NewWindow(t, "100292524a010020040080100400")
	win.SetSize(Width, Height)
	defer close()
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	fntdef, err := os.Open(utils.MustResolvePath("@mod/assets/fonts/Go/Regular.fnt"))
	require.NoError(t, err)
	fnt, err := font.Parse(fntdef)
	require.NoError(t, err)

	fontSize := 72 * pt
	style := font.Style{
		Size:       fontSize,
		Kerning:    true,
		LineHeight: 1,
		TabSize:    4,
	}

	text := font.Layout([]rune("|\n\t|\n \t|\n  \t|\n   \t|\n    \t|\n|\t|\t|\t|\t|\t|\t|\t_\t|"), fnt, font.LayoutSpecs{
		Ellipsis: true,
		Height:   Height,
		Width:    Width,
	}, style)

	loadFontMesh(text, fnt, style)

	pagePaths := make([]string, len(fnt.Pages))
	fntDir := utils.MustResolvePath("@mod/assets/fonts/Go/")
	for i, file := range fnt.Pages {
		pagePaths[i] = filepath.Join(fntDir, file)
	}
	bm, err := core.NewTexture3D(core.
		InitPaths(len(pagePaths), pagePaths[0], pagePaths[1:]...), data.RGBA8)
	require.NoError(t, err)
	bm.Bind(0)

	prog := test.NewProgram(t, "@mod/assets/shaders/text_test.vert", "@mod/assets/shaders/text_test.frag")
	prog.Bind()

	fontScale := fontSize / float32(fnt.Size)
	sclX, sclY := win.GetContentScale()
	mProj := mgl32.Scale3D(2./Width*sclX, 2./Height*sclY, 1.)
	mModel := mgl32.Translate3D(-Width/2, Height/2, 0).Mul4(mgl32.Scale3D(fontScale, fontScale, 1.))
	mMP := mProj.Mul4(mModel)
	prog.MustGetUniform("u_transform").Set(mMP)

	for !win.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(6*len(text)))
		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func TestBM(t *testing.T) {
	win, close := test.NewWindow(t, "15ca9aac901327aab510")
	win.SetSize(Width, Height)
	defer close()
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	fntdef, err := os.Open(utils.MustResolvePath("@mod/assets/fonts/Go/Regular.fnt"))
	require.NoError(t, err)
	fnt, err := font.Parse(fntdef)
	require.NoError(t, err)

	fontSize := 200 * pt
	style := font.Style{
		Size:       fontSize,
		Kerning:    true,
		LineHeight: 1,
		TabSize:    4,
	}

	text := font.Layout([]rune("Hello World!"), fnt, font.LayoutSpecs{
		Ellipsis: true,
		Height:   Height,
		Width:    Width,
	}, style)
	loadFontMesh(text, fnt, style)

	pagePaths := make([]string, len(fnt.Pages))
	fntDir := utils.MustResolvePath("@mod/assets/fonts/Go/")
	for i, file := range fnt.Pages {
		pagePaths[i] = filepath.Join(fntDir, file)
	}
	bm, err := core.NewTexture3D(core.
		InitPaths(len(pagePaths), pagePaths[0], pagePaths[1:]...), data.RGBA8)
	require.NoError(t, err)
	bm.Bind(0)

	prog := test.NewProgram(t, "@mod/assets/shaders/text_test.vert", "@mod/assets/shaders/text_test.frag")
	prog.Bind()

	fontScale := fontSize / float32(fnt.Size)
	sclX, sclY := win.GetContentScale()
	mProj := mgl32.Scale3D(2./Width*sclX, 2./Height*sclY, 1.)
	mSdfModel := mgl32.Translate3D(-Width/2, Height/2, 0).Mul4(mgl32.Scale3D(fontScale, fontScale, 1.))
	mSdfMP := mProj.Mul4(mSdfModel)

	prog.MustGetUniform("u_transform").Set(mSdfMP)

	for !win.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(6*len(text)))
		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func TestSDF(t *testing.T) {
	win, close := test.NewWindow(t, "5cabaac9d1325a6f510")
	win.SetSize(Width, Height)
	defer close()
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	fntdef, err := os.Open(utils.MustResolvePath("@mod/assets/fonts/Go/Regular-SDF.fnt"))
	require.NoError(t, err)
	fnt, err := font.Parse(fntdef)
	require.NoError(t, err)

	fontSize := 200 * pt
	style := font.Style{
		Size:       fontSize,
		Kerning:    true,
		LineHeight: 1,
		TabSize:    4,
	}

	text := font.Layout([]rune("Hello World!"), fnt, font.LayoutSpecs{
		Ellipsis: true,
		Height:   Height,
		Width:    Width,
	}, style)
	loadFontMesh(text, fnt, style)

	pagePaths := make([]string, len(fnt.Pages))
	fntDir := "@mod/assets/fonts/Go"
	for i, file := range fnt.Pages {
		pagePaths[i] = path.Join(fntDir, file)
	}
	bm, err := core.NewTexture3D(core.InitPaths(len(pagePaths), pagePaths[0], pagePaths[1:]...), data.RGBA8)
	require.NoError(t, err)
	bm.Bind(0)

	prog := test.NewProgram(t, "@mod/assets/shaders/text_test.vert", "@mod/assets/shaders/text_test.frag")
	prog.Bind()

	fontScale := fontSize / float32(fnt.Size)
	sclX, sclY := win.GetContentScale()
	mProj := mgl32.Scale3D(2./Width*sclX, 2./Height*sclY, 1.)
	mModel := mgl32.Translate3D(-Width/2, Height/2, 0).Mul4(mgl32.Scale3D(fontScale, fontScale, 1.))
	mMP := mProj.Mul4(mModel)

	prog.MustGetUniform("u_type").Set(1)
	prog.MustGetUniform("u_transform").Set(mMP)

	for !win.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(6*len(text)))
		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func TestMSDF(t *testing.T) {
	win, close := test.NewWindow(t, "54abaac9d1325a6f510")
	win.SetSize(Width, Height)
	defer close()
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	fntdef, err := os.Open(utils.MustResolvePath("@mod/assets/fonts/Go/Regular-MSDF.fnt"))
	require.NoError(t, err)
	fnt, err := font.Parse(fntdef)
	require.NoError(t, err)

	fontSize := 200 * pt
	style := font.Style{
		Size:       fontSize,
		Kerning:    true,
		LineHeight: 1,
		TabSize:    4,
	}

	text := font.Layout([]rune("Hello World!"), fnt, font.LayoutSpecs{
		Ellipsis: true,
		Height:   Height,
		Width:    Width,
	}, style)
	loadFontMesh(text, fnt, style)

	pagePaths := make([]string, len(fnt.Pages))
	fntDir := utils.MustResolvePath("@mod/assets/fonts/Go/")
	for i, file := range fnt.Pages {
		pagePaths[i] = filepath.Join(fntDir, file)
	}
	bm, err := core.NewTexture3D(core.
		InitPaths(len(pagePaths), pagePaths[0], pagePaths[1:]...), data.RGBA8)
	require.NoError(t, err)
	bm.Bind(0)

	prog := test.NewProgram(t, "@mod/assets/shaders/text_test.vert", "@mod/assets/shaders/text_test.frag")
	prog.Bind()

	fontScale := fontSize / float32(fnt.Size)
	sclX, sclY := win.GetContentScale()
	mProj := mgl32.Scale3D(2./Width*sclX, 2./Height*sclY, 1.)
	mModel := mgl32.Translate3D(-Width/2, Height/2, 0).Mul4(mgl32.Scale3D(fontScale, fontScale, 1.))
	mMP := mProj.Mul4(mModel)

	prog.MustGetUniform("u_type").Set(2)
	prog.MustGetUniform("u_transform").Set(mMP)

	for !win.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(6*len(text)))
		win.SwapBuffers()
		glfw.PollEvents()
	}
}

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
