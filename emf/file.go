package emf

import (
	"bytes"
	"image"

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

	device := f.Header.Original.Device
	width := int(device.CX) + 1
	height := int(device.CY) + 1

	emfdc := NewEmfContext(width, height)

	for idx := range f.Records {
		f.Records[idx].Draw(emfdc)
	}

	var img *image.RGBA
	var err error

	img, err = emfdc.DrawToImage()
	if err != nil {
		log.Error(err)
		return err
	}

	return ImageToPNG(img, output)
}
