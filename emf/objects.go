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

type XForm struct {
	M11, M12, M21, M22, Dx, Dy float32
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

// MS-WMF types
type ColorRef uint32

type SizeL struct {
	// MS-WMF says it's 32-bit unsigned integer
	// but there are files with negative values here
	Cx, Cy int32
}

type PointS struct {
	X, Y int16
}
type BitmapInfoHeader struct {
	HeaderSize                   uint32
	Width, Height                int32
	Planes, BitCount             uint16
	Compression, ImageSize       uint32
	XPelsPerMeter, YPelsPerMeter int32
	ColorUsed, ColorImportant    uint32
}

type DibHeaderInfo struct{}
