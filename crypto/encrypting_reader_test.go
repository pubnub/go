package crypto

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
	"testing/quick"
)

func encryptingReaderCanReadDifferentSizeOfChunks(in []byte, bufferSize uint8) bool {
	inPadded := padWithPKCS7VarBlock(in, 16)
	if bufferSize == 0 {
		return true
	}
	encryptingReader := NewBlockModeEncryptingReader(bytes.NewReader(in), &DoNothingBlockMode{})
	readDataBuffer := bytes.NewBuffer(nil)
	buffer := make([]byte, bufferSize)
	numberOfReadBytes := 0

	var e error
	var readBytes int

	for e == nil {
		readBytes, e = encryptingReader.Read(buffer)
		numberOfReadBytes += readBytes
		readDataBuffer.Write(buffer[:readBytes])
	}

	out := readDataBuffer.Bytes()[:numberOfReadBytes]

	return bytes.Equal(inPadded, out)
}

func Test_EncryptingReader_ReadDifferentSizeOfBuffers(t *testing.T) {
	c := quick.Config{MaxCount: 10000}
	if err := quick.Check(encryptingReaderCanReadDifferentSizeOfChunks, &c); err != nil {
		t.Error(err)
	}
}

func Test_EncryptingReader_ReadDifferentSizeOfBuffers1(t *testing.T) {
	assert.True(t, encryptingReaderCanReadDifferentSizeOfChunks([]byte{0x25, 0xb4, 0x40, 0x82, 0x6c, 0xbf, 0x30, 0xd1, 0xad, 0x8e, 0x31, 0x4f, 0x72, 0xd5, 0xbb, 0xd3, 0xc7, 0x2c, 0xf1, 0x60, 0x68, 0x76, 0x98, 0xda, 0x36, 0x3d, 0xe5, 0xc0, 0xfd, 0xd1, 0x57, 0x44, 0xca, 0xcf, 0xd3, 0xd1}, 0x21))
}
