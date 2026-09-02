package pubnub

// This file holds the shared types and constants used by the DataSync
// data-plane API (Entities, Relationships, Users, Channels and Memberships).
// System fields live at the top level of each resource while application-defined
// fields live under Payload.

const (
	// entitiesPath is the collection path used for list (GET) and create (POST).
	entitiesPath = "/v1/datasync/subkeys/%s/entities"
	// entitiesIDPath is the item path used for read/set/update/delete.
	entitiesIDPath = "/v1/datasync/subkeys/%s/entities/%s"

	// relationshipsPath is the collection path used for list (GET) and create (POST).
	relationshipsPath = "/v1/datasync/subkeys/%s/relationships"
	// relationshipsIDPath is the item path used for read/set/update/delete.
	relationshipsIDPath = "/v1/datasync/subkeys/%s/relationships/%s"

	// usersPath is the collection path used for list (GET) and create (POST).
	usersPath = "/v1/datasync/subkeys/%s/users"
	// usersIDPath is the item path used for read/set/update/delete.
	usersIDPath = "/v1/datasync/subkeys/%s/users/%s"

	// channelsPath is the collection path used for list (GET) and create (POST).
	channelsPath = "/v1/datasync/subkeys/%s/channels"
	// channelsIDPath is the item path used for read/set/update/delete.
	channelsIDPath = "/v1/datasync/subkeys/%s/channels/%s"

	// membershipsPath is the collection path used for list (GET) and create (POST).
	membershipsPath = "/v1/datasync/subkeys/%s/memberships"
	// membershipsIDPath is the item path used for read/set/update/delete.
	membershipsIDPath = "/v1/datasync/subkeys/%s/memberships/%s"

	// entityContentType is the media type sent on create (POST) and full
	// replacement (PUT) entity requests.
	entityContentType = "application/vnd.pubnub.objects.entity+json;version=1"
	// relationshipContentType is the media type sent on create (POST) and full
	// replacement (PUT) relationship requests.
	relationshipContentType = "application/vnd.pubnub.objects.relationship+json;version=1"
	// userContentType is the media type sent on create (POST) and full
	// replacement (PUT) user requests.
	userContentType = "application/vnd.pubnub.objects.user+json;version=1"
	// channelContentType is the media type sent on create (POST) and full
	// replacement (PUT) channel requests.
	channelContentType = "application/vnd.pubnub.objects.channel+json;version=1"
	// membershipContentType is the media type sent on create (POST) and full
	// replacement (PUT) membership requests.
	membershipContentType = "application/vnd.pubnub.objects.membership+json;version=1"
	// entityPatchContentType is the media type sent on JSON Patch (PATCH) requests
	// for entities, relationships, users, channels and memberships.
	entityPatchContentType = "application/json-patch+json"

	// entitiesDefaultLimit is the default page size for list requests when the
	// caller does not specify one. Valid server range is 1-100. Shared by
	// Entities, Relationships, Users, Channels and Memberships list endpoints.
	entitiesDefaultLimit = 20
)

func conflictingDataSyncFilters(filterFast, filter string) bool {
	return filterFast != "" && filter != ""
}

// PNEntity is the generic entity resource returned by the Entities API.
// System fields are top-level; application-defined fields live under Payload.
type PNEntity struct {
	ID                    string                 `json:"id"`
	Status                string                 `json:"status,omitempty"`
	EntityClass           string                 `json:"entityClass"`
	EntityClassVersion    int                    `json:"entityClassVersion"`
	EntityClassLevel      PNEntityClassLevel     `json:"entityClassLevel,omitempty"`
	EntityClassExtendsRef string                 `json:"entityClassExtendsRef,omitempty"`
	Payload               map[string]interface{} `json:"payload,omitempty"`
	CreatedAt             string                 `json:"createdAt,omitempty"`
	UpdatedAt             string                 `json:"updatedAt,omitempty"`
	ETag                  string                 `json:"eTag,omitempty"`
	ExpiresAt             string                 `json:"expiresAt,omitempty"`
}

// PNEntityPaginationMeta carries the cursor-based pagination metadata returned
// by list responses.
type PNEntityPaginationMeta struct {
	NextCursor string `json:"next_cursor,omitempty"`
	PrevCursor string `json:"prev_cursor,omitempty"`
	HasNext    bool   `json:"has_next"`
	HasPrev    bool   `json:"has_prev"`
	Limit      int    `json:"limit,omitempty"`
}

// PNEntityPaginationLinks carries the HATEOAS navigation links returned by list
// responses.
type PNEntityPaginationLinks struct {
	Self string `json:"self,omitempty"`
	Next string `json:"next,omitempty"`
	Prev string `json:"prev,omitempty"`
}

// PNEntityResponse is the envelope returned by single-entity operations
// (create, get, set, update).
type PNEntityResponse struct {
	Status int      `json:"status,omitempty"`
	Data   PNEntity `json:"data"`
}

// PNEntitiesResponse is the envelope returned by the list (GetEntities) operation.
type PNEntitiesResponse struct {
	Status int                      `json:"status,omitempty"`
	Data   []PNEntity               `json:"data"`
	Meta   *PNEntityPaginationMeta  `json:"meta,omitempty"`
	Links  *PNEntityPaginationLinks `json:"links,omitempty"`
}

// PNJSONPatchOperation is a single RFC 6902 JSON Patch operation applied by
// UpdateEntity / UpdateRelationship. Value is required for "add", "replace" and
// "test"; From is required for "move" and "copy".
type PNJSONPatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
	From  string      `json:"from,omitempty"`
}

// PNRelationship is the generic relationship resource returned by the
// Relationships API. It links two entities (A and B). System fields and the
// linked entity IDs are top-level; application-defined fields live under Payload.
type PNRelationship struct {
	ID                       string                 `json:"id"`
	EntityAID                string                 `json:"entityAId"`
	EntityBID                string                 `json:"entityBId"`
	Status                   string                 `json:"status,omitempty"`
	RelationshipClass        string                 `json:"relationshipClass"`
	RelationshipClassVersion int                    `json:"relationshipClassVersion"`
	Payload                  map[string]interface{} `json:"payload,omitempty"`
	CreatedAt                string                 `json:"createdAt,omitempty"`
	UpdatedAt                string                 `json:"updatedAt,omitempty"`
	ETag                     string                 `json:"eTag,omitempty"`
	ExpiresAt                string                 `json:"expiresAt,omitempty"`
}

// PNRelationshipResponse is the envelope returned by single-relationship
// operations (create, get, set, update).
type PNRelationshipResponse struct {
	Status int            `json:"status,omitempty"`
	Data   PNRelationship `json:"data"`
}

// PNRelationshipsResponse is the envelope returned by the list
// (GetRelationships) operation.
type PNRelationshipsResponse struct {
	Status int                      `json:"status,omitempty"`
	Data   []PNRelationship         `json:"data"`
	Meta   *PNEntityPaginationMeta  `json:"meta,omitempty"`
	Links  *PNEntityPaginationLinks `json:"links,omitempty"`
}

// PNUser is a DataSync user resource — a specialized Entity of the User class
// (or a subclass). System fields are top-level; profile fields live under Payload.
type PNUser = PNEntity

// PNUserResponse is the envelope returned by single-user operations
// (create, get, set, update).
type PNUserResponse = PNEntityResponse

// PNUsersResponse is the envelope returned by the list (GetUsers) operation.
type PNUsersResponse = PNEntitiesResponse

// PNDataSyncChannel is a DataSync channel resource — a specialized Entity of the
// Channel class (or a subclass). System fields are top-level; custom fields live
// under Payload. Named to avoid colliding with the Objects API PNChannel type.
type PNDataSyncChannel = PNEntity

// PNDataSyncChannelResponse is the envelope returned by single-channel operations
// (create, get, set, update).
type PNDataSyncChannelResponse = PNEntityResponse

// PNDataSyncChannelsResponse is the envelope returned by the list (GetChannels) operation.
type PNDataSyncChannelsResponse = PNEntitiesResponse

// PNDataSyncMembership is a DataSync membership resource — a specialized
// Relationship of class Membership linking a channel to a user. Named to avoid
// colliding with the Objects API PNMemberships type. Uses friendly channelId /
// userId fields rather than entityAId / entityBId.
type PNDataSyncMembership struct {
	ID                       string                 `json:"id"`
	ChannelID                string                 `json:"channelId"`
	UserID                   string                 `json:"userId"`
	Status                   string                 `json:"status,omitempty"`
	RelationshipClass        string                 `json:"relationshipClass"`
	RelationshipClassVersion int                    `json:"relationshipClassVersion"`
	Payload                  map[string]interface{} `json:"payload,omitempty"`
	CreatedAt                string                 `json:"createdAt,omitempty"`
	UpdatedAt                string                 `json:"updatedAt,omitempty"`
	ETag                     string                 `json:"eTag,omitempty"`
	ExpiresAt                string                 `json:"expiresAt,omitempty"`
}

// PNDataSyncMembershipResponse is the envelope returned by single-membership
// operations (create, get, set, update).
type PNDataSyncMembershipResponse struct {
	Status int                  `json:"status,omitempty"`
	Data   PNDataSyncMembership `json:"data"`
}

// PNDataSyncMembershipsResponse is the envelope returned by the list
// (GetMemberships) operation.
type PNDataSyncMembershipsResponse struct {
	Status int                      `json:"status,omitempty"`
	Data   []PNDataSyncMembership   `json:"data"`
	Meta   *PNEntityPaginationMeta  `json:"meta,omitempty"`
	Links  *PNEntityPaginationLinks `json:"links,omitempty"`
}
