package glwindow

import (
	"fmt"
	"syscall"
	"unsafe"
)

type rect struct {
	left   int32
	top    int32
	right  int32
	bottom int32
}

const DWMWA_EXTENDED_FRAME_BOUNDS = 0x9

var (
	dwmapi                = syscall.MustLoadDLL("Dwmapi")
	dwmGetWindowAttribute = dwmapi.MustFindProc("DwmGetWindowAttribute")
)

func (win extWindow) GetFrameSize() (left, top, right, bottom int) {
	hwndC := win.Window.GetWin32Window()
	hwnd := *((*uint64)(unsafe.Pointer(&hwndC)))

	r := rect{}

	ret, _, err := dwmGetWindowAttribute.Call(uintptr(hwnd), uintptr(DWMWA_EXTENDED_FRAME_BOUNDS), uintptr(unsafe.Pointer(&r)), unsafe.Sizeof(r))
	if ret != 0 {
		if err != nil {
			panic(err)
		} else {
			panic(fmt.Sprintf("%s returned non zero result: %X", dwmGetWindowAttribute.Name, ret))
		}
	}

	x, y := win.GetPos()
	w, h := win.GetSize()

	return x - int(r.left), y - int(r.top), int(r.right) - x - w, int(r.bottom) - y - h
}
