package crypto

import (
	"bytes"
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
