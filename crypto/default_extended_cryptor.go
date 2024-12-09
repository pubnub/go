package crypto

import (
	"bytes"
	"io"
)

type defaultExtendedCryptor struct {
	Cryptor
}

func liftToExtendedCryptor(cryptor Cryptor) ExtendedCryptor {
	if extendedCryptor, ok := cryptor.(ExtendedCryptor); ok {
		return extendedCryptor
	} else {
		return newDefaultExtendedCryptor(cryptor)
	}
}

func newDefaultExtendedCryptor(cryptor Cryptor) ExtendedCryptor {
	return &defaultExtendedCryptor{
		cryptor,
	}
}

func (c *defaultExtendedCryptor) EncryptStream(input io.Reader) (*EncryptedStreamData, error) {
	inputBytes, e := io.ReadAll(input)
	if e != nil {
		return nil, e
	}

	encryptedData, e := c.Cryptor.Encrypt(inputBytes)
	if e != nil {
		return nil, e
	}

	return &EncryptedStreamData{
		Reader:   bytes.NewReader(encryptedData.Data),
		Metadata: encryptedData.Metadata,
	}, nil
}

func (c *defaultExtendedCryptor) DecryptStream(input *EncryptedStreamData) (io.Reader, error) {
	inputBytes, e := io.ReadAll(input.Reader)
	if e != nil {
		return nil, e
	}

	decryptedData, e := c.Cryptor.Decrypt(&EncryptedData{Data: inputBytes, Metadata: input.Metadata})
	if e != nil {
		return nil, e
	}

	return bytes.NewReader(decryptedData), nil
}
