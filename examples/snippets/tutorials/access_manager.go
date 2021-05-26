package main

import (
	"fmt"

	pubnub "github.com/pubnub/go/v5"
)

func operationLevel() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "my_sub_key"
	config.PublishKey = "my_pub_key"
	config.SecretKey = "my_secret_key"
	config.Secure = false

	pn := pubnub.NewPubNub(config)

	res, status, err := pn.Grant().
		Read(false).
		Write(false).
		Manage(false).
		TTL(60).
		Execute()

	if err != nil {
		fmt.Println(err)
		// handle error
	}

	fmt.Println(res, status)
}

func operationLevel2() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "my_sub_key"
	config.PublishKey = "my_pub_key"
	config.SecretKey = "my_secret_key"

	pn := pubnub.NewPubNub(config)

	res, status, err := pn.Grant().
		Read(false).
		Write(false).
		Manage(false).
		Execute()

	if err != nil {
		fmt.Println(err)
		// handle error
	}

	fmt.Println(res, status)
}

func operationLevel3() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "my_sub_key"
	config.PublishKey = "my_pub_key"
	config.SecretKey = "my_secret_key"

	pn := pubnub.NewPubNub(config)

	res, status, err := pn.Grant().
		Read(true).
		Write(false).
		Channels([]string{"public_chat"}).
		TTL(60).
		Execute()

	if err != nil {
		fmt.Println(err)
		// handle error
	}

	fmt.Println(res, status)
}

func operationLevel4() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "my_sub_key"
	config.PublishKey = "my_pub_key"
	config.SecretKey = "my_secret_key"

	pn := pubnub.NewPubNub(config)

	res, status, err := pn.Grant().
		Read(true).
		Write(true).
		Channels([]string{"public_chat"}).
		AuthKeys([]string{"auth_keys"}).
		TTL(60).
		Execute()

	if err != nil {
		fmt.Println(err)
		// handle error
	}

	fmt.Println(res, status)
}

func permissionDenied() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "sub-c-b9ab9508-43cf-11e8-9967-869954283fb4"
	config.PublishKey = "pub-c-1bd448ed-05ba-4dbc-81a5-7d6ff5c6e2bb"
	config.SecretKey = "wrong-key"

	pn := pubnub.NewPubNub(config)
	doneAccessDenied := make(chan bool)

	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNAccessDeniedCategory:
					doneAccessDenied <- true
				}
			case <-listener.Message:
			case <-listener.Presence:
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().
		Channels([]string{"private_chat"}).
		Execute()

	<-doneAccessDenied
}

func grantChannelGroup() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "my_sub_key"
	config.PublishKey = "my_pub_key"
	config.SecretKey = "my_secret_key"

	pn := pubnub.NewPubNub(config)

	res, status, err := pn.Grant().
		Read(true).
		Write(false).
		ChannelGroups([]string{"gr1", "gr2", "gr3"}).
		AuthKeys([]string{"key1", "key2", "key3"}).
		TTL(60).
		Execute()

	if err != nil {
		fmt.Println(err)
		// handle error
	}

	fmt.Println(res, status)
}

func revokeChannelGroup() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "my_sub_key"
	config.PublishKey = "my_pub_key"
	config.SecretKey = "my_secret_key"

	pn := pubnub.NewPubNub(config)

	res, status, err := pn.Grant().
		Read(false).
		Write(false).
		Manage(false).
		ChannelGroups([]string{"gr1", "gr2", "gr3"}).
		AuthKeys([]string{"key1", "key2", "key3"}).
		TTL(60).
		Execute()

	if err != nil {
		fmt.Println(err)
		// handle error
	}

	fmt.Println(res, status)
}

func cipher() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "my_sub_key"
	config.PublishKey = "my_pub_key"
	config.SecretKey = "my_secret_key"

	pn := pubnub.NewPubNub(config)

	_ = pn
}

func main() {
	operationLevel()
	operationLevel2()
	operationLevel3()
	operationLevel4()
	permissionDenied()
	grantChannelGroup()
	revokeChannelGroup()
	cipher()
}
