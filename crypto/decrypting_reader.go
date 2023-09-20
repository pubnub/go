package crypto

import (
	"bufio"
	"bytes"
	"crypto/cipher"
	"errors"
	"io"
)

func newBlockModeDecryptingReader(r io.Reader, mode cipher.BlockMode) io.Reader {
	return &blockModeDecryptingReader{
		r:         bufio.NewReader(r),
		blockMode: mode,
		buffer:    bytes.NewBuffer(nil),
		err:       nil,
	}
}

type blockModeDecryptingReader struct {
	r         *bufio.Reader
	blockMode cipher.BlockMode
	buffer    *bytes.Buffer
	err       error
}

func (decryptingReader *blockModeDecryptingReader) readNextBlock() ([]byte, error) {
	reader := decryptingReader.r

	output := make([]byte, decryptingReader.blockMode.BlockSize())
	sizeOfCurrentlyRead, readErr := io.ReadFull(reader, output)
	if readErr != nil && readErr != io.EOF {
		return nil, readErr
	}

	if sizeOfCurrentlyRead == 0 && readErr == io.EOF {
		return nil, io.EOF
	}

	if _, e := reader.Peek(1); e == io.EOF {
		return output[:sizeOfCurrentlyRead], io.EOF
	}

	return output, nil
}

func (decryptingReader *blockModeDecryptingReader) decryptUntilPFull(p []byte) (n int, err error) {
	var copied int
	var block []byte
	var e error
	alreadyWrote := 0
	for alreadyWrote < len(p) {
		block, e = decryptingReader.readNextBlock()
		decryptingReader.err = e

		if errors.Is(e, io.ErrUnexpectedEOF) {
			return alreadyWrote, errors.New("encrypted data length is not a multiple of the block size")
		}

		if e != nil && e != io.EOF {
			return alreadyWrote, e
		}

		decryptingReader.blockMode.CryptBlocks(block, block)

		if bytes.Equal(block, bytes.Repeat([]byte{byte(decryptingReader.blockMode.BlockSize())}, len(block))) {
			break
		} else if e == io.EOF {
			unpadded, unpadErr := unpadPKCS7(block)
			if len(unpadded) > 0 {
				block = unpadded
				copied = copy(p[alreadyWrote:], block)
				alreadyWrote += copied
			}

			if unpadErr != nil {
				return alreadyWrote, unpadErr
			}
			if copied < len(block) {
				decryptingReader.buffer.Write(block[copied:])
			}
			break
		} else {
			copied = copy(p[alreadyWrote:], block)
			alreadyWrote += copied
			if copied < len(block) {
				decryptingReader.buffer.Write(block[copied:])
			}
		}
	}

	return alreadyWrote, nil
}

func readFromBufferUntilEmpty(buffer *bytes.Buffer, p []byte) int {
	buffered := buffer.Next(len(p))
	if len(buffered) == 0 {
		return 0
	}
	return copy(p, buffered)
}

func (decryptingReader *blockModeDecryptingReader) readFromBufferUntilEmpty(p []byte) (n int, err error) {
	buffered := decryptingReader.buffer.Next(len(p))
	var alreadyWrote int
	if len(buffered) > 0 {
		alreadyWrote = copy(p, buffered)
	}

	if decryptingReader.err != nil {
		return alreadyWrote, decryptingReader.err
	}

	return alreadyWrote, nil
}

func (decryptingReader *blockModeDecryptingReader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, errors.New("cannot read into empty buffer")
	}

	if (decryptingReader.err != nil) && decryptingReader.buffer.Len() == 0 {
		return 0, decryptingReader.err
	}

	alreadyWrote := readFromBufferUntilEmpty(decryptingReader.buffer, p)

	if decryptingReader.err != nil {
		return alreadyWrote, nil
	}

	n, e := decryptingReader.decryptUntilPFull(p[alreadyWrote:])
	return alreadyWrote + n, e
}
