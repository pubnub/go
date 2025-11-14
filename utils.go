package pubnub

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/pubnub/go/v8/crypto"
	"github.com/pubnub/go/v8/pnerr"
)

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

func encryptString(module crypto.CryptoModule, message string) (string, error) {
	encryptedData, e := module.Encrypt([]byte(encodeNonASCIIChars(message)))
	if e != nil {
		return "", e
	}
	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

func serializeEncryptAndSerialize(cryptoModule crypto.CryptoModule, msg interface{}, serialize bool, loggerMgr *loggerManager) (string, error) {
	var encrypted string
	var err error

	if serialize {
		if loggerMgr != nil {
			loggerMgr.LogSimple(PNLogLevelTrace, "Serialization: serializing message content", false)
		}
		jsonSerialized, errJSONMarshal := json.Marshal(msg)
		if errJSONMarshal != nil {
			return "", errJSONMarshal
		}
		if loggerMgr != nil {
			loggerMgr.LogSimple(PNLogLevelTrace, "Serialization: message serialized successfully", false)
		}
		encrypted, err = encryptString(cryptoModule, string(jsonSerialized))

	} else {
		if serializedMsg, ok := msg.(string); ok {
			encrypted, err = encryptString(cryptoModule, string(serializedMsg))
		} else {
			return "", pnerr.NewBuildRequestError("Message is not JSON serialized.")
		}
	}
	if err != nil {
		return "", err
	}
	if loggerMgr != nil {
		loggerMgr.LogSimple(PNLogLevelTrace, "Serialization: serializing encrypted content", false)
	}
	jsonSerialized, errJSONMarshal := json.Marshal(encrypted)
	if errJSONMarshal != nil {
		return "", errJSONMarshal
	}
	if loggerMgr != nil {
		loggerMgr.LogSimple(PNLogLevelTrace, "Serialization: encrypted content serialized successfully", false)
	}
	return string(jsonSerialized), nil
}

func serializeAndEncrypt(cryptoModule crypto.CryptoModule, msg interface{}, serialize bool, loggerMgr *loggerManager) (string, error) {
	var encrypted string
	var err error
	if serialize {
		if loggerMgr != nil {
			loggerMgr.LogSimple(PNLogLevelTrace, "Serialization: serializing message content", false)
		}
		jsonSerialized, errJSONMarshal := json.Marshal(msg)
		if errJSONMarshal != nil {
			return "", errJSONMarshal
		}
		if loggerMgr != nil {
			loggerMgr.LogSimple(PNLogLevelTrace, "Serialization: message serialized successfully", false)
		}
		encrypted, err = encryptString(cryptoModule, string(jsonSerialized))
	} else {
		if serializedMsg, ok := msg.(string); ok {
			encrypted, err = encryptString(cryptoModule, serializedMsg)
		} else {
			return "", pnerr.NewBuildRequestError("Message is not JSON serialized.")
		}
	}
	if err != nil {
		return "", err
	}

	return encrypted, nil
}

func encryptStreamAndCopyTo(module crypto.CryptoModule, reader io.Reader, writer io.Writer) error {
	encryptedStream, e := module.EncryptStream(reader)
	if e != nil {
		return e
	}
	_, e = io.Copy(writer, encryptedStream)
	if e != nil {
		return e
	}
	return nil
}

func decryptString(cryptoModule crypto.CryptoModule, message string) (retVal interface{}, err error) {
	value, decodeErr := base64.StdEncoding.DecodeString(message)
	if decodeErr != nil {
		return "***decrypt error***", fmt.Errorf("decrypt error on decode: %s", decodeErr)
	}

	val, e := cryptoModule.Decrypt(value)
	return string(val), e
}

// unmarshalWithLogging wraps json.Unmarshal with trace-level logging
func unmarshalWithLogging(data []byte, v interface{}, loggerMgr *loggerManager, operation string) error {
	if loggerMgr != nil {
		loggerMgr.LogSimple(PNLogLevelTrace, fmt.Sprintf("Deserialization: deserializing %s response", operation), false)
	}
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	if loggerMgr != nil {
		loggerMgr.LogSimple(PNLogLevelTrace, fmt.Sprintf("Deserialization: %s response deserialized successfully", operation), false)
	}
	return nil
}

func isCustomMessageTypeValid(customMessageType string) bool {
	if len(customMessageType) == 0 {
		return true
	}

	if len(customMessageType) < 3 || len(customMessageType) > 50 {
		return false
	}

	for _, c := range customMessageType {
		if !('a' <= c && 'z' >= c) && !('A' <= c && 'Z' >= c) && c != '-' && c != '_' {
			return false
		}
	}

	return true
}
