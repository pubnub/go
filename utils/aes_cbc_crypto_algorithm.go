package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"io"
)

type CBCCryptor struct {
	block cipher.Block
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

func NewCBCCryptor(cipherKey string) (*CBCCryptor, error) {
	block, e := aesCipher(cipherKey)
	if e != nil {
		return nil, e
	}

	return &CBCCryptor{
		block: block,
	}, nil
}

func (c *CBCCryptor) Encrypt(message []byte) (EncryptedData, error) {
	message = padWithPKCS7(message)
	iv := generateIV(aes.BlockSize)
	blockmode := cipher.NewCBCEncrypter(c.block, iv)

	encryptedBytes := make([]byte, len(message))
	blockmode.CryptBlocks(encryptedBytes, message)

	return EncryptedData{
		Metadata: iv,
		Data:     encryptedBytes,
	}, nil
}

func (c *CBCCryptor) Decrypt(encryptedData EncryptedData) (r []byte, e error) {
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

var crivId = []byte{'C', 'R', 'I', 'V'} //("CRIV")

func (c *CBCCryptor) Id() []byte {
	return crivId
}

func (c *CBCCryptor) HeaderVersion() CryptoHeaderVersion {
	return CryptoHeaderV1
}

func (c *CBCCryptor) EncryptStream(input io.Reader, output io.Writer) ([]byte, error) {
	iv := generateIV(aes.BlockSize)
	encrypter := cipher.NewCBCEncrypter(c.block, iv)

	err := encryptStream(encrypter, input, output)
	if err != nil {
		return nil, err
	}
	return iv, nil
}

func (c *CBCCryptor) DecryptStream(input io.Reader, metadata []byte, output io.Writer) error {
	decrypter := cipher.NewCBCDecrypter(c.block, metadata)
	return decryptStream(decrypter, input, output)
}
