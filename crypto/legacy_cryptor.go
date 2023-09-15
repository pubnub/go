package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

// 16 byte constant legacy IV
var valIV = "0123456789012345"

type legacyCryptor struct {
	block    cipher.Block
	randomIv bool
}

func NewLegacyCryptor(cipherKey string, useRandomIV bool) (ExtendedCryptor, error) {
	block, e := legacyAesCipher(cipherKey)
	if e != nil {
		return nil, e
	}

	return &legacyCryptor{
		block:    block,
		randomIv: useRandomIV,
	}, nil
}

var legacyId = string([]byte{0x00, 0x00, 0x00, 0x00})

func (c *legacyCryptor) Id() string {
	return legacyId
}

func (c *legacyCryptor) Encrypt(message []byte) (*EncryptedData, error) {
	message = padWithPKCS7(message)
	iv := make([]byte, aes.BlockSize)
	if c.randomIv {
		iv = generateIV(aes.BlockSize)
	} else {
		iv = []byte(valIV)
	}
	blockmode := cipher.NewCBCEncrypter(c.block, iv)

	encryptedBytes := make([]byte, len(message))
	blockmode.CryptBlocks(encryptedBytes, message)

	if c.randomIv {
		return &EncryptedData{Data: append(iv, encryptedBytes...), Metadata: nil}, nil
	}
	return &EncryptedData{Data: encryptedBytes, Metadata: nil}, nil
}

func (c *legacyCryptor) Decrypt(encryptedData *EncryptedData) (r []byte, e error) {
	iv := make([]byte, aes.BlockSize)
	data := encryptedData.Data
	if c.randomIv {
		iv = data[:aes.BlockSize]
		data = data[aes.BlockSize:]
	} else {
		iv = []byte(valIV)
	}

	decrypter := cipher.NewCBCDecrypter(c.block, iv)
	//to handle decryption errors
	defer func() {
		if rec := recover(); rec != nil {
			r, e = nil, fmt.Errorf("decrypt error: %s", rec)
		}
	}()

	decrypted := make([]byte, len(data))
	decrypter.CryptBlocks(decrypted, data)
	val, err := unpadPKCS7(decrypted)
	if err != nil {
		return nil, fmt.Errorf("decrypt error: %s", err)
	}

	return val, nil
}

func (c *legacyCryptor) EncryptStream(reader io.Reader) (*EncryptedStreamData, error) {
	iv := generateIV(aes.BlockSize)

	return &EncryptedStreamData{
		Metadata: nil,
		Reader:   io.MultiReader(bytes.NewReader(iv), NewBlockModeEncryptingReader(reader, cipher.NewCBCEncrypter(c.block, iv))),
	}, nil
}

func (c *legacyCryptor) DecryptStream(encryptedData *EncryptedStreamData) (io.Reader, error) {
	iv := make([]byte, aes.BlockSize)
	_, err := io.ReadFull(encryptedData.Reader, iv)
	if err != nil {
		return nil, err
	}

	return NewBlockModeDecryptingReader(encryptedData.Reader, cipher.NewCBCDecrypter(c.block, iv)), nil
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
