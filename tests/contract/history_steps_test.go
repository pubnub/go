package contract

import (
	"context"
	"fmt"
)

func historyResponseContainsMessagesWithProvidedTypes(ctx context.Context, firstType string, secondType string) error {
	historyState := getHistoryState(ctx)

	for _, fetchResponseItems := range historyState.fetchResponse.Messages {
		for _, item := range fetchResponseItems {
			if item.Type != firstType && item.Type != secondType {
				return fmt.Errorf("expected type to be %s or %s, but found %s", firstType, secondType, item.Type)
			}
		}
	}
	return nil
}

func historyResponseContainsMessagesWithProvidedMessageTypes(ctx context.Context, firstType int, secondType int) error {
	historyState := getHistoryState(ctx)

	for _, fetchResponseItems := range historyState.fetchResponse.Messages {
		for _, item := range fetchResponseItems {
			if item.MessageType != firstType && item.MessageType != secondType {
				return fmt.Errorf("expected message type to be %d or %d, but found %d", firstType, secondType, item.MessageType)
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

func historyResponseContainsMessagesWithoutType(ctx context.Context) error {
	historyState := getHistoryState(ctx)

	for _, fetchResponseItems := range historyState.fetchResponse.Messages {
		for _, item := range fetchResponseItems {
			if item.Type != "" {
				return fmt.Errorf("expected empty message type, but found %s", item.Type)
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

func iFetchMessageHistoryWithIncludeTypeSetToFalseForChannel(ctx context.Context, channel string) error {
	commonState := getCommonState(ctx)
	historyState := getHistoryState(ctx)

	r, s, err := commonState.pubNub.Fetch().Channels([]string{channel}).IncludeType(false).Execute()

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
