package utils

import (
	"errors"
	"io"
)

type Crypto struct {
	encryptor           CryptoAlgorithm
	decryptors          []CryptoAlgorithm
	defaultDecryptor    CryptoAlgorithm
	cryptoHeaderVersion CryptoHeaderVersion
}

func NewCrypto(encryptor CryptoAlgorithm, decryptors []CryptoAlgorithm, defaultDecryptor CryptoAlgorithm, version CryptoHeaderVersion) *Crypto {
	return &Crypto{
		encryptor:           encryptor,
		decryptors:          decryptors,
		defaultDecryptor:    defaultDecryptor,
		cryptoHeaderVersion: version,
	}
}

func (c *Crypto) Encrypt(message []byte) ([]byte, error) {
	r, e := c.encryptor.Encrypt(message)
	if e != nil {
		return nil, e
	}
	if c.encryptor.HeaderVersion() == headless {
		return r.Data, nil
	}
	if c.cryptoHeaderVersion == CryptoHeaderV1 {
		return returnWithV1Header(c.encryptor.Id(), r)
	}
	return nil, errors.New("unsupported crypto header version")
}

func (c *Crypto) Decrypt(data []byte) ([]byte, error) {
	cryptorId, encryptedData, err := ParseHeader(data)
	if err != nil {
		return nil, err
	}

	if cryptorId == nil {
		return c.defaultDecryptor.Decrypt(*encryptedData)
	}

	return c.decryptors[0].Decrypt(*encryptedData)
}

func (c *Crypto) EncryptStream(input io.Reader, output io.Writer) error {
	_, e := c.encryptor.EncryptStream(input, output)
	if e != nil {
		return e
	}
	if c.encryptor.HeaderVersion() == headless {
		return nil
	}

	if c.cryptoHeaderVersion == CryptoHeaderV1 {
		return nil //TODO it's too late to add header now
	}
	return nil //TODO
}

func (c *Crypto) DecryptStream(input io.Reader, output io.Writer) error {
	cryptorId, reader, metadata, err := ParseHeaderStream(input)
	if err != nil {
		return err
	}

	if cryptorId == nil {
		return c.defaultDecryptor.DecryptStream(reader, metadata, output)
	}

	return c.decryptors[0].DecryptStream(reader, metadata, output)
}
