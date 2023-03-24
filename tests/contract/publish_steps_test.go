package contract

import (
	"context"
	pubnub "github.com/pubnub/go/v7"
)

func iPublishMessageWithSpaceIdAndType(ctx context.Context, spaceId string, typ string) error {
	commonState := getCommonState(ctx)

	_, s, err := commonState.pubNub.Publish().
		Message("whatever").
		Channel("whatever").
		Type(typ).
		SpaceId(pubnub.SpaceId(spaceId)).Execute()
	commonState.err = err
	commonState.statusResponse = s
	return nil
}
