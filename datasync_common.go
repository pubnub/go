package pubnub

// This file holds the shared types and constants used by the DataSync
// data-plane API (Entities and Relationships). System fields live at the top
// level of each resource while application-defined fields live under Payload.

const (
	// entitiesPath is the collection path used for list (GET) and create (POST).
	entitiesPath = "/v1/datasync/subkeys/%s/entities"
	// entitiesIDPath is the item path used for read/update/patch/delete.
	entitiesIDPath = "/v1/datasync/subkeys/%s/entities/%s"

	// relationshipsPath is the collection path used for list (GET) and create (POST).
	relationshipsPath = "/v1/datasync/subkeys/%s/relationships"
	// relationshipsIDPath is the item path used for read/update/patch/delete.
	relationshipsIDPath = "/v1/datasync/subkeys/%s/relationships/%s"

	// entityContentType is the media type sent on create (POST) and full
	// replacement (PUT) entity requests.
	entityContentType = "application/vnd.pubnub.objects.entity+json;version=1"
	// relationshipContentType is the media type sent on create (POST) and full
	// replacement (PUT) relationship requests.
	relationshipContentType = "application/vnd.pubnub.objects.relationship+json;version=1"
	// entityPatchContentType is the media type sent on JSON Patch (PATCH) requests
	// for both entities and relationships.
	entityPatchContentType = "application/json-patch+json"

	// entitiesDefaultLimit is the default page size for list requests when the
	// caller does not specify one. Valid server range is 1-100. Shared by
	// Entities and Relationships list endpoints.
	entitiesDefaultLimit = 20
)

// PNEntity is the generic entity resource returned by the Entities API.
// System fields are top-level; application-defined fields live under Payload.
type PNEntity struct {
	ID                    string                 `json:"id"`
	Status                string                 `json:"status,omitempty"`
	EntityClass           string                 `json:"entityClass"`
	EntityClassVersion    int                    `json:"entityClassVersion"`
	EntityClassLevel      string                 `json:"entityClassLevel,omitempty"`
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
// (create, get, update, patch).
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
// PatchEntity / PatchRelationship. Value is required for "add", "replace" and
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
// operations (create, get, update, patch).
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
