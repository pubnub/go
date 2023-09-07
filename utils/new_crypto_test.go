package utils

import (
	"encoding/base64"
	"testing"
)

func TestJust(t *testing.T) {
	legacyCryptor, e1 := NewLegacyCryptor("enigmaenigmaenig", true)
	cbcCryptor, e2 := NewCBCCryptor("enigmaenigmaenig")

	if e1 != nil || e2 != nil {
		t.Errorf("error, %s, %s", e1, e2)
	}

	legacy := NewCrypto(legacyCryptor, []CryptoAlgorithm{cbcCryptor}, legacyCryptor, headless)
	newC := NewCrypto(cbcCryptor, []CryptoAlgorithm{cbcCryptor}, legacyCryptor, CryptoHeaderV1)

	r1, _ := legacy.Encrypt([]byte("sure i can"))
	r2, _ := newC.Encrypt([]byte("sure i can"))
	r3, _ := newC.Decrypt(r1)
	r4, _ := newC.Decrypt(r2)
	println(base64.StdEncoding.EncodeToString(r1))
	println(EncryptString("enigma", "sure i can", true))
	println(base64.StdEncoding.EncodeToString(r2))
	println(string(r3))
	println(string(r4))
}
