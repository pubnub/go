package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

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
				fmt.Println(status.Error)
				fmt.Println(status.ErrorData)
				fmt.Println(status.ClientRequest)
				fmt.Print(">>> ")
			case msg := <-listener.Message:
				fmt.Println(" --- MESSAGE: ")
				fmt.Println(fmt.Sprintf("msg.Channel: %s", msg.Channel))
				fmt.Println(fmt.Sprintf("msg.Message: %s", msg.Message))
				fmt.Println(fmt.Sprintf("msg.SubscribedChannel: %s", msg.SubscribedChannel))
				fmt.Println(fmt.Sprintf("msg.Timetoken: %d", msg.Timetoken))
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
		fmt.Println("")
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
	showTimeHelp()
	showHereNowHelp()
	showHistoryHelp()
	showWhereNowHelp()
	fmt.Println("\n ================")
	fmt.Println(" ||  COMMANDS  ||")
	fmt.Println(" ================\n")
	fmt.Println(" UNSUBSCRIBE \n\tq ")
	fmt.Println(" QUIT \n\tctrl+c ")
}

func showPublishHelp() {
	fmt.Println(" PUBLISH EXAMPLE: ")
	fmt.Println("	pub usePost \"my-message\" my-channel,my-another-channel")
	fmt.Println("	pub false \"my-message\" my-channel,my-another-channel")
}

func showTimeHelp() {
	fmt.Println(" TIME EXAMPLE: ")
	fmt.Println("	time")
}

func showHereNowHelp() {
	fmt.Println(" HERENOW EXAMPLE: ")
	fmt.Println("	herenow includeState includeUUIDs channel channel-group")
	fmt.Println("	herenow false false my-channel my-channel-group")
}

func showSubscribeHelp() {
	fmt.Println(" SUBSCRIBE EXAMPLE: ")
	fmt.Println("	sub withPresence channels channelGroups")
	fmt.Println("	sub true my-channel,my-another-channel my-channelgroup,my-another-channel-group")
}

func showHistoryHelp() {
	fmt.Println(" HISTORY EXAMPLE: ")
	fmt.Println("	hist Channel IncludeTimetoken Reverse Count Start End ")
	fmt.Println("	hist test true true 10 15210190573608384 15211140747622125 ")
}

func showWhereNowHelp() {
	fmt.Println(" WHERENOW EXAMPLE: ")
	fmt.Println("	wherenow uuid ")
	fmt.Println("	wherenow \"uuidToCheck\"")
}

func readCommand(cmd string) {
	command := strings.Split(cmd, " ")

	switch w := command[0]; w {
	case "pub":
		publishRequest(command[1:])
	case "sub":
		subscribeRequest(command[1:])
	case "time":
		timeRequest()
	case "herenow":
		hereNowRequest(command[1:])
	case "hist":
		historyRequest(command[1:])
	case "wherenow":
		whereNowRequest(command[1:])
	/*case "fetch":
		unsubscribeRequest(command[1:])
	case "delmessage":
		unsubscribeRequest(command[1:])
	case "usub":
		unsubscribeRequest(command[1:])
	case "setState":
		subscribeRequest(command[1:])
	case "getState":
		subscribeRequest(command[1:])
	case "addChCg:
		subscribeRequest(command[1:])
	case "remChCg":
		subscribeRequest(command[1:])
	case "listChCg":
		subscribeRequest(command[1:])
	case "delCg":
		subsc`ribeRequest(command[1:])
	case "grant":
		subsc`ribeRequest(command[1:])*/
	case "help":
		showHelp()
	case "q":
		pn.UnsubscribeAll()
	default:
		showHelp()
	}
}

func whereNowRequest(args []string) {
	uuidToUse := ""
	if len(args) > 0 {
		uuidToUse = args[0]
	}

	if len(uuidToUse) == 0 {
		res, status, err := pn.WhereNow().Execute()
		fmt.Println(res, status, err)
	} else {
		res, status, err := pn.WhereNow().Uuid(uuidToUse).Execute()
		fmt.Println(res, status, err)
	}
}

func historyRequest(args []string) {
	if len(args) == 0 {
		showHistoryHelp()
		return
	}

	var channel string
	if len(args) > 0 {
		channel = args[0]
	}

	var includeTimetoken bool
	if len(args) > 1 {
		includeTimetoken, _ = strconv.ParseBool(args[1])
	}

	var reverse bool
	if len(args) > 2 {
		reverse, _ = strconv.ParseBool(args[2])
	}

	var count int
	if len(args) > 3 {
		i, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			i = 100
		} else {
			count = int(i)
		}
	}

	var start int64
	if len(args) > 4 {
		i, err := strconv.ParseInt(args[4], 10, 64)
		if err != nil {
			i = 0
		} else {
			start = i
		}
	}

	var end int64
	if len(args) > 5 {
		i, err := strconv.ParseInt(args[5], 10, 64)
		if err != nil {
			i = 0
		} else {
			end = i
		}
	}

	if (end != 0) && (start != 0) {
		res, status, err := pn.History().
			Channel(channel).
			Count(count).
			Start(start).
			End(end).
			IncludeTimetoken(includeTimetoken).
			Reverse(reverse).
			Execute()
		ParseHistory(res, status, err)
	} else if start != 0 {
		res, status, err := pn.History().
			Channel(channel).
			Count(count).
			Start(start).
			IncludeTimetoken(includeTimetoken).
			Reverse(reverse).
			Execute()
		ParseHistory(res, status, err)
	} else if end != 0 {
		res, status, err := pn.History().
			Channel(channel).
			Count(count).
			End(end).
			IncludeTimetoken(includeTimetoken).
			Reverse(reverse).
			Execute()
		ParseHistory(res, status, err)
	} else {
		res, status, err := pn.History().
			Channel(channel).
			Count(count).
			IncludeTimetoken(includeTimetoken).
			Reverse(reverse).
			Execute()
		ParseHistory(res, status, err)
	}
}

func ParseHistory(res *pubnub.HistoryResponse, status pubnub.StatusResponse, err error) {
	for _, v := range res.Messages {
		fmt.Println(fmt.Sprintf("Timetoken %d", v.Timetoken))
		fmt.Println(fmt.Sprintf("Message %s", v.Message))
	}
	fmt.Println(fmt.Sprintf("EndTimetoken %d", res.EndTimetoken))
	fmt.Println(fmt.Sprintf("StartTimetoken %d", res.StartTimetoken))
}

func timeRequest() {
	fmt.Println(time.Now())
	res, status, err := pn.Time().Execute()
	fmt.Println(time.Now())
	fmt.Println(res, status, err)
}

func hereNowResponse(res *pubnub.HereNowResponse, status pubnub.StatusResponse, err error) {
	fmt.Println(res, status, err)
	for _, v := range res.Channels {
		fmt.Println(v.ChannelName)
		fmt.Println(v.Occupancy)
		fmt.Println(v.Occupants)

		for _, v := range v.Occupants {
			fmt.Println(v.Uuid)
		}
	}
	fmt.Println(res.TotalChannels)
	fmt.Println(res.TotalOccupancy)

}

func hereNowRequest(args []string) {
	if len(args) == 0 {
		res, status, err := pn.HereNow().Execute()
		hereNowResponse(res, status, err)
		return
	}
	var includeState bool
	if len(args) > 0 {
		includeState, _ = strconv.ParseBool(args[0])
	}

	var includeUUIDs bool
	if len(args) > 1 {
		includeUUIDs, _ = strconv.ParseBool(args[1])
	}

	var channels []string
	if len(args) > 2 {
		if len(args[2]) != 0 {
			channels = strings.Split(args[2], ",")
		}
	}

	var channelGroups []string
	if len(args) > 3 {
		if len(args[3]) != 0 {
			channelGroups = strings.Split(args[3], ",")
		}
	}

	if (len(channels) != 0) && (len(channelGroups) != 0) {
		res, status, err := pn.HereNow().Channels(channels).ChannelGroups(channelGroups).IncludeState(includeState).IncludeUuids(includeUUIDs).Execute()
		hereNowResponse(res, status, err)
	} else if len(channels) != 0 {
		res, status, err := pn.HereNow().Channels(channels).IncludeState(includeState).IncludeUuids(includeUUIDs).Execute()
		hereNowResponse(res, status, err)
	} else if len(channelGroups) != 0 {
		res, status, err := pn.HereNow().ChannelGroups(channelGroups).IncludeState(includeState).IncludeUuids(includeUUIDs).Execute()
		hereNowResponse(res, status, err)
	} else {
		res, status, err := pn.HereNow().IncludeState(includeState).IncludeUuids(includeUUIDs).Execute()
		hereNowResponse(res, status, err)
	}
}

func publishRequest(args []string) {
	if len(args) == 0 {
		showPublishHelp()
		return
	}

	usePost, _ := strconv.ParseBool(args[0])

	message := args[1]
	reg := regexp.MustCompile(`"([^"]*)"`)
	res := reg.ReplaceAllString(message, "${1}")

	if res == "" {
		showErr("Empty message!")
		return
	}

	if len(args) != 3 {
		showErr("Not found channels or message")
		return
	}

	channels := strings.Split(args[2], ",")

	for _, ch := range channels {
		fmt.Println("Publishing to channel: ", ch)
		res, status, err := pn.Publish().
			Channel(ch).
			Message(res).
			UsePost(usePost).
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

	withPresence, _ := strconv.ParseBool(args[0])

	channels := strings.Split(args[1], ",")
	if (len(args)) > 2 {
		groups := strings.Split(args[2], ",")
		pn.Subscribe().
			Channels(channels).
			ChannelGroups(groups).
			WithPresence(withPresence).
			Execute()
	} else {
		pn.Subscribe().
			Channels(channels).
			WithPresence(withPresence).
			Execute()
	}

}
