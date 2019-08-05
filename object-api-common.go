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

type PNUserMembership struct {
	Id      string                 `json:"id"`
	Name    string                 `json:"name"`
	Space   PNSpace                `json:"space"`
	Created string                 `json:"created"`
	Updated string                 `json:"updated"`
	ETag    string                 `json:"eTag"`
	Custom  map[string]interface{} `json:"custom"`
}

type PNSpaceMembership struct {
	Id      string                 `json:"id"`
	Name    string                 `json:"name"`
	User    PNUser                 `json:"user"`
	Created string                 `json:"created"`
	Updated string                 `json:"updated"`
	ETag    string                 `json:"eTag"`
	Custom  map[string]interface{} `json:"custom"`
}

type PNUserMembershipInput struct {
	Id     string                 `json:"id"`
	Custom map[string]interface{} `json:"custom"`
}

type PNUserMembershipRemove struct {
	Id string `json:"id"`
}

type PNSpaceMembershipInput struct {
	Id     string                 `json:"id"`
	Custom map[string]interface{} `json:"custom"`
}

type PNSpaceMembershipRemove struct {
	Id string `json:"id"`
}
