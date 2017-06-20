package utils

import (
	"encoding/json"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// PubNub - specific serializer
func ValueAsString(value interface{}) ([]byte, error) {
	switch t := value.(type) {
	case string:
		return []byte(fmt.Sprintf("\"%s\"", t)), nil
	default:
		val, err := json.Marshal(value)
		return val, err
	}
}

// Generate a random uuid string
func Uuid() string {
	return uuid.NewV4().String()
}
