package emf

import (
	"bytes"

	"github.com/lokks307/go-emf/w32"
	log "github.com/sirupsen/logrus"
)

type EmfFile struct {
	Header  *HeaderRecord
	Records []Recorder
	Eof     *EofRecord
}

func ReadFile(data []byte) *EmfFile {
	reader := bytes.NewReader(data)
	emfFile := &EmfFile{}

	for reader.Len() > 0 {
		rec, err := readRecord(reader)
		if err != nil {
			log.Error(err)
			break
		}

		switch rec := rec.(type) {
		case *HeaderRecord:
			emfFile.Header = rec
		case *EofRecord:
			emfFile.Eof = rec
		default:
			emfFile.Records = append(emfFile.Records, rec)
		}
	}

	return emfFile
}

type EmfContext struct {
	MemDC   w32.HDC
	Width   int
	Height  int
	Objects map[uint32]interface{}
	wo      *PointL
	vo      *PointL
	we      *SizeL
	ve      *SizeL
	mm      uint32
}

func (f *EmfFile) NewEmfContext(width, height int) *EmfContext {

	memDC := w32.CreateCompatibleDC(0)
	memBM := w32.CreateCompatibleBitmap(memDC, width, height)
	w32.SelectObject(memDC, w32.HGDIOBJ(memBM))

	return &EmfContext{
		MemDC:   memDC,
		Width:   width,
		Height:  height,
		mm:      MM_TEXT,
		Objects: make(map[uint32]interface{}),
	}
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
