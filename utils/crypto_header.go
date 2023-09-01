package utils

import (
	"bytes"
	"errors"
	"strconv"
)

type CryptoHeaderVersion int

const (
	Headless CryptoHeaderVersion = iota
	CryptoHeaderV1
)

const version_position = 4
const version_v1 = 1
const cryptor_id_position = 5
const cryptor_id_length = 4
const size_position = 9
const short_size_length = 1
const long_size_length = 3
const long_size_indicator = 0xFF
const max_short_size = 254

var sentinel = []byte{0x50, 0x4E, 0x45, 0x44}

func returnWithV1Header(cryptorId []byte, encryptedData EncryptedData) ([]byte, error) {
	cryptorDataSize := len(encryptedData.CryptorData)
	var cryptorDataBytesSize int

	if cryptorDataSize <= max_short_size {
		cryptorDataBytesSize = short_size_length
	} else {
		cryptorDataBytesSize = long_size_length
	}
	r := make([]byte, 0, len(sentinel)+1+cryptor_id_length+cryptorDataBytesSize+cryptorDataSize+len(encryptedData.Data))

	buffer := bytes.NewBuffer(r)

	buffer.Write(sentinel)
	buffer.WriteByte(version_v1)
	buffer.Write(cryptorId)
	if cryptorDataBytesSize == short_size_length {
		buffer.WriteByte(byte(cryptorDataSize))
	} else {
		buffer.WriteByte(long_size_indicator)
		buffer.Write([]byte(strconv.FormatInt(int64(cryptorDataSize), 10)))
	}
	buffer.Write(encryptedData.CryptorData)
	buffer.Write(encryptedData.Data)

	return buffer.Bytes(), nil
}

func slicesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func parseHeader(data []byte) ([]byte, *EncryptedData, error) {
	if !slicesEqual(data[:len(sentinel)], sentinel) {
		return nil, &EncryptedData{CryptorData: nil, Data: data}, nil
	}

	if data[version_position] != version_v1 {
		return nil, nil, errors.New("unsupported crypto header version")
	}

	cryptorId := data[cryptor_id_position : cryptor_id_position+cryptor_id_length]
	var headerSize int64
	var err error
	var offset int64
	if data[size_position] == long_size_indicator {
		offset = size_position + long_size_length
		headerSize, err = strconv.ParseInt(string(data[size_position:size_position+long_size_length]), 10, 32)
		if err != nil {
			return nil, nil, err
		}
	} else {
		offset = size_position + short_size_length
		headerSize = int64(data[size_position])
	}

	metadata := data[offset : offset+headerSize]
	offset += int64(len(metadata))

	return cryptorId, &EncryptedData{Data: data[offset:], CryptorData: metadata}, nil
}
