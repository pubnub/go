package crypto

import (
	"bytes"
	"io"
	"testing"
	"testing/quick"
)

func canDecryptEncryptStreamResult(in []byte) bool {
	cryptor, e := NewAesCbcCryptor("enigma")
	if e != nil {
		return false
	}

	output, err := cryptor.EncryptStream(bytes.NewReader(in))
	if err != nil {
		return false
	}

	encryptedData, err := io.ReadAll(output.Reader)

	decrypted, err := cryptor.Decrypt(&EncryptedData{
		Data:     encryptedData,
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
	c := quick.Config{MaxCount: 10000}
	if err := quick.Check(canDecryptEncryptStreamResult, &c); err != nil {
		t.Error(err)
	}
}

func Test_AesCBC_DecryptStream(t *testing.T) {
	c := quick.Config{MaxCount: 10000}
	if err := quick.Check(canDecryptStreamEncryptResult, &c); err != nil {
		t.Error(err)
	}
}
