package e2e

import (
	"errors"
	"fmt"
	"testing"
	"time"

	pubnub "github.com/pubnub/go/v7"
	"github.com/pubnub/go/v7/tests/stubs"
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
		Path:               fmt.Sprintf("/v2/presence/sub-key/%s/uuid/person-uuid", config.SubscribeKey),
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
	time.Sleep(500 * time.Millisecond)

	checkFor(assert, 5*time.Second, 500*time.Millisecond, func() error {
		res, _, err := pn.WhereNow().
			Execute()

		if err != nil {
			return err
		}

		if len(res.Channels) == 0 {
			return errors.New("res.Channels can't be empty")
		}

		return nil
	})
}

func TestWhereNowNoUUIDContext(t *testing.T) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(config)
	pn.Subscribe().Channels([]string{"ch1"}).Execute()
	time.Sleep(500 * time.Millisecond)

	checkFor(assert, 5*time.Second, 500*time.Millisecond, func() error {
		res, _, err := pn.WhereNowWithContext(backgroundContext).
			Execute()

		if err != nil {
			return err
		}

		if len(res.Channels) == 0 {
			return errors.New("res.Channels can't be empty")
		}

		return nil
	})
}
