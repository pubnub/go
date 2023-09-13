package crypto

import "io"

type EncryptedData struct {
	Metadata []byte
	Data     []byte
}

type CryptoAlgorithm interface {
	Id() string
	Encrypt(message []byte) (*EncryptedData, error)
	Decrypt(encryptedData *EncryptedData) ([]byte, error)
}

type EncryptedStreamData struct {
	Metadata []byte
	Reader   io.Reader
}

type ExtendedCryptoAlgorithm interface {
	CryptoAlgorithm
	EncryptStream(reader io.Reader) (*EncryptedStreamData, error)
	DecryptStream(encryptedData *EncryptedStreamData) (io.Reader, error)
}
