package pubnub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUnsubscribeBuilder(t *testing.T) {
	assert := assert.New(t)

	o := newUnsubscribeBuilder(pubnub)

	assert.NotNil(o)
	assert.NotNil(o.operation)
	assert.Equal(pubnub, o.pubnub)
}

func TestUnsubscribeBuilderChannels(t *testing.T) {
	assert := assert.New(t)

	o := newUnsubscribeBuilder(pubnub)
	channels := []string{"ch1", "ch2", "ch3"}

	result := o.Channels(channels)

	assert.Equal(o, result) // Should return same instance for chaining
	assert.Equal(channels, o.operation.Channels)
}

func TestUnsubscribeBuilderChannelGroups(t *testing.T) {
	assert := assert.New(t)

	o := newUnsubscribeBuilder(pubnub)
	channelGroups := []string{"cg1", "cg2", "cg3"}

	result := o.ChannelGroups(channelGroups)

	assert.Equal(o, result) // Should return same instance for chaining
	assert.Equal(channelGroups, o.operation.ChannelGroups)
}

func TestUnsubscribeBuilderQueryParam(t *testing.T) {
	assert := assert.New(t)

	o := newUnsubscribeBuilder(pubnub)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	result := o.QueryParam(queryParam)

	assert.Equal(o, result) // Should return same instance for chaining
	assert.Equal(queryParam, o.operation.QueryParam)
}

func TestUnsubscribeBuilderFluentInterface(t *testing.T) {
	assert := assert.New(t)

	o := newUnsubscribeBuilder(pubnub)
	channels := []string{"ch1", "ch2"}
	channelGroups := []string{"cg1"}
	queryParam := map[string]string{"param": "value"}

	// Test fluent interface chaining
	result := o.Channels(channels).
		ChannelGroups(channelGroups).
		QueryParam(queryParam)

	assert.Equal(o, result) // Should return same instance
	assert.Equal(channels, o.operation.Channels)
	assert.Equal(channelGroups, o.operation.ChannelGroups)
	assert.Equal(queryParam, o.operation.QueryParam)
}

func TestUnsubscribeBuilderDefaults(t *testing.T) {
	assert := assert.New(t)

	o := newUnsubscribeBuilder(pubnub)

	// Test default values
	assert.Nil(o.operation.Channels)
	assert.Nil(o.operation.ChannelGroups)
	assert.Nil(o.operation.QueryParam)
}

func TestUnsubscribeBuilderEmptyValues(t *testing.T) {
	assert := assert.New(t)

	o := newUnsubscribeBuilder(pubnub)

	// Test with empty slices and maps
	o.Channels([]string{})
	o.ChannelGroups([]string{})
	o.QueryParam(map[string]string{})

	assert.NotNil(o.operation.Channels)
	assert.Len(o.operation.Channels, 0)
	assert.NotNil(o.operation.ChannelGroups)
	assert.Len(o.operation.ChannelGroups, 0)
	assert.NotNil(o.operation.QueryParam)
	assert.Len(o.operation.QueryParam, 0)
}

func TestUnsubscribeBuilderNilValues(t *testing.T) {
	assert := assert.New(t)

	o := newUnsubscribeBuilder(pubnub)

	// Test with nil values
	o.Channels(nil)
	o.ChannelGroups(nil)
	o.QueryParam(nil)

	assert.Nil(o.operation.Channels)
	assert.Nil(o.operation.ChannelGroups)
	assert.Nil(o.operation.QueryParam)
}

func TestUnsubscribeBuilderOverwriteValues(t *testing.T) {
	assert := assert.New(t)

	o := newUnsubscribeBuilder(pubnub)

	// Set initial values
	channels1 := []string{"ch1", "ch2"}
	channelGroups1 := []string{"cg1"}
	queryParam1 := map[string]string{"q1": "v1"}

	o.Channels(channels1).
		ChannelGroups(channelGroups1).
		QueryParam(queryParam1)

	assert.Equal(channels1, o.operation.Channels)
	assert.Equal(channelGroups1, o.operation.ChannelGroups)
	assert.Equal(queryParam1, o.operation.QueryParam)

	// Overwrite with new values
	channels2 := []string{"ch3", "ch4"}
	channelGroups2 := []string{"cg2", "cg3"}
	queryParam2 := map[string]string{"q2": "v2", "q3": "v3"}

	o.Channels(channels2).
		ChannelGroups(channelGroups2).
		QueryParam(queryParam2)

	assert.Equal(channels2, o.operation.Channels)
	assert.Equal(channelGroups2, o.operation.ChannelGroups)
	assert.Equal(queryParam2, o.operation.QueryParam)
}
