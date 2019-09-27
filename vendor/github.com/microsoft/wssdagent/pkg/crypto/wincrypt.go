package crypto

import "unsafe"

//go:generate go run ../mksyscall_windows.go -output zsyscall_windows.go wincrypt.go

//sys encryptData(datain *_DATA_BLOB, varw *uint16, vara *uint16, varb *uint16, varc *uint16, vard *uint16, dataout *_DATA_BLOB) (hr bool) = crypt32.CryptProtectData
//sys decryptData(datain *_DATA_BLOB, varw *uint16, vara *uint16, varb *uint16, varc *uint16, vard *uint16, dataout *_DATA_BLOB) (hr bool) = crypt32.CryptUnprotectData

type _DATA_BLOB struct {
	cbData uint32
	pbData *byte
}

func NewBlob(d []byte) *_DATA_BLOB {
	if len(d) == 0 {
		return &_DATA_BLOB{}
	}
	return &_DATA_BLOB{
		pbData: &d[0],
		cbData: uint32(len(d)),
	}
}

func ToByteArray(b *_DATA_BLOB) *[]byte {
	d := make([]byte, b.cbData)
	copy(d, (*[1 << 30]byte)(unsafe.Pointer(b.pbData))[:])
	return &d
}