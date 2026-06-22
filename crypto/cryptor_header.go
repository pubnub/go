package crypto

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	math "math"
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
var errTruncatedHeader = errors.New("decryption error: truncated header")

func headerV1(cryptorId string, metadata []byte) ([]byte, error) {
	cryptorDataSize := len(metadata)
	var cryptorDataBytesSize int

	if cryptorDataSize < longSizeIndicator {
		cryptorDataBytesSize = shortSizeLength
	} else if cryptorDataSize < math.MaxUint16 {
		cryptorDataBytesSize = longSizeLength
	} else {
		return nil, fmt.Errorf("size of cryptor metadata too large %d", cryptorDataSize)
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
		sizeBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(sizeBytes, uint16(cryptorDataSize))
		_, e = buffer.Write(sizeBytes)
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

func peekHeaderCryptorId(data []byte) (cryptorId *string, e error) {
	if len(data) < len(sentinel) || !bytes.Equal(data[:len(sentinel)], sentinel[:]) {
		return &legacyId, nil
	}

	if len(data) < cryptorIdPosition+cryptorIdLength {
		return nil, errTruncatedHeader
	}

	if data[versionPosition] != versionV1 {
		return nil, unsupportedHeaderVersion(int(data[versionPosition]))
	}

	id := string(data[cryptorIdPosition : cryptorIdPosition+cryptorIdLength])
	return &id, nil
}

func parseHeader(data []byte) (cryptorId *string, encrData *EncryptedData, e error) {
	id, err := peekHeaderCryptorId(data)
	if err != nil {
		return nil, nil, err
	}
	if (*id) == legacyId {
		return id, &EncryptedData{Data: data, Metadata: nil}, nil
	}
	if len(data) < sizePosition+shortSizeLength {
		return nil, nil, errTruncatedHeader
	}
	var headerSize int64
	position := int64(sizePosition)
	if data[sizePosition] == longSizeIndicator {
		if len(data) < sizePosition+longSizeLength {
			return nil, nil, errTruncatedHeader
		}
		position += longSizeLength

		headerSize = int64(binary.BigEndian.Uint16(data[sizePosition+shortSizeLength : sizePosition+longSizeLength]))
	} else {
		position += shortSizeLength
		headerSize = int64(data[sizePosition])
	}

	if int64(len(data)) < position+headerSize {
		return nil, nil, errTruncatedHeader
	}

	metadata := data[position : position+headerSize]
	position += int64(len(metadata))

	return id, &EncryptedData{Data: data[position:], Metadata: metadata}, nil
}

func parseHeaderStream(bufData *bufio.Reader) (cryptorId *string, encrypted *EncryptedStreamData, e error) {
	peeked, err := bufData.Peek(sentinelLength)
	if err != nil {
		return &legacyId, &EncryptedStreamData{
			Reader:   bufData,
			Metadata: nil,
		}, nil
	}
	if !bytes.Equal(peeked[:sentinelLength], sentinel[:]) {
		return &legacyId, &EncryptedStreamData{
			Reader:   bufData,
			Metadata: nil,
		}, nil
	}

	peeked, err = bufData.Peek(sizePosition + shortSizeLength)
	if err != nil {
		return nil, nil, errTruncatedHeader
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
	longSize := peeked[sizePosition] == longSizeIndicator
	position := int64(sizePosition)
	if longSize {
		peeked, err = bufData.Peek(sizePosition + longSizeLength)
		if err != nil {
			return nil, nil, errTruncatedHeader
		}
		position += longSizeLength
		metadataSize = int64(binary.BigEndian.Uint16(peeked[sizePosition+shortSizeLength : sizePosition+longSizeLength]))
	} else {
		position += shortSizeLength
		metadataSize = int64(peeked[sizePosition])
	}
	if longSize {
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
		return nil, nil, errTruncatedHeader
	}

	return id, &EncryptedStreamData{
		Reader:   bufData,
		Metadata: m,
	}, nil
}
