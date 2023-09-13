package crypto

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strconv"
)

type CryptoHeaderVersion int

const (
	headless CryptoHeaderVersion = iota
	CryptoHeaderV1
)

const versionPosition = 4
const versionV1 = 1
const cryptorIdPosition = 5
const cryptorIdLength = 4
const sizePosition = 9
const shortSizeLength = 1
const longSizeLength = 3
const longSizeIndicator = 0xFF
const maxShortSize = 254
const sentinelLength = 4

var sentinel = [sentinelLength]byte{0x50, 0x4E, 0x45, 0x44}

func headerV1(cryptorId string, metadata []byte) ([]byte, error) {
	cryptorDataSize := len(metadata)
	var cryptorDataBytesSize int

	if cryptorDataSize <= maxShortSize {
		cryptorDataBytesSize = shortSizeLength
	} else {
		cryptorDataBytesSize = longSizeLength
	}
	r := make([]byte, 0, len(sentinel)+1+cryptorIdLength+cryptorDataBytesSize+cryptorDataSize)

	buffer := bytes.NewBuffer(r)
	_, e := buffer.Write(sentinel[:])
	if e != nil {
		return nil, e
	}
	e = buffer.WriteByte(versionV1)
	if e != nil {
		return nil, e
	}

	_, e = buffer.Write([]byte(cryptorId))
	if e != nil {
		return nil, e
	}
	if cryptorDataBytesSize == shortSizeLength {
		e = buffer.WriteByte(byte(cryptorDataSize))
		if e != nil {
			return nil, e
		}
	} else {
		e = buffer.WriteByte(longSizeIndicator)
		if e != nil {
			return nil, e
		}
		_, e = buffer.Write([]byte(strconv.FormatInt(int64(cryptorDataSize), 10)))
		if e != nil {
			return nil, e
		}
	}
	_, e = buffer.Write(metadata)
	if e != nil {
		return nil, e
	}
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

func peekHeaderCryptorId(data []byte) (cryptorId *string, e error) {
	if !slicesEqual(data[:len(sentinel)], sentinel[:]) {
		return nil, nil
	}

	if data[versionPosition] != versionV1 {
		return nil, errors.New("unsupported crypto header version")
	}

	id := string(data[cryptorIdPosition : cryptorIdPosition+cryptorIdLength])
	return &id, nil
}

func parseHeader(data []byte) (cryptorId *string, encryptedData *EncryptedData, e error) {
	if !slicesEqual(data[:len(sentinel)], sentinel[:]) {
		return nil, &EncryptedData{Metadata: nil, Data: data}, nil
	}

	if data[versionPosition] != versionV1 {
		return nil, nil, errors.New("unsupported crypto header version")
	}

	id := string(data[cryptorIdPosition : cryptorIdPosition+cryptorIdLength])
	var headerSize int64
	position := int64(sizePosition)
	if data[sizePosition] == longSizeIndicator {
		position += longSizeLength
		var err error
		headerSize, err = strconv.ParseInt(string(data[sizePosition:sizePosition+longSizeLength]), 10, 32)
		if err != nil {
			return nil, nil, err
		}
	} else {
		position += shortSizeLength
		headerSize = int64(data[sizePosition])
	}

	metadata := data[position : position+headerSize]
	position += int64(len(metadata))

	return &id, &EncryptedData{Data: data[position:], Metadata: metadata}, nil
}

func peekHeaderStreamCryptorId(data *bufio.Reader) (cryptorId *string, e error) {
	peeked, err := data.Peek(sentinelLength + 1 + cryptorIdLength + longSizeLength)
	if err != nil {
		return nil, err
	}

	if !slicesEqual(peeked[:len(sentinel)], sentinel[:]) {
		return nil, nil
	}

	if peeked[versionPosition] != versionV1 {
		return nil, errors.New("unsupported crypto header version")
	}

	id := string(peeked[cryptorIdPosition : cryptorIdPosition+cryptorIdLength])
	return &id, nil
}

func parseHeaderStream(data io.Reader) (cryptorId []byte, d io.Reader, metadata []byte, e error) {
	bufData := bufio.NewReader(data)

	peeked, err := bufData.Peek(sentinelLength + 1 + cryptorIdLength + longSizeLength)
	if err != nil {
		return nil, nil, nil, err
	}

	if !slicesEqual(peeked[:len(sentinel)], sentinel[:]) {
		return nil, bufData, nil, nil
	}

	if peeked[versionPosition] != versionV1 {
		return nil, nil, nil, errors.New("unsupported crypto header version")
	}

	id := peeked[cryptorIdPosition : cryptorIdPosition+cryptorIdLength]
	var headerSize int64
	position := int64(sizePosition)
	if peeked[sizePosition] == longSizeIndicator {
		position += longSizeLength
		var e error
		headerSize, e = strconv.ParseInt(string(peeked[sizePosition:sizePosition+longSizeLength]), 10, 32)
		if e != nil {
			return nil, nil, nil, e
		}
	} else {
		position += shortSizeLength
		headerSize = int64(peeked[sizePosition])
	}
	if headerSize > 254 {
		_, e := bufData.Discard(sentinelLength + 1 + cryptorIdLength + longSizeLength)
		if e != nil {
			return nil, nil, nil, e
		}
	} else {
		_, e := bufData.Discard(sentinelLength + 1 + cryptorIdLength + shortSizeLength)
		if e != nil {
			return nil, nil, nil, e
		}
	}
	m := make([]byte, headerSize)
	_, e = io.ReadFull(bufData, metadata)
	if e != nil {
		return nil, nil, nil, e
	}
	return id, bufData, m, nil
}
