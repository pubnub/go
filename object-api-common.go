package pubnub

// PNUser is the response to createUser request. It contains a map of type PNUserItem
type PNUser struct {
	Id         string                 `json:"id"`
	Name       string                 `json:"name"`
	ExternalId string                 `json:"externalId"`
	ProfileUrl string                 `json:"profileUrl"`
	Email      string                 `json:"email"`
	Created    string                 `json:"created"`
	Updated    string                 `json:"updated"`
	ETag       string                 `json:"eTag"`
	Custom     map[string]interface{} `json:"custom"`
}

type PNSpace struct {
	Id          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Created     string                 `json:"created"`
	Updated     string                 `json:"updated"`
	ETag        string                 `json:"eTag"`
	Custom      map[string]interface{} `json:"custom"`
}
