package pubnub

// This file holds the shared types and constants used by the generic Entities
// data-plane API (see the /entities Endpoint Reference). Entities are the
// schema-versioned, class-based base objects of the data plane. System fields
// live at the top level of the resource while application-defined fields live
// under Payload.

const (
	// entitiesPath is the collection path used for list (GET) and create (POST).
	entitiesPath = "/v1/datasync/subkeys/%s/entities"
	// entitiesIDPath is the item path used for read/update/patch/delete.
	entitiesIDPath = "/v1/datasync/subkeys/%s/entities/%s"

	// entityContentType is the media type sent on create (POST) and full
	// replacement (PUT) requests.
	entityContentType = "application/vnd.pubnub.objects.entity+json;version=1"
	// entityPatchContentType is the media type sent on JSON Patch (PATCH) requests.
	entityPatchContentType = "application/json-patch+json"

	// entitiesDefaultLimit is the default page size for list requests when the
	// caller does not specify one. Valid server range is 1-100.
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
// PatchEntity. Value is required for "add", "replace" and "test"; From is
// required for "move" and "copy".
type PNJSONPatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
	From  string      `json:"from,omitempty"`
}
