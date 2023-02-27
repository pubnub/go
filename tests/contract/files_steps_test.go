package contract

import (
	"context"
	pubnub "github.com/pubnub/go/v7"
	"os"
)

func iSendAFileWithSpaceidAndMessageType(ctx context.Context, spaceId string, messageType string) error {
	commonState := getCommonState(ctx)

	file, err := os.Open("test_file.txt")
	defer file.Close()
	if err != nil {
		return err
	}

	_, s, err := commonState.pubNub.SendFile().
		Message("This is a message").
		MessageType(pubnub.MessageType(messageType)).
		SpaceId(pubnub.SpaceId(spaceId)).
		File(file).
		Name("name").
		Channel("channel").
		Execute()

	commonState.err = err
	commonState.statusResponse = s

	return nil
}
