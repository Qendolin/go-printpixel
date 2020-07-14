package glwindow

import "github.com/go-gl/glfw/v3.3/glfw"

type CharCallback func(w Extended, char rune)
type CharModsCallback func(w Extended, char rune, mods glfw.ModifierKey)
type CloseCallback func(w Extended)
type ContentScaleCallback func(w Extended, x float32, y float32)
type CursorEnterCallback func(w Extended, entered bool)
type CursorPosCallback func(w Extended, x float64, y float64)
type DropCallback func(w Extended, names []string)
type FocusCallback func(w Extended, focused bool)
type FramebufferSizeCallback func(w Extended, width int, height int)
type IconifyCallback func(w Extended, iconified bool)
type KeyCallback func(w Extended, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey)
type MaximizeCallback func(w Extended, maximized bool)
type MouseButtonCallback func(w Extended, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey)
type PosCallback func(w Extended, x int, y int)
type RefreshCallback func(w Extended)
type ScrollCallback func(w Extended, xoff float64, yoff float64)
type SizeCallback func(w Extended, width int, height int)

type callbacks struct {
	char            CharCallback
	charMods        CharModsCallback
	close           CloseCallback
	contentScale    ContentScaleCallback
	cursorEnter     CursorEnterCallback
	cursorPos       CursorPosCallback
	drop            DropCallback
	focus           FocusCallback
	framebufferSize FramebufferSizeCallback
	iconify         IconifyCallback
	key             KeyCallback
	maximize        MaximizeCallback
	mouseButton     MouseButtonCallback
	pos             PosCallback
	refresh         RefreshCallback
	scroll          ScrollCallback
	size            SizeCallback
}

func (w *extWindow) SetCharCallback(cb CharCallback) (previous CharCallback) {
	p := w.cbs.char
	w.cbs.char = cb
	w.GetGLFWWindow().SetCharCallback(func(_ *glfw.Window, char rune) {
		w.cbs.char(w, char)
	})
	return p
}

func (w *extWindow) SetCharModsCallback(cb CharModsCallback) (previous CharModsCallback) {
	p := w.cbs.charMods
	w.cbs.charMods = cb
	w.GetGLFWWindow().SetCharModsCallback(func(_ *glfw.Window, char rune, mods glfw.ModifierKey) {
		w.cbs.charMods(w, char, mods)
	})
	return p
}

func (w *extWindow) SetCloseCallback(cb CloseCallback) (previous CloseCallback) {
	p := w.cbs.close
	w.cbs.close = cb
	w.GetGLFWWindow().SetCloseCallback(func(_ *glfw.Window) {
		w.cbs.close(w)
	})
	return p
}

func (w *extWindow) SetContentScaleCallback(cb ContentScaleCallback) (previous ContentScaleCallback) {
	p := w.cbs.contentScale
	w.cbs.contentScale = cb
	w.GetGLFWWindow().SetContentScaleCallback(func(_ *glfw.Window, x float32, y float32) {
		w.cbs.contentScale(w, x, y)
	})
	return p
}

func (w *extWindow) SetCursorEnterCallback(cb CursorEnterCallback) (previous CursorEnterCallback) {
	p := w.cbs.cursorEnter
	w.cbs.cursorEnter = cb
	w.GetGLFWWindow().SetCursorEnterCallback(func(_ *glfw.Window, entered bool) {
		w.cbs.cursorEnter(w, entered)
	})
	return p
}

func (w *extWindow) SetCursorPosCallback(cb CursorPosCallback) (previous CursorPosCallback) {
	p := w.cbs.cursorPos
	w.cbs.cursorPos = cb
	w.GetGLFWWindow().SetCursorPosCallback(func(_ *glfw.Window, x float64, y float64) {
		w.cbs.cursorPos(w, x, y)
	})
	return p
}

func (w *extWindow) SetDropCallback(cb DropCallback) (previous DropCallback) {
	p := w.cbs.drop
	w.cbs.drop = cb
	w.GetGLFWWindow().SetDropCallback(func(_ *glfw.Window, names []string) {
		w.cbs.drop(w, names)
	})
	return p
}

func (w *extWindow) SetFocusCallback(cb FocusCallback) (previous FocusCallback) {
	p := w.cbs.focus
	w.cbs.focus = cb
	w.GetGLFWWindow().SetFocusCallback(func(_ *glfw.Window, focused bool) {
		w.cbs.focus(w, focused)
	})
	return p
}

func (w *extWindow) SetFramebufferSizeCallback(cb FramebufferSizeCallback) (previous FramebufferSizeCallback) {
	p := w.cbs.framebufferSize
	w.cbs.framebufferSize = cb
	w.GetGLFWWindow().SetFramebufferSizeCallback(func(_ *glfw.Window, width int, height int) {
		w.cbs.framebufferSize(w, width, height)
	})
	return p
}

func (w *extWindow) SetIconifyCallback(cb IconifyCallback) (previous IconifyCallback) {
	p := w.cbs.iconify
	w.cbs.iconify = cb
	w.GetGLFWWindow().SetIconifyCallback(func(_ *glfw.Window, iconified bool) {
		w.cbs.iconify(w, iconified)
	})
	return p
}

func (w *extWindow) SetKeyCallback(cb KeyCallback) (previous KeyCallback) {
	p := w.cbs.key
	w.cbs.key = cb
	w.GetGLFWWindow().SetKeyCallback(func(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		w.cbs.key(w, key, scancode, action, mods)
	})
	return p
}

func (w *extWindow) SetMaximizeCallback(cb MaximizeCallback) (previous MaximizeCallback) {
	p := w.cbs.maximize
	w.cbs.maximize = cb
	w.GetGLFWWindow().SetMaximizeCallback(func(_ *glfw.Window, maximized bool) {
		w.cbs.maximize(w, maximized)
	})
	return p
}

func (w *extWindow) SetMouseButtonCallback(cb MouseButtonCallback) (previous MouseButtonCallback) {
	p := w.cbs.mouseButton
	w.cbs.mouseButton = cb
	w.GetGLFWWindow().SetMouseButtonCallback(func(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		w.cbs.mouseButton(w, button, action, mods)
	})
	return p
}

func (w *extWindow) SetPosCallback(cb PosCallback) (previous PosCallback) {
	p := w.cbs.pos
	w.cbs.pos = cb
	w.GetGLFWWindow().SetPosCallback(func(_ *glfw.Window, x int, y int) {
		w.cbs.pos(w, x, y)
	})
	return p
}

func (w *extWindow) SetRefreshCallback(cb RefreshCallback) (previous RefreshCallback) {
	p := w.cbs.refresh
	w.cbs.refresh = cb
	w.GetGLFWWindow().SetRefreshCallback(func(_ *glfw.Window) {
		w.cbs.refresh(w)
	})
	return p
}

func (w *extWindow) SetScrollCallback(cb ScrollCallback) (previous ScrollCallback) {
	p := w.cbs.scroll
	w.cbs.scroll = cb
	w.GetGLFWWindow().SetScrollCallback(func(_ *glfw.Window, xoff float64, yoff float64) {
		w.cbs.scroll(w, xoff, yoff)
	})
	return p
}

func (w *extWindow) SetSizeCallback(cb SizeCallback) (previous SizeCallback) {
	p := w.cbs.size
	w.cbs.size = cb
	w.GetGLFWWindow().SetSizeCallback(func(_ *glfw.Window, width int, height int) {
		w.cbs.size(w, width, height)
	})
	return p
}
