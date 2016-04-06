package messaging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenUuid(t *testing.T) {
	assert := assert.New(t)

	uuid, err := GenUuid()
	assert.Nil(err)
	assert.Len(uuid, 32)
}

func TestGetSubscribeLoopActionEmptyLists(t *testing.T) {
	assert := assert.New(t)
	pubnub := Pubnub{
		channels: *newSubscriptionEntity(),
		groups:   *newSubscriptionEntity(),
	}

	errCh := make(chan []byte)

	action := pubnub.getSubscribeLoopAction("", "", errCh)
	assert.Equal(subscribeLoopDoNothing, action)

	action = pubnub.getSubscribeLoopAction("channel", "", errCh)
	assert.Equal(subscribeLoopStart, action)

	action = pubnub.getSubscribeLoopAction("", "group", errCh)
	assert.Equal(subscribeLoopStart, action)
}

func TestGetSubscribeLoopActionListWithSingleChannel(t *testing.T) {
	assert := assert.New(t)
	pubnub := Pubnub{
		channels: *newSubscriptionEntity(),
		groups:   *newSubscriptionEntity(),
	}

	existingSuccessChannel := make(chan []byte)
	existingErrorChannel := make(chan []byte)
	errCh := make(chan []byte)
	await := make(chan bool)

	pubnub.channels.Add("existing_channel",
		existingSuccessChannel, existingErrorChannel)

	action := pubnub.getSubscribeLoopAction("", "", errCh)
	assert.Equal(subscribeLoopDoNothing, action)

	action = pubnub.getSubscribeLoopAction("channel", "", errCh)
	assert.Equal(subscribeLoopRestart, action)

	action = pubnub.getSubscribeLoopAction("", "group", errCh)
	assert.Equal(subscribeLoopRestart, action)

	// existing
	go func() {
		<-errCh
		await <- true
	}()
	action = pubnub.getSubscribeLoopAction("existing_channel", "", errCh)
	<-await
	assert.Equal(subscribeLoopDoNothing, action)
}

func TestGetSubscribeLoopActionListWithSingleGroup(t *testing.T) {
	assert := assert.New(t)
	pubnub := Pubnub{
		channels: *newSubscriptionEntity(),
		groups:   *newSubscriptionEntity(),
	}

	existingSuccessChannel := make(chan []byte)
	existingErrorChannel := make(chan []byte)
	errCh := make(chan []byte)
	await := make(chan bool)

	pubnub.groups.Add("existing_group",
		existingSuccessChannel, existingErrorChannel)

	action := pubnub.getSubscribeLoopAction("", "", errCh)
	assert.Equal(subscribeLoopDoNothing, action)

	action = pubnub.getSubscribeLoopAction("channel", "", errCh)
	assert.Equal(subscribeLoopRestart, action)

	action = pubnub.getSubscribeLoopAction("", "group", errCh)
	assert.Equal(subscribeLoopRestart, action)

	// existing
	go func() {
		<-errCh
		await <- true
	}()
	action = pubnub.getSubscribeLoopAction("", "existing_group", errCh)
	<-await
	assert.Equal(subscribeLoopDoNothing, action)
}

func TestGetSubscribeLoopActionListWithMultipleChannels(t *testing.T) {
	assert := assert.New(t)
	pubnub := Pubnub{
		channels: *newSubscriptionEntity(),
		groups:   *newSubscriptionEntity(),
	}

	existingSuccessChannel := make(chan []byte)
	existingErrorChannel := make(chan []byte)
	errCh := make(chan []byte)
	await := make(chan bool)

	pubnub.channels.Add("ex_ch1",
		existingSuccessChannel, existingErrorChannel)
	pubnub.channels.Add("ex_ch2",
		existingSuccessChannel, existingErrorChannel)

	action := pubnub.getSubscribeLoopAction("ch1,ch2", "", errCh)
	assert.Equal(subscribeLoopRestart, action)

	action = pubnub.getSubscribeLoopAction("", "gr1,gr2", errCh)
	assert.Equal(subscribeLoopRestart, action)

	go func() {
		<-errCh
		await <- true
	}()
	action = pubnub.getSubscribeLoopAction("ch1,ex_ch1,ch2", "", errCh)
	<-await
	assert.Equal(subscribeLoopRestart, action)

	go func() {
		<-errCh
		<-errCh
		await <- true
	}()
	action = pubnub.getSubscribeLoopAction("ex_ch1,ex_ch2", "", errCh)
	<-await
	assert.Equal(subscribeLoopDoNothing, action)
}
