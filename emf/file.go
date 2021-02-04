package emf

import (
	"bytes"
	"image"

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

func (f *EmfFile) DrawToImg(mode int) (image.Image, error) {
	emfdc := NewEmfContext(f.Header.Original.Bounds, f.Header.Original.Device)

	for idx := range f.Records {
		f.Records[idx].Draw(emfdc)
	}

	var img interface{}
	var err error

	if mode == DRAW_COLOR_IMAGE {
		img, err = emfdc.DrawToColorImage(PAGE_AREA)
		if err != nil {
			log.Error(err)
			return nil, err
		}

	} else {
		img, err = emfdc.DrawToGrayImage(PAGE_AREA)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}

	var imgx image.Image
	switch t := img.(type) {
	case *image.NRGBA:
		imgx = t
	case *image.Gray:
		imgx = t
	}

	return imgx, nil
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
