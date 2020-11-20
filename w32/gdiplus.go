package w32

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

var (
	gdiplus = syscall.NewLazyDLL("gdiplus.dll")

	gdipCreateBitmapFromFile     = gdiplus.NewProc("GdipCreateBitmapFromFile")
	gdipCreateHBITMAPFromBitmap  = gdiplus.NewProc("GdipCreateHBITMAPFromBitmap")
	gdipCreateBitmapFromResource = gdiplus.NewProc("GdipCreateBitmapFromResource")
	gdipCreateBitmapFromStream   = gdiplus.NewProc("GdipCreateBitmapFromStream")
	gdipDisposeImage             = gdiplus.NewProc("GdipDisposeImage")
	gdiplusShutdown              = gdiplus.NewProc("GdiplusShutdown")
	gdiplusStartup               = gdiplus.NewProc("GdiplusStartup")
)

func GdipCreateBitmapFromFile(filename string) (*uintptr, error) {
	var bitmap *uintptr
	ret, _, _ := gdipCreateBitmapFromFile.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(filename))),
		uintptr(unsafe.Pointer(&bitmap)),
	)

	if ret != Ok {
		return nil, errors.New(fmt.Sprintf(
			"GdipCreateBitmapFromFile failed with status '%s' for file '%s'",
			GetGpStatus(int32(ret)),
			filename,
		))
	}

	return bitmap, nil
}

func GdipCreateBitmapFromResource(instance HINSTANCE, resId *uint16) (*uintptr, error) {
	var bitmap *uintptr
	ret, _, _ := gdipCreateBitmapFromResource.Call(
		uintptr(instance),
		uintptr(unsafe.Pointer(resId)),
		uintptr(unsafe.Pointer(&bitmap)),
	)

	if ret != Ok {
		return nil, errors.New(fmt.Sprintf("GdiCreateBitmapFromResource failed with status '%s'", GetGpStatus(int32(ret))))
	}

	return bitmap, nil
}

func GdipCreateBitmapFromStream(stream *IStream) (*uintptr, error) {
	var bitmap *uintptr
	ret, _, _ := gdipCreateBitmapFromStream.Call(
		uintptr(unsafe.Pointer(stream)),
		uintptr(unsafe.Pointer(&bitmap)),
	)

	if ret != Ok {
		return nil, errors.New(fmt.Sprintf("GdipCreateBitmapFromStream failed with status '%s'", GetGpStatus(int32(ret))))
	}

	return bitmap, nil
}

func GdipCreateHBITMAPFromBitmap(bitmap *uintptr, background uint32) (HBITMAP, error) {
	var hbitmap HBITMAP
	ret, _, _ := gdipCreateHBITMAPFromBitmap.Call(
		uintptr(unsafe.Pointer(bitmap)),
		uintptr(unsafe.Pointer(&hbitmap)),
		uintptr(background),
	)

	if ret != Ok {
		return 0, errors.New(fmt.Sprintf("GdipCreateHBITMAPFromBitmap failed with status '%s'", GetGpStatus(int32(ret))))
	}

	return hbitmap, nil
}

func GdipDisposeImage(image *uintptr) {
	gdipDisposeImage.Call(uintptr(unsafe.Pointer(image)))
}

func GdiplusShutdown(token uintptr) {
	gdiplusShutdown.Call(token)
}

func GdiplusStartup(input *GdiplusStartupInput, output *GdiplusStartupOutput) (token uintptr, status uint32) {
	ret, _, _ := gdiplusStartup.Call(
		uintptr(unsafe.Pointer(&token)),
		uintptr(unsafe.Pointer(input)),
		uintptr(unsafe.Pointer(output)),
	)
	status = uint32(ret)
	return
}
