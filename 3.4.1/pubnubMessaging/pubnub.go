// Package pubnubMessaging provides the implemetation to connect to pubnub api.
// Build Date: Sep 9, 2013
// Version: 3.4.1
package pubnubMessaging

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
    "time"
    "net"
    "crypto/tls"
    "net/url"
    "reflect"
    "bytes"
    "strconv"
    "crypto/aes"
    "crypto/cipher"
    "crypto/hmac"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "io"
)

// This string is appended to all presence channels 
// to differentiate from the subscribe requests.
const _presenceSuffix = "-pnpres"

// The string is used when the server returns a malformed or non-JSON response.
const _invalidJson = "Invalid JSON"

// The string is returned as a message when the http request times out.
const _operationTimeout = "Operation Timeout"

// The string is returned as a message when the http request is aborted.
const _connectionAborted = "Connection aborted"

// The string is encountered when the http request couldn't connect to the origin.
const _noSuchHost = "no such host"

// The string is returned as a message when network connection is not avaialbe.
const _networkUnavailable = "Network unavailable"

// The string is encountered when the http request faces connectivity issues.
const _closedNetworkConnection = "closed network connection"

// The string is encountered when the http request faces connectivity issues.
const _connectionResetByPeer = "connection reset by peer"

// The string is returned as a message when the http request encounters network connectivity issues.
const _connectionResetByPeerU = "Connection reset by peer"

// The string is encountered the http request times out.
const _timeout = "timeout"

// The string is returned as a message when the http request times out.
const _timeoutU = "Timeout"

// The string is retured when the client faces issues in initializing the transport.
const _errorInInitializing = "Error in initializing connection: "

// The string is used when the server returns a non 200 response on publish 
const _publishFailed = "Publish Failed"

// The time after which the Publish/HereNow/DetailedHitsory/Unsubscribe/
// UnsibscribePresence/Time  request will timeout.
// In seconds.
const _nonSubscribeTimeout = 5 //sec

// On Subscribe/Presence timeout, the number of times the reconnect attempts are made.
const _maxRetries = 50 //times

// The delay in the reconnect attempts on timeout.
// In seconds.
const _retryInterval = 10 //sec

// The HTTP transport Dial timeout.
// In seconds.
const _connectTimeout = 10 //sec

// Root url value of pubnub api without the http/https protocol.
var _origin = "pubsub.pubnub.com"

// The time after which the Subscribe/Presence request will timeout.
// In seconds.
var _subscribeTimeout int64 = 310 //sec

// If _resumeOnReconnect is TRUE, then upon reconnect, 
// it should use the last successfully retrieved timetoken. 
// This has the effect of continuing, or “catching up” to missed traffic.
// If resumeOnReconnect is FALSE, then upon reconnect, 
// it should use a 0 (zero) timetoken. 
// This has the effect of continuing from “this moment onward”. 
// Any messages received since the previous timeout or network error are skipped.
var _resumeOnReconnect = true 

// Global variable to reuse a commmon connection instance for non subscribe requests 
// Publish/HereNow/DetailedHitsory/Unsubscribe/UnsibscribePresence/Time.
var _conn net.Conn

// Global variable to reuse a commmon connection instance for Subscribe/Presence requests.
var _subscribeConn net.Conn

// Global variable to reuse a commmon transport instance for Subscribe/Presence requests.
var _subscribeTransport http.RoundTripper

// Global variable to reuse a commmon transport instance for non subscribe requests 
// Publish/HereNow/DetailedHitsory/Unsubscribe/UnsibscribePresence/Time.
var _transport http.RoundTripper

// No of retries made since disconnection.
var _retryCount = 0

// Global variable to store the proxy server if set.
var _proxyServer string

// Global variable to store the proxy port if set.
var _proxyPort int

// Global variable to store the proxy username if set.
var _proxyUser string

// Global variable to store the proxy password if set.
var _proxyPassword string

// Global variable to check if the proxy server if used.    
var _proxyServerEnabled = false

// 16 byte IV  
var _IV = "0123456789012345"

type ResponseStruct struct {
    Message []interface{}
    Timetoken string
    ChannelName string
}

// Pubnub structure.  
// Origin stores the root url value of pubnub api in the current instance.
// PublishKey stores the user specific Publish Key in the current instance.
// SubscribeKey stores the user specific Subscribe Key in the current instance.
// SecretKey stores the user specific Secret Key in the current instance.
// CipherKey stores the user specific Cipher Key in the current instance.
// SSL is true if enabled, else is false for the current instance.
// Uuid is the unique identifier, it can be a custom value or is automatically generated.
// SubscribedChannels keeps a list of subscribed Pubnub channels by the user in the a comma separated string.
// TimeToken is the current value of the servertime. This will be used to appened in each request.
// SentTimeToken: This is the timetoken sent to the server with the request
// ResetTimeToken: In case of a new request or an error this variable is set to true so that the 
// timeToken will be set to 0 in the next request.
// PresenceChannels: All the presence responses will be routed to this channel. It stores the response channels for 
// each pubnub channel as map using the pubnub channel name as the key. 
// SubscribeChannels: All the subscribe responses will be routed to this channel. It stores the response channels for 
// each pubnub channel as map using the pubnub channel name as the key.
// PresenceErrorChannels: All the presence error responses will be routed to this channel. It stores the response channels for 
// each pubnub channel as map using the pubnub channel name as the key.
// SubscribeErrorChannels: All the subscribe error responses will be routed to this channel. It stores the response channels for 
// each pubnub channel as map using the pubnub channel name as the key.
// NewSubscribedChannels keeps a list of the new subscribed Pubnub channels by the user in the a comma 
// separated string, before they are appended to the Pubnub SubscribedChannels.
type Pubnub struct {
    Origin                   string
    PublishKey               string
    SubscribeKey             string
    SecretKey                string
    CipherKey                string
    Ssl                      bool
    Uuid                     string
    SubscribedChannels       string 
    TimeToken                string
    SentTimeToken            string
    ResetTimeToken           bool  
    PresenceChannels         map[string] chan []byte
    SubscribeChannels        map[string] chan []byte
    PresenceErrorChannels    map[string] chan []byte
    SubscribeErrorChannels   map[string] chan []byte
    NewSubscribedChannels    string
}

// VersionInfo returns the version of the this code along with the build date. 
func VersionInfo() string{
    return "Version: 3.4.1; Build Date: Sep 9, 2013;"
}

// PubnubInit initializes pubnub struct with the user provided values.
// And then initiates the origin by appending the protocol based upon the sslOn argument.
// Then it uses the customuuid or generates the uuid.
// 
// It accepts the following parameters:
// publishKey is the user specific Publish Key. Mandatory.
// subscribeKey is the user specific Subscribe Key. Mandatory.
// secretKey is the user specific Secret Key. Accepts empty string if not used.
// cipherKey stores the user specific Cipher Key. Accepts empty string if not used. 
// sslOn is true if enabled, else is false.  
// customUuid is the unique identifier, it can be a custom value or sent as empty for automatic generation. 
//
// returns the pointer to Pubnub instance.
func PubnubInit(publishKey string, subscribeKey string, secretKey string, cipherKey string, sslOn bool, customUuid string) *Pubnub {
    newPubnub := &Pubnub{
        Origin:                _origin,
        PublishKey:            publishKey,
        SubscribeKey:          subscribeKey,
        SecretKey:             secretKey,
        CipherKey:             cipherKey,
        Ssl:                   sslOn,
        Uuid:                  "",
        SubscribedChannels:    "",
        ResetTimeToken:        true,
        TimeToken:             "0",
        SentTimeToken:             "0",
        NewSubscribedChannels: "",
        PresenceChannels:       make(map[string] chan []byte),
        SubscribeChannels:       make(map[string] chan []byte),
        PresenceErrorChannels:       make(map[string] chan []byte),
        SubscribeErrorChannels:       make(map[string] chan []byte),
    }

    if newPubnub.Ssl {
        newPubnub.Origin = "https://" + newPubnub.Origin
    } else {
        newPubnub.Origin = "http://" + newPubnub.Origin
    }
    
    //Generate the uuid is custmUuid is not provided
    if strings.TrimSpace(customUuid) == "" {
        uuid, err := GenUuid()
        if err == nil {
            newPubnub.Uuid = uuid
        } else {
            fmt.Println(err)
        }
    } else {
        newPubnub.Uuid = customUuid
     }
    return newPubnub
}

// SetProxy sets the global variables for the parameters.
// It also sets the _proxyServerEnabled value to true.
// 
// It accepts the following parameters:
// proxyServer proxy server name or ip.
// proxyPort proxy port.
// proxyUser proxyUserName.
// proxyPassword proxyPassword.
func SetProxy(proxyServer string, proxyPort int, proxyUser string, proxyPassword string){
    _proxyServer = proxyServer
    _proxyPort = proxyPort
    _proxyUser = proxyUser
    _proxyPassword = proxyPassword
    _proxyServerEnabled = true
}

// SetResumeOnReconnect sets the value of _resumeOnReconnect.
func SetResumeOnReconnect(val bool){
    _resumeOnReconnect = val
}

// SetSubscribeTimeout sets the value of _subscribeTimeout.
func SetSubscribeTimeout(val int64){
    _subscribeTimeout = val
}

// SetOrigin sets the value of _origin. Should be called before PubnubInit
func SetOrigin(val string){
    _origin = val
}

// Abort is the struct Pubnub's instance method that closes the open connections for both subscribe 
// and non-subscribe requests.
//
// It also sends a leave request for all the subscribed channel and
// sets the pub.SubscribedChannels as empty to break the loop in the func StartSubscribeLoop  
func (pub *Pubnub) Abort() {
    if(pub.SubscribedChannels != ""){
        value, _, err := pub.SendLeaveRequest(pub.SubscribedChannels)
        if err != nil {
            pub.SendResponseToChannel(nil, pub.SubscribedChannels, 9, err.Error(), "")
        }else{
            pub.SendResponseToChannel(nil, pub.SubscribedChannels, 5, string(value), "")
        }            
        
        pub.SubscribedChannels = ""
    }
    
    if(_conn != nil) {
        _conn.Close()
    }
    if(_subscribeConn!= nil) {
        _subscribeConn.Close()
    }
}

// GetTime is the struct Pubnub's instance method that calls the ExecuteTime
// method to process the time request.
//. 
// It accepts the following parameters:
// callbackChannel on which to send the response.
// errorChannel on which to send the error response. 
func (pub *Pubnub) GetTime(callbackChannel chan []byte, errorChannel chan []byte) {
    pub.ExecuteTime(callbackChannel, errorChannel, 0)
}

// ExecuteTime  is the struct Pubnub's instance method that creates a time request and sends back the 
// response to the channel.
// Closes the channel when the response is sent.
// In case we get an invalid json response the routine retries till the _maxRetries to get a valid
// response. 
//
// callbackChannel on which to send the response.
// errorChannel on which the error response is sent.
// retryCount to track the retry logic. 
func (pub *Pubnub) ExecuteTime(callbackChannel chan []byte, errorChannel chan []byte, retryCount int) {
    count := retryCount
     
    timeUrl := ""
    timeUrl += "/time"
    timeUrl += "/0"

    value, _, err := pub.HttpRequest(timeUrl, false)

    if err != nil {        
        pub.SendResponseToChannel(errorChannel, "", 10, err.Error(), "")
    } else {
        _, _, _, errJson := ParseJson(value, pub.CipherKey)
        if(errJson != nil && strings.Contains(errJson.Error(), _invalidJson)){
            pub.SendResponseToChannel(errorChannel, "", 10, errJson.Error(), "") 
            if (count<_maxRetries) {
                count++
                pub.ExecuteTime(callbackChannel, errorChannel, count)    
            }            
        } else {
            callbackChannel <- []byte(fmt.Sprintf("%s", value))
        }
    }
}

// SendPublishRequest is the struct Pubnub's instance method that posts a publish request and 
// sends back the response to the channel.
//
// It accepts the following parameters:
// channel: pubnub channel to publish to
// publishUrlString: The url to which the message is to be appended.
// jsonBytes: the message to be sent.
// callbackChannel: Channel on which to send the response.
// errorChannel on which the error response is sent.
func (pub *Pubnub) SendPublishRequest(channel string, publishUrlString string, jsonBytes []byte, callbackChannel chan []byte, errorChannel chan []byte) {
    var publishUrl *url.URL
    publishUrl, urlErr := url.Parse(publishUrlString)
    if urlErr != nil {
        errorChannel <- []byte(fmt.Sprintf("%s", urlErr))
    } else {
        publishUrl.Path += string(jsonBytes)
        value, responseCode, err := pub.HttpRequest(publishUrl.String(), false)

        if ((responseCode != 200) || (err != nil)) {
            if ((value != nil) && (responseCode > 0)) {
                var s []interface{}
                errJson := json.Unmarshal(value, &s)
                
                if ((errJson==nil) && (len(s) >0)){
                    if message, ok := s[1].(string); ok { 
                        pub.SendResponseToChannel(errorChannel, channel, 9, message, strconv.Itoa(responseCode))
                    } else {
                        pub.SendResponseToChannel(errorChannel, channel, 9, string(value), strconv.Itoa(responseCode))
                    }  
                } else {                
                    pub.SendResponseToChannel(errorChannel, channel, 9, string(value), strconv.Itoa(responseCode))
                }    
            } else if ((err != nil) && (responseCode > 0))  {
                pub.SendResponseToChannel(errorChannel, channel, 9, err.Error(),  strconv.Itoa(responseCode))
            } else if (err != nil) {
                pub.SendResponseToChannel(errorChannel, channel, 9, err.Error(), "")    
            } else {
                pub.SendResponseToChannel(errorChannel, channel, 9, _publishFailed, strconv.Itoa(responseCode))
            }    
        } else {
            _, _, _, errJson := ParseJson(value, pub.CipherKey)
            if(errJson != nil && strings.Contains(errJson.Error(), _invalidJson)){
                pub.SendResponseToChannel(errorChannel, channel, 9, errJson.Error(), "") 
            } else {
                callbackChannel <- []byte(fmt.Sprintf("%s", value))
            }
        }    
    }
}

// InvalidMessage takes the message in form of a interface and checks if the message is nil or empty.
// Returns true if the message is nil or empty.
// Returns false is the message is acceptable.
func InvalidMessage(message interface{}) bool{
    if(message == nil){
        return true
    }
    
    dataInterface := message.(interface{})
    
    switch vv := dataInterface.(type){
        case string:
            if (strings.TrimSpace(vv) != ""){
                return false
            }
        case []interface{}:
            if (vv != nil) {
                return false
            }
        default :
            if (vv != nil) {
                return false
            }
        }     
    return true    
}

// InvalidChannel takes the Pubnub channel and the channel as parameters. 
// Multiple Pubnub channels are accepted separated by comma.
// It splits the Pubnub channel string by a comma and checks if the channel empty.
// Returns true if any one of the channel is empty. And sends a response on the Pubnub channel stating 
// that there is an "Invalid Channel".
// Returns false if all the channels is acceptable.
func InvalidChannel(channel string, c chan []byte) bool{
    if (strings.TrimSpace(channel) == "") {
        return true
    } else {
        channelArray := strings.Split(channel, ",")
    
        for i:=0; i < len(channelArray); i++ {
            if (strings.TrimSpace(channelArray[i]) == "") {    
                c <- []byte(fmt.Sprintf("Invalid Channel: %s", channel))
                return true
            }
        }    
    }
    return false
}

// Publish is the struct Pubnub's instance method that creates a publish request and calls 
// SendPublishRequest to post the request. 
//
// It calls the InvalidChannel and InvalidMessage methods to validate the Pubnub channels and message.
// Calls the GetHmacSha256 to generate a signature if a secretKey is to be used.
// Creates the publish url
// Calls json marshal
// Calls the EncryptString method is the cipherkey is used and calls json marshal
// Closes the channel after the response is received
//
// It accepts the following parameters:
// channel: The Pubnub channel to which the message is to be posted.
// message: message to be posted.
// callbackChannel: Channel on which to send the response back.
// errorChannel on which the error response is sent.
func (pub *Pubnub) Publish(channel string, message interface{}, callbackChannel chan []byte, errorChannel chan []byte) {
    if(pub.PublishKey == ""){
        pub.SendResponseToChannel(errorChannel, channel, 9, "Publish key required.", "")
        return
    } 

    if(InvalidChannel(channel, callbackChannel)){
        return 
    }

    if(InvalidMessage(message)){
        pub.SendResponseToChannel(errorChannel, channel, 9, "Invalid Message.", "")
        return 
    }

    signature := ""
    if pub.SecretKey != "" {
        signature = GetHmacSha256(pub.SecretKey, fmt.Sprintf("%s/%s/%s/%s/%s", pub.PublishKey, pub.SubscribeKey, pub.SecretKey, channel, message))
    } else {
        signature = "0"
    }
    var publishUrlBuffer bytes.Buffer
    publishUrlBuffer.WriteString("/publish")
    publishUrlBuffer.WriteString("/")
    publishUrlBuffer.WriteString(pub.PublishKey)
    publishUrlBuffer.WriteString("/")
    publishUrlBuffer.WriteString(pub.SubscribeKey)
    publishUrlBuffer.WriteString("/")
    publishUrlBuffer.WriteString(signature)
    publishUrlBuffer.WriteString("/")
    publishUrlBuffer.WriteString(channel)
    publishUrlBuffer.WriteString("/0/")

    jsonSerialized, err := json.Marshal(message)
    if err != nil {
        pub.SendResponseToChannel(errorChannel, channel, 9, fmt.Sprintf("error in serializing: %s", err), "")
    } else {
        if pub.CipherKey != "" {
            //Encrypt and Serialize
            jsonEncBytes, errEnc := json.Marshal(EncryptString(pub.CipherKey, fmt.Sprintf("%s", jsonSerialized)))
            if errEnc != nil {
                pub.SendResponseToChannel(errorChannel, channel, 9, fmt.Sprintf("error in serializing: %s", errEnc), "")
            } else {
                pub.SendPublishRequest(channel, publishUrlBuffer.String(), jsonEncBytes, callbackChannel, errorChannel)
            }
        } else {
            pub.SendPublishRequest(channel, publishUrlBuffer.String(), jsonSerialized, callbackChannel, errorChannel)
        }
    }
}

// SendResponseToChannel is the struct Pubnub's instance method that sends a reponse on the channel 
// provided as an argument or to the subscribe / presence channel is the argument is nil. 
//
// Constructs the response based on the action (1-9). In case the action is 5 sends the response 
// as in the parameter response. 
//
// It accepts the following parameters:
// c: Channel on which to send the response back. Can be nil. If nil, assumes that if the channel name 
// is suffixed with "-pnpres" it is a presence channel else subscribe channel and sends the response to all the 
// respective channels. Then it fetches the corresonding channel from the pub.PresenceChannels or pub.SubscribeChannels
// in case of callback and pub.PresenceErrorChannels or pub.SubscribeErrorChannels in case of error 
//
// channels: Pubnub Channels to send a response to. Comma separated string for multiple channels.
// action: (1-9) 
// response: can be nil, is used only in the case action is '5'. 
// response2: Additional error info.
func (pub *Pubnub) SendResponseToChannel(c chan []byte, channels string, action int, response string, response2 string){
    message := ""
    intResponse := "0"
    sendReponseAsIs := false
    sendErrorResponse := false
    errorWithoutChannel := false
    switch action {
        case 1:
            message = "already subscribed"
            intResponse = "0"          
        case 2:
            message = "connected"
            intResponse = "1"
        case 3:
            message = "unsubscribed"
            intResponse = "1"
        case 4:
            message = "not subscribed"
            intResponse = "0"
        case 5:
            sendReponseAsIs = true
        case 6:
            message = "reconnected"
            intResponse = "1"
        case 7:
            message = "disconnected due to internet connection issues, trying to reconnect. Retry count:" + response
            response = ""
            intResponse = "0"
            sendErrorResponse = true
        case 8:
            message = "aborted due to max retry limit"
            intResponse = "0"
            sendErrorResponse = true
        case 9:
            sendErrorResponse = true  
            sendReponseAsIs = true
            intResponse = "0" 
        case 10:
            errorWithoutChannel = true
        case 11:
            message = "timed out."
            response = ""
            intResponse = "0"
            sendErrorResponse = true
    }
    
    var value string
    channelArray := strings.Split(channels, ",")
    
    for i := 0; i < len(channelArray); i++ {
        responseChannel := c
        presence := "Subscription to channel "
        channel := strings.TrimSpace(channelArray[i])
        if(channel == ""){
            continue
        }
        
        if(response == ""){
            response = message
        }
            
        if (sendErrorResponse){
            isPresence := false
            if(responseChannel == nil){
                responseChannel, isPresence = pub.GetChannelForPubnubChannel(channel, true)
            }        
            if(isPresence){
                presence = "Presence notifications for channel "
            }
            if(sendReponseAsIs){
                presence = ""
            }
            if ((response2 != "") && (response2 != "0")){
                value = fmt.Sprintf("[%s, \"%s%s\", %s, \"%s\"]", intResponse, presence, response, response2, strings.Replace(channel, _presenceSuffix, "", -1))
            } else {
                value = fmt.Sprintf("[%s, \"%s%s\", \"%s\"]", intResponse, presence, response, strings.Replace(channel, _presenceSuffix, "", -1))
            }

            if responseChannel!= nil {
                responseChannel <- []byte(value)    
            }            
        } else {
            isPresence := false
            if(responseChannel == nil){
                responseChannel, isPresence = pub.GetChannelForPubnubChannel(channel, false)
            }    
            if(isPresence) {
                channel = strings.Replace(channel, _presenceSuffix, "", -1)
                presence = "Presence notifications for channel "
            }
            
            if(sendReponseAsIs){
                value = strings.Replace(response, _presenceSuffix, "", -1)
            } else {
                value = fmt.Sprintf("[%s, \"%s'%s' %s\", \"%s\"]", intResponse, presence, channel, message, channel)
            }
            if responseChannel!= nil {
                responseChannel <- []byte(value)    
            }            
        }
    }  
    if(errorWithoutChannel){
        responseChannel := c
        value = fmt.Sprintf("[%s, \"%s\"]", intResponse, response)
        if responseChannel!= nil {
            responseChannel <- []byte(value)    
        }            
    }      
}

// GetChannelForPubnubChannel parses the pubnub channel and returns the the callback or the erro channel
// 
// Accepts the pubnub channel name channel as string, the string is parsed to check if it is a
// Subscribe or a Presence channel.
// 
// and isErrorChannel as bool. If it is true PresenceErrorChannels or SubscribeErrorChannels
// will be used to fetch the corresponding channel.
//
// Returns channel to send a response on 
// and bool if true means it is a pubnub Presence channel. Else it is a pubnub Subscribe channel. 
func (pub *Pubnub) GetChannelForPubnubChannel(channel string, isErrorChannel bool) (chan []byte, bool) {
    isPresence := strings.Contains(channel, _presenceSuffix)
    if(isPresence){
        channel = strings.Replace(string(channel), _presenceSuffix, "", -1)
        if(isErrorChannel){
            c, found := pub.PresenceErrorChannels[channel]
            if(found){
                return c, true
            }
        } else {
            c, found := pub.PresenceChannels[channel]
            if(found){
                return c, true
            }
        }
    }else {
        if(isErrorChannel){
            c, found := pub.SubscribeErrorChannels[channel]
            if(found){
                return c, false
            }
        } else {    
            c, found := pub.SubscribeChannels[channel]
            if(found){
                return c, false
            }
        }            
    }
    return nil, false    
}

// GetSubscribedChannels is the struct Pubnub's instance method that iterates through the Pubnub 
// SubscribedChannels and appends the new channels.  
//
// It splits the Pubnub channels in the parameter by a comma and compares them to the existing 
// subscribed Pubnub channels. 
// If a new Pubnub channels is found it is appended to the Pubnub SubscribedChannels. The return 
// parameter channelsModified is set to true 
// If an subscribed pubnub channel is already present in the Pubnub SubscribedChannels it is added to 
// the alreadySubscribedChannels string and a response is sent back to the channel  
//
// It accepts the following parameters:
// channels: Pubnub Channels to send a response to. Comma separated string for multiple channels.
// c: Channel on which to send the response back. Can be nil. If nil assumes that if the channel name 
// is suffixed with "-pnpres" it is a presence channel else subscribe channel and send the response to 
// the respective channel.
// isPresenceSubscribe: can be nil, is used only in the case action is '5'.
// errorChannel: channel to send the error response to.
//
// Returns:
// subChannels: the Pubnub subscribed channels as a comma separated string.  
// newSubChannels: the new Pubnub subscribed channels as a comma separated string.
// b: The return parameter channelsModified is set to true if new channels are added.
func (pub *Pubnub) GetSubscribedChannels(channels string, callbackChannel chan []byte, isPresenceSubscribe bool, errorChannel chan []byte) (subChannels string, newSubChannels string, b bool) {
    channelArray := strings.Split(channels, ",")
    subscribedChannels := pub.SubscribedChannels
    newSubscribedChannels := ""
    channelsModified := false
    alreadySubscribedChannels := ""
        
    for i := 0; i < len(channelArray); i++ {
        channelToSub := strings.TrimSpace(channelArray[i])
        if(isPresenceSubscribe){
            channelToSub += _presenceSuffix
        } 
        
        if pub.NotDuplicate(channelToSub) {
            if len(subscribedChannels)>0 {
                subscribedChannels += ","
            }
            subscribedChannels += channelToSub
                 
            if len(newSubscribedChannels)>0 {
                newSubscribedChannels += ","
            }     
            newSubscribedChannels += channelToSub
            channelsModified = true
        }else{
            if len(alreadySubscribedChannels)>0 {
                alreadySubscribedChannels += ","
            }            
            alreadySubscribedChannels += channelToSub
        }
    }
    
    if len(alreadySubscribedChannels)>0 {
        pub.SendResponseToChannel(errorChannel, alreadySubscribedChannels, 1, "", "")
    }    

    return subscribedChannels, newSubscribedChannels, channelsModified
}

// CheckForTimeoutAndRetries parses the error in case of subscribe error response. Its an Pubnub instance method.
// If any of the strings "Error in initializating connection", "timeout", "no such host" 
// are found it assumes that a network connection is lost.
// Sends a response to the subscribe/presence channel.
//
// If max retries limit is reached it empties the Pubnub SubscribedChannels thus initiating 
// the subscribe/presence subscription closure.
//
// It accepts the following parameters:
// err: error object
// errChannel: channel to send a response to.
//
// Returns:
// b: Bool variable true incase the connection is lost.
// bTimeOut: bool variable true in case Timeout condition is met.
func (pub *Pubnub) CheckForTimeoutAndRetries(err error, errChannel chan []byte) (bool, bool){
    bRet := false
    bTimeOut := false
    errorInitConn :=strings.Contains(err.Error(), _errorInInitializing)
    if (errorInitConn){
        SleepForAWhile(true)
        pub.SendResponseToChannel(nil, pub.SubscribedChannels, 9, err.Error(), "Retry count: " + strconv.Itoa(_retryCount))
        bRet = true
    }else if  (strings.Contains(err.Error(), _timeoutU)) {
        SleepForAWhile(false)
        pub.SendResponseToChannel(nil, pub.SubscribedChannels, 11, strconv.Itoa(_retryCount), "")        
        bRet = true
        bTimeOut = true
    }else if (strings.Contains(err.Error(), _noSuchHost) || strings.Contains(err.Error(), _networkUnavailable)){
        SleepForAWhile(true)
        pub.SendResponseToChannel(nil, pub.SubscribedChannels, 7, strconv.Itoa(_retryCount), "")
        bRet = true
    }
    
    if(_retryCount >= _maxRetries){
        pub.SendResponseToChannel(nil, pub.SubscribedChannels, 8, "", "")
        pub.SubscribedChannels = ""
        _retryCount = 0
    }
    
    if(_retryCount > 0){    
        return bRet, bTimeOut
    }
    return bRet, bTimeOut
}

// StartSubscribeLoop starts a continuous loop that handles the reponse from pubnub 
// subscribe/presence subscriptions.
//
// It creates subscribe request url and posts it. 
// When the response is received it: 
// Checks for errors and timeouts, closes the existing connections and continues the loop if true.
// else parses the response. stores the time token if it is a timeout from server.
  
// Checks For Timeout And Retries: 
// If sent timetoken is 0 and the data is empty the connected response is sent back to the channel.
// If no error is received the response is sent to the presence or subscribe pubnub channels. 
// if the channel name is suffixed with "-pnpres" it is a presence channel else subscribe channel 
// and send the response the the respective channel.
//
// It accepts the following parameters:
// channels: channels to subscribe.
// errorChannel: Channel to send the error response to.
//
// TODO: Refactor
func (pub *Pubnub) StartSubscribeLoop(channels string, errorChannel chan []byte){
    channelCount := len(pub.SubscribedChannels)
    channelsModified := false
    for {
         if len(pub.SubscribedChannels) > 0 {
            if(len(pub.SubscribedChannels) != channelCount){
                channelsModified = true
            }
            sentTimeToken := pub.TimeToken
            subscribeUrl, sentTimeToken := pub.CreateSubscribeUrl(sentTimeToken)
            value, responseCode, err := pub.HttpRequest(subscribeUrl, true)
            
            if ((responseCode != 200) || (err != nil)) {
                
                if(err!=nil){
                    bNonTimeout, bTimeOut := pub.CheckForTimeoutAndRetries(err, errorChannel)
                    if (strings.Contains(err.Error(), _connectionAborted)){
                        pub.CloseExistingConnection()	
                        pub.SendResponseToChannel(nil, pub.SubscribedChannels, 9, err.Error(), strconv.Itoa(responseCode))
                    } else if(bNonTimeout){
                        pub.CloseExistingConnection()
                        if(bTimeOut){
                            _, returnTimeToken, _, errJson := ParseJson(value, pub.CipherKey)
                            if(errJson == nil){
                               pub.TimeToken = returnTimeToken
                            }
                        }
                        if (!_resumeOnReconnect) {
                            pub.ResetTimeToken = true
                        }
                    } else {
                        pub.CloseExistingConnection()
                        pub.SendResponseToChannel(nil, pub.SubscribedChannels, 9, err.Error(), strconv.Itoa(responseCode))
                        SleepForAWhile(true)
                    }
                }
                continue
            } else if string(value) != "" {                
                if string(value) == "[]" {
                    SleepForAWhile(false)
                    continue
                }      
                        
                data, returnTimeToken, channelName, errJson := ParseJson(value, pub.CipherKey)
                pub.TimeToken = returnTimeToken
                if (data == "[]") {
                    if(!channelsModified){
                        
                        channelsModified = false
                    }
                    if(sentTimeToken == "0"){
                        pub.SendResponseToChannel(nil, pub.SubscribedChannels, 2, "", "")
                        pub.NewSubscribedChannels = ""
                    }
                    _retryCount = 0
                    continue
                }            
                pub.ParseHttpResponse(value, data, channelName, returnTimeToken, errJson, errorChannel)
            } 
        } else {
            break
        }    
    }    
}

// CreateSubscribeUrl creates a subscribe url to send to the origin
// If the resetTimeToken flag is true 
// it sends 0 to init the subscription. 
// Else sends the last timetoken.
//
// Accepts the sentTimeToken as a string parameter.
// retunrs the Url and the senttimetoken based on the logic above .
func (pub *Pubnub) CreateSubscribeUrl(sentTimeToken string) (string, string){
    var subscribeUrlBuffer bytes.Buffer
    subscribeUrlBuffer.WriteString("/subscribe")
    subscribeUrlBuffer.WriteString("/")
    subscribeUrlBuffer.WriteString(pub.SubscribeKey)
    subscribeUrlBuffer.WriteString("/")
    subscribeUrlBuffer.WriteString(pub.SubscribedChannels)
    subscribeUrlBuffer.WriteString("/0")
                
    if (pub.ResetTimeToken) {
        subscribeUrlBuffer.WriteString("/0")
        sentTimeToken = "0"   
        pub.SentTimeToken = "0"     
        pub.ResetTimeToken = false
    }else{
        subscribeUrlBuffer.WriteString("/")
        if(strings.TrimSpace(pub.TimeToken) == ""){
            pub.TimeToken = "0"
            pub.SentTimeToken = "0"
        } else {
            pub.SentTimeToken = sentTimeToken
        }    
        subscribeUrlBuffer.WriteString(pub.TimeToken)
    }
                
    if pub.Uuid != "" {
        subscribeUrlBuffer.WriteString("?uuid=")
        subscribeUrlBuffer.WriteString(pub.Uuid)
    }
    return subscribeUrlBuffer.String(), sentTimeToken
}

// ParseHttpResponse parses the http response from the orgin for the subscribe resquest 
// if errJson is not nil it sends an error response on the error channel.
// In case of subscribe response it parses the returned data and splits if multiple messages are received.
// 
// Accespts the following parameters
// value: is the actual response.  
// data: is the json deserialized string, 
// channelName: the pubnub channel of the response 
// returnTimeToken: return time token from the origin, 
// errJson: error if received from server, can be nil. 
// errorChannel: channel to send an error response to.
func (pub *Pubnub) ParseHttpResponse(value []byte, data string, channelName string, returnTimeToken string, errJson error, errorChannel chan []byte){
    if errJson != nil {
        pub.SendResponseToChannel(nil, channelName, 9, fmt.Sprintf("%s", errJson), "")
        SleepForAWhile(false)
    } else {
        _retryCount = 0
        if (channelName == ""){                        
            channelName = pub.SubscribedChannels
        }
        pub.SplitMessagesAndSendJsonResponse(data, returnTimeToken, channelName, errorChannel)
    }                 
}

// SplitMessagesAndSendJsonResponse unmarshals the data and sends a response if the 
// data type is a non array. Else calls the CreateAndSendJsonResponse to split the messages.
// 
// parameters:
// data: the data to parse and split, 
// returnTimeToken: the return timetoken in the response
// channels: pubnub channels in the response.
func (pub *Pubnub) SplitMessagesAndSendJsonResponse (data string, returnTimeToken string, channels string, errorChannel chan []byte) {
    channelSlice := strings.Split(channels, ",")
    channelLen := len(channelSlice)
    isPresence := false
    if(channelLen == 1){
        isPresence = strings.Contains(channels, _presenceSuffix)
    }
    
    if((channelLen == 1) && (isPresence)){
        pub.SplitPresenceMessages([]byte(data), returnTimeToken, channelSlice[0], errorChannel)
    } else if((channelLen == 1) && (!isPresence)) {
        pub.SplitSubscribeMessages(data, returnTimeToken, channelSlice[0], errorChannel)
    } else {
        var returnedMessages interface{}
        errUnmarshalMessages := json.Unmarshal([]byte(data), &returnedMessages)
        
        if errUnmarshalMessages == nil {
            v := returnedMessages.(interface{})
            
            switch vv := v.(type) {
                case string:
                   length := len(vv)
                   if(length > 0){
                          pub.SendJsonResponse(vv, returnTimeToken, channels)
                   }
                case []interface{}:
                      pub.CreateAndSendJsonResponse(vv, returnTimeToken, channels)
            }
        }
    }
}

// SplitPresenceMessages splits the multiple messages 
// unmarshals the data into the custom structure, 
// calls the SendJsonResponse funstion to creates the json again.
//
// Parameters:
// data: data to unmarshal,
// returnTimeToken: the returned timetoken in the pubnub response, 
// channel: pubnub channel,
// errorChannel: error channel to send a error response back.
func (pub *Pubnub) SplitPresenceMessages(data []byte, returnTimeToken string, channel string, errorChannel chan []byte){
    var occupants []struct {
        Action string `json:"action"`
        Uuid string `json:"uuid"`
        Timestamp float64 `json:"timestamp"`
        Occupancy int `json:"occupancy"`
    }
    errUnmarshalMessages := json.Unmarshal(data, &occupants)
    if(errUnmarshalMessages !=nil){    
        pub.SendResponseToChannel(nil, channel, 9, _invalidJson, "")
    } else {
        for i := range occupants {
            intf := make([]interface{}, 1)
            intf[0] = occupants[i]
            pub.SendJsonResponse(intf, returnTimeToken, channel)
        }        
    }    
}
    
// SplitSubscribeMessages splits the multiple messages 
// unmarshals the data into the custom structure, 
// calls the SendJsonResponse funstion to creates the json again.
//
// Parameters:
// data: data to unmarshal,
// returnTimeToken: the returned timetoken in the pubnub response, 
// channel: pubnub channel,
// errorChannel: error channel to send a error response back.
func (pub *Pubnub) SplitSubscribeMessages(data string, returnTimeToken string, channel string, errorChannel chan []byte){
    var occupants []interface {}
    errUnmarshalMessages := json.Unmarshal([]byte(data), &occupants)
    if(errUnmarshalMessages !=nil){    
        pub.SendResponseToChannel(nil, channel, 9, _invalidJson, "")
    } else {
        for i := range occupants {
            intf := make([]interface{}, 1)
            intf[0] = occupants[i]
            pub.SendJsonResponse(intf, returnTimeToken, channel)           
        }        
    }    
}

// CreateAndSendJsonResponse marshals the data for each split message and calls
// the SendJsonResponse multiple times to send response back to the channel
//  
// Accepts:
// rawData: the data to parse and split, 
// returnTimeToken: the return timetoken in the response
// channels: pubnub channels in the response.
func (pub *Pubnub) CreateAndSendJsonResponse(rawData interface{}, returnTimeToken string, channels string){
    channelSlice := strings.Split(channels, ",")
    dataInterface := rawData.(interface{})
    switch vv := dataInterface.(type){
        case []interface{}:
            for i, u := range vv {
                intf := make([]interface{}, 1)
                if (reflect.TypeOf(u).Kind() == reflect.String){
                    intf[0] = u
                } else {
                    intf[0] = vv[i]
                }
                channel := ""
                
                if(i <= len(channelSlice)-1){
                    channel = channelSlice[i]
                } else {
                    channel = channelSlice[0]
                } 
                
                pub.SendJsonResponse(intf, returnTimeToken, channel)
            }
    } 
}

// SendJsonResponse creates a json response and sends back to the response channel
// 
// Accepts:
// message: response to send back, 
// returnTimeToken: the timetoken for the response, 
// channelName: the pubnub channel for the response.
func (pub *Pubnub) SendJsonResponse (message interface{}, returnTimeToken string, channelName string){
                        
    if(channelName != "") {
        response := []interface{} {message, fmt.Sprintf("%s",pub.TimeToken), channelName}
        jsonData, err := json.Marshal(response)
        if (err != nil) {
            pub.SendResponseToChannel(nil, channelName, 9, _invalidJson, err.Error())
        }
        pub.SendResponseToChannel(nil, channelName, 5, string(jsonData), "")    
    }
}

// GetSubscribedChannelName is the struct Pubnub's instance method. 
// In case of single subscribe request the channelname will be empty.
// This methos iterates through the pubnub SubscribedChannels to find the name of the channel.
func (pub *Pubnub) GetSubscribedChannelName() (string){
    channelArray := strings.Split(pub.SubscribedChannels, ",")
    for i := 0; i < len(channelArray); i++ {
        if (strings.Contains(channelArray[i], _presenceSuffix)) {
            continue
        }else{
            return channelArray[i]
        }            
    }
    return ""
}

// CloseExistingConnection: Closes the open subscribe/presence connection.
func (pub *Pubnub) CloseExistingConnection(){
    if(_subscribeConn != nil){
        _subscribeConn.Close()
    }    
}

// Subscribe is the struct Pubnub's instance method which checks for the InvalidChannels 
// and returns if true.
// Initaiates the presence and subscribe response channels. 
// It creates a map for callback and error response channels for 
// each pubnub channel using the pubnub channel name as the key.
// If muliple channels are passed then the same callback or error channel is used.
//
// If there is no existing subscribe/presence loop running then it starts a 
// new loop with the new pubnub channels.
// Else closes the existing connections and starts a new loop
//
// It accepts the following parameters:
// channels: comma separated pubnub channel list.
// timetoken: if timetoken is present the subscribe request is sent using this timetoken 
// callbackChannel: Channel on which to send the response back.
// isPresenceSubscribe: tells the method that presence subscription is requested.
// errorChannel: channel to send an error response to.
func (pub *Pubnub) Subscribe(channels string, timetoken string, callbackChannel chan []byte, isPresenceSubscribe bool, errorChannel chan []byte) {
    if(InvalidChannel(channels, callbackChannel)){
        return 
    }
    
    subscribedChannels, newSubscribedChannels, channelsModified := pub.GetSubscribedChannels(channels, callbackChannel, isPresenceSubscribe, errorChannel)

    var channelArr = strings.Split(channels, ",")
    
    for i, u := range channelArr {
        if isPresenceSubscribe {
            pub.PresenceChannels[u] = callbackChannel
            pub.PresenceErrorChannels[u] = errorChannel
        } else {
            pub.SubscribeChannels[u] = callbackChannel
            pub.SubscribeErrorChannels[u] = errorChannel
        }
        i++
    }
    pub.NewSubscribedChannels = newSubscribedChannels
    if(pub.SubscribedChannels == ""){
        if(strings.TrimSpace(timetoken) != ""){
            pub.TimeToken = timetoken
            pub.ResetTimeToken = false
        } else {
            pub.ResetTimeToken = true
        }
        pub.SubscribedChannels = subscribedChannels
        go pub.StartSubscribeLoop(channels, errorChannel)
    }else if (channelsModified){  
        pub.CloseExistingConnection()
        if(strings.TrimSpace(timetoken) != ""){
            pub.TimeToken = timetoken
            pub.ResetTimeToken = false
        } else {
            pub.ResetTimeToken = true
        }
        pub.SubscribedChannels = subscribedChannels
    }
}    

// SleepForAWhile pauses the subscribe/presence loop for the _retryInterval. 
func SleepForAWhile(retry bool){
    if(retry) {
        _retryCount++
    }
    time.Sleep(_retryInterval * time.Second)
}

// NotDuplicate is the struct Pubnub's instance method which checks for the channel name 
// to check in the existing pubnub SubscribedChannels.
// 
// It accepts the following parameters:
// channel: the Pubnub channel name to check in the existing pubnub SubscribedChannels.
//
// returns:
// true if the channel is found.
// false if not found.
func (pub *Pubnub) NotDuplicate(channel string) (b bool){
    var channels = strings.Split(pub.SubscribedChannels, ",")
    for i, u := range channels {
        if channel == u {
            return false
        } 
        i++
    }
    return true 
}

// RemoveFromSubscribeList is the struct Pubnub's instance method which checks for the 
// channel name to check in the existing pubnub SubscribedChannels and removes it if found 
// 
// It accepts the following parameters:
// c: Channel on which to send the response back.
// channel: the pubnub channel name to check in the existing pubnub SubscribedChannels.
//
// returns:
// true if the channel is found and removed.
// false if not found.
func (pub *Pubnub) RemoveFromSubscribeList(c chan []byte, channel string) (b bool){
    var channels = strings.Split(pub.SubscribedChannels, ",")
    newChannels := ""
    found := false
    for _, u := range channels {
        if channel == u {
            found = true
            pub.SendResponseToChannel(c, u, 3, "", "")
        } else {
            if len(newChannels)>0 {
                newChannels += ","
            }          
            newChannels += u            
        }
    }
    if found {
        pub.SubscribedChannels = newChannels
    }
    return found
}

// Unsubscribe is the struct Pubnub's instance method which unsubscribes a pubnub subscribe 
// channel(s) from the subscribe loop.
//
// If all the pubnub channels are not removed the method StartSubscribeLoop will take care 
// of it by starting a new loop. 
// Closes the channel c when the processing is complete 
// 
// It accepts the following parameters:
// channels: the pubnub channel(s) in a comma separated string.
// callbackChannel: Channel on which to send the response back.
// errorChannel: channel to send an error response to.
func (pub *Pubnub) Unsubscribe(channels string, callbackChannel chan []byte, errorChannel chan []byte) {
    channelArray := strings.Split(channels, ",")
    unsubscribeChannels := ""
    channelRemoved := false
    
    for i := 0; i < len(channelArray); i++ {
        if i>0 {
            unsubscribeChannels += ","
        }
        channelToUnsub := strings.TrimSpace(channelArray[i]);
        unsubscribeChannels += channelToUnsub
        removed := pub.RemoveFromSubscribeList(callbackChannel, channelToUnsub)
        if !removed {
            pub.SendResponseToChannel(callbackChannel, channelToUnsub, 4, "", "")
        } else {
            channelRemoved = true
        }
    }
    
    if(channelRemoved) {
        pub.CloseExistingConnection()
        
        if (strings.TrimSpace(pub.SubscribedChannels) == "") {
            value, _, err := pub.SendLeaveRequest(channels)        
            if err != nil {
                pub.SendResponseToChannel(errorChannel, channels, 9, err.Error(), "")
            }else{
                pub.SendResponseToChannel(callbackChannel, channels, 5, string(value), "")
            }
        }
    }
}

// PresenceUnsubscribe is the struct Pubnub's instance method which unsubscribes a pubnub 
// presence channel(s) from the subscribe loop. 
//
// If all the pubnub channels are not removed the method StartSubscribeLoop will take care 
// of it by starting a new loop.
// When the pubnub channel(s) are removed it creates and posts a leave request. 
// 
// It accepts the following parameters:
// channels: the pubnub channel(s) in a comma separated string.
// callbackChannel: Channel on which to send the response back.
// errorChannel: channel to send an error response to.
func (pub *Pubnub) PresenceUnsubscribe(channels string, callbackChannel chan []byte, errorChannel chan []byte) {
    channelArray := strings.Split(channels, ",")
    presenceChannels := ""
    channelRemoved := false
    
    for i := 0; i < len(channelArray); i++ {
        if i>0 {
            presenceChannels += ","
        }
        channelToUnsub := strings.TrimSpace(channelArray[i]) + _presenceSuffix
        presenceChannels += channelToUnsub
        removed := pub.RemoveFromSubscribeList(callbackChannel, channelToUnsub) 
        if !removed {
            pub.SendResponseToChannel(errorChannel, channelToUnsub, 4, "", "")
        }else {
            channelRemoved = true
        }
    }
    
    if(channelRemoved) {
        pub.CloseExistingConnection() 
        if (strings.TrimSpace(pub.SubscribedChannels) == "") {
            value, _, err := pub.SendLeaveRequest(presenceChannels)        
            if err != nil {
                pub.SendResponseToChannel(errorChannel, channels, 9, err.Error(), "")
            }else{
                pub.SendResponseToChannel(callbackChannel, channels, 5, string(value), "")
            }            
        }    
    }    
}

// SendLeaveRequest: Sends a leave request to the origin
//
// It accepts the following parameters:
// channels: Channels to leave
//
// returns:
// the HttpRequest response contents as byte array.
// response error code,
// error if any.
func (pub *Pubnub) SendLeaveRequest(channels string) ([]byte, int, error){
    var subscribeUrlBuffer bytes.Buffer
    subscribeUrlBuffer.WriteString("/v2/presence")
    subscribeUrlBuffer.WriteString("/sub-key/")
    subscribeUrlBuffer.WriteString(pub.SubscribeKey)
    subscribeUrlBuffer.WriteString("/channel/")
    subscribeUrlBuffer.WriteString(channels)
    subscribeUrlBuffer.WriteString("/leave?uuid=")
    subscribeUrlBuffer.WriteString(pub.Uuid)
    
    return pub.HttpRequest(subscribeUrlBuffer.String(), false)
}

// History is the struct Pubnub's instance method which creates and post the History request 
// for a single pubnub channel.
//
// It parses the response to get the data and return it to the channel.
// 
// It accepts the following parameters:
// channel: a single value of the pubnub channel.
// limit: number of history messages to return.
// start: start time from where to begin the history messages.
// end: end time till where to get the history messages.
// reverse: to fetch the messages in ascending order
// callbackChannel on which to send the response.
// errorChannel on which the error response is sent.
// retryCount to track the retry logic. 
func (pub *Pubnub) History(channel string, limit int, start int64, end int64, reverse bool, callbackChannel chan []byte, errorChannel chan []byte) {
    pub.ExecuteHistory(channel, limit, start, end, reverse, callbackChannel, errorChannel, 0)
}

// ExecuteHistory is the struct Pubnub's instance method which creates and post the History request 
// for a single pubnub channel.
//
// It parses the response to get the data and return it to the channel.
// In case we get an invalid json response the routine retries till the _maxRetries to get a valid
// response. 
// 
// It accepts the following parameters:
// channel: a single value of the pubnub channel.
// limit: number of history messages to return.
// start: start time from where to begin the history messages.
// end: end time till where to get the history messages.
// reverse: to fetch the messages in ascending order
// callbackChannel on which to send the response.
// errorChannel on which the error response is sent.
// retryCount to track the retry logic.
func (pub *Pubnub) ExecuteHistory(channel string, limit int, start int64, end int64, reverse bool, callbackChannel chan []byte, errorChannel chan []byte, retryCount int) {
    count := retryCount
    if(InvalidChannel(channel, callbackChannel)){
        return 
    }

    if(limit < 0){
        limit = 100
    }
    
    var parameters bytes.Buffer
    parameters.WriteString("&reverse=")
    parameters.WriteString(fmt.Sprintf("%t", reverse))
    
    if(start > 0){
        parameters.WriteString("&start=")
        parameters.WriteString(fmt.Sprintf("%d", start))
    }
    if(end > 0){
        parameters.WriteString("&end=")
        parameters.WriteString(fmt.Sprintf("%d", end))
    }
    
    var historyUrlBuffer bytes.Buffer
    historyUrlBuffer.WriteString("/v2/history")
    historyUrlBuffer.WriteString("/sub-key/")
    historyUrlBuffer.WriteString(pub.SubscribeKey)
    historyUrlBuffer.WriteString("/channel/")
    historyUrlBuffer.WriteString(channel)
    historyUrlBuffer.WriteString("?count=")
    historyUrlBuffer.WriteString(fmt.Sprintf("%d", limit))
    historyUrlBuffer.WriteString(parameters.String())
        
    value, _, err := pub.HttpRequest(historyUrlBuffer.String(), false)

    if err != nil {
        pub.SendResponseToChannel(errorChannel, channel, 9, err.Error(), "")
    } else {
        data, returnOne, returnTwo, errJson := ParseJson(value, pub.CipherKey)
        if(errJson != nil && strings.Contains(errJson.Error(), _invalidJson)){
            pub.SendResponseToChannel(errorChannel, channel, 9, errJson.Error(), "") 
            if (count<_maxRetries) {
                count++
                pub.ExecuteHistory(channel, limit, start, end, reverse, callbackChannel, errorChannel, count)    
            }                  
        } else {
            var buffer bytes.Buffer
            buffer.WriteString("[")
            buffer.WriteString(data)
            buffer.WriteString(",\"" + returnOne + "\",\"" + returnTwo + "\"]")
               
            callbackChannel <- []byte(fmt.Sprintf("%s", buffer.Bytes()))
        }
    }
}

// HereNow is the struct Pubnub's instance method which creates and posts the herenow 
// request to get the connected users details.  
//
// It accepts the following parameters:
// channel: a single value of the pubnub channel. 
// callbackChannel on which to send the response.
// errorChannel on which the error response is sent.
func (pub *Pubnub) HereNow(channel string, callbackChannel chan []byte, errorChannel chan []byte) {
    pub.ExecuteHereNow(channel, callbackChannel, errorChannel, 0)
}

// ExecuteHereNow  is the struct Pubnub's instance method that creates a time request and sends back the 
// response to the channel.
// 
// In case we get an invalid json response the routine retries till the _maxRetries to get a valid
// response. 
//
// callbackChannel on which to send the response.
// errorChannel on which the error response is sent.
// retryCount to track the retry logic.
func (pub *Pubnub) ExecuteHereNow(channel string, callbackChannel chan []byte, errorChannel chan []byte, retryCount int) {
    count := retryCount
    
    if(InvalidChannel(channel, callbackChannel)){
        return
    }

    var hereNowUrl bytes.Buffer
    hereNowUrl.WriteString("/v2/presence")
    hereNowUrl.WriteString("/sub-key/")
    hereNowUrl.WriteString(pub.SubscribeKey)
    hereNowUrl.WriteString("/channel/")
    hereNowUrl.WriteString(channel)

    value, _, err := pub.HttpRequest(hereNowUrl.String(), false)

    if err != nil {
        pub.SendResponseToChannel(errorChannel, channel, 9, err.Error(), "")
    } else {
        //parsejson
        _, _, _, errJson := ParseJson(value, pub.CipherKey)
        if(errJson != nil && strings.Contains(errJson.Error(), _invalidJson)){
            pub.SendResponseToChannel(errorChannel, channel, 9, errJson.Error(), "")
            if (count<_maxRetries) {
                count++
                pub.ExecuteHereNow(channel, callbackChannel, errorChannel, count)    
            }           
        } else {        
            callbackChannel <- []byte(fmt.Sprintf("%s", value))
        }    
    }
}

// GetData parses the interface data and decrypts the messages if the cipher key is provided.  
//
// It accepts the following parameters:
// interface: the interface to parse.
// cipherKey: the key to decrypt the messages (can be empty).
//
// returns the decrypted and/or unescaped data json data as string.
func GetData(rawData interface{}, cipherKey string) (string){
    dataInterface := rawData.(interface{})
    switch vv := dataInterface.(type){
        case string:
            jsonData, err := json.Marshal(fmt.Sprintf("%s", vv[0]))
            if(err == nil){
                return string(jsonData)
            }else{
                return fmt.Sprintf("%s", vv[0])
            }
        case []interface{}:
            retval := ParseInterface(vv, cipherKey)
            if(retval != ""){
                return retval
            }
            
    } 
    return fmt.Sprintf("%s", rawData)    
}

// ParseInterface umarshals the response data, marshals the data again in a 
// different format and returns the json string. It also unescapes the data. 
//
// parameters:
// vv: interface array to parse and extract data from.
// cipher key: used to decrypt data. cipher key can be empty.
// 
// returns the json marshalled string.  
func ParseInterface(vv []interface{}, cipherKey string) (string){
    for i, u := range vv {
        if (reflect.TypeOf(u).Kind() == reflect.String){
            var intf interface{} 
            
            if(cipherKey != ""){
                intf = ParseCipherInterface(u, cipherKey)
                var returnedMessages interface{}

                errUnmarshalMessages := json.Unmarshal([]byte(intf.(string)), &returnedMessages)
    
                if errUnmarshalMessages == nil {
                    vv[i] = returnedMessages
                } else {                
                    vv[i] = intf
                }
            } else {
                intf = u
                unescapeVal, unescapeErr := url.QueryUnescape(intf.(string))
                if(unescapeErr != nil){
                    vv[i] = intf
                }else{
                    vv[i] = unescapeVal
                }
            }    
        }
    }
    length := len(vv) 
    if(length > 0){
        jsonData, err := json.Marshal(vv)
        if (err == nil){
            return string(jsonData)
        }else{  
            return fmt.Sprintf("%s", vv)
        }
    } 
    return ""
}

// ParseCipherInterface handles the decryption in case a cipher key is used
// in case of error it returns data as is.
// 
// parameters 
// data: the data to decrypt as interface.
// cipherKey: cipher key to use to decrypt.
//
// returns the decrypted data as interface.
func ParseCipherInterface(data interface{}, cipherKey string) (interface{}){
    var intf interface{} 
    decrypted, errDecryption := DecryptString(cipherKey, data.(string))
    if(errDecryption != nil){
        intf = data
    }else{
        intf = decrypted
    }
    return intf
}

// ParseJson parses the json data. 
// It extracts the actual data (value 0),
// Timetoken/from time in case of detailed history (value 1), 
// pubnub channelname/timetoken/to time in case of detailed history (value 2).
//
// It accepts the following parameters:
// contents: the contents to parse.
// cipherKey: the key to decrypt the messages (can be empty).
//
// returns:
// data: as string.
// Timetoken/from time in case of detailed history as string.
// pubnub channelname/timetoken/to time in case of detailed history (value 2).
// error if any.
func ParseJson (contents []byte, cipherKey string) (string, string, string, error){
    var s interface{}
    returnData := ""
    returnOne := ""
    returnTwo := ""
    
    err := json.Unmarshal(contents, &s)
    
    if err == nil {
        v := s.(interface{})
        
        switch vv := v.(type) {
           case string:
               length := len(vv)
               if(length > 0){
                   returnData = vv
               }
           case []interface{}:
               length := len(vv)
               if(length > 0){
                   returnData = GetData(vv[0], cipherKey)
               }
               if(length > 1){
                    returnOne = ParseInterfaceData(vv[1])
               }
               if(length > 2){
                   returnTwo = ParseInterfaceData(vv[2])
               }
        }
    } else {
        err = fmt.Errorf(_invalidJson)
    }
    return returnData, returnOne, returnTwo, err
}

// ParseInterfaceData formats the data to string as per the type of the data.
//
// It accepts the following parameters:
// myInterface: the interface data to parse and convert to string.
//
// returns: the data in string format.
func ParseInterfaceData(myInterface interface{}) string{
    switch v := myInterface.(type) {
        case int:
            return strconv.Itoa(v)
        case float64:
            return strconv.FormatFloat(v, 'f', -1, 64)
        case string:
            return string(v)
        }
    return fmt.Sprintf("%s", myInterface)
}

// HttpRequest is the struct Pubnub's instance method.
// It creates a connection to the pubnub origin by calling the Connect method which 
// returns the response or the error while connecting.  
//
// It accepts the following parameters:
// requestUrl: the url to connect to.
// isSubscribe: true if it is a subscribe request.
//
// returns:
// the response contents as byte array.
// response error code if any.
// error if any.
func (pub *Pubnub) HttpRequest(requestUrl string, isSubscribe bool) ([]byte, int, error) {
    contents, responseStatusCode, err := Connect(pub.Origin + requestUrl, isSubscribe)
    
    if err != nil {
        if  (strings.Contains(err.Error(), _timeout)) {
            return nil, responseStatusCode, fmt.Errorf(_operationTimeout)
        } else if (strings.Contains(fmt.Sprintf("%s", err.Error()), _closedNetworkConnection)) {
            return nil, responseStatusCode, fmt.Errorf(_connectionAborted)
        } else if (strings.Contains(fmt.Sprintf("%s", err.Error()), _noSuchHost)) {
            return nil, responseStatusCode, fmt.Errorf(_networkUnavailable)    
        } else if (strings.Contains(fmt.Sprintf("%s", err.Error()), _connectionResetByPeer)) {
            return nil, responseStatusCode, fmt.Errorf(_connectionResetByPeerU)    
        } else {
            return nil, responseStatusCode, err 
        }
    } else {
        if ((_retryCount > 0) && (isSubscribe)){
            pub.SendResponseToChannel(nil, pub.SubscribedChannels, 6, "", "")
        }
    }
    
    return contents, responseStatusCode, err
}
// SetOrGetTransport creates the transport and sets it for reuse.
// Creates a different transport for subscribe and non-subscribe requests. 
// Also sets the proxy details if provided
// It sets the timeouts based on the subscribe and non-subscribe requests.
// 
// It accepts the following parameters:
// isSubscribe: true if it is a subscribe request.
//
// returns:
// the transport.   
func SetOrGetTransport(isSubscribe bool) (http.RoundTripper){
    transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, 
        Dial: func(netw, addr string) (net.Conn, error) {
            c, err := net.DialTimeout(netw, addr, _connectTimeout * time.Second)
            
            if(c != nil){
                if(isSubscribe){
                    deadline := time.Now().Add(time.Duration(_subscribeTimeout) * time.Second)
                    c.SetDeadline(deadline)            
                    _subscribeConn = c
                } else {
                    deadline := time.Now().Add(_nonSubscribeTimeout * time.Second)
                    c.SetDeadline(deadline)
                    _conn = c
                }
            } else {
                err = fmt.Errorf(_errorInInitializing ,err.Error())
            }
                
            if err != nil {
                return nil, err
            }
            
            return c, nil
    }}

    if(_proxyServerEnabled){
        proxyUrl, err := url.Parse(fmt.Sprintf("http://%s:%s@%s:%d", _proxyUser, _proxyPassword, _proxyServer, _proxyPort))
        if(err == nil){ 
            transport.Proxy = http.ProxyURL(proxyUrl)
        } else {
            fmt.Println("Error in connecting to proxy: ", err)
        }
    }
    return transport
}

// CreateHttpClient creates the http.Client by creating or reusing the transport for 
// subscribe and non-subscribe requests. 
// 
// It accepts the following parameters:
// isSubscribe: true if it is a subscribe request.
//
// returns:
// the pointer to the http.Client
// error is any.   
func CreateHttpClient (isSubscribe bool) (*http.Client, error) {
    var transport http.RoundTripper
    
    if (isSubscribe){
        if (_subscribeTransport == nil){
            trans := SetOrGetTransport(isSubscribe)
            _subscribeTransport = trans
        }
        transport = _subscribeTransport
    }else{
        if (_transport == nil){
            trans := SetOrGetTransport(isSubscribe)
            _transport = trans
        }
        transport = _transport
    }
    
    var err error
    var httpClient *http.Client
    if(transport != nil) {
        httpClient = &http.Client{Transport: transport, CheckRedirect: nil}
    } else {
        err = fmt.Errorf("Error in initializating transport")
    }
    return httpClient, err
}

// Connect creates a http request to the pubnub origin and returns the 
// response or the error while connecting. 
// 
// It accepts the following parameters:
// requestUrl: the url to connect to.
// isSubscribe: true if it is a subscribe request.
//
// returns:
// the response as byte array.
// response errorcode if any.
// error if any.  
func Connect (requestUrl string, isSubscribe bool) ([]byte, int, error) {
    var contents []byte
    httpClient, err := CreateHttpClient(isSubscribe)
    
    if(err == nil) {
        req, err := http.NewRequest("GET", requestUrl, nil) 
        if(err == nil) {
            response, err := httpClient.Do(req)  
             if (err == nil) {
                defer response.Body.Close()
                bodyContents, e := ioutil.ReadAll(response.Body)
                if(e == nil){
                    contents = bodyContents
                    return contents, response.StatusCode, nil
                } else {
                    return nil, response.StatusCode, e
                }
            }else {
                if(response!=nil){
                    return nil, response.StatusCode, err    
                } else {
                    return nil, 0, err
                }
            }
        }else {
            return nil, 0, err
        }
    }
   
    return nil, 0, err
}

// PKCS7Padding pads the data as per the PKCS7 standard
// It accepts the following parameters:
// data: data to pad as byte array.
// returns the padded data as byte array.
func PKCS7Padding(data []byte) []byte {
    dataLen := len(data)
    var bit16 int
    if dataLen%16 == 0 {
        bit16 = dataLen
    } else {
        bit16 = int(dataLen/16+1) * 16
    }

    paddingNum := bit16 - dataLen
    bitCode := byte(paddingNum)

    padding := make([]byte, paddingNum)
    for i := 0; i < paddingNum; i++ {
        padding[i] = bitCode
    }
    return append(data, padding...)
}

// UnPKCS7Padding unpads the data as per the PKCS7 standard
// It accepts the following parameters:
// data: data to unpad as byte array.
// returns the unpadded data as byte array.
func UnPKCS7Padding(data []byte) []byte {
    dataLen := len(data)
    if dataLen == 0 {
        return data
    }
    endIndex := int(data[dataLen-1])
    if 16 > endIndex {
        if 1 < endIndex {
            for i := dataLen - endIndex; i < dataLen; i++ {
                if data[dataLen-1] != data[i] {
                    fmt.Println(" : ", data[dataLen-1], " ：", i, "  ：", data[i])
                }
            }
        }
        return data[:dataLen-endIndex]
    }
    return data
}

// GetHmacSha256 creates the cipher key hashed against SHA256.
// It accepts the following parameters:
// secretKey: the secret key.
// input: input to hash.
//
// returns the hash.
func GetHmacSha256(secretKey string, input string) string {
    hmacSha256 := hmac.New(sha256.New, []byte(secretKey))
    io.WriteString(hmacSha256, input)
    
    return fmt.Sprintf("%x", hmacSha256.Sum(nil))
}

// GenUuid generates a unique UUID
// returns the unique UUID or error.
func GenUuid() (string, error) {
    uuid := make([]byte, 16)
    n, err := rand.Read(uuid)
    if n != len(uuid) || err != nil {
        return "", err
    }
    // TODO: verify the two lines implement RFC 4122 correctly
    uuid[8] = 0x80 // variant bits see page 5
    uuid[4] = 0x40 // version 4 Pseudo Random, see page 7

    return hex.EncodeToString(uuid), nil
}

// EncodeNonAsciiChars creates unicode string of the non-ascii chars. 
// It accepts the following parameters:
// message: to parse.
//
// returns the encoded string.
func EncodeNonAsciiChars(message string) string {
    runeOfMessage := []rune(message)
    lenOfRune := len(runeOfMessage)
    encodedString := ""    
    for i := 0; i < lenOfRune; i++ {
        intOfRune := uint16(runeOfMessage[i])
        if(intOfRune>127){
            hexOfRune := strconv.FormatUint(uint64(intOfRune), 16)
            dataLen := len(hexOfRune)
            paddingNum := 4 - dataLen
            prefix := ""
            for i := 0; i < paddingNum; i++ {
                prefix += "0"
            }
            hexOfRune = prefix + hexOfRune
            encodedString += bytes.NewBufferString(`\u` + hexOfRune).String()
        } else {
            encodedString += string(runeOfMessage[i])
        }
    }
    return encodedString
}

// EncryptString creates the base64 encoded encrypted string using the cipherKey.
// It accepts the following parameters:
// cipherKey: cipher key to use to encrypt. 
// message: to encrypted.
//
// returns the base64 encoded encrypted string.
func EncryptString(cipherKey string, message string) string {
    block, _ := AesCipher(cipherKey)
    message = EncodeNonAsciiChars(message)
    value := []byte(message)
    value = PKCS7Padding(value)
    blockmode := cipher.NewCBCEncrypter(block, []byte(_IV))
    cipherBytes := make([]byte, len(value))
    blockmode.CryptBlocks(cipherBytes, value)
    
    return base64.StdEncoding.EncodeToString(cipherBytes)
}

// DecryptString decodes encrypted string using the cipherKey  
// 
// It accepts the following parameters:
// cipherKey: cipher key to use to decrypt. 
// message: to encrypted.
//
// returns the unencoded encrypted string,
// error if any.
func DecryptString(cipherKey string, message string) (retVal interface{}, err error) { 
    block, aesErr := AesCipher(cipherKey)
    if(aesErr != nil){
        return "***Decrypt Error***", fmt.Errorf("Decrypt error aes cipher: ", aesErr) 
    }
    
    value, decodeErr := base64.StdEncoding.DecodeString(message)
    if(decodeErr != nil){
        return "***Decrypt Error***", fmt.Errorf("Decrypt error on decode: ", decodeErr) 
    }
    decrypter := cipher.NewCBCDecrypter(block, []byte(_IV))
    //to handle decryption errors
    defer func(){
        if r := recover(); r != nil {
            retVal, err = "***Decrypt Error***", fmt.Errorf("Decrypt error:", r)
        }
    }()
    decrypted := make([]byte, len(value))
    decrypter.CryptBlocks(decrypted, value)
    return fmt.Sprintf("%s", string(UnPKCS7Padding(decrypted))), nil
}

// AesCipher returns the cipher block
// 
// It accepts the following parameters:
// cipherKey: cipher key. 
//
// returns the cipher block,
// error if any.
func AesCipher(cipherKey string) (cipher.Block, error) {
    key := EncryptCipherKey(cipherKey)
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    return block, nil
}
    
// EncryptCipherKey creates the 256 bit hex of the cipher key
// 
// It accepts the following parameters:
// cipherKey: cipher key to use to decrypt. 
//
// returns the 256 bit hex of the cipher key.
func EncryptCipherKey(cipherKey string) []byte {
    hash := sha256.New()
    hash.Write([]byte(cipherKey))

    sha256String := hash.Sum(nil)[:16]
    return []byte(hex.EncodeToString(sha256String))
}