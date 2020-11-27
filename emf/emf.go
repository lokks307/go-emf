package emf

import (
	"image"

	"github.com/lokks307/go-emf/w32"
	log "github.com/sirupsen/logrus"
)

// FIXME: handle following stockobject
// DEFAULT_PALETTE, DC_BRUSH, DC_PEN

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
		default:
			log.Error("Unknown type of object")
		}
	}
}

type EmfContext struct {
	MDC      w32.HDC
	Width    int
	Height   int
	Objects  map[uint32]interface{}
	BitCount int
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
	memDC := w32.CreateCompatibleDC(0)
	hBitmap := w32.CreateCompatibleBitmap(memDC, width, height)

	if hBitmap == 0 {
		log.Error("failed to create CreateCompatibleBitmap")
	}

	w32.SetTextAlign(memDC, TA_LEFT|TA_TOP)
	w32.SetBkColor(memDC, 0x00FFFFFF)

	w32.SelectObject(memDC, w32.HGDIOBJ(hBitmap))
	w32.Rectangle(memDC, 0, 0, width, height) // too fill white background

	w32.SetGraphicsMode(memDC, w32.GM_ADVANCED)

	defer func() {
		w32.DeleteObject(w32.HGDIOBJ(hBitmap))
	}()

	return &EmfContext{
		MDC:      memDC,
		Objects:  make(map[uint32]interface{}),
		Width:    width,
		Height:   height,
		BitCount: w32.GetDeviceCaps(memDC, w32.COLORRES),
	}

}

func (e *EmfContext) DrawToImage() (*image.RGBA, error) {
	return DeviceContextToImage(e.MDC, e.GetWidth(), e.GetHeight())
}
