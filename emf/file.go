package emf

import (
	"bytes"

	log "github.com/sirupsen/logrus"
)

const (
	CROP_AREA = iota
	PAGE_AREA
)

const (
	DRAW_COLOR_IMAGE = iota
	DRAW_GRAY_IMAGE
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

func (f *EmfFile) DrawToGrayPNG(output string) error {
	return f.drawToPNG(output, DRAW_GRAY_IMAGE)
}

func (f *EmfFile) DrawToColorPNG(output string) error {
	return f.drawToPNG(output, DRAW_COLOR_IMAGE)
}

func (f *EmfFile) drawToPNG(output string, mode int) error {
	emfdc := NewEmfContext(f.Header.Original.Bounds, f.Header.Original.Device)

	for idx := range f.Records {
		f.Records[idx].Draw(emfdc)
	}

	if mode == DRAW_COLOR_IMAGE {
		img, err := emfdc.DrawToColorImage(PAGE_AREA)
		if err != nil {
			log.Error(err)
			return err
		}

		return ImageToPNG(img, output)
	} else {
		img, err := emfdc.DrawToGrayImage(PAGE_AREA)
		if err != nil {
			log.Error(err)
			return err
		}

		return ImageToPNG(img, output)
	}
}
