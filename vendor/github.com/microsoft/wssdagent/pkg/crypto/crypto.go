package crypto

func EncryptSecret(pswd []byte) (*[]byte, error) {
	blob := NewBlob(pswd)
	var outblob *_DATA_BLOB
	bo := encryptData(blob, 0,0,0,0,0, outblob)
	if bo {
		return nil,nil
	}
	test := ToByteArray(outblob)
	return &test, nil
}