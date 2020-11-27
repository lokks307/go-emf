package emf

import (
	"bytes"
	"encoding/binary"
	"os"

	"github.com/lokks307/go-emf/w32"
	log "github.com/sirupsen/logrus"
)

type CommonBitmapInfo struct {
	Bounds     w32.RECT
	XDest      int32
	YDest      int32
	CxDest     int32
	CyDest     int32
	BitBltROP  uint32
	XSrc       int32
	YSrc       int32
	XformSrc   w32.XFORM
	BkColorSrc w32.COLORREF
	UsageSrc   uint32
	OffBmiSrc  uint32
	CbBmiSrc   uint32
	OffBitsSrc uint32
	CbBitsSrc  uint32
}

type BitBltRecord struct {
	Record           // 8 bytes
	CommonBitmapInfo // 92 bytes
	BmiSrc           w32.BITMAPINFO
	BitsSrc          []byte
}

func readBitBltRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &BitBltRecord{}
	r.Record = Record{Type: EMR_BITBLT, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.CommonBitmapInfo); err != nil {
		return r, err
	}

	// BitmapBuffer

	if r.OffBmiSrc > 0 {

		sizeUndefinedSpace1 := r.OffBmiSrc - 100
		if sizeUndefinedSpace1 > 0 {
			reader.Seek(int64(sizeUndefinedSpace1), os.SEEK_CUR) // skipping UndefinedSpace1
		}

		if err := binary.Read(reader, binary.LittleEndian, &r.BmiSrc.BITMAPINFOHEADER); err != nil {
			return nil, err
		}

		sizeUndefinedSpace2 := r.OffBitsSrc - r.OffBmiSrc - r.CbBmiSrc
		if sizeUndefinedSpace2 > 0 {
			reader.Seek(int64(sizeUndefinedSpace2), os.SEEK_CUR) // skipping UndefinedSpace2
		}

		r.BitsSrc = make([]byte, r.CbBitsSrc)
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

	if r.OffBmiSrc > 0 {

		BitsData := PixelConvert(r.BitsSrc, int(r.BmiSrc.BiWidth), int(-r.BmiSrc.BiHeight), int(r.BmiSrc.BiBitCount), ctx.BitCount)
		r.BmiSrc.BiBitCount = uint16(ctx.BitCount)

		hbitmap := w32.CreateBitmap(int(r.XSrc), int(r.YSrc), w32.UINT(r.BmiSrc.BiPlanes), w32.UINT(r.BmiSrc.BiBitCount), BitsData)
		srcDC := w32.CreateCompatibleDC(ctx.MDC)
		w32.SelectObject(srcDC, w32.HGDIOBJ(hbitmap))

		if !w32.BitBlt(
			ctx.MDC, int(r.XDest), int(r.YDest), int(r.CxDest), int(r.CyDest), // dest
			srcDC, int(r.XSrc), int(r.YSrc), // src
			w32.DWORD(r.BitBltROP)) {
			log.Error("failed to run BitBlt")
		}

	}
}

type StretchbltRecord struct {
	Record           // 8 bytes
	CommonBitmapInfo // 92 bytes
	CxSrc            int32
	CySrc            int32
	BmiSrc           w32.BITMAPINFO
	BitsSrc          []byte
}

func readStretchBltRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &StretchbltRecord{}
	r.Record = Record{Type: EMR_STRETCHBLT, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.CommonBitmapInfo); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.CxSrc); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.CySrc); err != nil {
		return r, err
	}

	// BitmapBuffer

	if r.OffBmiSrc > 0 {

		sizeUndefinedSpace1 := r.OffBmiSrc - 108
		if sizeUndefinedSpace1 > 0 {
			reader.Seek(int64(sizeUndefinedSpace1), os.SEEK_CUR) // skipping UndefinedSpace1
		}

		if err := binary.Read(reader, binary.LittleEndian, &r.BmiSrc.BITMAPINFOHEADER); err != nil {
			return nil, err
		}

		sizeUndefinedSpace2 := r.OffBitsSrc - r.OffBmiSrc - r.CbBmiSrc
		if sizeUndefinedSpace2 > 0 {
			reader.Seek(int64(sizeUndefinedSpace2), os.SEEK_CUR) // skipping UndefinedSpace2
		}

		r.BitsSrc = make([]byte, r.CbBitsSrc)
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

	if r.OffBmiSrc > 0 {

		BitsData := PixelConvert(r.BitsSrc, int(r.BmiSrc.BiWidth), int(-r.BmiSrc.BiHeight), int(r.BmiSrc.BiBitCount), ctx.BitCount)
		r.BmiSrc.BiBitCount = uint16(ctx.BitCount)

		hbitmap := w32.CreateBitmap(int(r.XSrc), int(r.YSrc), w32.UINT(r.BmiSrc.BiPlanes), w32.UINT(r.BmiSrc.BiBitCount), BitsData)
		srcDC := w32.CreateCompatibleDC(ctx.MDC)
		w32.SelectObject(srcDC, w32.HGDIOBJ(hbitmap))

		if !w32.StretchBlt(
			ctx.MDC, int(r.XDest), int(r.YDest), int(r.CxDest), int(r.CyDest), // dest
			srcDC, int(r.XSrc), int(r.YSrc), int(r.CxSrc), int(r.CySrc), // src
			w32.DWORD(r.BitBltROP)) {
			log.Error("failed to run StretchBlt")
		}
	}
}

type StretchDIBitsInfo struct {
	Bounds     w32.RECT
	XDest      int32
	YDest      int32
	XSrc       int32
	YSrc       int32
	CxSrc      int32
	CySrc      int32
	OffBmiSrc  uint32
	CbBmiSrc   uint32
	OffBitsSrc uint32
	CbBitsSrc  uint32
	UsageSrc   uint32
	BitBltROP  uint32
	CxDest     int32
	CyDest     int32 // 72 bytes
}
type StretchDIBitsRecord struct {
	Record            // 8 bytes
	StretchDIBitsInfo // 72 bytes
	BmiSrc            w32.BITMAPINFO
	BitsSrc           []byte
}

func readStretchDIBitsRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &StretchDIBitsRecord{}
	r.Record = Record{Type: EMR_STRETCHDIBITS, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.StretchDIBitsInfo); err != nil {
		return nil, err
	}

	if r.OffBmiSrc > 0 {

		sizeUndefinedSpace1 := r.OffBmiSrc - 80
		if sizeUndefinedSpace1 > 0 {
			reader.Seek(int64(sizeUndefinedSpace1), os.SEEK_CUR) // skipping UndefinedSpace1
		}

		if err := binary.Read(reader, binary.LittleEndian, &r.BmiSrc.BITMAPINFOHEADER); err != nil {
			return nil, err
		}

		sizeUndefinedSpace2 := r.OffBitsSrc - r.OffBmiSrc - r.CbBmiSrc
		if sizeUndefinedSpace2 > 0 {
			reader.Seek(int64(sizeUndefinedSpace2), os.SEEK_CUR) // skipping UndefinedSpace2
		}

		r.BitsSrc = make([]byte, r.CbBitsSrc)
		if _, err := reader.Read(r.BitsSrc); err != nil {
			return nil, err
		}

	}

	return r, nil
}

func (r *StretchDIBitsRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_STRETCHDIBITS")

	hrgn := w32.CreateRectRgn(int(r.Bounds.Left), int(r.Bounds.Top), int(r.Bounds.Right), int(r.Bounds.Bottom))
	w32.SelectObject(ctx.MDC, w32.HGDIOBJ(hrgn))

	if r.OffBmiSrc > 0 {

		BitsData := PixelConvert(r.BitsSrc, int(r.BmiSrc.BiWidth), int(-r.BmiSrc.BiHeight), int(r.BmiSrc.BiBitCount), ctx.BitCount)
		r.BmiSrc.BiBitCount = uint16(ctx.BitCount)

		if w32.StretchDIBits(
			ctx.MDC, int(r.XDest), int(r.YDest), int(r.CxDest), int(r.CyDest), // dest
			int(r.XSrc), int(r.YSrc), int(r.CxSrc), int(r.CySrc), BitsData, &r.BmiSrc, // src
			w32.UINT(r.UsageSrc), w32.DWORD(r.BitBltROP)) == 0 {
			log.Error("failed to run StretchDIBits")
		}
	}
}
