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

	legacy := NewCrypto(legacyCryptor, []Cryptor{cbcCryptor}, legacyCryptor, Headless)
	newC := NewCrypto(cbcCryptor, []Cryptor{cbcCryptor}, legacyCryptor, CryptoHeaderV1)

	r1, _ := legacy.Encrypt([]byte("aaaa"))
	r2, _ := newC.Encrypt([]byte("aaaa"))
	r3, _ := newC.Decrypt(r1)
	r4, _ := newC.Decrypt(r2)
	println(base64.StdEncoding.EncodeToString(r1))
	println(EncryptString("enigma", "aaaa", true))
	println(base64.StdEncoding.EncodeToString(r2))
	println(string(r3))
	println(string(r4))
}
