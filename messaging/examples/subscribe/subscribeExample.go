package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pubnub/go/messaging"
	"os"
)

func main() {
	messaging.SetLogOutput(os.Stderr)

	publishKey := flag.String("pub", "demo", "publish key")
	subscribeKey := flag.String("sub", "demo", "subscribe key")
	secretKey := flag.String("secret", "demo", "secret key")

	channels := flag.String("channels", "qwer,qwer-pnpres", "channels to subscribe to")
	groups := flag.String("groups", "zzz,zzz-pnpres", "channel groups to subscribe to")

	pubnub := messaging.NewPubnub(*publishKey, *subscribeKey, *secretKey, "", false, "")

	go populateGroup(pubnub, "zzz", "asdf")

	successChannel, errorChannel, eventsChannel := messaging.CreateSubscriptionChannels()

	go pubnub.Subscribe(*channels, successChannel, errorChannel, eventsChannel)
	go pubnub.ChannelGroupSubscribe(*groups, successChannel, errorChannel, eventsChannel)

	subscribeHandler(successChannel, errorChannel, eventsChannel)
}

func subscribeHandler(
	successChannel chan messaging.SuccessResponse,
	errorChannel chan messaging.ErrorResponse,
	eventsChannel chan messaging.ConnectionEvent) {

	for {
		select {
		case response := <-successChannel:
			var name string

			switch response.Type {
			case messaging.ChannelResponse:
				name = response.Channel
			case messaging.ChannelGroupResponse:
				name = response.Source
			case messaging.WildcardResponse:
				name = response.Source
			}

			if response.Presence {
				fmt.Printf("New presence event on %s %s (%s)\n", name,
					messaging.StringResponseType(response.Type), response.Timetoken)

				var presenceEvent messaging.PresenceEvent

				err := json.Unmarshal(response.Data, &presenceEvent)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}

				fmt.Printf("%s action of %s user. New occupancy %d\n\n",
					presenceEvent.Action, presenceEvent.Uuid, presenceEvent.Occupancy)
			} else {
				fmt.Printf("New message on %s %s (%s)\n", name,
					messaging.StringResponseType(response.Type), response.Timetoken)
				fmt.Printf("Received raw data: %s\n\n", response.Data)
			}

		case err := <-errorChannel:
			switch er := err.(type) {
			case messaging.ServerSideErrorResponse:
				fmt.Printf("Server-side error: %s", err.Error())
				if er.Data.Payload != nil {
					payload, err := json.Marshal(er.Data.Payload)
					if err != nil {
						fmt.Println(err.Error())
						continue
					}

					fmt.Printf("Additional payload: %s", payload)
				}
			case messaging.ClientSideErrorResponse:
				fmt.Printf("Client-side error: %s\n", err.Error())
			}

		case event := <-eventsChannel:
			fmt.Printf("%s event on",
				messaging.StringConnectionAction(event.Action))

			switch event.Type {
			case messaging.ChannelResponse:
				fmt.Printf("%s channel", event.Channel)
			case messaging.ChannelGroupResponse:
				fmt.Printf("%s channel group", event.Source)
			case messaging.WildcardResponse:
				fmt.Printf("%s wildcard channel", event.Source)
			}

		case <-messaging.SubscribeTimeout():
			fmt.Printf("Subscirbe request timeout")
		}
	}
}

func populateGroup(pubnub *messaging.Pubnub, group, channels string) {
	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	pubnub.ChannelGroupAddChannel(group, channels, successChannel, errorChannel)

	select {
	case <-successChannel:
	case <-errorChannel:
	}
}
