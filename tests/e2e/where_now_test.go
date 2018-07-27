package e2e

import (
	"testing"

	pubnub "github.com/pubnub/go"
	"github.com/pubnub/go/tests/stubs"
	"github.com/stretchr/testify/assert"
)

func TestWhereNowNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(configCopy())

	res, _, err := pn.WhereNow().
		UUID("person-uuid").
		Execute()

	assert.Nil(err)
	assert.Equal(0, len(res.Channels))
}

func TestWhereNowMultipleChannels(t *testing.T) {
	assert := assert.New(t)
	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/v2/presence/sub-key/sub-c-e41d50d4-43ce-11e8-a433-9e6b275e7b64/uuid/person-uuid",
		Query:              "",
		ResponseBody:       "{\"status\": 200, \"message\": \"OK\", \"payload\": {\"channels\": [\"a\",\"b\"]}, \"service\": \"Presence\"}",
		IgnoreQueryKeys:    []string{"pnsdk", "uuid"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(config)
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.WhereNow().
		UUID("person-uuid").
		Execute()

	assert.Nil(err)

	assert.Equal(2, len(res.Channels))
	assert.Equal("a", res.Channels[0])
	assert.Equal("b", res.Channels[1])
}

func TestWhereNowNoUUID(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(config)
	pn.Subscribe().Channels([]string{"ch1"}).Execute()

	res, _, err := pn.WhereNow().
		Execute()

	assert.Nil(err)

	assert.NotEqual(0, len(res.Channels))
}

func TestWhereNowNoUUIDContext(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(config)
	pn.Subscribe().Channels([]string{"ch1"}).Execute()

	res, _, err := pn.WhereNowWithContext(backgroundContext).
		Execute()

	assert.Nil(err)

	assert.NotEqual(0, len(res.Channels))
}
