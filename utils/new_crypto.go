package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
)

type Crypto struct {
	encryptor           Cryptor
	decryptors          []Cryptor
	defaultDecryptor    Cryptor
	cryptoHeaderVersion CryptoHeaderVersion
}

func NewCrypto(encryptor Cryptor, decryptors []Cryptor, defaultDecryptor Cryptor, version CryptoHeaderVersion) *Crypto {
	return &Crypto{
		encryptor:           encryptor,
		decryptors:          decryptors,
		defaultDecryptor:    defaultDecryptor,
		cryptoHeaderVersion: version,
	}
}

func (c *Crypto) Encrypt(message []byte) ([]byte, error) {
	if c.cryptoHeaderVersion == Headless {
		r, e := c.encryptor.Encrypt(message)
		if e != nil {
			return nil, e
		}
		return r.Data, nil
	}
	if c.cryptoHeaderVersion == CryptoHeaderV1 {
		r, e := c.encryptor.Encrypt(message)
		if e != nil {
			return nil, e
		}
		return returnWithV1Header(c.encryptor.Id(), r)
	}
	return nil, errors.New("unsupported crypto header version")
}

func (c *Crypto) Decrypt(data []byte) ([]byte, error) {
	cryptorId, encryptedData, err := parseHeader(data)
	if err != nil {
		return nil, err
	}

	if cryptorId == nil {
		return c.defaultDecryptor.Decrypt(*encryptedData)
	}

	return c.decryptors[0].Decrypt(*encryptedData)

}

type Cryptor interface {
	Id() []byte
	Encrypt(message []byte) (EncryptedData, error)
	Decrypt(encryptedData EncryptedData) ([]byte, error)
}

type LegacyCryptor struct {
	cipherKey   cipher.Block
	useRandomIV bool
}

func NewLegacyCryptor(cipherKey string, useRandomIV bool) (*LegacyCryptor, error) {
	block, e := legacyAesCipher(cipherKey)
	if e != nil {
		return nil, e
	}

	return &LegacyCryptor{
		cipherKey:   block,
		useRandomIV: useRandomIV,
	}, nil
}

func (c *LegacyCryptor) Id() []byte {
	return []byte("lega")
}

func (c *LegacyCryptor) Encrypt(message []byte) (EncryptedData, error) {
	d := encrypt(c.cipherKey, message, c.useRandomIV)

	return EncryptedData{
		CryptorData: nil,
		Data:        d,
	}, nil
}

func (c *LegacyCryptor) Decrypt(encryptedData EncryptedData) ([]byte, error) {
	return decrypt(c.cipherKey, encryptedData.Data, c.useRandomIV)
}

type EncryptedData struct {
	CryptorData []byte
	Data        []byte
}

type CBCCryptor struct {
	cipherKey cipher.Block
}

func aesCipher(cipherKey string) (cipher.Block, error) {
	block, err := aes.NewCipher([]byte(cipherKey))
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
		cipherKey: block,
	}, nil
}

func (c *CBCCryptor) Encrypt(message []byte) (EncryptedData, error) {
	message = padWithPKCS7(message)
	iv := generateIV(aes.BlockSize)
	blockmode := cipher.NewCBCEncrypter(c.cipherKey, iv)

	encryptedBytes := make([]byte, len(message))
	blockmode.CryptBlocks(encryptedBytes, message)

	return EncryptedData{
		CryptorData: iv,
		Data:        encryptedBytes,
	}, nil
}

func (c *CBCCryptor) Decrypt(encryptedData EncryptedData) (r []byte, e error) {
	decrypter := cipher.NewCBCDecrypter(c.cipherKey, encryptedData.CryptorData)
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

func (c *CBCCryptor) Id() []byte {
	return []byte("CRIV")
}
