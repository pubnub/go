package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
)

type aesCbcCryptor struct {
	block cipher.Block
}

func NewAesCbcCryptor(cipherKey string) (ExtendedCryptor, error) {
	block, e := aesCipher(cipherKey)
	if e != nil {
		return nil, e
	}

	return &aesCbcCryptor{
		block: block,
	}, nil
}

var crivId = "ACRH"

func (c *aesCbcCryptor) Id() string {
	return crivId
}

func (c *aesCbcCryptor) Encrypt(message []byte) (*EncryptedData, error) {
	message = padWithPKCS7(message)
	iv := generateIV(aes.BlockSize)
	blockmode := cipher.NewCBCEncrypter(c.block, iv)

	encryptedBytes := make([]byte, len(message))
	blockmode.CryptBlocks(encryptedBytes, message)

	return &EncryptedData{
		Metadata: iv,
		Data:     encryptedBytes,
	}, nil
}

func (c *aesCbcCryptor) Decrypt(encryptedData *EncryptedData) (r []byte, e error) {
	decrypter := cipher.NewCBCDecrypter(c.block, encryptedData.Metadata)
	//to handle decryption errors
	defer func() {
		if rec := recover(); rec != nil {
			r, e = nil, fmt.Errorf("decrypt error: %s", rec)
		}
	}()

	decrypted := make([]byte, len(encryptedData.Data))
	decrypter.CryptBlocks(decrypted, encryptedData.Data)
	val, err := unpadPKCS7(decrypted)
	if err != nil {
		return nil, fmt.Errorf("decrypt error: %s", err)
	}

	return val, nil
}

func (c *aesCbcCryptor) EncryptStream(reader io.Reader) (*EncryptedStreamData, error) {
	iv := generateIV(aes.BlockSize)

	return &EncryptedStreamData{
		Metadata: iv,
		Reader:   newBlockModeEncryptingReader(reader, cipher.NewCBCEncrypter(c.block, iv)),
	}, nil
}

func (c *aesCbcCryptor) DecryptStream(encryptedData *EncryptedStreamData) (io.Reader, error) {
	if encryptedData.Metadata == nil {
		return nil, errors.New("missing metadata")
	}
	return newBlockModeDecryptingReader(encryptedData.Reader, cipher.NewCBCDecrypter(c.block, encryptedData.Metadata)), nil
}

func aesCipher(cipherKey string) (cipher.Block, error) {
	hash := sha256.New()
	hash.Write([]byte(cipherKey))

	block, err := aes.NewCipher(hash.Sum(nil))
	if err != nil {
		return nil, err
	}
	return block, nil
}
