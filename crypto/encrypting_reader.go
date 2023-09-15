package crypto

import (
	"bufio"
	"bytes"
	"crypto/cipher"
	"errors"
	"io"
)

func NewBlockModeEncryptingReader(r io.Reader, mode cipher.BlockMode) io.Reader {
	return &blockModeEncryptingReader{
		r:         bufio.NewReader(r),
		blockMode: mode,
		buffer:    bytes.NewBuffer(nil),
		err:       nil,
	}
}

type blockModeEncryptingReader struct {
	r         *bufio.Reader
	blockMode cipher.BlockMode
	buffer    *bytes.Buffer
	err       error
}

func (encryptingReader *blockModeEncryptingReader) readNextBlockPadded() (read []byte, err error) {
	output := make([]byte, encryptingReader.blockMode.BlockSize())
	sizeOfCurrentlyRead, readErr := io.ReadFull(encryptingReader.r, output)
	if readErr != nil && readErr != io.EOF && !errors.Is(readErr, io.ErrUnexpectedEOF) {
		return nil, readErr
	}

	if sizeOfCurrentlyRead == 0 && readErr == io.EOF {
		return padWithPKCS7VarBlock(nil, encryptingReader.blockMode.BlockSize()), io.EOF
	}

	if errors.Is(readErr, io.ErrUnexpectedEOF) {
		return padWithPKCS7VarBlock(output[:sizeOfCurrentlyRead], encryptingReader.blockMode.BlockSize()), io.EOF
	}

	return output, nil
}

func (encryptingReader *blockModeEncryptingReader) encryptUntilPFull(p []byte) (int, error) {
	var alreadyWrote int
	for alreadyWrote <= len(p) {
		block, e := encryptingReader.readNextBlockPadded()
		if e != nil && e != io.EOF {
			return alreadyWrote, e
		}
		encryptingReader.err = e

		encryptingReader.blockMode.CryptBlocks(block, block)
		copied := copy(p[alreadyWrote:], block)
		alreadyWrote += copied
		if copied < len(block) {
			encryptingReader.buffer.Write(block[copied:])
		}
		if e == io.EOF {
			break
		}
	}

	return alreadyWrote, nil
}

func (encryptingReader *blockModeEncryptingReader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, errors.New("cannot read into empty buffer")
	}

	if (encryptingReader.err != nil) && encryptingReader.buffer.Len() == 0 {
		return 0, encryptingReader.err
	}

	alreadyWrote := readFromBufferUntilEmpty(encryptingReader.buffer, p)

	if encryptingReader.err != nil {
		return alreadyWrote, nil
	}

	n, e := encryptingReader.encryptUntilPFull(p[alreadyWrote:])

	return alreadyWrote + n, e
}
