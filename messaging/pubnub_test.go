package messaging

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"strings"
	"testing"
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
		channels:   *newSubscriptionEntity(),
		groups:     *newSubscriptionEntity(),
		infoLogger: log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile),
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
		channels:   *newSubscriptionEntity(),
		groups:     *newSubscriptionEntity(),
		infoLogger: log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile),
	}

	existingSuccessChannel := make(chan []byte)
	existingErrorChannel := make(chan []byte)
	errCh := make(chan []byte)
	await := make(chan bool)

	pubnub.channels.Add("existing_channel",
		existingSuccessChannel, existingErrorChannel, pubnub.infoLogger)

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
		channels:   *newSubscriptionEntity(),
		groups:     *newSubscriptionEntity(),
		infoLogger: log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile),
	}

	existingSuccessChannel := make(chan []byte)
	existingErrorChannel := make(chan []byte)
	errCh := make(chan []byte)
	await := make(chan bool)

	pubnub.groups.Add("existing_group",
		existingSuccessChannel, existingErrorChannel, pubnub.infoLogger)

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
		channels:   *newSubscriptionEntity(),
		groups:     *newSubscriptionEntity(),
		infoLogger: log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile),
	}

	existingSuccessChannel := make(chan []byte)
	existingErrorChannel := make(chan []byte)
	errCh := make(chan []byte)
	await := make(chan bool)

	pubnub.channels.Add("ex_ch1",
		existingSuccessChannel, existingErrorChannel, pubnub.infoLogger)
	pubnub.channels.Add("ex_ch2",
		existingSuccessChannel, existingErrorChannel, pubnub.infoLogger)

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

var (
	testMessage1 = `PRISE EN MAIN - Le Figaro a pu approcher les nouveaux smartphones de Google. Voici nos premières observations. Le premier «smartphone conçu Google». Voilà comment a été présenté mardi le Pixel mardi. Il ne s'agit pas tout à fait de la première`
	testMessage2 = `Everybody copies everybody. It doesn't mean they're "out of ideas" or "in a technological cul-de-sac" - or at least it doesn't necessarily mean that - it does mean they want to make money and keep users.`
)

func BenchmarkEncodeNonASCIIChars(b *testing.B) {
	for i := 0; i < b.N; i++ {
		encodeNonASCIIChars(testMessage1)
		encodeNonASCIIChars(testMessage2)
	}
}

func TestEncodeNonASCIIChars(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    testMessage1,
			expected: "PRISE EN MAIN - Le Figaro a pu approcher les nouveaux smartphones de Google. Voici nos premi\\u00e8res observations. Le premier \\u00absmartphone con\\u00e7u Google\\u00bb. Voil\\u00e0 comment a \\u00e9t\\u00e9 pr\\u00e9sent\\u00e9 mardi le Pixel mardi. Il ne s'agit pas tout \\u00e0 fait de la premi\\u00e8re",
		},
		{
			input:    testMessage2,
			expected: testMessage2, // no non-ascii characters here, so the string should be unchanged
		},
		{
			input:    "",
			expected: "",
		},
	}
	for _, tc := range cases {
		assert.Equal(t, encodeNonASCIIChars(tc.input), tc.expected)
	}
}

func TestFilterExpression(t *testing.T) {
	assert := assert.New(t)
	pubnub := Pubnub{
		infoLogger: log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
	var filterExp = "aoi_x >= 0 AND aoi_x <= 2 AND aoi_y >= 0 AND aoi_y<= 2"
	pubnub.SetFilterExpression(filterExp)
	assert.Equal(pubnub.FilterExpression(), filterExp)
}

func TestCheckCallbackNilException(t *testing.T) {
	assert := assert.New(t)
	// Handle errors in defer func with recover.
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("pkg: %v", r)
				//fmt.Println(err)
				assert.True(strings.Contains(err.Error(), "Callback is nil for GrantSubscribe"))
			}
		}

	}()

	pubnub := Pubnub{
		infoLogger: log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
	var callbackChannel = make(chan []byte)
	close(callbackChannel)
	callbackChannel = nil
	pubnub.checkCallbackNil(callbackChannel, false, "GrantSubscribe")

}

func TestCheckCallbackNil(t *testing.T) {
	assert := assert.New(t)
	// Handle errors in defer func with recover.
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("pkg: %v", r)
				//fmt.Println(err)
				assert.True(strings.Contains(err.Error(), "Callback is nil for GrantSubscribe"))
			} else {
				assert.True(true)
			}
		}

	}()
	pubnub := Pubnub{
		infoLogger: log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
	var callbackChannel = make(chan []byte)
	pubnub.checkCallbackNil(callbackChannel, false, "GrantSubscribe")

}
