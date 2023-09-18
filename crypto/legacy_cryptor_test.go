package crypto

import (
	"bytes"
	"io"
	"testing"
	"testing/quick"
)

func legacyCanDecryptEncryptStreamResult(in []byte) bool {
	cryptor, e := NewLegacyCryptor("enigma", true)
	if e != nil {
		return false
	}
	output, err := cryptor.EncryptStream(bytes.NewReader(in))
	if err != nil {
		return false
	}

	encr, err := io.ReadAll(output.Reader)

	decrypted, err := cryptor.Decrypt(&EncryptedData{
		Data:     encr,
		Metadata: output.Metadata,
	})
	if err != nil {
		return false
	}
	return bytes.Equal(in, decrypted)
}

func legacyCanDecryptStreamEncryptResult(in []byte) bool {
	cryptor, e := NewLegacyCryptor("enigma", true)
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

func Test_Legacy_EncryptStream(t *testing.T) {
	if err := quick.Check(legacyCanDecryptEncryptStreamResult, defaultPropertyTestConfig); err != nil {
		t.Error(err)
	}
}

func Test_Legacy_DecryptStream(t *testing.T) {
	if err := quick.Check(legacyCanDecryptStreamEncryptResult, defaultPropertyTestConfig); err != nil {
		t.Error(err)
	}
}
