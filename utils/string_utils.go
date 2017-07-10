package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
)

func ChannelsAsString(channels []string) ([]byte, error) {
	// TODO: channels should be optionally encoded
	return []byte(strings.Join(channels, ",")), nil
}

// PubNub - specific serializer
func ValueAsString(value interface{}) ([]byte, error) {
	switch t := value.(type) {
	case string:
		return []byte(fmt.Sprintf("\"%s\"", t)), nil
	default:
		val, err := json.Marshal(value)
		fmt.Printf("Marshaled %s to %s\n", value, val)
		return val, err
	}
}

// Generate a random uuid string
func Uuid() string {
	return uuid.NewV4().String()
}
