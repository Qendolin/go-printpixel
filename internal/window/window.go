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
}

type extWindow struct {
	*glfw.Window
	lastSwap time.Time
	cbs      callbacks
}

func (w *extWindow) GetGLFWWindow() *glfw.Window {
	return w.Window
}

func (w *extWindow) SwapBuffers() {
	w.lastSwap = time.Now()
	w.Window.SwapBuffers()
}

func (w *extWindow) Delta() time.Duration {
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

	return &extWindow{
		Window:   glfwWin,
		lastSwap: time.Now(),
	}, nil
}
