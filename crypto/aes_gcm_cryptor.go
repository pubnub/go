package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
	"io"
)

// aesGcmCryptorID is the 4-byte identifier written into the PNED header.
//
// Metadata layout is fixed at 28 bytes: the 12-byte nonce followed by the
// 16-byte authentication tag. The ciphertext is stored after the header and
// is not padded. AES-GCM does not use PKCS#7.
const aesGcmCryptorID = "ACGM"

const (
	gcmNonceSize    = 12
	gcmTagSize      = 16
	gcmMetadataSize = gcmNonceSize + gcmTagSize
)

// aesGcmCryptor decrypts AES-256-GCM payloads. The key is SHA-256(cipherKey),
// which is 32 bytes, and the additional authenticated data is empty.
//
// Factories register this cryptor only in the decryption set. Encrypt is
// implemented so this package can produce ciphertext for tests and for a
// later opt-in encryptor. It is not selected when a module encrypts.
type aesGcmCryptor struct {
	aead cipher.AEAD
}

func newAesGcmCryptor(cipherKey string) (ExtendedCryptor, error) {
	sum := sha256.Sum256([]byte(cipherKey))
	block, e := aes.NewCipher(sum[:])
	if e != nil {
		return nil, e
	}

	aead, e := cipher.NewGCM(block)
	if e != nil {
		return nil, e
	}
	if aead.NonceSize() != gcmNonceSize || aead.Overhead() != gcmTagSize {
		return nil, errors.New("encryption error")
	}

	return &aesGcmCryptor{aead: aead}, nil
}

func (c *aesGcmCryptor) Id() string {
	return aesGcmCryptorID
}

func (c *aesGcmCryptor) Encrypt(message []byte) (*EncryptedData, error) {
	return c.encryptWithNonce(message, generateIV(gcmNonceSize))
}

func (c *aesGcmCryptor) encryptWithNonce(message, nonce []byte) (*EncryptedData, error) {
	if len(message) == 0 {
		return nil, errors.New("encryption error: can't encrypt empty data")
	}
	if len(nonce) != gcmNonceSize {
		return nil, errors.New("encryption error")
	}

	// Seal appends the authentication tag after the ciphertext. The wire
	// format stores that tag in metadata, next to the nonce, and keeps only
	// the ciphertext in the body.
	sealed := c.aead.Seal(nil, nonce, message, nil)
	ciphertextLen := len(sealed) - gcmTagSize

	metadata := make([]byte, gcmMetadataSize)
	copy(metadata, nonce)
	copy(metadata[gcmNonceSize:], sealed[ciphertextLen:])

	ciphertext := make([]byte, ciphertextLen)
	copy(ciphertext, sealed[:ciphertextLen])

	return &EncryptedData{
		Metadata: metadata,
		Data:     ciphertext,
	}, nil
}

func (c *aesGcmCryptor) Decrypt(encryptedData *EncryptedData) (r []byte, e error) {
	defer func() {
		if rec := recover(); rec != nil {
			r, e = nil, errors.New("decryption error")
		}
	}()

	if encryptedData == nil || len(encryptedData.Metadata) != gcmMetadataSize {
		return nil, errors.New("decryption error")
	}

	nonce := encryptedData.Metadata[:gcmNonceSize]
	tag := encryptedData.Metadata[gcmNonceSize:]
	sealed := make([]byte, len(encryptedData.Data)+gcmTagSize)
	copy(sealed, encryptedData.Data)
	copy(sealed[len(encryptedData.Data):], tag)

	plain, e := c.aead.Open(nil, nonce, sealed, nil)
	if e != nil {
		return nil, errors.New("decryption error")
	}
	return plain, nil
}

func (c *aesGcmCryptor) EncryptStream(reader io.Reader) (*EncryptedStreamData, error) {
	message, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}

	encrypted, e := c.Encrypt(message)
	if e != nil {
		return nil, e
	}

	return &EncryptedStreamData{
		Metadata: encrypted.Metadata,
		Reader:   bytes.NewReader(encrypted.Data),
	}, nil
}

func (c *aesGcmCryptor) DecryptStream(encryptedData *EncryptedStreamData) (io.Reader, error) {
	if encryptedData == nil {
		return nil, errors.New("decryption error")
	}

	ciphertext, e := io.ReadAll(encryptedData.Reader)
	if e != nil {
		return nil, e
	}

	plain, e := c.Decrypt(&EncryptedData{
		Metadata: encryptedData.Metadata,
		Data:     ciphertext,
	})
	if e != nil {
		return nil, e
	}

	return bytes.NewReader(plain), nil
}
