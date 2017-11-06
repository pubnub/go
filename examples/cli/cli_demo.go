package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	pubnub "github.com/pubnub/go"
)

var config = pubnub.NewConfig()
var pn = pubnub.NewPubNub(config)
var quitSubscribe = false

func main() {
	config.PublishKey = "pub-c-071e1a3f-607f-4351-bdd1-73a8eb21ba7c"
	config.SubscribeKey = "sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f"

	// for subscribe event
	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				fmt.Println(" --- STATUS: ")
				fmt.Println(status)
				fmt.Print(">>> ")
			case msg := <-listener.Message:
				fmt.Println(" --- MESSAGE: ")
				fmt.Println(msg)
				fmt.Print(">>> ")
			case presence := <-listener.Presence:
				fmt.Println(" --- PRESENCE: ")
				fmt.Println(presence)
				fmt.Print(">>> ")
			}
		}
	}()

	pn.AddListener(listener)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">>> ")
		text, _ := reader.ReadString('\n')

		text = text[:len(text)-1]

		if len(text) != 0 {
			readCommand(text)
		}
	}
}

func showErr(err string) {
	fmt.Println("\x1b[31;1m", errors.New(err), "\x1b[0m")
}

func showHelp() {
	fmt.Println("")
	fmt.Println("\n ============")
	fmt.Println(" ||  HELP  ||")
	fmt.Println(" ============\n")
	showPublishHelp()
	showSubscribeHelp()
	fmt.Println("\n ================")
	fmt.Println(" ||  COMMANDS  ||")
	fmt.Println(" ================\n")
	fmt.Println(" UNSUBSCRIBE \n\tq ")
	fmt.Println(" QUIT \n\tctrl+c ")
}

func showPublishHelp() {
	fmt.Println(" PUBLISH EXAMPLE: ")
	fmt.Println(" 	pub \"my-message\" my-channel,my-another-channel")
}

func showSubscribeHelp() {
	fmt.Println(" SUBSCRIBE EXAMPLE: ")
	fmt.Println(" 	sub my-channel,my-another-channel my-channelgroup,my-another-channel-group")
}

func readCommand(cmd string) {
	command := strings.Split(cmd, " ")

	switch w := command[0]; w {
	case "pub":
		publishRequest(command[1:])
	case "sub":
		subscribeRequest(command[1:])
	case "help":
		showHelp()
	case "q":
		pn.UnsubscribeAll()
	default:
		showHelp()
	}
}

func publishRequest(args []string) {
	if len(args) == 0 {
		showPublishHelp()
		return
	}

	message := args[0]
	reg := regexp.MustCompile(`"([^"]*)"`)
	res := reg.ReplaceAllString(message, "${1}")

	if res == "" {
		showErr("Empty message!")
		return
	}

	if len(args) != 2 {
		showErr("Not found channels or message")
		return
	}

	channels := strings.Split(args[1], ",")

	for _, ch := range channels {
		fmt.Println("Publishing to channel: ", ch)
		res, status, err := pn.Publish().
			Channel(ch).
			Message(res).
			Execute()

		if err != nil {
			showErr("Error while publishing: " + err.Error())
		}

		fmt.Println(res, status)
	}
}

func subscribeRequest(args []string) {
	if len(args) == 0 {
		showSubscribeHelp()
		return
	}

	channels := strings.Split(args[0], ",")
	groups := strings.Split(args[1], ",")

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels:      channels,
		ChannelGroups: groups,
	})
}
