package emf

import (
	"bytes"

	log "github.com/sirupsen/logrus"
)

const (
	CROP_AREA = iota
	PAGE_AREA
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
			return emfFile
		default:
			emfFile.Records = append(emfFile.Records, rec)
		}
	}

	return emfFile
}

func (f *EmfFile) DrawToPNG(output string) error {
	emfdc := NewEmfContext(f.Header.Original.Bounds, f.Header.Original.Device)

	for idx := range f.Records {
		f.Records[idx].Draw(emfdc)
	}

	img, w, h, err := emfdc.DrawToImage(PAGE_AREA)
	if err != nil {
		log.Error(err)
		return err
	}

	return ImageToPNG(img, w, h, output)
}
