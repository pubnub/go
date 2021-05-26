package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	//"io/ioutil"
	"encoding/json"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	pubnub "github.com/pubnub/go/v5"

	"net/http"
	_ "net/http/pprof"
)

var config *pubnub.Config
var pn *pubnub.PubNub
var quitSubscribe = false

const outputPrefix = "\x1b[32;1m Example >>>> \x1b[0m"
const outputSuffix = "\x1b[32;2m Example <<<< \x1b[0m"

func main() {
	connect()
	// go pubnub.NewPubNub(pubnub.NewConfig())
	// pubnub.NewPubNub(pubnub.NewConfig())
}

func connect() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	config = pubnub.NewConfig()
	config.UseHTTP2 = false

	config.PNReconnectionPolicy = pubnub.PNExponentialPolicy

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

	config.Log = infoLogger
	config.Log.SetPrefix("PubNub :->  ")
	config.PublishKey = "demo"
	config.SubscribeKey = "demo"

	config.CipherKey = "enigma"
	config.UseRandomInitializationVector = true

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
			case msg := <-listener.Signal:
				fmt.Print(fmt.Sprintf("%s Subscribe Response:", outputPrefix))
				fmt.Println(" --- SIGNAL: ")
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
			case uuidEvent := <-listener.UUIDEvent:
				fmt.Print(fmt.Sprintf("%s Subscribe Response:", outputPrefix))
				fmt.Println(" --- UUIDEvent: ")
				fmt.Println(fmt.Sprintf("%s %s", outputPrefix, uuidEvent))
				fmt.Println(fmt.Sprintf("%s uuidEvent.Channel: %s", outputPrefix, uuidEvent.Channel))
				fmt.Println(fmt.Sprintf("%s uuidEvent.SubscribedChannel: %s", outputPrefix, uuidEvent.SubscribedChannel))
				fmt.Println(fmt.Sprintf("%s uuidEvent.Event: %s", outputPrefix, uuidEvent.Event))
				fmt.Println(fmt.Sprintf("%s uuidEvent.UUID: %s", outputPrefix, uuidEvent.UUID))
				fmt.Println(fmt.Sprintf("%s uuidEvent.Description: %s", outputPrefix, uuidEvent.Description))
				fmt.Println(fmt.Sprintf("%s uuidEvent.Timestamp: %s", outputPrefix, uuidEvent.Timestamp))
				fmt.Println(fmt.Sprintf("%s uuidEvent.Name: %s", outputPrefix, uuidEvent.Name))
				fmt.Println(fmt.Sprintf("%s uuidEvent.ExternalID: %s", outputPrefix, uuidEvent.ExternalID))
				fmt.Println(fmt.Sprintf("%s uuidEvent.ProfileURL: %s", outputPrefix, uuidEvent.ProfileURL))
				fmt.Println(fmt.Sprintf("%s uuidEvent.Email: %s", outputPrefix, uuidEvent.Email))
				// fmt.Println(fmt.Sprintf("%s uuidEvent.Created: %s", outputPrefix, uuidEvent.Created))
				fmt.Println(fmt.Sprintf("%s uuidEvent.Updated: %s", outputPrefix, uuidEvent.Updated))
				fmt.Println(fmt.Sprintf("%s uuidEvent.ETag: %s", outputPrefix, uuidEvent.ETag))
				fmt.Println(fmt.Sprintf("%s uuidEvent.Custom: %v", outputPrefix, uuidEvent.Custom))

			case channelEvent := <-listener.ChannelEvent:
				fmt.Print(fmt.Sprintf("%s Subscribe Response:", outputPrefix))
				fmt.Println(" --- ChannelEvent: ")
				fmt.Println(fmt.Sprintf("%s %s", outputPrefix, channelEvent))
				fmt.Println(fmt.Sprintf("%s channelEvent.Channel: %s", outputPrefix, channelEvent.Channel))
				fmt.Println(fmt.Sprintf("%s channelEvent.SubscribedChannel: %s", outputPrefix, channelEvent.SubscribedChannel))
				fmt.Println(fmt.Sprintf("%s channelEvent.Event: %s", outputPrefix, channelEvent.Event))
				fmt.Println(fmt.Sprintf("%s channelEvent.Channel: %s", outputPrefix, channelEvent.Channel))
				fmt.Println(fmt.Sprintf("%s channelEvent.Description: %s", outputPrefix, channelEvent.Description))
				fmt.Println(fmt.Sprintf("%s channelEvent.Timestamp: %s", outputPrefix, channelEvent.Timestamp))
				// fmt.Println(fmt.Sprintf("%s channelEvent.Created: %s", outputPrefix, channelEvent.Created))
				fmt.Println(fmt.Sprintf("%s channelEvent.Updated: %s", outputPrefix, channelEvent.Updated))
				fmt.Println(fmt.Sprintf("%s channelEvent.ETag: %s", outputPrefix, channelEvent.ETag))
				fmt.Println(fmt.Sprintf("%s channelEvent.Custom: %v", outputPrefix, channelEvent.Custom))

			case membershipEvent := <-listener.MembershipEvent:
				fmt.Print(fmt.Sprintf("%s Subscribe Response:", outputPrefix))
				fmt.Println(" --- MembershipEvent: ")
				fmt.Println(fmt.Sprintf("%s %s", outputPrefix, membershipEvent))
				fmt.Println(fmt.Sprintf("%s membershipEvent.Channel: %s", outputPrefix, membershipEvent.Channel))
				fmt.Println(fmt.Sprintf("%s membershipEvent.SubscribedChannel: %s", outputPrefix, membershipEvent.SubscribedChannel))
				fmt.Println(fmt.Sprintf("%s membershipEvent.Event: %s", outputPrefix, membershipEvent.Event))
				fmt.Println(fmt.Sprintf("%s membershipEvent.Channel: %s", outputPrefix, membershipEvent.Channel))
				fmt.Println(fmt.Sprintf("%s membershipEvent.UUID: %s", outputPrefix, membershipEvent.UUID))
				fmt.Println(fmt.Sprintf("%s membershipEvent.Description: %s", outputPrefix, membershipEvent.Description))
				fmt.Println(fmt.Sprintf("%s membershipEvent.Timestamp: %s", outputPrefix, membershipEvent.Timestamp))
				fmt.Println(fmt.Sprintf("%s membershipEvent.Custom: %v", outputPrefix, membershipEvent.Custom))

			case messageActionsEvent := <-listener.MessageActionsEvent:
				fmt.Print(fmt.Sprintf("%s Subscribe Response:", outputPrefix))
				fmt.Println(" --- MessageActionsEvent: ")
				fmt.Println(fmt.Sprintf("%s %s", outputPrefix, messageActionsEvent))
				fmt.Println(fmt.Sprintf("%s messageActionsEvent.Channel: %s", outputPrefix, messageActionsEvent.Channel))
				fmt.Println(fmt.Sprintf("%s messageActionsEvent.SubscribedChannel: %s", outputPrefix, messageActionsEvent.SubscribedChannel))
				fmt.Println(fmt.Sprintf("%s messageActionsEvent.Event: %s", outputPrefix, messageActionsEvent.Event))
				fmt.Println(fmt.Sprintf("%s messageActionsEvent.Data.ActionType: %s", outputPrefix, messageActionsEvent.Data.ActionType))
				fmt.Println(fmt.Sprintf("%s messageActionsEvent.Data.ActionValue: %s", outputPrefix, messageActionsEvent.Data.ActionValue))
				fmt.Println(fmt.Sprintf("%s messageActionsEvent.Data.ActionTimetoken: %s", outputPrefix, messageActionsEvent.Data.ActionTimetoken))
				fmt.Println(fmt.Sprintf("%s messageActionsEvent.Data.MessageTimetoken: %s", outputPrefix, messageActionsEvent.Data.MessageTimetoken))
			case file := <-listener.File:
				fmt.Print(fmt.Sprintf("%s Subscribe Response:", outputPrefix))
				fmt.Println(" --- MessageActionsEvent: ")
				fmt.Println(fmt.Sprintf("file.File.PNMessage.Text: %s", file.File.PNMessage.Text))
				fmt.Println(fmt.Sprintf("file.File.PNFile.Name: %s", file.File.PNFile.Name))
				fmt.Println(fmt.Sprintf("file.File.PNFile.ID: %s", file.File.PNFile.ID))
				fmt.Println(fmt.Sprintf("file.File.PNFile.URL: %s", file.File.PNFile.URL))
				fmt.Println(fmt.Sprintf("file.Channel: %s", file.Channel))
				fmt.Println(fmt.Sprintf("file.Timetoken: %d", file.Timetoken))
				fmt.Println(fmt.Sprintf("file.SubscribedChannel: %s", file.SubscribedChannel))
				fmt.Println(fmt.Sprintf("file.Publisher: %s", file.Publisher))
				out, _ := os.Create("out.txt")
				resDLFile, statusDLFile, errDLFile := pn.DownloadFile().Channel("demo-channel").CipherKey("enigma").ID(file.File.PNFile.ID).Name(file.File.PNFile.Name).Execute()
				fmt.Println(statusDLFile, errDLFile)
				if resDLFile != nil {
					_, err := io.Copy(out, resDLFile.File)

					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}()

	pn.AddListener(listener)
	showHelp()

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
	fmt.Println("============")
	fmt.Println(" ||  HELP  ||")
	fmt.Println("============")
	fmt.Println("")
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
	showSetStateHelp()
	showGetStateHelp()
	showAddToCgHelp()
	showRemFromCgHelp()
	showListAllChOfCgHelp()
	showDelCgHelp()
	showGrantHelp()
	showGrantTokenHelp()
	showSubscribeWithStateHelp()
	showPresenceTimeoutHelp()
	showPresenceHelp()
	showMessageCountsHelp()
	showSignalHelp()
	showSetUUIDMetadataHelp()
	showGetAllUUIDMetadataHelp()
	showEditMembershipsHelp()
	showUpdateMembersHelp()
	showGetSpaceMembershipsHelp()
	showGetMembersHelp()
	showGetAllChannelMetadataHelp()
	showUpdateChannelMetadataHelp()
	showDeleteSpaceHelp()
	showCreateSpaceHelp()
	showGetSpaceHelp()
	showDeleteUserHelp()
	showUpdateUUIDMetadataHelp()
	showGetUserHelp()
	showAddActionHelp()
	showGetActionsHelp()
	showDeleteActionHelp()
	fmt.Println("")
	fmt.Println("================")
	fmt.Println(" ||  COMMANDS  ||")
	fmt.Println("================")
	fmt.Println("")
	fmt.Println(" UNSUBSCRIBE ALL \n\tq ")
	fmt.Println(" QUIT \n\tctrl+c ")
}

func showAddActionHelp() {
	fmt.Println(" AddAction EXAMPLE: ")
	fmt.Println("	addaction channel timetoken actiontype actionval")
	fmt.Println("	addaction my-channel 15210190573608384 reaction smiley_face")

}

func showGetActionsHelp() {
	fmt.Println(" GetActions EXAMPLE: ")
	fmt.Println("	getactions channel start end limit")
	fmt.Println("	getactions my-channel 15692395344923130 15210190573608384 10")

}

func showDeleteActionHelp() {
	fmt.Println(" DeleteAction EXAMPLE: ")
	fmt.Println("	remaction channel messagetTimetoken actionTimetoken")
	fmt.Println("	remaction my-channel 15210190573608384 15692395344923130 ")

}

func showEditMembershipsHelp() {
	fmt.Println(" EditMemberships EXAMPLE: ")
	fmt.Println("	managememberships channelMetadataid id a/u/r limit count")
	fmt.Println("	managememberships id0 id1 a 100 true")

}
func showUpdateMembersHelp() {
	fmt.Println(" UpdateMembers EXAMPLE: ")
	fmt.Println("	managemem memebers id a/u/r limit count")
	fmt.Println("	managemem id0 id0 a 100 true")
}

func showGetSpaceMembershipsHelp() {
	fmt.Println(" GetChannelMetadata EXAMPLE: ")
	fmt.Println("	getchannelmetadata channelMetadataid limit count start")
	fmt.Println("	getchannelmetadata id0 100 true Mymx")

}
func showGetMembersHelp() {
	fmt.Println(" GetMembers EXAMPLE: ")
	fmt.Println("	getmem userid limit count start")
	fmt.Println("	getmem id0 100 true Mymx")

}
func showGetAllChannelMetadataHelp() {
	fmt.Println(" GetAllChannelMetadata EXAMPLE: ")
	fmt.Println("	getallchannelmetadata limit count start")
	fmt.Println("	getallchannelmetadata 100 true MjWn")

}
func showUpdateChannelMetadataHelp() {
	fmt.Println(" UpdateChannelMetadata EXAMPLE: ")
	fmt.Println("	updatechannelmetadata id name desc")
	fmt.Println("	updatechannelmetadata id0 name desc")

}
func showDeleteSpaceHelp() {
	fmt.Println(" DeleteChannelMetadata EXAMPLE: ")
	fmt.Println("	delchannelmetadata id")
	fmt.Println("	delchannelmetadata id0")

}
func showCreateSpaceHelp() {
	fmt.Println(" SetChannelMetadata EXAMPLE: ")
	fmt.Println("	setchannelmetadata id name desc")
	fmt.Println("	setchannelmetadata id0 name desc")

}
func showGetSpaceHelp() {
	fmt.Println(" GetChannelMetadata EXAMPLE: ")
	fmt.Println("	getchannelmetadata id")
	fmt.Println("	getchannelmetadata id0")

}
func showDeleteUserHelp() {
	fmt.Println(" DeleteUUIDMetadata EXAMPLE: ")
	fmt.Println("	deleteuuidmetadata id")
	fmt.Println("	deleteuuidmetadata id0")
}

func showUpdateUUIDMetadataHelp() {
	fmt.Println(" UpdateUUIDMetadata EXAMPLE: ")
	fmt.Println("	updateuuidmetadata id name extid url email")
	fmt.Println("	updateuuidmetadata id0 name extid purl email")
}

func showGetUserHelp() {
	fmt.Println(" GetMetadata EXAMPLE: ")
	fmt.Println("	getuuidmetadata id")
	fmt.Println("	getuuidmetadata id0")
}

func showMessageCountsHelp() {
	fmt.Println(" MessageCounts EXAMPLE: ")
	fmt.Println("	messagecounts Channel(s) timetoken1,timetoken2")
	fmt.Println("	messagecounts my-channel,my-channel1 15210190573608384,15211140747622125")
}

func showGetAllUUIDMetadataHelp() {
	fmt.Println(" GetAllUUIDMetadata EXAMPLE: ")
	fmt.Println("	getalluuidmetadata limit count start")
	fmt.Println("	getalluuidmetadata 100 true MjWn")
}

func showSetUUIDMetadataHelp() {
	fmt.Println(" SetUUIDMetadata EXAMPLE: ")
	fmt.Println("	setuuidmetadata id name extid url email")
	fmt.Println("	setuuidmetadata id0 name extid purl email")
}

func showSignalHelp() {
	fmt.Println(" Signal EXAMPLE: ")
	fmt.Println("	signal Channel Message")
	fmt.Println("	signal my-channel \"my-signal\"")

}

func showGetStateHelp() {
	fmt.Println(" GET STATE EXAMPLE: ")
	fmt.Println("	getstate Channel ")
	fmt.Println("	getstate my-channel ")
}

func showSetStateHelp() {
	fmt.Println(" SET STATE EXAMPLE: ")
	fmt.Println("	setstate Channel state ")
	fmt.Println("	setstate my-channel {\"k\":\"v\"} ")
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

func showSubscribeWithStateHelp() {
	fmt.Println(" SUBSCRIBE With State EXAMPLE: ")
	fmt.Println("	subs withPresence channels channelGroups state")
	fmt.Println("	subs true my-channel,my-another-channel my-channelgroup,my-another-channel-group {\"k2\":\"v2\"}")
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

func showAddToCgHelp() {
	fmt.Println(" Add Channels to Channel Group EXAMPLE: ")
	fmt.Println("	addcg Channel ChannelGroup")
	fmt.Println("	addcg my-channel1,my-channel2 cg")
}

func showRemFromCgHelp() {
	fmt.Println(" Remove Channels from Channel Group EXAMPLE: ")
	fmt.Println("	remcg Channel ChannelGroup ")
	fmt.Println("	remcg my-channel1 cg")
}

func showListAllChOfCgHelp() {
	fmt.Println(" List Channels of Channel Group EXAMPLE: ")
	fmt.Println("	listcg ChannelGroup")
	fmt.Println("	listcg cg ")
}

func showDelCgHelp() {
	fmt.Println(" Delete Channel Group EXAMPLE: ")
	fmt.Println("	delcg ChannelGroup ")
	fmt.Println("	delcg cg ")
}

func showGrantTokenHelp() {
	fmt.Println(" GRANTTOKEN EXAMPLE: ")
	fmt.Println("	granttoken Channels ChannelGroups Users Spaces ttl ")
	fmt.Println("	granttoken ch1,ch2 cg1,cg2 u1,u2 s1,s2 ^ch-[0-9a-f]*$ ^:cg-[0-9a-f]*$ ^u-[0-9a-f]*$ ^s-[0-9a-f]*$ 10")
}

func showGrantHelp() {
	fmt.Println(" GRANT2 EXAMPLE: ")
	fmt.Println("	grant Channels ChannelGroups manage read write ttl ")
	fmt.Println("	grant my-channel cg false false false 10")
}

func showPresenceTimeoutHelp() {
	fmt.Println(" Presence Timeout: ")
	fmt.Println("	setpto presenceTimeout presenceHeartbeatInterval ")
	fmt.Println("	setpto 120 59")
}

func showPresenceHelp() {
	fmt.Println(" Presence: ")
	fmt.Println("	presence Connected Channels ChannelGroups")
	fmt.Println("	presence true my-channel,my-another-channel my-channelgroup,my-another-channel-group")
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
	case "subs":
		subscribeRequest(command[1:])
	case "setstate":
		setStateRequest(command[1:])
	case "getstate":
		getStateRequest(command[1:])
	case "addcg":
		addToChannelGroup(command[1:])
	case "remcg":
		removeFromChannelGroup(command[1:])
	case "listcg":
		listChannelsOfChannelGroup(command[1:])
	case "delcg":
		delChannelGroup(command[1:])
	case "grant":
		grant(command[1:])
	case "granttoken":
		granttoken(command[1:])
	case "help":
		showHelp()
	case "pt":
		publishTest()
	case "setpto":
		setPresenceTimeout(command[1:])
	case "presence":
		runPresenceRequest(command[1:])
	case "messagecounts":
		messageCounts(command[1:])
	case "signal":
		signal(command[1:])
	case "setuuidmetadata":
		setUUIDMetadata(command[1:])
	case "getalluuidmetadata":
		getAllUUIDMetadata(command[1:])
	case "getuuidmetadata":
		getUUIDMetadata(command[1:])
	case "updateuuidmetadata":
		updateUUIDMetadata(command[1:])
	case "deleteuuidmetadata":
		deleteUser(command[1:])
	case "getallchannelmetadata":
		getSpaces(command[1:])
	case "setchannelmetadata":
		createChannelMetadata(command[1:])
	case "delchannelmetadata":
		deleteChannelMetadata(command[1:])
	case "updatechannelmetadata":
		updateChannelMetadata(command[1:])
	case "getchannelmetadata":
		getChannelMetadata(command[1:])
	case "getmemberships":
		getSpaceMemberships(command[1:])
	case "getmem":
		getMembers(command[1:])
	case "managememberships":
		manageMemberships(command[1:])
	case "managemem":
		manageMembers(command[1:])
	case "settoken":
		setToken(command[1:])
	case "settokens":
		setTokens(command[1:])
	case "gettoken":
		getToken(command[1:])
	case "gettokens":
		getTokens(command[1:])
	case "gettokenres":
		getTokenRes(command[1:])
	case "addaction":
		addMessageAction(command[1:])
	case "addactions":
		addMessageActions(command[1:])
	case "getactions":
		getMessageActions(command[1:])
	case "getactionsrec":
		getMessageActionsRec(command[1:])
	case "getactionsrec2":
		getMessageActionsRec(command[1:])
	case "uploadfile":
		uploadFile(command[1:])
	case "delfile":
		delFile(command[1:])
	case "listfiles":
		listFiles(command[1:])
	case "getfileurl":
		getFileURL(command[1:])
	case "downloadfile":
		downloadFile(command[1:])
	case "q":
		pn.UnsubscribeAll()
	case "d":
		pn.Destroy()
		fmt.Println("after Destroy")
	default:
		showHelp()
	}
}

func uploadFile(args []string) {
	channel := args[0]
	message := args[1]
	name := args[2]
	filepath := args[3]
	file, err := os.Open(filepath)

	defer file.Close()

	cipherKey := args[4]
	res, status, err := pn.SendFile().Channel(channel).Message(message).CipherKey(cipherKey).Name(name).File(file).Execute()
	fmt.Println("status", status)
	fmt.Println("err", err)
	fmt.Println("res", res)

}

func delFile(args []string) {
	ch := args[0]
	id := args[1]
	name := args[2]

	res, status, err := pn.DeleteFile().Channel(ch).ID(id).Name(name).Execute()
	fmt.Println("status", status)
	fmt.Println("err", err)
	fmt.Println("res", res)
}

func listFiles(args []string) {
	ch := args[0]
	res, status, err := pn.ListFiles().Channel(ch).Execute()
	fmt.Println("status", status)
	fmt.Println("err", err)
	fmt.Println("res", res)
}

func getFileURL(args []string) {
	ch := args[0]
	id := args[1]
	name := args[2]

	res, status, err := pn.GetFileURL().Channel(ch).ID(id).Name(name).Execute()
	fmt.Println("status", status)
	fmt.Println("err", err)
	fmt.Println("res", res)
}

func downloadFile(args []string) {
	ch := args[0]
	id := args[1]
	name := args[2]

	resDLFile, status, err := pn.DownloadFile().Channel(ch).ID(id).Name(name).Execute()
	fmt.Println("status", status)
	fmt.Println("err", err)
	fmt.Println("res", resDLFile)

	if resDLFile != nil {
		out, _ := os.Create("out.txt")
		defer out.Close()

		_, err := io.Copy(out, resDLFile.File)

		if err != nil {
			fmt.Println(err)
		}
	}
}

func getMessageActionsRec2(args []string) {
	channel := args[0]
	getMessageActionsRecursive(channel, "", false, 0)
}

func getMessageActionsRec(args []string) {
	channel := args[0]
	getMessageActionsRecursive(channel, "", true, 0)
}

func getMessageActionsRecursive(channel string, start string, more bool, counter int) {
	var res *pubnub.PNGetMessageActionsResponse
	if start == "" {
		res, _, _ = pn.GetMessageActions().Channel(channel).Execute()
	} else {
		res, _, _ = pn.GetMessageActions().Channel(channel).Start(start).Execute()
	}
	if (res != nil) && (len(res.Data) > 0) {
		printMessageActions(res, counter+1)
		if more {
			if res.More.Start != "" {
				getMessageActionsRecursive(channel, res.More.Start, more, len(res.Data))
			}
		} else {
			if len(res.Data) > 0 {
				getMessageActionsRecursive(channel, res.Data[0].ActionTimetoken, more, len(res.Data))
			}
		}
	}

}

func removeMessageActions(args []string) {
	if len(args) < 3 {
		showDeleteActionHelp()
		return
	}
	channel := args[0]
	tt := args[1]
	att := args[2]
	res, status, err := pn.RemoveMessageAction().Channel(channel).MessageTimetoken(tt).ActionTimetoken(att).Execute()
	fmt.Println("status", status)
	fmt.Println("err", err)
	fmt.Println("res", res)
}

func getMessageActions(args []string) {
	if len(args) < 1 {
		showGetActionsHelp()
		return
	}
	channel := args[0]
	if len(args) == 4 {
		var limit int

		n, err := strconv.ParseInt(args[3], 10, 64)
		if err == nil {
			limit = int(n)
		}

		res, status, err := pn.GetMessageActions().Channel(channel).Start(args[1]).End(args[2]).Limit(limit).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		printMessageActions(res, 0)

	} else if len(args) == 3 {
		res, status, err := pn.GetMessageActions().Channel(channel).Start(args[1]).End(args[2]).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		printMessageActions(res, 0)
	} else if len(args) == 2 {
		res, status, err := pn.GetMessageActions().Channel(channel).Start(args[1]).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		printMessageActions(res, 0)
	} else {
		res, status, err := pn.GetMessageActions().Channel(channel).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		printMessageActions(res, 0)
	}

}

func printMessageActions(res *pubnub.PNGetMessageActionsResponse, counter int) {
	if res != nil {
		for i, k := range res.Data {
			fmt.Println(fmt.Sprintf("No: %d, Val: %s", i+counter, k))
		}
		fmt.Println("More:", res.More)
	}

}

func addMessageActions(args []string) {
	if len(args) < 5 {
		showAddActionHelp()
		return
	}
	// addaction my-channel 15210190573608384 reaction smiley_face
	channel := args[0]
	tt := args[1]
	actionType := args[2]
	actionVal := args[3]

	var count int

	n, err := strconv.ParseInt(args[4], 10, 64)
	if err == nil {
		count = int(n)
	}

	for i := 0; i < count; i++ {
		ma := pubnub.MessageAction{
			ActionType:  actionType,
			ActionValue: actionVal + "_" + strconv.Itoa(i),
		}

		res, status, err := pn.AddMessageAction().Channel(channel).MessageTimetoken(tt).Action(ma).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		fmt.Println("res", res)
	}
}

func addMessageAction(args []string) {
	if len(args) < 4 {
		showAddActionHelp()
		return
	}
	// addaction my-channel 15210190573608384 reaction smiley_face
	channel := args[0]
	tt := args[1]
	actionType := args[2]
	actionVal := args[3]
	ma := pubnub.MessageAction{
		ActionType:  actionType,
		ActionValue: actionVal,
	}

	res, status, err := pn.AddMessageAction().Channel(channel).MessageTimetoken(tt).Action(ma).Execute()
	fmt.Println("status", status)
	fmt.Println("err", err)
	fmt.Println("res", res)
}

func setToken(args []string) {
	pn.SetToken(args[0])
}

func setTokens(args []string) {
	var tokens []string
	tokens = strings.Split(args[0], ",")

	pn.SetTokens(tokens)
}

func getToken(args []string) {
	res := pn.GetTokens()
	fmt.Println(res)
}

func getTokens(args []string) {
	res := pn.GetTokens()
	fmt.Println(res)
}

func getTokenRes(args []string) {
	res := pn.GetTokens()
	fmt.Println(res)

}

func manageMembers(args []string) {
	if len(args) < 5 {
		showEditMembershipsHelp()
		return
	}
	channelMetadataID := args[0]
	id0 := args[1]
	//id1 := args[2]
	//id2 := args[3]
	action := args[2]
	var limit int

	n, err := strconv.ParseInt(args[3], 10, 64)
	if err == nil {
		limit = int(n)
	}
	count, _ := strconv.ParseBool(args[4])
	var start string
	if len(args) > 5 {
		start = args[5]
	}

	incl := []pubnub.PNChannelMembersInclude{
		pubnub.PNChannelMembersIncludeCustom,
		pubnub.PNChannelMembersIncludeUUID,
		pubnub.PNChannelMembersIncludeUUIDCustom,
	}

	custom := make(map[string]interface{})
	custom["a1"] = "b1"
	custom["c1"] = "d1"
	uuid := pubnub.PNChannelMembersUUID{
		ID: id0,
	}

	in := pubnub.PNChannelMembersSet{
		UUID:   uuid,
		Custom: custom,
	}

	custom2 := make(map[string]interface{})
	custom2["a2"] = "b2"
	custom2["c2"] = "d2"

	up := pubnub.PNChannelMembersSet{
		UUID:   uuid,
		Custom: custom2,
	}

	inArr := []pubnub.PNChannelMembersSet{
		in,
		up,
	}
	re := pubnub.PNChannelMembersRemove{
		UUID: uuid,
	}

	reArr := []pubnub.PNChannelMembersRemove{
		re,
	}

	if action == "a" {
		reArr = []pubnub.PNChannelMembersRemove{}
		inArr = []pubnub.PNChannelMembersSet{}
	} else if action == "u" {
		reArr = []pubnub.PNChannelMembersRemove{}
		inArr = []pubnub.PNChannelMembersSet{}
	} else if action == "r" {
		inArr = []pubnub.PNChannelMembersSet{}
	}

	if start != "" {
		res, status, err := pn.ManageChannelMembers().Channel(channelMetadataID).Set(inArr).Remove(reArr).Include(incl).Limit(limit).Count(count).Start(start).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		fmt.Println("res", res)

	} else {
		res, status, err := pn.ManageChannelMembers().Channel(channelMetadataID).Set(inArr).Remove(reArr).Include(incl).Limit(limit).Count(count).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		fmt.Println("res", res)
	}
}

func manageMemberships(args []string) {
	if len(args) < 5 {
		showUpdateMembersHelp()
		return
	}
	userID := args[0]
	id0 := args[1]
	//id1 := args[2]
	//id2 := args[3]
	action := args[2]
	var limit int

	n, err := strconv.ParseInt(args[3], 10, 64)
	if err == nil {
		limit = int(n)
	}
	count, _ := strconv.ParseBool(args[4])
	var start string
	if len(args) > 5 {
		start = args[5]
	}

	incl := []pubnub.PNMembershipsInclude{
		pubnub.PNMembershipsIncludeCustom,
		pubnub.PNMembershipsIncludeChannel,
		pubnub.PNMembershipsIncludeChannelCustom,
	}

	custom3 := make(map[string]interface{})
	custom3["a3"] = "b3"
	custom3["c3"] = "d3"

	channel := pubnub.PNMembershipsChannel{
		ID: id0,
	}

	in := pubnub.PNMembershipsSet{
		Channel: channel,
		Custom:  custom3,
	}

	custom4 := make(map[string]interface{})
	custom4["a4"] = "b4"
	custom4["c4"] = "d4"

	up := pubnub.PNMembershipsSet{
		Channel: channel,
		Custom:  custom4,
	}

	inArr := []pubnub.PNMembershipsSet{
		in,
		up,
	}

	re := pubnub.PNMembershipsRemove{
		Channel: channel,
	}

	reArr := []pubnub.PNMembershipsRemove{
		re,
	}

	if action == "a" {
		reArr = []pubnub.PNMembershipsRemove{}
		inArr = []pubnub.PNMembershipsSet{}
	} else if action == "u" {
		reArr = []pubnub.PNMembershipsRemove{}
		inArr = []pubnub.PNMembershipsSet{}
	} else if action == "r" {
		inArr = []pubnub.PNMembershipsSet{}
	}

	if start != "" {
		res, status, err := pn.ManageMemberships().Set(inArr).Remove(reArr).Include(incl).Limit(limit).Count(count).Start(start).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		fmt.Println("res", res)

	} else {

		res, status, err := pn.ManageMemberships().UUID(userID).Set(inArr).Remove(reArr).Include(incl).Limit(limit).Count(count).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		fmt.Println("res", res)
	}
}

func getSpaceMemberships(args []string) {
	if len(args) < 3 {
		showGetSpaceMembershipsHelp()
		return
	}
	id := args[0]

	var limit int

	n, err := strconv.ParseInt(args[1], 10, 64)
	if err == nil {
		limit = int(n)
	}
	count, _ := strconv.ParseBool(args[2])
	var start string
	if len(args) > 3 {
		start = args[3]
	}

	incl := []pubnub.PNMembershipsInclude{
		pubnub.PNMembershipsIncludeCustom,
		pubnub.PNMembershipsIncludeChannel,
		pubnub.PNMembershipsIncludeChannelCustom,
	}
	if start != "" {
		res, status, err := pn.GetMemberships().UUID(id).Include(incl).Limit(limit).Count(count).Start(start).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		fmt.Println("res", res)
	} else {
		res, status, err := pn.GetMemberships().UUID(id).Include(incl).Limit(limit).Count(count).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		fmt.Println("res", res)
	}
}

func getMembers(args []string) {
	if len(args) < 3 {
		showGetMembersHelp()
		return
	}
	id := args[0]

	var limit int

	n, err := strconv.ParseInt(args[1], 10, 64)
	if err == nil {
		limit = int(n)
	}
	count, _ := strconv.ParseBool(args[2])
	var start string
	if len(args) > 3 {
		start = args[3]
	}

	incl := []pubnub.PNChannelMembersInclude{
		pubnub.PNChannelMembersIncludeCustom,
		pubnub.PNChannelMembersIncludeUUIDCustom,
		pubnub.PNChannelMembersIncludeUUID,
	}
	sort := []string{"updated:desc"}
	if start != "" {
		res, status, err := pn.GetChannelMembers().Channel(id).Include(incl).Sort(sort).Limit(limit).Count(count).Start(start).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		fmt.Println("res", res)
	} else {
		res, status, err := pn.GetChannelMembers().Channel(id).Sort(sort).Include(incl).Limit(limit).Count(count).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		fmt.Println("res", res)
	}
}

func getSpaces(args []string) {
	if len(args) < 2 {
		showGetAllChannelMetadataHelp()
		return
	}
	var limit int

	n, err := strconv.ParseInt(args[0], 10, 64)
	if err == nil {
		limit = int(n)
	}
	count, _ := strconv.ParseBool(args[1])
	var start string
	if len(args) > 2 {
		start = args[2]
	}

	incl := []pubnub.PNChannelMetadataInclude{
		pubnub.PNChannelMetadataIncludeCustom,
	}

	if start != "" {
		res, status, err := pn.GetAllChannelMetadata().Include(incl).Limit(limit).Count(count).Start(start).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		fmt.Println("res", res)
	} else {
		res, status, err := pn.GetAllChannelMetadata().Include(incl).Limit(limit).Count(count).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		fmt.Println("res", res)
	}
}

func updateChannelMetadata(args []string) {
	if len(args) < 2 {
		showUpdateChannelMetadataHelp()
		return
	}
	id := args[0]
	name := args[1]
	desc := args[2]

	custom := make(map[string]interface{})
	custom["a"] = "b"
	custom["c"] = "d"

	incl := []pubnub.PNChannelMetadataInclude{
		pubnub.PNChannelMetadataIncludeCustom,
	}

	res, status, err := pn.SetChannelMetadata().Channel(id).Name(name).Description(desc).Include(incl).Custom(custom).Execute()
	fmt.Println("status", status)
	fmt.Println("err", err)
	fmt.Println("res", res)
}

func deleteChannelMetadata(args []string) {
	if len(args) < 1 {
		showDeleteSpaceHelp()
		return
	}
	id := args[0]

	res, status, err := pn.RemoveChannelMetadata().Channel(id).Execute()
	fmt.Println("status", status)
	fmt.Println("err", err)
	fmt.Println("res", res)
}

func createChannelMetadata(args []string) {
	if len(args) < 3 {
		showCreateSpaceHelp()
		return
	}
	id := args[0]
	name := args[1]
	desc := args[2]

	custom := make(map[string]interface{})
	custom["a"] = "b"

	incl := []pubnub.PNChannelMetadataInclude{
		pubnub.PNChannelMetadataIncludeCustom,
		pubnub.PNChannelMetadataIncludeCustom,
	}

	res, status, err := pn.SetChannelMetadata().Channel(id).Name(name).Description(desc).Include(incl).Custom(custom).Execute()
	fmt.Println("status", status)
	fmt.Println("err", err)
	fmt.Println("res", res)
}

func getChannelMetadata(args []string) {
	if len(args) < 1 {
		showGetSpaceHelp()
		return
	}
	id := args[0]

	incl := []pubnub.PNChannelMetadataInclude{
		pubnub.PNChannelMetadataIncludeCustom,
	}

	res, status, err := pn.GetChannelMetadata().Channel(id).Include(incl).Execute()
	fmt.Println("status", status)
	fmt.Println("err", err)
	fmt.Println("res", res)
}

func deleteUser(args []string) {
	if len(args) < 1 {
		showDeleteUserHelp()
		return
	}
	id := args[0]

	custom := make(map[string]interface{})
	custom["a"] = "b"
	custom["c"] = "d"

	res, status, err := pn.RemoveUUIDMetadata().UUID(id).Execute()
	fmt.Println("status", status)
	fmt.Println("err", err)
	fmt.Println("res", res)
}

func updateUUIDMetadata(args []string) {
	if len(args) < 5 {
		showUpdateUUIDMetadataHelp()
		return
	}
	id := args[0]
	name := args[1]
	extid := args[2]
	purl := args[3]
	email := args[4]

	custom := make(map[string]interface{})
	custom["a"] = "b"
	custom["c"] = "d"
	incl := []pubnub.PNUUIDMetadataInclude{
		pubnub.PNUUIDMetadataIncludeCustom,
	}

	res, status, err := pn.SetUUIDMetadata().Include(incl).UUID(id).Name(name).ExternalID(extid).ProfileURL(purl).Email(email).Custom(custom).Execute()
	fmt.Println("status", status)
	fmt.Println("err", err)
	fmt.Println("res", res)
}

func getUUIDMetadata(args []string) {
	if len(args) < 1 {
		showGetUserHelp()
		return
	}
	id := args[0]

	incl := []pubnub.PNUUIDMetadataInclude{
		pubnub.PNUUIDMetadataIncludeCustom,
	}

	res, status, err := pn.GetUUIDMetadata().Include(incl).UUID(id).Execute()
	fmt.Println("status", status)
	fmt.Println("err", err)
	fmt.Println("res", res)
}

func getAllUUIDMetadata(args []string) {
	if len(args) < 2 {
		showGetAllUUIDMetadataHelp()
		return
	}
	var limit int

	n, err := strconv.ParseInt(args[0], 10, 64)
	if err == nil {
		limit = int(n)
	}
	count, _ := strconv.ParseBool(args[1])
	var start string
	if len(args) > 2 {
		start = args[2]
	}

	incl := []pubnub.PNUUIDMetadataInclude{
		pubnub.PNUUIDMetadataIncludeCustom,
	}

	if start != "" {
		res, status, err := pn.GetAllUUIDMetadata().Include(incl).Start(start).Limit(limit).Count(count).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		fmt.Println("res", res)

	} else {
		res, status, err := pn.GetAllUUIDMetadata().Include(incl).Limit(limit).Filter("name == 'name 891'").Count(count).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		fmt.Println("res", res)
	}

}

func setUUIDMetadata(args []string) {
	if len(args) < 5 {
		showSetUUIDMetadataHelp()
		return
	}
	id := args[0]
	name := args[1]
	extid := args[2]
	purl := args[3]
	email := args[4]

	custom := make(map[string]interface{})
	custom["a"] = "b"
	custom["c"] = "d"

	incl := []pubnub.PNUUIDMetadataInclude{
		pubnub.PNUUIDMetadataIncludeCustom,
	}

	res, status, err := pn.SetUUIDMetadata().Include(incl).UUID(id).Name(name).ExternalID(extid).ProfileURL(purl).Email(email).Custom(custom).Execute()
	fmt.Println("status", status)
	fmt.Println("err", err)
	fmt.Println("res", res)
}

func signal(args []string) {
	if len(args) < 2 {
		showSignalHelp()
		return
	}

	var channel string
	channel = args[0]

	message := args[1]

	res, status, err := pn.Signal().Channel(channel).Message(message).Execute()
	fmt.Println("status", status)
	fmt.Println(err)
	fmt.Println("res", res)

}

func messageCounts(args []string) {
	if len(args) < 2 {
		showMessageCountsHelp()
		return
	}

	var channels []string
	channels = strings.Split(args[0], ",")

	var channelsTimetoken []int64
	if len(args) > 1 {
		strSlice := strings.Split(args[1], ",")
		channelsTimetoken = make([]int64, len(strSlice))
		for i := range strSlice {
			n, err := strconv.ParseInt(strSlice[i], 10, 64)
			if err == nil {
				channelsTimetoken[i] = n
			} else {
				fmt.Println(err)
			}
		}
	}

	res, status, err := pn.MessageCounts().Channels(channels).ChannelsTimetoken(channelsTimetoken).Execute()
	fmt.Println(status)
	fmt.Println(err)
	if err == nil {
		for ch, v := range res.Channels {
			fmt.Printf("%s %d", ch, v)
			fmt.Println("")
		}
	}

}

func runPresenceRequest(args []string) {
	if len(args) < 2 {
		showPresenceHelp()
	}
	var connected bool
	connected, _ = strconv.ParseBool(args[0])

	var channels []string
	channels = strings.Split(args[1], ",")

	var groups []string
	if len(args) > 2 {
		groups = strings.Split(args[2], ",")
	}
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}
	state := map[string]interface{}{
		"state": "stateval",
	}
	pn.Presence().Connected(connected).Channels(channels).QueryParam(queryParam).State(state).ChannelGroups(groups).Execute()
}

func setPresenceTimeout(args []string) {
	if len(args) < 0 {
		showPresenceTimeoutHelp()
	}

	var timeout int
	if len(args) > 0 {
		i, err := strconv.ParseInt(args[0], 10, 64)
		if err == nil {
			timeout = int(i)
		}
	}

	var interval int
	if len(args) > 1 {
		i, err := strconv.ParseInt(args[1], 10, 64)
		if err == nil {
			interval = int(i)
		}
	}

	if interval <= 0 {
		pn.Config.SetPresenceTimeout(timeout)
	} else {
		pn.Config.SetPresenceTimeoutWithCustomInterval(timeout, interval)
	}
}

func granttoken(args []string) {
	if len(args) < 9 {
		fmt.Println(len(args))
		showGrantHelp()
		return
	}

	var channels []string
	if len(args) > 0 {
		channels = strings.Split(args[0], ",")
	}
	var groups []string
	if len(args) > 1 {
		groups = strings.Split(args[1], ",")
	}
	var users []string
	if len(args) > 2 {
		users = strings.Split(args[2], ",")
	}
	var channelMetadatas []string
	if len(args) > 3 {
		channelMetadatas = strings.Split(args[3], ",")
	}
	var channelsPat []string
	if len(args) > 0 {
		channelsPat = strings.Split(args[4], ",")
	}
	var groupsPat []string
	if len(args) > 1 {
		groupsPat = strings.Split(args[5], ",")
	}
	var usersPat []string
	if len(args) > 2 {
		usersPat = strings.Split(args[6], ",")
	}
	var channelMetadatasPat []string
	if len(args) > 3 {
		channelMetadatasPat = strings.Split(args[7], ",")
	}
	var ttl int
	if len(args) > 4 {
		i, err := strconv.ParseInt(args[8], 10, 64)
		if err != nil {
			ttl = 1440
		} else {
			ttl = int(i)
		}
	}

	// ch1 := randomnized("ch1")
	// cg1 := "cg"
	// cg2 := "cg1"
	// u1 := "u"
	// s1 := "s"

	ch := make(map[string]pubnub.ChannelPermissions, len(channels))
	for _, k := range channels {
		ch[k] = pubnub.ChannelPermissions{
			Read:   true,
			Write:  true,
			Delete: false,
		}
	}

	s := make(map[string]pubnub.UserSpacePermissions, len(channelMetadatas))
	for _, k := range channelMetadatas {
		s[k] = pubnub.UserSpacePermissions{
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: true,
			Create: true,
		}
	}

	u := make(map[string]pubnub.UserSpacePermissions, len(users))
	for _, k := range users {
		u[k] = pubnub.UserSpacePermissions{
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: false,
			Create: false,
		}
	}

	cg := make(map[string]pubnub.GroupPermissions, len(groups))
	for _, k := range groups {
		cg[k] = pubnub.GroupPermissions{
			Read:   true,
			Manage: false,
		}
	}

	chPat := make(map[string]pubnub.ChannelPermissions, len(channelsPat))
	for _, k := range channelsPat {
		chPat[k] = pubnub.ChannelPermissions{
			Read:   true,
			Write:  true,
			Delete: true,
		}
	}

	sPat := make(map[string]pubnub.UserSpacePermissions, len(channelMetadatasPat))
	for _, k := range channelMetadatasPat {
		sPat[k] = pubnub.UserSpacePermissions{
			Read:   true,
			Write:  true,
			Manage: false,
			Delete: true,
			Create: true,
		}
	}

	uPat := make(map[string]pubnub.UserSpacePermissions, len(usersPat))
	for _, k := range usersPat {
		uPat[k] = pubnub.UserSpacePermissions{
			Read:   true,
			Write:  true,
			Manage: true,
			Delete: true,
			Create: false,
		}
	}

	cgPat := make(map[string]pubnub.GroupPermissions, len(groupsPat))
	for _, k := range groupsPat {
		cgPat[k] = pubnub.GroupPermissions{
			Read:   true,
			Manage: true,
		}
	}

	// u := map[string]pubnub.ResourcePermissions{
	// 	u1: pubnub.ResourcePermissions{
	// 		Read:   true,
	// 		Write:  true,
	// 		Manage: true,
	// 		Delete: true,
	// 		Create: false,
	// 	},
	// }

	// s := map[string]pubnub.ResourcePermissions{
	// 	s1: pubnub.ResourcePermissions{
	// 		Read:   true,
	// 		Write:  true,
	// 		Manage: true,
	// 		Delete: true,
	// 		Create: true,
	// 	},
	// }

	// cg := map[string]pubnub.ResourcePermissions{
	// 	cg1: pubnub.ResourcePermissions{
	// 		Read:   true,
	// 		Write:  true,
	// 		Manage: true,
	// 		Delete: false,
	// 		Create: true,
	// 	},
	// 	cg2: pubnub.ResourcePermissions{
	// 		Read:   true,
	// 		Write:  true,
	// 		Manage: false,
	// 		Delete: false,
	// 		Create: true,
	// 	},
	// }

	res, _, err := pn.GrantToken().TTL(ttl).
		//Channels(ch).
		//ChannelGroups(cg).
		Users(u).
		Spaces(s).
		//ChannelsPattern(chPat).
		//ChannelGroupsPattern(cgPat).
		UsersPattern(uPat).
		SpacesPattern(sPat).
		Execute()

	fmt.Println(res)
	fmt.Println(err)
}

func grant(args []string) {
	if len(args) < 6 {
		fmt.Println(len(args))
		showGrantHelp()
		return
	}

	var channels []string
	if len(args) > 0 {
		channels = strings.Split(args[0], ",")
	}
	var groups []string
	if len(args) > 1 {
		groups = strings.Split(args[1], ",")
	}
	var manage bool
	if len(args) > 2 {
		manage, _ = strconv.ParseBool(args[2])
	}
	var read bool
	if len(args) > 3 {
		read, _ = strconv.ParseBool(args[3])
	}
	var write bool
	if len(args) > 4 {
		write, _ = strconv.ParseBool(args[4])
	}
	var ttl int
	if len(args) > 5 {
		i, err := strconv.ParseInt(args[5], 10, 64)
		if err != nil {
			ttl = 1440
		} else {
			ttl = int(i)
		}
	}

	res, _, err := pn.Grant().
		ChannelGroups(groups).
		Channels(channels).
		Manage(manage).
		Read(read).
		TTL(ttl).
		Write(write).
		Execute()

	fmt.Println(res)
	fmt.Println(err)
}

func addToChannelGroup(args []string) {
	if len(args) < 2 {
		fmt.Println(len(args))
		showAddToCgHelp()
		return
	}
	var channels []string
	if len(args) > 0 {
		channels = strings.Split(args[0], ",")
	}

	var cg string
	if len(args) > 1 {
		cg = args[1]
	}

	_, _, err := pn.AddChannelToChannelGroup().
		Channels(channels).
		ChannelGroup(cg).
		Execute()

	fmt.Println(err)
}

func removeFromChannelGroup(args []string) {
	if len(args) < 2 {
		fmt.Println(len(args))
		showRemFromCgHelp()
		return
	}
	var channels []string
	if len(args) > 0 {
		channels = strings.Split(args[0], ",")
	}

	var cg string
	if len(args) > 1 {
		cg = args[1]
	}

	_, _, err := pn.RemoveChannelFromChannelGroup().
		Channels(channels).
		ChannelGroup(cg).
		Execute()

	fmt.Println(err)
}

func listChannelsOfChannelGroup(args []string) {
	if len(args) < 1 {
		fmt.Println(len(args))
		showListAllChOfCgHelp()
		return
	}

	var cg string
	if len(args) > 0 {
		cg = args[0]
	}

	res, _, err := pn.ListChannelsInChannelGroup().
		ChannelGroup(cg).
		Execute()
	fmt.Println("ChannelGroup", res.ChannelGroup)
	for _, ch := range res.Channels {
		fmt.Println(ch)
	}
	fmt.Println(err)
}

func delChannelGroup(args []string) {
	if len(args) < 1 {
		fmt.Println(len(args))
		showDelCgHelp()
		return
	}

	var cg string
	if len(args) > 0 {
		cg = args[0]
	}

	_, _, err := pn.DeleteChannelGroup().
		ChannelGroup(cg).
		Execute()

	fmt.Println(err)

}

func setStateRequest(args []string) {
	if len(args) < 2 {
		fmt.Println(len(args))
		showSetStateHelp()
		return
	}

	var channel string
	if len(args) > 0 {
		channel = args[0]
	}

	var state map[string]interface{}
	if len(args) > 1 {
		var v interface{}
		err := json.Unmarshal([]byte(args[1]), &v)
		if err == nil {
			if st, ok := v.(map[string]interface{}); ok {
				state = st
			} else {

				fmt.Println("!ok", reflect.TypeOf(v))
				showSetStateHelp()
				return
			}
		} else {
			fmt.Println("err", err)
			showSetStateHelp()
			return
		}
	}

	res, status, err := pn.SetState().Channels([]string{channel}).State(state).UUID("nuuid").Execute()

	fmt.Println("status===>", status)
	if err != nil {
		fmt.Println("err=>>>", err)
	} else {
		fmt.Println(res.State)
		fmt.Println(res.Message)
	}
}

func getStateRequest(args []string) {
	if len(args) < 1 {
		fmt.Println(len(args))
		showGetStateHelp()
		return
	}

	var channel string
	if len(args) > 0 {
		channel = args[0]
	}

	res, status, err := pn.GetState().Channels([]string{channel}).UUID("").Execute()

	fmt.Println("status===>", status)
	if err != nil {
		fmt.Println("err=>>>", err)
	} else {
		for j, k := range res.State {
			fmt.Println("channel:", j, ", state:", k)
		}
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
			start = 0
		} else {
			start = i
		}
	}

	var end int64
	if len(args) > 2 {
		i, err := strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			end = 0
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
		res, status, err := pn.WhereNow().UUID(uuidToUse).Execute()
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
			count = 100
		} else {
			count = int(i)
		}
	}

	var start int64
	if len(args) > 3 {
		i, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			start = 0
		} else {
			start = i
		}
	}

	var end int64
	if len(args) > 4 {
		i, err := strconv.ParseInt(args[4], 10, 64)
		if err != nil {
			end = 0
		} else {
			end = i
		}
	}

	var withMessageActions = false
	if len(args) > 5 {
		withMessageActions, _ = strconv.ParseBool(args[5])
	}

	var withMeta bool
	if len(args) > 6 {
		withMeta, _ = strconv.ParseBool(args[6])
	}

	if (end != 0) && (start != 0) {
		res, status, err := pn.Fetch().
			Channels(channels).
			Count(count).
			Start(start).
			End(end).
			Reverse(reverse).
			IncludeMessageActions(withMessageActions).
			IncludeMeta(withMeta).
			Execute()
		parseFetch(res, status, err)
	} else if start != 0 {
		res, status, err := pn.Fetch().
			Channels(channels).
			Count(count).
			Start(start).
			Reverse(reverse).
			IncludeMessageActions(withMessageActions).
			IncludeMeta(withMeta).
			Execute()
		parseFetch(res, status, err)
	} else if end != 0 {
		res, status, err := pn.Fetch().
			Channels(channels).
			Count(count).
			End(end).
			Reverse(reverse).
			IncludeMessageActions(withMessageActions).
			IncludeMeta(withMeta).
			Execute()
		parseFetch(res, status, err)
	} else {
		res, status, err := pn.Fetch().
			Channels(channels).
			Count(count).
			Reverse(reverse).
			IncludeMessageActions(withMessageActions).
			IncludeMeta(withMeta).
			Execute()
		parseFetch(res, status, err)
	}
}

func parseFetch(res *pubnub.FetchResponse, status pubnub.StatusResponse, err error) {
	fmt.Println(fmt.Sprintf("%s ParseFetch:", outputPrefix))
	if status.Error == nil {
		for channel, messages := range res.Messages {
			fmt.Println("channel:", channel)
			for _, messageInt := range messages {
				message := pubnub.FetchResponseItem(messageInt)
				fmt.Println("message.Message:", message.Message)
				fmt.Println("message.Timetoken:", message.Timetoken)
				fmt.Println("message.Meta:", message.Meta)
				fmt.Println("message.Actions:", message.MessageActions)

				for action := range message.MessageActions {
					actionTypes := message.MessageActions[action].ActionsTypeValues
					fmt.Println("action1:", action)
					if len(actionTypes) > 0 {
						for actionVal, actionType := range actionTypes {
							fmt.Println("actionVal:", actionVal)
							r00 := actionType
							if r00 != nil {
								fmt.Println("UUID", r00[0].UUID)
								fmt.Println("ActionTimetoken", r00[0].ActionTimetoken)
							} else {
								fmt.Println("r0 nil")
							}
						}
					} else {
						fmt.Println("actionTypes nil")
					}
				}
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
			count = 100
		} else {
			count = int(i)
		}
	}

	var start int64
	if len(args) > 4 {
		i, err := strconv.ParseInt(args[4], 10, 64)
		if err != nil {
			start = 0
		} else {
			start = i
		}
	}

	var end int64
	if len(args) > 5 {
		i, err := strconv.ParseInt(args[5], 10, 64)
		if err != nil {
			end = 0
		} else {
			end = i
		}
	}

	var withMeta bool
	if len(args) > 6 {
		withMeta, _ = strconv.ParseBool(args[6])
	}

	if (end != 0) && (start != 0) {
		res, status, err := pn.History().
			Channel(channel).
			Count(count).
			Start(start).
			End(end).
			IncludeTimetoken(includeTimetoken).
			Reverse(reverse).
			IncludeMeta(withMeta).
			Execute()
		parseHistory(res, status, err)
	} else if start != 0 {
		res, status, err := pn.History().
			Channel(channel).
			Count(count).
			Start(start).
			IncludeTimetoken(includeTimetoken).
			Reverse(reverse).
			IncludeMeta(withMeta).
			Execute()
		parseHistory(res, status, err)
	} else if end != 0 {
		res, status, err := pn.History().
			Channel(channel).
			Count(count).
			End(end).
			IncludeTimetoken(includeTimetoken).
			Reverse(reverse).
			IncludeMeta(withMeta).
			Execute()
		parseHistory(res, status, err)
	} else {
		res, status, err := pn.History().
			Channel(channel).
			Count(count).
			IncludeTimetoken(includeTimetoken).
			Reverse(reverse).
			IncludeMeta(withMeta).
			Execute()
		parseHistory(res, status, err)
	}
}

func parseHistory(res *pubnub.HistoryResponse, status pubnub.StatusResponse, err error) {
	fmt.Println(fmt.Sprintf("%s ParseHistory:", outputPrefix))
	if res != nil {
		if res.Messages != nil {
			for _, v := range res.Messages {
				fmt.Println(fmt.Sprintf("%s Timetoken %d", outputPrefix, v.Timetoken))
				fmt.Println(fmt.Sprintf("%s Message %s", outputPrefix, v.Message))

				fmt.Println(fmt.Sprintf("%s Meta %s", outputPrefix, v.Meta))
			}
		} else {
			fmt.Println(fmt.Sprintf("res.Messages null"))
		}
		fmt.Println(fmt.Sprintf("%s EndTimetoken %d", outputPrefix, res.EndTimetoken))
		fmt.Println(fmt.Sprintf("%s StartTimetoken %d", outputPrefix, res.StartTimetoken))
		fmt.Println(fmt.Sprintf("%s", outputSuffix))
	} else {
		fmt.Println(fmt.Sprintf("%s StatusResponse %s %e", outputPrefix, status.Error, err))
	}
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
			fmt.Println(v.UUID)
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
		res, status, err := pn.HereNow().Channels(channels).ChannelGroups(channelGroups).IncludeState(includeState).IncludeUUIDs(includeUUIDs).Execute()
		hereNowResponse(res, status, err)
	} else if len(channels) != 0 {
		res, status, err := pn.HereNow().Channels(channels).IncludeState(includeState).IncludeUUIDs(includeUUIDs).Execute()
		hereNowResponse(res, status, err)
	} else if len(channelGroups) != 0 {
		res, status, err := pn.HereNow().ChannelGroups(channelGroups).IncludeState(includeState).IncludeUUIDs(includeUUIDs).Execute()
		hereNowResponse(res, status, err)
	} else {
		res, status, err := pn.HereNow().IncludeState(includeState).IncludeUUIDs(includeUUIDs).Execute()
		hereNowResponse(res, status, err)
	}
}

func publishTest() {
	ch := "my-channel"
	for i := 0; i < 1000; i++ {
		go publish(i, ch)
	}
}

func publish(i int, ch string) {
	msg := fmt.Sprintf("Message: %d", i)
	fmt.Println(fmt.Sprintf("%s Publishing to channel: %s", outputPrefix, ch))

	res, status, err := pn.Publish().
		Channel(ch).
		Message(msg).
		Execute()

	if err != nil {
		showErr("Error while publishing: " + err.Error())
	}

	fmt.Println(fmt.Sprintf("%s\nPublish Response: msg %s\n%s %s\n%s", outputPrefix, msg, res, status, outputSuffix))
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
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	meta := map[string]string{
		"m1": "n1",
		"m2": "n2",
	}

	payload := map[string]interface{}{
		"pn_apns": map[string]interface{}{
			"aps": map[string]interface{}{
				"alert": "hi",
				"badge": 2,
				"sound": "melody",
			},
		},
		"pn_gcm": map[string]interface{}{
			"c": "1",
		},
		"b": "2",
	}
	fmt.Println(payload)

	for _, ch := range channels {
		fmt.Println(fmt.Sprintf("%s Publishing to channel: %s", outputPrefix, ch))

		res, status, err := pn.Publish().
			Channel(ch).
			Message("Text with 😜 emoji 🎉" + message).
			UsePost(usePost).
			ShouldStore(store).
			Meta(meta).
			DoNotReplicate(repl).QueryParam(queryParam).
			Execute()

		if err != nil {
			showErr("Error while publishing: " + err.Error())
		}

		fmt.Println(fmt.Sprintf("%s Publish Response:", outputPrefix))

		fmt.Println(fmt.Sprintf("%s %s", res, status))
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
			TTL(1).
			Execute()

		if err != nil {
			showErr("Error while publishing: " + err.Error())
		}

		fmt.Println(fmt.Sprintf("%s Publish Response:", outputPrefix))

		fmt.Println(fmt.Sprintf("%s %s", res, status))
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
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	channels := strings.Split(args[1], ",")
	if (len(args)) > 3 {
		fmt.Println("sub with state")
		var state map[string]interface{}
		var v interface{}
		err := json.Unmarshal([]byte(args[3]), &v)
		if err == nil {
			if st, ok := v.(map[string]interface{}); ok {
				state = st
				fmt.Println("state ok")
			} else {

				fmt.Println("!ok", reflect.TypeOf(v))
				showSubscribeWithStateHelp()
				return
			}
		} else {
			fmt.Println("err", err)
			showSubscribeWithStateHelp()
			return
		}
		groups := strings.Split(args[2], ",")
		pn.Subscribe().
			Channels(channels).
			ChannelGroups(groups).
			WithPresence(withPresence).
			State(state).
			Execute()

	} else if (len(args)) > 2 {
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
			QueryParam(queryParam).
			Execute()
	}

}
