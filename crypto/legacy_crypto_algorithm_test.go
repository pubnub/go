package crypto

import (
	"bytes"
	"io"
	"testing"
	"testing/quick"
)

func legacyCanDecryptEncryptStreamResult(in []byte) bool {
	cryptoAlgorithm, e := NewLegacyCryptoAlgorithm("enigma", true)
	if e != nil {
		return false
	}
	output, err := cryptoAlgorithm.EncryptStream(bytes.NewReader(in))
	if err != nil {
		return false
	}

	encryptedData, err := io.ReadAll(output.Reader)

	decrypted, err := cryptoAlgorithm.Decrypt(&EncryptedData{
		Data:     encryptedData,
		Metadata: output.Metadata,
	})
	if err != nil {
		return false
	}
	return bytes.Equal(in, decrypted[16:])
}

func legacyCanDecryptStreamEncryptResult(in []byte) bool {
	cryptor, e := NewLegacyCryptoAlgorithm("enigma", true)
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
	c := quick.Config{MaxCount: 10000}

	if err := quick.Check(legacyCanDecryptEncryptStreamResult, &c); err != nil {
		t.Error(err)
	}
}

func Test_Legacy_DecryptStream(t *testing.T) {
	c := quick.Config{MaxCount: 10000}
	if err := quick.Check(legacyCanDecryptStreamEncryptResult, &c); err != nil {
		t.Error(err)
	}
}
