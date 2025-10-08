package pubnub

// PNPublishMessage is the part of the message struct used in Publish File.
// Text can be any JSON-serializable type: string, map[string]interface{}, []interface{}, number, bool, etc.
type PNPublishMessage struct {
	Text interface{} `json:"text"`
}

// PNPublishMessageRaw is used when UseRawMessage is true - the message content is sent directly without "text" wrapper.
// Text can be any JSON-serializable type: string, map[string]interface{}, []interface{}, number, bool, etc.
type PNPublishMessageRaw struct {
	Text interface{} `json:"-"`
}

// PNFileInfoForPublish is the part of the message struct used in Publish File
type PNFileInfoForPublish struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// PNPublishFileMessage is the message struct used in Publish File
type PNPublishFileMessage struct {
	PNMessage *PNPublishMessage     `json:"message"`
	PNFile    *PNFileInfoForPublish `json:"file"`
}

// PNPublishFileMessageRaw is used when UseRawMessage is true - the message is sent as raw content without "text" wrapper
type PNPublishFileMessageRaw struct {
	PNMessage *PNPublishMessageRaw  `json:"message"`
	PNFile    *PNFileInfoForPublish `json:"file"`
}

// PNFileInfo is the File Upload API struct returned on for each file.
type PNFileInfo struct {
	Name    string `json:"name"`
	ID      string `json:"id"`
	Size    int    `json:"size"`
	Created string `json:"created"`
}

// PNFileData is used in the responses to show File ID
type PNFileData struct {
	ID string `json:"id"`
}

// PNFileUploadRequest is used to store the info related to file upload to S3
type PNFileUploadRequest struct {
	URL        string        `json:"url"`
	Method     string        `json:"method"`
	FormFields []PNFormField `json:"form_fields"`
}

// PNFormField is part of the struct used in file upload to S3
type PNFormField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// PNFileDetails is used in the responses to show File Info
type PNFileDetails struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	URL  string
}

// PNFileMessageAndDetails is used to store the file message and file info
type PNFileMessageAndDetails struct {
	PNMessage PNPublishMessage `json:"message"`
	PNFile    PNFileDetails    `json:"file"`
}

// ParseFileInfo extracts file information and message content from PubNub file payloads.
// It handles multiple message formats:
//   - Regular format with "text" wrapper: {"message": {"text": <value>}} - extracts value from "text" field
//   - Raw format without wrapper: {"message": <any JSON type>} - uses the message value directly
//   - JSON objects without "text" field are treated as raw format
//
// Returns PNFileDetails containing file metadata and PNPublishMessage with the message content.
func ParseFileInfo(filesPayload map[string]interface{}) (PNFileDetails, PNPublishMessage) {
	resp := &PNFileMessageAndDetails{}
	resp.PNMessage = PNPublishMessage{}
	resp.PNFile = PNFileDetails{}

	//"message":{"text":"test file"},"file":{"name":"test_file_upload_name_32899","id":"9076246e-5036-42af-b3a3-767b514c93c8"}}
	if o, ok := filesPayload["file"]; ok {
		if o != nil {
			if data, ok := o.(map[string]interface{}); ok {
				if d, ok := data["id"]; ok {
					if idStr, ok := d.(string); ok {
						resp.PNFile.ID = idStr
					}
				}
				if d, ok := data["name"]; ok {
					if nameStr, ok := d.(string); ok {
						resp.PNFile.Name = nameStr
					}
				}
			}
		}
	}
	if m, ok := filesPayload["message"]; ok {
		if m != nil {
			// Handle multiple message formats
			if data, ok := m.(map[string]interface{}); ok {
				// Message is a JSON object - check if it has "text" field
				if d, ok := data["text"]; ok {
					// Format: {"message": {"text": <value>}} - extract value from "text" field
					resp.PNMessage.Text = d
				} else {
					// Format: {"message": {"key": "value", ...}} - use entire object as message
					resp.PNMessage.Text = m
				}
			} else {
				// Format: {"message": <primitive>} - use value directly (string, number, bool, array)
				resp.PNMessage.Text = m
			}
		}
	}
	return resp.PNFile, resp.PNMessage
}
