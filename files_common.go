package pubnub

// PNPublishMessage is the part of the message struct used in Publish File
type PNPublishMessage struct {
	Text string `json:"text"`
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

// ParseFileInfo is a function extract file info and add to the struct PNFileMessageAndDetails
func ParseFileInfo(filesPayload map[string]interface{}) (PNFileDetails, PNPublishMessage) {
	var data map[string]interface{}
	resp := &PNFileMessageAndDetails{}
	resp.PNMessage = PNPublishMessage{}
	resp.PNFile = PNFileDetails{}

	//"message":{"text":"test file"},"file":{"name":"test_file_upload_name_32899","id":"9076246e-5036-42af-b3a3-767b514c93c8"}}
	if o, ok := filesPayload["file"]; ok {
		data = o.(map[string]interface{})
		if d, ok := data["id"]; ok {
			resp.PNFile.ID = d.(string)
		}
		if d, ok := data["name"]; ok {
			resp.PNFile.Name = d.(string)
		}
	}
	if m, ok := filesPayload["message"]; ok {
		data = m.(map[string]interface{})
		if d, ok := data["text"]; ok {
			resp.PNMessage.Text = d.(string)
		}
	}
	return resp.PNFile, resp.PNMessage
}
