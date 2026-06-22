package crypto

import (
	"bufio"
	"bytes"
	"io"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCryptorHeader_CreateHeaderWithLargeMetadata(t *testing.T) {
	metadata := make([]byte, 512)
	header, _ := headerV1("abcd", metadata)
	cryptorDataSize := header[sizePosition : sizePosition+longSizeLength]
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

func TestCryptorHeader_ParseHeaderWithLargeMetadata(t *testing.T) {
	metadata := make([]byte, 512)
	encrypted := []byte("encrypted")
	header, _ := headerV1("abcd", metadata)

	id, data, e := parseHeader(append(header, encrypted...))

	assert.NoError(t, e)
	assert.Equal(t, "abcd", *id)
	assert.Equal(t, metadata, data.Metadata)
	assert.Equal(t, encrypted, data.Data)
}

func TestCryptorHeader_ParseHeaderStreamWithLargeMetadata(t *testing.T) {
	metadata := make([]byte, 512)
	encrypted := []byte("encrypted")
	header, _ := headerV1("abcd", metadata)
	reader := bufio.NewReader(bytes.NewReader(append(header, encrypted...)))

	id, data, e := parseHeaderStream(reader)

	assert.NoError(t, e)
	assert.Equal(t, "abcd", *id)
	assert.Equal(t, metadata, data.Metadata)
	remaining, e := io.ReadAll(data.Reader)
	assert.NoError(t, e)
	assert.Equal(t, encrypted, remaining)
}

func TestCryptorHeader_ParseHeaderStreamWithLongSizeIndicatorAndSmallMetadata(t *testing.T) {
	encrypted := []byte("encrypted")
	data := append([]byte{'P', 'N', 'E', 'D', 1, 'a', 'b', 'c', 'd', 0xff, 0x00, 0x01, 'm'}, encrypted...)
	reader := bufio.NewReader(bytes.NewReader(data))

	id, streamData, e := parseHeaderStream(reader)

	assert.NoError(t, e)
	assert.Equal(t, "abcd", *id)
	assert.Equal(t, []byte{'m'}, streamData.Metadata)
	remaining, e := io.ReadAll(streamData.Reader)
	assert.NoError(t, e)
	assert.Equal(t, encrypted, remaining)
}

func TestCryptorHeader_ParseMalformedHeaderReturnsError(t *testing.T) {
	_, _, e := parseHeader([]byte("PNED"))

	assert.Error(t, e)
}

func TestCryptorHeader_ParseHeaderWithOversizeMetadataReturnsError(t *testing.T) {
	data := append([]byte{'P', 'N', 'E', 'D', 1, 'a', 'b', 'c', 'd', 0xff, 0xff, 0xff}, []byte("short")...)

	_, _, e := parseHeader(data)

	assert.Error(t, e)
}

func TestCryptoModule_DecryptMalformedHeaderReturnsError(t *testing.T) {
	module, e := NewAesCbcCryptoModule("key", true)
	assert.NoError(t, e)

	_, e = module.Decrypt([]byte("PNED"))

	assert.Error(t, e)
}

func TestCryptoModule_DecryptStreamMalformedHeaderReturnsError(t *testing.T) {
	module, e := NewAesCbcCryptoModule("key", true)
	assert.NoError(t, e)

	_, e = module.DecryptStream(bytes.NewReader([]byte("PNED")))

	assert.Error(t, e)
}

func FuzzCryptorHeader_ParseHeader(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte("PNED"))
	f.Add([]byte{'P', 'N', 'E', 'D', 1, 'a', 'b', 'c', 'd', 0xff})
	f.Add([]byte{'P', 'N', 'E', 'D', 1, 'a', 'b', 'c', 'd', 0xff, 0x02, 0x00})
	f.Add([]byte("legacy payload"))

	f.Fuzz(func(t *testing.T, data []byte) {
		_, _, _ = parseHeader(data)
		_, _, _ = parseHeaderStream(bufio.NewReader(bytes.NewReader(data)))
	})
}
