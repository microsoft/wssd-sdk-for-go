package crypto

import (
	"fmt"
	"unsafe"
)

func EncryptSecret(secret []byte) (*[]byte, error) {
	inBlob := NewBlob(secret)
	outblob := NewBlob([]byte(""))
	winBool := encryptData(inBlob, nil, nil, nil, nil, nil, outblob)
	if !winBool {
		return nil, fmt.Errorf("Could Not Encrypt Secret")
	}

	defer localFree(uintptr(unsafe.Pointer(outblob.pbData)))
	return ToByteArray(outblob), nil
}

func DecryptSecret(secret []byte) (*[]byte, error) {
	inBlob := NewBlob(secret)
	outblob := NewBlob([]byte(""))
	winBool := decryptData(inBlob, nil, nil, nil, nil, nil, outblob)
	if !winBool {
		return nil, fmt.Errorf("Could Not Decrypt Secret")
	}
	
	defer localFree(uintptr(unsafe.Pointer(outblob.pbData)))
	return ToByteArray(outblob), nil
}