package crypto

import (
	"bytes"
	"io"
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
	decryptingReader := newBlockModeDecryptingReader(bytes.NewReader(inPadded), &DoNothingBlockMode{})
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

	if e != nil && e != io.EOF {
		return false
	}

	out := readDataBuffer.Bytes()[:numberOfReadBytes]

	return bytes.Equal(in, out)
}

func Test_DecryptingReader_ReadDifferentSizeOfBuffers(t *testing.T) {
	if err := quick.Check(canReadDifferentSizeOfChunks, defaultPropertyTestConfig); err != nil {
		t.Error(err)
	}
}
