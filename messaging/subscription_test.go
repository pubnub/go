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

func TestSubscriptionRemoveNonExistingItem(t *testing.T) {
	items := testEntityWithOneItem()

	assert.False(t, items.Remove("blah"))
}

func TestSubscriptionClear(t *testing.T) {
	assert := assert.New(t)

	items := testEntityWithOneItem()

	assert.Len(items.items, 1)
	items.Clear()

	assert.Zero(len(items.items))
}

func TestSubscriptionAbort(t *testing.T) {
	assert := assert.New(t)

	items := testEntityWithOneItem()

	assert.False(items.abortedMarker)
	assert.Len(items.items, 1)

	items.ApplyAbort()

	assert.False(items.abortedMarker)
	assert.Len(items.items, 1)

	items.Abort()

	assert.True(items.abortedMarker)
	assert.Len(items.items, 1)

	items.ApplyAbort()

	assert.True(items.abortedMarker)
	assert.Zero(len(items.items))
}

func TestSubscriptionResetConnected(t *testing.T) {
	assert := assert.New(t)

	items := testEntityWithOneItem()
	items.items["qwer"].Connected = true

	assert.Equal([]string{"qwer"}, items.ConnectedNames())

	items.ResetConnected()

	assert.Equal([]string{}, items.ConnectedNames())
}

func testEntityWithOneItem() *subscriptionEntity {
	items := newSubscriptionEntity()

	items.items["qwer"] = &subscriptionItem{}

	return items
}
