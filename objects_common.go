package pubnub

// PNUser is the response to createUser request. It contains a map of type PNUserItem
type PNUser struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	ExternalID string                 `json:"externalId"`
	ProfileURL string                 `json:"profileUrl"`
	Email      string                 `json:"email"`
	Created    string                 `json:"created"`
	Updated    string                 `json:"updated"`
	ETag       string                 `json:"eTag"`
	Custom     map[string]interface{} `json:"custom"`
}

type PNSpace struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Created     string                 `json:"created"`
	Updated     string                 `json:"updated"`
	ETag        string                 `json:"eTag"`
	Custom      map[string]interface{} `json:"custom"`
}

type PNMembers struct {
	ID      string                 `json:"id"`
	User    PNUser                 `json:"user"`
	Created string                 `json:"created"`
	Updated string                 `json:"updated"`
	ETag    string                 `json:"eTag"`
	Custom  map[string]interface{} `json:"custom"`
}

type PNMemberships struct {
	ID      string                 `json:"id"`
	Space   PNSpace                `json:"space"`
	Created string                 `json:"created"`
	Updated string                 `json:"updated"`
	ETag    string                 `json:"eTag"`
	Custom  map[string]interface{} `json:"custom"`
}

type PNMembersInput struct {
	ID     string                 `json:"id"`
	Custom map[string]interface{} `json:"custom"`
}

type PNMembersRemove struct {
	ID string `json:"id"`
}

type PNMembershipsInput struct {
	ID     string                 `json:"id"`
	Custom map[string]interface{} `json:"custom"`
}

type PNMembershipsRemove struct {
	ID string `json:"id"`
}

type PNObjectsResponse struct {
	Event       PNObjectsEvent         `json:"event"` // enum value
	EventType   PNObjectsEventType     `json:"type"`  // enum value
	Name        string                 `json:"name"`
	UserID      string                 `json:"userId"`      // the user id if user related
	SpaceID     string                 `json:"spaceId"`     // the space id if space related
	Description string                 `json:"description"` // the description of what happened
	Timestamp   string                 `json:"timestamp"`   // the timetoken of the event
	ExternalID  string                 `json:"externalId"`
	ProfileURL  string                 `json:"profileUrl"`
	Email       string                 `json:"email"`
	Created     string                 `json:"created"`
	Updated     string                 `json:"updated"`
	ETag        string                 `json:"eTag"`
	Custom      map[string]interface{} `json:"custom"`
	Data        map[string]interface{} `json:"data"`
}
