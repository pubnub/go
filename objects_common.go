package pubnub

// PNUUID is the Objects API user struct
type PNUUID struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	ExternalID string                 `json:"externalId"`
	ProfileURL string                 `json:"profileUrl"`
	Email      string                 `json:"email"`
	Updated    string                 `json:"updated"`
	ETag       string                 `json:"eTag"`
	Custom     map[string]interface{} `json:"custom"`
}

// PNChannel is the Objects API space struct
type PNChannel struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Updated     string                 `json:"updated"`
	ETag        string                 `json:"eTag"`
	Custom      map[string]interface{} `json:"custom"`
}

// PNChannelMembers is the Objects API Members struct
type PNChannelMembers struct {
	ID      string                 `json:"id"`
	UUID    PNUUID                 `json:"uuid"`
	Created string                 `json:"created"`
	Updated string                 `json:"updated"`
	ETag    string                 `json:"eTag"`
	Custom  map[string]interface{} `json:"custom"`
}

// PNMemberships is the Objects API Memberships struct
type PNMemberships struct {
	ID      string                 `json:"id"`
	Channel PNChannel              `json:"channel"`
	Created string                 `json:"created"`
	Updated string                 `json:"updated"`
	ETag    string                 `json:"eTag"`
	Custom  map[string]interface{} `json:"custom"`
}

// PNChannelMembersUUID is the Objects API Members input struct used to add members
type PNChannelMembersUUID struct {
	ID string `json:"id"`
}

// PNChannelMembersSet is the Objects API Members input struct used to add members
type PNChannelMembersSet struct {
	UUID   PNChannelMembersUUID   `json:"uuid"`
	Custom map[string]interface{} `json:"custom"`
}

// PNChannelMembersRemove is the Objects API Members struct used to remove members
type PNChannelMembersRemove struct {
	UUID PNChannelMembersUUID `json:"uuid"`
}

// PNMembershipsChannel is the Objects API Memberships input struct used to add members
type PNMembershipsChannel struct {
	ID string `json:"id"`
}

// PNMembershipsSet is the Objects API Memberships input struct used to add members
type PNMembershipsSet struct {
	Channel PNMembershipsChannel   `json:"channel"`
	Custom  map[string]interface{} `json:"custom"`
}

// PNMembershipsRemove is the Objects API Memberships struct used to remove members
type PNMembershipsRemove struct {
	Channel PNMembershipsChannel `json:"channel"`
}

// PNObjectsResponse is the Objects API collective Response struct of all methods.
type PNObjectsResponse struct {
	Event       PNObjectsEvent         `json:"event"` // enum value
	EventType   PNObjectsEventType     `json:"type"`  // enum value
	Name        string                 `json:"name"`
	ID          string                 `json:"id"`          // the uuid if user related
	Channel     string                 `json:"channel"`     // the channel if space related
	Description string                 `json:"description"` // the description of what happened
	Timestamp   string                 `json:"timestamp"`   // the timetoken of the event
	ExternalID  string                 `json:"externalId"`
	ProfileURL  string                 `json:"profileUrl"`
	Email       string                 `json:"email"`
	Updated     string                 `json:"updated"`
	ETag        string                 `json:"eTag"`
	Custom      map[string]interface{} `json:"custom"`
	Data        map[string]interface{} `json:"data"`
}

// PNManageMembershipsBody is the Objects API input to add, remove or update membership
type PNManageMembershipsBody struct {
	Set    []PNMembershipsSet    `json:"set"`
	Remove []PNMembershipsRemove `json:"delete"`
}

// PNManageChannelMembersBody is the Objects API input to add, remove or update members
type PNManageChannelMembersBody struct {
	Set    []PNChannelMembersSet    `json:"set"`
	Remove []PNChannelMembersRemove `json:"delete"`
}
