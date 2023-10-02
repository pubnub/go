package crypto

import (
	"bufio"
	"bytes"
	"errors"
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

func NewAesCbcCryptoModule(cipherKey string, randomIv bool) (CryptoModule, error) {
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
	if len(message) == 0 {
		return nil, errors.New("encryption error: can't encrypt empty data")
	}
	encryptedData, e := m.encryptor.Encrypt(message)
	if e != nil {
		return nil, fmt.Errorf("encryption error: %s", e.Error())
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
	if len(data) == 0 {
		return nil, errors.New("decryption error: can't decrypt empty data")
	}

	id, encryptedData, e := parseHeader(data)
	if e != nil {
		return nil, e
	}

	decryptor := m.decryptors[*id]

	if decryptor == nil {
		return nil, fmt.Errorf("unknown crypto error: unknown cryptor id %s", *id)
	}

	if len(encryptedData.Data) == 0 {
		return nil, errors.New("decryption error: can't decrypt empty data")
	}

	var r []byte
	if r, e = m.decryptors[*id].Decrypt(encryptedData); e != nil {
		return nil, fmt.Errorf("decryption error: %s", e.Error())
	}

	return r, nil
}

func (m *module) EncryptStream(input io.Reader) (io.Reader, error) {
	bufferedReader := bufio.NewReader(input)
	peekedBytes, e := bufferedReader.Peek(1)
	if len(peekedBytes) == 0 {
		return nil, errors.New("encryption error: can't encrypt empty data")
	}

	encryptedStreamData, e := m.encryptor.EncryptStream(bufferedReader)
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
	bufData := bufio.NewReader(input)

	id, encryptedStreamData, e := parseHeaderStream(bufData)
	if e != nil {
		return nil, e
	}

	peeked, e := bufData.Peek(1)
	if e != nil {
		return nil, fmt.Errorf("decryption error: %w", e)
	}

	if len(peeked) == 0 {
		return nil, errors.New("decryption error: can't decrypt empty data")
	}

	decryptor := m.decryptors[*id]

	if decryptor == nil {
		return nil, fmt.Errorf("unknown crypto error: unknown cryptor id %s", *id)
	}

	if e != nil {
		return nil, fmt.Errorf("decryption error: %s", e.Error())
	}

	var r io.Reader
	if r, e = m.decryptors[*id].DecryptStream(encryptedStreamData); e != nil {
		return nil, fmt.Errorf("decryption error: %s", e.Error())
	}

	return r, nil
}
