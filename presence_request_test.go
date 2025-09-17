package pubnub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPresenceBuilder(t *testing.T) {
	assert := assert.New(t)

	o := newPresenceBuilder(pubnub)

	assert.NotNil(o)
	assert.NotNil(o.opts)
	assert.Equal(pubnub, o.opts.pubnub)
}

func TestNewPresenceBuilderWithContext(t *testing.T) {
	assert := assert.New(t)

	o := newPresenceBuilderWithContext(pubnub, backgroundContext)

	assert.NotNil(o)
	assert.NotNil(o.opts)
	assert.Equal(pubnub, o.opts.pubnub)
	assert.Equal(backgroundContext, o.opts.ctx)
}

func TestPresenceBuilderChannels(t *testing.T) {
	assert := assert.New(t)

	o := newPresenceBuilder(pubnub)
	channels := []string{"ch1", "ch2", "ch3"}

	result := o.Channels(channels)

	assert.Equal(o, result) // Should return same instance for chaining
	assert.Equal(channels, o.opts.channels)
}

func TestPresenceBuilderChannelGroups(t *testing.T) {
	assert := assert.New(t)

	o := newPresenceBuilder(pubnub)
	channelGroups := []string{"cg1", "cg2", "cg3"}

	result := o.ChannelGroups(channelGroups)

	assert.Equal(o, result) // Should return same instance for chaining
	assert.Equal(channelGroups, o.opts.channelGroups)
}

func TestPresenceBuilderConnected(t *testing.T) {
	assert := assert.New(t)

	o := newPresenceBuilder(pubnub)

	// Test setting connected to true
	result := o.Connected(true)
	assert.Equal(o, result) // Should return same instance for chaining
	assert.True(o.opts.connected)

	// Test setting connected to false
	o.Connected(false)
	assert.False(o.opts.connected)
}

func TestPresenceBuilderState(t *testing.T) {
	assert := assert.New(t)

	o := newPresenceBuilder(pubnub)
	state := map[string]interface{}{
		"status": "online",
		"age":    25,
	}

	result := o.State(state)

	assert.Equal(o, result) // Should return same instance for chaining
	assert.Equal(state, o.opts.state)
}

func TestPresenceBuilderQueryParam(t *testing.T) {
	assert := assert.New(t)

	o := newPresenceBuilder(pubnub)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	result := o.QueryParam(queryParam)

	assert.Equal(o, result) // Should return same instance for chaining
	assert.Equal(queryParam, o.opts.queryParam)
}

func TestPresenceBuilderFluentInterface(t *testing.T) {
	assert := assert.New(t)

	o := newPresenceBuilder(pubnub)
	channels := []string{"ch1", "ch2"}
	channelGroups := []string{"cg1"}
	state := map[string]interface{}{"status": "active"}
	queryParam := map[string]string{"param": "value"}

	// Test fluent interface chaining
	result := o.Channels(channels).
		ChannelGroups(channelGroups).
		Connected(true).
		State(state).
		QueryParam(queryParam)

	assert.Equal(o, result) // Should return same instance
	assert.Equal(channels, o.opts.channels)
	assert.Equal(channelGroups, o.opts.channelGroups)
	assert.True(o.opts.connected)
	assert.Equal(state, o.opts.state)
	assert.Equal(queryParam, o.opts.queryParam)
}

func TestPresenceBuilderDefaults(t *testing.T) {
	assert := assert.New(t)

	o := newPresenceBuilder(pubnub)

	// Test default values
	assert.Nil(o.opts.channels)
	assert.Nil(o.opts.channelGroups)
	assert.False(o.opts.connected)
	assert.Nil(o.opts.state)
	assert.Nil(o.opts.queryParam)
}

func TestPresenceBuilderEmptyValues(t *testing.T) {
	assert := assert.New(t)

	o := newPresenceBuilder(pubnub)

	// Test with empty slices and maps
	o.Channels([]string{})
	o.ChannelGroups([]string{})
	o.State(map[string]interface{}{})
	o.QueryParam(map[string]string{})

	assert.NotNil(o.opts.channels)
	assert.Len(o.opts.channels, 0)
	assert.NotNil(o.opts.channelGroups)
	assert.Len(o.opts.channelGroups, 0)
	assert.NotNil(o.opts.state)
	assert.Len(o.opts.state, 0)
	assert.NotNil(o.opts.queryParam)
	assert.Len(o.opts.queryParam, 0)
}

func TestPresenceBuilderNilValues(t *testing.T) {
	assert := assert.New(t)

	o := newPresenceBuilder(pubnub)

	// Test with nil values
	o.Channels(nil)
	o.ChannelGroups(nil)
	o.State(nil)
	o.QueryParam(nil)

	assert.Nil(o.opts.channels)
	assert.Nil(o.opts.channelGroups)
	assert.Nil(o.opts.state)
	assert.Nil(o.opts.queryParam)
}
