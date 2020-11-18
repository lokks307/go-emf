package emf

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	"github.com/lokks307/go-emf/w32"
	log "github.com/sirupsen/logrus"
	"github.com/tdewolff/canvas"
)

type Recorder interface {
	Draw(*EmfContext)
}
type Record struct {
	Type, Size uint32
}

func (r *Record) Draw(ctx *EmfContext) {
	/* do nothing */
}

func readRecord(reader *bytes.Reader) (Recorder, error) {
	var rec Record

	if err := binary.Read(reader, binary.LittleEndian, &rec); err != nil {
		return nil, err
	}

	log.Tracef("Record type = %02x\n", rec.Type)

	fn, ok := records[rec.Type]
	if !ok {
		return nil, fmt.Errorf("Unknown record %#v found", rec.Type)
	}

	if fn != nil {
		return fn(reader, rec.Size)
	}

	// default implementation skips record data
	_, err := reader.Seek(int64(rec.Size-8), os.SEEK_CUR)
	return &rec, err
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
	Bounds, Frame           RectL
	RecordSignature         uint32
	Version, Bytes, Records uint32
	Handles                 uint16
	Reserved                uint16
	NDescription            uint32
	OffDescription          uint32
	NPalEntries             uint32
	Device, Millimeters     SizeL
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

type SetwindowextexRecord struct {
	Record
	Extent SizeL
}

func readSetwindowextexRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetwindowextexRecord{}
	r.Record = Record{Type: EMR_SETWINDOWEXTEX, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Extent); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetwindowextexRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETWINDOWEXTEX")

	if !w32.SetWindowExtEx(ctx.MemDC, int(r.Extent.Cx), int(r.Extent.Cy), nil) {
		log.Error("failed to run SetWindowExtEx")
	}
}

type SetwindoworgexRecord struct {
	Record
	Origin PointL
}

func readSetwindoworgexRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetwindoworgexRecord{}
	r.Record = Record{Type: EMR_SETWINDOWORGEX, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Origin); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetwindoworgexRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETWINDOWORGEX")

	if !w32.SetWindowOrgEx(ctx.MemDC, int(r.Origin.X), int(r.Origin.Y), nil) {
		log.Error("failed to run SetWindowOrgEx")
	}
}

type SetviewportextexRecord struct {
	Record
	Extent SizeL
}

func readSetviewportextexRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetviewportextexRecord{}
	r.Record = Record{Type: EMR_SETVIEWPORTEXTEX, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Extent); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetviewportextexRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETVIEWPORTEXTEX")

	if !w32.SetViewportExtEx(ctx.MemDC, int(r.Extent.Cx), int(r.Extent.Cy), nil) {
		log.Error("failed to run SetViewportExtEx")
	}
}

type SetviewportorgexRecord struct {
	Record
	Origin PointL
}

func readSetviewportorgexRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetviewportorgexRecord{}
	r.Record = Record{Type: EMR_SETVIEWPORTORGEX, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Origin); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetviewportorgexRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETVIEWPORTORGEX")

	if !w32.SetViewportOrgEx(ctx.MemDC, int(r.Origin.X), int(r.Origin.Y), nil) {
		log.Error("failed to run SetViewportOrgEx")
	}
}

type EofRecord struct {
	Record
	NPalEntries, OffPalEntries, SizeLast uint32
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

type SetmapmodeRecord struct {
	Record
	MapMode uint32
}

func readSetmapmodeRecord(reader *bytes.Reader, size uint32) (Recorder, error) {

	r := &SetmapmodeRecord{}
	r.Record = Record{Type: EMR_SETMAPMODE, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.MapMode); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetmapmodeRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETMAPMODE")

	if w32.SetMapMode(ctx.MemDC, int(r.MapMode)) == 0 {
		log.Error("failed to run SetMapMode")
	}
}

type SetbkmodeRecord struct {
	Record
	BkMode uint32
}

func readSetbkmodeRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetbkmodeRecord{}
	r.Record = Record{Type: EMR_SETBKMODE, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.BkMode); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetbkmodeRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETBKMODE")

	if w32.SetBkMode(ctx.MemDC, int(r.BkMode)) == 0 {
		log.Error("failed to run SetBkMode")
	}

}

type SetpolyfillmodeRecord struct {
	Record
	PolygonFillMode uint32
}

func readSetpolyfillmodeRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetpolyfillmodeRecord{}
	r.Record = Record{Type: EMR_SETPOLYFILLMODE, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.PolygonFillMode); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetpolyfillmodeRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETPOLYFILLMODE")

	if w32.SetPolyFillMode(ctx.MemDC, int(r.PolygonFillMode)) == 0 {
		log.Error("failed to run SetPolyFillMode")
	}
}

type SettextalignRecord struct {
	Record
	TextAlignmentMode uint32
}

func readSettextalignRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SettextalignRecord{}
	r.Record = Record{Type: EMR_SETTEXTALIGN, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.TextAlignmentMode); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SettextalignRecord) Draw(ctx *EmfContext) {
	log.Trace("DRAW EMR_SETTEXTALIGN")

	if w32.SetTextAlign(ctx.MemDC, uint(r.TextAlignmentMode)) == 0 {
		log.Error("failed to run SetTextAlign")
	}
}

type SetstretchbltmodeRecord struct {
	Record
	StretchMode uint32
}

func readSetstretchbltmodeRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetstretchbltmodeRecord{}
	r.Record = Record{Type: EMR_SETSTRETCHBLTMODE, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.StretchMode); err != nil {
		return nil, err
	}

	return r, nil
}

type SettextcolorRecord struct {
	Record
	Color ColorRef
}

func readSettextcolorRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SettextcolorRecord{}
	r.Record = Record{Type: EMR_SETTEXTCOLOR, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Color); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SettextcolorRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETTEXTCOLOR")

	if w32.SetTextColor(ctx.MemDC, w32.COLORREF(r.Color)) == w32.COLORREF(w32.CLR_INVALID) {
		log.Error("failed to run SetTextColor")
	}
}

type SetbkcolorRecord struct {
	Record
	Color ColorRef
}

func readSetbkcolorRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetbkcolorRecord{}
	r.Record = Record{Type: EMR_SETBKCOLOR, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Color); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetbkcolorRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETBKCOLOR")

	if w32.SetBkColor(ctx.MemDC, w32.COLORREF(r.Color)) == w32.COLORREF(w32.CLR_INVALID) {
		log.Error("failed to run SetBkColor")
	}
}

type MovetoexRecord struct {
	Record
	Offset PointL
}

func readMovetoexRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &MovetoexRecord{}
	r.Record = Record{Type: EMR_MOVETOEX, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Offset); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *MovetoexRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_MOVETOEX")

	if !w32.MoveToEx(ctx.MemDC, int(r.Offset.X), int(r.Offset.Y), nil) {
		log.Error("failed to run MoveToEx")
	}
}

type IntersectcliprectRecord struct {
	Record
	Clip RectL
}

func readIntersectcliprectRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &IntersectcliprectRecord{}
	r.Record = Record{Type: EMR_INTERSECTCLIPRECT, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Clip); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *IntersectcliprectRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_INTERSECTCLIPRECT")

	if w32.IntersectClipRect(ctx.MemDC, int(r.Clip.Left), int(r.Clip.Top), int(r.Clip.Right), int(r.Clip.Bottom)) == w32.ERROR {
		log.Error("failed to run IntersectClipRect")
	}
}

type SavedcRecord struct {
	Record
}

func readSavedcRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	return &SavedcRecord{Record: Record{Type: EMR_SAVEDC, Size: size}}, nil
}

func (r *SavedcRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SAVEDC")

	if w32.SaveDC(ctx.MemDC) == 0 {
		log.Error("failed to run SaveDC")
	}
}

type RestoredcRecord struct {
	Record
	SavedDC int32
}

func readRestoredcRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &RestoredcRecord{}
	r.Record = Record{Type: EMR_RESTOREDC, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.SavedDC); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *RestoredcRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_RESTOREDC")

	if !w32.RestoreDC(ctx.MemDC, int(r.SavedDC)) {
		log.Error("failed to run RestoreDC")
	}
}

type SetworldtransformRecord struct {
	Record
	XForm XForm
}

func readSetworldtransformRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetworldtransformRecord{}
	r.Record = Record{Type: EMR_SETWORLDTRANSFORM, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.XForm); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SetworldtransformRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETWORLDTRANSFORM")

	x32Xform := w32.XFORM(r.XForm)

	if !w32.SetWorldTransform(ctx.MemDC, &x32Xform) {
		log.Error("failed to run SetWorldTransform")
	}
}

type ModifyworldtransformRecord struct {
	Record
	XForm                    XForm
	ModifyWorldTransformMode uint32
}

func readModifyworldtransformRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &ModifyworldtransformRecord{}
	r.Record = Record{Type: EMR_MODIFYWORLDTRANSFORM, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.XForm); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.ModifyWorldTransformMode); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *ModifyworldtransformRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_MODIFYWORLDTRANSFORM")

	x32Xform := w32.XFORM(r.XForm)

	if !w32.ModifyWorldTransform(ctx.MemDC, &x32Xform, w32.DWORD(r.ModifyWorldTransformMode)) {
		log.Error("failed to run ModifyWorldTransform")
	}
}

type SelectobjectRecord struct {
	Record
	ihObject uint32
}

func readSelectobjectRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SelectobjectRecord{}
	r.Record = Record{Type: EMR_SELECTOBJECT, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.ihObject); err != nil {
		return nil, err
	}

	return r, nil
}

var StockObjects = map[uint32]interface{}{
	WHITE_BRUSH:         LogBrushEx{BrushStyle: BS_SOLID, ColorRef: 0x00FFFFFF},
	LTGRAY_BRUSH:        LogBrushEx{BrushStyle: BS_SOLID, ColorRef: 0x00C0C0C0},
	GRAY_BRUSH:          LogBrushEx{BrushStyle: BS_SOLID, ColorRef: 0x00808080},
	DKGRAY_BRUSH:        LogBrushEx{BrushStyle: BS_SOLID, ColorRef: 0x00404040},
	BLACK_BRUSH:         LogBrushEx{BrushStyle: BS_SOLID, ColorRef: 0x00000000},
	NULL_BRUSH:          LogBrushEx{BrushStyle: BS_NULL},
	WHITE_PEN:           LogPen{PenStyle: PS_COSMETIC | PS_SOLID, ColorRef: 0x00FFFFFF, Width: PointL{1, 0}},
	BLACK_PEN:           LogPen{PenStyle: PS_COSMETIC | PS_SOLID, ColorRef: 0x00000000, Width: PointL{1, 0}},
	NULL_PEN:            LogPen{PenStyle: PS_NULL},
	OEM_FIXED_FONT:      LogFont{CharSet: OEM_CHARSET, PitchAndFamily: (FF_DONTCARE<<4 + FIXED_PITCH)},
	ANSI_FIXED_FONT:     LogFont{CharSet: ANSI_CHARSET, PitchAndFamily: (FF_DONTCARE<<4 + FIXED_PITCH)},
	ANSI_VAR_FONT:       LogFont{CharSet: ANSI_CHARSET, PitchAndFamily: (FF_DONTCARE<<4 + VARIABLE_PITCH)},
	SYSTEM_FONT:         LogFont{Height: 11},
	DEVICE_DEFAULT_FONT: LogFont{Height: 11},
	SYSTEM_FIXED_FONT:   LogFont{Height: 11},
	DEFAULT_GUI_FONT:    LogFont{Height: 11},
}

// FIXME: handle following stockobject
// DEFAULT_PALETTE, DC_BRUSH, DC_PEN

func (r *SelectobjectRecord) Draw(ctx *EmfContext) {

	log.Trace("Draw EMR_SELECTOBJECT")

	gdiObject, ok := StockObjects[r.ihObject]
	if !ok {
		gdiObject, ok = ctx.Objects[r.ihObject]
		if !ok {
			log.Errorf("Object 0x%x not found\n", r.ihObject)
			return
		}
	}

	switch gdiObject.(type) {
	case LogPen:
		pen := gdiObject.(LogPen)
		plpen := w32.LOGPEN{
			LopnStyle: pen.PenStyle,
			LopnWidth: w32.POINT(pen.Width),
			LopnColor: w32.COLORREF(pen.ColorRef),
		}
		hpen := w32.CreatePenIndirect(&plpen)
		w32.SelectObject(ctx.MemDC, w32.HGDIOBJ(hpen))
	case LogPenEx:
		pen := gdiObject.(LogPenEx)
		plpen := w32.LOGPEN{
			LopnStyle: pen.PenStyle,
			LopnWidth: w32.POINT{X: int32(pen.Width), Y: 0},
			LopnColor: w32.COLORREF(pen.ColorRef),
		}
		hpen := w32.CreatePenIndirect(&plpen)
		w32.SelectObject(ctx.MemDC, w32.HGDIOBJ(hpen))
	case LogBrushEx:
		brushex := gdiObject.(LogBrushEx)
		plpen := w32.LOGBRUSH{
			LbStyle: brushex.BrushStyle,
			LbColor: w32.COLORREF(brushex.ColorRef),
			LbHatch: 0,
		}
		hbrush := w32.CreateBrushIndirect(&plpen)
		w32.SelectObject(ctx.MemDC, w32.HGDIOBJ(hbrush))
	case LogFont:
		w32logfont := gdiObject.(w32.LOGFONT)
		hlogfont := w32.CreateFontIndirectW(&w32logfont)
		w32.SelectObject(ctx.MemDC, w32.HGDIOBJ(hlogfont))
	}
}

type CreatepenRecord struct {
	Record
	ihPen  uint32
	LogPen LogPen
}

func readCreatepenRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &CreatepenRecord{}
	r.Record = Record{Type: EMR_CREATEPEN, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.ihPen); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.LogPen); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *CreatepenRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_CREATEPEN")

	ctx.Objects[r.ihPen] = r.LogPen
}

type CreatebrushindirectRecord struct {
	Record
	ihBrush  uint32
	LogBrush LogBrushEx
}

func readCreatebrushindirectRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &CreatebrushindirectRecord{}
	r.Record = Record{Type: EMR_CREATEBRUSHINDIRECT, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.ihBrush); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.LogBrush); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *CreatebrushindirectRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_CREATEBRUSHINDIRECT")

	ctx.Objects[r.ihBrush] = r.LogBrush
}

type DeleteobjectRecord struct {
	Record
	ihObject uint32
}

func readDeleteobjectRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &DeleteobjectRecord{}
	r.Record = Record{Type: EMR_DELETEOBJECT, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.ihObject); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *DeleteobjectRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_DELETEOBJECT")

	delete(ctx.Objects, r.ihObject)
}

type RectangleRecord struct {
	Record
	Box RectL
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

	if !w32.Rectangle(ctx.MemDC, int(r.Box.Left), int(r.Box.Top), int(r.Box.Right), int(r.Box.Bottom)) {
		log.Error("failed to run Rectangle")
	}
}

type ArcRecord struct {
	Record
	Box   RectL
	Start PointL
	End   PointL
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

	if !w32.Arc(ctx.MemDC, int(r.Box.Left), int(r.Box.Top), int(r.Box.Right), int(r.Box.Bottom),
		int(r.Start.X), int(r.Start.Y), int(r.End.X), int(r.End.Y)) {
		log.Error("failed to run Arc")
	}
}

type LinetoRecord struct {
	Record
	Point PointL
}

func readLinetoRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &LinetoRecord{}
	r.Record = Record{Type: EMR_LINETO, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Point); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *LinetoRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_LINETO")

	if !w32.LineTo(ctx.MemDC, int(r.Point.X), int(r.Point.Y)) {
		log.Error("failed to run LineTo")
	}
}

type BeginpathRecord struct {
	Record
}

func readBeginpathRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	return &BeginpathRecord{Record{Type: EMR_BEGINPATH, Size: size}}, nil
}

func (r *BeginpathRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_BEGINPATH")

	if !w32.BeginPath(ctx.MemDC) {
		log.Error("failed to run BeginPath")
	}
}

type EndpathRecord struct {
	Record
}

func readEndpathRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	return &EndpathRecord{Record{Type: EMR_ENDPATH, Size: size}}, nil
}

func (r *EndpathRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_ENDPATH")

	if !w32.EndPath(ctx.MemDC) {
		log.Error("failed to run EndPath")
	}
}

type ClosefigureRecord struct {
	Record
}

func readClosefigureRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	return &ClosefigureRecord{Record{Type: EMR_CLOSEFIGURE, Size: size}}, nil
}

func (r *ClosefigureRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_CLOSEFIGURE")

	if !w32.CloseFigure(ctx.MemDC) {
		log.Error("failed to run CloseFigure")
	}
}

type FillpathRecord struct {
	Record
	Bounds RectL
}

func readFillpathRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &FillpathRecord{}
	r.Record = Record{Type: EMR_FILLPATH, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Bounds); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *FillpathRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_FILLPATH")

	hrgn := w32.CreateRectRgn(int(r.Bounds.Left), int(r.Bounds.Top), int(r.Bounds.Right), int(r.Bounds.Bottom))
	w32.SelectObject(ctx.MemDC, w32.HGDIOBJ(hrgn))

	if !w32.FillPath(ctx.MemDC) {
		log.Error("failed to run FillPath")
	}
}

type StrokeandfillpathRecord struct {
	Record
	Bounds RectL
}

func readStrokeandfillpathRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &StrokeandfillpathRecord{}
	r.Record = Record{Type: EMR_STROKEANDFILLPATH, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Bounds); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *StrokeandfillpathRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_STROKEANDFILLPATH")

	hrgn := w32.CreateRectRgn(int(r.Bounds.Left), int(r.Bounds.Top), int(r.Bounds.Right), int(r.Bounds.Bottom))
	w32.SelectObject(ctx.MemDC, w32.HGDIOBJ(hrgn))

	if !w32.StrokeAndFillPath(ctx.MemDC) {
		log.Error("failed to run StrokeAndFillPath")
	}
}

type StrokepathRecord struct {
	Record
	Bounds RectL
}

func readStrokepathRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &StrokepathRecord{}
	r.Record = Record{Type: EMR_STROKEPATH, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Bounds); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *StrokepathRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_STROKEPATH")

	if !w32.StrokePath(ctx.MemDC) {
		log.Error("failed to run StrokePath")
	}
}

type SelectclippathRecord struct {
	Record
	RegionMode uint32
}

func readSelectclippathRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SelectclippathRecord{}
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
	reader.Seek(int64(r.Size-8), os.SEEK_CUR)
	return r, nil
}

type ExtcreatefontindirectwRecord struct {
	Record
	ihFonts uint32
	elw     LogFont
}

func readExtcreatefontindirectwRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &ExtcreatefontindirectwRecord{}
	r.Record = Record{Type: EMR_EXTCREATEFONTINDIRECTW, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.ihFonts); err != nil {
		return nil, err
	}

	var err error

	r.elw, err = readLogFont(reader)
	if err != nil {
		return nil, err
	}

	// skip the rest because we read only limited amount of data (LogFont) here
	reader.Seek(int64(r.Size-(12+92)), os.SEEK_CUR)

	return r, nil
}

func (r *ExtcreatefontindirectwRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_EXTCREATEFONTINDIRECTW")

	ctx.Objects[r.ihFonts] = r.elw
}

type ExttextoutwRecord struct {
	Record
	Bounds           RectL
	iGraphicsMode    uint32
	exScale, eyScale float32
	wEmrText         EmrText
}

func readExttextoutwRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &ExttextoutwRecord{}
	r.Record = Record{Type: EMR_EXTTEXTOUTW, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Bounds); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.iGraphicsMode); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.exScale); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.eyScale); err != nil {
		return nil, err
	}

	var err error
	r.wEmrText, err = readEmrText(reader, reader.Len()+36)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (m *ExttextoutwRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_EXTTEXTOUTW")

	width := float64(m.Bounds.Right - m.Bounds.Left)
	height := float64(m.Bounds.Bottom - m.Bounds.Top)
	text := canvas.NewTextBox(*ctx.FontFace, m.wEmrText.OutputString, width, height, canvas.Left, canvas.Center, 0, 100)
	ctx.DrawText(float64(m.Bounds.Left), float64(m.Bounds.Top), text)
}

type Polybezier16Record struct {
	Record
	Bounds  RectL
	Count   uint32
	aPoints []PointS
}

func readPolybezier16Record(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &Polybezier16Record{}
	r.Record = Record{Type: EMR_POLYBEZIER16, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Bounds); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.Count); err != nil {
		return nil, err
	}

	r.aPoints = make([]PointS, r.Count)
	if err := binary.Read(reader, binary.LittleEndian, &r.aPoints); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Polybezier16Record) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_POLYBEZIER16")

	ctx.MoveTo(float64(r.aPoints[0].X), float64(r.aPoints[0].Y))
	for i := 1; i < int(r.Count); i = i + 3 {
		ctx.CubeTo(
			float64(r.aPoints[i].X), float64(r.aPoints[i].Y),
			float64(r.aPoints[i+1].X), float64(r.aPoints[i+1].Y),
			float64(r.aPoints[i+2].X), float64(r.aPoints[i+2].Y),
		)
	}
	ctx.Stroke()
}

type Polygon16Record struct {
	Record
	Bounds  RectL
	Count   uint32
	aPoints []PointS
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

	r.aPoints = make([]PointS, r.Count)
	if err := binary.Read(reader, binary.LittleEndian, &r.aPoints); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Polygon16Record) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_POLYGON16")

	ctx.MoveTo(float64(r.aPoints[0].X), float64(r.aPoints[0].Y))
	for i := 1; i < int(r.Count); i++ {
		ctx.LineTo(float64(r.aPoints[i].X), float64(r.aPoints[i].Y))
	}
	ctx.Close()
	ctx.FillStroke()
}

type Polyline16Record struct {
	Record
	Bounds  RectL
	Count   uint32
	aPoints []PointS
}

func readPolyline16Record(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &Polyline16Record{}
	r.Record = Record{Type: EMR_POLYLINE16, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Bounds); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.Count); err != nil {
		return nil, err
	}

	r.aPoints = make([]PointS, r.Count)
	if err := binary.Read(reader, binary.LittleEndian, &r.aPoints); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Polyline16Record) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_POLYLINE16")

	ctx.MoveTo(float64(r.aPoints[0].X), float64(r.aPoints[0].Y))
	for i := 1; i < int(r.Count); i++ {
		ctx.LineTo(float64(r.aPoints[i].X), float64(r.aPoints[i].Y))
	}
	// ctx.Stroke()
}

type Polybezierto16Record struct {
	Record
	Bounds  RectL
	Count   uint32
	aPoints []PointS
}

func readPolybezierto16Record(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &Polybezierto16Record{}
	r.Record = Record{Type: EMR_POLYBEZIERTO16, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Bounds); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.Count); err != nil {
		return nil, err
	}

	r.aPoints = make([]PointS, r.Count)
	if err := binary.Read(reader, binary.LittleEndian, &r.aPoints); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Polybezierto16Record) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_POLYBEZIERTO16")

	for i := 0; i < int(r.Count); i = i + 3 {
		ctx.CubeTo(
			float64(r.aPoints[i].X), float64(r.aPoints[i].Y),
			float64(r.aPoints[i+1].X), float64(r.aPoints[i+1].Y),
			float64(r.aPoints[i+2].X), float64(r.aPoints[i+2].Y),
		)
	}

	ctx.Stroke()
}

type Polylineto16Record struct {
	Record
	Bounds  RectL
	Count   uint32
	aPoints []PointS
}

func readPolylineto16Record(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &Polylineto16Record{}
	r.Record = Record{Type: EMR_POLYLINETO16, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Bounds); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.Count); err != nil {
		return nil, err
	}

	r.aPoints = make([]PointS, r.Count)
	if err := binary.Read(reader, binary.LittleEndian, &r.aPoints); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Polylineto16Record) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_POLYLINETO16")

	for i := 0; i < int(r.Count); i++ {
		ctx.LineTo(float64(r.aPoints[i].X), float64(r.aPoints[i].Y))
	}
}

type Polypolygon16Record struct {
	Record
	Bounds            RectL
	NumberOfPolygons  uint32
	Count             uint32
	PolygonPointCount []uint32
	aPoints           []PointS
}

func readPolypolygon16Record(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &Polypolygon16Record{}
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

	r.aPoints = make([]PointS, r.Count)
	if err := binary.Read(reader, binary.LittleEndian, &r.aPoints); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Polypolygon16Record) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_POLYPOLYGON16")

	idx := 0
	for p := 0; p < int(r.NumberOfPolygons); p++ {
		pCount := int(r.PolygonPointCount[p])
		ctx.MoveTo(float64(r.aPoints[idx].X), float64(r.aPoints[idx].Y))
		for i := 1; i < pCount; i++ {
			ctx.LineTo(float64(r.aPoints[idx+i].X), float64(r.aPoints[idx+i].Y))
		}
		idx += pCount
		ctx.Close()
	}
	ctx.FillStroke()
}

type ExtcreatepenRecord struct {
	Record
	ihPen           uint32
	offBmi, cbBmi   uint32
	offBits, cbBits uint32
	elp             LogPenEx
	BmiSrc          DibHeaderInfo
	BitsSrc         []byte
}

func readExtcreatepenRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &ExtcreatepenRecord{}
	r.Record = Record{Type: EMR_EXTCREATEPEN, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.ihPen); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.offBmi); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.cbBmi); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.offBits); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.cbBits); err != nil {
		return nil, err
	}

	var err error
	r.elp, err = readLogPenEx(reader)
	if err != nil {
		return nil, err
	}

	// offset for bitmap info less than possible minimum
	// assuming there is no bitmap
	if r.offBmi < 52 {
		return r, nil
	}

	// BitmapBuffer
	// skipping UndefinedSpace
	reader.Seek(int64(r.offBmi-52-(r.elp.NumStyleEntries*4)), os.SEEK_CUR)

	// record does not contain bitmap
	if r.cbBmi == 0 {
		return r, nil
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.BmiSrc); err != nil {
		return nil, err
	}

	r.BitsSrc = make([]byte, r.cbBits)
	if _, err := reader.Read(r.BitsSrc); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *ExtcreatepenRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_EXTCREATEPEN")

	ctx.Objects[r.ihPen] = r.elp
}

type SeticmmodeRecord struct {
	Record
	ICMMode uint32
}

func readSeticmmodeRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SeticmmodeRecord{}
	r.Record = Record{Type: EMR_SETICMMODE, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.ICMMode); err != nil {
		return nil, err
	}

	return r, nil
}
