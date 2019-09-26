// Code generated mksyscall_windows.exe DO NOT EDIT

package crypto

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var _ unsafe.Pointer

// Do the interface allocations only once for common
// Errno values.
const (
	errnoERROR_IO_PENDING = 997
)

var (
	errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}

var (
	modcrypt32 = windows.NewLazySystemDLL("crypt32.dll")

	procCryptProtectData   = modcrypt32.NewProc("CryptProtectData")
	procCryptUnprotectData = modcrypt32.NewProc("CryptUnprotectData")
)

func encryptData(datain *_DATA_BLOB, varw uint16, vara uint16, varb uint16, varc uint16, vard uint16, dataout *_DATA_BLOB) (hr bool) {
	r0, _, _ := syscall.Syscall9(procCryptProtectData.Addr(), 7, uintptr(unsafe.Pointer(datain)), uintptr(varw), uintptr(vara), uintptr(varb), uintptr(varc), uintptr(vard), uintptr(unsafe.Pointer(dataout)), 0, 0)
	hr = r0 != 0
	return
}

func decryptData(datain *_DATA_BLOB, varw uint16, vara uint16, varb uint16, varc uint16, vard uint16, dataout *_DATA_BLOB) (hr bool) {
	r0, _, _ := syscall.Syscall9(procCryptUnprotectData.Addr(), 7, uintptr(unsafe.Pointer(datain)), uintptr(varw), uintptr(vara), uintptr(varb), uintptr(varc), uintptr(vard), uintptr(unsafe.Pointer(dataout)), 0, 0)
	hr = r0 != 0
	return
}
