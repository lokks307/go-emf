package w32

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	libkernel32 = windows.NewLazySystemDLL("kernel32.dll")

	globalAlloc  = libkernel32.NewProc("GlobalAlloc")
	globalFree   = libkernel32.NewProc("GlobalFree")
	globalLock   = libkernel32.NewProc("GlobalLock")
	globalUnlock = libkernel32.NewProc("GlobalUnlock")
)

func GlobalAlloc(uFlags uint32, dwBytes uintptr) HGLOBAL {
	ret, _, _ := syscall.Syscall(globalAlloc.Addr(), 2,
		uintptr(uFlags),
		dwBytes,
		0)

	return HGLOBAL(ret)
}
func GlobalFree(hMem HGLOBAL) HGLOBAL {
	ret, _, _ := syscall.Syscall(globalFree.Addr(), 1,
		uintptr(hMem),
		0,
		0)

	return HGLOBAL(ret)
}

func GlobalLock(hMem HGLOBAL) unsafe.Pointer {
	ret, _, _ := syscall.Syscall(globalLock.Addr(), 1,
		uintptr(hMem),
		0,
		0)

	return unsafe.Pointer(ret)
}

func GlobalUnlock(hMem HGLOBAL) bool {
	ret, _, _ := syscall.Syscall(globalUnlock.Addr(), 1,
		uintptr(hMem),
		0,
		0)

	return ret != 0
}
