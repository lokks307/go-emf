package emf

import (
	"errors"
	"image"
	"image/png"
	"os"
	"unsafe"

	"github.com/lokks307/go-emf/w32"
	log "github.com/sirupsen/logrus"
)

func ImageToPNG(img *image.NRGBA, output string) error {

	var err error
	var outf *os.File

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

func PixelConvertFromMonochrome(src []byte, width, height, destBppBit int) []byte {
	if 1 == destBppBit {
		return src
	}

	destBppByte := destBppBit / 8

	destPadding := (4 - (width * destBppByte % 4)) % 4

	numPixels := width * height

	dest := make([]byte, numPixels*destBppByte+destPadding*height)

	// TODO : complete codes

	return dest
}

func PixelConvert(src []byte, width, height, srcBppBit, destBppBit int) []byte {

	if srcBppBit == destBppBit {
		return src
	}

	srcBppByte := srcBppBit / 8
	destBppByte := destBppBit / 8

	srcPadding := (4 - (width * srcBppByte % 4)) % 4
	destPadding := (4 - (width * destBppByte % 4)) % 4

	numPixels := len(src) / srcBppByte

	dest := make([]byte, numPixels*destBppByte+destPadding*height)

	var R, G, B int

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcPos := (y*width+x)*srcBppByte + srcPadding*y
			destPos := (y*width+x)*destBppByte + destPadding*y

			switch srcBppBit {
			case 16:
				R = int(src[srcPos]&0x0F) * 16
				G = int((src[srcPos+1]&0xF0)>>4) * 16
				B = int(src[srcPos+1]&0x0F) * 16
			case 24:
				fallthrough
			case 32:
				R = int(src[srcPos])
				G = int(src[srcPos+1])
				B = int(src[srcPos+2])
			}

			switch destBppBit {
			case 16:
			case 24:
				fallthrough
			case 32:
				dest[destPos] = byte(R)
				dest[destPos+1] = byte(G)
				dest[destPos+2] = byte(B)
			}
		}
	}

	return dest
}

func DeviceContextToImage(srcDC w32.HDC, width, height int) (*image.RGBA, error) {

	destDC := w32.CreateCompatibleDC(srcDC)

	if destDC == 0 {
		return nil, errors.New("CreateCompatibleDC failed")
	}
	defer w32.DeleteDC(destDC)

	bitmap := w32.CreateCompatibleBitmap(srcDC, width, height)

	oobj := w32.SelectObject(destDC, w32.HGDIOBJ(bitmap)) // attach bitmap to destDC
	if oobj == 0 {
		return nil, errors.New("SelectObject failed")
	}
	defer w32.SelectObject(destDC, oobj)

	if bitmap == 0 {
		return nil, errors.New("CreateCompatibleBitmap failed")
	}
	defer w32.DeleteObject(w32.HGDIOBJ(bitmap))

	var header w32.BITMAPINFO
	header.BiSize = uint32(unsafe.Sizeof(header))
	header.BiPlanes = 1
	header.BiBitCount = 32
	header.BiWidth = int32(width)
	header.BiHeight = int32(-height)
	header.BiCompression = w32.BI_RGB
	header.BiSizeImage = 0

	bitmapDataSize := uintptr(((int64(width)*int64(header.BiBitCount) + 31) / 32) * 4 * int64(height))
	hmem := w32.GlobalAlloc(w32.GMEM_MOVEABLE, bitmapDataSize)
	memptr := w32.GlobalLock(hmem)
	defer func() {
		w32.GlobalUnlock(hmem)
		w32.GlobalFree(hmem)
	}()

	if !w32.BitBlt(destDC, 0, 0, width, height, srcDC, 0, 0, w32.SRCCOPY) { // copy srcDC to destDC(bitmap)
		return nil, errors.New("BitBlt failed")
	}

	if w32.GetDIBits(destDC, bitmap, 0, w32.UINT(height), memptr, &header, w32.DIB_RGB_COLORS) == 0 { // bitmap on destDC to memptr
		return nil, errors.New("GetDIBits failed")
	}

	rect := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rect)

	i := 0
	src := uintptr(memptr)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			v0 := *(*uint8)(unsafe.Pointer(src))
			v1 := *(*uint8)(unsafe.Pointer(src + 1))
			v2 := *(*uint8)(unsafe.Pointer(src + 2))

			// BGRA => RGBA, and set A to 255
			img.Pix[i], img.Pix[i+1], img.Pix[i+2], img.Pix[i+3] = v2, v1, v0, 255

			i += 4
			src += 4
		}
	}

	return img, nil
}
