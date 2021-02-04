package emf

import (
	"bytes"
	"encoding/binary"
	"os"
	"unicode/utf16"

	"github.com/lokks307/go-emf/w32"
)

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
	OffString    uint32
	Options      uint32
	Rectangle    w32.RECT
	OffDx        uint32
	OutputString []uint16
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
	if err := binary.Read(reader, binary.LittleEndian, &r.OffString); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.Options); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.Rectangle); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.OffDx); err != nil {
		return r, err
	}

	reader.Seek(int64(int(r.OffString)-(offset-reader.Len())), os.SEEK_CUR) // UndefinedSpace1

	r.OutputString = make([]uint16, r.Chars)
	if err := binary.Read(reader, binary.LittleEndian, &r.OutputString); err != nil {
		return r, err
	}

	reader.Seek(int64(int(r.OffDx)-(offset-reader.Len())), os.SEEK_CUR) // UndefinedSpace2

	r.OutputDx = make([]int32, r.Chars)
	if err := binary.Read(reader, binary.LittleEndian, &r.OutputDx); err != nil {
		return r, err
	}

	return r, nil
}

func (t *EmrText) GetString() string {
	return string(utf16.Decode(t.OutputString))
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
	ret := w32.COLORREF((uint32(m.Red)) | (uint32(m.Green) << 8) | (uint32(m.Blue) << 16))
	//log.Infof("RGB 0x%08x", ret)
	return ret
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
		PenStyle: m.PenStyle,
		Width:    m.Width,
		ColorRef: m.ColorRef.ColorRef(),
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
		PenStyle:        m.PenStyle,
		Width:           m.Width,
		BrushStyle:      m.BrushStyle,
		ColorRef:        m.ColorRef.ColorRef(),
		BrushHatch:      m.BrushHatch,
		NumStyleEntries: m.NumStyleEntries,
		StyleEntry:      m.StyleEntry,
	}
}

func readLogPalette(reader *bytes.Reader) (w32.LOGPALETTE, error) {
	r := w32.LOGPALETTE{}
	if err := binary.Read(reader, binary.LittleEndian, &r.Version); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.NumberOfEntries); err != nil {
		return r, err
	}
	r.PaletteEntries = make([]w32.COLORREF, r.NumberOfEntries)
	if err := binary.Read(reader, binary.LittleEndian, &r.PaletteEntries); err != nil {
		return r, err
	}

	return r, nil
}
