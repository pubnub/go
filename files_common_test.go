package pubnub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFileInfo(t *testing.T) {
	assert := assert.New(t)
	resp := make(map[string]interface{})
	resp["message"] = nil
	resp["file"] = map[string]interface{}{"name": "test_file_upload_name_32899", "id": "9076246e-5036-42af-b3a3-767b514c93c8"}
	f, m := ParseFileInfo(resp)
	assert.Equal(f.ID, "9076246e-5036-42af-b3a3-767b514c93c8")
	assert.Nil(m.Text, "Text should be nil when message is nil")
}

func TestParseFileInfoNotNil(t *testing.T) {
	assert := assert.New(t)
	resp := make(map[string]interface{})
	resp["message"] = map[string]interface{}{"text": "test file"}
	resp["file"] = map[string]interface{}{"name": "test_file_upload_name_32899", "id": "9076246e-5036-42af-b3a3-767b514c93c8"}

	f, m := ParseFileInfo(resp)
	assert.Equal(f.ID, "9076246e-5036-42af-b3a3-767b514c93c8")
	assert.Equal(m.Text, "test file")
}

// TestParseFileInfoWithRawStringMessage tests parsing file info when message is a raw string
func TestParseFileInfoWithRawStringMessage(t *testing.T) {
	assert := assert.New(t)
	resp := make(map[string]interface{})
	resp["message"] = "raw message"
	resp["file"] = map[string]interface{}{"name": "test_file.txt", "id": "file-id-123"}

	f, m := ParseFileInfo(resp)
	assert.Equal("file-id-123", f.ID)
	assert.Equal("test_file.txt", f.Name)
	assert.Equal("raw message", m.Text)
}

// TestParseFileInfoWithJSONObjectMessage tests parsing when message is a JSON object
func TestParseFileInfoWithJSONObjectMessage(t *testing.T) {
	assert := assert.New(t)
	resp := make(map[string]interface{})
	resp["message"] = map[string]interface{}{
		"type":     "document",
		"priority": "high",
		"metadata": map[string]interface{}{
			"author": "John Doe",
			"tags":   []string{"important", "quarterly"},
		},
	}
	resp["file"] = map[string]interface{}{"name": "document.pdf", "id": "file-id-456"}

	f, m := ParseFileInfo(resp)
	assert.Equal("file-id-456", f.ID)
	assert.Equal("document.pdf", f.Name)

	// Verify the message is correctly parsed as a map
	messageMap, ok := m.Text.(map[string]interface{})
	assert.True(ok, "Message should be a map[string]interface{}")
	assert.Equal("document", messageMap["type"])
	assert.Equal("high", messageMap["priority"])

	metadata, ok := messageMap["metadata"].(map[string]interface{})
	assert.True(ok)
	assert.Equal("John Doe", metadata["author"])
}

// TestParseFileInfoWithJSONObjectTextWrapper tests parsing when message has "text" wrapper with JSON object
func TestParseFileInfoWithJSONObjectTextWrapper(t *testing.T) {
	assert := assert.New(t)
	resp := make(map[string]interface{})
	resp["message"] = map[string]interface{}{
		"text": map[string]interface{}{
			"user_id": float64(123), // Use float64 to simulate JSON unmarshaling behavior
			"action":  "upload",
		},
	}
	resp["file"] = map[string]interface{}{"name": "image.jpg", "id": "file-id-789"}

	f, m := ParseFileInfo(resp)
	assert.Equal("file-id-789", f.ID)
	assert.Equal("image.jpg", f.Name)

	// Verify the message is correctly extracted from "text" wrapper
	messageMap, ok := m.Text.(map[string]interface{})
	assert.True(ok, "Message should be a map[string]interface{}")

	// JSON unmarshaling converts numbers to float64
	userID, ok := messageMap["user_id"].(float64)
	assert.True(ok, "user_id should be float64")
	assert.Equal(float64(123), userID)
	assert.Equal("upload", messageMap["action"])
}

// TestParseFileInfoWithArrayMessage tests parsing when message is an array
func TestParseFileInfoWithArrayMessage(t *testing.T) {
	assert := assert.New(t)
	resp := make(map[string]interface{})
	resp["message"] = []interface{}{"item1", "item2", "item3"}
	resp["file"] = map[string]interface{}{"name": "list.txt", "id": "file-id-999"}

	f, m := ParseFileInfo(resp)
	assert.Equal("file-id-999", f.ID)
	assert.Equal("list.txt", f.Name)

	// Verify the message is correctly parsed as an array
	messageArray, ok := m.Text.([]interface{})
	assert.True(ok, "Message should be a []interface{}")
	assert.Equal(3, len(messageArray))
	assert.Equal("item1", messageArray[0])
	assert.Equal("item2", messageArray[1])
	assert.Equal("item3", messageArray[2])
}

// TestParseFileInfoWithNumericMessage tests parsing when message is a number
func TestParseFileInfoWithNumericMessage(t *testing.T) {
	assert := assert.New(t)
	resp := make(map[string]interface{})
	resp["message"] = float64(42)
	resp["file"] = map[string]interface{}{"name": "number.txt", "id": "file-id-111"}

	f, m := ParseFileInfo(resp)
	assert.Equal("file-id-111", f.ID)
	assert.Equal("number.txt", f.Name)

	// Verify the message is correctly parsed as a number
	messageNumber, ok := m.Text.(float64)
	assert.True(ok, "Message should be a float64")
	assert.Equal(float64(42), messageNumber)
}

// TestParseFileInfoWithBooleanMessage tests parsing when message is a boolean
func TestParseFileInfoWithBooleanMessage(t *testing.T) {
	assert := assert.New(t)
	resp := make(map[string]interface{})
	resp["message"] = true
	resp["file"] = map[string]interface{}{"name": "bool.txt", "id": "file-id-222"}

	f, m := ParseFileInfo(resp)
	assert.Equal("file-id-222", f.ID)
	assert.Equal("bool.txt", f.Name)

	// Verify the message is correctly parsed as a boolean
	messageBool, ok := m.Text.(bool)
	assert.True(ok, "Message should be a bool")
	assert.Equal(true, messageBool)
}

// TestParseFileInfoWithEmptyMap tests parsing when message is an empty JSON object
func TestParseFileInfoWithEmptyMap(t *testing.T) {
	assert := assert.New(t)
	resp := make(map[string]interface{})
	resp["message"] = map[string]interface{}{} // Empty map
	resp["file"] = map[string]interface{}{"name": "empty.json", "id": "file-id-444"}

	f, m := ParseFileInfo(resp)
	assert.Equal("file-id-444", f.ID)
	assert.Equal("empty.json", f.Name)

	// Empty map without "text" field should be treated as raw format (entire empty object)
	messageMap, ok := m.Text.(map[string]interface{})
	assert.True(ok, "Message should be a map[string]interface{}")
	assert.Equal(0, len(messageMap), "Empty map should be preserved")
}

// TestParseFileInfoWithComplexNestedMessage tests parsing with deeply nested structures
func TestParseFileInfoWithComplexNestedMessage(t *testing.T) {
	assert := assert.New(t)
	resp := make(map[string]interface{})
	resp["message"] = map[string]interface{}{
		"user": map[string]interface{}{
			"id":   123,
			"name": "Test User",
			"settings": map[string]interface{}{
				"theme":         "dark",
				"notifications": true,
			},
		},
		"file_info": map[string]interface{}{
			"uploaded_at": "2023-01-01T00:00:00Z",
			"size":        1024,
			"tags":        []interface{}{"important", "document"},
		},
	}
	resp["file"] = map[string]interface{}{"name": "complex.json", "id": "file-id-333"}

	f, m := ParseFileInfo(resp)
	assert.Equal("file-id-333", f.ID)
	assert.Equal("complex.json", f.Name)

	// Verify the complex nested structure is preserved
	messageMap, ok := m.Text.(map[string]interface{})
	assert.True(ok, "Message should be a map[string]interface{}")

	user, ok := messageMap["user"].(map[string]interface{})
	assert.True(ok)
	assert.Equal("Test User", user["name"])

	settings, ok := user["settings"].(map[string]interface{})
	assert.True(ok)
	assert.Equal("dark", settings["theme"])
	assert.Equal(true, settings["notifications"])

	fileInfo, ok := messageMap["file_info"].(map[string]interface{})
	assert.True(ok)

	tags, ok := fileInfo["tags"].([]interface{})
	assert.True(ok)
	assert.Equal(2, len(tags))
	assert.Equal("important", tags[0])
	assert.Equal("document", tags[1])
}
