package emf

import (
	"errors"
	"image"
	"unsafe"

	"github.com/lokks307/go-emf/w32"
	log "github.com/sirupsen/logrus"
)

var stockObjectData = map[uint32]interface{}{
	WHITE_BRUSH:         w32.LOGBRUSH{BrushStyle: BS_SOLID, Color: 0x00FFFFFF},
	LTGRAY_BRUSH:        w32.LOGBRUSH{BrushStyle: BS_SOLID, Color: 0x00C0C0C0},
	GRAY_BRUSH:          w32.LOGBRUSH{BrushStyle: BS_SOLID, Color: 0x00808080},
	DKGRAY_BRUSH:        w32.LOGBRUSH{BrushStyle: BS_SOLID, Color: 0x00404040},
	BLACK_BRUSH:         w32.LOGBRUSH{BrushStyle: BS_SOLID, Color: 0x00000000},
	NULL_BRUSH:          w32.LOGBRUSH{BrushStyle: BS_NULL},
	WHITE_PEN:           w32.LOGPEN{PenStyle: PS_COSMETIC | PS_SOLID, ColorRef: 0x00FFFFFF, Width: w32.POINT{X: 1, Y: 0}},
	BLACK_PEN:           w32.LOGPEN{PenStyle: PS_COSMETIC | PS_SOLID, ColorRef: 0x00000000, Width: w32.POINT{X: 1, Y: 0}},
	NULL_PEN:            w32.LOGPEN{PenStyle: PS_NULL},
	OEM_FIXED_FONT:      w32.LOGFONT{CharSet: OEM_CHARSET, PitchAndFamily: (FF_DONTCARE<<4 + FIXED_PITCH)},
	ANSI_FIXED_FONT:     w32.LOGFONT{CharSet: ANSI_CHARSET, PitchAndFamily: (FF_DONTCARE<<4 + FIXED_PITCH)},
	ANSI_VAR_FONT:       w32.LOGFONT{CharSet: ANSI_CHARSET, PitchAndFamily: (FF_DONTCARE<<4 + VARIABLE_PITCH)},
	SYSTEM_FONT:         w32.LOGFONT{Height: 11},
	DEVICE_DEFAULT_FONT: w32.LOGFONT{Height: 11},
	SYSTEM_FIXED_FONT:   w32.LOGFONT{Height: 11},
	DEFAULT_GUI_FONT:    w32.LOGFONT{Height: 11},
}

var StockObjects map[uint32]interface{}

func init() {
	StockObjects = make(map[uint32]interface{})

	for key := range stockObjectData {

		switch object := stockObjectData[key].(type) {
		case w32.LOGPEN:
			StockObjects[key] = w32.CreatePenIndirect(&object)
		case w32.LOGBRUSH:
			StockObjects[key] = w32.CreateBrushIndirect(&object)
		case w32.LOGFONT:
			StockObjects[key] = w32.CreateFontIndirectW(&object)
		}
	}
}

type EmfContext struct {
	MDC     w32.HDC
	Width   int
	Height  int
	Objects map[uint32]interface{}
}

func (e *EmfContext) Release() {
	if !w32.DeleteDC(e.MDC) {
		log.Error("Error on DeleteDC")
	}
}

func (e *EmfContext) GetWidth() int {
	return e.Width
}

func (e *EmfContext) GetHeight() int {
	return e.Height
}

func NewEmfContext(width, height int) *EmfContext {
	MDC := w32.CreateCompatibleDC(0)
	MBM := w32.CreateCompatibleBitmap(MDC, width, height)
	w32.SelectObject(MDC, w32.HGDIOBJ(MBM))

	e := &EmfContext{
		MDC:     MDC,
		Objects: make(map[uint32]interface{}),
		Width:   width,
		Height:  height,
	}

	w32.SetTextAlign(e.MDC, TA_LEFT|TA_TOP)
	w32.Rectangle(e.MDC, 0, 0, width, height) // too fill white background
	return e
}

func (e *EmfContext) DrawToImage() (*image.RGBA, error) {

	width, height := e.GetWidth(), e.GetHeight()
	memory_device := w32.CreateCompatibleDC(e.MDC)

	if memory_device == 0 {
		return nil, errors.New("CreateCompatibleDC failed")
	}
	defer w32.DeleteDC(memory_device)

	bitmap := w32.CreateCompatibleBitmap(e.MDC, width, height)
	if bitmap == 0 {
		return nil, errors.New("CreateCompatibleBitmap failed")
	}
	defer w32.DeleteObject(w32.HGDIOBJ(bitmap))

	var header w32.BITMAPINFOHEADER
	header.BiSize = uint32(unsafe.Sizeof(header))
	header.BiPlanes = 1
	header.BiBitCount = 32
	header.BiWidth = int32(width)
	header.BiHeight = int32(-height)
	header.BiCompression = w32.BI_RGB
	header.BiSizeImage = 0

	bitmapDataSize := uintptr(((int64(width)*int64(header.BiBitCount) + 31) / 32) * 4 * int64(height))
	hmem := w32.GlobalAlloc(w32.GMEM_MOVEABLE, bitmapDataSize)
	defer w32.GlobalFree(hmem)
	memptr := w32.GlobalLock(hmem)
	defer w32.GlobalUnlock(hmem)

	old := w32.SelectObject(memory_device, w32.HGDIOBJ(bitmap))
	if old == 0 {
		return nil, errors.New("SelectObject failed")
	}
	defer w32.SelectObject(memory_device, old)

	if !w32.BitBlt(memory_device, 0, 0, width, height, e.MDC, 0, 0, w32.SRCCOPY) {
		return nil, errors.New("BitBlt failed")
	}

	if w32.GetDIBits(e.MDC, bitmap, 0, w32.UINT(height), memptr, (*w32.BITMAPINFO)(unsafe.Pointer(&header)), w32.DIB_RGB_COLORS) == 0 {
		return nil, errors.New("GetDIBits failed")
	}

	rect := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rect)

	i := 0
	src := uintptr(memptr)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			v0 := *(*uint8)(unsafe.Pointer(src))
			v1 := *(*uint8)(unsafe.Pointer(src + 1))
			v2 := *(*uint8)(unsafe.Pointer(src + 2))

			// BGRA => RGBA, and set A to 255
			img.Pix[i], img.Pix[i+1], img.Pix[i+2], img.Pix[i+3] = v2, v1, v0, 255

			i += 4
			src += 4
		}
	}

	return img, nil
}
