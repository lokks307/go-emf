package emf

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/lokks307/go-emf/w32"
	log "github.com/sirupsen/logrus"
)

type Recorder interface {
	Draw(*EmfContext)
}
type Record struct {
	Type uint32
	Size uint32
}

func (r *Record) Draw(ctx *EmfContext) {
	log.Trace("Unsupported Draw")
}

func readRecord(reader *bytes.Reader) (Recorder, error) {
	var defaultRecord Record

	if err := binary.Read(reader, binary.LittleEndian, &defaultRecord); err != nil {
		return nil, err
	}

	log.Tracef("Record type = %02x\n", defaultRecord.Type)

	fn, ok := records[defaultRecord.Type]
	if !ok {
		return nil, fmt.Errorf("Unknown record %#v found", defaultRecord.Type)
	}

	if fn != nil {
		return fn(reader, defaultRecord.Size)
	}

	// default implementation skips record data
	_, err := reader.Seek(int64(defaultRecord.Size-8), os.SEEK_CUR)
	return &defaultRecord, err
}

type HeaderRecord struct {
	Record
	Original HeaderOriginal
	Ext1     HeaderExtension1
	Ext2     HeaderExtension2
}

type HeaderExtension1 struct {
	CbPixelFormat, OffPixelFormat, BOpenGL uint32
}

type HeaderExtension2 struct {
	MicrometersX, MicrometersY uint32
}

type HeaderOriginal struct {
	Bounds          w32.RECT
	Frame           w32.RECT
	RecordSignature uint32
	Version         uint32
	Bytes           uint32
	Records         uint32
	Handles         uint16
	Reserved        uint16
	NDescription    uint32
	OffDescription  uint32
	NPalEntries     uint32
	Device          w32.SIZE
	Millimeters     w32.SIZE
}

func (HeaderOriginal) Size() uint32 {
	return 88
}

func readHeaderRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	hdr := &HeaderRecord{}
	hdr.Record = Record{Type: EMR_HEADER, Size: size}
	headerSize := hdr.Original.Size()

	if size < headerSize {
		return nil, errors.New("invalid minimum header size")
	}
	headerSize = size

	if err := binary.Read(reader, binary.LittleEndian, &hdr.Original); err != nil {
		return nil, err
	}

	numBytesDescription := hdr.Original.OffDescription + 2*hdr.Original.NDescription
	if hdr.Original.OffDescription >= hdr.Original.Size() && numBytesDescription <= size {
		headerSize = hdr.Original.OffDescription
	}

	if headerSize >= 100 {
		if err := binary.Read(reader, binary.LittleEndian, &hdr.Ext1); err != nil {
			return nil, err
		}

		if hdr.Ext1.OffPixelFormat >= 100 && (hdr.Ext1.OffPixelFormat+hdr.Ext1.CbPixelFormat) <= size {
			if hdr.Ext1.OffPixelFormat < headerSize {
				headerSize = hdr.Ext1.OffPixelFormat
			}
		}
	}

	if headerSize >= 108 {
		if err := binary.Read(reader, binary.LittleEndian, &hdr.Ext2); err != nil {
			return nil, err
		}
	}

	reader.Seek(int64(size), os.SEEK_SET)

	return hdr, nil
}

type SetWindowExtExRecord struct {
	Record
	Extent w32.SIZE
}

func readSetWindowExtExRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetWindowExtExRecord{}
	r.Record = Record{Type: EMR_SETWINDOWEXTEX, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Extent); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetWindowExtExRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETWINDOWEXTEX")

	if !w32.SetWindowExtEx(ctx.MDC, int(r.Extent.CX), int(r.Extent.CY), nil) {
		log.Error("failed to run SetWindowExtEx")
	}
}

type SetWindowOrgExRecord struct {
	Record
	Origin w32.POINT
}

func readSetWindowOrgExRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetWindowOrgExRecord{}
	r.Record = Record{Type: EMR_SETWINDOWORGEX, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Origin); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetWindowOrgExRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETWINDOWORGEX")

	if !w32.SetWindowOrgEx(ctx.MDC, int(r.Origin.X), int(r.Origin.Y), nil) {
		log.Error("failed to run SetWindowOrgEx")
	}
}

type SetWiewporTextExRecord struct {
	Record
	Extent w32.SIZE
}

func readSetViewportExtExRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetWiewporTextExRecord{}
	r.Record = Record{Type: EMR_SETVIEWPORTEXTEX, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Extent); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetWiewporTextExRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETVIEWPORTEXTEX")

	if !w32.SetViewportExtEx(ctx.MDC, int(r.Extent.CX), int(r.Extent.CY), nil) {
		log.Error("failed to run SetViewportExtEx")
	}
}

type SetWiewportOrgExRecord struct {
	Record
	Origin w32.POINT
}

func readSetViewportOrgExRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetWiewportOrgExRecord{}
	r.Record = Record{Type: EMR_SETVIEWPORTORGEX, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Origin); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetWiewportOrgExRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETVIEWPORTORGEX")

	if !w32.SetViewportOrgEx(ctx.MDC, int(r.Origin.X), int(r.Origin.Y), nil) {
		log.Error("failed to run SetViewportOrgEx")
	}
}

type EofRecord struct {
	Record
	NPalEntries   uint32
	OffPalEntries uint32
	SizeLast      uint32
}

func readEofRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &EofRecord{}
	r.Record = Record{Type: EMR_EOF, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.NPalEntries); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.OffPalEntries); err != nil {
		return nil, err
	}

	if r.NPalEntries > 0 {
		log.Error("nPalEntries found - ", r.NPalEntries)
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.SizeLast); err != nil {
		return nil, err
	}

	return r, nil
}

type SetMapModeRecord struct {
	Record
	MapMode uint32
}

func readSetMapModeRecord(reader *bytes.Reader, size uint32) (Recorder, error) {

	r := &SetMapModeRecord{}
	r.Record = Record{Type: EMR_SETMAPMODE, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.MapMode); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetMapModeRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETMAPMODE")

	if w32.SetMapMode(ctx.MDC, int(r.MapMode)) == 0 {
		log.Error("failed to run SetMapMode")
	}
}

type SetBkModeRecord struct {
	Record
	BkMode uint32
}

func readSetBkModeRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetBkModeRecord{}
	r.Record = Record{Type: EMR_SETBKMODE, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.BkMode); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetBkModeRecord) Draw(ctx *EmfContext) {
	log.Tracef("Draw EMR_SETBKMODE 0x%04x", r.BkMode)

	if w32.SetBkMode(ctx.MDC, int(r.BkMode)) == 0 {
		log.Error("failed to run SetBkMode")
	}

}

type SetPolyfillModeRecord struct {
	Record
	PolygonFillMode uint32
}

func readSetPolyfillModeRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetPolyfillModeRecord{}
	r.Record = Record{Type: EMR_SETPOLYFILLMODE, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.PolygonFillMode); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetPolyfillModeRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETPOLYFILLMODE")

	if w32.SetPolyFillMode(ctx.MDC, int(r.PolygonFillMode)) == 0 {
		log.Error("failed to run SetPolyFillMode")
	}
}

type SetTextAlignRecord struct {
	Record
	TextAlignmentMode uint32
}

func readSetTextAlignRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetTextAlignRecord{}
	r.Record = Record{Type: EMR_SETTEXTALIGN, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.TextAlignmentMode); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetTextAlignRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETTEXTALIGN")

	// FIXME: it does not work properly

	// if w32.SetTextAlign(ctx.MDC, w32.UINT(r.TextAlignmentMode)) == w32.GDI_ERROR {
	// 	log.Error("failed to run SetTextAlign")
	// }
}

type SetStretchBltModeRecord struct {
	Record
	StretchMode uint32
}

func readSetStretchBltModeRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetStretchBltModeRecord{}
	r.Record = Record{Type: EMR_SETSTRETCHBLTMODE, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.StretchMode); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetStretchBltModeRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETSTRETCHBLTMODE")

	if ret := w32.SetStretchBltMode(ctx.MDC, int(r.StretchMode)); ret == 0 {
		log.Error("failed to run SetStretchBltMode")
	}
}

type SetTextColorRecord struct {
	Record
	Color WMFCOLORREF
}

func readSetTextColorRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetTextColorRecord{}
	r.Record = Record{Type: EMR_SETTEXTCOLOR, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Color); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetTextColorRecord) Draw(ctx *EmfContext) {
	log.Tracef("Draw EMR_SETTEXTCOLOR 0x%08x", r.Color.ColorRef())

	if w32.SetTextColor(ctx.MDC, r.Color.ColorRef()) == w32.COLORREF(w32.CLR_INVALID) {
		log.Error("failed to run SetTextColor")
	}
}

type SetBkColorRecord struct {
	Record
	Color WMFCOLORREF
}

func readSetBkColorRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetBkColorRecord{}
	r.Record = Record{Type: EMR_SETBKCOLOR, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Color); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetBkColorRecord) Draw(ctx *EmfContext) {
	log.Tracef("Draw EMR_SETBKCOLOR")

	if w32.SetBkColor(ctx.MDC, r.Color.ColorRef()) == w32.COLORREF(w32.CLR_INVALID) {
		log.Error("failed to run SetBkColor")
	}
}

type XYNumDenon struct {
	XNum   uint32
	XDenon uint32
	YNum   uint32
	YDenon uint32
}
type ScaleWindowExtExRecord struct {
	Record
	XYNumDenon
}

func readScaleWindowExtExRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &ScaleWindowExtExRecord{}
	r.Record = Record{Type: EMR_SCALEWINDOWEXTEX, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.XYNumDenon); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *ScaleWindowExtExRecord) Draw(ctx *EmfContext) {
	log.Tracef("Draw EMR_SCALEWINDOWEXTEX")

	if !w32.ScaleWindowExtEx(ctx.MDC, int(r.XNum), int(r.XDenon), int(r.YNum), int(r.YDenon), nil) {
		log.Error("failed to run ScaleWindowExtEx")
	}
}

type SetMetaRgnRecord struct {
	Record
}

func readSetMetaRgnRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetMetaRgnRecord{}
	r.Record = Record{Type: EMR_SETMETARGN, Size: size}

	return r, nil
}

func (r *SetMetaRgnRecord) Draw(ctx *EmfContext) {
	log.Tracef("Draw EMR_SETMETARGN")

	if w32.SetMetaRgn(ctx.MDC) == w32.ERROR {
		log.Error("failed to run SetMetaRgn")
	}
}

type OffSetClipRgnRecord struct {
	Record
	Offset w32.POINT
}

func readOffSetClipRgnRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &OffSetClipRgnRecord{}
	r.Record = Record{Type: EMR_OFFSETCLIPRGN, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Offset); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *OffSetClipRgnRecord) Draw(ctx *EmfContext) {
	log.Tracef("Draw EMR_OFFSETCLIPRGN")

	if w32.OffsetClipRgn(ctx.MDC, int(r.Offset.X), int(r.Offset.Y)) == w32.ERROR {
		log.Error("failed to run OffsetClipRgn")
	}
}

type BreakExCn struct {
	NBreakExtra uint32
	NBreakCount uint32
}
type SetTextJustificationRecord struct {
	Record
	BreakExCn
}

func readSetTextJustificationRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetTextJustificationRecord{}
	r.Record = Record{Type: EMR_SETTEXTJUSTIFICATION, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.BreakExCn); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetTextJustificationRecord) Draw(ctx *EmfContext) {
	log.Tracef("Draw EMR_SETTEXTJUSTIFICATION")

	if !w32.SetTextJustification(ctx.MDC, int(r.NBreakExtra), int(r.NBreakCount)) {
		log.Error("failed to run SetTextJustification")
	}
}

type MoveToExRecord struct {
	Record
	Offset w32.POINT
}

func readMoveToExRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &MoveToExRecord{}
	r.Record = Record{Type: EMR_MOVETOEX, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Offset); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *MoveToExRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_MOVETOEX")

	if !w32.MoveToEx(ctx.MDC, int(r.Offset.X), int(r.Offset.Y), nil) {
		log.Error("failed to run MoveToEx")
	}
}

type FillRgnRecord struct {
	Record
	Bounds      w32.RECT
	RgnDataSize uint32
	IhBrush     uint32
	RgnData     RegionData
}

func readFillRgnRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &FillRgnRecord{}
	r.Record = Record{Type: EMR_FILLRGN, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Bounds); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.RgnDataSize); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.IhBrush); err != nil {
		return nil, err
	}

	r.RgnData.Data = make([]w32.RECT, r.RgnData.CountRects)
	if err := binary.Read(reader, binary.LittleEndian, &r.RgnData); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *FillRgnRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_FILLRGN")

	gdiObject, ok := StockObjects[r.IhBrush]
	if !ok {
		gdiObject, ok = ctx.Objects[r.IhBrush]
		if !ok {
			log.Errorf("Object 0x%x not found\n", r.IhBrush)
			return
		}
	}

	hbrush := gdiObject.(w32.HBRUSH)

	for idx := range r.RgnData.Data {
		hrgn := w32.CreateRectRgnIndirect(&r.RgnData.Data[idx])
		if !w32.FillRgn(ctx.MDC, hrgn, hbrush) {
			log.Error("faile to run FillRgn")
		}
	}
}

type IntersectClipRectRecord struct {
	Record
	Clip w32.RECT
}

func readIntersectClipRectRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &IntersectClipRectRecord{}
	r.Record = Record{Type: EMR_INTERSECTCLIPRECT, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Clip); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *IntersectClipRectRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_INTERSECTCLIPRECT")

	if w32.IntersectClipRect(ctx.MDC, int(r.Clip.Left), int(r.Clip.Top), int(r.Clip.Right), int(r.Clip.Bottom)) == w32.ERROR {
		log.Error("failed to run IntersectClipRect")
	}
}

type SaveDCRecord struct {
	Record
}

func readSaveDCRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	return &SaveDCRecord{Record: Record{Type: EMR_SAVEDC, Size: size}}, nil
}

func (r *SaveDCRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SAVEDC")

	if w32.SaveDC(ctx.MDC) == 0 {
		log.Error("failed to run SaveDC")
	}
}

type RestoreDCRecord struct {
	Record
	SavedDC int32
}

func readRestoreDCRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &RestoreDCRecord{}
	r.Record = Record{Type: EMR_RESTOREDC, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.SavedDC); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *RestoreDCRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_RESTOREDC")

	if !w32.RestoreDC(ctx.MDC, int(r.SavedDC)) {
		log.Error("failed to run RestoreDC")
	}
}

type SetWorldTransformRecord struct {
	Record
	XForm w32.XFORM
}

func readSetWorldTransformRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetWorldTransformRecord{}
	r.Record = Record{Type: EMR_SETWORLDTRANSFORM, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.XForm); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetWorldTransformRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETWORLDTRANSFORM")

	// if !w32.SetWorldTransform(ctx.MDC, &r.XForm) {
	// 	log.Error("failed to run SetWorldTransform")
	// }
}

type ModifyWorldTransformRecord struct {
	Record
	XForm                    w32.XFORM
	ModifyWorldTransformMode uint32
}

func readModifyWorldTransformRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &ModifyWorldTransformRecord{}
	r.Record = Record{Type: EMR_MODIFYWORLDTRANSFORM, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.XForm); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.ModifyWorldTransformMode); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *ModifyWorldTransformRecord) Draw(ctx *EmfContext) {
	log.Tracef("Draw EMR_MODIFYWORLDTRANSFORM 0x%02x", r.ModifyWorldTransformMode)

	// if !w32.ModifyWorldTransform(ctx.MDC, &r.XForm, w32.DWORD(r.ModifyWorldTransformMode)) {
	// 	log.Error("failed to run ModifyWorldTransform")
	// }
}

type SelectObjectRecord struct {
	Record
	IhObject uint32
}

func readSelectObjectRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SelectObjectRecord{}
	r.Record = Record{Type: EMR_SELECTOBJECT, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.IhObject); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SelectObjectRecord) Draw(ctx *EmfContext) {

	log.Tracef("Draw EMR_SELECTOBJECT 0x%08x", r.IhObject)

	gdiObject, ok := StockObjects[r.IhObject]
	if !ok {
		gdiObject, ok = ctx.Objects[r.IhObject]
		if !ok {
			log.Errorf("Object 0x%x not found\n", r.IhObject)
			return
		}
	}

	switch object := gdiObject.(type) {
	case w32.HPEN:
		w32.SelectObject(ctx.MDC, w32.HGDIOBJ(object))
	case w32.HBRUSH:
		w32.SelectObject(ctx.MDC, w32.HGDIOBJ(object))
	case w32.HFONT:
		w32.SelectObject(ctx.MDC, w32.HGDIOBJ(object))
	default:
		log.Error("Unknown type of object")
	}
}

type CreatePenRecord struct {
	Record
	IhPen  uint32
	LogPen WMFLOGPEN
}

func readCreatePenRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &CreatePenRecord{}
	r.Record = Record{Type: EMR_CREATEPEN, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.IhPen); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.LogPen); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *CreatePenRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_CREATEPEN")

	w32logpen := r.LogPen.LogPen()

	ctx.Objects[r.IhPen] = w32.CreatePenIndirect(&w32logpen)
}

type CreateBrushIndirectRecord struct {
	Record
	IhBrush  uint32
	LogBrush WMFLOGBRUSH
}

func readCreateBrushIndirectRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &CreateBrushIndirectRecord{}
	r.Record = Record{Type: EMR_CREATEBRUSHINDIRECT, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.IhBrush); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.LogBrush); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *CreateBrushIndirectRecord) Draw(ctx *EmfContext) {
	log.Tracef("Draw EMR_CREATEBRUSHINDIRECT 0x%08x", r.IhBrush)

	w32logbrush := r.LogBrush.LogBrush()

	ctx.Objects[r.IhBrush] = w32.CreateBrushIndirect(&w32logbrush)
}

type CreatePaletteRecord struct {
	Record
	IhPal      uint32
	LogPalette w32.LOGPALETTE
}

func readCreatePaletteRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &CreatePaletteRecord{}
	r.Record = Record{Type: EMR_CREATEBRUSHINDIRECT, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.IhPal); err != nil {
		return nil, err
	}

	var err error
	r.LogPalette, err = readLogPalette(reader)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *CreatePaletteRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_CREATEPALETTE")

	ctx.Objects[r.IhPal] = w32.CreatePalette(&r.LogPalette)
}

type SelectPaletteRecord struct {
	Record
	IhPal uint32
}

func readSelectPaletteRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SelectPaletteRecord{}
	r.Record = Record{Type: EMR_SELECTPALETTE, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.IhPal); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SelectPaletteRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SELECTPALETTE")

	gdiObject, ok := StockObjects[r.IhPal]
	if !ok {
		gdiObject, ok = ctx.Objects[r.IhPal]
		if !ok {
			log.Errorf("Object 0x%x not found\n", r.IhPal)
			return
		}
	}

	switch object := gdiObject.(type) {
	case w32.HPALETTE:
		w32.SelectPalette(ctx.MDC, object, w32.FALSE)
	default:
		log.Error("Unknown type of object")
	}

}

type DeleteObjectRecord struct {
	Record
	IhObject uint32
}

func readDeleteObjectRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &DeleteObjectRecord{}
	r.Record = Record{Type: EMR_DELETEOBJECT, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.IhObject); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *DeleteObjectRecord) Draw(ctx *EmfContext) {
	log.Tracef("Draw EMR_DELETEOBJECT 0x%08x", r.IhObject)

	delete(ctx.Objects, r.IhObject)
}

type RectangleRecord struct {
	Record
	Box w32.RECT
}

func readRectangleRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &RectangleRecord{}
	r.Record = Record{Type: EMR_RECTANGLE, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Box); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *RectangleRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_RECTANGLE")

	if !w32.Rectangle(ctx.MDC, int(r.Box.Left), int(r.Box.Top), int(r.Box.Right), int(r.Box.Bottom)) {
		log.Error("failed to run Rectangle")
	}
}

type ArcRecord struct {
	Record
	Box   w32.RECT
	Start w32.POINT
	End   w32.POINT
}

func readArcRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &ArcRecord{}
	r.Record = Record{Type: EMR_ARC, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Box); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.Start); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.End); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *ArcRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_ARC")

	if !w32.Arc(ctx.MDC, int(r.Box.Left), int(r.Box.Top), int(r.Box.Right), int(r.Box.Bottom),
		int(r.Start.X), int(r.Start.Y), int(r.End.X), int(r.End.Y)) {
		log.Error("failed to run Arc")
	}
}

type LineToRecord struct {
	Record
	Point w32.POINT
}

func readLineToRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &LineToRecord{}
	r.Record = Record{Type: EMR_LINETO, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Point); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *LineToRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_LINETO")

	if !w32.LineTo(ctx.MDC, int(r.Point.X), int(r.Point.Y)) {
		log.Error("failed to run LineTo")
	}
}

type BeginPathRecord struct {
	Record
}

func readBeginPathRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	return &BeginPathRecord{Record{Type: EMR_BEGINPATH, Size: size}}, nil
}

func (r *BeginPathRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_BEGINPATH")

	if !w32.BeginPath(ctx.MDC) {
		log.Error("failed to run BeginPath")
	}
}

type EndPathRecord struct {
	Record
}

func readEndPathRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	return &EndPathRecord{Record{Type: EMR_ENDPATH, Size: size}}, nil
}

func (r *EndPathRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_ENDPATH")

	if !w32.EndPath(ctx.MDC) {
		log.Error("failed to run EndPath")
	}
}

type CloseFigureRecord struct {
	Record
}

func readCloseFigureRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	return &CloseFigureRecord{Record{Type: EMR_CLOSEFIGURE, Size: size}}, nil
}

func (r *CloseFigureRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_CLOSEFIGURE")

	if !w32.CloseFigure(ctx.MDC) {
		log.Error("failed to run CloseFigure")
	}
}

type FillPathRecord struct {
	Record
	Bounds w32.RECT
}

func readFillPathRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &FillPathRecord{}
	r.Record = Record{Type: EMR_FILLPATH, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Bounds); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *FillPathRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_FILLPATH")

	hrgn := w32.CreateRectRgn(int(r.Bounds.Left), int(r.Bounds.Top), int(r.Bounds.Right), int(r.Bounds.Bottom))
	w32.SelectObject(ctx.MDC, w32.HGDIOBJ(hrgn))

	if !w32.FillPath(ctx.MDC) {
		log.Error("failed to run FillPath")
	}
}

type StrokeAndFillPathRecord struct {
	Record
	Bounds w32.RECT
}

func readStrokeAndFillPathRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &StrokeAndFillPathRecord{}
	r.Record = Record{Type: EMR_STROKEANDFILLPATH, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Bounds); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *StrokeAndFillPathRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_STROKEANDFILLPATH")

	hrgn := w32.CreateRectRgn(int(r.Bounds.Left), int(r.Bounds.Top), int(r.Bounds.Right), int(r.Bounds.Bottom))
	w32.SelectObject(ctx.MDC, w32.HGDIOBJ(hrgn))

	if !w32.StrokeAndFillPath(ctx.MDC) {
		log.Error("failed to run StrokeAndFillPath")
	}
}

type StrokePathRecord struct {
	Record
	Bounds w32.RECT
}

func readStrokePathRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &StrokePathRecord{}
	r.Record = Record{Type: EMR_STROKEPATH, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Bounds); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *StrokePathRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_STROKEPATH")

	hrgn := w32.CreateRectRgn(int(r.Bounds.Left), int(r.Bounds.Top), int(r.Bounds.Right), int(r.Bounds.Bottom))
	w32.SelectObject(ctx.MDC, w32.HGDIOBJ(hrgn))

	if !w32.StrokePath(ctx.MDC) {
		log.Error("failed to run StrokePath")
	}
}

type SelectClipPathRecord struct {
	Record
	RegionMode uint32
}

func readSelectClipPathRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SelectClipPathRecord{}
	r.Record = Record{Type: EMR_SELECTCLIPPATH, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.RegionMode); err != nil {
		return nil, err
	}
	return r, nil
}

type CommentRecord struct {
	Record
}

func readCommentRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &CommentRecord{}
	r.Record = Record{Type: EMR_COMMENT, Size: size}
	// skip record data
	reader.Seek(int64(size-8), os.SEEK_CUR)
	return r, nil
}

type ExtCreateFontIndirectWRecord struct {
	Record
	IhFonts uint32
	Elw     w32.LOGFONTEXDV
	isExDV  bool
}

func readExtCreateFontIndirectWRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &ExtCreateFontIndirectWRecord{}
	r.Record = Record{Type: EMR_EXTCREATEFONTINDIRECTW, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.IhFonts); err != nil {
		return nil, err
	}

	sizeElw := size - 12

	if sizeElw > w32.LOGFONTPANOSESIZE { // size of Elw = size - 12
		r.isExDV = true
	}

	var remainSize uint32

	if !r.isExDV {
		if err := binary.Read(reader, binary.LittleEndian, &r.Elw.LOGFONT); err != nil {
			return nil, err
		}

		remainSize = sizeElw - w32.LOGFONTSIZE

	} else {
		if err := binary.Read(reader, binary.LittleEndian, &r.Elw.LOGFONTEX); err != nil {
			return nil, err
		}

		if err := binary.Read(reader, binary.LittleEndian, &r.Elw.Signature); err != nil {
			return nil, err
		}

		if err := binary.Read(reader, binary.LittleEndian, &r.Elw.NumAxes); err != nil {
			return nil, err
		}

		if r.Elw.NumAxes > 0 {
			r.Elw.Values = make([]int32, r.Elw.NumAxes)
			if err := binary.Read(reader, binary.LittleEndian, &r.Elw.Values); err != nil {
				return nil, err
			}
		}

		remainSize = sizeElw - w32.LOGFONTEXSIZE - 8 - r.Elw.NumAxes*4
	}

	if remainSize > 0 {
		reader.Seek(int64(remainSize), os.SEEK_CUR)
	}

	return r, nil
}

func (r *ExtCreateFontIndirectWRecord) Draw(ctx *EmfContext) {
	log.Tracef("Draw EMR_EXTCREATEFONTINDIRECTW 0x%08x %s", r.IhFonts, r.Elw.GetFaceName())

	// if r.isExDV {
	// 	ctx.Objects[r.IhFonts] = w32.CreateFontIndirectExW(&r.Elw)
	// } else {
	// 	ctx.Objects[r.IhFonts] = w32.CreateFontIndirectW(&r.Elw.LOGFONT)
	// }

	ctx.Objects[r.IhFonts] = w32.CreateFontIndirectW(&r.Elw.LOGFONT)
}

type ExtTextOutWRecord struct {
	Record
	Bounds        w32.RECT
	IGraphicsMode uint32
	ExScale       float32
	EyScale       float32
	WEmrText      EmrText
}

func readExtTextOutWRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &ExtTextOutWRecord{}
	r.Record = Record{Type: EMR_EXTTEXTOUTW, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Bounds); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.IGraphicsMode); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.ExScale); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.EyScale); err != nil {
		return nil, err
	}

	var err error
	r.WEmrText, err = readEmrText(reader, reader.Len()+36)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *ExtTextOutWRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_EXTTEXTOUTW ", r.WEmrText.GetString())

	if strings.TrimSpace(r.WEmrText.GetString()) == "" {
		return
	}

	//w32.SetGraphicsMode(ctx.MDC, int(r.IGraphicsMode))

	hrgn := w32.CreateRectRgn(int(r.Bounds.Left), int(r.Bounds.Top), int(r.Bounds.Right), int(r.Bounds.Bottom))
	w32.SelectObject(ctx.MDC, w32.HGDIOBJ(hrgn))

	dx := make([]w32.INT, len(r.WEmrText.OutputDx))
	for idx := range r.WEmrText.OutputDx {
		dx[idx] = w32.INT(r.WEmrText.OutputDx[idx])
	}

	if !w32.ExtTextOutW(ctx.MDC, int(r.Bounds.Left), int(r.Bounds.Top),
		w32.UINT(r.WEmrText.Options), &r.WEmrText.Rectangle, r.WEmrText.GetString(), w32.UINT(r.WEmrText.Chars), dx) {
		log.Error("failed to run ExtTextOutW")
	}
}

type PolyBezier16Record struct {
	Record
	Bounds  w32.RECT
	Count   uint32
	APoints []PointS
}

func readPolyBezier16Record(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &PolyBezier16Record{}
	r.Record = Record{Type: EMR_POLYBEZIER16, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Bounds); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.Count); err != nil {
		return nil, err
	}

	r.APoints = make([]PointS, r.Count)
	if err := binary.Read(reader, binary.LittleEndian, &r.APoints); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *PolyBezier16Record) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_POLYBEZIER16")

	hrgn := w32.CreateRectRgn(int(r.Bounds.Left), int(r.Bounds.Top), int(r.Bounds.Right), int(r.Bounds.Bottom))
	w32.SelectObject(ctx.MDC, w32.HGDIOBJ(hrgn))

	bezerPoints := make([]w32.POINT, r.Count)
	for idx := range r.APoints {
		bezerPoints[idx] = w32.POINT{
			X: int32(r.APoints[idx].X),
			Y: int32(r.APoints[idx].Y),
		}
	}

	if !w32.PolyBezier(ctx.MDC, bezerPoints, w32.DWORD(r.Count)) {
		log.Error("failed to run PolyBezier")
	}
}

type Polygon16Record struct {
	Record
	Bounds  w32.RECT
	Count   uint32
	APoints []PointS
}

func readPolygon16Record(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &Polygon16Record{}
	r.Record = Record{Type: EMR_POLYGON16, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Bounds); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.Count); err != nil {
		return nil, err
	}

	r.APoints = make([]PointS, r.Count)
	if err := binary.Read(reader, binary.LittleEndian, &r.APoints); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Polygon16Record) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_POLYGON16")

	vertexPoints := make([]w32.POINT, r.Count)
	for idx := range r.APoints {
		vertexPoints[idx] = w32.POINT{
			X: int32(r.APoints[idx].X),
			Y: int32(r.APoints[idx].Y),
		}
	}

	if !w32.Polygon(ctx.MDC, vertexPoints, int(r.Count)) {
		log.Error("failed to run Polygon")
	}
}

type PolyLine16Record struct {
	Record
	Bounds  w32.RECT
	Count   uint32
	APoints []PointS
}

func readPolyLine16Record(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &PolyLine16Record{}
	r.Record = Record{Type: EMR_POLYLINE16, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Bounds); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.Count); err != nil {
		return nil, err
	}

	r.APoints = make([]PointS, r.Count)
	if err := binary.Read(reader, binary.LittleEndian, &r.APoints); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *PolyLine16Record) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_POLYLINE16")

	points := make([]w32.POINT, r.Count)
	for idx := range r.APoints {
		points[idx] = w32.POINT{
			X: int32(r.APoints[idx].X),
			Y: int32(r.APoints[idx].Y),
		}
	}

	if !w32.Polyline(ctx.MDC, points, int(r.Count)) {
		log.Error("failed to run Polygon")
	}
}

type PolyBezierTo16Record struct {
	Record
	Bounds  w32.RECT
	Count   uint32
	APoints []PointS
}

func readPolyBezierTo16Record(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &PolyBezierTo16Record{}
	r.Record = Record{Type: EMR_POLYBEZIERTO16, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Bounds); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.Count); err != nil {
		return nil, err
	}

	r.APoints = make([]PointS, r.Count)
	if err := binary.Read(reader, binary.LittleEndian, &r.APoints); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *PolyBezierTo16Record) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_POLYBEZIERTO16")

	bezerPoints := make([]w32.POINT, r.Count)
	for idx := range r.APoints {
		bezerPoints[idx] = w32.POINT{
			X: int32(r.APoints[idx].X),
			Y: int32(r.APoints[idx].Y),
		}
	}

	if !w32.PolyBezierTo(ctx.MDC, bezerPoints, w32.DWORD(r.Count)) {
		log.Error("failed to run PolyBezier")
	}
}

type PolyLineTo16Record struct {
	Record
	Bounds  w32.RECT
	Count   uint32
	APoints []PointS
}

func readPolyLineTo16Record(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &PolyLineTo16Record{}
	r.Record = Record{Type: EMR_POLYLINETO16, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Bounds); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.Count); err != nil {
		return nil, err
	}

	r.APoints = make([]PointS, r.Count)
	if err := binary.Read(reader, binary.LittleEndian, &r.APoints); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *PolyLineTo16Record) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_POLYLINETO16")

	points := make([]w32.POINT, r.Count)
	for idx := range r.APoints {
		points[idx] = w32.POINT{
			X: int32(r.APoints[idx].X),
			Y: int32(r.APoints[idx].Y),
		}
	}

	if !w32.PolylineTo(ctx.MDC, points, w32.DWORD(r.Count)) {
		log.Error("failed to run Polygon")
	}
}

type PolyPolygon16Record struct {
	Record
	Bounds            w32.RECT
	NumberOfPolygons  uint32
	Count             uint32
	PolygonPointCount []uint32
	APoints           []PointS
}

func readPolyPolygon16Record(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &PolyPolygon16Record{}
	r.Record = Record{Type: EMR_POLYPOLYGON16, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Bounds); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.NumberOfPolygons); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.Count); err != nil {
		return nil, err
	}

	r.PolygonPointCount = make([]uint32, r.NumberOfPolygons)
	if err := binary.Read(reader, binary.LittleEndian, &r.PolygonPointCount); err != nil {
		return nil, err
	}

	r.APoints = make([]PointS, r.Count)
	if err := binary.Read(reader, binary.LittleEndian, &r.APoints); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *PolyPolygon16Record) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_POLYPOLYGON16")

	points := make([]w32.POINT, r.Count)
	asz := make([]int, r.Count)
	for idx := range r.APoints {
		points[idx] = w32.POINT{
			X: int32(r.APoints[idx].X),
			Y: int32(r.APoints[idx].Y),
		}
		asz[idx] = int(r.PolygonPointCount[idx])
	}

	if !w32.PolyPolygon(ctx.MDC, points, asz, int(r.Count)) {
		log.Error("failed to run Polygon")
	}

}

type ExtCreatePenRecord struct {
	Record
	IhPen   uint32
	OffBmi  uint32
	CbBmi   uint32
	OffBits uint32
	CbBits  uint32
	Elp     w32.LOGPENEX
	BmiSrc  w32.BITMAPINFOHEADER
	BitsSrc []byte
}

func readExtCreatePenRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &ExtCreatePenRecord{}
	r.Record = Record{Type: EMR_EXTCREATEPEN, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.IhPen); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.OffBmi); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.CbBmi); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.OffBits); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.CbBits); err != nil {
		return nil, err
	}

	var err error
	r.Elp, err = readLogPenEx(reader)
	if err != nil {
		return nil, err
	}

	// offset for bitmap info less than possible minimum
	// assuming there is no bitmap
	if r.OffBmi < 52 {
		return r, nil
	}

	// BitmapBuffer

	reader.Seek(int64(r.OffBmi-52-(r.Elp.NumStyleEntries*4)), os.SEEK_CUR) // skipping UndefinedSpace

	if r.CbBmi > 0 {

		if err := binary.Read(reader, binary.LittleEndian, &r.BmiSrc); err != nil {
			return nil, err
		}

		r.BitsSrc = make([]byte, r.CbBits)
		if _, err := reader.Read(r.BitsSrc); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (r *ExtCreatePenRecord) Draw(ctx *EmfContext) {
	log.Tracef("Draw EMR_EXTCREATEPEN 0x%08x", r.IhPen)

	logbrush := w32.LOGBRUSH{
		BrushStyle: r.Elp.BrushStyle,
		Color:      r.Elp.ColorRef,
		BrushHatch: r.Elp.BrushHatch,
	}

	styleEntry := make([]w32.DWORD, len(r.Elp.StyleEntry))
	for idx := range r.Elp.StyleEntry {
		styleEntry[idx] = w32.DWORD(r.Elp.StyleEntry[idx])
	}

	ctx.Objects[r.IhPen] = w32.ExtCreatePen(w32.DWORD(r.Elp.PenStyle), w32.DWORD(r.Elp.Width), &logbrush, w32.DWORD(r.Elp.NumStyleEntries), styleEntry)
}

type SetICMMmodeRecord struct {
	Record
	ICMMode uint32
}

func readSetICMModeRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetICMMmodeRecord{}
	r.Record = Record{Type: EMR_SETICMMODE, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.ICMMode); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetICMMmodeRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETICMMODE")
}

type SetBrushOrgExRecord struct {
	Record
	Origin w32.POINT
}

func readSetBrushOrgExRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetBrushOrgExRecord{}
	r.Record = Record{Type: EMR_SETBRUSHORGEX, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Origin); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetBrushOrgExRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETBRUSHORGEX")

	if !w32.SetBrushOrgEx(ctx.MDC, int(r.Origin.X), int(r.Origin.Y), nil) {
		log.Error("failed to run SetBrushOrgEx")
	}
}

type SetPixelvRecord struct {
	Record
	Pixel w32.POINT
	Color WMFCOLORREF
}

func readSetPixelvRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetPixelvRecord{}
	r.Record = Record{Type: EMR_SETPIXELV, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Pixel); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.Color); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetPixelvRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETPIXELV")

	if !w32.SetPixelV(ctx.MDC, int(r.Pixel.X), int(r.Pixel.Y), r.Color.ColorRef()) {
		log.Error("failed to run SetPixelV")
	}
}

type SetMapperFlagsRecord struct {
	Record
	Flags uint32
}

func readSetMapperFlagsRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetMapperFlagsRecord{}
	r.Record = Record{Type: EMR_SETMAPPERFLAGS, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Flags); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetMapperFlagsRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETMAPPERFLAGS")

	if w32.SetMapperFlags(ctx.MDC, w32.DWORD(r.Flags)) == w32.GDI_ERROR {
		log.Error("failed to run SetMapperFlags")
	}
}

type SetROP2Record struct {
	Record
	ROP2Mode uint32
}

func readSetROP2Record(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetROP2Record{}
	r.Record = Record{Type: EMR_SETROP2, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.ROP2Mode); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetROP2Record) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETROP2")

	if w32.SetROP2(ctx.MDC, int(r.ROP2Mode)) == 0 {
		log.Error("failed to run SetROP2")
	}
}

type SetMiterLimitRecord struct {
	Record
	MiterLimit uint32
}

func readSetMiterLimitRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetMiterLimitRecord{}
	r.Record = Record{Type: EMR_SETMITERLIMIT, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.MiterLimit); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetMiterLimitRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETMITERLIMIT ", r.MiterLimit)

	var old float32

	if !w32.SetMiterLimit(ctx.MDC, float32(r.MiterLimit), &old) {
		log.Error("failed to run SetMiterLimit")
	}
}

type ExtSelectClipRgnRecord struct {
	Record
	RgnDataSize uint32
	RegionMode  uint32
	RgnData     RegionData
}

func readRegionData(reader *bytes.Reader) (RegionData, error) {
	r := RegionData{}
	if err := binary.Read(reader, binary.LittleEndian, &r.RegionDataHeader); err != nil {
		return r, err
	}

	r.Data = make([]w32.RECT, r.RegionDataHeader.CountRects)
	if err := binary.Read(reader, binary.LittleEndian, &r.Data); err != nil {
		return r, err
	}

	return r, nil
}

func readExtSelectClipRgnRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &ExtSelectClipRgnRecord{}
	r.Record = Record{Type: EMR_EXTSELECTCLIPRGN, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.RgnDataSize); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.RegionMode); err != nil {
		return nil, err
	}

	if r.RegionMode != RGN_COPY {

		var err error
		r.RgnData, err = readRegionData(reader)
		if err != nil {
			return nil, err
		}

	} else {
		reader.Seek(int64(size-16), os.SEEK_CUR) // skip
	}

	return r, nil
}

func (r *ExtSelectClipRgnRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_EXTSELECTCLIPRGN")

	if r.RegionMode == RGN_COPY {
		if w32.ExtSelectClipRgn(ctx.MDC, 0, RGN_COPY) == 0 { // default cliping region = null region
			log.Error("failed to run ExtSelectClipRgn")
		}
	} else {
		for _, rect := range r.RgnData.Data {
			hrgn := w32.CreateRectRgn(int(rect.Left), int(rect.Top), int(rect.Right), int(rect.Bottom))

			if w32.ExtSelectClipRgn(ctx.MDC, hrgn, int(r.RegionMode)) == 0 {
				log.Error("failed to run ExtSelectClipRgn")
			}
		}
	}
}
