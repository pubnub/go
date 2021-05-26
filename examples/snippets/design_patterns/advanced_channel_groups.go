package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	pubnub "github.com/pubnub/go/v5"
)

var pn *pubnub.PubNub

func init() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "sub-c-b9ab9508-43cf-11e8-9967-869954283fb4"
	config.PublishKey = "pub-c-1bd448ed-05ba-4dbc-81a5-7d6ff5c6e2bb"

	pn = pubnub.NewPubNub(config)
}

func main() {
	data := map[string]string{}

	data["author"] = "user-a"
	data["status"] = "I am reading about Advanced Channel Groups!"
	data["timestamp"] = strconv.Itoa(int(time.Now().Unix()))

	publishRes, _, err := pn.Publish().
		Message(data).
		Channel("ch-user-a-status").
		Execute()

	if err != nil {
		fmt.Println("error :", err)
	}

	fmt.Println(publishRes.Timestamp)

	_, _, err = pn.AddChannelToChannelGroup().
		Channels([]string{"ch-user-a-present"}).
		ChannelGroup("cg-user-a-friends").
		Execute()

	if err != nil {
		fmt.Println("Operation failed: ", err)
	}

	fmt.Println("Channel added to channel group")

	_, _, err = pn.AddChannelToChannelGroup().
		Channels([]string{"ch-user-a-present"}).
		ChannelGroup("cg-user-a-status-feed").
		Execute()

	if err != nil {
		fmt.Println("Operation failed: ", err)
	}

	fmt.Println("Channel added to channel group")

	// ************************************
	// * User A and User B become friends
	// ************************************

	// Add User B to User A's groups: Add ch-user-b-present to cg-user-a-friends
	_, _, err = pn.AddChannelToChannelGroup().
		Channels([]string{"ch-user-b-present"}).
		ChannelGroup("cg-user-a-friends").
		Execute()

	if err != nil {
		fmt.Println("Operation failed: ", err)
	}

	fmt.Println("Channel added to channel group")

	// Add User B to User A's groups: ch-user-b-status to cg-user-a-status-feed
	_, _, err = pn.AddChannelToChannelGroup().
		Channels([]string{"ch-user-b-status"}).
		ChannelGroup("cg-user-a-status-feed").
		Execute()

	if err != nil {
		fmt.Println("Operation failed: ", err)
	}

	fmt.Println("Channel added to channel group")

	// Add User A to User B's groups: Add ch-user-a-present to cg-user-b-friends
	_, _, err = pn.AddChannelToChannelGroup().
		Channels([]string{"ch-user-a-present"}).
		ChannelGroup("cg-user-b-friends").
		Execute()

	if err != nil {
		fmt.Println("Operation failed: ", err)
	}

	fmt.Println("Channel added to channel group")

	// Add User B to User A's groups: ch-user-a-status to cg-user-b-status-feed
	_, _, err = pn.AddChannelToChannelGroup().
		Channels([]string{"ch-user-a-status"}).
		ChannelGroup("cg-user-b-status-feed").
		Execute()

	if err != nil {
		fmt.Println("Operation failed: ", err)
	}

	fmt.Println("Channel added to channel group")

	// Get the List of Friends
	listChannelsRes, _, err := pn.ListChannelsInChannelGroup().
		ChannelGroup("cg-user-a-friends").
		Execute()

	if err != nil {
		fmt.Println("Operation failed ", err)
	}

	fmt.Println("FRIEND LIST: ")
	for k, v := range listChannelsRes.Channels {
		fmt.Println("channel: ", k, v)
	}

	// Which Friends are online right now
	hereNowRes, _, err := pn.HereNow().
		ChannelGroups([]string{"cg-user-a-friends"}).
		Execute()

	if err != nil {
		fmt.Println("Operation failed ", err)
	}

	fmt.Println("ONLINE NOW: ", hereNowRes.TotalOccupancy)

	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case <-listener.Status:
			case <-listener.Message:
			case presence := <-listener.Presence:
				log.Println("PRESENCE ", presence)
			}
		}
	}()

	// Watch Friends come online / go offline
	pn.Subscribe().
		Channels([]string{"cg-user-a-friends"}).
		WithPresence(true).
		Execute()
}
