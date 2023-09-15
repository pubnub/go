package crypto

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

// CryptoModule is an interface for encrypting and decrypting data.
type CryptoModule interface {
	Encrypt(input []byte) ([]byte, error)
	Decrypt(input []byte) ([]byte, error)
	EncryptStream(input io.Reader) (io.Reader, error)
	DecryptStream(input io.Reader) (io.Reader, error)
}

type module struct {
	encryptor  ExtendedCryptor
	decryptors map[string]ExtendedCryptor
}

func NewLegacyCryptoModule(cipherKey string, randomIv bool) (CryptoModule, error) {
	legacy, e := NewLegacyCryptor(cipherKey, randomIv)
	if e != nil {
		return nil, e
	}
	aesCbc, e := NewAesCbcCryptor(cipherKey)

	if e != nil {
		return nil, e
	}

	return NewCryptoModule(legacy, []Cryptor{aesCbc}), nil
}

func NewDefaultCryptoModule(cipherKey string, randomIv bool) (CryptoModule, error) {
	aesCbc, e := NewAesCbcCryptor(cipherKey)
	if e != nil {
		return nil, e
	}

	legacy, e := NewLegacyCryptor(cipherKey, randomIv)
	if e != nil {
		return nil, e
	}

	return NewCryptoModule(aesCbc, []Cryptor{legacy}), nil
}

func NewCryptoModule(defaultCryptor Cryptor, decryptors []Cryptor) CryptoModule {

	decryptorsMap := make(map[string]ExtendedCryptor, len(decryptors)+1)
	for _, d := range decryptors {
		decryptorsMap[d.Id()] = liftToExtendedCryptor(d)
	}

	encryptor := liftToExtendedCryptor(defaultCryptor)
	decryptorsMap[encryptor.Id()] = encryptor

	return &module{
		encryptor:  encryptor,
		decryptors: decryptorsMap,
	}
}

func (m *module) Encrypt(message []byte) ([]byte, error) {
	encryptedData, e := m.encryptor.Encrypt(message)
	if e != nil {
		return nil, e
	}

	if m.encryptor.Id() == legacyId {
		return encryptedData.Data, nil
	}

	header, e := headerV1(m.encryptor.Id(), encryptedData.Metadata)
	if e != nil {
		return nil, e
	}

	return append(header, encryptedData.Data...), nil
}

func (m *module) Decrypt(data []byte) ([]byte, error) {
	id, encryptedData, e := parseHeader(data)
	if e != nil {
		return nil, e
	}

	decryptor := m.decryptors[*id]

	if decryptor == nil {
		return nil, fmt.Errorf("decryption error: unknown cryptor id %s", *id)
	}

	return m.decryptors[*id].Decrypt(encryptedData)
}

func (m *module) EncryptStream(input io.Reader) (io.Reader, error) {
	encryptedStreamData, e := m.encryptor.EncryptStream(input)
	if e != nil {
		return nil, e
	}
	if m.encryptor.Id() == legacyId {
		return encryptedStreamData.Reader, nil
	}
	header, e := headerV1(m.encryptor.Id(), encryptedStreamData.Metadata)
	if e != nil {
		return nil, e
	}

	headerReader := bytes.NewReader(header)

	return io.MultiReader(headerReader, encryptedStreamData.Reader), nil
}

func (m *module) DecryptStream(input io.Reader) (io.Reader, error) {

	id, encryptedStreamData, e := parseHeaderStream(bufio.NewReader(input))
	if e != nil {
		return nil, e
	}
	decryptor := m.decryptors[*id]

	if decryptor == nil {
		return nil, fmt.Errorf("decryption error: unknown cryptor id %s", *id)
	}

	return m.decryptors[*id].DecryptStream(encryptedStreamData)
}
