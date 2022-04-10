package glw

import (
	"image"
	"log"
	"math/big"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var id uint64 = 0

func NewId() uint64 {
	id++
	return id
}

var (
	windows     = sync.Map{}
	windowCount int32
)

func Get(id uint64) Window {
	w, ok := windows.Load(id)
	if !ok {
		return nil
	}
	return (w).(Window)
}

func Put(win Window) {
	windows.Store(win.Id(), win)
	atomic.AddInt32(&windowCount, 1)
}

func Remove(win Window) {
	windows.Delete(win.Id())
}

func DestroyAll() {
	glfw.Terminate()
	windows = sync.Map{}
}

type Window interface {
	Destroy()
	GetAttrib(attrib glfw.Hint) int
	GetClipboardString() string
	GetContentScale() (x float32, y float32)
	GetCursorPos() (x, y float64)
	// GetMarginFrameSize returns the frame size including invisible margins (that are used to resize)
	GetMarginFrameSize() (left, top, right, bottom int)
	GetFrameSize() (left, top, right, bottom int)
	GetFramebufferSize() (width, height int)
	GetInputMode(mode glfw.InputMode) int
	GetKey(key glfw.Key) glfw.Action
	GetMonitor() *glfw.Monitor
	GetMouseButton(button glfw.MouseButton) glfw.Action
	GetOpacity() float32
	GetPos() (x, y int)
	GetSize() (width, height int)
	GetWidth() int
	GetHeight() int
	GetUserPointer() unsafe.Pointer
	Handle() unsafe.Pointer
	Hide()
	Iconify()
	MakeContextCurrent()
	Maximize()
	RequestAttention()
	Restore()
	SetAspectRatio(numer, denom int)
	SetAttrib(attrib glfw.Hint, value int)
	SetCharCallback(cbfun CharCallback) (previous CharCallback)
	SetCharModsCallback(cbfun CharModsCallback) (previous CharModsCallback)
	SetClipboardString(str string)
	SetCloseCallback(cbfun CloseCallback) (previous CloseCallback)
	SetContentScaleCallback(cbfun ContentScaleCallback) ContentScaleCallback
	SetCursor(c *glfw.Cursor)
	SetCursorEnterCallback(cbfun CursorEnterCallback) (previous CursorEnterCallback)
	SetCursorPos(xpos, ypos float64)
	SetCursorPosCallback(cbfun CursorPosCallback) (previous CursorPosCallback)
	SetDropCallback(cbfun DropCallback) (previous DropCallback)
	SetFocusCallback(cbfun FocusCallback) (previous FocusCallback)
	SetFramebufferSizeCallback(cbfun FramebufferSizeCallback) (previous FramebufferSizeCallback)
	SetIcon(images []image.Image)
	SetIconifyCallback(cbfun IconifyCallback) (previous IconifyCallback)
	SetInputMode(mode glfw.InputMode, value int)
	SetKeyCallback(cbfun KeyCallback) (previous KeyCallback)
	SetMaximizeCallback(cbfun MaximizeCallback) MaximizeCallback
	SetMonitor(monitor *glfw.Monitor, xpos, ypos, width, height, refreshRate int)
	SetFullscreen(monitor *glfw.Monitor)
	SetMouseButtonCallback(cbfun MouseButtonCallback) (previous MouseButtonCallback)
	SetOpacity(opacity float32)
	SetPos(xpos, ypos int)
	SetPosCallback(cbfun PosCallback) (previous PosCallback)
	SetRefreshCallback(cbfun RefreshCallback) (previous RefreshCallback)
	SetScrollCallback(cbfun ScrollCallback) (previous ScrollCallback)
	SetShouldClose(value bool)
	Close()
	SetSize(width, height int)
	SetSizeCallback(cbfun SizeCallback) (previous SizeCallback)
	SetSizeLimits(minw, minh, maxw, maxh int)
	SetTitle(title string)
	SetUserPointer(pointer unsafe.Pointer)
	ShouldClose() bool
	Show()
	SwapBuffers()
	Delta() time.Duration
	GetGLFWWindow() *glfw.Window
	Id() uint64
	HasProblem(problem) bool
}

type problem string

const (
	CannotMaximize problem = "the window cannot be maximized, this is likely a system bug"
	NoDebugContext problem = "a debug context was requested but not granted, it is likely not supported by the device"
)

type extWindow struct {
	*glfw.Window
	lastSwap     time.Time
	d            time.Duration
	cbs          callbacks
	id           uint64
	problems     map[problem]bool
	debugHandler DebugHandler
}

func (w *extWindow) GetGLFWWindow() *glfw.Window {
	return w.Window
}

func (w *extWindow) GetMarginFrameSize() (left, top, right, bottom int) {
	return w.Window.GetFrameSize()
}

func (w *extWindow) SwapBuffers() {
	w.d = time.Since(w.lastSwap)
	w.lastSwap = time.Now()
	w.Window.SwapBuffers()
}

func (w *extWindow) Delta() time.Duration {
	return w.d
}

func (w *extWindow) Id() uint64 {
	return w.id
}

func (w *extWindow) HasProblem(p problem) bool {
	return w.problems[p]
}

func (w *extWindow) Close() {
	w.SetShouldClose(true)
}

func (w *extWindow) GetWidth() int {
	width, _ := w.GetSize()
	return width
}

func (w *extWindow) GetHeight() int {
	_, height := w.GetSize()
	return height
}

func (w *extWindow) SetFullscreen(mon *glfw.Monitor) {
	vidMode := mon.GetVideoMode()
	w.SetMonitor(mon, 0, 0, vidMode.Width, vidMode.Height, 0)
}

// Destroy destroys the specified window and its context. On calling this function, no further callbacks will be called for that window.
// Will also call glfw.Terminate() when the last window is destroyed.
func (w *extWindow) Destroy() {
	runtime.LockOSThread()
	w.Window.Destroy()
	Remove(w)
	runtime.UnlockOSThread()
	w.id = 0

	atomic.AddInt32(&windowCount, -1)
	if windowCount == 0 {
		DestroyAll()
	}
}

func New(conf Config) (Window, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	err := glfw.Init()
	if err != nil {
		return nil, err
	}

	glfw.DefaultWindowHints()
	conf.load()

	glfwWin, err := glfw.CreateWindow(conf.Width, conf.Height, conf.Title, conf.Monitor, conf.SharedContext)
	if err != nil {
		return nil, err
	}
	glfwWin.MakeContextCurrent()
	x := &extWindow{
		Window:       glfwWin,
		lastSwap:     time.Now(),
		id:           NewId(),
		problems:     map[problem]bool{},
		debugHandler: conf.DebugHandler,
	}
	id := x.id
	glfwWin.SetUserPointer(gl.Ptr(&id))
	Put(x)
	if err = gl.Init(); err != nil {
		return nil, err
	}

	applyExtendedConfig(x, conf)

	if conf.Maximized && glfwWin.GetAttrib(glfw.Maximized) != glfw.True {
		// on some systems it's not possible to maximize the window
		x.problems[CannotMaximize] = true
	}
	var contextFlags int32
	gl.GetIntegerv(gl.CONTEXT_FLAGS, &contextFlags)
	if conf.DebugContext && (contextFlags&gl.CONTEXT_FLAG_DEBUG_BIT) == 0 {
		// debug contexts (GL_DEBUG_OUTPUT etc.) are only core since version 4.3
		x.problems[NoDebugContext] = true
	}

	x.SetFramebufferSizeCallback(ResizeGlViewport)

	return x, nil
}

func applyExtendedConfig(win *extWindow, conf Config) {
	win.SetSizeLimits(conf.MinWidth, conf.MinHeight, conf.MaxWidth, conf.MaxHeight)
	var aspectNumer, aspectDenom int
	if conf.AspectRatio <= 0 {
		aspectDenom = DontCare
		aspectNumer = DontCare
	} else {
		aspect := new(big.Rat).SetFloat64(conf.AspectRatio)
		aspectDenom = int(aspect.Denom().Int64())
		aspectNumer = int(aspect.Num().Int64())
	}
	if !conf.Maximized {
		win.SetAspectRatio(aspectNumer, aspectDenom)
		win.SetSize(conf.Width, conf.Height)
		win.SetPos(conf.X, conf.Y)
	}
	if conf.Vsync {
		glfw.SwapInterval(1)
	}
	if conf.DebugContext {
		var major, minor int32
		gl.GetIntegerv(gl.MAJOR_VERSION, &major)
		gl.GetIntegerv(gl.MINOR_VERSION, &minor)
		log.Printf(" === System Information === \n")
		log.Printf("OpenGL Version: %v (%v.%v)\n", gl.GoStr(gl.GetString(gl.VERSION)), major, minor)
		log.Printf("GLSL Version: %v\n", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))
		log.Printf("Renderer: %v\n", gl.GoStr(gl.GetString(gl.RENDERER)))
		log.Printf("Vendor: %v\n", gl.GoStr(gl.GetString(gl.VENDOR)))
		log.Println()
		gl.DebugMessageCallback(DefaultDebugMessageCallback(win), nil)
		gl.Enable(gl.DEBUG_OUTPUT)
		gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)
	}
}

func ResizeGlViewport(win Window, w, h int) {
	gl.Viewport(0, 0, int32(w), int32(h))
}
