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
// Metadata is the 12-byte nonce. The 16-byte authentication tag is appended
// to the ciphertext. AES-GCM does not use PKCS#7.
const aesGcmCryptorID = "ACGM"

const (
	gcmNonceSize = 12
	gcmTagSize   = 16
)

// aesGcmCryptor decrypts AES-256-GCM payloads. The key is SHA-256(cipherKey),
// which is 32 bytes, and the additional authenticated data is empty.
//
// Factories register this cryptor only in the decryption set. Encrypt returns
// an error until AES-GCM encryption is enabled.
type aesGcmCryptor struct {
	aead cipher.AEAD
}

// errAesGcmEncryptionUnavailable is returned until AES-GCM encryption is enabled.
var errAesGcmEncryptionUnavailable = errors.New("encryption error: AES-GCM encryption is not available")

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
	// Temporary: this cryptor is decryption-only. Remove this return when
	// AES-GCM encryption is enabled.
	return nil, errAesGcmEncryptionUnavailable
}

func (c *aesGcmCryptor) encryptWithNonce(message, nonce []byte) (*EncryptedData, error) {
	if len(message) == 0 {
		return nil, errors.New("encryption error: can't encrypt empty data")
	}
	if len(nonce) != gcmNonceSize {
		return nil, errors.New("encryption error")
	}

	// Seal appends the authentication tag after the ciphertext. That tag stays
	// on the body. Metadata carries only the nonce.
	sealed := c.aead.Seal(nil, nonce, message, nil)
	metadata := make([]byte, gcmNonceSize)
	copy(metadata, nonce)

	return &EncryptedData{
		Metadata: metadata,
		Data:     sealed,
	}, nil
}

func (c *aesGcmCryptor) Decrypt(encryptedData *EncryptedData) (r []byte, e error) {
	defer func() {
		if rec := recover(); rec != nil {
			r, e = nil, errors.New("decryption error")
		}
	}()

	if encryptedData == nil || len(encryptedData.Metadata) != gcmNonceSize || len(encryptedData.Data) < gcmTagSize {
		return nil, errors.New("decryption error")
	}

	plain, e := c.aead.Open(nil, encryptedData.Metadata, encryptedData.Data, nil)
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
