package crypto

import (
	"bufio"
	"io"
)

type Module struct {
	encryptor         Cryptor
	decryptors        map[string]Cryptor
	fallbackDecryptor Cryptor
}

func NewLegacyCryptoModule(cipherKey string, randomIv bool) (*Module, error) {
	legacyCryptoAlgorithm, e := NewLegacyCryptoAlgorithm(cipherKey, randomIv)
	if e != nil {
		return nil, e
	}
	aesCBCCryptoAlgorithm, e := NewAesCBCCryptoAlgorithm(cipherKey)

	if e != nil {
		return nil, e
	}

	decryptors := []CryptoAlgorithm{aesCBCCryptoAlgorithm}
	return newCryptoModule(legacyCryptoAlgorithm, decryptors, legacyCryptoAlgorithm), nil
}

func NewCryptoModule(cipherKey string, randomIv bool) (*Module, error) {
	aesCBCCryptoAlgorithm, e := NewAesCBCCryptoAlgorithm(cipherKey)
	if e != nil {
		return nil, e
	}

	legacyCryptoAlgorithm, e := NewLegacyCryptoAlgorithm(cipherKey, randomIv)
	if e != nil {
		return nil, e
	}

	decryptors := []CryptoAlgorithm{aesCBCCryptoAlgorithm}
	return newCryptoModule(aesCBCCryptoAlgorithm, decryptors, legacyCryptoAlgorithm), nil
}

func newCryptoModule(encryptingAlgorithm CryptoAlgorithm, decryptingAlgorithms []CryptoAlgorithm, fallbackDecryptingAlgorithm CryptoAlgorithm) *Module {

	decryptors := make(map[string]Cryptor, len(decryptingAlgorithms))
	for _, decryptor := range decryptingAlgorithms {
		decryptors[decryptor.Id()] = NewCryptor(decryptor)
	}

	encryptor := NewCryptor(encryptingAlgorithm)
	fallbackDecryptor := NewCryptor(fallbackDecryptingAlgorithm)

	return &Module{
		encryptor:         encryptor,
		decryptors:        decryptors,
		fallbackDecryptor: fallbackDecryptor,
	}
}

func (c *Module) Encrypt(message []byte) ([]byte, error) {
	return c.encryptor.Encrypt(message)
}

func (c *Module) Decrypt(data []byte) ([]byte, error) {
	cryptorId, e := peekHeaderCryptorId(data)
	if e != nil {
		return nil, e
	}
	if cryptorId == nil {
		return c.fallbackDecryptor.Decrypt(data)
	}
	return c.decryptors[*cryptorId].Decrypt(data)
}

func (c *Module) EncryptStream(input io.Reader) (io.Reader, error) {
	return c.encryptor.EncryptStream(input)
}

func (c *Module) DecryptStream(input io.Reader) (io.Reader, error) {
	data := bufio.NewReader(input)
	cryptorId, e := peekHeaderStreamCryptorId(data)
	if e != nil {
		return nil, e
	}
	if cryptorId == nil {
		return c.fallbackDecryptor.DecryptStream(data)
	}
	return c.decryptors[*cryptorId].DecryptStream(data)
}
