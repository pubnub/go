package messaging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubscriptionEntity(t *testing.T) {
	channels := *newSubscriptionEntity()

	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	channels.Add("qwer", successChannel, errorChannel)
	channels.Add("asdf", successChannel, errorChannel)
	channels.Add("zxcv", successChannel, errorChannel)

	assert.Equal(t, "", channels.ConnectedNamesString(), "should be equal")
	assert.Len(t, channels.NamesString(), 14, "should be 14")
	assert.Contains(t, channels.NamesString(), "asdf", "should contain asdf")
	assert.Contains(t, channels.NamesString(), "qwer", "should contain qwer")
	assert.Contains(t, channels.NamesString(), "zxcv", "should contain zxcv")
}

func TestSubscriptionPanicOnUndefinedResponseType(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "Undefined response type: 0", r)
		}
	}()

	event := connectionEvent{}
	event.Bytes()
}
