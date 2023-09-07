package utils

import (
	"crypto/rand"
)

func generateIV(blocksize int) []byte {
	iv := make([]byte, blocksize)
	if _, err := rand.Read(iv); err != nil {
		panic(err)
	}
	return iv
}
