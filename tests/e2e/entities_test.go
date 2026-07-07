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

	res, status, err := pn.CreateEntity().
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

	res, status, err := pn.GetEntity().ID("entity-abc").Execute()

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

	res, status, err := pn.GetEntities().EntityClass("vehicle").Execute()

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

	res, status, err := pn.UpdateEntity().
		ID("entity-abc").
		EntityClassVersion(2).
		Status("active").
		IfMatch("BfklQ...").
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

	res, status, err := pn.PatchEntity().
		ID("entity-abc").
		Replace("/status", "inactive").
		Add("/payload/mileage", 42000).
		IfMatch("BfklQ...").
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

	res, status, err := pn.DeleteEntity().ID("entity-abc").IfMatch("Dpqr3...").Execute()

	assert.Nil(err)
	assert.Equal(200, status.StatusCode)
	assert.NotNil(res)
}
