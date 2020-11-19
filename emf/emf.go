package emf

import (
	"github.com/lokks307/go-emf/w32"
	log "github.com/sirupsen/logrus"
)

var stockObjectData = map[uint32]interface{}{
	WHITE_BRUSH:         w32.LOGBRUSH{LbStyle: BS_SOLID, LbColor: 0x00FFFFFF},
	LTGRAY_BRUSH:        w32.LOGBRUSH{LbStyle: BS_SOLID, LbColor: 0x00C0C0C0},
	GRAY_BRUSH:          w32.LOGBRUSH{LbStyle: BS_SOLID, LbColor: 0x00808080},
	DKGRAY_BRUSH:        w32.LOGBRUSH{LbStyle: BS_SOLID, LbColor: 0x00404040},
	BLACK_BRUSH:         w32.LOGBRUSH{LbStyle: BS_SOLID, LbColor: 0x00000000},
	NULL_BRUSH:          w32.LOGBRUSH{LbStyle: BS_NULL},
	WHITE_PEN:           w32.LOGPEN{LopnStyle: PS_COSMETIC | PS_SOLID, LopnColor: 0x00FFFFFF, LopnWidth: w32.POINT{1, 0}},
	BLACK_PEN:           w32.LOGPEN{LopnStyle: PS_COSMETIC | PS_SOLID, LopnColor: 0x00000000, LopnWidth: w32.POINT{1, 0}},
	NULL_PEN:            w32.LOGPEN{LopnStyle: PS_NULL},
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
	for key := range stockObjectData {

		switch stockObjectData[key].(type) {
		case w32.LOGPEN:
			pen := stockObjectData[key].(w32.LOGPEN)
			StockObjects[key] = w32.CreatePenIndirect(&pen)
		case w32.LOGBRUSH:
			brushex := stockObjectData[key].(w32.LOGBRUSH)
			StockObjects[key] = w32.CreateBrushIndirect(&brushex)
		case w32.LOGFONT:
			w32logfont := stockObjectData[key].(w32.LOGFONT)
			StockObjects[key] = w32.CreateFontIndirectW(&w32logfont)
		}
	}
}

type EmfContext struct {
	MemDC   w32.HDC
	Objects map[uint32]interface{}
	wo      *w32.POINT
	vo      *w32.POINT
	we      *SizeL
	ve      *SizeL
	mm      uint32
}

func (f *EmfFile) NewEmfContext(width, height int) *EmfContext {

	memDC := w32.CreateCompatibleDC(0)
	memBM := w32.CreateCompatibleBitmap(memDC, width, height)
	w32.SelectObject(memDC, w32.HGDIOBJ(memBM))

	emfDC := &EmfContext{
		MemDC:   memDC,
		mm:      MM_TEXT,
		Objects: make(map[uint32]interface{}),
	}

	// init align

	w32.SetTextAlign(emfDC.MemDC, TA_LEFT|TA_TOP)

	return emfDC
}

func (f *EmfFile) DrawToPDF(outPath string) {

	bounds := f.Header.Original.Bounds

	width := int(bounds.Width()) + 1
	height := int(bounds.Height()) + 1

	ctx := f.NewEmfContext(width, height)

	// if bounds.Left != 0 || bounds.Top != 0 {
	// 	ctx.Translate(-float64(bounds.Left), -float64(bounds.Top))
	// }

	for idx := range f.Records {
		log.Tracef("%d-th record", idx)
		f.Records[idx].Draw(ctx)
	}

}
