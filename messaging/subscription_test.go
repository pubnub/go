package messaging

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSubscriptionEntity(t *testing.T) {
	channels := *NewSubscriptionEntity()

	successChannel := make(chan SuccessResponse)
	errorChannel := make(chan ErrorResponse)
	eventChannel := make(chan ConnectionEvent)

	channels.Add("qwer", successChannel, errorChannel, eventChannel)
	channels.Add("asdf", successChannel, errorChannel, eventChannel)
	channels.Add("zxcv", successChannel, errorChannel, eventChannel)

	assert.Equal(t, "", channels.ConnectedNamesString(), "should be equal")
	assert.Len(t, channels.NamesString(), 14, "should be 14")
	assert.Contains(t, channels.NamesString(), "asdf", "should contain asdf")
	assert.Contains(t, channels.NamesString(), "qwer", "should contain qwer")
	assert.Contains(t, channels.NamesString(), "zxcv", "should contain zxcv")
}
