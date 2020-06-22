package window

import (
	"image"
	"time"
	"unsafe"

	"github.com/Qendolin/go-printpixel/internal/context"
	"github.com/go-gl/glfw/v3.3/glfw"
)

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
	SetCharCallback(cbfun glfw.CharCallback) (previous glfw.CharCallback)
	SetCharModsCallback(cbfun glfw.CharModsCallback) (previous glfw.CharModsCallback)
	SetClipboardString(str string)
	SetCloseCallback(cbfun glfw.CloseCallback) (previous glfw.CloseCallback)
	SetContentScaleCallback(cbfun glfw.ContentScaleCallback) glfw.ContentScaleCallback
	SetCursor(c *glfw.Cursor)
	SetCursorEnterCallback(cbfun glfw.CursorEnterCallback) (previous glfw.CursorEnterCallback)
	SetCursorPos(xpos, ypos float64)
	SetCursorPosCallback(cbfun glfw.CursorPosCallback) (previous glfw.CursorPosCallback)
	SetDropCallback(cbfun glfw.DropCallback) (previous glfw.DropCallback)
	SetFocusCallback(cbfun glfw.FocusCallback) (previous glfw.FocusCallback)
	SetFramebufferSizeCallback(cbfun glfw.FramebufferSizeCallback) (previous glfw.FramebufferSizeCallback)
	SetIcon(images []image.Image)
	SetIconifyCallback(cbfun glfw.IconifyCallback) (previous glfw.IconifyCallback)
	SetInputMode(mode glfw.InputMode, value int)
	SetKeyCallback(cbfun glfw.KeyCallback) (previous glfw.KeyCallback)
	SetMaximizeCallback(cbfun glfw.MaximizeCallback) glfw.MaximizeCallback
	SetMonitor(monitor *glfw.Monitor, xpos, ypos, width, height, refreshRate int)
	SetMouseButtonCallback(cbfun glfw.MouseButtonCallback) (previous glfw.MouseButtonCallback)
	SetOpacity(opacity float32)
	SetPos(xpos, ypos int)
	SetPosCallback(cbfun glfw.PosCallback) (previous glfw.PosCallback)
	SetRefreshCallback(cbfun glfw.RefreshCallback) (previous glfw.RefreshCallback)
	SetScrollCallback(cbfun glfw.ScrollCallback) (previous glfw.ScrollCallback)
	SetShouldClose(value bool)
	SetSize(width, height int)
	SetSizeCallback(cbfun glfw.SizeCallback) (previous glfw.SizeCallback)
	SetSizeLimits(minw, minh, maxw, maxh int)
	SetTitle(title string)
	SetUserPointer(pointer unsafe.Pointer)
	ShouldClose() bool
	Show()
	SwapBuffers()
	Delta() time.Duration
	GetGLFWWindow() *glfw.Window
}

type glfwWindow struct {
	*glfw.Window
	lastSwap time.Time
}

func (w *glfwWindow) GetGLFWWindow() *glfw.Window {
	return w.Window
}

func (w *glfwWindow) SwapBuffers() {
	w.lastSwap = time.Now()
	w.Window.SwapBuffers()
}

func (w *glfwWindow) Delta() time.Duration {
	return time.Since(w.lastSwap)
}

func New(hints Hints, title string, width, height int, monitor *glfw.Monitor) (Extended, error) {

	if context.Status()&context.StatusGlfwInitialized == 0 {
		return nil, context.ErrGlfwNotInitialized
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

	if hints.Vsync.Value {
		glfw.SwapInterval(1)
	}

	return &glfwWindow{
		Window:   glfwWin,
		lastSwap: time.Now(),
	}, nil
}
