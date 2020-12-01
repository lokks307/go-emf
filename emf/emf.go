package emf

import (
	"image"

	im "github.com/disintegration/imaging"
	"github.com/lokks307/go-emf/w32"
	log "github.com/sirupsen/logrus"
)

// FIXME: handle following stockobject
// DEFAULT_PALETTE, DC_BRUSH, DC_PEN

var StockObjects map[uint32]interface{}

func init() {
	StockObjects = map[uint32]interface{}{
		WHITE_BRUSH:         w32.HBRUSH(w32.GetStockObject(WHITE_BRUSH)),
		LTGRAY_BRUSH:        w32.HBRUSH(w32.GetStockObject(LTGRAY_BRUSH)),
		GRAY_BRUSH:          w32.HBRUSH(w32.GetStockObject(GRAY_BRUSH)),
		DKGRAY_BRUSH:        w32.HBRUSH(w32.GetStockObject(DKGRAY_BRUSH)),
		BLACK_BRUSH:         w32.HBRUSH(w32.GetStockObject(BLACK_BRUSH)),
		NULL_BRUSH:          w32.HBRUSH(w32.GetStockObject(NULL_BRUSH)),
		DC_BRUSH:            w32.HBRUSH(w32.GetStockObject(DC_BRUSH)),
		WHITE_PEN:           w32.HPEN(w32.GetStockObject(WHITE_PEN)),
		BLACK_PEN:           w32.HPEN(w32.GetStockObject(BLACK_PEN)),
		DC_PEN:              w32.HPEN(w32.GetStockObject(DC_PEN)),
		NULL_PEN:            w32.HPEN(w32.GetStockObject(NULL_BRUSH)),
		OEM_FIXED_FONT:      w32.HFONT(w32.GetStockObject(OEM_FIXED_FONT)),
		ANSI_FIXED_FONT:     w32.HFONT(w32.GetStockObject(ANSI_FIXED_FONT)),
		ANSI_VAR_FONT:       w32.HFONT(w32.GetStockObject(ANSI_VAR_FONT)),
		SYSTEM_FONT:         w32.HFONT(w32.GetStockObject(SYSTEM_FONT)),
		DEVICE_DEFAULT_FONT: w32.HFONT(w32.GetStockObject(DEVICE_DEFAULT_FONT)),
		SYSTEM_FIXED_FONT:   w32.HFONT(w32.GetStockObject(SYSTEM_FIXED_FONT)),
		DEFAULT_GUI_FONT:    w32.HFONT(w32.GetStockObject(DEFAULT_GUI_FONT)),
		DEFAULT_PALETTE:     w32.HPALETTE(w32.GetStockObject(DEFAULT_PALETTE)),
	}
}

type EmfContext struct {
	MDC          w32.HDC
	Width        int
	Height       int
	Objects      map[uint32]interface{}
	BitCount     int
	GraphicsMode int
	XForm        w32.XFORM
	View         w32.RECT
	Window       w32.SIZE
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

func NewEmfContext(view w32.RECT, window w32.SIZE) *EmfContext {
	memDC := w32.CreateCompatibleDC(0)
	hBitmap := w32.CreateCompatibleBitmap(memDC, int(window.CX), int(window.CY))

	if hBitmap == 0 {
		log.Error("failed to create CreateCompatibleBitmap")
	}

	log.Info("EMF-View = ", view)
	log.Info("EMF-Window = ", window)

	defer func() {
		w32.DeleteObject(w32.HGDIOBJ(hBitmap))
	}()

	emf := &EmfContext{
		MDC:          memDC,
		Objects:      make(map[uint32]interface{}),
		BitCount:     w32.GetDeviceCaps(memDC, w32.COLORRES),
		GraphicsMode: w32.GM_COMPATIBLE,
		View:         view,
		Window:       window,
	}

	emf.SetDefaultXForm()
	emf.ScaleView()

	// w32.SetTextAlign(memDC, TA_LEFT|TA_TOP)
	w32.SetBkColor(emf.MDC, 0x00FFFFFF)

	w32.SetGraphicsMode(emf.MDC, emf.GraphicsMode)

	w32.SelectObject(emf.MDC, w32.HGDIOBJ(hBitmap))
	w32.Rectangle(emf.MDC, 0, 0, int(window.CX), int(window.CY)) // too fill white background

	return emf
}

func (e *EmfContext) SetXForm(xform w32.XFORM) {
	e.XForm = xform
}

func (e *EmfContext) SetDefaultXForm() {
	e.SetXForm(w32.XFORM{
		M11: 1.0,
		M12: 0.0,
		M21: 0.0,
		M22: 1.0,
		Dx:  0.0,
		Dy:  0.0,
	})
}

func (e *EmfContext) ScaleXForm(m11, m22, dx, dy float32) {
	e.XForm.M11 *= m11
	e.XForm.M22 *= m22
	e.XForm.Dx = dx
	e.XForm.Dy = dy
}

func (e *EmfContext) ScaleView() {
	w32.SetWindowExtEx(e.MDC, int(float32(e.Window.CX)*e.XForm.M11), int(float32(e.Window.CY)*e.XForm.M22), nil)
	w32.SetViewportExtEx(e.MDC, int(float32(e.View.Right-e.View.Left)*e.XForm.M11), int(float32(e.View.Bottom-e.View.Top)*e.XForm.M22), nil)
	// w32.SetWindowOrgEx(e.MDC, int(e.XForm.Dx), int(e.XForm.Dy), nil)
	// w32.SetViewportOrgEx(e.MDC, int(-e.XForm.Dx), int(-e.XForm.Dy), nil)
}

func (e *EmfContext) DrawToImage(pMode int) (*image.NRGBA, error) {

	bound := e.View
	device := e.Window

	if bound.Left < 0 {
		bound.Left = 0
	}

	if bound.Top < 0 {
		bound.Top = 0
	}

	if bound.Right > device.CX {
		bound.Right = device.CX
	}

	if bound.Bottom > device.CY {
		bound.Bottom = device.CY
	}

	if img, err := DeviceContextToImage(e.MDC, int(device.CX), int(device.CY)); err != nil {
		return nil, err
	} else {
		if pMode == CROP_AREA {
			return im.Crop(img, image.Rect(int(bound.Left), int(bound.Top), int(bound.Right), int(bound.Bottom))), nil
		} else {
			return im.Crop(img, image.Rect(0, 0, int(device.CX), int(device.CY))), nil
		}
	}
}
