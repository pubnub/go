package crypto

import (
	"bytes"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCryptorHeader_CreateHeaderWithLargeMetadata(t *testing.T) {
	metadata := make([]byte, 512)
	header, _ := headerV1("abcd", metadata)
	cryptorDataSize := header[sentinelLength+1+cryptorIdLength+longSizeLength-2 : sentinelLength+1+cryptorIdLength+longSizeLength]
	assert.True(t, bytes.Equal([]byte{0xff, 0x02, 0x00}, cryptorDataSize))
}

func TestCryptorHeader_CreateHeaderWithSmallMetadata(t *testing.T) {
	metadata := make([]byte, 128)
	header, _ := headerV1("abcd", metadata)
	assert.Equal(t, byte(0x80), header[sizePosition])
}

func TestCryptorHeader_WithTooLargeMetadata(t *testing.T) {
	metadata := make([]byte, math.MaxUint16+255)
	_, e := headerV1("abcd", metadata)
	if e == nil {
		assert.Fail(t, "expected error")
	}
}
