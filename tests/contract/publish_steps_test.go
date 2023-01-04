package contract

import (
	"context"
	pubnub "github.com/pubnub/go/v7"
)

func iPublishMessageWithSpaceIdAndMessageType(ctx context.Context, spaceId string, messageType string) error {
	commonState := getCommonState(ctx)

	_, s, err := commonState.pubNub.Publish().
		Message("whatever").
		Channel("whatever").
		MessageType(pubnub.MessageType(messageType)).
		SpaceId(pubnub.SpaceId(spaceId)).Execute()
	commonState.err = err
	commonState.statusResponse = s
	return nil
}
