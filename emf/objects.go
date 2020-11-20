package emf

import (
	"bytes"
	"encoding/binary"
	"os"
	"unicode/utf16"

	"github.com/lokks307/go-emf/w32"
)

type LogPaletteEntry struct {
	_, Blue, Green, Red uint8
}

func readLogPenEx(reader *bytes.Reader) (w32.LOGPENEX, error) {
	r := w32.LOGPENEX{}
	if err := binary.Read(reader, binary.LittleEndian, &r.PenStyle); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.Width); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.BrushStyle); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.ColorRef); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.BrushHatch); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.NumStyleEntries); err != nil {
		return r, err
	}

	if r.PenStyle == PS_USERSTYLE && r.NumStyleEntries > 0 {
		r.StyleEntry = make([]uint32, r.NumStyleEntries)
		if err := binary.Read(reader, binary.LittleEndian, &r.StyleEntry); err != nil {
			return r, err
		}
	}

	return r, nil
}

type EmrText struct {
	Reference    w32.POINT
	Chars        uint32
	offString    uint32
	Options      uint32
	Rectangle    w32.RECT
	offDx        uint32
	OutputString string
	OutputDx     []int32
}

func readEmrText(reader *bytes.Reader, offset int) (EmrText, error) {
	r := EmrText{}
	if err := binary.Read(reader, binary.LittleEndian, &r.Reference); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.Chars); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.offString); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.Options); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.Rectangle); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.offDx); err != nil {
		return r, err
	}

	// UndefinedSpace1
	reader.Seek(int64(int(r.offString)-(offset-reader.Len())), os.SEEK_CUR)
	b := make([]uint16, r.Chars)
	if err := binary.Read(reader, binary.LittleEndian, &b); err != nil {
		return r, err
	}
	r.OutputString = string(utf16.Decode(b))

	// UndefinedSpace2
	reader.Seek(int64(int(r.offDx)-(offset-reader.Len())), os.SEEK_CUR)
	r.OutputDx = make([]int32, r.Chars)
	if err := binary.Read(reader, binary.LittleEndian, &r.OutputDx); err != nil {
		return r, err
	}

	return r, nil
}

func readLogFont(reader *bytes.Reader) (w32.LOGFONT, error) {
	r := w32.LOGFONT{}
	if err := binary.Read(reader, binary.LittleEndian, &r); err != nil {
		return r, err
	}
	return r, nil
}

type PointS struct {
	X, Y int16
}

type RegionDataHeader struct {
	Size       uint32
	Type       uint32
	CountRects uint32
	RgnSize    uint32
	Bounds     w32.RECT
}

type RegionData struct {
	RegionDataHeader
	Data []w32.RECT
}

type WMFCOLORREF struct {
	Red      byte
	Green    byte
	Blue     byte
	Reseverd byte
}

func (m WMFCOLORREF) ColorRef() w32.COLORREF {

	red := uint32(m.Red)
	green := uint32(m.Green) << 4
	blue := uint32(m.Blue) << 8

	return w32.COLORREF(red | green | blue)
}

type WMFLOGBRUSH struct {
	BrushStyle uint32
	Color      WMFCOLORREF
	BrushHatch uint32
}

func (m WMFLOGBRUSH) LogBrush() w32.LOGBRUSH {
	return w32.LOGBRUSH{
		BrushStyle: m.BrushStyle,
		Color:      m.Color.ColorRef(),
		BrushHatch: m.BrushHatch,
	}
}

type WMFLOGPEN struct {
	PenStyle uint32
	Width    w32.POINT
	ColorRef WMFCOLORREF
}

func (m WMFLOGPEN) LogPen() w32.LOGPEN {
	return w32.LOGPEN{
		BrushStyle: m.BrushStyle,
		Width:      m.Width,
		ColorRef:   m.ColorRef.ColorRef(),
	}
}

type WMFLOGPENEX struct {
	PenStyle        uint32
	Width           uint32
	BrushStyle      uint32
	ColorRef        WMFCOLORREF
	BrushHatch      uint32
	NumStyleEntries uint32
	StyleEntry      []uint32
}

func (m WMFLOGPENEX) LogPenEx() w32.LOGPENEX {
	return w32.LOGPENEX{
		PenStyle: m.PenStyle,
		Width: m.Width,
		BrushStyle: m.BrushStyle
		ColorRef: m.ColorRef.ColorRef(),
		BrushHatch: m.BrushHatch,
		NumStyleEntries: m.NumStyleEntries,
		StyleEntry: m.StyleEntry
	}
}
