package messaging

import (
	"github.com/stretchr/testify/assert"
	//"log"
	//"os"
	"strings"
	"testing"
)

func TestGetChannelsAndGroupsChannels(t *testing.T) {
	assert := assert.New(t)
	response := `{"t":{"t":"14586613280736475","r":4},"m":[{"a":"1","f":0,"i":"UUID_SubscriptionConnectedForSimple","s":1,"p":{"t":"14593254434932405","r":4},"k":"sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f","c":"Channel_SubscriptionConnectedForSimple","b":"Channel_SubscriptionConnectedForSimple","d":"Test message"},{"a":"1","f":0,"i":"UUID_SubscriptionConnectedForSimple","s":2,"p":{"t":"14593254434932405","r":4},"k":"sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f","c":"Channel_SubscriptionConnectedForSimple","b":"Channel_SubscriptionConnectedForSimple","d":"Test message2"}]}`
	pubnub := NewPubnub("demo", "demo", "demo", "enigma", true, "testuuid")
	pubnub.channels = *newSubscriptionEntity()
	pubnub.groups = *newSubscriptionEntity()
	var callbackChannel = make(chan []byte)
	var errorChannel = make(chan []byte)

	channel := "ch"
	channelGroup := "cg"
	pubnub.channels.Add(channel, callbackChannel, errorChannel, pubnub.infoLogger)
	pubnub.groups.Add(channelGroup, callbackChannel, errorChannel, pubnub.infoLogger)

	subEnvelope, _, _, _ := pubnub.ParseSubscribeResponse([]byte(response), "")
	channelNames, channelGroupNames := subEnvelope.getChannelsAndGroups(pubnub)

	strch := strings.Join(channelNames, ",")
	strcg := strings.Join(channelGroupNames, ",")

	//log.SetOutput(os.Stdout)
	//log.Printf("strch:%s", strch)
	//log.Printf("strcg:%s", strcg)
	assert.Equal("Channel_SubscriptionConnectedForSimple,Channel_SubscriptionConnectedForSimple", strch)
	assert.Equal("", strcg)
}

func TestGetChannelsAndGroupsChannelAndChannelGroup(t *testing.T) {
	assert := assert.New(t)
	response := `{"t":{"t":"14586613280736475","r":4},"m":[{"a":"1","f":0,"i":"UUID_SubscriptionConnectedForSimple","s":2,"p":{"t":"14593254434932405","r":4},"k":"sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f","c":"Channel_SubscriptionConnectedForSimple","b":"Channel_SubscriptionConnectedForSimple_CG","d":"Test message2"}]}`
	pubnub := NewPubnub("demo", "demo", "demo", "enigma", true, "testuuid")
	pubnub.channels = *newSubscriptionEntity()
	pubnub.groups = *newSubscriptionEntity()
	var callbackChannel = make(chan []byte)
	var errorChannel = make(chan []byte)

	channel := "ch"
	channelGroup := "cg"
	pubnub.channels.Add(channel, callbackChannel, errorChannel, pubnub.infoLogger)
	pubnub.groups.Add(channelGroup, callbackChannel, errorChannel, pubnub.infoLogger)

	subEnvelope, _, _, _ := pubnub.ParseSubscribeResponse([]byte(response), "")
	channelNames, channelGroupNames := subEnvelope.getChannelsAndGroups(pubnub)

	strch := strings.Join(channelNames, ",")
	strcg := strings.Join(channelGroupNames, ",")

	//log.SetOutput(os.Stdout)
	//log.Printf("strch:%s", strch)
	//log.Printf("strcg:%s", strcg)
	assert.Equal("Channel_SubscriptionConnectedForSimple", strch)
	assert.Equal("Channel_SubscriptionConnectedForSimple_CG", strcg)
}

func TestGetChannelsAndGroupsWildcard(t *testing.T) {
	assert := assert.New(t)
	response := `{"t":{"t":"14586613280736475","r":4},"m":[{"a":"1","f":0,"i":"UUID_SubscriptionConnectedForSimple","s":2,"p":{"t":"14593254434932405","r":4},"k":"sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f","c":"Channel_SubscriptionConnectedForSimple","b":"Channel_SubscriptionConnectedForSimple.*","d":"Test message2"}]}`
	pubnub := NewPubnub("demo", "demo", "demo", "enigma", true, "testuuid")
	pubnub.channels = *newSubscriptionEntity()
	pubnub.groups = *newSubscriptionEntity()
	var callbackChannel = make(chan []byte)
	var errorChannel = make(chan []byte)

	channel := "ch"
	channelGroup := "cg"
	pubnub.channels.Add(channel, callbackChannel, errorChannel, pubnub.infoLogger)
	pubnub.groups.Add(channelGroup, callbackChannel, errorChannel, pubnub.infoLogger)

	subEnvelope, _, _, _ := pubnub.ParseSubscribeResponse([]byte(response), "")
	channelNames, channelGroupNames := subEnvelope.getChannelsAndGroups(pubnub)

	strch := strings.Join(channelNames, ",")
	strcg := strings.Join(channelGroupNames, ",")

	//log.SetOutput(os.Stdout)
	//log.Printf("strch:%s", strch)
	//log.Printf("strcg:%s", strcg)
	assert.Equal("Channel_SubscriptionConnectedForSimple", strch)
	assert.Equal("", strcg)
}
