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

		w32.DeleteObject(w32.HGDIOBJ(hbitmap))
		w32.DeleteDC(srcDC)

	}
}

type MaskAdditionInfo struct {
	XMask       int32
	YMask       int32
	UsageMask   uint32
	OffBmiMask  uint32
	CbBmiMask   uint32
	OffBitsMask uint32
	CbBitsMask  uint32
}

type MaskBltRecord struct {
	Record           // 8 bytes
	CommonBitmapInfo // 92 bytes
	MaskAdditionInfo // 28 bytes
	BmiSrc           w32.BITMAPINFO
	BitsSrc          []byte
	BmiMask          w32.BITMAPINFO
	BitsMask         []byte
}

func readMaskBltRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &MaskBltRecord{}
	r.Record = Record{Type: EMR_MASKBLT, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.CommonBitmapInfo); err != nil {
		return r, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.MaskAdditionInfo); err != nil {
		return r, err
	}

	// BitmapBuffer

	sizeUndefinedSpace1 := r.OffBmiSrc - 128
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

	sizeUndefinedSpace3 := r.OffBmiMask - r.OffBitsSrc - r.CbBitsSrc
	if sizeUndefinedSpace3 > 0 {
		reader.Seek(int64(sizeUndefinedSpace3), os.SEEK_CUR) // skipping UndefinedSpace1
	}

	if err := binary.Read(reader, binary.LittleEndian, &r.BmiMask.BITMAPINFOHEADER); err != nil {
		return nil, err
	}

	sizeUndefinedSpace4 := r.OffBitsMask - r.OffBmiMask - r.CbBitsMask
	if sizeUndefinedSpace4 > 0 {
		reader.Seek(int64(sizeUndefinedSpace4), os.SEEK_CUR) // skipping UndefinedSpace1
	}

	r.BitsMask = make([]byte, r.CbBitsMask)
	if _, err := reader.Read(r.BitsMask); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *MaskBltRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_MASKBLT")

	BitsData := PixelConvert(r.BitsSrc, int(r.BmiSrc.BiWidth), int(-r.BmiSrc.BiHeight), int(r.BmiSrc.BiBitCount), ctx.BitCount)
	r.BmiSrc.BiBitCount = uint16(ctx.BitCount)

	hbitmap := w32.CreateBitmap(int(r.XSrc), int(r.YSrc), w32.UINT(r.BmiSrc.BiPlanes), w32.UINT(r.BmiSrc.BiBitCount), BitsData)
	srcDC := w32.CreateCompatibleDC(ctx.MDC)
	w32.SelectObject(srcDC, w32.HGDIOBJ(hbitmap))

	maskBitmap := w32.CreateBitmap(int(r.BmiMask.BiWidth), int(-r.BmiMask.BiHeight), w32.UINT(r.BmiMask.BiPlanes), w32.UINT(r.BmiMask.BiBitCount), r.BitsMask)

	if !w32.MaskBlt(
		ctx.MDC, int(r.XDest), int(r.YDest), int(r.CxDest), int(r.CyDest), // dest
		srcDC, int(r.XSrc), int(r.YSrc), // src
		maskBitmap, int(r.XMask), int(r.YMask), // mask
		w32.DWORD(r.BitBltROP)) {
		log.Error("failed to run MaskBlt")
	}

	w32.DeleteObject(w32.HGDIOBJ(hbitmap))
	w32.DeleteDC(srcDC)

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

		w32.DeleteObject(w32.HGDIOBJ(hbitmap))
		w32.DeleteDC(srcDC)
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

type SetDIBitsToDeviceInfo struct {
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
	IStartScan uint32
	CScans     uint32 // 68 bytes
}

type SetDIBitsToDeviceRecord struct {
	Record                // 8 bytes
	SetDIBitsToDeviceInfo // 68 bytes
	BmiSrc                w32.BITMAPINFO
	BitsSrc               []byte
}

func readSetDIBitsToDeviceRecord(reader *bytes.Reader, size uint32) (Recorder, error) {
	r := &SetDIBitsToDeviceRecord{}
	r.Record = Record{Type: EMR_SETDIBITSTODEVICE, Size: size}

	if err := binary.Read(reader, binary.LittleEndian, &r.SetDIBitsToDeviceInfo); err != nil {
		return nil, err
	}

	if r.OffBmiSrc > 0 {

		sizeUndefinedSpace1 := r.OffBmiSrc - 76
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

func (r *SetDIBitsToDeviceRecord) Draw(ctx *EmfContext) {
	log.Trace("Draw EMR_SETDIBITSTODEVICE")

	if r.OffBmiSrc > 0 {

		BitsData := PixelConvert(r.BitsSrc, int(r.BmiSrc.BiWidth), int(-r.BmiSrc.BiHeight), int(r.BmiSrc.BiBitCount), ctx.BitCount)
		r.BmiSrc.BiBitCount = uint16(ctx.BitCount)

		if w32.SetDIBitsToDevice(
			ctx.MDC, int(r.XDest), int(r.YDest), w32.DWORD(r.CxSrc), w32.DWORD(r.CySrc), // dest
			int(r.XSrc), int(r.YSrc), w32.UINT(r.IStartScan), w32.UINT(r.CScans), BitsData, &r.BmiSrc, // src
			w32.UINT(r.UsageSrc)) == 0 {
			log.Error("failed to run SetDIBitsToDevice")
		}
	}
}
