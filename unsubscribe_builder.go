package pubnub

type unsubscribeBuilder struct {
	operation *UnsubscribeOperation
	pubnub    *PubNub
}

func newUnsubscribeBuilder(pubnub *PubNub) *unsubscribeBuilder {
	builder := unsubscribeBuilder{
		pubnub:    pubnub,
		operation: &UnsubscribeOperation{},
	}

	return &builder
}

// Channels set the channels for the Unsubscribe request.
func (b *unsubscribeBuilder) Channels(channels []string) *unsubscribeBuilder {
	b.operation.Channels = channels

	return b
}

// ChannelGroups set the Channel Groups for the Unsubscribe request.
func (b *unsubscribeBuilder) ChannelGroups(groups []string) *unsubscribeBuilder {
	b.operation.ChannelGroups = groups

	return b
}

// Execute runs the Unsubscribe request and unsubscribes from the specified channels.
func (b *unsubscribeBuilder) Execute() {
	b.pubnub.subscriptionManager.adaptUnsubscribe(b.operation)
}
