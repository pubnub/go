package crypto

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
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
		return &legacyId, nil
	}

	if data[versionPosition] != versionV1 {
		return nil, unsupportedHeaderVersion(int(data[versionPosition]))
	}

	id := string(data[cryptorIdPosition : cryptorIdPosition+cryptorIdLength])
	return &id, nil
}

func parseHeader(data []byte) (cryptorId *string, encryptedData *EncryptedData, e error) {
	id, err := peekHeaderCryptorId(data)
	if err != nil {
		return nil, nil, err
	}
	if (*id) == legacyId {
		return id, &EncryptedData{Data: data, Metadata: nil}, nil
	}
	var headerSize int64
	position := int64(sizePosition)
	if data[sizePosition] == longSizeIndicator {
		position += longSizeLength

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

	return id, &EncryptedData{Data: data[position:], Metadata: metadata}, nil
}

func parseHeaderStream(bufData *bufio.Reader) (cryptorId *string, encrypted *EncryptedStreamData, e error) {
	peeked, err := bufData.Peek(sentinelLength + 1 + cryptorIdLength + longSizeLength)
	if err != nil {
		return nil, nil, fmt.Errorf("decryption error: %w", err)
	}

	id, err := peekHeaderCryptorId(peeked)
	if err != nil {
		return nil, nil, err
	}

	if (*id) == legacyId {
		return id, &EncryptedStreamData{
			Reader:   bufData,
			Metadata: nil,
		}, nil
	}

	var metadataSize int64
	position := int64(sizePosition)
	if peeked[sizePosition] == longSizeIndicator {
		position += longSizeLength
		var e error
		metadataSize, e = strconv.ParseInt(string(peeked[sizePosition:sizePosition+longSizeLength]), 10, 32)
		if e != nil {
			return nil, nil, e
		}
	} else {
		position += shortSizeLength
		metadataSize = int64(peeked[sizePosition])
	}
	if metadataSize > 254 {
		_, e := bufData.Discard(sentinelLength + 1 + cryptorIdLength + longSizeLength)
		if e != nil {
			return nil, nil, e
		}
	} else {
		_, e := bufData.Discard(sentinelLength + 1 + cryptorIdLength + shortSizeLength)
		if e != nil {
			return nil, nil, e
		}
	}
	m := make([]byte, metadataSize)
	_, e = io.ReadFull(bufData, m)
	if e != nil {
		return nil, nil, e
	}
	return id, &EncryptedStreamData{
		Reader:   bufData,
		Metadata: m,
	}, nil
}
