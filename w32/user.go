package w32

import (
	"golang.org/x/sys/windows"
)

var (
	user32 = windows.NewLazySystemDLL("user32.dll")

	getDC            = user32.NewProc("GetDC")
	getDesktopWindow = user32.NewProc("GetDesktopWindow")
)

func GetDC(hWnd HWND) HDC {
	ret, _, _ := getDC.Call(
		uintptr(hWnd),
	)
	return HDC(ret)
}

func GetDesktopWindow() HWND {
	ret, _, _ := getDesktopWindow.Call()
	return HWND(ret)
}
