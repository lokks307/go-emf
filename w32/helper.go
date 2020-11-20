package w32

import (
	"syscall"
	"unicode/utf16"
	"unsafe"
)

var (
	modkernel32      = syscall.NewLazyDLL("kernel32.dll")
	procGetLastError = modkernel32.NewProc("GetLastError")
)

func UTF16PtrToString(cstr *uint16) string {
	if cstr != nil {
		us := make([]uint16, 0, 256)
		for p := uintptr(unsafe.Pointer(cstr)); ; p += 2 {
			u := *(*uint16)(unsafe.Pointer(p))
			if u == 0 {
				return string(utf16.Decode(us))
			}
			us = append(us, u)
		}
	}

	return ""
}

func GetGpStatus(s int32) string {
	switch s {
	case Ok:
		return "Ok"
	case GenericError:
		return "GenericError"
	case InvalidParameter:
		return "InvalidParameter"
	case OutOfMemory:
		return "OutOfMemory"
	case ObjectBusy:
		return "ObjectBusy"
	case InsufficientBuffer:
		return "InsufficientBuffer"
	case NotImplemented:
		return "NotImplemented"
	case Win32Error:
		return "Win32Error"
	case WrongState:
		return "WrongState"
	case Aborted:
		return "Aborted"
	case FileNotFound:
		return "FileNotFound"
	case ValueOverflow:
		return "ValueOverflow"
	case AccessDenied:
		return "AccessDenied"
	case UnknownImageFormat:
		return "UnknownImageFormat"
	case FontFamilyNotFound:
		return "FontFamilyNotFound"
	case FontStyleNotFound:
		return "FontStyleNotFound"
	case NotTrueTypeFont:
		return "NotTrueTypeFont"
	case UnsupportedGdiplusVersion:
		return "UnsupportedGdiplusVersion"
	case GdiplusNotInitialized:
		return "GdiplusNotInitialized"
	case PropertyNotFound:
		return "PropertyNotFound"
	case PropertyNotSupported:
		return "PropertyNotSupported"
	case ProfileNotFound:
		return "ProfileNotFound"
	}
	return "Unknown Status Value"
}

func MakeIntResource(id uint16) *uint16 {
	return (*uint16)(unsafe.Pointer(uintptr(id)))
}

func LOWORD(dw uint32) uint16 {
	return uint16(dw)
}

func HIWORD(dw uint32) uint16 {
	return uint16(dw >> 16 & 0xffff)
}

func BoolToBOOL(value bool) BOOL {
	if value {
		return 1
	}
	return 0
}

// these constants can be passed to VerQueryValueString as the item
const (
	CompanyName      = "CompanyName"
	FileDescription  = "FileDescription"
	FileVersion      = "FileVersion"
	LegalCopyright   = "LegalCopyright"
	LegalTrademarks  = "LegalTrademarks"
	OriginalFilename = "OriginalFilename"
	ProductVersion   = "ProductVersion"
	PrivateBuild     = "PrivateBuild"
	SpecialBuild     = "SpecialBuild"
)

func MAKEWPARAM(low, high uint16) uintptr {
	return uintptr(low) | uintptr(high)<<16
}

func MAKELPARAM(low, high uint16) uintptr {
	return MAKEWPARAM(low, high)
}

func GetLastError() uint32 {
	ret, _, _ := procGetLastError.Call()
	return uint32(ret)
}
