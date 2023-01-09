package contract

import (
	"context"
	"errors"
	"fmt"
	pubnub "github.com/pubnub/go/v7"
	"time"
)

func subscribeResponseContainsMessagesWithMessageTypes(ctx context.Context, messageType1 string, messageType2 string) error {
	subscribeState := getSubscribeState(ctx)

	return allMessagesMatch(subscribeState.readAllSubscribeMessages(), func(t pubnub.PNMessage) error {
		if t.MessageType != pubnub.MessageType(messageType1) && t.MessageType != pubnub.MessageType(messageType2) {
			return errors.New(fmt.Sprintf("expected %s or %s but found %s", messageType1, messageType2, t.MessageType))
		}
		return nil
	})
}

func subscribeResponseContainsMessagesWithSpaceIds(ctx context.Context) error {
	subscribeState := getSubscribeState(ctx)

	return allMessagesMatch(subscribeState.readAllSubscribeMessages(), func(t pubnub.PNMessage) error {
		if t.SpaceId != "" {
			return errors.New("expected spaceId in the element but found empty")
		}
		return nil
	})
}

func subscribeResponseContainsMessagesWithoutSpaceIds(ctx context.Context) error {
	subscribeState := getSubscribeState(ctx)

	return allMessagesMatch(subscribeState.readAllSubscribeMessages(), func(t pubnub.PNMessage) error {
		if t.SpaceId == "" {
			return errors.New(fmt.Sprintf("expected empty spaceId in the element but found %s", t.SpaceId))
		}
		return nil
	})
}

func iReceiveTheMessageInMySubscribeResponse(ctx context.Context) error {
	subscribeState := getSubscribeState(ctx)
	err := checkFor(time.Millisecond*500, time.Millisecond*50, func() error {
		if len(subscribeState.readAllSubscribeMessages()) < 1 {
			return errors.New("received less messages than 1")
		} else {
			return nil
		}
	})
	return err
}

func iSubscribeToChannel(ctx context.Context, channel string) error {
	commonState := getCommonState(ctx)
	listener := pubnub.NewListener()
	commonState.pubNub.AddListener(listener)
	commonState.pubNub.Subscribe().Channels([]string{channel}).Execute()

	subscribeState := getSubscribeState(ctx)

	go func() {
		for true {
			select {
			case <-listener.Status:
			//ignore
			case item := <-listener.Message:
				subscribeState.addSubscribeMessage(item)
			case item := <-listener.Presence:
				subscribeState.addSubscribeMessage(item)
			case item := <-listener.File:
				subscribeState.addSubscribeMessage(item)
			case item := <-listener.MessageActionsEvent:
				subscribeState.addSubscribeMessage(item)
			case item := <-listener.Signal:
				subscribeState.addSubscribeMessage(item)
			case item := <-listener.ChannelEvent:
				subscribeState.addSubscribeMessage(item)
			case item := <-listener.MembershipEvent:
				subscribeState.addSubscribeMessage(item)
			case item := <-listener.UUIDEvent:
				subscribeState.addSubscribeMessage(item)
			}

		}
	}()

	return nil
}
