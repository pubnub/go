package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"io"
)

func encrypt(mode cipher.BlockMode, input []byte) []byte {
	cipherBytes := make([]byte, len(input))
	mode.CryptBlocks(cipherBytes, input)
	return cipherBytes
}

func decrypt(mode cipher.BlockMode, input []byte) (output []byte, err error) {
	//to handle decryption errors
	defer func() {
		if r := recover(); r != nil {
			output, err = nil, fmt.Errorf("decrypt error: %s", r)
		}
	}()
	decrypted := make([]byte, len(input))
	mode.CryptBlocks(decrypted, input)
	val, err := unpadPKCS7(decrypted)
	if err != nil {
		return nil, fmt.Errorf("decrypt error: %s", err)
	}

	return val, nil
}

func encryptStream(mode cipher.BlockMode, input io.Reader, output io.Writer) error {
	readDataSlice := make([]byte, aes.BlockSize)
	for {
		sizeOfCurrentlyRead, readErr := io.ReadFull(input, readDataSlice)

		if readErr != nil && readErr != io.EOF && !errors.Is(readErr, io.ErrUnexpectedEOF) {
			return readErr
		}

		if sizeOfCurrentlyRead < aes.BlockSize {
			pad := bytes.Repeat([]byte{byte(aes.BlockSize - sizeOfCurrentlyRead)}, aes.BlockSize-sizeOfCurrentlyRead)
			copy(readDataSlice[sizeOfCurrentlyRead:], pad)
		}

		mode.CryptBlocks(readDataSlice, readDataSlice)
		_, writeErr := output.Write(readDataSlice)
		if writeErr != nil {
			return writeErr
		}

		if readErr == io.EOF || errors.Is(readErr, io.ErrUnexpectedEOF) {
			break
		}

	}

	return nil
}

func decryptStream(mode cipher.BlockMode, input io.Reader, output io.Writer) error {
	readDataSlice := make([]byte, aes.BlockSize)

	for {
		_, readErr := io.ReadFull(input, readDataSlice)

		if readErr != nil && readErr != io.EOF {
			return readErr
		}

		if readErr == io.EOF {
			break
		}

		mode.CryptBlocks(readDataSlice, readDataSlice)
		unpadded, unpadError := unpadPKCS7(readDataSlice)
		if unpadError != nil {
			return unpadError
		}
		_, writeErr := output.Write(unpadded)
		if writeErr != nil {
			return writeErr
		}

		if errors.Is(readErr, io.ErrUnexpectedEOF) {
			break
		}

	}

	return nil
}

// padWithPKCS7 pads the data as per the PKCS7 standard
// It accepts the following parameters:
// data: data to pad as byte array.
// returns the padded data as byte array.
func padWithPKCS7(data []byte) []byte {
	blocklen := 16
	padlen := 1
	for ((len(data) + padlen) % blocklen) != 0 {
		padlen = padlen + 1
	}

	pad := bytes.Repeat([]byte{byte(padlen)}, padlen)
	return append(data, pad...)
}

// unpadPKCS7 unpads the data as per the PKCS7 standard
// It accepts the following parameters:
// data: data to unpad as byte array.
// returns the unpadded data as byte array.
func unpadPKCS7(data []byte) ([]byte, error) {
	blocklen := 16
	if len(data)%blocklen != 0 || len(data) == 0 {
		return nil, fmt.Errorf("invalid data len %d", len(data))
	}
	padlen := int(data[len(data)-1])
	if padlen > blocklen || padlen == 0 {
		return nil, fmt.Errorf("padding is invalid")
	}
	// check padding
	pad := data[len(data)-padlen:]
	for i := 0; i < padlen; i++ {
		if pad[i] != byte(padlen) {
			return nil, fmt.Errorf("padding is invalid")
		}
	}

	return data[:len(data)-padlen], nil
}
