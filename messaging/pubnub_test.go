package messaging

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	//"os"
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

func TestExtractMessage(t *testing.T) {
	assert := assert.New(t)

	pubnub := Pubnub{
		infoLogger: log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
	response := `{"t":{"t":"14586613280736475","r":4},"m":[{"a":"1","f":0,"i":"UUID_SubscriptionConnectedForSimple","s":1,"p":{"t":"14593254434932405","r":4},"k":"sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f","c":"Channel_SubscriptionConnectedForSimple","b":"Channel_SubscriptionConnectedForSimple","d":"Test message"},{"a":"1","f":0,"i":"UUID_SubscriptionConnectedForSimple","s":2,"p":{"t":"14593254434932405","r":4},"k":"sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f","c":"Channel_SubscriptionConnectedForSimple","b":"Channel_SubscriptionConnectedForSimple","d":"Test message2"}]}`

	subEnvelope, newTimetoken, region, _ := pubnub.ParseSubscribeResponse([]byte(response), "")
	count := 0
	if subEnvelope.Messages != nil {
		for _, msg := range subEnvelope.Messages {
			count++
			var message = pubnub.extractMessage(msg)
			var msgStr = string(message)
			if count == 1 {
				assert.Equal("\"Test message\"", msgStr)
			} else {
				assert.Equal("\"Test message2\"", msgStr)
			}
		}
	}
	assert.Equal(newTimetoken, "14586613280736475")
	assert.Equal("4", region)
	assert.Equal(2, count)

}

func TestExtractMessageCipherNonEncryptedMessage(t *testing.T) {
	assert := assert.New(t)

	pubnub := Pubnub{
		cipherKey:  "enigma",
		infoLogger: log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
	response := `{"t":{"t":"14586613280736475","r":4},"m":[{"a":"1","f":0,"i":"UUID_SubscriptionConnectedForSimple","s":1,"p":{"t":"14593254434932405","r":4},"k":"sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f","c":"Channel_SubscriptionConnectedForSimple","b":"Channel_SubscriptionConnectedForSimple","d":"Test message"},{"a":"1","f":0,"i":"UUID_SubscriptionConnectedForSimple","s":2,"p":{"t":"14593254434932405","r":4},"k":"sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f","c":"Channel_SubscriptionConnectedForSimple","b":"Channel_SubscriptionConnectedForSimple","d":"Test message2"}]}`

	subEnvelope, newTimetoken, region, _ := pubnub.ParseSubscribeResponse([]byte(response), "")
	count := 0
	if subEnvelope.Messages != nil {
		for _, msg := range subEnvelope.Messages {
			count++
			var message = pubnub.extractMessage(msg)
			var msgStr = string(message)
			if count == 1 {
				assert.Equal("\"Test message\"", msgStr)
			} else {
				assert.Equal("\"Test message2\"", msgStr)
			}
		}
	}
	assert.Equal(newTimetoken, "14586613280736475")
	assert.Equal("4", region)
	assert.Equal(2, count)

}

func TestExtractMessageCipher(t *testing.T) {
	assert := assert.New(t)

	pubnub := Pubnub{
		cipherKey:  "enigma",
		infoLogger: log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
	response := `{"t":{"t":"14586613280736475","r":4},"m":[{"a":"1","f":0,"i":"UUID_SubscriptionConnectedForSimple","s":1,"p":{"t":"14593254434932405","r":4},"k":"sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f","c":"Channel_SubscriptionConnectedForSimple","b":"Channel_SubscriptionConnectedForSimple","d":"HSoHp4g0o/uHfiS1PYXzWw=="},{"a":"1","f":0,"i":"UUID_SubscriptionConnectedForSimple","s":2,"p":{"t":"14593254434932405","r":4},"k":"sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f","c":"Channel_SubscriptionConnectedForSimple","b":"Channel_SubscriptionConnectedForSimple","d":"xXch1+uwbgGgLOudCKzFSw=="}]}`

	subEnvelope, newTimetoken, region, _ := pubnub.ParseSubscribeResponse([]byte(response), "")
	count := 0
	if subEnvelope.Messages != nil {
		for _, msg := range subEnvelope.Messages {
			count++
			var message = pubnub.extractMessage(msg)
			var msgStr = string(message)
			if count == 1 {
				assert.Equal("\"Test message\"", msgStr)
			} else {
				assert.Equal("\"message2\"", msgStr)
			}
		}
	}
	assert.Equal(newTimetoken, "14586613280736475")
	assert.Equal("4", region)
	assert.Equal(2, count)

}

func TestGetDataCipher(t *testing.T) {
	assert := assert.New(t)

	pubnub := Pubnub{
		cipherKey:  "enigma",
		infoLogger: log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
	response := `[["h5Uhyc8uf3h11w5C68QsVenCf7Llvdq5XWLa1RSgdfU=","AA9MBpymUzq/bfLCtIKFB+J6L+s3UGm6xPGh9kuXsoQ=","SfGYYp58jU2FGBNNsRk0kZ8KWRjZ6OsG3OxSySd7FF0=","ek+lrKjHCJPp5wYpxWlZcg806w/SWU5dzNYmjqDVb6o=","HrIrwvdGrm3/TM4kCf0EGl5SzcD+JqOXesWtzzc8+UA="],14610686757083461,14610686757935083]`
	var contents = []byte(response)
	var s interface{}
	err := json.Unmarshal(contents, &s)
	if err == nil {
		v := s.(interface{})
		switch vv := v.(type) {
		case []interface{}:
			length := len(vv)
			if length > 0 {
				msgStr := pubnub.getData(vv[0], pubnub.cipherKey)
				//pubnub.infoLogger.Printf(msgStr)
				assert.Equal("[\"Test Message 5\",\"Test Message 6\",\"Test Message 7\",\"Test Message 8\",\"Test Message 9\"]", msgStr)
			}
		default:
			assert.Fail("default fall through")
		}
	} else {
		assert.Fail("Unmarshal failed")
	}
}

func TestGetData(t *testing.T) {
	assert := assert.New(t)

	pubnub := Pubnub{
		//cipherKey:  "enigma",
		infoLogger: log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
	response := "[[\"Test Message 5\",\"Test Message 6\",\"Test Message 7\",\"Test Message 8\",\"Test Message 9\"],14610686757083461,14610686757935083]"
	var contents = []byte(response)
	var s interface{}
	err := json.Unmarshal(contents, &s)
	if err == nil {
		v := s.(interface{})
		switch vv := v.(type) {
		case []interface{}:
			length := len(vv)
			if length > 0 {
				msgStr := pubnub.getData(vv[0], pubnub.cipherKey)
				//pubnub.infoLogger.Printf(msgStr)
				assert.Equal("[\"Test Message 5\",\"Test Message 6\",\"Test Message 7\",\"Test Message 8\",\"Test Message 9\"]", msgStr)
			}
		default:
			assert.Fail("default fall through")
		}
	} else {
		assert.Fail("Unmarshal failed %s", err.Error())
	}
}

func TestGetDataCipherNonEnc(t *testing.T) {
	assert := assert.New(t)

	pubnub := Pubnub{
		cipherKey:  "enigma",
		infoLogger: log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
	response := "[[\"Test Message 5\",\"Test Message 6\",\"Test Message 7\",\"Test Message 8\",\"Test Message 9\"],14610686757083461,14610686757935083]"
	var contents = []byte(response)
	var s interface{}
	err := json.Unmarshal(contents, &s)
	if err == nil {
		v := s.(interface{})
		switch vv := v.(type) {
		case []interface{}:
			length := len(vv)
			if length > 0 {
				msgStr := pubnub.getData(vv[0], pubnub.cipherKey)
				//pubnub.infoLogger.Printf(msgStr)
				assert.Equal("[\"Test Message 5\",\"Test Message 6\",\"Test Message 7\",\"Test Message 8\",\"Test Message 9\"]", msgStr)
			}
		default:
			assert.Fail("default fall through")
		}
	} else {
		assert.Fail("Unmarshal failed %s", err.Error())
	}
}

func TestGetDataCipherSingle(t *testing.T) {
	assert := assert.New(t)

	pubnub := Pubnub{
		cipherKey:  "enigma",
		infoLogger: log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
	response := `["h5Uhyc8uf3h11w5C68QsVenCf7Llvdq5XWLa1RSgdfU=",14610686757083461,14610686757935083]`
	var contents = []byte(response)
	var s interface{}
	err := json.Unmarshal(contents, &s)
	if err == nil {
		v := s.(interface{})
		switch vv := v.(type) {
		case []interface{}:
			length := len(vv)
			if length > 0 {
				msgStr := pubnub.parseInterface(vv, pubnub.cipherKey)
				//pubnub.infoLogger.Printf(msgStr)
				assert.Equal("[\"Test Message 5\",1.461068675708346e+16,1.4610686757935084e+16]", msgStr)
			}
		default:
			assert.Fail("default fall through")
		}
	} else {
		assert.Fail("Unmarshal failed")
	}
}

func TestGetDataSingle(t *testing.T) {
	assert := assert.New(t)

	pubnub := Pubnub{
		//cipherKey:  "enigma",
		infoLogger: log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
	response := "[\"Test Message 5\",14610686757083461,14610686757935083]"
	var contents = []byte(response)
	var s interface{}
	err := json.Unmarshal(contents, &s)
	if err == nil {
		v := s.(interface{})
		switch vv := v.(type) {
		case []interface{}:
			msgStr := pubnub.parseInterface(vv, pubnub.cipherKey)
			assert.Equal("[\"Test Message 5\",1.461068675708346e+16,1.4610686757935084e+16]", msgStr)
		default:
			assert.Fail("default fall through")
		}
	} else {
		assert.Fail("Unmarshal failed %s", err.Error())
	}
}

func TestGetDataCipherNonEncSingle(t *testing.T) {
	assert := assert.New(t)

	pubnub := Pubnub{
		cipherKey:  "enigma",
		infoLogger: log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
	response := "[\"Test Message 5\",14610686757083461,14610686757935083]"
	var contents = []byte(response)
	var s interface{}
	err := json.Unmarshal(contents, &s)
	if err == nil {
		v := s.(interface{})
		switch vv := v.(type) {
		case []interface{}:
			length := len(vv)
			if length > 0 {
				msgStr := pubnub.parseInterface(vv, pubnub.cipherKey)
				//pubnub.infoLogger.Printf(msgStr)
				assert.Equal("[\"Test Message 5\",1.461068675708346e+16,1.4610686757935084e+16]", msgStr)
			}
		default:
			assert.Fail("default fall through")
		}
	} else {
		assert.Fail("Unmarshal failed %s", err.Error())
	}
}
