package main

import (
	"bufio"
	"errors"
	"fmt"
	//"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	pubnub "github.com/pubnub/go"
)

var config *pubnub.Config
var pn *pubnub.PubNub
var quitSubscribe = false

const outputPrefix = "\x1b[32;1m Example >>>> \x1b[0m"
const outputSuffix = "\x1b[32;2m Example <<<< \x1b[0m"

func main() {
	config = pubnub.NewConfig()
	config.PNReconnectionPolicy = pubnub.PNLinearPolicy
	//config.EnableLogging = false

	var infoLogger *log.Logger

	logfileName := "pubnubMessaging.log"
	f, err := os.OpenFile(logfileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {

		fmt.Println("error opening file: ", err.Error())
		fmt.Println("Logging disabled")
	} else {
		fmt.Println("Logging enabled writing to ", logfileName)
		infoLogger = log.New(f, "", log.Ldate|log.Ltime|log.Lshortfile)
	}
	//config.Log = log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile)
	//config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	config.Log = infoLogger
	config.Log.SetPrefix("PubNub :->  ")
	//config.SuppressLeaveEvents = true

	config.PublishKey = "pub-c-4f1dbd79-ab94-487d-b779-5881927db87c"
	config.SubscribeKey = "sub-c-f2489488-2dbd-11e8-a27a-a2b5bab5b996"
	config.SecretKey = "sec-c-NjlmYzVkMjEtOWIxZi00YmJlLThjZDktMjI4NGQwZDUxZDQ0"
	//config.CipherKey = "enigma"
	pn = pubnub.NewPubNub(config)

	// for subscribe event
	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				fmt.Print(fmt.Sprintf("%s Subscribe Response:", outputPrefix))
				fmt.Println(" --- STATUS: ")
				fmt.Println(fmt.Sprintf("%s status.Error %s", outputPrefix, status.Error))
				fmt.Println(fmt.Sprintf("%s status.Category %s", outputPrefix, status.Category))
				fmt.Println(fmt.Sprintf("%s status.Operation %s", outputPrefix, status.Operation))
				fmt.Println(fmt.Sprintf("%s status.StatusCode %d", outputPrefix, status.StatusCode))
				fmt.Println(fmt.Sprintf("%s %s", outputPrefix, status.ErrorData))
				fmt.Println(fmt.Sprintf("%s %s", outputPrefix, status.ClientRequest))
				fmt.Println("")
				fmt.Println(fmt.Sprintf("%s", outputSuffix))
			case msg := <-listener.Message:
				fmt.Print(fmt.Sprintf("%s Subscribe Response:", outputPrefix))
				fmt.Println(" --- MESSAGE: ")
				fmt.Println(fmt.Sprintf("%s msg.Channel: %s", outputPrefix, msg.Channel))
				fmt.Println(fmt.Sprintf("%s msg.Message: %s", outputPrefix, msg.Message))
				fmt.Println(fmt.Sprintf("%s msg.SubscribedChannel: %s", outputPrefix, msg.SubscribedChannel))
				fmt.Println(fmt.Sprintf("%s msg.Timetoken: %d", outputPrefix, msg.Timetoken))
				fmt.Println("")
				fmt.Println(fmt.Sprintf("%s", outputSuffix))
			case presence := <-listener.Presence:
				fmt.Print(fmt.Sprintf("%s Subscribe Response:", outputPrefix))
				fmt.Println(" --- PRESENCE: ")
				fmt.Println(fmt.Sprintf("%s %s", outputPrefix, presence))
				fmt.Println("")
				fmt.Println(fmt.Sprintf("%s", outputSuffix))
			}
		}
	}()

	pn.AddListener(listener)
	showHelp()

	/*config2 := pubnub.NewConfig()
	config2.PublishKey = "pub-c-c6a4792f-af77-4028-88b4-5995da3aa7b4"
	config2.SubscribeKey = "sub-c-4cf48c6c-2025-11e8-b192-4eac351dc434"
	config2.SubscribeRequestTimeout = 59
	config2.Uuid = "GlobalSubscriber"
	config2.PNReconnectionPolicy = pubnub.PNLinearPolicy
	config2.Log = infoLogger
	config2.Log.SetPrefix("PubNub2:")

	pn2 := pubnub.NewPubNub(config2)
	pn2.AddListener(listener)
	channel := "ch1"

	pn2.Subscribe().Channels([]string{channel}).Execute()*/

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(fmt.Sprintf("%s ", outputPrefix))
		text, _ := reader.ReadString('\n')

		text = text[:len(text)-1]

		if len(text) != 0 {
			readCommand(text)
		}
		fmt.Println("")
	}

}

func showErr(err string) {
	fmt.Println(fmt.Sprintf("%s \x1b[31;1m %s \x1b[0m", outputPrefix, errors.New(err)))
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
	showDelMessagesHelp()
	showWhereNowHelp()
	showUnsubscribeHelp()
	showFetchHelp()
	showFireHelp()
	fmt.Println("\n ================")
	fmt.Println(" ||  COMMANDS  ||")
	fmt.Println(" ================\n")
	fmt.Println(" UNSUBSCRIBE ALL \n\tq ")
	fmt.Println(" QUIT \n\tctrl+c ")
}

func showFetchHelp() {
	fmt.Println(" FETCH EXAMPLE: ")
	fmt.Println("	fetch Channel Reverse Max Start End ")
	fmt.Println("	fetch my-channel,test true 10 15210190573608384 15211140747622125 ")
}

func showFireHelp() {
	fmt.Println(" FIRE EXAMPLE: ")
	fmt.Println("	fire usePost \"my-message\" my-channel")
	fmt.Println("	fire false \"my-message\" my-channel")
}

func showPublishHelp() {
	fmt.Println(" PUBLISH EXAMPLE: ")
	fmt.Println("	pub usePost store noreplicate \"my-message\" my-channel")
	fmt.Println("	pub false true false \"my-message\" my-channel")
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
	fmt.Println("	hist my-channel true true 10 15210190573608384 15211140747622125 ")
}

func showDelMessagesHelp() {
	fmt.Println(" Delete Messages EXAMPLE: ")
	fmt.Println("	delmessages Channel Start End ")
	fmt.Println("	delmessages my-channel 15210190573608384 15211140747622125 ")
}

func showWhereNowHelp() {
	fmt.Println(" WHERENOW EXAMPLE: ")
	fmt.Println("	wherenow uuid ")
	fmt.Println("	wherenow \"uuidToCheck\"")
}

func showUnsubscribeHelp() {
	fmt.Println(" UNSUBSCRIBE EXAMPLE: ")
	fmt.Println("	unsub channels channelGroups")
	fmt.Println("	unsub my-channel,my-another-channel my-channelgroup,my-another-channel-group")

}

func readCommand(cmd string) {
	command := strings.Split(cmd, " ")

	switch w := command[0]; w {
	case "pub":
		publishRequest(command[1:])
	case "fire":
		fireRequest(command[1:])
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
	case "unsub":
		unsubscribeRequest(command[1:])
	case "fetch":
		fetchRequest(command[1:])
	case "delmessages":
		delMessageRequest(command[1:])
	case "setState":
		subscribeRequest(command[1:])
	/*case "getState":
		subscribeRequest(command[1:])
	case "addChCg:
		subscribeRequest(command[1:])
	case "remChCg":
		subscribeRequest(command[1:])
	case "listChCg":
		subscribeRequest(command[1:])
	case "delCg":
		subscribeRequest(command[1:])
	case "grant":
		subscribeRequest(command[1:])*/
	case "help":
		showHelp()
	case "q":
		pn.UnsubscribeAll()
	default:
		showHelp()
	}
}

func delMessageRequest(args []string) {
	if len(args) == 0 {
		showDelMessagesHelp()
		return
	}

	var channel string
	if len(args) > 0 {
		channel = args[0]
	}

	var start int64
	if len(args) > 1 {
		i, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			i = 0
		} else {
			start = i
		}
	}

	var end int64
	if len(args) > 2 {
		i, err := strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			i = 0
		} else {
			end = i
		}
	}

	if (end != 0) && (start != 0) {
		res, status, err := pn.DeleteMessages().Channel(channel).End(end).Start(start).Execute()
		fmt.Println(res, status, err)
	} else if start != 0 {
		res, status, err := pn.DeleteMessages().Channel(channel).Start(start).Execute()
		fmt.Println(res, status, err)
	} else if end != 0 {
		res, status, err := pn.DeleteMessages().Channel(channel).End(end).Execute()
		fmt.Println(res, status, err)
	} else {
		res, status, err := pn.DeleteMessages().Channel(channel).Execute()
		fmt.Println(res, status, err)
	}
	fmt.Println(fmt.Sprintf("%s", outputSuffix))

}

func whereNowRequest(args []string) {
	uuidToUse := ""
	if len(args) > 0 {
		uuidToUse = args[0]
	}

	fmt.Println(fmt.Sprintf("%s whereNowRequest:", outputPrefix))
	if len(uuidToUse) == 0 {
		res, status, err := pn.WhereNow().Execute()
		fmt.Println(res, status, err)
	} else {
		res, status, err := pn.WhereNow().Uuid(uuidToUse).Execute()
		fmt.Println(res, status, err)
	}
	fmt.Println(fmt.Sprintf("%s", outputSuffix))
}

func fetchRequest(args []string) {
	if len(args) == 0 {
		showFetchHelp()
		return
	}

	var channels []string
	if len(args) > 0 {
		channels = strings.Split(args[0], ",")
	}

	var reverse bool
	if len(args) > 1 {
		reverse, _ = strconv.ParseBool(args[1])
	}

	var count int
	if len(args) > 2 {
		i, err := strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			i = 100
		} else {
			count = int(i)
		}
	}

	var start int64
	if len(args) > 3 {
		i, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			i = 0
		} else {
			start = i
		}
	}

	var end int64
	if len(args) > 4 {
		i, err := strconv.ParseInt(args[4], 10, 64)
		if err != nil {
			i = 0
		} else {
			end = i
		}
	}

	if (end != 0) && (start != 0) {
		res, status, err := pn.Fetch().
			Channels(channels).
			Count(count).
			Start(start).
			End(end).
			Reverse(reverse).
			Execute()
		ParseFetch(res, status, err)
	} else if start != 0 {
		res, status, err := pn.Fetch().
			Channels(channels).
			Count(count).
			Start(start).
			Reverse(reverse).
			Execute()
		ParseFetch(res, status, err)
	} else if end != 0 {
		res, status, err := pn.Fetch().
			Channels(channels).
			Count(count).
			End(end).
			Reverse(reverse).
			Execute()
		ParseFetch(res, status, err)
	} else {
		res, status, err := pn.Fetch().
			Channels(channels).
			Count(count).
			Reverse(reverse).
			Execute()
		ParseFetch(res, status, err)
	}
}

func ParseFetch(res *pubnub.FetchResponse, status pubnub.StatusResponse, err error) {
	fmt.Println(fmt.Sprintf("%s ParseFetch:", outputPrefix))
	if status.Error == nil {
		for channel, messages := range res.Messages {
			fmt.Println("channel", channel)
			for _, messageInt := range messages {
				message := pubnub.FetchResponseItem(messageInt)
				fmt.Println(message.Message)
				fmt.Println(message.Timetoken)
			}
		}
	} else {
		fmt.Println("ParseFetch", err)
		fmt.Println("ParseFetch", status.StatusCode)
	}
	fmt.Println(fmt.Sprintf("%s", outputSuffix))
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
	fmt.Println(fmt.Sprintf("%s ParseHistory:", outputPrefix))
	for _, v := range res.Messages {
		fmt.Println(fmt.Sprintf("%s Timetoken %d", outputPrefix, v.Timetoken))
		fmt.Println(fmt.Sprintf("%s Message %s", outputPrefix, v.Message))
	}
	fmt.Println(fmt.Sprintf("%s EndTimetoken %d", outputPrefix, res.EndTimetoken))
	fmt.Println(fmt.Sprintf("%s StartTimetoken %d", outputPrefix, res.StartTimetoken))
	fmt.Println(fmt.Sprintf("%s", outputSuffix))
}

func timeRequest() {

	res, status, err := pn.Time().Execute()
	fmt.Println(fmt.Sprintf("%s timeResponse:", outputPrefix))
	fmt.Println(res, status, err)
	fmt.Println(fmt.Sprintf("%s", outputSuffix))
}

func hereNowResponse(res *pubnub.HereNowResponse, status pubnub.StatusResponse, err error) {
	fmt.Println(fmt.Sprintf("%s hereNowResponse:", outputPrefix))
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
	fmt.Println(fmt.Sprintf("%s", outputSuffix))
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
	if len(args) < 5 {
		showErr("channels or message not found")
		showPublishHelp()
		return
	}

	usePost, _ := strconv.ParseBool(args[0])
	var store bool
	if len(args) > 1 {
		store, _ = strconv.ParseBool(args[1])
	}
	var repl bool
	if len(args) > 2 {
		repl, _ = strconv.ParseBool(args[2])
	}

	message := args[3]
	reg := regexp.MustCompile(`"([^"]*)"`)
	res := reg.ReplaceAllString(message, "${1}")

	if res == "" {
		showErr("Empty message!")
		return
	}

	channels := strings.Split(args[4], ",")

	for _, ch := range channels {
		fmt.Println(fmt.Sprintf("%s Publishing to channel: %s", outputPrefix, ch))
		res, status, err := pn.Publish().
			Channel(ch).
			Message(res).
			UsePost(usePost).
			ShouldStore(store).
			DoNotReplicate(repl).
			Execute()

		if err != nil {
			showErr("Error while publishing: " + err.Error())
		}

		fmt.Println(fmt.Sprintf("%s Publish Response:", outputPrefix))

		fmt.Println(fmt.Sprintf("%%s %s", res, status))
		fmt.Println(fmt.Sprintf("%s", outputSuffix))
	}
}

func fireRequest(args []string) {
	if len(args) < 3 {
		showErr("channels or message not found")
		showFireHelp()
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

	channels := strings.Split(args[2], ",")

	for _, ch := range channels {
		fmt.Println(fmt.Sprintf("%s Publishing to channel: %s", outputPrefix, ch))
		res, status, err := pn.Fire().
			Channel(ch).
			Message(res).
			UsePost(usePost).
			Ttl(1).
			Execute()

		if err != nil {
			showErr("Error while publishing: " + err.Error())
		}

		fmt.Println(fmt.Sprintf("%s Publish Response:", outputPrefix))

		fmt.Println(fmt.Sprintf("%%s %s", res, status))
		fmt.Println(fmt.Sprintf("%s", outputSuffix))
	}
}

func unsubscribeRequest(args []string) {
	if len(args) == 0 {
		showUnsubscribeHelp()
		return
	}

	channels := strings.Split(args[0], ",")
	if (len(args)) > 2 {
		groups := strings.Split(args[1], ",")
		pn.Unsubscribe().
			Channels(channels).
			ChannelGroups(groups).
			Execute()
	} else {
		pn.Unsubscribe().
			Channels(channels).
			Execute()
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
