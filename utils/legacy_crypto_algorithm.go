package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"io"
)

type LegacyCryptoAlgorithm struct {
	block                cipher.Block
	useRandomIvForSlices bool
}

func NewLegacyCryptor(cipherKey string, useRandomIV bool) (*LegacyCryptoAlgorithm, error) {
	block, e := legacyAesCipher(cipherKey)
	if e != nil {
		return nil, e
	}

	return &LegacyCryptoAlgorithm{
		block:                block,
		useRandomIvForSlices: useRandomIV,
	}, nil
}

var legacyId = []byte{0x00, 0x00, 0x00, 0x00}

func (c *LegacyCryptoAlgorithm) Id() []byte {
	return legacyId
}

func (c *LegacyCryptoAlgorithm) HeaderVersion() CryptoHeaderVersion {
	return headless
}

func (c *LegacyCryptoAlgorithm) Encrypt(value []byte) (EncryptedData, error) {
	value = padWithPKCS7(value)
	iv := make([]byte, aes.BlockSize)
	if c.useRandomIvForSlices {
		iv = generateIV(aes.BlockSize)
	} else {
		iv = []byte(valIV)
	}

	blockmode := cipher.NewCBCEncrypter(c.block, iv)
	cipherBytes := encrypt(blockmode, value)

	if c.useRandomIvForSlices {
		return EncryptedData{Data: append(iv, cipherBytes...), Metadata: nil}, nil
	}
	return EncryptedData{Data: cipherBytes, Metadata: iv}, nil
}

func (c *LegacyCryptoAlgorithm) Decrypt(encryptedData EncryptedData) ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	if c.useRandomIvForSlices {
		iv = generateIV(aes.BlockSize) //TODO
	} else {
		iv = []byte(valIV)
	}

	decrypter := cipher.NewCBCDecrypter(c.block, iv)

	return decrypt(decrypter, encryptedData.Data)
}

func (c *LegacyCryptoAlgorithm) EncryptStream(input io.Reader, output io.Writer) ([]byte, error) {
	iv := generateIV(aes.BlockSize)
	_, err := output.Write(iv)
	if err != nil {
		return nil, err
	}
	blockmode := cipher.NewCBCEncrypter(c.block, iv)
	err = encryptStream(blockmode, input, output)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *LegacyCryptoAlgorithm) DecryptStream(input io.Reader, _ []byte, output io.Writer) error {
	iv := make([]byte, aes.BlockSize)
	_, err := io.ReadFull(input, iv)
	if err != nil {
		return nil
	}

	blockmode := cipher.NewCBCEncrypter(c.block, iv)
	return decryptStream(blockmode, input, output)
}

// EncryptCipherKey DEPRECATED
// EncryptCipherKey creates the 256 bit hex of the cipher key
//
// It accepts the following parameters:
// cipherKey: cipher key to use to decrypt.
//
// returns the 256 bit hex of the cipher key.
func EncryptCipherKey(cipherKey string) []byte {
	hash := sha256.New()
	hash.Write([]byte(cipherKey))

	sha256String := hash.Sum(nil)[:16]
	return []byte(hex.EncodeToString(sha256String))
}

// legacyAesCipher returns the cipher block
//
// It accepts the following parameters:
// cipherKey: cipher key.
//
// returns the cipher block,
// error if any.
func legacyAesCipher(cipherKey string) (cipher.Block, error) {
	key := EncryptCipherKey(cipherKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return block, nil
}
