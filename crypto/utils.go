package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"fmt"
)

func generateIV(blocksize int) []byte {
	iv := make([]byte, blocksize)
	if _, err := rand.Read(iv); err != nil {
		panic(err)
	}
	return iv
}

func padWithPKCS7VarBlock(data []byte, blocklen int) []byte {
	padlen := 1
	for ((len(data) + padlen) % blocklen) != 0 {
		padlen = padlen + 1
	}

	pad := bytes.Repeat([]byte{byte(padlen)}, padlen)
	return append(data, pad...)
}

// padWithPKCS7 pads the data as per the PKCS7 standard
// It accepts the following parameters:
// data: data to pad as byte array.
// returns the padded data as byte array.
func padWithPKCS7(data []byte) []byte {
	return padWithPKCS7VarBlock(data, 16)
}

// unpadPKCS7 unpads the data as per the PKCS7 standard
// It accepts the following parameters:
// data: data to unpad as byte array.
// returns the unpadded data as byte array.
func unpadPKCS7(data []byte) ([]byte, error) {
	blocklen := 16
	if len(data) == 0 || len(data)%blocklen != 0 {
		return nil, errors.New("decryption error")
	}

	// Validate the PKCS#7 padding in constant time. The padding length is the
	// last byte and the last padlen bytes must all equal padlen. Branching on
	// the pad length or on which byte is wrong would leak a padding oracle, so
	// every byte of the final block is examined and the result is folded into
	// "good", which stays 1 only when the padding is fully valid.
	n := len(data)
	padlen := int(data[n-1])
	good := subtle.ConstantTimeLessOrEq(1, padlen) & subtle.ConstantTimeLessOrEq(padlen, blocklen)
	for j := 1; j <= blocklen; j++ {
		inPad := subtle.ConstantTimeLessOrEq(j, padlen)
		isPadByte := subtle.ConstantTimeByteEq(data[n-j], byte(padlen))
		good &= isPadByte | (1 ^ inPad)
	}

	if good != 1 {
		return nil, errors.New("decryption error")
	}

	return data[:n-padlen], nil
}

func unsupportedHeaderVersion(version int) error {
	return fmt.Errorf("unknown crypto error: unsupported crypto header version %d", version)
}
