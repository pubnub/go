package contract

import (
	"context"
	"fmt"
	pubnub "github.com/pubnub/go/v7"
)

func historyResponseContainsMessagesWithProvidedMessageTypes(ctx context.Context, firstMessageType string, secondMessageType string) error {
	historyState := getHistoryState(ctx)

	for _, fetchResponseItems := range historyState.fetchResponse.Messages {
		for _, item := range fetchResponseItems {
			if item.MessageType != pubnub.MessageType(firstMessageType) && item.MessageType != pubnub.MessageType(secondMessageType) {
				return fmt.Errorf("expected message type to be %s or %s, but found %s", firstMessageType, secondMessageType, item.MessageType)
			}
		}
	}
	return nil
}

func historyResponseContainsMessagesWithMessageTypes(ctx context.Context) error {
	historyState := getHistoryState(ctx)

	for _, fetchResponseItems := range historyState.fetchResponse.Messages {
		for _, item := range fetchResponseItems {
			if item.MessageType == "" {
				return fmt.Errorf("expected non empty message type")
			}
		}
	}
	return nil
}

func historyResponseContainsMessagesWithSpaceIds(ctx context.Context) error {
	historyState := getHistoryState(ctx)

	for _, fetchResponseItems := range historyState.fetchResponse.Messages {
		for _, item := range fetchResponseItems {
			if item.SpaceId == "" {
				return fmt.Errorf("expected non empty space id")
			}
		}
	}
	return nil

}

func historyResponseContainsMessagesWithoutMessageTypes(ctx context.Context) error {
	historyState := getHistoryState(ctx)

	for _, fetchResponseItems := range historyState.fetchResponse.Messages {
		for _, item := range fetchResponseItems {
			if item.MessageType != "" {
				return fmt.Errorf("expected empty message type, but found %s", item.MessageType)
			}
		}
	}
	return nil
}

func historyResponseContainsMessagesWithoutSpaceIds(ctx context.Context) error {
	historyState := getHistoryState(ctx)

	for _, fetchResponseItems := range historyState.fetchResponse.Messages {
		for _, item := range fetchResponseItems {
			if item.SpaceId != "" {
				return fmt.Errorf("expected empty space id, but found %s", item.SpaceId)
			}
		}
	}
	return nil
}

func iFetchMessageHistoryForChannel(ctx context.Context, channel string) error {
	commonState := getCommonState(ctx)
	historyState := getHistoryState(ctx)
	r, s, err := commonState.pubNub.Fetch().Channels([]string{channel}).Execute()

	commonState.err = err
	commonState.statusResponse = s
	historyState.fetchResponse = r
	return nil
}

func iFetchMessageHistoryWithIncludeMessageTypeSetToFalseForChannel(ctx context.Context, channel string) error {
	commonState := getCommonState(ctx)
	historyState := getHistoryState(ctx)

	r, s, err := commonState.pubNub.Fetch().Channels([]string{channel}).IncludeMessageType(false).Execute()

	commonState.err = err
	commonState.statusResponse = s
	historyState.fetchResponse = r
	return nil
}

func iFetchMessageHistoryWithIncludeSpaceIdSetToTrueForChannel(ctx context.Context, channel string) error {
	commonState := getCommonState(ctx)
	historyState := getHistoryState(ctx)

	r, s, err := commonState.pubNub.Fetch().Channels([]string{channel}).IncludeSpaceId(true).Execute()

	commonState.err = err
	commonState.statusResponse = s
	historyState.fetchResponse = r
	return nil
}
