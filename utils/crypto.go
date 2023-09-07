package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// 16 byte IV
var valIV = "0123456789012345"

// EncryptString DEPRECATED
// EncryptString creates the base64 encoded encrypted string using the
// cipherKey.
// It accepts the following parameters:
// cipherKey: cipher key to use to encrypt.
// message: to encrypted.
// useRandomInitializationVector: if true the IV is random and is sent along with the message
//
// returns the base64 encoded encrypted string.
func EncryptString(cipherKey string, message string, useRandomInitializationVector bool) string {
	cryptor, e := NewLegacyCryptor(cipherKey, useRandomInitializationVector)
	if e != nil {
		panic(e)
	}
	encryptedData, e := cryptor.Encrypt([]byte(encodeNonASCIIChars(message)))
	if e != nil {
		panic(e)
	}
	return base64.StdEncoding.EncodeToString(encryptedData.Data)
}

// DecryptString DEPRECATED
// DecryptString decodes encrypted string using the cipherKey
//
// It accepts the following parameters:
// cipherKey: cipher key to use to decrypt.
// message: to encrypted.
// useRandomInitializationVector: if true the IV is random and is prepended to the text. The IV is extracted and then the cipher is decoded.
//
// returns the unencoded encrypted string,
// error if any.
func DecryptString(cipherKey string, message string, useRandomInitializationVector bool) (retVal interface{}, err error) {
	value, decodeErr := base64.StdEncoding.DecodeString(message)
	if decodeErr != nil {
		return "***decrypt error***", fmt.Errorf("decrypt error on decode: %s", decodeErr)
	}

	cryptor, e := NewLegacyCryptor(cipherKey, useRandomInitializationVector)
	if e != nil {
		return nil, e
	}
	return cryptor.Decrypt(EncryptedData{Data: value, Metadata: nil})
}

// encodeNonAsciiChars creates unicode string of the non-ascii chars.
// It accepts the following parameters:
// message: to parse.
//
// returns the encoded string.
func encodeNonASCIIChars(message string) string {
	runeOfMessage := []rune(message)
	lenOfRune := len(runeOfMessage)
	encodedString := bytes.NewBuffer(make([]byte, 0, lenOfRune))
	for i := 0; i < lenOfRune; i++ {
		intOfRune := uint16(runeOfMessage[i])
		if intOfRune > 127 {
			hexOfRune := strconv.FormatUint(uint64(intOfRune), 16)
			dataLen := len(hexOfRune)
			paddingNum := 4 - dataLen
			encodedString.WriteString(`\u`)
			for i := 0; i < paddingNum; i++ {
				encodedString.WriteString("0")
			}
			encodedString.WriteString(hexOfRune)
		} else {
			encodedString.WriteString(string(runeOfMessage[i]))
		}
	}
	return encodedString.String()
}

// getHmacSha256 creates the cipher key hashed against SHA256.
// It accepts the following parameters:
// secretKey: the secret key.
// input: input to hash.
//
// returns the hash.
func GetHmacSha256(secretKey string, input string) string {
	hmacSha256 := hmac.New(sha256.New, []byte(secretKey))
	hmacSha256.Write([]byte(input))
	rawSig := base64.StdEncoding.EncodeToString(hmacSha256.Sum(nil))
	signature := strings.Replace(strings.Replace(rawSig, "+", "-", -1), "/", "_", -1)
	return signature
}

// EncryptFile DEPRECATED
func EncryptFile(cipherKey string, _ []byte, filePart io.Writer, file *os.File) {
	cryptor, e := NewLegacyCryptor(cipherKey, true)
	if e != nil {
		panic(e)
	}
	_, e = cryptor.EncryptStream(file, filePart)
	if e != nil {
		panic(e)
	}
}

// DecryptFile DEPRECATED
func DecryptFile(cipherKey string, _ int64, reader io.Reader, w io.WriteCloser) {
	cryptor, e := NewLegacyCryptor(cipherKey, true)
	if e != nil {
		panic(e)
	}
	e = cryptor.DecryptStream(reader, nil, w)
	if e != nil {
		panic(e)
	}
	e = w.Close()
}
