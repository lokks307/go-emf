package w32

import (
	"strconv"
	"syscall"
	"unsafe"
)

var (
	gdi32 = syscall.NewLazyDLL("gdi32.dll")

	getDeviceCaps             = gdi32.NewProc("GetDeviceCaps")
	deleteObject              = gdi32.NewProc("DeleteObject")
	createFontIndirectW       = gdi32.NewProc("CreateFontIndirectW")
	abortDoc                  = gdi32.NewProc("AbortDoc")
	bitBlt                    = gdi32.NewProc("BitBlt")
	maskBlt                   = gdi32.NewProc("MaskBlt")
	patBlt                    = gdi32.NewProc("PatBlt")
	closeEnhMetaFile          = gdi32.NewProc("CloseEnhMetaFile")
	copyEnhMetaFile           = gdi32.NewProc("CopyEnhMetaFileW")
	createBrushIndirect       = gdi32.NewProc("CreateBrushIndirect")
	createCompatibleDC        = gdi32.NewProc("CreateCompatibleDC")
	createCompatibleBitmap    = gdi32.NewProc("CreateCompatibleBitmap")
	createBitmap              = gdi32.NewProc("CreateBitmap")
	createDC                  = gdi32.NewProc("CreateDCW")
	createDIBSection          = gdi32.NewProc("CreateDIBSection")
	createEnhMetaFile         = gdi32.NewProc("CreateEnhMetaFileW")
	createIC                  = gdi32.NewProc("CreateICW")
	deleteDC                  = gdi32.NewProc("DeleteDC")
	deleteEnhMetaFile         = gdi32.NewProc("DeleteEnhMetaFile")
	ellipse                   = gdi32.NewProc("Ellipse")
	endDoc                    = gdi32.NewProc("EndDoc")
	endPage                   = gdi32.NewProc("EndPage")
	extCreatePen              = gdi32.NewProc("ExtCreatePen")
	getEnhMetaFile            = gdi32.NewProc("GetEnhMetaFileW")
	getEnhMetaFileHeader      = gdi32.NewProc("GetEnhMetaFileHeader")
	getObject                 = gdi32.NewProc("GetObjectW")
	getStockObject            = gdi32.NewProc("GetStockObject")
	getTextExtentExPoint      = gdi32.NewProc("GetTextExtentExPointW")
	getTextExtentPoint32      = gdi32.NewProc("GetTextExtentPoint32W")
	getTextMetrics            = gdi32.NewProc("GetTextMetricsW")
	lineTo                    = gdi32.NewProc("LineTo")
	moveToEx                  = gdi32.NewProc("MoveToEx")
	playEnhMetaFile           = gdi32.NewProc("PlayEnhMetaFile")
	rectangle                 = gdi32.NewProc("Rectangle")
	resetDC                   = gdi32.NewProc("ResetDCW")
	selectObject              = gdi32.NewProc("SelectObject")
	setBkMode                 = gdi32.NewProc("SetBkMode")
	setBrushOrgEx             = gdi32.NewProc("SetBrushOrgEx")
	setStretchBltMode         = gdi32.NewProc("SetStretchBltMode")
	setTextColor              = gdi32.NewProc("SetTextColor")
	setBkColor                = gdi32.NewProc("SetBkColor")
	startDoc                  = gdi32.NewProc("StartDocW")
	startPage                 = gdi32.NewProc("StartPage")
	stretchBlt                = gdi32.NewProc("StretchBlt")
	setDIBitsToDevice         = gdi32.NewProc("SetDIBitsToDevice")
	choosePixelFormat         = gdi32.NewProc("ChoosePixelFormat")
	describePixelFormat       = gdi32.NewProc("DescribePixelFormat")
	getEnhMetaFilePixelFormat = gdi32.NewProc("GetEnhMetaFilePixelFormat")
	getPixelFormat            = gdi32.NewProc("GetPixelFormat")
	setPixelFormat            = gdi32.NewProc("SetPixelFormat")
	setPixelV                 = gdi32.NewProc("SetPixelV")
	swapBuffers               = gdi32.NewProc("SwapBuffers")
	textOutW                  = gdi32.NewProc("TextOutW")
	createSolidBrush          = gdi32.NewProc("CreateSolidBrush")
	getDIBits                 = gdi32.NewProc("GetDIBits")
	pie                       = gdi32.NewProc("Pie")
	setDCPenColor             = gdi32.NewProc("SetDCPenColor")
	setDCBrushColor           = gdi32.NewProc("SetDCBrushColor")
	createPen                 = gdi32.NewProc("CreatePen")
	arc                       = gdi32.NewProc("Arc")
	arcTo                     = gdi32.NewProc("ArcTo")
	angleArc                  = gdi32.NewProc("AngleArc")
	chord                     = gdi32.NewProc("Chord")
	polygon                   = gdi32.NewProc("Polygon")
	polyline                  = gdi32.NewProc("Polyline")
	polyBezier                = gdi32.NewProc("PolyBezier")
	intersectClipRect         = gdi32.NewProc("IntersectClipRect")
	selectClipRgn             = gdi32.NewProc("SelectClipRgn")
	createRectRgn             = gdi32.NewProc("CreateRectRgn")
	combineRgn                = gdi32.NewProc("CombineRgn")
	enumFontFamiliesEx        = gdi32.NewProc("EnumFontFamiliesExW")
	setWindowExtEx            = gdi32.NewProc("SetWindowExtEx")
	setWindowOrgEx            = gdi32.NewProc("SetWindowOrgEx")
	setViewportExtEx          = gdi32.NewProc("SetViewportExtEx")
	setViewportOrgEx          = gdi32.NewProc("SetViewportOrgEx")
	setMapMode                = gdi32.NewProc("SetMapMode")
	setPolyFillMode           = gdi32.NewProc("SetPolyFillMode")
	setTextAlign              = gdi32.NewProc("SetTextAlign")
	saveDC                    = gdi32.NewProc("SaveDC")
	restoreDC                 = gdi32.NewProc("RestoreDC")
	setWorldTransform         = gdi32.NewProc("SetWorldTransform")
	modifyWorldTransform      = gdi32.NewProc("ModifyWorldTransform")
	beginPath                 = gdi32.NewProc("BeginPath")
	endPath                   = gdi32.NewProc("EndPath")
	abortPath                 = gdi32.NewProc("AbortPath")
	closeFigure               = gdi32.NewProc("CloseFigure")
	fillPath                  = gdi32.NewProc("FillPath")
	strokeAndFillPath         = gdi32.NewProc("StrokeAndFillPath")
	strokePath                = gdi32.NewProc("StrokePath")
	createPenIndirect         = gdi32.NewProc("CreatePenIndirect")
	createFontIndirectExW     = gdi32.NewProc("CreateFontIndirectExW")
	extTextOutW               = gdi32.NewProc("ExtTextOutW")
	polyBezierTo              = gdi32.NewProc("PolyBezierTo")
	polylineTo                = gdi32.NewProc("PolylineTo")
	polyPolygon               = gdi32.NewProc("PolyPolygon")
	stretchDIBits             = gdi32.NewProc("StretchDIBits")
	setMapperFlags            = gdi32.NewProc("SetMapperFlags")
	setROP2                   = gdi32.NewProc("SetROP2")
	scaleWindowExtEx          = gdi32.NewProc("ScaleWindowExtEx")
	setMetaRgn                = gdi32.NewProc("SetMetaRgn")
	offsetClipRgn             = gdi32.NewProc("OffsetClipRgn")
	setTextJustification      = gdi32.NewProc("SetTextJustification")
	fillRgn                   = gdi32.NewProc("FillRgn")
	createRectRgnIndirect     = gdi32.NewProc("CreateRectRgnIndirect")
	setGraphicsMode           = gdi32.NewProc("SetGraphicsMode")
	createPalette             = gdi32.NewProc("CreatePalette")
	selectPalette             = gdi32.NewProc("SelectPalette")
	setBitmapBits             = gdi32.NewProc("SetBitmapBits")
	createDIBitmap            = gdi32.NewProc("CreateDIBitmap")
	setMiterLimit             = gdi32.NewProc("SetMiterLimit")
	extSelectClipRgn          = gdi32.NewProc("ExtSelectClipRgn")
)

func GetDeviceCaps(hdc HDC, index int) int {
	ret, _, _ := getDeviceCaps.Call(
		uintptr(hdc),
		uintptr(index),
	)
	return int(ret)
}

func DeleteObject(hObject HGDIOBJ) bool {
	ret, _, _ := deleteObject.Call(uintptr(hObject))
	return ret != 0
}

func CreateFontIndirectW(logFont *LOGFONT) HFONT {
	ret, _, _ := createFontIndirectW.Call(uintptr(unsafe.Pointer(logFont)))
	return HFONT(ret)
}

func CreateFontIndirectExW(logFontExDv *LOGFONTEXDV) HFONT {
	ret, _, _ := createFontIndirectExW.Call(uintptr(unsafe.Pointer(logFontExDv)))
	return HFONT(ret)
}

func AbortDoc(hdc HDC) int {
	ret, _, _ := abortDoc.Call(uintptr(hdc))
	return int(ret)
}

func BitBlt(hdc HDC, x, y, cx, cy int, hdcSrc HDC, x1, y1 int, rop DWORD) bool {
	ret, _, _ := bitBlt.Call(
		uintptr(hdc),
		uintptr(x),
		uintptr(y),
		uintptr(cx),
		uintptr(cy),
		uintptr(hdcSrc),
		uintptr(x1),
		uintptr(y1),
		uintptr(rop),
	)
	return ret != 0
}

func MaskBlt(
	dest HDC, destX, destY, destWidth, destHeight int,
	source HDC, sourceX, sourceY int,
	mask HBITMAP, maskX, maskY int,
	operation DWORD,
) bool {
	ret, _, _ := maskBlt.Call(
		uintptr(dest),
		uintptr(destX),
		uintptr(destY),
		uintptr(destWidth),
		uintptr(destHeight),
		uintptr(source),
		uintptr(sourceX),
		uintptr(sourceY),
		uintptr(mask),
		uintptr(maskX),
		uintptr(maskX),
		uintptr(operation),
	)
	return ret != 0
}

func PatBlt(hdc HDC, nXLeft, nYLeft, nWidth, nHeight int, rop DWORD) bool {
	ret, _, _ := patBlt.Call(
		uintptr(hdc),
		uintptr(nXLeft),
		uintptr(nYLeft),
		uintptr(nWidth),
		uintptr(nHeight),
		uintptr(rop),
	)
	return ret != 0
}

func CloseEnhMetaFile(hdc HDC) HENHMETAFILE {
	ret, _, _ := closeEnhMetaFile.Call(uintptr(hdc))
	return HENHMETAFILE(ret)
}

func CopyEnhMetaFile(hemfSrc HENHMETAFILE, lpszFile *uint16) HENHMETAFILE {
	ret, _, _ := copyEnhMetaFile.Call(
		uintptr(hemfSrc),
		uintptr(unsafe.Pointer(lpszFile)),
	)
	return HENHMETAFILE(ret)
}

func CreateBrushIndirect(lplb *LOGBRUSH) HBRUSH {
	ret, _, _ := createBrushIndirect.Call(uintptr(unsafe.Pointer(lplb)))
	return HBRUSH(ret)
}

func CreateCompatibleDC(hdc HDC) HDC {
	ret, _, _ := createCompatibleDC.Call(uintptr(hdc))
	return HDC(ret)
}

func CreateCompatibleBitmap(hdc HDC, width, height int) HBITMAP {
	ret, _, _ := createCompatibleBitmap.Call(
		uintptr(hdc),
		uintptr(width),
		uintptr(height),
	)
	return HBITMAP(ret)
}

func CreateBitmap(nWidth, nHeight int, nPlanes, nBitCount UINT, lpBits []byte) HBITMAP {
	if len(lpBits) == 0 {
		return HBITMAP(0)
	}

	ret, _, _ := createBitmap.Call(
		uintptr(nWidth),
		uintptr(nHeight),
		uintptr(nPlanes),
		uintptr(nBitCount),
		uintptr(unsafe.Pointer(&lpBits[0])),
	)
	return HBITMAP(ret)
}

func CreateDC(lpszDriver, lpszDevice, lpszOutput *uint16, lpInitData *DEVMODE) HDC {
	ret, _, _ := createDC.Call(
		uintptr(unsafe.Pointer(lpszDriver)),
		uintptr(unsafe.Pointer(lpszDevice)),
		uintptr(unsafe.Pointer(lpszOutput)),
		uintptr(unsafe.Pointer(lpInitData)),
	)
	return HDC(ret)
}

func CreateDIBSection(hdc HDC, pbmi *BITMAPINFO, iUsage uint, ppvBits *unsafe.Pointer, hSection HANDLE, dwOffset uint) HBITMAP {
	ret, _, _ := createDIBSection.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(pbmi)),
		uintptr(iUsage),
		uintptr(unsafe.Pointer(ppvBits)),
		uintptr(hSection),
		uintptr(dwOffset),
	)
	return HBITMAP(ret)
}

func CreateEnhMetaFile(hdcRef HDC, lpFilename *uint16, lpRect *RECT, lpDescription *uint16) HDC {
	ret, _, _ := createEnhMetaFile.Call(
		uintptr(hdcRef),
		uintptr(unsafe.Pointer(lpFilename)),
		uintptr(unsafe.Pointer(lpRect)),
		uintptr(unsafe.Pointer(lpDescription)),
	)
	return HDC(ret)
}

func CreateIC(lpszDriver, lpszDevice, lpszOutput *uint16, lpdvmInit *DEVMODE) HDC {
	ret, _, _ := createIC.Call(
		uintptr(unsafe.Pointer(lpszDriver)),
		uintptr(unsafe.Pointer(lpszDevice)),
		uintptr(unsafe.Pointer(lpszOutput)),
		uintptr(unsafe.Pointer(lpdvmInit)),
	)
	return HDC(ret)
}

func DeleteDC(hdc HDC) bool {
	ret, _, _ := deleteDC.Call(uintptr(hdc))
	return ret != 0
}

func DeleteEnhMetaFile(hemf HENHMETAFILE) bool {
	ret, _, _ := deleteEnhMetaFile.Call(uintptr(hemf))
	return ret != 0
}

func Ellipse(hdc HDC, nLeftRect, nTopRect, nRightRect, nBottomRect int) bool {
	ret, _, _ := ellipse.Call(
		uintptr(hdc),
		uintptr(nLeftRect),
		uintptr(nTopRect),
		uintptr(nRightRect),
		uintptr(nBottomRect),
	)
	return ret != 0
}

func EndDoc(hdc HDC) int {
	ret, _, _ := endDoc.Call(uintptr(hdc))
	return int(ret)
}

func EndPage(hdc HDC) int {
	ret, _, _ := endPage.Call(uintptr(hdc))
	return int(ret)
}

func ExtCreatePen(dwPenStyle, dwWidth DWORD, lplb *LOGBRUSH, dwStyleCount DWORD, lpStyle []DWORD) HPEN {

	if len(lpStyle) == 0 {
		lpStyle = append(lpStyle, DWORD(0))
	}

	ret, _, _ := extCreatePen.Call(
		uintptr(dwPenStyle),
		uintptr(dwWidth),
		uintptr(unsafe.Pointer(lplb)),
		uintptr(dwStyleCount),
		uintptr(unsafe.Pointer(&lpStyle[0])),
	)
	return HPEN(ret)
}

func GetEnhMetaFile(lpszMetaFile *uint16) HENHMETAFILE {
	ret, _, _ := getEnhMetaFile.Call(uintptr(unsafe.Pointer(lpszMetaFile)))
	return HENHMETAFILE(ret)
}

func GetEnhMetaFileHeader(hemf HENHMETAFILE, cbBuffer uint, lpemh *ENHMETAHEADER) uint {
	ret, _, _ := getEnhMetaFileHeader.Call(
		uintptr(hemf),
		uintptr(cbBuffer),
		uintptr(unsafe.Pointer(lpemh)),
	)
	return uint(ret)
}

func GetObject(hgdiobj HGDIOBJ, cbBuffer uintptr, lpvObject unsafe.Pointer) int {
	ret, _, _ := getObject.Call(
		uintptr(hgdiobj),
		uintptr(cbBuffer),
		uintptr(lpvObject),
	)
	return int(ret)
}

func GetStockObject(fnObject int) HGDIOBJ {
	ret, _, _ := getStockObject.Call(uintptr(fnObject))
	return HGDIOBJ(ret)
}

func GetTextExtentExPoint(hdc HDC, lpszStr *uint16, cchString, nMaxExtent int, lpnFit, alpDx *int, lpSize *SIZE) bool {
	ret, _, _ := getTextExtentExPoint.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(lpszStr)),
		uintptr(cchString),
		uintptr(nMaxExtent),
		uintptr(unsafe.Pointer(lpnFit)),
		uintptr(unsafe.Pointer(alpDx)),
		uintptr(unsafe.Pointer(lpSize)),
	)
	return ret != 0
}

func GetTextExtentPoint32(hdc HDC, text string) (SIZE, bool) {
	var s SIZE
	str, err := syscall.UTF16FromString(text)
	if err != nil {
		return s, false
	}
	ret, _, _ := getTextExtentPoint32.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(&str[0])),
		uintptr(len(str)-1), // -1 for the trailing '\0'
		uintptr(unsafe.Pointer(&s)),
	)
	return s, ret != 0
}

func GetTextMetrics(hdc HDC, lptm *TEXTMETRIC) bool {
	ret, _, _ := getTextMetrics.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(lptm)),
	)
	return ret != 0
}

func LineTo(hdc HDC, nXEnd, nYEnd int) bool {
	ret, _, _ := lineTo.Call(
		uintptr(hdc),
		uintptr(nXEnd),
		uintptr(nYEnd),
	)
	return ret != 0
}

func MoveToEx(hdc HDC, x, y int, lpPoint *POINT) bool {
	ret, _, _ := moveToEx.Call(
		uintptr(hdc),
		uintptr(x),
		uintptr(y),
		uintptr(unsafe.Pointer(lpPoint)),
	)
	return ret != 0
}

func PlayEnhMetaFile(hdc HDC, hemf HENHMETAFILE, lpRect *RECT) bool {
	ret, _, _ := playEnhMetaFile.Call(
		uintptr(hdc),
		uintptr(hemf),
		uintptr(unsafe.Pointer(lpRect)),
	)
	return ret != 0
}

func Rectangle(hdc HDC, nLeftRect, nTopRect, nRightRect, nBottomRect int) bool {
	ret, _, _ := rectangle.Call(
		uintptr(hdc),
		uintptr(nLeftRect),
		uintptr(nTopRect),
		uintptr(nRightRect),
		uintptr(nBottomRect),
	)
	return ret != 0
}

func ResetDC(hdc HDC, lpInitData *DEVMODE) HDC {
	ret, _, _ := resetDC.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(lpInitData)),
	)
	return HDC(ret)
}

func SelectObject(hdc HDC, hgdiobj HGDIOBJ) HGDIOBJ {
	ret, _, _ := selectObject.Call(
		uintptr(hdc),
		uintptr(hgdiobj),
	)
	return HGDIOBJ(ret)
}

func SetBkMode(hdc HDC, iBkMode int) int {
	ret, _, _ := setBkMode.Call(
		uintptr(hdc),
		uintptr(iBkMode),
	)
	return int(ret)
}

func SetBrushOrgEx(hdc HDC, nXOrg, nYOrg int, lppt *POINT) bool {
	ret, _, _ := setBrushOrgEx.Call(
		uintptr(hdc),
		uintptr(nXOrg),
		uintptr(nYOrg),
		uintptr(unsafe.Pointer(lppt)),
	)
	return ret != 0
}

func SetStretchBltMode(hdc HDC, iStretchMode int) int {
	ret, _, _ := setStretchBltMode.Call(
		uintptr(hdc),
		uintptr(iStretchMode),
	)
	return int(ret)
}

func SetTextColor(hdc HDC, crColor COLORREF) COLORREF {
	ret, _, _ := setTextColor.Call(
		uintptr(hdc),
		uintptr(crColor),
	)
	return COLORREF(ret)
}

func SetBkColor(hdc HDC, crColor COLORREF) COLORREF {
	ret, _, _ := setBkColor.Call(
		uintptr(hdc),
		uintptr(crColor),
	)
	return COLORREF(ret)
}

func StartDoc(hdc HDC, lpdi *DOCINFO) int {
	ret, _, _ := startDoc.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(lpdi)),
	)
	return int(ret)
}

func StartPage(hdc HDC) int {
	ret, _, _ := startPage.Call(uintptr(hdc))
	return int(ret)
}

func StretchBlt(hdcDest HDC, nXOriginDest, nYOriginDest, nWidthDest, nHeightDest int, hdcSrc HDC, nXOriginSrc, nYOriginSrc, nWidthSrc, nHeightSrc int, dwRop DWORD) bool {
	ret, _, _ := stretchBlt.Call(
		uintptr(hdcDest),
		uintptr(nXOriginDest),
		uintptr(nYOriginDest),
		uintptr(nWidthDest),
		uintptr(nHeightDest),
		uintptr(hdcSrc),
		uintptr(nXOriginSrc),
		uintptr(nYOriginSrc),
		uintptr(nWidthSrc),
		uintptr(nHeightSrc),
		uintptr(dwRop),
	)
	return ret != 0
}

func SetDIBitsToDevice(hdc HDC, xDest, yDest int, dwWidth, dwHeight DWORD, xSrc, ySrc int, uStartScan, cScanLines UINT, lpvBits []byte, lpbmi *BITMAPINFO, fuColorUse UINT) int {
	if len(lpvBits) == 0 {
		return 0
	}

	ret, _, _ := setDIBitsToDevice.Call(
		uintptr(hdc),
		uintptr(xDest),
		uintptr(yDest),
		uintptr(dwWidth),
		uintptr(dwHeight),
		uintptr(xSrc),
		uintptr(ySrc),
		uintptr(uStartScan),
		uintptr(cScanLines),
		uintptr(unsafe.Pointer(&lpvBits[0])),
		uintptr(unsafe.Pointer(lpbmi)),
		uintptr(fuColorUse),
	)
	return int(ret)
}

func ChoosePixelFormat(hdc HDC, pfd *PIXELFORMATDESCRIPTOR) int {
	ret, _, _ := choosePixelFormat.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(pfd)),
	)
	return int(ret)
}

func DescribePixelFormat(hdc HDC, iPixelFormat int, nBytes uint, pfd *PIXELFORMATDESCRIPTOR) int {
	ret, _, _ := describePixelFormat.Call(
		uintptr(hdc),
		uintptr(iPixelFormat),
		uintptr(nBytes),
		uintptr(unsafe.Pointer(pfd)),
	)
	return int(ret)
}

func GetEnhMetaFilePixelFormat(hemf HENHMETAFILE, cbBuffer uint32, pfd *PIXELFORMATDESCRIPTOR) uint {
	ret, _, _ := getEnhMetaFilePixelFormat.Call(
		uintptr(hemf),
		uintptr(cbBuffer),
		uintptr(unsafe.Pointer(pfd)),
	)
	return uint(ret)
}

func GetPixelFormat(hdc HDC) int {
	ret, _, _ := getPixelFormat.Call(uintptr(hdc))
	return int(ret)
}

func SetPixelFormat(hdc HDC, iPixelFormat int, pfd *PIXELFORMATDESCRIPTOR) bool {
	ret, _, _ := setPixelFormat.Call(
		uintptr(hdc),
		uintptr(iPixelFormat),
		uintptr(unsafe.Pointer(pfd)),
	)
	return ret == TRUE
}

func SwapBuffers(hdc HDC) bool {
	ret, _, _ := swapBuffers.Call(uintptr(hdc))
	return ret == TRUE
}

func TextOutW(hdc HDC, x, y int, s string) bool {
	str, err := syscall.UTF16FromString(s)
	if err != nil {
		return false
	}
	ret, _, _ := textOutW.Call(
		uintptr(hdc),
		uintptr(x),
		uintptr(y),
		uintptr(unsafe.Pointer(&str[0])),
		uintptr(len(str)-1), // -1 for the trailing '\0'
	)
	return ret != 0
}

func CreateSolidBrush(color uint32) HBRUSH {
	ret, _, _ := createSolidBrush.Call(uintptr(color))
	return HBRUSH(ret)
}

func GetDIBits(
	dc HDC,
	bmp HBITMAP,
	startScan, scanLines UINT,
	bits unsafe.Pointer,
	info *BITMAPINFO,
	usage UINT,
) int {
	ret, _, _ := getDIBits.Call(
		uintptr(dc),
		uintptr(bmp),
		uintptr(startScan),
		uintptr(scanLines),
		uintptr(bits),
		uintptr(unsafe.Pointer(info)),
		uintptr(usage),
	)
	return int(ret)
}

func Pie(hdc HDC, left, top, right, bottom, xr1, yr1, xr2, yr2 int) bool {
	ret, _, _ := pie.Call(
		uintptr(hdc),
		uintptr(left),
		uintptr(top),
		uintptr(right),
		uintptr(bottom),
		uintptr(xr1),
		uintptr(yr1),
		uintptr(xr2),
		uintptr(yr2),
	)
	return ret != 0
}

func SetDCPenColor(hdc HDC, color COLORREF) COLORREF {
	ret, _, _ := setDCPenColor.Call(uintptr(hdc), uintptr(color))
	return COLORREF(ret)
}

func SetDCBrushColor(hdc HDC, color COLORREF) COLORREF {
	ret, _, _ := setDCBrushColor.Call(uintptr(hdc), uintptr(color))
	return COLORREF(ret)
}

func CreatePen(style int, width int, color COLORREF) HPEN {
	ret, _, _ := createPen.Call(
		uintptr(style),
		uintptr(width),
		uintptr(color),
	)
	return HPEN(ret)
}

func Arc(hdc HDC, x1, y1, x2, y2, x3, y3, x4, y4 int) bool {
	ret, _, _ := arc.Call(
		uintptr(hdc),
		uintptr(x1),
		uintptr(y1),
		uintptr(x2),
		uintptr(y2),
		uintptr(x3),
		uintptr(y3),
		uintptr(x4),
		uintptr(y4),
	)
	return ret != 0
}

func ArcTo(hdc HDC, left, top, right, bottom, xr1, yr1, xr2, yr2 int) bool {
	ret, _, _ := arcTo.Call(
		uintptr(hdc),
		uintptr(left),
		uintptr(top),
		uintptr(right),
		uintptr(bottom),
		uintptr(xr1),
		uintptr(yr1),
		uintptr(xr2),
		uintptr(yr2),
	)
	return ret != 0
}

func AngleArc(hdc HDC, x, y, r int, startAngle, sweepAngle float32) bool {
	ret, _, _ := angleArc.Call(
		uintptr(hdc),
		uintptr(x),
		uintptr(y),
		uintptr(r),
		uintptr(startAngle),
		uintptr(sweepAngle),
	)
	return ret != 0
}

func Chord(hdc HDC, x1, y1, x2, y2, x3, y3, x4, y4 int) bool {
	ret, _, _ := chord.Call(
		uintptr(hdc),
		uintptr(x1),
		uintptr(y1),
		uintptr(x2),
		uintptr(y2),
		uintptr(x3),
		uintptr(y3),
		uintptr(x4),
		uintptr(y4),
	)
	return ret != 0
}

func Polygon(hdc HDC, apt []POINT, cpt int) bool {
	if len(apt) == 0 {
		return false
	}

	ret, _, _ := polygon.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(&apt[0])),
		uintptr(cpt),
	)
	return ret != 0
}

func Polyline(hdc HDC, apt []POINT, cpt int) bool {
	if len(apt) == 0 {
		return false
	}

	ret, _, _ := polyline.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(&apt[0])),
		uintptr(cpt),
	)
	return ret != 0
}

func PolyBezier(hdc HDC, apt []POINT, cpt DWORD) bool {
	if len(apt) == 0 {
		return false
	}

	ret, _, _ := polyBezier.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(&apt[0])),
		uintptr(cpt),
	)
	return ret != 0
}

func IntersectClipRect(hdc HDC, left, top, right, bottom int) int {
	ret, _, _ := intersectClipRect.Call(
		uintptr(hdc),
		uintptr(left),
		uintptr(top),
		uintptr(right),
		uintptr(bottom),
	)
	return int(ret)
}

func SelectClipRgn(hdc HDC, region HRGN) int {
	ret, _, _ := selectClipRgn.Call(uintptr(hdc), uintptr(region))
	return int(ret)
}

func CreateRectRgn(x1, y1, x2, y2 int) HRGN {
	ret, _, _ := createRectRgn.Call(
		uintptr(x1),
		uintptr(y1),
		uintptr(x2),
		uintptr(y2),
	)
	return HRGN(ret)
}

func CombineRgn(dest, src1, src2 HRGN, mode int) int {
	ret, _, _ := combineRgn.Call(
		uintptr(dest),
		uintptr(src1),
		uintptr(src2),
		uintptr(mode),
	)
	return int(ret)
}

type FontType int

const (
	RASTER_FONTTYPE   FontType = 1
	DEVICE_FONTTYPE   FontType = 2
	TRUETYPE_FONTTYPE FontType = 4
)

func (t FontType) String() string {
	switch t {
	case RASTER_FONTTYPE:
		return "RASTER_FONTTYPE"
	case DEVICE_FONTTYPE:
		return "DEVICE_FONTTYPE"
	case TRUETYPE_FONTTYPE:
		return "TRUETYPE_FONTTYPE"
	}
	return strconv.Itoa(int(t))
}

func EnumFontFamiliesEx(hdc HDC, font LOGFONT, f func(font *LOGFONTEX, metric *ENUMTEXTMETRIC, fontType FontType) bool) {
	callback := syscall.NewCallback(func(font, metric uintptr, typ uint32, _ uintptr) uintptr {
		if f(
			(*LOGFONTEX)(unsafe.Pointer(font)),
			(*ENUMTEXTMETRIC)(unsafe.Pointer(metric)),
			FontType(typ),
		) {
			return 1
		}
		return 0
	})
	enumFontFamiliesEx.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(&font)),
		callback,
		0,
		0,
	)
}

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

func SetTextAlign(hdc HDC, align UINT) uint {
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

func AbortPath(hdc HDC) bool {
	ret, _, _ := abortPath.Call(
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

func ExtTextOutW(hdc HDC, x, y int, options UINT, lprect *RECT, lpString string, c UINT, lpDx []INT) bool {

	if len(lpDx) == 0 {
		lpDx = make([]INT, 1)
	}

	lpStringUint16 := syscall.StringToUTF16Ptr(lpString)

	ret, _, _ := extTextOutW.Call(
		uintptr(hdc),
		uintptr(x),
		uintptr(y),
		uintptr(options),
		uintptr(unsafe.Pointer(lprect)),
		uintptr(unsafe.Pointer(lpStringUint16)),
		uintptr(c),
		uintptr(unsafe.Pointer(&lpDx[0])),
	)
	return ret != 0
}

func PolyBezierTo(hdc HDC, apt []POINT, cpt DWORD) bool {
	if len(apt) == 0 {
		return false
	}

	ret, _, _ := polyBezierTo.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(&apt[0])),
		uintptr(cpt),
	)
	return ret != 0
}

func PolylineTo(hdc HDC, apt []POINT, cpt DWORD) bool {
	if len(apt) == 0 {
		return false
	}

	ret, _, _ := polylineTo.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(&apt[0])),
		uintptr(cpt),
	)
	return ret != 0
}

func PolyPolygon(hdc HDC, apt []POINT, asz []int, csz int) bool {
	if len(apt) == 0 || len(asz) == 0 {
		return false
	}

	ret, _, _ := polyPolygon.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(&apt[0])),
		uintptr(unsafe.Pointer(&asz[0])),
		uintptr(csz),
	)
	return ret != 0
}

func StretchDIBits(hdc HDC, xDest, yDest, destWidth, destHeight, xSrc, ySrc, srcWidth, srcHeight int, lpBits []byte, lpbmi *BITMAPINFO, iUsage UINT, rop DWORD) int {
	if len(lpBits) == 0 {
		return 0
	}

	ret, _, _ := stretchDIBits.Call(
		uintptr(hdc),
		uintptr(xDest),
		uintptr(yDest),
		uintptr(destWidth),
		uintptr(destHeight),
		uintptr(xSrc),
		uintptr(ySrc),
		uintptr(srcWidth),
		uintptr(srcHeight),
		uintptr(unsafe.Pointer(&lpBits[0])),
		uintptr(unsafe.Pointer(lpbmi)),
		uintptr(iUsage),
		uintptr(rop),
	)
	return int(ret)
}

func SetPixelV(hdc HDC, x, y int, color COLORREF) bool {
	ret, _, _ := setPixelV.Call(
		uintptr(hdc),
		uintptr(x),
		uintptr(y),
		uintptr(color),
	)
	return ret != 0
}

func SetMapperFlags(hdc HDC, flags DWORD) DWORD {
	ret, _, _ := setMapperFlags.Call(
		uintptr(hdc),
		uintptr(flags),
	)
	return DWORD(ret)
}

func SetROP2(hdc HDC, rop2 int) int {
	ret, _, _ := setROP2.Call(
		uintptr(hdc),
		uintptr(rop2),
	)
	return int(ret)
}

func ScaleWindowExtEx(hdc HDC, xn, xd, yn, yd int, lpsz *SIZE) bool {
	ret, _, _ := scaleWindowExtEx.Call(
		uintptr(hdc),
		uintptr(xn),
		uintptr(xd),
		uintptr(yn),
		uintptr(yd),
		uintptr(unsafe.Pointer(lpsz)),
	)
	return ret != 0
}

func SetMetaRgn(hdc HDC) int {
	ret, _, _ := setMetaRgn.Call(
		uintptr(hdc),
	)
	return int(ret)
}

func OffsetClipRgn(hdc HDC, x, y int) int {
	ret, _, _ := offsetClipRgn.Call(
		uintptr(hdc),
		uintptr(x),
		uintptr(y),
	)
	return int(ret)
}

func SetTextJustification(hdc HDC, extra, count int) bool {
	ret, _, _ := setTextJustification.Call(
		uintptr(hdc),
		uintptr(extra),
		uintptr(count),
	)
	return ret != 0
}

func FillRgn(hdc HDC, hrgn HRGN, hbr HBRUSH) bool {
	ret, _, _ := fillRgn.Call(
		uintptr(hdc),
		uintptr(hrgn),
		uintptr(hbr),
	)
	return ret != 0
}

func CreateRectRgnIndirect(lprect *RECT) HRGN {
	ret, _, _ := createRectRgnIndirect.Call(
		uintptr(unsafe.Pointer(lprect)),
	)
	return HRGN(ret)
}

func SetGraphicsMode(hdc HDC, iMode int) int {
	ret, _, _ := setGraphicsMode.Call(
		uintptr(hdc),
		uintptr(iMode),
	)
	return int(ret)
}

func CreatePalette(plpal *LOGPALETTE) HPALETTE {
	ret, _, _ := createPalette.Call(
		uintptr(unsafe.Pointer(plpal)),
	)
	return HPALETTE(ret)
}

func SelectPalette(hdc HDC, hpal HPALETTE, bForceBkgd BOOL) HPALETTE {
	ret, _, _ := fillRgn.Call(
		uintptr(hdc),
		uintptr(hpal),
		uintptr(bForceBkgd),
	)

	return HPALETTE(ret)
}

func SetBitmapBits(bitmap HBITMAP, cb DWORD, pvBits []byte) LONG {
	if len(pvBits) == 0 {
		return LONG(0)
	}

	ret, _, _ := setBitmapBits.Call(
		uintptr(bitmap),
		uintptr(cb),
		uintptr(unsafe.Pointer(&pvBits[0])),
	)

	return LONG(ret)
}

func CreateDIBitmap(hdc HDC, pbmih *BITMAPINFOHEADER, flInit DWORD, pjBits []byte, pbmi *BITMAPINFO, iUsage UINT) HBITMAP {
	if len(pjBits) == 0 {
		return HBITMAP(0)
	}

	ret, _, _ := createDIBitmap.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(pbmih)),
		uintptr(flInit),
		uintptr(unsafe.Pointer(&pjBits[0])),
		uintptr(unsafe.Pointer(pbmi)),
		uintptr(iUsage),
	)

	return HBITMAP(ret)

}

func SetMiterLimit(hdc HDC, limit float32, old *float32) bool {
	if old == nil {
		old = new(float32)
	}

	ret, _, _ := setMiterLimit.Call(
		uintptr(hdc),
		uintptr(limit),
		uintptr(unsafe.Pointer(old)),
	)

	return ret != 0
}

func ExtSelectClipRgn(hdc HDC, hgrn HRGN, mode int) int {
	ret, _, _ := extSelectClipRgn.Call(
		uintptr(hdc),
		uintptr(hgrn),
		uintptr(mode),
	)
	return int(ret)
}
