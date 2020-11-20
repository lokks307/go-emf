package emf

import (
	"bytes"
	"image"
	"image/png"
	"os"

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

func (f *EmfFile) DrawToPNG(output string) error {

	bounds := f.Header.Original.Bounds
	width := int(bounds.Width()) + 1
	height := int(bounds.Height()) + 1

	emfdc := NewEmfContext(width, height)

	for idx := range f.Records {
		f.Records[idx].Draw(emfdc)
	}

	var img *image.RGBA
	var err error
	var outf *os.File

	img, err = emfdc.DrawToImage()
	if err != nil {
		log.Error(err)
		return err
	}

	outf, err = os.Create(output)
	if err != nil {
		log.Error(err)
		return err
	}
	defer outf.Close()

	err = png.Encode(outf, img)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
