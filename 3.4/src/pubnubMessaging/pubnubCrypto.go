package pubnubMessaging

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/hmac"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "fmt"
    "io"
    "bytes"
    "strconv"
)

var _IV = "0123456789012345"

func PKCS7Padding(data []byte) []byte {
    dataLen := len(data)
    var bit16 int
    if dataLen%16 == 0 {
        bit16 = dataLen
    } else {
        bit16 = int(dataLen/16+1) * 16
    }

    paddingNum := bit16 - dataLen
    bitCode := byte(paddingNum)

    padding := make([]byte, paddingNum)
    for i := 0; i < paddingNum; i++ {
        padding[i] = bitCode
    }
    return append(data, padding...)
}

func UnPKCS7Padding(data []byte) []byte {
    dataLen := len(data)
    if dataLen == 0 {
        return data
    }
    endIndex := int(data[dataLen-1])
    if 16 > endIndex {
        if 1 < endIndex {
            for i := dataLen - endIndex; i < dataLen; i++ {
                if data[dataLen-1] != data[i] {
                    fmt.Println(" : ", data[dataLen-1], " ：", i, "  ：", data[i])
                }
            }
        }
        return data[:dataLen-endIndex]
    }
    return data
}

func GetHmacSha256(secretKey string, input string) string {
    hmacSha256 := hmac.New(sha256.New, []byte(secretKey))
    io.WriteString(hmacSha256, input)
    
    return fmt.Sprintf("%x", hmacSha256.Sum(nil))
}

func GenUuid() (string, error) {
    uuid := make([]byte, 16)
    n, err := rand.Read(uuid)
    if n != len(uuid) || err != nil {
        return "", err
    }
    // TODO: verify the two lines implement RFC 4122 correctly
    uuid[8] = 0x80 // variant bits see page 5
    uuid[4] = 0x40 // version 4 Pseudo Random, see page 7

    return hex.EncodeToString(uuid), nil
}

func EncodeNonAsciiChars(message string) string {
    runeOfMessage := []rune(message)
    lenOfRune := len(runeOfMessage)
    encodedString := ""    
    for i := 0; i < lenOfRune; i++ {
        intOfRune := uint16(runeOfMessage[i])
        if(intOfRune>127){
            hexOfRune := strconv.FormatUint(uint64(intOfRune), 16)
            dataLen := len(hexOfRune)
            paddingNum := 4 - dataLen
            prefix := ""
            for i := 0; i < paddingNum; i++ {
                prefix += "0"
            }
            hexOfRune = prefix + hexOfRune
            encodedString += bytes.NewBufferString(`\u` + hexOfRune).String()
        } else {
            encodedString += string(runeOfMessage[i])
        }
    }
    return encodedString
}

func EncryptString(cipherKey string, message string) string {
    block, _ := AesCipher(cipherKey)
    message = EncodeNonAsciiChars(message)
    value := []byte(message)
    value = PKCS7Padding(value)
    blockmode := cipher.NewCBCEncrypter(block, []byte(_IV))
    cipherBytes := make([]byte, len(value))
    blockmode.CryptBlocks(cipherBytes, value)
    return fmt.Sprintf("%s", Encode(cipherBytes))
}

func DecryptString(cipherKey string, message string) string { //need add error catching
    block, _ := AesCipher(cipherKey)
    value, _ := base64.StdEncoding.DecodeString(message)
    
    decrypter := cipher.NewCBCDecrypter(block, []byte(_IV))
    decrypted := make([]byte, len(value))
    decrypter.CryptBlocks(decrypted, value)
    return fmt.Sprintf("%s", string(UnPKCS7Padding(decrypted)))
}

func AesCipher(cipherKey string) (cipher.Block, error) {
    key := EncryptCipherKey(cipherKey)
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    return block, nil
}

func EncryptCipherKey(cipherKey string) []byte {
    hash := sha256.New()
    hash.Write([]byte(cipherKey))

    sha256String := hash.Sum(nil)[:16]
    return []byte(hex.EncodeToString(sha256String))
}

//Encodes a value using base64
func Encode(value []byte) []byte {
    encoded := make([]byte, base64.StdEncoding.EncodedLen(len(value)))
    base64.StdEncoding.Encode(encoded, value)
    return encoded
}

//Decodes a value using base64 
func Decode(value []byte) ([]byte, error) {
    decoded := make([]byte, base64.StdEncoding.DecodedLen(len(value)))
    b, err := base64.StdEncoding.Decode(decoded, value)
    if err != nil {
        return nil, err
    }
    return decoded[:b], nil
}
