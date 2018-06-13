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

//
func (b *unsubscribeBuilder) Channels(channels []string) *unsubscribeBuilder {
	b.operation.Channels = channels

	return b
}

func (b *unsubscribeBuilder) ChannelGroups(groups []string) *unsubscribeBuilder {
	b.operation.ChannelGroups = groups

	return b
}

func (b *unsubscribeBuilder) Execute() {
	b.pubnub.subscriptionManager.adaptUnsubscribe(b.operation)
}
