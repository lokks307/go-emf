package emf

import (
	"bytes"
	"encoding/binary"
	"os"

	"github.com/lokks307/go-emf/w32"
	log "github.com/sirupsen/logrus"
)

type CommonBitmapRecord struct {
	Record
	Bounds                       w32.RECT
	xDest, yDest, cxDest, cyDest int32
	BitBltRasterOperation        uint32
	xSrc, ySrc                   int32
	XformSrc                     w32.XFORM
	BkColorSrc                   w32.COLORREF
	UsageSrc                     uint32
	offBmiSrc, cbBmiSrc          uint32
	offBitsSrc, cbBitsSrc        uint32

	cxSrc, cySrc int32 // only for EMR_STRETCHBLT

	BmiSrc  w32.BITMAPINFO
	BitsSrc []byte
}
type BitBltRecord struct {
	CommonBitmapRecord
}

func readBitBltRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &BitBltRecord{}
	r.Record = Record{Type: EMR_BITBLT, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Bounds); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.xDest); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.yDest); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.cxDest); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.cyDest); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.BitBltRasterOperation); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.xSrc); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.ySrc); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.XformSrc); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.BkColorSrc); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.UsageSrc); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.offBmiSrc); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.cbBmiSrc); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.offBitsSrc); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.cbBitsSrc); err != nil {
		return r, err
	}

	// BitmapBuffer

	if r.offBmiSrc > 0 {

		sizeUndefinedSpace1 := r.offBmiSrc - 80
		if sizeUndefinedSpace1 > 0 {
			reader.Seek(int64(sizeUndefinedSpace1), os.SEEK_CUR) // skipping UndefinedSpace1
		}

		if err := binary.Read(reader, binary.LittleEndian, &r.BmiSrc.BITMAPINFOHEADER); err != nil {
			return nil, err
		}

		sizeUndefinedSpace2 := r.offBitsSrc - r.offBmiSrc - r.cbBmiSrc
		if sizeUndefinedSpace2 > 0 {
			reader.Seek(int64(sizeUndefinedSpace2), os.SEEK_CUR) // skipping UndefinedSpace2
		}

		r.BitsSrc = make([]byte, r.cbBitsSrc)
		if _, err := reader.Read(r.BitsSrc); err != nil {
			return nil, err
		}

	}

	return r, nil
}

func (r *BitBltRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_BITBLT")

	hrgn := w32.CreateRectRgn(int(r.Bounds.Left), int(r.Bounds.Top), int(r.Bounds.Right), int(r.Bounds.Bottom))
	w32.SelectObject(ctx.MDC, w32.HGDIOBJ(hrgn))

	if r.offBmiSrc > 0 {

		BitsData := PixelConvert(r.BitsSrc, int(r.cxSrc), int(r.cySrc), int(r.BmiSrc.BiBitCount), ctx.BitCount)
		r.BmiSrc.BiBitCount = uint16(ctx.BitCount)

		hbitmap := w32.CreateBitmap(int(r.xSrc), int(r.ySrc), w32.UINT(r.BmiSrc.BiPlanes), w32.UINT(r.BmiSrc.BiBitCount), BitsData)
		srcDC := w32.CreateCompatibleDC(ctx.MDC)
		w32.SelectObject(srcDC, w32.HGDIOBJ(hbitmap))

		if !w32.BitBlt(ctx.MDC, int(r.xDest), int(r.yDest), int(r.cxDest), int(r.cyDest), srcDC, int(r.xSrc), int(r.ySrc), w32.DWORD(r.BitBltRasterOperation)) {
			log.Error("failed to run BitBlt")
		}

	}
}

type StretchbltRecord struct {
	CommonBitmapRecord
}

func readStretchBltRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &StretchbltRecord{}
	r.Record = Record{Type: EMR_STRETCHBLT, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Bounds); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.xDest); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.yDest); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.cxDest); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.cyDest); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.BitBltRasterOperation); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.xSrc); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.ySrc); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.XformSrc); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.BkColorSrc); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.UsageSrc); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.offBmiSrc); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.cbBmiSrc); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.offBitsSrc); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.cbBitsSrc); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.cxSrc); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.cySrc); err != nil {
		return r, err
	}

	// BitmapBuffer

	if r.offBmiSrc > 0 {

		sizeUndefinedSpace1 := r.offBmiSrc - 80
		if sizeUndefinedSpace1 > 0 {
			reader.Seek(int64(sizeUndefinedSpace1), os.SEEK_CUR) // skipping UndefinedSpace1
		}

		if err := binary.Read(reader, binary.LittleEndian, &r.BmiSrc.BITMAPINFOHEADER); err != nil {
			return nil, err
		}

		sizeUndefinedSpace2 := r.offBitsSrc - r.offBmiSrc - r.cbBmiSrc
		if sizeUndefinedSpace2 > 0 {
			reader.Seek(int64(sizeUndefinedSpace2), os.SEEK_CUR) // skipping UndefinedSpace2
		}

		r.BitsSrc = make([]byte, r.cbBitsSrc)
		if _, err := reader.Read(r.BitsSrc); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (r *StretchbltRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_STRETCHBLT")

	hrgn := w32.CreateRectRgn(int(r.Bounds.Left), int(r.Bounds.Top), int(r.Bounds.Right), int(r.Bounds.Bottom))
	w32.SelectObject(ctx.MDC, w32.HGDIOBJ(hrgn))

	if r.offBmiSrc > 0 {

		BitsData := PixelConvert(r.BitsSrc, int(r.cxSrc), int(r.cySrc), int(r.BmiSrc.BiBitCount), ctx.BitCount)
		r.BmiSrc.BiBitCount = uint16(ctx.BitCount)

		hbitmap := w32.CreateBitmap(int(r.xSrc), int(r.ySrc), w32.UINT(r.BmiSrc.BiPlanes), w32.UINT(r.BmiSrc.BiBitCount), BitsData)
		srcDC := w32.CreateCompatibleDC(ctx.MDC)
		w32.SelectObject(srcDC, w32.HGDIOBJ(hbitmap))

		if !w32.StretchBlt(ctx.MDC, int(r.xDest), int(r.yDest), int(r.cxDest), int(r.cyDest), srcDC, int(r.xSrc), int(r.ySrc), int(r.cxSrc), int(r.cySrc), w32.DWORD(r.BitBltRasterOperation)) {
			log.Error("failed to run StretchBlt")
		}
	}
}

type StretchdibitsRecord struct {
	CommonBitmapRecord
}

func readStretchDIBitsRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &StretchdibitsRecord{}
	r.Record = Record{Type: EMR_STRETCHDIBITS, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.Bounds); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.xDest); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.yDest); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.xSrc); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.ySrc); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.cxSrc); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.cySrc); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.offBmiSrc); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.cbBmiSrc); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.offBitsSrc); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.cbBitsSrc); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.UsageSrc); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.BitBltRasterOperation); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.cxDest); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.cyDest); err != nil {
		return nil, err
	}

	// BitmapBuffer

	if r.offBmiSrc > 0 {

		sizeUndefinedSpace1 := r.offBmiSrc - 80
		if sizeUndefinedSpace1 > 0 {
			reader.Seek(int64(sizeUndefinedSpace1), os.SEEK_CUR) // skipping UndefinedSpace1
		}

		if err := binary.Read(reader, binary.LittleEndian, &r.BmiSrc.BITMAPINFOHEADER); err != nil {
			return nil, err
		}

		sizeUndefinedSpace2 := r.offBitsSrc - r.offBmiSrc - r.cbBmiSrc
		if sizeUndefinedSpace2 > 0 {
			reader.Seek(int64(sizeUndefinedSpace2), os.SEEK_CUR) // skipping UndefinedSpace2
		}

		r.BitsSrc = make([]byte, r.cbBitsSrc)
		if _, err := reader.Read(r.BitsSrc); err != nil {
			return nil, err
		}

	}

	return r, nil
}

func (r *StretchdibitsRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_STRETCHDIBITS")

	hrgn := w32.CreateRectRgn(int(r.Bounds.Left), int(r.Bounds.Top), int(r.Bounds.Right), int(r.Bounds.Bottom))
	w32.SelectObject(ctx.MDC, w32.HGDIOBJ(hrgn))

	if r.offBmiSrc > 0 {

		BitsData := PixelConvert(r.BitsSrc, int(r.cxSrc), int(r.cySrc), int(r.BmiSrc.BiBitCount), ctx.BitCount)
		r.BmiSrc.BiBitCount = uint16(ctx.BitCount)

		if w32.StretchDIBits(ctx.MDC, int(r.xDest), int(r.yDest), int(r.cxDest), int(r.cyDest), int(r.xSrc), int(r.ySrc), int(r.cxSrc), int(r.cySrc), BitsData, &r.BmiSrc, w32.UINT(r.UsageSrc), w32.DWORD(r.BitBltRasterOperation)) == 0 {
			log.Error("failed to run StretchDIBits")
		}
	}
}
