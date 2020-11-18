package w32

import "unsafe"

type XFORM struct {
	M11, M12, M21, M22, Dx, Dy float32
}

var (
	setWindowExtEx        = gdi32.NewProc("SetWindowExtEx")
	setWindowOrgEx        = gdi32.NewProc("SetWindowOrgEx")
	setViewportExtEx      = gdi32.NewProc("SetViewportExtEx")
	setViewportOrgEx      = gdi32.NewProc("SetViewportOrgEx")
	setMapMode            = gdi32.NewProc("SetMapMode")
	setPolyFillMode       = gdi32.NewProc("SetPolyFillMode")
	setTextAlign          = gdi32.NewProc("SetTextAlign")
	saveDC                = gdi32.NewProc("SaveDC")
	restoreDC             = gdi32.NewProc("RestoreDC")
	setWorldTransform     = gdi32.NewProc("SetWorldTransform")
	modifyWorldTransform  = gdi32.NewProc("ModifyWorldTransform")
	beginPath             = gdi32.NewProc("BeginPath")
	endPath               = gdi32.NewProc("EndPath")
	closeFigure           = gdi32.NewProc("CloseFigure")
	fillPath              = gdi32.NewProc("FillPath")
	strokeAndFillPath     = gdi32.NewProc("StrokeAndFillPath")
	strokePath            = gdi32.NewProc("StrokePath")
	createPenIndirect     = gdi32.NewProc("CreatePenIndirect")
	createFontIndirectW   = gdi32.NewProc("CreateFontIndirectW")
	createFontIndirectExW = gdi32.NewProc("CreateFontIndirectExW")
)

func SetWindowExtEx(hdc HDC, x, y int, lpSize *POINT) bool {
	ret, _, _ := setWindowExtEx.Call(
		uintptr(hdc),
		uintptr(x),
		uintptr(y),
		uintptr(unsafe.Pointer(lpSize)),
	)
	return ret != 0
}

func SetWindowOrgEx(hdc HDC, x, y int, lpPoint *POINT) bool {
	ret, _, _ := setWindowOrgEx.Call(
		uintptr(hdc),
		uintptr(x),
		uintptr(y),
		uintptr(unsafe.Pointer(lpPoint)),
	)
	return ret != 0
}

func SetViewportExtEx(hdc HDC, x, y int, lpSize *POINT) bool {
	ret, _, _ := setViewportExtEx.Call(
		uintptr(hdc),
		uintptr(x),
		uintptr(y),
		uintptr(unsafe.Pointer(lpSize)),
	)
	return ret != 0
}

func SetViewportOrgEx(hdc HDC, x, y int, lpPont *POINT) bool {
	ret, _, _ := setViewportOrgEx.Call(
		uintptr(hdc),
		uintptr(x),
		uintptr(y),
		uintptr(unsafe.Pointer(lpPont)),
	)
	return ret != 0
}

func SetMapMode(hdc HDC, iMode int) int {
	ret, _, _ := setMapMode.Call(
		uintptr(hdc),
		uintptr(iMode),
	)
	return int(ret)
}

func SetPolyFillMode(hdc HDC, mode int) int {
	ret, _, _ := setPolyFillMode.Call(
		uintptr(hdc),
		uintptr(mode),
	)
	return int(ret)
}

func SetTextAlign(hdc HDC, align uint) uint {
	ret, _, _ := setTextAlign.Call(
		uintptr(hdc),
		uintptr(align),
	)
	return uint(ret)
}

func SaveDC(hdc HDC) int {
	ret, _, _ := saveDC.Call(
		uintptr(hdc),
	)
	return int(ret)
}

func RestoreDC(hdc HDC, nSavedDC int) bool {
	ret, _, _ := restoreDC.Call(
		uintptr(hdc),
		uintptr(nSavedDC),
	)
	return ret != 0
}

func SetWorldTransform(hdc HDC, lpxf *XFORM) bool {
	ret, _, _ := setWorldTransform.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(lpxf)),
	)
	return ret != 0
}

func ModifyWorldTransform(hdc HDC, lpxf *XFORM, mode DWORD) bool {
	ret, _, _ := modifyWorldTransform.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(lpxf)),
		uintptr(mode),
	)
	return ret != 0
}

func BeginPath(hdc HDC) bool {
	ret, _, _ := beginPath.Call(
		uintptr(hdc),
	)
	return ret != 0
}

func EndPath(hdc HDC) bool {
	ret, _, _ := endPath.Call(
		uintptr(hdc),
	)
	return ret != 0
}

func CloseFigure(hdc HDC) bool {
	ret, _, _ := closeFigure.Call(
		uintptr(hdc),
	)
	return ret != 0
}

func FillPath(hdc HDC) bool {
	ret, _, _ := fillPath.Call(
		uintptr(hdc),
	)
	return ret != 0
}

func StrokeAndFillPath(hdc HDC) bool {
	ret, _, _ := strokeAndFillPath.Call(
		uintptr(hdc),
	)
	return ret != 0
}

func StrokePath(hdc HDC) bool {
	ret, _, _ := strokePath.Call(
		uintptr(hdc),
	)
	return ret != 0
}

func CreatePenIndirect(plpen *LOGPEN) HPEN {
	ret, _, _ := createPenIndirect.Call(
		uintptr(unsafe.Pointer(plpen)),
	)
	return HPEN(ret)
}

func CreateFontIndirectW(lplf *LOGFONT) HFONT {
	ret, _, _ := createFontIndirectW.Call(
		uintptr(unsafe.Pointer(lplf)),
	)
	return HFONT(ret)
}
