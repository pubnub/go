package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// 16 byte IV
var valIV = "0123456789012345"

// EncryptString creates the base64 encoded encrypted string using the
// cipherKey.
// It accepts the following parameters:
// cipherKey: cipher key to use to encrypt.
// message: to encrypted.
// useRandomInitializationVector: if true the IV is random and is sent along with the message
//
// returns the base64 encoded encrypted string.
func EncryptString(cipherKey string, message string, useRandomInitializationVector bool) string {
	block, _ := aesCipher(cipherKey)

	message = encodeNonASCIIChars(message)
	value := []byte(message)
	value = padWithPKCS7(value)
	iv := make([]byte, aes.BlockSize)
	if useRandomInitializationVector {
		iv = generateIV(aes.BlockSize)
	} else {
		iv = []byte(valIV)
	}
	blockmode := cipher.NewCBCEncrypter(block, iv)

	cipherBytes := make([]byte, len(value))
	blockmode.CryptBlocks(cipherBytes, value)
	if useRandomInitializationVector {
		return base64.StdEncoding.EncodeToString(append(iv, cipherBytes...))
	}
	return base64.StdEncoding.EncodeToString(cipherBytes)
}

type A struct {
	I         string
	Interface *B
}
type B struct {
	Value string
}

// DecryptString decodes encrypted string using the cipherKey
//
// It accepts the following parameters:
// cipherKey: cipher key to use to decrypt.
// message: to encrypted.
// useRandomInitializationVector: if true the IV is random and is prepended to the text. The IV is extracted and then the cipher is decoded.
//
// returns the unencoded encrypted string,
// error if any.
func DecryptString(cipherKey string, message string, useRandomInitializationVector bool) (
	retVal interface{}, err error) {
	if message == "" {
		return "**decrypt error***", errors.New("message is empty")
	}

	block, aesErr := aesCipher(cipherKey)
	if aesErr != nil {
		return "***decrypt error***", fmt.Errorf("decrypt error aes cipher: %s", aesErr)
	}

	value, decodeErr := base64.StdEncoding.DecodeString(message)
	if decodeErr != nil {
		return "***decrypt error***", fmt.Errorf("decrypt error on decode: %s", decodeErr)
	}
	iv := make([]byte, aes.BlockSize)
	if useRandomInitializationVector {
		iv = value[:16]
		value = value[16:]
	} else {
		iv = []byte(valIV)
	}

	decrypter := cipher.NewCBCDecrypter(block, iv)
	//to handle decryption errors
	defer func() {
		if r := recover(); r != nil {
			retVal, err = "***decrypt error***", fmt.Errorf("decrypt error: %s", r)
		}
	}()
	decrypted := make([]byte, len(value))
	decrypter.CryptBlocks(decrypted, value)
	val, err := unpadPKCS7(decrypted)
	if err != nil {
		return "***decrypt error***", fmt.Errorf("decrypt error: %s", err)
	}

	return fmt.Sprintf("%s", string(val)), nil
}

// aesCipher returns the cipher block
//
// It accepts the following parameters:
// cipherKey: cipher key.
//
// returns the cipher block,
// error if any.
func aesCipher(cipherKey string) (cipher.Block, error) {
	key := EncryptCipherKey(cipherKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return block, nil
}

// EncryptCipherKey creates the 256 bit hex of the cipher key
//
// It accepts the following parameters:
// cipherKey: cipher key to use to decrypt.
//
// returns the 256 bit hex of the cipher key.
func EncryptCipherKey(cipherKey string) []byte {
	hash := sha256.New()
	hash.Write([]byte(cipherKey))

	sha256String := hash.Sum(nil)[:16]
	return []byte(hex.EncodeToString(sha256String))
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

func generateIV(blocksize int) []byte {
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		panic(err)
	}
	return iv
}

func EncryptFile(cipherKey string, iv []byte, filePart io.Writer, file *os.File) {
	key := EncryptCipherKey(cipherKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	if bytes.Equal(iv, []byte{}) {
		iv = generateIV(aes.BlockSize)
	}
	_, e := filePart.Write(iv)
	if e != nil {
		panic(e)
	}
	blockSize := 16
	bufferSize := 16
	p := make([]byte, bufferSize)

	mode := cipher.NewCBCEncrypter(block, iv)
	cryptoRan := false
	fii, _ := file.Stat()
	contentLenIn := fii.Size()
	var contentRead int64

	for {
		n2, err2 := io.ReadFull(file, p)
		contentRead += int64(n2)
		if err2 != nil {
			if err2 == io.EOF {
				ciphertext := make([]byte, blockSize)
				copy(ciphertext[:n2], p[:n2])
				break
			}

			if err2 == io.ErrUnexpectedEOF {
				if !cryptoRan {
					text := make([]byte, blockSize)
					ciphertext := make([]byte, blockSize)
					copy(text[:n2], p[:n2])
					pad := bytes.Repeat([]byte{byte(blockSize - n2)}, blockSize-n2)
					copy(text[n2:], pad)
					mode.CryptBlocks(ciphertext, text)
					filePart.Write(ciphertext)
				} else {
					text := make([]byte, blockSize)
					ciphertext := make([]byte, blockSize)
					copy(text[:n2], p[:n2])
					pad := bytes.Repeat([]byte{byte(blockSize - n2)}, blockSize-n2)
					copy(text[n2:], pad)
					mode.CryptBlocks(ciphertext, text)
					filePart.Write(ciphertext)

				}

			}
			break
		}

		ciphertext := make([]byte, blockSize)
		cryptoRan = true
		if contentRead >= contentLenIn {
			pad := bytes.Repeat([]byte{byte(blockSize - n2)}, blockSize-n2)
			copy(p[n2:], pad)
		}

		mode.CryptBlocks(ciphertext, p)
		filePart.Write(ciphertext)

	}
}

func DecryptFile(cipherKey string, contentLenEnc int64, reader io.Reader, w io.WriteCloser) {
	key := EncryptCipherKey(cipherKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	blockSize := 16
	bufferSize := 16
	p := make([]byte, bufferSize)
	ivBuff := make([]byte, blockSize)
	emptyByteVar := make([]byte, blockSize)

	iv2 := make([]byte, blockSize)
	count := 0

	var mode cipher.BlockMode

	cryptoRan := false
	var contentDownloaded int64

	go func() {
	ExitReadLabel:
		for {
			n2, err2 := io.ReadFull(reader, p)
			if err2 != nil {
				if err2 == io.EOF {
					ciphertext := make([]byte, blockSize)
					copy(ciphertext, p[:n2])
					ciphertext, _ = unpadPKCS7(ciphertext)
					w.Write(ciphertext)
					w.Close()
					break ExitReadLabel
				}

				if err2 == io.ErrUnexpectedEOF {
					if bytes.Equal(iv2, emptyByteVar) {
						copy(iv2, ivBuff[0:blockSize])
						mode = cipher.NewCBCDecrypter(block, iv2)
					}
					if !cryptoRan {
						text := make([]byte, blockSize)
						ciphertext := make([]byte, blockSize)
						copy(text, p[:n2])
						mode.CryptBlocks(ciphertext, text)
						ciphertext, _ = unpadPKCS7(ciphertext)
						w.Write(ciphertext)

						w.Close()
						break ExitReadLabel
					} else {
						ciphertext := make([]byte, blockSize)
						copy(ciphertext, p[:n2])
						ciphertext, _ = unpadPKCS7(ciphertext)
						w.Write(ciphertext)
						w.Close()
						break ExitReadLabel
					}

				}
				break ExitReadLabel
			} else {
				contentDownloaded += int64(n2)
				if count < blockSize/bufferSize {
					if err != nil {
						panic(err)
					}
					copy(ivBuff[bufferSize*count:], p)
				} else {

					if bytes.Equal(iv2, emptyByteVar) {
						copy(iv2, ivBuff[0:blockSize])
						mode = cipher.NewCBCDecrypter(block, iv2)
					}

					ciphertext := make([]byte, blockSize)

					text := make([]byte, blockSize)
					copy(text, p[:n2])

					mode.CryptBlocks(ciphertext, p)
					cryptoRan = true
					if contentDownloaded >= contentLenEnc {
						ciphertext, _ = unpadPKCS7(ciphertext)
						w.Write(ciphertext)
					} else {
						w.Write(ciphertext)
					}
				}
			}
			count++

		}
	}()
}
