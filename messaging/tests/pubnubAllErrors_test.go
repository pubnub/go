package tests

import (
	"testing"

	"github.com/pubnub/go/messaging"
	"github.com/stretchr/testify/assert"
)

func TestHistoryNetworkError(t *testing.T) {
	assert := assert.New(t)

	stop := NewAbortedTransport()
	defer stop()

	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, "")

	successGet := make(chan []byte)
	errorGet := make(chan []byte)
	go pubnubInstance.History("ch", 100, 0, 0, false, successGet, errorGet)
	select {
	case value := <-successGet:
		assert.Fail("Success response while expecting error", string(value))
	case err := <-errorGet:
		assert.Contains(string(err), abortedTransport.PnMessage)
	case <-timeouts(5):
		assert.Fail("WhereNow timeout 5s")
	}
}

func TestHistoryJSONError(t *testing.T) {
	assert := assert.New(t)

	stop := NewBadJSONTransport()
	defer stop()

	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, "")

	successGet := make(chan []byte)
	errorGet := make(chan []byte)
	go pubnubInstance.History("ch", 100, 0, 0, false, successGet, errorGet)
	select {
	case value := <-successGet:
		assert.Fail("Success response while expecting error", string(value))
	case err := <-errorGet:
		assert.Contains(string(err), "Invalid JSON")
	case <-timeouts(5):
		assert.Fail("WhereNow timeout 5s")
	}
}

func TestGroupNetworkError(t *testing.T) {
	assert := assert.New(t)

	stop := NewAbortedTransport()
	defer stop()

	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, "")

	successGet := make(chan []byte)
	errorGet := make(chan []byte)
	go pubnubInstance.ChannelGroupAddChannel("group", "ch", successGet, errorGet)
	select {
	case value := <-successGet:
		assert.Fail("Success response while expecting error", string(value))
	case err := <-errorGet:
		assert.Contains(string(err), abortedTransport.PnMessage)
		assert.Contains(string(err), "ch")
	case <-timeouts(5):
		assert.Fail("WhereNow timeout 5s")
	}
}

func TestGroupJSONError(t *testing.T) {
	assert := assert.New(t)

	stop := NewBadJSONTransport()
	defer stop()

	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, "")

	successGet := make(chan []byte)
	errorGet := make(chan []byte)
	go pubnubInstance.ChannelGroupAddChannel("group", "ch", successGet, errorGet)
	select {
	case value := <-successGet:
		assert.Fail("Success response while expecting error", string(value))
	case err := <-errorGet:
		assert.Contains(string(err), "Invalid JSON")
	case <-timeouts(5):
		assert.Fail("WhereNow timeout 5s")
	}
}

func TestWhereNowNetworkError(t *testing.T) {
	assert := assert.New(t)

	stop := NewAbortedTransport()
	defer stop()

	uuid := "UUID_WhereNow"
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, "")

	successGet := make(chan []byte)
	errorGet := make(chan []byte)
	go pubnubInstance.WhereNow(uuid, successGet, errorGet)
	select {
	case value := <-successGet:
		assert.Fail("Success response while expecting error", string(value))
	case err := <-errorGet:
		assert.Contains(string(err), abortedTransport.PnMessage)
		assert.Contains(string(err), uuid)
	case <-timeouts(5):
		assert.Fail("WhereNow timeout 5s")
	}
}

func TestWhereNowJSONError(t *testing.T) {
	assert := assert.New(t)

	stop := NewBadJSONTransport()
	defer stop()

	uuid := "UUID_WhereNow"
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, "")

	successGet := make(chan []byte)
	errorGet := make(chan []byte)
	go pubnubInstance.WhereNow(uuid, successGet, errorGet)
	select {
	case value := <-successGet:
		assert.Fail("Success response while expecting error", string(value))
	case err := <-errorGet:
		assert.Contains(string(err), "Invalid JSON")
		assert.Contains(string(err), uuid)
	case <-timeouts(5):
		assert.Fail("WhereNow timeout 5s")
	}
}

func TestGlobalHereNowNetworkError(t *testing.T) {
	assert := assert.New(t)

	stop := NewAbortedTransport()
	defer stop()

	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, "")

	successGet := make(chan []byte)
	errorGet := make(chan []byte)
	go pubnubInstance.GlobalHereNow(true, false, successGet, errorGet)
	select {
	case value := <-successGet:
		assert.Fail("Success response while expecting error", string(value))
	case err := <-errorGet:
		assert.Contains(string(err), abortedTransport.PnMessage)
	case <-timeouts(5):
		assert.Fail("WhereNow timeout 5s")
	}
}

func TestGlobalHereNowJSONError(t *testing.T) {
	assert := assert.New(t)

	stop := NewBadJSONTransport()
	defer stop()

	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, "")

	successGet := make(chan []byte)
	errorGet := make(chan []byte)
	go pubnubInstance.GlobalHereNow(true, false, successGet, errorGet)
	select {
	case value := <-successGet:
		assert.Fail("Success response while expecting error", string(value))
	case err := <-errorGet:
		assert.Contains(string(err), "Invalid JSON")
	case <-timeouts(5):
		assert.Fail("WhereNow timeout 5s")
	}
}

func TestHereNowNetworkError(t *testing.T) {
	assert := assert.New(t)

	stop := NewAbortedTransport()
	defer stop()

	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, "")

	successGet := make(chan []byte)
	errorGet := make(chan []byte)
	go pubnubInstance.HereNow("ch", true, false, successGet, errorGet)
	select {
	case value := <-successGet:
		assert.Fail("Success response while expecting error", string(value))
	case err := <-errorGet:
		assert.Contains(string(err), abortedTransport.PnMessage)
	case <-timeouts(5):
		assert.Fail("WhereNow timeout 5s")
	}
}

func TestHereNowJSONError(t *testing.T) {
	assert := assert.New(t)

	stop := NewBadJSONTransport()
	defer stop()

	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, "")

	successGet := make(chan []byte)
	errorGet := make(chan []byte)
	go pubnubInstance.HereNow("ch", true, false, successGet, errorGet)
	select {
	case value := <-successGet:
		assert.Fail("Success response while expecting error", string(value))
	case err := <-errorGet:
		assert.Contains(string(err), "Invalid JSON")
	case <-timeouts(5):
		assert.Fail("WhereNow timeout 5s")
	}
}

func TestGetUserStateNetworkError(t *testing.T) {
	assert := assert.New(t)

	stop := NewAbortedTransport()
	defer stop()

	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, "")

	successGet := make(chan []byte)
	errorGet := make(chan []byte)
	go pubnubInstance.GetUserState("ch", successGet, errorGet)
	select {
	case value := <-successGet:
		assert.Fail("Success response while expecting error", string(value))
	case err := <-errorGet:
		assert.Contains(string(err), abortedTransport.PnMessage)
	case <-timeouts(5):
		assert.Fail("WhereNow timeout 5s")
	}
}

func TestGetUserStateJSONError(t *testing.T) {
	assert := assert.New(t)

	stop := NewBadJSONTransport()
	defer stop()

	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, "")

	successGet := make(chan []byte)
	errorGet := make(chan []byte)
	go pubnubInstance.GetUserState("ch", successGet, errorGet)
	select {
	case value := <-successGet:
		assert.Fail("Success response while expecting error", string(value))
	case err := <-errorGet:
		assert.Contains(string(err), "Invalid JSON")
		assert.Contains(string(err), "ch")
	case <-timeouts(5):
		assert.Fail("WhereNow timeout 5s")
	}
}

func TestSetUserStateNetworkError(t *testing.T) {
	assert := assert.New(t)

	stop := NewAbortedTransport()
	defer stop()

	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, "")

	successGet := make(chan []byte)
	errorGet := make(chan []byte)
	go pubnubInstance.SetUserStateJSON("ch", "{}", successGet, errorGet)
	select {
	case value := <-successGet:
		assert.Fail("Success response while expecting error", string(value))
	case err := <-errorGet:
		assert.Contains(string(err), abortedTransport.PnMessage)
	case <-timeouts(5):
		assert.Fail("SetState timeout 5s")
	}
}

func TestSetUserStateJSONError(t *testing.T) {
	assert := assert.New(t)

	stop := NewBadJSONTransport()
	defer stop()

	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, "")

	successGet := make(chan []byte)
	errorGet := make(chan []byte)
	go pubnubInstance.SetUserStateJSON("ch", "{}", successGet, errorGet)
	select {
	case value := <-successGet:
		assert.Fail("Success response while expecting error", string(value))
	case err := <-errorGet:
		assert.Contains(string(err), "Invalid JSON")
		assert.Contains(string(err), "ch")
	case <-timeouts(5):
		assert.Fail("SetState timeout 5s")
	}
}

func TestUnsubscribeNetworkError(t *testing.T) {
	assert := assert.New(t)

	stop, _ := NewVCRBoth(
		"fixtures/unsubscribe/networkError", []string{"uuid"})
	defer stop()
	messaging.SetNonSubscribeTransport(abortedTransport)

	channel := "Channel_UnsubscribeNetError"
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, "")

	subscribeSuccess := make(chan []byte)
	subscribeError := make(chan []byte)

	go pubnubInstance.Subscribe(channel, "", subscribeSuccess, false, subscribeError)
	ExpectConnectedEvent(t, channel, "", subscribeSuccess, subscribeError)

	successGet := make(chan []byte)
	errorGet := make(chan []byte)

	go pubnubInstance.Unsubscribe(channel, successGet, errorGet)
	select {
	case value := <-successGet:
		assert.Contains(string(value), "unsubscribed")
		assert.Contains(string(value), channel)
	case err := <-errorGet:
		assert.Fail("Error response while expecting success", string(err))
	case <-timeouts(5):
		assert.Fail("WhereNow timeout 5s")
	}
	select {
	case value := <-successGet:
		assert.Fail("Success response while expecting error", string(value))
	case err := <-errorGet:
		assert.Contains(string(err), abortedTransport.PnMessage)
		assert.Contains(string(err), channel)
	case <-timeouts(5):
		assert.Fail("WhereNow timeout 5s")
	}

	messaging.SetNonSubscribeTransport(nil)
}

func TestGroupUnsubscribeNetworkError(t *testing.T) {
	assert := assert.New(t)

	stop, sleep := NewVCRBoth(
		"fixtures/unsubscribe/groupNetworkError", []string{"uuid"})
	defer stop()

	group := "Channel_GroupUnsubscribeNetError"
	pubnubInstance := messaging.NewPubnub(PubKey, SubKey, SecKey, "", false, "")

	createChannelGroups(pubnubInstance, []string{group})
	defer removeChannelGroups(pubnubInstance, []string{group})

	sleep(2)

	subscribeSuccess := make(chan []byte)
	subscribeError := make(chan []byte)

	go pubnubInstance.ChannelGroupSubscribe(group, subscribeSuccess, subscribeError)
	ExpectConnectedEvent(t, "", group, subscribeSuccess, subscribeError)

	successGet := make(chan []byte)
	errorGet := make(chan []byte)

	messaging.SetNonSubscribeTransport(abortedTransport)

	go pubnubInstance.ChannelGroupUnsubscribe(group, successGet, errorGet)
	select {
	case value := <-successGet:
		assert.Contains(string(value), "unsubscribed")
		assert.Contains(string(value), group)
	case err := <-errorGet:
		assert.Fail("Error response while expecting success", string(err))
	case <-timeouts(5):
		assert.Fail("WhereNow timeout 5s")
	}
	select {
	case value := <-successGet:
		assert.Fail("Success response while expecting error", string(value))
	case err := <-errorGet:
		assert.Contains(string(err), abortedTransport.PnMessage)
		assert.Contains(string(err), group)
	case <-timeouts(5):
		assert.Fail("WhereNow timeout 5s")
	}

	messaging.SetNonSubscribeTransport(nil)
}
