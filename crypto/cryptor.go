package crypto

import (
	"bytes"
	"io"
)

type Cryptor interface {
	Id() string
	Encrypt(input []byte) ([]byte, error)
	Decrypt(input []byte) ([]byte, error)
	EncryptStream(input io.Reader) (io.Reader, error)
	DecryptStream(input io.Reader) (io.Reader, error)
}

func NewCryptor(algorithm CryptoAlgorithm) Cryptor {
	var c Cryptor
	if legacy, ok := algorithm.(*LegacyCryptoAlgorithm); ok {
		c = newLegacyCryptor(legacy)
	} else if extended, ok := algorithm.(ExtendedCryptoAlgorithm); ok {
		c = newExtendedCryptorV1(extended)
	} else {
		c = newCryptorV1(algorithm.(CryptoAlgorithm))
	}
	return c
}

func newExtendedCryptorV1(algorithm ExtendedCryptoAlgorithm) Cryptor {
	var c Cryptor
	c = &extendedCryptorV1{
		cryptorV1: cryptorV1{
			algorithm: algorithm,
		},
		algorithm: algorithm,
	}
	return c
}

func newCryptorV1(algorithm CryptoAlgorithm) Cryptor {
	var c Cryptor
	c = &cryptorV1{
		algorithm: algorithm,
	}
	return c
}

func newLegacyCryptor(algorithm ExtendedCryptoAlgorithm) Cryptor {
	var c Cryptor
	c = &legacyCryptor{
		algorithm: algorithm,
	}
	return c
}

type extendedCryptorV1 struct {
	cryptorV1
	algorithm ExtendedCryptoAlgorithm
}

//func (c *extendedCryptorV1) Encrypt(message []byte) ([]byte, error) {
//	return c.cryptorV1.Encrypt(message)
//}
//
//func (c *extendedCryptorV1) Decrypt(message []byte) ([]byte, error) {
//	return c.cryptorV1.Decrypt(message)
//}

func (c *extendedCryptorV1) EncryptStream(input io.Reader) (io.Reader, error) {
	encryptedStreamData, e := c.algorithm.EncryptStream(input)
	if e != nil {
		return nil, e
	}
	header, e := headerV1(c.algorithm.Id(), encryptedStreamData.Metadata)
	if e != nil {
		return nil, e
	}

	headerReader := bytes.NewReader(header)

	return io.MultiReader(headerReader, encryptedStreamData.Reader), nil
}

func (c *extendedCryptorV1) DecryptStream(input io.Reader) (io.Reader, error) {
	_, readerWithOnlyData, metadata, e := parseHeaderStream(input)
	if e != nil {
		return nil, e
	}
	return c.algorithm.DecryptStream(&EncryptedStreamData{Reader: readerWithOnlyData, Metadata: metadata})
}

type cryptorV1 struct {
	algorithm CryptoAlgorithm
}

func (c *cryptorV1) Id() string {
	return c.algorithm.Id()
}

func (c *cryptorV1) Encrypt(message []byte) ([]byte, error) {
	encryptedData, e := c.algorithm.Encrypt(message)
	if e != nil {
		return nil, e
	}
	header, e := headerV1(c.algorithm.Id(), encryptedData.Metadata)
	if e != nil {
		return nil, e
	}

	return append(header, encryptedData.Data...), nil
}

func (c *cryptorV1) Decrypt(message []byte) ([]byte, error) {
	_, encryptedData, e := parseHeader(message)
	if e != nil {
		return nil, e
	}
	return c.algorithm.Decrypt(encryptedData)
}

func (c *cryptorV1) EncryptStream(input io.Reader) (io.Reader, error) {
	inputBytes, e := io.ReadAll(input)
	if e != nil {
		return nil, e
	}

	encryptedData, e := c.algorithm.Encrypt(inputBytes)
	if e != nil {
		return nil, e
	}

	header, e := headerV1(c.algorithm.Id(), encryptedData.Metadata)
	if e != nil {
		return nil, e
	}

	return io.MultiReader(bytes.NewReader(header), bytes.NewReader(encryptedData.Data)), nil
}

func (c *cryptorV1) DecryptStream(input io.Reader) (io.Reader, error) {
	inputBytes, e := io.ReadAll(input)
	if e != nil {
		return nil, e
	}

	_, encryptedData, e := parseHeader(inputBytes)
	if e != nil {
		return nil, e
	}

	decryptedData, e := c.algorithm.Decrypt(encryptedData)
	if e != nil {
		return nil, e
	}

	return bytes.NewReader(decryptedData), nil
}

type legacyCryptor struct {
	algorithm ExtendedCryptoAlgorithm
}

func (c *legacyCryptor) Id() string {
	return c.algorithm.Id()
}

func (c *legacyCryptor) Encrypt(message []byte) ([]byte, error) {
	encryptedData, e := c.algorithm.Encrypt(message)
	if e != nil {
		return nil, e
	}

	return encryptedData.Data, nil
}

func (c *legacyCryptor) Decrypt(message []byte) ([]byte, error) {
	return c.algorithm.Decrypt(&EncryptedData{Data: message, Metadata: nil})
}

func (c *legacyCryptor) EncryptStream(input io.Reader) (io.Reader, error) {
	encryptedStreamData, e := c.algorithm.EncryptStream(input)
	if e != nil {
		return nil, e
	}

	return encryptedStreamData.Reader, nil
}

func (c *legacyCryptor) DecryptStream(input io.Reader) (io.Reader, error) {
	return c.algorithm.DecryptStream(&EncryptedStreamData{Reader: input, Metadata: nil})
}
