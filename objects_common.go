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

type PNMembers struct {
	Id      string                 `json:"id"`
	User    PNUser                 `json:"user"`
	Created string                 `json:"created"`
	Updated string                 `json:"updated"`
	ETag    string                 `json:"eTag"`
	Custom  map[string]interface{} `json:"custom"`
}

type PNMemberships struct {
	Id      string                 `json:"id"`
	Space   PNSpace                `json:"space"`
	Created string                 `json:"created"`
	Updated string                 `json:"updated"`
	ETag    string                 `json:"eTag"`
	Custom  map[string]interface{} `json:"custom"`
}

type PNMembersInput struct {
	Id     string                 `json:"id"`
	Custom map[string]interface{} `json:"custom"`
}

type PNMembersRemove struct {
	Id string `json:"id"`
}

type PNMembershipsInput struct {
	Id     string                 `json:"id"`
	Custom map[string]interface{} `json:"custom"`
}

type PNMembershipsRemove struct {
	Id string `json:"id"`
}

type PNObjectsResponse struct {
	Event       PNObjectsEvent         `json:"event"` // enum value
	EventType   PNObjectsEventType     `json:"type"`  // enum value
	Name        string                 `json:"name"`
	UserId      string                 `json:"userId"`      // the user id if user related
	SpaceId     string                 `json:"spaceId"`     // the space id if space related
	Description string                 `json:"description"` // the description of what happened
	Timestamp   string                 `json:"timestamp"`   // the timetoken of the event
	ExternalId  string                 `json:"externalId"`
	ProfileUrl  string                 `json:"profileUrl"`
	Email       string                 `json:"email"`
	Created     string                 `json:"created"`
	Updated     string                 `json:"updated"`
	ETag        string                 `json:"eTag"`
	Custom      map[string]interface{} `json:"custom"`
	Data        map[string]interface{} `json:"data"`
}
