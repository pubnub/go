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
	assert.Equal(m.Text, "")
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
