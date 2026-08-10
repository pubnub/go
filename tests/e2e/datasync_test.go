package e2e

import (
	"fmt"
	"testing"

	pubnub "github.com/pubnub/go/v9"
	"github.com/pubnub/go/v9/tests/stubs"
	"github.com/stretchr/testify/assert"
)

// entitiesTestConfig returns a self-contained config with a fixed subscribe key
// so the stubbed paths are deterministic regardless of the environment.
func entitiesTestConfig() *pubnub.Config {
	cfg := pubnub.NewConfigWithUserId(pubnub.UserId(pubnub.GenerateUUID()))
	cfg.SubscribeKey = "sub-key"
	cfg.PublishKey = "pub-key"
	return cfg
}

func TestEntitiesCreateStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "POST",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/entities", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       `{"status":201,"data":{"id":"entity-abc","entityClass":"vehicle","entityClassVersion":1,"status":"active","payload":{"make":"Toyota"},"eTag":"BfklQ..."}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 201,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.CreateEntity().
		ID("entity-abc").
		EntityClass("vehicle").
		EntityClassVersion(1).
		Status("active").
		Payload(map[string]interface{}{"make": "Toyota"}).
		Execute()

	assert.Nil(err)
	assert.Equal(201, status.StatusCode)
	assert.Equal("entity-abc", res.Data.ID)
	assert.Equal("vehicle", res.Data.EntityClass)
	assert.Equal("BfklQ...", res.Data.ETag)
}

func TestEntitiesGetStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/entities/entity-abc", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       `{"status":200,"data":{"id":"entity-abc","entityClass":"vehicle","entityClassVersion":1,"status":"active","payload":{"make":"Toyota"},"eTag":"BfklQ..."}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.GetEntity().ID("entity-abc").Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.Equal("entity-abc", res.Data.ID)
	assert.Equal("Toyota", res.Data.Payload["make"])
}

func TestEntitiesListStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/entities", cfg.SubscribeKey),
		Query:              "entity_class=vehicle&limit=20",
		ResponseBody:       `{"data":[{"id":"i789","status":"active","entityClass":"vehicle","entityClassVersion":1,"entityClassLevel":"SubKey","eTag":"1","payload":{"custom":"fields"}}],"links":{"self":"/self","next":"/next","prev":null},"meta":{"has_next":true,"has_prev":false,"next_cursor":"TjIw","prev_cursor":null,"limit":20}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.GetEntities().EntityClass("vehicle").Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.Len(res.Data, 1)
	assert.Equal("i789", res.Data[0].ID)
	assert.NotNil(res.Meta)
	assert.True(res.Meta.HasNext)
	assert.Equal("TjIw", res.Meta.NextCursor)
}

func TestEntitiesUpdateStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "PUT",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/entities/entity-abc", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       `{"status":200,"data":{"id":"entity-abc","entityClass":"vehicle","entityClassVersion":2,"status":"active","eTag":"Cxyz1..."}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.UpdateEntity().
		ID("entity-abc").
		EntityClassVersion(2).
		Status("active").
		IfMatchETag("BfklQ...").
		Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.Equal(2, res.Data.EntityClassVersion)
	assert.Equal("Cxyz1...", res.Data.ETag)
}

func TestEntitiesPatchStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "PATCH",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/entities/entity-abc", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       `{"status":200,"data":{"id":"entity-abc","entityClass":"vehicle","entityClassVersion":1,"status":"inactive","payload":{"mileage":42000},"eTag":"Dpqr3..."}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.PatchEntity().
		ID("entity-abc").
		Replace("/status", "inactive").
		Add("/payload/mileage", 42000).
		IfMatchETag("BfklQ...").
		Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.Equal("inactive", res.Data.Status)
	assert.Equal(float64(42000), res.Data.Payload["mileage"])
}

func TestEntitiesDeleteStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "DELETE",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/entities/entity-abc", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       "",
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.RemoveEntity().ID("entity-abc").IfMatchETag("Dpqr3...").Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.NotNil(res)
}

func TestRelationshipsCreateStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "POST",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/relationships", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       `{"status":201,"data":{"id":"r-123","entityAId":"u123","entityBId":"s456","relationshipClass":"ProductOwner","relationshipClassVersion":1,"status":"active","payload":{"role":"admin"},"eTag":"1"}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 201,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.CreateRelationship().
		ID("r-123").
		EntityAID("u123").
		EntityBID("s456").
		RelationshipClass("ProductOwner").
		RelationshipClassVersion(1).
		Status("active").
		Payload(map[string]interface{}{"role": "admin"}).
		Execute()

	assert.Nil(err)
	assert.Equal(201, status.StatusCode)
	assert.Equal("r-123", res.Data.ID)
	assert.Equal("u123", res.Data.EntityAID)
	assert.Equal("ProductOwner", res.Data.RelationshipClass)
	assert.Equal("1", res.Data.ETag)
}

func TestRelationshipsGetStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/relationships/r-123", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       `{"status":200,"data":{"id":"r-123","entityAId":"u123","entityBId":"s456","relationshipClass":"ProductOwner","relationshipClassVersion":1,"status":"active","payload":{"role":"admin"},"eTag":"1"}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.GetRelationship().ID("r-123").Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.Equal("r-123", res.Data.ID)
	assert.Equal("admin", res.Data.Payload["role"])
}

func TestRelationshipsListStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/relationships", cfg.SubscribeKey),
		Query:              "relationship_class=ProductOwner&entity_a_id=u123&limit=20",
		ResponseBody:       `{"data":[{"id":"r-123","entityAId":"u123","entityBId":"s456","relationshipClass":"ProductOwner","relationshipClassVersion":1,"status":"active","eTag":"1","payload":{"custom":"fields"}}],"links":{"self":"/self","next":"/next","prev":null},"meta":{"has_next":true,"has_prev":false,"next_cursor":"TjIw","prev_cursor":null,"limit":20}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.GetRelationships().
		RelationshipClass("ProductOwner").
		EntityAID("u123").
		Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.Len(res.Data, 1)
	assert.Equal("r-123", res.Data[0].ID)
	assert.NotNil(res.Meta)
	assert.True(res.Meta.HasNext)
	assert.Equal("TjIw", res.Meta.NextCursor)
}

func TestRelationshipsUpdateStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "PUT",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/relationships/r-123", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       `{"status":200,"data":{"id":"r-123","entityAId":"u123","entityBId":"s456","relationshipClass":"ProductOwner","relationshipClassVersion":2,"status":"active","eTag":"2"}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.UpdateRelationship().
		ID("r-123").
		RelationshipClassVersion(2).
		Status("active").
		IfMatchETag("1").
		Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.Equal(2, res.Data.RelationshipClassVersion)
	assert.Equal("2", res.Data.ETag)
}

func TestRelationshipsPatchStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "PATCH",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/relationships/r-123", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       `{"status":200,"data":{"id":"r-123","entityAId":"u123","entityBId":"s456","relationshipClass":"ProductOwner","relationshipClassVersion":1,"status":"active","payload":{"custom":{"role":"admin"}},"eTag":"2"}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.PatchRelationship().
		ID("r-123").
		Add("/payload/custom/role", "admin").
		IfMatchETag("1").
		Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.Equal("2", res.Data.ETag)
	custom := res.Data.Payload["custom"].(map[string]interface{})
	assert.Equal("admin", custom["role"])
}

func TestRelationshipsDeleteStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "DELETE",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/relationships/r-123", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       "",
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.RemoveRelationship().ID("r-123").IfMatchETag("2").Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.NotNil(res)
}

func TestUsersCreateStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "POST",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/users", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       `{"status":201,"data":{"id":"alice","entityClass":"User","entityClassVersion":1,"entityClassLevel":"SubKey","status":"active","payload":{"name":"Alice"},"eTag":"1"}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 201,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.CreateUser().
		ID("alice").
		EntityClassVersion(1).
		Status("active").
		Payload(map[string]interface{}{"name": "Alice"}).
		Execute()

	assert.Nil(err)
	assert.Equal(201, status.StatusCode)
	assert.Equal("alice", res.Data.ID)
	assert.Equal("User", res.Data.EntityClass)
	assert.Equal("1", res.Data.ETag)
}

func TestUsersGetStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/users/alice", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       `{"status":200,"data":{"id":"alice","entityClass":"User","entityClassVersion":1,"status":"active","payload":{"name":"Alice"},"eTag":"1"}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.GetUser().ID("alice").Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.Equal("alice", res.Data.ID)
	assert.Equal("Alice", res.Data.Payload["name"])
}

func TestUsersListStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/users", cfg.SubscribeKey),
		Query:              "limit=20",
		ResponseBody:       `{"data":[{"id":"alice","status":"active","entityClass":"User","entityClassVersion":1,"entityClassLevel":"SubKey","eTag":"1","payload":{"name":"Alice"}}],"links":{"self":"/self","next":"/next","prev":null},"meta":{"has_next":true,"has_prev":false,"next_cursor":"TjIw","prev_cursor":null,"limit":20}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.GetUsers().Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.Len(res.Data, 1)
	assert.Equal("alice", res.Data[0].ID)
	assert.NotNil(res.Meta)
	assert.True(res.Meta.HasNext)
	assert.Equal("TjIw", res.Meta.NextCursor)
}

func TestUsersUpdateStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "PUT",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/users/alice", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       `{"status":200,"data":{"id":"alice","entityClass":"User","entityClassVersion":1,"status":"active","eTag":"2"}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.UpdateUser().
		ID("alice").
		EntityClassVersion(1).
		Status("active").
		IfMatchETag("1").
		Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.Equal("2", res.Data.ETag)
}

func TestUsersPatchStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "PATCH",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/users/alice", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       `{"status":200,"data":{"id":"alice","entityClass":"User","entityClassVersion":1,"status":"active","payload":{"profileUrl":"https://example.com/alice"},"eTag":"2"}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.PatchUser().
		ID("alice").
		Replace("/payload/profileUrl", "https://example.com/alice").
		IfMatchETag("1").
		Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.Equal("https://example.com/alice", res.Data.Payload["profileUrl"])
}

func TestUsersDeleteStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "DELETE",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/users/alice", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       "",
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.RemoveUser().ID("alice").IfMatchETag("2").Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.NotNil(res)
}

func TestChannelsCreateStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "POST",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/channels", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       `{"status":201,"data":{"id":"general","entityClass":"Channel","entityClassVersion":1,"entityClassLevel":"SubKey","status":"active","payload":{"name":"General"},"eTag":"1"}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 201,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.CreateChannel().
		ID("general").
		EntityClassVersion(1).
		Status("active").
		Payload(map[string]interface{}{"name": "General"}).
		Execute()

	assert.Nil(err)
	assert.Equal(201, status.StatusCode)
	assert.Equal("general", res.Data.ID)
	assert.Equal("Channel", res.Data.EntityClass)
	assert.Equal("1", res.Data.ETag)
}

func TestChannelsGetStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/channels/general", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       `{"status":200,"data":{"id":"general","entityClass":"Channel","entityClassVersion":1,"status":"active","payload":{"name":"General"},"eTag":"1"}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.GetChannel().ID("general").Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.Equal("general", res.Data.ID)
	assert.Equal("General", res.Data.Payload["name"])
}

func TestChannelsListStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/channels", cfg.SubscribeKey),
		Query:              "limit=20",
		ResponseBody:       `{"data":[{"id":"general","status":"active","entityClass":"Channel","entityClassVersion":1,"entityClassLevel":"SubKey","eTag":"1","payload":{"name":"General"}}],"links":{"self":"/self","next":"/next","prev":null},"meta":{"has_next":true,"has_prev":false,"next_cursor":"TjIw","prev_cursor":null,"limit":20}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.GetChannels().Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.Len(res.Data, 1)
	assert.Equal("general", res.Data[0].ID)
	assert.NotNil(res.Meta)
	assert.True(res.Meta.HasNext)
	assert.Equal("TjIw", res.Meta.NextCursor)
}

func TestChannelsUpdateStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "PUT",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/channels/general", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       `{"status":200,"data":{"id":"general","entityClass":"Channel","entityClassVersion":1,"status":"active","eTag":"2"}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.UpdateChannel().
		ID("general").
		EntityClassVersion(1).
		Status("active").
		IfMatchETag("1").
		Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.Equal("2", res.Data.ETag)
}

func TestChannelsPatchStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "PATCH",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/channels/general", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       `{"status":200,"data":{"id":"general","entityClass":"Channel","entityClassVersion":1,"status":"active","payload":{"topic":"announcements"},"eTag":"2"}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.PatchChannel().
		ID("general").
		Replace("/payload/topic", "announcements").
		IfMatchETag("1").
		Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.Equal("announcements", res.Data.Payload["topic"])
}

func TestChannelsDeleteStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "DELETE",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/channels/general", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       "",
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.RemoveChannel().ID("general").IfMatchETag("2").Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.NotNil(res)
}

func TestMembershipsCreateStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "POST",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/memberships", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       `{"status":201,"data":{"id":"m-123","channelId":"general","userId":"alice","relationshipClass":"Membership","relationshipClassVersion":1,"status":"active","payload":{"role":"moderator"},"eTag":"1"}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 201,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.CreateMembership().
		ID("m-123").
		ChannelID("general").
		UserID("alice").
		RelationshipClassVersion(1).
		Status("active").
		Payload(map[string]interface{}{"role": "moderator"}).
		Execute()

	assert.Nil(err)
	assert.Equal(201, status.StatusCode)
	assert.Equal("m-123", res.Data.ID)
	assert.Equal("general", res.Data.ChannelID)
	assert.Equal("alice", res.Data.UserID)
	assert.Equal("1", res.Data.ETag)
}

func TestMembershipsGetStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/memberships/m-123", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       `{"status":200,"data":{"id":"m-123","channelId":"general","userId":"alice","relationshipClass":"Membership","relationshipClassVersion":1,"status":"active","payload":{"role":"moderator"},"eTag":"1"}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.GetMembership().ID("m-123").Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.Equal("m-123", res.Data.ID)
	assert.Equal("moderator", res.Data.Payload["role"])
}

func TestMembershipsListStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/memberships", cfg.SubscribeKey),
		Query:              "channel_id=general&limit=20",
		ResponseBody:       `{"data":[{"id":"m-123","channelId":"general","userId":"alice","relationshipClass":"Membership","relationshipClassVersion":1,"status":"active","eTag":"1","payload":{"role":"moderator"}}],"links":{"self":"/self","next":"/next","prev":null},"meta":{"has_next":true,"has_prev":false,"next_cursor":"TjIw","prev_cursor":null,"limit":20}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.GetMemberships().ChannelID("general").Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.Len(res.Data, 1)
	assert.Equal("m-123", res.Data[0].ID)
	assert.NotNil(res.Meta)
	assert.True(res.Meta.HasNext)
	assert.Equal("TjIw", res.Meta.NextCursor)
}

func TestMembershipsUpdateStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "PUT",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/memberships/m-123", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       `{"status":200,"data":{"id":"m-123","channelId":"general","userId":"alice","relationshipClass":"Membership","relationshipClassVersion":1,"status":"active","eTag":"2"}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.UpdateMembership().
		ID("m-123").
		RelationshipClassVersion(1).
		Status("active").
		IfMatchETag("1").
		Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.Equal("2", res.Data.ETag)
}

func TestMembershipsPatchStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "PATCH",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/memberships/m-123", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       `{"status":200,"data":{"id":"m-123","channelId":"general","userId":"alice","relationshipClass":"Membership","relationshipClassVersion":1,"status":"active","payload":{"role":"admin"},"eTag":"2"}}`,
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.PatchMembership().
		ID("m-123").
		Replace("/payload/role", "admin").
		IfMatchETag("1").
		Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.Equal("admin", res.Data.Payload["role"])
}

func TestMembershipsDeleteStubbed(t *testing.T) {
	assert := assert.New(t)
	cfg := entitiesTestConfig()

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "DELETE",
		Path:               fmt.Sprintf("/v1/datasync/subkeys/%s/memberships/m-123", cfg.SubscribeKey),
		Query:              "",
		ResponseBody:       "",
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(cfg)
	pn.SetClient(interceptor.GetClient())

	res, status, err := pn.DataSync.RemoveMembership().ID("m-123").IfMatchETag("2").Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.NotNil(res)
}
