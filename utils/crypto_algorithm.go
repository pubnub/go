package utils

import (
	"bytes"
	"io"
)

type CryptoAlgorithm interface {
	HeaderVersion() CryptoHeaderVersion
	Id() []byte
	Encrypt(message []byte) (EncryptedData, error)
	Decrypt(encryptedData EncryptedData) ([]byte, error)

	EncryptStream(input io.Reader, output io.Writer) ([]byte, error)
	DecryptStream(input io.Reader, metadata []byte, output io.Writer) error
}

type EncryptedData struct {
	Metadata []byte
	Data     []byte
}

type EncryptingWriter interface {
	io.Writer
	WriteMetadata(in []byte) (n int, err error)
}

type EncryptingWriterV1 struct {
	w              *io.Writer
	metadataBuffer *bytes.Buffer
	metadataDone   bool
}

func (e *EncryptingWriterV1) Write(in []byte) (n int, err error) {
	e.metadataDone = true
	return 0, nil
}

func (e *EncryptingWriterV1) WriteMetadata(in []byte) error {
	_, er := e.metadataBuffer.Write(in)
	if er != nil {
		return er
	}
	return nil
}
