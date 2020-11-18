package emf

import (
	"bytes"
	"encoding/binary"
	"math"
	"os"
	"unicode/utf16"

	log "github.com/sirupsen/logrus"
	"github.com/tdewolff/canvas"
)

type LogPaletteEntry struct {
	_, Blue, Green, Red uint8
}

type LogPen struct {
	PenStyle uint32
	Width    PointL
	ColorRef ColorRef
}

type LogPenEx struct {
	PenStyle        uint32
	Width           uint32
	BrushStyle      uint32
	ColorRef        ColorRef
	BrushHatch      uint32
	NumStyleEntries uint32
	StyleEntry      []uint32
}

func readLogPenEx(reader *bytes.Reader) (LogPenEx, error) {
	r := LogPenEx{}
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

type LogBrushEx struct {
	BrushStyle uint32
	ColorRef   ColorRef
	BrushHatch uint32
}

type XForm struct {
	M11, M12, M21, M22, Dx, Dy float32
}

type EmrText struct {
	Reference    PointL
	Chars        uint32
	offString    uint32
	Options      uint32
	Rectangle    RectL
	offDx        uint32
	OutputString string
	OutputDx     []uint32
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
	r.OutputDx = make([]uint32, r.Chars)
	if err := binary.Read(reader, binary.LittleEndian, &r.OutputDx); err != nil {
		return r, err
	}

	return r, nil
}

type LogFont struct {
	Height, Width                        int32
	Escapement, Orientation, Weight      uint32
	Italic, Underline, StrikeOut         uint8
	CharSet, OutPrecision, ClipPrecision uint8
	Quality                              uint8
	PitchAndFamily                       uint8
	Facename                             string
}

func readLogFont(reader *bytes.Reader) (LogFont, error) {
	r := LogFont{}
	if err := binary.Read(reader, binary.LittleEndian, &r.Height); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.Width); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.Escapement); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.Orientation); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.Weight); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.Italic); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.Underline); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.StrikeOut); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.CharSet); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.OutPrecision); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.ClipPrecision); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.Quality); err != nil {
		return r, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &r.PitchAndFamily); err != nil {
		return r, err
	}

	b := make([]uint16, 32)
	if err := binary.Read(reader, binary.LittleEndian, &b); err != nil {
		return r, err
	}

	// raw, _, _ := transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder(), b)

	bTrim := make([]uint16, 0)
	for i := range b {
		if b[i] == 0x0000 {
			break
		}
		bTrim = append(bTrim, b[i])
	}

	for i := range bTrim {
		log.Tracef("%04x ", bTrim[i])
	}

	r.Facename = string(utf16.Decode(bTrim))

	log.Tracef("font name raw=%s\n", r.Facename)

	return r, nil
}

func (m LogFont) GetFontFace() *canvas.FontFace {

	var fontStyle canvas.FontStyle

	switch m.Weight {
	case 100:
		fontStyle = canvas.FontExtraLight
	case 200:
		fontStyle = canvas.FontLight
	case 300:
		fontStyle = canvas.FontBook
	case 500:
		fontStyle = canvas.FontMedium
	case 600:
		fontStyle = canvas.FontSemibold
	case 700:
		fontStyle = canvas.FontBold
	case 800:
		fontStyle = canvas.FontBlack
	case 900:
		fontStyle = canvas.FontExtraBlack
	default:
		fontStyle = canvas.FontRegular
	}

	if m.Italic == 0x01 {
		fontStyle = fontStyle | canvas.FontItalic
	}

	var fontDeco []canvas.FontDecorator

	if m.StrikeOut == 0x01 {
		fontDeco = append(fontDeco, canvas.FontStrikethrough)
	}

	if m.Underline == 0x01 {
		fontDeco = append(fontDeco, canvas.FontUnderline)
	}

	ff := canvas.NewFontFamily(m.Facename)
	_ = ff.LoadLocalFont(m.Facename, fontStyle)

	fontFace := ff.Face(math.Abs(float64(m.Height)), canvas.Black, fontStyle, canvas.FontNormal, fontDeco...)

	return &fontFace
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

type PointL struct {
	X, Y int32
}

type RectL struct {
	Left, Top, Right, Bottom int32
}

func (r RectL) Width() int32  { return r.Right - r.Left }
func (r RectL) Height() int32 { return r.Bottom - r.Top }

func (r RectL) Center() PointL {
	return PointL{
		X: r.Left + r.Width()/2,
		Y: r.Top + r.Height()/2,
	}
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
