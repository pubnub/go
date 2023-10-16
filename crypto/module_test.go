package crypto

import (
	"encoding/base64"
	"testing"
)

func TestToProduceData(t *testing.T) {
	cipherKey := "myCipherKey"
	textToEncrypt := "Hello world encrypted with "

	legacyModuleStaticIv, _ := NewLegacyCryptoModule(cipherKey, false)
	legacyModuleRandomIv, _ := NewLegacyCryptoModule(cipherKey, true)
	aesCbcModule, _ := NewAesCbcCryptoModule(cipherKey, true)
	r1, _ := legacyModuleStaticIv.Encrypt([]byte(textToEncrypt + "legacyModuleStaticIv"))
	r2, _ := legacyModuleRandomIv.Encrypt([]byte(textToEncrypt + "legacyModuleRandomIv"))
	r3, _ := aesCbcModule.Encrypt([]byte(textToEncrypt + "aesCbcModule"))

	println(base64.StdEncoding.EncodeToString(r1))
	println(base64.StdEncoding.EncodeToString(r2))
	println(base64.StdEncoding.EncodeToString(r3))
}
