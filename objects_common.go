package pubnub

// PNUser is the Objects API user struct
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

// PNSpace is the Objects API space struct
type PNSpace struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Created     string                 `json:"created"`
	Updated     string                 `json:"updated"`
	ETag        string                 `json:"eTag"`
	Custom      map[string]interface{} `json:"custom"`
}

// PNMembers is the Objects API Members struct
type PNMembers struct {
	ID      string                 `json:"id"`
	User    PNUser                 `json:"user"`
	Created string                 `json:"created"`
	Updated string                 `json:"updated"`
	ETag    string                 `json:"eTag"`
	Custom  map[string]interface{} `json:"custom"`
}

// PNMemberships is the Objects API Memberships struct
type PNMemberships struct {
	ID      string                 `json:"id"`
	Space   PNSpace                `json:"space"`
	Created string                 `json:"created"`
	Updated string                 `json:"updated"`
	ETag    string                 `json:"eTag"`
	Custom  map[string]interface{} `json:"custom"`
}

// PNMembersInput is the Objects API Members input struct used to add members
type PNMembersInput struct {
	ID     string                 `json:"id"`
	Custom map[string]interface{} `json:"custom"`
}

// PNMembersRemove is the Objects API Members struct used to remove members
type PNMembersRemove struct {
	ID string `json:"id"`
}

// PNMembershipsInput is the Objects API Memberships input struct used to add members
type PNMembershipsInput struct {
	ID     string                 `json:"id"`
	Custom map[string]interface{} `json:"custom"`
}

// PNMembershipsRemove is the Objects API Memberships struct used to remove members
type PNMembershipsRemove struct {
	ID string `json:"id"`
}

// PNObjectsResponse is the Objects API collective Response struct of all methods.
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
