package crypto

import (
	"bytes"
	"testing"
	"testing/quick"
)

type DoNothingBlockMode struct {
}

func (d *DoNothingBlockMode) BlockSize() int {
	return 16
}

func (d *DoNothingBlockMode) CryptBlocks(dst, src []byte) {
	copy(dst, src)
}

func canReadDifferentSizeOfChunks(in []byte, bufferSize uint8) bool {
	inPadded := padWithPKCS7VarBlock(in, 16)
	if bufferSize == 0 {
		return true
	}
	decryptingReader := NewBlockModeDecryptingReader(bytes.NewReader(inPadded), &DoNothingBlockMode{})
	readDataBuffer := bytes.NewBuffer(nil)
	buffer := make([]byte, bufferSize)
	numberOfReadBytes := 0

	var e error
	var readBytes int

	for e == nil {
		readBytes, e = decryptingReader.Read(buffer)
		numberOfReadBytes += readBytes
		readDataBuffer.Write(buffer[:readBytes])
	}

	out := readDataBuffer.Bytes()[:numberOfReadBytes]

	return bytes.Equal(in, out)
}

func Test_DecryptingReader_ReadDifferentSizeOfBuffers(t *testing.T) {
	c := quick.Config{MaxCount: 10000}
	if err := quick.Check(canReadDifferentSizeOfChunks, &c); err != nil {
		t.Error(err)
	}
}
