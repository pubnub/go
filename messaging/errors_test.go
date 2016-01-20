package messaging

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorSingleMessage(t *testing.T) {
	assert := assert.New(t)

	err := errorResponse{
		Message: "hello",
		Reason:  responseNotSubscribed,
		Type:    channelResponse,
	}

	expected := "[0, \"Subscription to channel 'blah' not subscribed\", \"blah\"]"
	assert.Equal(expected, string(err.BytesForSource("blah")))
}

func TestErrorDoubleMessages(t *testing.T) {
	assert := assert.New(t)

	err := errorResponse{
		Message:         "hello",
		DetailedMessage: "world",
		Reason:          responseAlreadySubscribed,
		Type:            channelResponse,
	}

	expected := "[0, \"Subscription to channel 'blah' already subscribed\", world, \"blah\"]"
	assert.Equal(expected, string(err.BytesForSource("blah")))
}

func TestErrorResponseAsIs(t *testing.T) {
	assert := assert.New(t)

	err := errorResponse{
		Message: "hello",
		Reason:  responseAsIsError,
	}

	expected := "[0, \"hello\", \"blah\"]"
	assert.Equal(expected, string(err.BytesForSource("blah")))
}
