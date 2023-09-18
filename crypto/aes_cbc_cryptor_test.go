package crypto

import (
	"bytes"
	"io"
	"testing"
	"testing/quick"
)

var defaultPropertyTestConfig = &quick.Config{
	MaxCount: 1000,
}

func canDecryptEncryptStreamResult(in []byte) bool {
	cryptor, e := NewAesCbcCryptor("enigma")
	if e != nil {
		return false
	}

	output, err := cryptor.EncryptStream(bytes.NewReader(in))
	if err != nil {
		return false
	}

	encrData, err := io.ReadAll(output.Reader)

	decrypted, err := cryptor.Decrypt(&EncryptedData{
		Data:     encrData,
		Metadata: output.Metadata,
	})
	if err != nil {
		return false
	}
	return bytes.Equal(in, decrypted)
}

func canDecryptStreamEncryptResult(in []byte) bool {
	cryptor, e := NewAesCbcCryptor("enigma")
	if e != nil {
		return false
	}

	output, err := cryptor.Encrypt(in)
	if err != nil {
		println(err.Error())
		return false
	}

	decryptingReader, err := cryptor.DecryptStream(&EncryptedStreamData{
		Reader:   bytes.NewReader(output.Data),
		Metadata: output.Metadata,
	})
	if err != nil {
		println(err.Error())
		return false
	}

	decrypted, err := io.ReadAll(decryptingReader)
	if err != nil {
		println(err.Error())
		return false
	}

	return bytes.Equal(in, decrypted)
}

func Test_AesCBC_EncryptStream(t *testing.T) {
	if err := quick.Check(canDecryptEncryptStreamResult, defaultPropertyTestConfig); err != nil {
		t.Error(err)
	}
}

func Test_AesCBC_DecryptStream(t *testing.T) {
	if err := quick.Check(canDecryptStreamEncryptResult, defaultPropertyTestConfig); err != nil {
		t.Error(err)
	}
}
