package utils

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/stretchr/testify/assert"

	//"net/url"
	"os"
	"testing"
	"unicode/utf16"
)

func TestSignSha256(t *testing.T) {
	assert := assert.New(t)

	signInput := "sub-c-7ba2ac4c-4836-11e6-85a4-0619f8945a4f\npub-c-98863562-19a6-4760-bf0b-d537d1f5c582\ngrant\nchannel=asyncio-pam-FI2FCS0A&pnsdk=PubNub-Python-Asyncio%252F4.0.0&r=1&timestamp=1468409553&uuid=a4dbf92e-e5cb-428f-b6e6-35cce03500a2&w=1"

	res := GetHmacSha256("my_key", signInput)

	assert.Equal("Dq92jnwRTCikdeP2nUs1__gyJthF8NChwbs5aYy2r_I=", res)
}

// func TestSignSha256New2(t *testing.T) {
// 	assert := assert.New(t)
// 	v := &url.Values{}
// 	v.Set("PoundsSterling=", "%C2%A313.37")
// 	v.Set("auth=", "joker")
// 	v.Set("r=", "1")
// 	v.Set("timestamp=", "123456789")
// 	v.Set("ttl=", "60")
// 	v.Set("w=", "1")

// 	d := PreparePamParams(v)

// 	signInput := fmt.Sprintf("%s\n%s\n%s\n%s\n", "demo", "demo", "/v2/auth/grant/sub-key/demo", d)

// 	res := GetHmacSha256("wMfbo9G0xVUG8yfTfYw5qIdfJkTd7A", signInput)

// 	assert.Equal("v2rgQQ1eFzk8omugFV9V1_eKRUvvMv9jyC9Z-L1ogdw=", res)
// }

// func TestSignSha256New(t *testing.T) {
// 	assert := assert.New(t)
// 	v := &url.Values{}
// 	v.Set("store=", "1")
// 	v.Set("seqn=", "1")
// 	v.Set("auth=", "myAuth")
// 	v.Set("timestamp=", "1535125017")
// 	v.Set("pnsdk=", "PubNub-Go/4.1.2")
// 	v.Set("uuid=", "myUuid")

// 	d := PreparePamParams(v)

// 	signInput := fmt.Sprintf("%s\n%s\n%s\n%s\n", "demoSubscribeKey", "demoPublishKey", "/publish/demoPublishKey/demoSubscribeKey/0/my-channel/0/%22my-message%22", d)

// 	res := GetHmacSha256("secretKey", signInput)

// 	assert.Equal("Dq92jnwRTCikdeP2nUs1__gyJthF8NChwbs5aYy2r_I=", res)
// }

func TestPad(t *testing.T) {
	assert := assert.New(t)

	var badMsg interface{}
	b := []byte(`{
	"kind": "click",
	"user": {"key" : "user@test.com"},
	"creationDate": 9223372036854775808346,
	"key": "54651fa39868621628000002",
	"url": "http://www.google.com"
	}`)

	json.Unmarshal(b, &badMsg)
	jsonSerialized, _ := json.Marshal(badMsg)

	actual := EncryptString("enigma", fmt.Sprintf("%s", jsonSerialized), false)
	expected := "yzJ2MMyt8So18nNXm4m3Dqzb1G+as9LDqdlZ+p8iEGi358F5h25wmKrj9FTOPdMQ0TMy/Xhf3hS3+ZRUlv/zLD6/0Ns/c834HQMUmG+6DN9SQy9II3bkUGZu9Bn6Ng/ZmJTrHV7QnkLnjD+pGOHEvqrPEduR5pfA2n9mA3qQNhqFgnsIvffxGB0AqM57NdD3Tlr2ig8A2VI4Lh3DmX7f1Q=="

	assert.Equal(expected, actual)
}

func TestRandomIVYayEncryption(t *testing.T) {
	EncryptionAndDecryptionWithRandomIVCommon(t, []byte("yay!"))
}

func TestRandomIVSomeBytesEncryption(t *testing.T) {
	b := []byte(`{
		"kind": "click",
		"user": {"key" : "user@test.com"},
		"creationDate": 9223372036854775808346,
		"key": "54651fa39868621628000002",
		"url": "http://www.google.com"
		}`)
	EncryptionAndDecryptionWithRandomIVCommon(t, []byte(b))
}

func TestRandomIVCustomStructEncryption(t *testing.T) {
	message := customStruct{
		Foo: "hi!",
		Bar: []int{1, 2, 3, 4, 5},
	}
	b1, _ := json.Marshal(message)

	EncryptionAndDecryptionWithRandomIVCommon(t, []byte(b1))
}

func EncryptionAndDecryptionWithRandomIVCommon(t *testing.T, msg []byte) {
	assert := assert.New(t)
	encmsg := EncryptString("enigma", fmt.Sprintf("%s", msg), true)
	decrypted, _ := DecryptString("enigma", encmsg, true)
	decMessage := fmt.Sprintf("%s", decrypted)
	assert.Equal(string(msg), decMessage)

}

func TestUnpad(t *testing.T) {
	assert := assert.New(t)

	message := "yzJ2MMyt8So18nNXm4m3Dl0XuYAOJFj2JXG8P3BGlCsDsqM44ReH15MRGbEkJZCSqgMiX1wUK44Qz8gsTcmGcZm/7KtOa+kRnvgDpNkTuBUrDqSjmYeuBLqRIEIfoGrRNljbFmP1W9Zv8iVbJMmovF+gmNNiIzlC3J9dHK51/OgW7s2EASMQJr3UJZ26PoFmmXY/wYN+2EyRnT4PBRCocQ=="
	decrypted, _ := DecryptString("enigma", message, false)

	decMessage := fmt.Sprintf("%s", decrypted)

	assert.Contains(decMessage, `"user":{"key":"user@test.com"}`)
	assert.Contains(decMessage, `"key":"54651fa39868621628000002"`)
}

// TestYayDecryptionBasic tests the yay decryption.
// Assumes that the input message is deserialized
// Decrypted string should match yay!
func TestYayDecryptionBasic(t *testing.T) {
	assert := assert.New(t)

	message := "q/xJqqN6qbiZMXYmiQC1Fw=="

	decrypted, decErr := DecryptString("enigma", message, false)
	assert.NoError(decErr)

	assert.Equal("yay!", decrypted)
}

// TestYayEncryptionBasic tests the yay encryption.
// Assumes that the input message is not serialized
// Decrypted string should match q/xJqqN6qbiZMXYmiQC1Fw==
func TestYayEncryptionBasic(t *testing.T) {
	assert := assert.New(t)

	message := "yay!"
	encrypted := EncryptString("enigma", message, false)

	assert.Equal("q/xJqqN6qbiZMXYmiQC1Fw==", encrypted)
}

// TestYayDecryption tests the yay decryption.
// Assumes that the input message is serialized
// Decrypted string should match yay!
func TestYayDecryption(t *testing.T) {
	assert := assert.New(t)

	message := "Wi24KS4pcTzvyuGOHubiXg=="

	decrypted, decErr := DecryptString("enigma", message, false)
	assert.NoError(decErr)

	b, err := json.Marshal("yay!")
	assert.NoError(err)

	assert.Equal(string(b), decrypted)
}

// TestYayEncryption tests the yay encryption.
// Assumes that the input message is serialized
// Decrypted string should match q/xJqqN6qbiZMXYmiQC1Fw==
func TestYayEncryption(t *testing.T) {
	assert := assert.New(t)

	message := "yay!"
	b, err := json.Marshal(message)
	assert.NoError(err)

	encrypted := EncryptString("enigma", string(b), false)
	assert.Equal("Wi24KS4pcTzvyuGOHubiXg==", encrypted)
}

// TestArrayDecryption tests the slice decryption.
// Assumes that the input message is deserialized
// And the output message has to been deserialized.
// Decrypted string should match Ns4TB41JjT2NCXaGLWSPAQ==
func TestArrayDecryption(t *testing.T) {
	assert := assert.New(t)

	message := "Ns4TB41JjT2NCXaGLWSPAQ=="
	decrypted, decErr := DecryptString("enigma", message, false)
	assert.NoError(decErr)
	slice := []string{}
	b, err := json.Marshal(slice)
	assert.NoError(err)

	assert.Equal(string(b), decrypted)
}

// TestArrayEncryption tests the slice encryption.
// Assumes that the input message is not serialized
// Decrypted string should match Ns4TB41JjT2NCXaGLWSPAQ==
func TestArrayEncryption(t *testing.T) {
	assert := assert.New(t)

	message := []string{}

	b, err := json.Marshal(message)
	assert.NoError(err)

	encrypted := EncryptString("enigma", string(b), false)
	assert.Equal("Ns4TB41JjT2NCXaGLWSPAQ==", encrypted)
}

// TestObjectDecryption tests the empty object decryption.
// Assumes that the input message is deserialized
// And the output message has to been deserialized.
// Decrypted string should match IDjZE9BHSjcX67RddfCYYg==
func TestObjectDecryption(t *testing.T) {
	assert := assert.New(t)

	message := "IDjZE9BHSjcX67RddfCYYg=="
	decrypted, decErr := DecryptString("enigma", message, false)
	assert.NoError(decErr)

	emptyStruct := emptyStruct{}

	b, err := json.Marshal(emptyStruct)
	assert.NoError(err)
	assert.Equal(string(b), decrypted)
}

// TestObjectEncryption tests the empty object encryption.
// The output is not serialized
// Encrypted string should match the serialized object
func TestObjectEncryption(t *testing.T) {
	assert := assert.New(t)

	message := emptyStruct{}

	b, err := json.Marshal(message)
	assert.NoError(err)

	encrypted := EncryptString("enigma", string(b), false)

	assert.Equal("IDjZE9BHSjcX67RddfCYYg==", encrypted)
}

// TestMyObjectDecryption tests the custom object decryption.
// Assumes that the input message is deserialized
// And the output message has to been deserialized.
// Decrypted string should match BMhiHh363wsb7kNk7krTtDcey/O6ZcoKDTvVc4yDhZY=
func TestMyObjectDecryption(t *testing.T) {
	assert := assert.New(t)

	message := "BMhiHh363wsb7kNk7krTtDcey/O6ZcoKDTvVc4yDhZY="
	decrypted, decErr := DecryptString("enigma", message, false)

	assert.NoError(decErr)
	customStruct := customStruct{
		Foo: "hi!",
		Bar: []int{1, 2, 3, 4, 5},
	}
	b, err := json.Marshal(customStruct)
	assert.NoError(err)
	assert.Equal(string(b), decrypted)
}

// TestMyObjectEncryption tests the custom object encryption.
// The output is not serialized
// Encrypted string should match the serialized object
func TestMyObjectEncryption(t *testing.T) {
	assert := assert.New(t)

	message := customStruct{
		Foo: "hi!",
		Bar: []int{1, 2, 3, 4, 5},
	}

	b1, err := json.Marshal(message)
	assert.NoError(err)

	encrypted := EncryptString("enigma", string(b1), false)
	assert.Equal("BMhiHh363wsb7kNk7krTtDcey/O6ZcoKDTvVc4yDhZY=", encrypted)
}

// TestPubNubDecryption2 tests the Pubnub Messaging API 2 decryption.
// Assumes that the input message is deserialized
// Decrypted string should match Pubnub Messaging API 2
func TestPubNubDecryption2(t *testing.T) {
	assert := assert.New(t)

	message := "f42pIQcWZ9zbTbH8cyLwB/tdvRxjFLOYcBNMVKeHS54="
	decrypted, decErr := DecryptString("enigma", message, false)
	assert.NoError(decErr)

	b, err := json.Marshal("Pubnub Messaging API 2")
	assert.NoError(err)
	assert.Equal(string(b), decrypted)
}

// TestPubNubEncryption2 tests the Pubnub Messaging API 2 encryption.
// Assumes that the input message is not serialized
// Decrypted string should match f42pIQcWZ9zbTbH8cyLwB/tdvRxjFLOYcBNMVKeHS54=
func TestPubNubEncryption2(t *testing.T) {
	assert := assert.New(t)

	message := "Pubnub Messaging API 2"
	b, err := json.Marshal(message)
	assert.NoError(err)

	encrypted := EncryptString("enigma", string(b), false)
	assert.Equal("f42pIQcWZ9zbTbH8cyLwB/tdvRxjFLOYcBNMVKeHS54=", encrypted)
}

// TestPubNubDecryption tests the Pubnub Messaging API 1 decryption.
// Assumes that the input message is deserialized
// Decrypted string should match Pubnub Messaging API 1
func TestPubNubDecryption(t *testing.T) {
	assert := assert.New(t)

	message := "f42pIQcWZ9zbTbH8cyLwByD/GsviOE0vcREIEVPARR0="
	decrypted, decErr := DecryptString("enigma", message, false)
	assert.NoError(decErr)

	b, err := json.Marshal("Pubnub Messaging API 1")
	assert.NoError(err)
	assert.Equal(string(b), decrypted)
}

// TestPubNubEncryption tests the Pubnub Messaging API 1 encryption.
// Assumes that the input message is not serialized
// Decrypted string should match f42pIQcWZ9zbTbH8cyLwByD/GsviOE0vcREIEVPARR0=
func TestPubNubEncryption(t *testing.T) {
	assert := assert.New(t)

	message := "Pubnub Messaging API 1"
	b, err := json.Marshal(message)
	assert.NoError(err)
	encrypted := EncryptString("enigma", string(b), false)
	assert.Equal("f42pIQcWZ9zbTbH8cyLwByD/GsviOE0vcREIEVPARR0=", encrypted)
}

// TestStuffCanDecryption tests the StuffCan decryption.
// Assumes that the input message is deserialized
// Decrypted string should match {\"this stuff\":{\"can get\":\"complicated!\"}}
func TestStuffCanDecryption(t *testing.T) {
	assert := assert.New(t)

	message := "zMqH/RTPlC8yrAZ2UhpEgLKUVzkMI2cikiaVg30AyUu7B6J0FLqCazRzDOmrsFsF"
	decrypted, decErr := DecryptString("enigma", message, false)
	assert.NoError(decErr)
	assert.Equal("{\"this stuff\":{\"can get\":\"complicated!\"}}", decrypted)
}

// TestStuffCanEncryption tests the StuffCan encryption.
// Assumes that the input message is not serialized
// Decrypted string should match zMqH/RTPlC8yrAZ2UhpEgLKUVzkMI2cikiaVg30AyUu7B6J0FLqCazRzDOmrsFsF
func TestStuffCanEncryption(t *testing.T) {
	assert := assert.New(t)

	message := "{\"this stuff\":{\"can get\":\"complicated!\"}}"
	encrypted := EncryptString("enigma", message, false)
	assert.Equal("zMqH/RTPlC8yrAZ2UhpEgLKUVzkMI2cikiaVg30AyUu7B6J0FLqCazRzDOmrsFsF", encrypted)
}

// TestHashDecryption tests the hash decryption.
// Assumes that the input message is deserialized
// Decrypted string should match {\"foo\":{\"bar\":\"foobar\"}}
func TestHashDecryption(t *testing.T) {
	assert := assert.New(t)

	message := "GsvkCYZoYylL5a7/DKhysDjNbwn+BtBtHj2CvzC4Y4g="
	decrypted, decErr := DecryptString("enigma", message, false)
	assert.NoError(decErr)
	assert.Equal("{\"foo\":{\"bar\":\"foobar\"}}", decrypted)
}

// TestHashEncryption tests the hash encryption.
// Assumes that the input message is not serialized
// Decrypted string should match GsvkCYZoYylL5a7/DKhysDjNbwn+BtBtHj2CvzC4Y4g=
func TestHashEncryption(t *testing.T) {
	assert := assert.New(t)

	message := "{\"foo\":{\"bar\":\"foobar\"}}"

	encrypted := EncryptString("enigma", message, false)
	assert.Equal("GsvkCYZoYylL5a7/DKhysDjNbwn+BtBtHj2CvzC4Y4g=", encrypted)
}

// TestUnicodeDecryption tests the Unicode decryption.
// Assumes that the input message is deserialized
// Decrypted string should match 漢語
func TestUnicodeDecryption(t *testing.T) {
	assert := assert.New(t)

	message := "+BY5/miAA8aeuhVl4d13Kg=="
	decrypted, decErr := DecryptString("enigma", message, false)
	assert.NoError(decErr)
	var msg interface{}
	json.Unmarshal([]byte(decrypted.(string)), &msg)
	assert.Equal("漢語", msg)
}

// UTF16ToString returns the UTF-8 encoding of the UTF-16 sequence s,
// with a terminating NUL removed.
func UTF16ToString(s []uint16) []rune {
	for i, v := range s {
		if v == 0 {
			s = s[0:i]
			break
		}
	}
	return utf16.Decode(s)
}

// TestUnicodeEncryption tests the Unicode encryption.
// Assumes that the input message is not serialized
// Decrypted string should match +BY5/miAA8aeuhVl4d13Kg==
func TestUnicodeEncryption(t *testing.T) {
	assert := assert.New(t)

	message := "漢語"
	b, err := json.Marshal(message)
	assert.NoError(err)

	encrypted := EncryptString("enigma", string(b), false)
	assert.Equal("+BY5/miAA8aeuhVl4d13Kg==", encrypted)
}

// TestGermanDecryption tests the German decryption.
// Assumes that the input message is deserialized
// Decrypted string should match ÜÖ
func TestGermanDecryption(t *testing.T) {
	assert := assert.New(t)

	message := "stpgsG1DZZxb44J7mFNSzg=="
	decrypted, decErr := DecryptString("enigma", message, false)
	assert.NoError(decErr)

	var msg interface{}

	json.Unmarshal([]byte(decrypted.(string)), &msg)
	assert.Equal("ÜÖ", msg)
}

// TestGermanEncryption tests the German encryption.
// Assumes that the input message is not serialized
// Decrypted string should match stpgsG1DZZxb44J7mFNSzg==
func TestGermanEncryption(t *testing.T) {
	assert := assert.New(t)

	message := "ÜÖ"
	b, err := json.Marshal(message)
	assert.NoError(err)

	encrypted := EncryptString("enigma", string(b), false)
	assert.Equal("stpgsG1DZZxb44J7mFNSzg==", encrypted)
}

// EmptyStruct provided the empty struct to test the encryption.
type emptyStruct struct {
}

// CustomStruct to test the custom structure encryption and decryption
// The variables "foo" and "bar" as used in the other languages are not
// accepted by golang and give an empty value when serialized, used "Foo"
// and "Bar" instead.
type customStruct struct {
	Foo string
	Bar []int
}

func CreateLoggerForTests() *log.Logger {
	var infoLogger *log.Logger
	logfileName := "pubnubMessagingTests.log"
	f, err := os.OpenFile(logfileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error opening file: ", err.Error())
		fmt.Println("Logging disabled")
	} else {
		//fmt.Println("Logging enabled writing to ", logfileName)
		infoLogger = log.New(f, "", log.Ldate|log.Ltime|log.Lshortfile)
	}
	return infoLogger
}
