package main

import (
	"bufio"
	"errors"
	"fmt"
	//"io/ioutil"
	"encoding/json"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	pubnub "github.com/sprucehealth/pubnub-go"
)

import _ "net/http/pprof"
import "net/http"

var config *pubnub.Config
var pn *pubnub.PubNub
var quitSubscribe = false

const outputPrefix = "\x1b[32;1m Example >>>> \x1b[0m"
const outputSuffix = "\x1b[32;2m Example <<<< \x1b[0m"

func main() {
	connect()
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
	//config.Log = log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile)
	//config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	config.Log = infoLogger
	config.Log.SetPrefix("PubNub :->  ")
	config.PublishKey = "pub-c-3ed95c83-12e6-4cda-9d69-c47ba2abb57e"   //"demo"   //"demo"
	config.SubscribeKey = "sub-c-26a73b0a-c3f2-11e9-8b24-569e8a5c3af3" //"demo" //"sub-c-10b61350-bec7-11e9-a375-f698c1d99dce" //"demo" //
	//config.SecretKey = //"pam"    //"demo"

	//config.PublishKey = "pub-c-cdea0ef1-c571-4b72-b43f-ff1dc8aa4c5d"
	//config.SubscribeKey = "sub-c-4757f09c-c3f2-11e9-9d00-8a58a5558306"
	//config.SecretKey = "sec-c-YTYxNzVjYzctNDY2MS00N2NmLTg2NjYtNGRlNWY1NjMxMDBm"

	//config.AuthKey = "akey"

	config.CipherKey = "enigma"
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
			case userEvent := <-listener.UserEvent:
				fmt.Print(fmt.Sprintf("%s Subscribe Response:", outputPrefix))
				fmt.Println(" --- UserEvent: ")
				fmt.Println(fmt.Sprintf("%s %s", outputPrefix, userEvent))
				fmt.Println(fmt.Sprintf("%s userEvent.Channel: %s", outputPrefix, userEvent.Channel))
				fmt.Println(fmt.Sprintf("%s userEvent.SubscribedChannel: %s", outputPrefix, userEvent.SubscribedChannel))
				fmt.Println(fmt.Sprintf("%s userEvent.Event: %s", outputPrefix, userEvent.Event))
				fmt.Println(fmt.Sprintf("%s userEvent.UserID: %s", outputPrefix, userEvent.UserID))
				fmt.Println(fmt.Sprintf("%s userEvent.Description: %s", outputPrefix, userEvent.Description))
				fmt.Println(fmt.Sprintf("%s userEvent.Timestamp: %s", outputPrefix, userEvent.Timestamp))
				fmt.Println(fmt.Sprintf("%s userEvent.Name: %s", outputPrefix, userEvent.Name))
				fmt.Println(fmt.Sprintf("%s userEvent.ExternalID: %s", outputPrefix, userEvent.ExternalID))
				fmt.Println(fmt.Sprintf("%s userEvent.ProfileURL: %s", outputPrefix, userEvent.ProfileURL))
				fmt.Println(fmt.Sprintf("%s userEvent.Email: %s", outputPrefix, userEvent.Email))
				fmt.Println(fmt.Sprintf("%s userEvent.Created: %s", outputPrefix, userEvent.Created))
				fmt.Println(fmt.Sprintf("%s userEvent.Updated: %s", outputPrefix, userEvent.Updated))
				fmt.Println(fmt.Sprintf("%s userEvent.ETag: %s", outputPrefix, userEvent.ETag))
				fmt.Println(fmt.Sprintf("%s userEvent.Custom: %v", outputPrefix, userEvent.Custom))

			case spaceEvent := <-listener.SpaceEvent:
				fmt.Print(fmt.Sprintf("%s Subscribe Response:", outputPrefix))
				fmt.Println(" --- SpaceEvent: ")
				fmt.Println(fmt.Sprintf("%s %s", outputPrefix, spaceEvent))
				fmt.Println(fmt.Sprintf("%s spaceEvent.Channel: %s", outputPrefix, spaceEvent.Channel))
				fmt.Println(fmt.Sprintf("%s spaceEvent.SubscribedChannel: %s", outputPrefix, spaceEvent.SubscribedChannel))
				fmt.Println(fmt.Sprintf("%s spaceEvent.Event: %s", outputPrefix, spaceEvent.Event))
				fmt.Println(fmt.Sprintf("%s spaceEvent.SpaceID: %s", outputPrefix, spaceEvent.SpaceID))
				fmt.Println(fmt.Sprintf("%s spaceEvent.Description: %s", outputPrefix, spaceEvent.Description))
				fmt.Println(fmt.Sprintf("%s spaceEvent.Timestamp: %s", outputPrefix, spaceEvent.Timestamp))
				fmt.Println(fmt.Sprintf("%s spaceEvent.Created: %s", outputPrefix, spaceEvent.Created))
				fmt.Println(fmt.Sprintf("%s spaceEvent.Updated: %s", outputPrefix, spaceEvent.Updated))
				fmt.Println(fmt.Sprintf("%s spaceEvent.ETag: %s", outputPrefix, spaceEvent.ETag))
				fmt.Println(fmt.Sprintf("%s spaceEvent.Custom: %v", outputPrefix, spaceEvent.Custom))

			case membershipEvent := <-listener.MembershipEvent:
				fmt.Print(fmt.Sprintf("%s Subscribe Response:", outputPrefix))
				fmt.Println(" --- MembershipEvent: ")
				fmt.Println(fmt.Sprintf("%s %s", outputPrefix, membershipEvent))
				fmt.Println(fmt.Sprintf("%s membershipEvent.Channel: %s", outputPrefix, membershipEvent.Channel))
				fmt.Println(fmt.Sprintf("%s membershipEvent.SubscribedChannel: %s", outputPrefix, membershipEvent.SubscribedChannel))
				fmt.Println(fmt.Sprintf("%s membershipEvent.Event: %s", outputPrefix, membershipEvent.Event))
				fmt.Println(fmt.Sprintf("%s membershipEvent.SpaceID: %s", outputPrefix, membershipEvent.SpaceID))
				fmt.Println(fmt.Sprintf("%s membershipEvent.UserID: %s", outputPrefix, membershipEvent.UserID))
				fmt.Println(fmt.Sprintf("%s membershipEvent.Description: %s", outputPrefix, membershipEvent.Description))
				fmt.Println(fmt.Sprintf("%s membershipEvent.Timestamp: %s", outputPrefix, membershipEvent.Timestamp))
				fmt.Println(fmt.Sprintf("%s membershipEvent.Custom: %v", outputPrefix, membershipEvent.Custom))
			}
		}
	}()

	pn.AddListener(listener)
	showHelp()

	/*config2 := pubnub.NewConfig()
	config2.SubscribeRequestTimeout = 59
	config2.UUID = "GlobalSubscriber"
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
	showSubscribeWithStateHelp()
	showPresenceTimeoutHelp()
	showPresenceHelp()
	showMessageCountsHelp()
	showSignalHelp()
	showCreateUserHelp()
	showGetUsersHelp()
	showEditMembershipsHelp()
	showUpdateMembersHelp()
	showGetSpaceMembershipsHelp()
	showGetMembersHelp()
	showGetSpacesHelp()
	showUpdateSpaceHelp()
	showDeleteSpaceHelp()
	showCreateSpaceHelp()
	showGetSpaceHelp()
	showDeleteUserHelp()
	showUpdateUserHelp()
	showGetUserHelp()
	fmt.Println("")
	fmt.Println("================")
	fmt.Println(" ||  COMMANDS  ||")
	fmt.Println("================")
	fmt.Println("")
	fmt.Println(" UNSUBSCRIBE ALL \n\tq ")
	fmt.Println(" QUIT \n\tctrl+c ")
}

func showEditMembershipsHelp() {
	fmt.Println(" EditMemberships EXAMPLE: ")
	fmt.Println("	updatespacemem spaceid id a/u/r limit count")
	fmt.Println("	updatespacemem id0 id1 a 100 true")

}
func showUpdateMembersHelp() {
	fmt.Println(" UpdateMembers EXAMPLE: ")
	fmt.Println("	updatemem memebers id a/u/r limit count")
	fmt.Println("	updatemem id0 id0 a 100 true")
}

func showGetSpaceMembershipsHelp() {
	fmt.Println(" GetSpaceMemberships EXAMPLE: ")
	fmt.Println("	getspacemem spaceid limit count start")
	fmt.Println("	getspacemem id0 100 true Mymx")

}
func showGetMembersHelp() {
	fmt.Println(" GetMembers EXAMPLE: ")
	fmt.Println("	getmem userid limit count start")
	fmt.Println("	getmem id0 100 true Mymx")

}
func showGetSpacesHelp() {
	fmt.Println(" GetSpaces EXAMPLE: ")
	fmt.Println("	getspaces limit count start")
	fmt.Println("	getspaces 100 true MjWn")

}
func showUpdateSpaceHelp() {
	fmt.Println(" UpdateSpace EXAMPLE: ")
	fmt.Println("	updatespace id name desc")
	fmt.Println("	updatespace id0 name desc")

}
func showDeleteSpaceHelp() {
	fmt.Println(" DeleteSpace EXAMPLE: ")
	fmt.Println("	delspace id")
	fmt.Println("	delspace id0")

}
func showCreateSpaceHelp() {
	fmt.Println(" CreateSpace EXAMPLE: ")
	fmt.Println("	createspace id name desc")
	fmt.Println("	createspace id0 name desc")

}
func showGetSpaceHelp() {
	fmt.Println(" GetSpace EXAMPLE: ")
	fmt.Println("	getspace id")
	fmt.Println("	getspace id0")

}
func showDeleteUserHelp() {
	fmt.Println(" DeleteUser EXAMPLE: ")
	fmt.Println("	deleteuser id")
	fmt.Println("	deleteuser id0")
}

func showUpdateUserHelp() {
	fmt.Println(" UpdateUser EXAMPLE: ")
	fmt.Println("	updateuser id name extid url email")
	fmt.Println("	updateuser id0 name extid purl email")
}

func showGetUserHelp() {
	fmt.Println(" GetUser EXAMPLE: ")
	fmt.Println("	getuser id")
	fmt.Println("	getuser id0")
}

func showMessageCountsHelp() {
	fmt.Println(" MessageCounts EXAMPLE: ")
	fmt.Println("	messagecounts Channel(s) timetoken1,timetoken2")
	fmt.Println("	messagecounts my-channel,my-channel1 15210190573608384,15211140747622125")
}

func showGetUsersHelp() {
	fmt.Println(" GetUsers EXAMPLE: ")
	fmt.Println("	getusers limit count start")
	fmt.Println("	getusers 100 true MjWn")
}

func showCreateUserHelp() {
	fmt.Println(" CreateUser EXAMPLE: ")
	fmt.Println("	createuser id name extid url email")
	fmt.Println("	createuser id0 name extid purl email")
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

func showGrantHelp() {
	fmt.Println(" GRANT EXAMPLE: ")
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
	case "createuser":
		createUser(command[1:])
	case "getusers":
		getUsers(command[1:])
	case "getuser":
		getUser(command[1:])
	case "updateuser":
		updateUser(command[1:])
	case "deleteuser":
		deleteUser(command[1:])
	case "getspaces":
		getSpaces(command[1:])
	case "createspace":
		createSpace(command[1:])
	case "delspace":
		deleteSpace(command[1:])
	case "updatespace":
		updateSpace(command[1:])
	case "getspace":
		getSpace(command[1:])
	case "getspacemem":
		getSpaceMemberships(command[1:])
	case "getmem":
		getMembers(command[1:])
	case "updatespacemem":
		manageMemberships(command[1:])
	case "updatemem":
		manageMembers(command[1:])
	case "q":
		pn.UnsubscribeAll()
	case "d":
		pn.Destroy()
	default:
		showHelp()
	}
}

func manageMembers(args []string) {
	if len(args) < 5 {
		showEditMembershipsHelp()
		return
	}
	spaceID := args[0]
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

	incl := []pubnub.PNMembersInclude{
		pubnub.PNMembersCustom,
		pubnub.PNMembersUser,
		pubnub.PNMembersUserCustom,
	}

	custom := make(map[string]interface{})
	custom["a1"] = "b1"
	custom["c1"] = "d1"

	in := pubnub.PNMembersInput{
		ID:     id0,
		Custom: custom,
	}

	inArr := []pubnub.PNMembersInput{
		in,
	}

	custom2 := make(map[string]interface{})
	custom2["a2"] = "b2"
	custom2["c2"] = "d2"

	up := pubnub.PNMembersInput{
		ID:     id0,
		Custom: custom2,
	}

	upArr := []pubnub.PNMembersInput{
		up,
	}

	re := pubnub.PNMembersRemove{
		ID: id0,
	}

	reArr := []pubnub.PNMembersRemove{
		re,
	}

	if action == "a" {
		reArr = []pubnub.PNMembersRemove{}
		upArr = []pubnub.PNMembersInput{}
	} else if action == "u" {
		reArr = []pubnub.PNMembersRemove{}
		inArr = []pubnub.PNMembersInput{}
	} else if action == "r" {
		upArr = []pubnub.PNMembersInput{}
		inArr = []pubnub.PNMembersInput{}
	}

	if start != "" {
		res, status, err := pn.ManageMembers().SpaceID(spaceID).Add(inArr).Update(upArr).Remove(reArr).Include(incl).Limit(limit).Count(count).Start(start).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		fmt.Println("res", res)

	} else {
		res, status, err := pn.ManageMembers().SpaceID(spaceID).Add(inArr).Update(upArr).Remove(reArr).Include(incl).Limit(limit).Count(count).Execute()
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
		pubnub.PNMembershipsCustom,
		pubnub.PNMembershipsSpace,
		pubnub.PNMembershipsSpaceCustom,
	}

	custom3 := make(map[string]interface{})
	custom3["a3"] = "b3"
	custom3["c3"] = "d3"

	in := pubnub.PNMembershipsInput{
		ID:     id0,
		Custom: custom3,
	}

	inArr := []pubnub.PNMembershipsInput{
		in,
	}

	custom4 := make(map[string]interface{})
	custom4["a4"] = "b4"
	custom4["c4"] = "d4"

	up := pubnub.PNMembershipsInput{
		ID:     id0,
		Custom: custom4,
	}

	upArr := []pubnub.PNMembershipsInput{
		up,
	}

	re := pubnub.PNMembershipsRemove{
		ID: id0,
	}

	reArr := []pubnub.PNMembershipsRemove{
		re,
	}

	if action == "a" {
		reArr = []pubnub.PNMembershipsRemove{}
		upArr = []pubnub.PNMembershipsInput{}
	} else if action == "u" {
		reArr = []pubnub.PNMembershipsRemove{}
		inArr = []pubnub.PNMembershipsInput{}
	} else if action == "r" {
		upArr = []pubnub.PNMembershipsInput{}
		inArr = []pubnub.PNMembershipsInput{}
	}

	if start != "" {
		res, status, err := pn.ManageMemberships().Add(inArr).Update(upArr).Remove(reArr).Include(incl).Limit(limit).Count(count).Start(start).Execute()
		//res, status, err := pn.UpdateMembers().UserID(userID).Add(inArr).Include(incl).Limit(limit).Count(count).Start(start).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		fmt.Println("res", res)

	} else {

		res, status, err := pn.ManageMemberships().UserID(userID).Add(inArr).Update(upArr).Remove(reArr).Include(incl).Limit(limit).Count(count).Execute()
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
		pubnub.PNMembershipsCustom,
		pubnub.PNMembershipsSpace,
		pubnub.PNMembershipsSpaceCustom,
	}
	if start != "" {
		res, status, err := pn.GetMemberships().UserID(id).Include(incl).Limit(limit).Count(count).Start(start).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		fmt.Println("res", res)
	} else {
		res, status, err := pn.GetMemberships().UserID(id).Include(incl).Limit(limit).Count(count).Execute()
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

	incl := []pubnub.PNMembersInclude{
		pubnub.PNMembersCustom,
		pubnub.PNMembersUser,
		pubnub.PNMembersUserCustom,
	}
	if start != "" {
		res, status, err := pn.GetMembers().SpaceID(id).Include(incl).Limit(limit).Count(count).Start(start).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		fmt.Println("res", res)
	} else {
		res, status, err := pn.GetMembers().SpaceID(id).Include(incl).Limit(limit).Count(count).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		fmt.Println("res", res)
	}
}

func getSpaces(args []string) {
	if len(args) < 2 {
		showGetSpacesHelp()
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

	incl := []pubnub.PNUserSpaceInclude{
		pubnub.PNUserSpaceCustom,
	}

	if start != "" {
		res, status, err := pn.GetSpaces().Include(incl).Limit(limit).Count(count).Start(start).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		fmt.Println("res", res)
	} else {
		res, status, err := pn.GetSpaces().Include(incl).Limit(limit).Count(count).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		fmt.Println("res", res)
	}
}

func updateSpace(args []string) {
	if len(args) < 2 {
		showUpdateSpaceHelp()
		return
	}
	id := args[0]
	name := args[1]
	desc := args[2]

	custom := make(map[string]interface{})
	custom["a"] = "b"
	custom["c"] = "d"

	incl := []pubnub.PNUserSpaceInclude{
		pubnub.PNUserSpaceCustom,
	}

	res, status, err := pn.UpdateSpace().ID(id).Name(name).Description(desc).Include(incl).Custom(custom).Execute()
	fmt.Println("status", status)
	fmt.Println("err", err)
	fmt.Println("res", res)
}

func deleteSpace(args []string) {
	if len(args) < 1 {
		showDeleteSpaceHelp()
		return
	}
	id := args[0]

	res, status, err := pn.DeleteSpace().ID(id).Execute()
	fmt.Println("status", status)
	fmt.Println("err", err)
	fmt.Println("res", res)
}

func createSpace(args []string) {
	if len(args) < 3 {
		showCreateSpaceHelp()
		return
	}
	id := args[0]
	name := args[1]
	desc := args[2]

	custom := make(map[string]interface{})
	custom["a"] = "b"

	incl := []pubnub.PNUserSpaceInclude{
		pubnub.PNUserSpaceCustom,
		pubnub.PNUserSpaceCustom,
	}

	//res, status, err := pn.CreateSpace().ID("id0").Name("name").Description("desc").Include([]string{"custom"}).Custom(custom).Execute()
	res, status, err := pn.CreateSpace().ID(id).Name(name).Description(desc).Include(incl).Custom(custom).Execute()
	fmt.Println("status", status)
	fmt.Println("err", err)
	fmt.Println("res", res)
}

func getSpace(args []string) {
	if len(args) < 1 {
		showGetSpaceHelp()
		return
	}
	id := args[0]

	incl := []pubnub.PNUserSpaceInclude{
		pubnub.PNUserSpaceCustom,
	}

	res, status, err := pn.GetSpace().ID(id).Include(incl).Execute()
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

	res, status, err := pn.DeleteUser().ID(id).Execute()
	fmt.Println("status", status)
	fmt.Println("err", err)
	fmt.Println("res", res)
}

func updateUser(args []string) {
	if len(args) < 5 {
		showUpdateUserHelp()
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
	incl := []pubnub.PNUserSpaceInclude{
		pubnub.PNUserSpaceCustom,
	}

	res, status, err := pn.UpdateUser().Include(incl).ID(id).Name(name).ExternalID(extid).ProfileURL(purl).Email(email).Custom(custom).Execute()
	fmt.Println("status", status)
	fmt.Println("err", err)
	fmt.Println("res", res)
}

func getUser(args []string) {
	if len(args) < 1 {
		showGetUserHelp()
		return
	}
	id := args[0]

	incl := []pubnub.PNUserSpaceInclude{
		pubnub.PNUserSpaceCustom,
	}

	res, status, err := pn.GetUser().Include(incl).ID(id).Execute()
	fmt.Println("status", status)
	fmt.Println("err", err)
	fmt.Println("res", res)
}

func getUsers(args []string) {
	if len(args) < 2 {
		showGetUsersHelp()
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

	incl := []pubnub.PNUserSpaceInclude{
		pubnub.PNUserSpaceCustom,
	}

	if start != "" {
		res, status, err := pn.GetUsers().Include(incl).Start(start).Limit(limit).Count(count).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		fmt.Println("res", res)

	} else {
		res, status, err := pn.GetUsers().Include(incl).Limit(limit).Count(count).Execute()
		fmt.Println("status", status)
		fmt.Println("err", err)
		fmt.Println("res", res)
	}

}

func createUser(args []string) {
	if len(args) < 5 {
		showCreateUserHelp()
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

	incl := []pubnub.PNUserSpaceInclude{
		pubnub.PNUserSpaceCustom,
	}

	res, status, err := pn.CreateUser().Include(incl).ID(id).Name(name).ExternalID(extid).ProfileURL(purl).Email(email).Custom(custom).Execute()
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

func grant(args []string) {
	if len(args) < 6 {
		fmt.Println(len(args))
		showAddToCgHelp()
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

	if (end != 0) && (start != 0) {
		res, status, err := pn.Fetch().
			Channels(channels).
			Count(count).
			Start(start).
			End(end).
			Reverse(reverse).
			Execute()
		parseFetch(res, status, err)
	} else if start != 0 {
		res, status, err := pn.Fetch().
			Channels(channels).
			Count(count).
			Start(start).
			Reverse(reverse).
			Execute()
		parseFetch(res, status, err)
	} else if end != 0 {
		res, status, err := pn.Fetch().
			Channels(channels).
			Count(count).
			End(end).
			Reverse(reverse).
			Execute()
		parseFetch(res, status, err)
	} else {
		res, status, err := pn.Fetch().
			Channels(channels).
			Count(count).
			Reverse(reverse).
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

	if (end != 0) && (start != 0) {
		res, status, err := pn.History().
			Channel(channel).
			Count(count).
			Start(start).
			End(end).
			IncludeTimetoken(includeTimetoken).
			Reverse(reverse).
			Execute()
		parseHistory(res, status, err)
	} else if start != 0 {
		res, status, err := pn.History().
			Channel(channel).
			Count(count).
			Start(start).
			IncludeTimetoken(includeTimetoken).
			Reverse(reverse).
			Execute()
		parseHistory(res, status, err)
	} else if end != 0 {
		res, status, err := pn.History().
			Channel(channel).
			Count(count).
			End(end).
			IncludeTimetoken(includeTimetoken).
			Reverse(reverse).
			Execute()
		parseHistory(res, status, err)
	} else {
		res, status, err := pn.History().
			Channel(channel).
			Count(count).
			IncludeTimetoken(includeTimetoken).
			Reverse(reverse).
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

	for _, ch := range channels {
		fmt.Println(fmt.Sprintf("%s Publishing to channel: %s", outputPrefix, ch))
		res, status, err := pn.Publish().
			Channel(ch).
			Message(res).
			UsePost(usePost).
			ShouldStore(store).
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
