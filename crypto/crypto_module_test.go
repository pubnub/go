package crypto

import (
	"encoding/base64"
	"testing"
)

func TestJust(t *testing.T) {
	cipherKey := "enigma"

	legacy, e := NewLegacyCryptoModule(cipherKey, true)
	if e != nil {
		t.Errorf(e.Error())
	}
	newC, e := NewCryptoModule(cipherKey, true)
	if e != nil {
		t.Errorf(e.Error())
	}

	r1, _ := legacy.Encrypt([]byte("sure i can"))
	r2, _ := newC.Encrypt([]byte("sure i can"))
	r3, _ := newC.Decrypt(r1)
	r4, _ := newC.Decrypt(r2)
	println(base64.StdEncoding.EncodeToString(r1))
	println(base64.StdEncoding.EncodeToString(r2))
	println(string(r3))
	println(string(r4))
}
