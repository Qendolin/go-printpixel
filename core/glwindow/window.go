package glwindow

import (
	"image"
	"time"
	"unsafe"

	"github.com/Qendolin/go-printpixel/core/glcontext"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var id uint64 = 0

func NewId() uint64 {
	id++
	return id
}

var windows = map[uint64]Extended{}

func Get(id uint64) Extended {
	return windows[id]
}

func Put(win Extended) {
	windows[win.Id()] = win
}

type Extended interface {
	Destroy()
	GetAttrib(attrib glfw.Hint) int
	GetClipboardString() string
	GetContentScale() (x float32, y float32)
	GetCursorPos() (x, y float64)
	GetFrameSize() (left, top, right, bottom int)
	GetVisibleFrameSize() (left, top, right, bottom int) //TODO: rename
	GetFramebufferSize() (width, height int)
	GetInputMode(mode glfw.InputMode) int
	GetKey(key glfw.Key) glfw.Action
	GetMonitor() *glfw.Monitor
	GetMouseButton(button glfw.MouseButton) glfw.Action
	GetOpacity() float32
	GetPos() (x, y int)
	GetSize() (width, height int)
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
	SetMouseButtonCallback(cbfun MouseButtonCallback) (previous MouseButtonCallback)
	SetOpacity(opacity float32)
	SetPos(xpos, ypos int)
	SetPosCallback(cbfun PosCallback) (previous PosCallback)
	SetRefreshCallback(cbfun RefreshCallback) (previous RefreshCallback)
	SetScrollCallback(cbfun ScrollCallback) (previous ScrollCallback)
	SetShouldClose(value bool)
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
)

type extWindow struct {
	*glfw.Window
	lastSwap time.Time
	d        time.Duration
	cbs      callbacks
	id       uint64
	problems map[problem]bool
}

func (w *extWindow) GetGLFWWindow() *glfw.Window {
	return w.Window
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

func New(hints Hints, title string, width, height int, monitor *glfw.Monitor) (Extended, error) {
	if glcontext.Status()&glcontext.StatusGlfwInitialized == 0 {
		return nil, glcontext.ErrGlfwNotInitialized
	}

	glfw.DefaultWindowHints()
	hints.apply()

	if monitor == nil && (hints.Fullscreen.Value) {
		monitor = glfw.GetPrimaryMonitor()
	}

	glfwWin, err := glfw.CreateWindow(width, height, title, monitor, nil)
	if err != nil {
		return nil, err
	}
	glfwWin.MakeContextCurrent()
	x := &extWindow{
		Window:   glfwWin,
		lastSwap: time.Now(),
		id:       NewId(),
	}
	id := x.id
	glfwWin.SetUserPointer(gl.Ptr(&id))
	Put(x)

	if hints.Vsync.Value {
		glfw.SwapInterval(1)
	}
	if hints.Maximized.Value && glfwWin.GetAttrib(glfw.Maximized) != glfw.True {
		// on some systems it's not possible to maximize the window
		x.problems[CannotMaximize] = true
	}

	glfwWin.SetFramebufferSizeCallback(func(_ *glfw.Window, w, h int) {
		if glcontext.Status()&glcontext.StatusGlInitialized != 0 {
			gl.Viewport(0, 0, int32(w), int32(h))
		}
	})

	return x, nil
}
