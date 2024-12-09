package crypto

import "io"

type EncryptedData struct {
	Metadata []byte
	Data     []byte
}

type Cryptor interface {
	Id() string
	Encrypt(message []byte) (*EncryptedData, error)
	Decrypt(encryptedData *EncryptedData) ([]byte, error)
}

type EncryptedStreamData struct {
	Metadata []byte
	Reader   io.Reader
}

type ExtendedCryptor interface {
	Cryptor
	EncryptStream(reader io.Reader) (*EncryptedStreamData, error)
	DecryptStream(encryptedData *EncryptedStreamData) (io.Reader, error)
}
