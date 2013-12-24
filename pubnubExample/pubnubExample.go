// Package main provides the example implemetation to connect to pubnub api.
// Runs on the console.
package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/pubnub/go/pubnubMessaging"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
	"unicode/utf8"
)

// connectChannels: the conected pubnub channels, multiple channels are stored separated by comma.
var connectChannels = ""

// ssl: true if the ssl is enabled else false.
var ssl bool

// cipher: stores the cipher key set by the user.
var cipher = ""

// uuid stores the custom uuid set by the user.
var uuid = ""

// a boolean to capture user preference of displaying errors.
var displayError = true

// pub instance of pubnubMessaging package.
var pub *pubnubMessaging.Pubnub

// main method to initiate the application in the console.
// Calls the init method to read user input. And starts the read loop to parse user input.
func main() {
	b := Init()
	if b {
		ReadLoop()
	}
	fmt.Println("Exit")
}

// Init asks the user the basic settings to initialize to the pubnub struct.
// Settings include the pubnub channel(s) to connect to.
// Ssl settings
// Cipher key
// Secret Key
// Custom Uuid
// Proxy details
//
// The method returns false if the channel name is not provided.
//
// returns: a bool, true if the user completed the initail settings.
func Init() (b bool) {
	fmt.Println("PubNub Api for go;", pubnubMessaging.VersionInfo())
	fmt.Println("Please enter the channel name(s). Enter multiple channels separated by comma without spaces.")
	reader := bufio.NewReader(os.Stdin)

	line, _, err := reader.ReadLine()
	if err != nil {
		fmt.Println(err)
	} else {
		connectChannels = string(line)
		if strings.TrimSpace(connectChannels) != "" {
			fmt.Println("Channel: ", connectChannels)
			fmt.Println("Enable SSL? Enter y for Yes, n for No.")
			var enableSsl string
			fmt.Scanln(&enableSsl)

			if enableSsl == "y" || enableSsl == "Y" {
				ssl = true
				fmt.Println("SSL enabled")
			} else {
				ssl = false
				fmt.Println("SSL disabled")
			}

			fmt.Println("Please enter a CIPHER key, leave blank if you don't want to use this.")
			fmt.Scanln(&cipher)
			fmt.Println("Cipher: ", cipher)

			fmt.Println("Please enter a Custom UUID, leave blank for default.")
			fmt.Scanln(&uuid)
			fmt.Println("UUID: ", uuid)

			fmt.Println("Display error messages? Enter y for Yes, n for No. Default is Yes")
			var enableErrorMessages = "y"
			fmt.Scanln(&enableErrorMessages)

			if enableErrorMessages == "y" || enableErrorMessages == "Y" {
				displayError = true
				fmt.Println("Error messages will be displayed")
			} else {
				displayError = false
				fmt.Println("Error messages will not be displayed")
			}

			fmt.Println("Enable resume on reconnect? Enter y for Yes, n for No. Default is Yes")
			var enableResumeOnReconnect = "y"
			fmt.Scanln(&enableResumeOnReconnect)

			if enableResumeOnReconnect == "y" || enableResumeOnReconnect == "Y" {
				pubnubMessaging.SetResumeOnReconnect(true)
				fmt.Println("Resume on reconnect enabled")
			} else {
				pubnubMessaging.SetResumeOnReconnect(false)
				fmt.Println("Resume on reconnect disabled")
			}

			fmt.Println("Set subscribe timeout? Enter numerals.")
			var subscribeTimeout = ""
			fmt.Scanln(&subscribeTimeout)
			val, err := strconv.Atoi(subscribeTimeout)
			if err != nil {
				fmt.Println("Entered value is invalid. Using default value.")
			} else {
				pubnubMessaging.SetSubscribeTimeout(int64(val))
			}
			pubnubMessaging.SetOrigin("pubsub.pubnub.com")
			var pubInstance = pubnubMessaging.NewPubnub("demo", "demo", "", cipher, ssl, uuid)
			pub = pubInstance

			SetupProxy()

			return true
		}
		fmt.Println("Channel cannot be empty.")
	}
	return false
}

// SetupProxy asks the user the Proxy details and calls the SetProxy of the pubnubMessaging
// package with the details.
func SetupProxy() {
	fmt.Println("Using Proxy? Enter y to setup.")
	var enableProxy string
	fmt.Scanln(&enableProxy)

	if enableProxy == "y" || enableProxy == "Y" {
		proxyServer := askServer()
		proxyPort := askPort()
		proxyUser := askUser()
		proxyPassword := askPassword()

		pubnubMessaging.SetProxy(proxyServer, proxyPort, proxyUser, proxyPassword)

		fmt.Println("Proxy sever set")
	} else {
		fmt.Println("Proxy not used")
	}
}

// AskServer asks the user to enter the proxy server name or IP.
// It validates the input and returns the value if validated.
func askServer() string {
	var proxyServer string

	fmt.Println("Enter proxy servername or IP.")
	fmt.Scanln(&proxyServer)

	if strings.TrimSpace(proxyServer) == "" {
		fmt.Println("Proxy servername or IP is empty.")
		askServer()
	}
	return proxyServer
}

// AskPort asks the user to enter the proxy port number.
// It validates the input and returns the value if validated.
func askPort() int {
	var proxyPort string

	fmt.Println("Enter proxy port.")
	fmt.Scanln(&proxyPort)

	port, err := strconv.Atoi(proxyPort)
	if (err != nil) || ((port <= 0) || (port > 65536)) {
		fmt.Println("Proxy port is invalid.")
		askPort()
	}
	return port
}

// AskUser asks the user to enter the proxy username.
// returns the value, can be empty.
func askUser() string {
	var proxyUser string

	fmt.Println("Enter proxy username (optional)")
	fmt.Scanln(&proxyUser)

	return proxyUser
}

// AskPassword asks the user to enter the proxy password.
// returns the value, can be empty.
func askPassword() string {
	var proxyPassword string

	fmt.Println("Enter proxy password (optional)")
	fmt.Scanln(&proxyPassword)

	return proxyPassword
}

// AskChannel asks the user to channel name.
// If the channel(s) are not provided the channel(s) provided by the user
// at the beginning will be used.
// returns the read channel(s), or error
func askChannel() (string, error) {
	fmt.Println("Please enter the channel name. Leave empty to use the channel(s) provided at the beginning.")
	reader := bufio.NewReader(os.Stdin)
	channels, _, errReadingChannel := reader.ReadLine()
	if errReadingChannel != nil {
		fmt.Println("Error channel(s): ", errReadingChannel.Error())
		return "", errReadingChannel
	}
	if strings.TrimSpace(string(channels)) == "" {
		fmt.Println("Using channel(s): ", connectChannels)
		return connectChannels, nil
	}
	return string(channels), nil
}

// AskPort asks the user to enter the proxy port number.
// It validates the input and returns the value if validated.
func askNumber(what string) int64 {
	var input string

	fmt.Println("Enter " + what)
	fmt.Scanln(&input)

	//val, err := strconv.(input, 10, 32)
	bi := big.NewInt(0)
	if _, ok := bi.SetString(input, 10); !ok {
		//if (err != nil) {
		fmt.Println(what + " is invalid. Please enter numerals.")
		askNumber(what)
	}
	fmt.Println(bi.Int64())
	return bi.Int64()
}

// UTF16BytesToString converts UTF-16 encoded bytes, in big or little endian byte order,
// to a UTF-8 encoded string.
func utf16BytesToString(b []byte, o binary.ByteOrder) string {
	utf := make([]uint16, (len(b)+(2-1))/2)
	for i := 0; i+(2-1) < len(b); i += 2 {
		utf[i/2] = o.Uint16(b[i:])
	}
	if len(b)/2 < len(utf) {
		utf[len(utf)-1] = utf8.RuneError
	}
	return string(utf16.Decode(utf))
}

// ReadLoop starts an infinite loop to read the user's input.
// Based on the input the respective go routine is called as a parallel process.
func ReadLoop() {
	fmt.Println("")
	fmt.Println("ENTER 1 FOR Subscribe")
	fmt.Println("ENTER 2 FOR Subscribe with timetoken")
	fmt.Println("ENTER 3 FOR Publish")
	fmt.Println("ENTER 4 FOR Presence")
	fmt.Println("ENTER 5 FOR Detailed History")
	fmt.Println("ENTER 6 FOR Here_Now")
	fmt.Println("ENTER 7 FOR Unsubscribe")
	fmt.Println("ENTER 8 FOR Presence-Unsubscribe")
	fmt.Println("ENTER 9 FOR Time")
	fmt.Println("ENTER 10 FOR Disconnect/Retry")
	fmt.Println("ENTER 11 FOR Exit")
	fmt.Println("")
	reader := bufio.NewReader(os.Stdin)

	for {
		var action string
		fmt.Scanln(&action)

		breakOut := false
		switch action {
		case "1":
			fmt.Println("Running Subscribe")
			channels, errReadingChannel := askChannel()
			if errReadingChannel != nil {
				fmt.Println("errReadingChannel: ", errReadingChannel)
			} else {
				go subscribeRoutine(channels, "")
			}
		case "2":
			fmt.Println("Running Subscribe with timetoken")
			channels, errReadingChannel := askChannel()
			if errReadingChannel != nil {
				fmt.Println("errReadingChannel: ", errReadingChannel)
			} else {
				timetoken := askNumber("Timetoken")
				go subscribeRoutine(channels, strconv.FormatInt(timetoken, 10))
			}
		case "111":
			//for test
			fmt.Println("Running Subscribe2")
			channels, errReadingChannel := askChannel()
			if errReadingChannel != nil {
				fmt.Println("errReadingChannel: ", errReadingChannel)
			} else {
				go subscribeRoutine2(channels, "")
			}
		case "3":
			channels, errReadingChannel := askChannel()
			if errReadingChannel != nil {
				fmt.Println("errReadingChannel: ", errReadingChannel)
			} else {
				fmt.Println("Please enter the message")
				message, _, err := reader.ReadLine()
				if err != nil {
					fmt.Println(err)
				} else {
					go publishRoutine(channels, string(message))
				}
			}
		case "4":
			fmt.Println("Running Presence")
			go presenceRoutine()
		case "333":
			//for test
			fmt.Println("Running Presence2")
			go presenceRoutine2()
		case "5":
			fmt.Println("Running detailed history")
			go detailedHistoryRoutine()
		case "6":
			fmt.Println("Running here now")
			go hereNowRoutine()
		case "7":
			fmt.Println("Running Unsubscribe")
			channels, errReadingChannel := askChannel()
			if errReadingChannel != nil {
				fmt.Println("errReadingChannel: ", errReadingChannel)
			} else {
				go unsubscribeRoutine(channels)
			}
		case "8":
			fmt.Println("Running Unsubscribe Presence")
			go unsubscribePresenceRoutine()
		case "9":
			fmt.Println("Running Time")
			go timeRoutine()
		case "10":
			fmt.Println("Disconnect/Retry")
			pub.CloseExistingConnection()
		case "11":
			fmt.Println("Exiting")
			pub.Abort()
			time.Sleep(3 * time.Second)
			breakOut = true
		default:
			fmt.Println("Invalid choice!")
		}
		if breakOut {
			break
		} else {
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

// ParseResponseSubscribe parses the response of the Subscribed pubnub channel.
// It prints the response as-is in the console.
func parseErrorResponse(channel chan []byte) {
	for {
		value, ok := <-channel
		if !ok {
			fmt.Println("")
			break
		}
		if string(value) != "[]" {
			if displayError {
				fmt.Println(fmt.Sprintf("Error Callback: %s", value))
				fmt.Println("")
			}
		}
	}
}

//for test
func parseErrorResponse2(channel chan []byte) {
	for {
		value, ok := <-channel
		if !ok {
			fmt.Println("")
			break
		}
		if string(value) != "[]" {
			fmt.Println(fmt.Sprintf("Error Callback2: %s", value))
			fmt.Println("")
		}
	}
}

// ParseResponseSubscribe parses the response of the Subscribed pubnub channel.
// It prints the response as-is in the console.
func parseResponseSubscribe(channel chan []byte) {
	for {
		value, ok := <-channel
		if !ok {
			fmt.Println("")
			break
		}
		if string(value) != "[]" {
			fmt.Println(fmt.Sprintf("Subscribe: %s", value))
			fmt.Println("")
		}
	}
}

// ParseResponseSubscribe parses the response of the Subscribed pubnub channel.
// It prints the response as-is in the console.
func parseResponseSubscribe2(channel chan []byte) {
	for {
		value, ok := <-channel
		if !ok {
			fmt.Println("")
			break
		}
		if string(value) != "[]" {
			fmt.Println(fmt.Sprintf("Subscribe2: %s", value))
			fmt.Println("")
		}
	}
}

// ParseResponsePresence parses the response of the presence subscription pubnub channel.
// It prints the response as-is in the console.
func parseResponsePresence(channel chan []byte) {
	for {
		value, ok := <-channel
		if !ok {
			break
		}
		if string(value) != "[]" {
			fmt.Println(fmt.Sprintf("Presence: %s ", value))
			fmt.Println("")
		}
	}
}

// ParseResponsePresence parses the response of the presence subscription pubnub channel.
// It prints the response as-is in the console.
// for test
func parseResponsePresence2(channel chan []byte) {
	for {
		value, ok := <-channel
		if !ok {
			break
		}
		if string(value) != "[]" {
			fmt.Println(fmt.Sprintf("Presence2: %s ", value))
			fmt.Println("")
		}
	}
}

// ParseResponse parses the response of all the other activities apart
// from subscribe and presence on the pubnub channel.
// It prints the response as-is in the console.
func parseResponse(channel chan []byte) {
	for {
		value, ok := <-channel
		if !ok {
			break
		}
		if string(value) != "[]" {
			fmt.Println(fmt.Sprintf("Response: %s ", value))
			fmt.Println("")
		}
	}
}

// ParseUnsubResponse parses the response.
// It prints the response as-is in the console.
func parseUnsubResponse(channel chan []byte) {
	for {
		value, ok := <-channel
		if !ok {
			break
		}
		if string(value) != "[]" {
			fmt.Println(fmt.Sprintf("Unsub Response: %s ", value))
			fmt.Println("")
		}
	}
}

// SubscribeRoutine calls the Subscribe routine of the pubnubMessaging package
// as a parallel process.
func subscribeRoutine(channels string, timetoken string) {
	var errorChannel = make(chan []byte)
	var subscribeChannel = make(chan []byte)
	go pub.Subscribe(channels, timetoken, subscribeChannel, false, errorChannel)
	go parseResponseSubscribe(subscribeChannel)
	go parseErrorResponse(errorChannel)
}

// SubscribeRoutine calls the Subscribe routine of the pubnubMessaging package
// as a parallel process.
func subscribeRoutine2(channels string, timetoken string) {
	var errorChannel = make(chan []byte)
	var subscribeChannel = make(chan []byte)
	go pub.Subscribe(channels, timetoken, subscribeChannel, false, errorChannel)
	go parseResponseSubscribe2(subscribeChannel)
	go parseErrorResponse2(errorChannel)
}

// PublishRoutine asks the user the message to send to the pubnub channel(s) and
// calls the Publish routine of the pubnubMessaging package as a parallel
// process. If we have multiple pubnub channels then this method will spilt the
// channel by comma and send the message on all the pubnub channels.
func publishRoutine(channels string, message string) {
	var errorChannel = make(chan []byte)
	channelArray := strings.Split(channels, ",")

	for i := 0; i < len(channelArray); i++ {
		ch := strings.TrimSpace(channelArray[i])
		fmt.Println("Publish to channel: ", ch)
		channel := make(chan []byte)
		go pub.Publish(ch, message, channel, errorChannel)
		go parseResponse(channel)
		go parseErrorResponse(errorChannel)
	}
}

// PresenceRoutine calls the Subscribe routine of the pubnubMessaging package,
// by setting the last argument as true, as a parallel process.
func presenceRoutine() {
	var errorChannel = make(chan []byte)
	var presenceChannel = make(chan []byte)
	go pub.Subscribe(connectChannels, "", presenceChannel, true, errorChannel)
	go parseResponsePresence(presenceChannel)
	go parseErrorResponse(errorChannel)
}

// for test
func presenceRoutine2() {
	var errorChannel = make(chan []byte)
	var presenceChannel = make(chan []byte)
	go pub.Subscribe(connectChannels, "", presenceChannel, true, errorChannel)
	go parseResponsePresence2(presenceChannel)
	go parseErrorResponse2(errorChannel)
}

// DetailedHistoryRoutine calls the History routine of the pubnubMessaging package as a parallel
// process. If we have multiple pubnub channels then this method will spilt the _connectChannels
// by comma and send the message on all the pubnub channels.
func detailedHistoryRoutine() {
	var errorChannel = make(chan []byte)
	channelArray := strings.Split(connectChannels, ",")
	for i := 0; i < len(channelArray); i++ {
		ch := strings.TrimSpace(channelArray[i])
		fmt.Println("DetailedHistory for channel: ", ch)

		channel := make(chan []byte)

		//go _pub.History(ch, 100, 13662867154115803, 13662867243518473, false, channel)
		go pub.History(ch, 100, 0, 0, false, channel, errorChannel)
		go parseResponse(channel)
		go parseErrorResponse(errorChannel)
	}
}

// HereNowRoutine calls the HereNow routine of the pubnubMessaging package as a parallel
// process. If we have multiple pubnub channels then this method will spilt the _connectChannels
// by comma and send the message on all the pubnub channels.
func hereNowRoutine() {
	var errorChannel = make(chan []byte)
	channelArray := strings.Split(connectChannels, ",")
	for i := 0; i < len(channelArray); i++ {
		channel := make(chan []byte)
		ch := strings.TrimSpace(channelArray[i])
		fmt.Println("HereNow for channel: ", ch)

		go pub.HereNow(ch, channel, errorChannel)
		go parseResponse(channel)
		go parseErrorResponse(errorChannel)
	}
}

// UnsubscribeRoutine calls the Unsubscribe routine of the pubnubMessaging package as a parallel
// process. All the channels in the _connectChannels string will be unsubscribed.
func unsubscribeRoutine(channels string) {
	var errorChannel = make(chan []byte)
	channel := make(chan []byte)

	go pub.Unsubscribe(channels, channel, errorChannel)
	parseUnsubResponse(channel)
}

// UnsubscribePresenceRoutine calls the PresenceUnsubscribe routine of the pubnubMessaging package as a parallel
// process. All the channels in the _connectChannels string will be unsubscribed.
func unsubscribePresenceRoutine() {
	var errorChannel = make(chan []byte)
	channel := make(chan []byte)

	go pub.PresenceUnsubscribe(connectChannels, channel, errorChannel)
	parseResponse(channel)
}

// TimeRoutine calls the GetTime routine of the pubnubMessaging package as a parallel
// process.
func timeRoutine() {
	var errorChannel = make(chan []byte)
	channel := make(chan []byte)
	go pub.GetTime(channel, errorChannel)
	go parseResponse(channel)
	go parseErrorResponse(errorChannel)
}
