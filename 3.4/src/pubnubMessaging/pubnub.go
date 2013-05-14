// Package pubnubMessaging provides the implemetation to connect to pubnub api.
// Build Date: May 14, 2013
// Version: 3.4
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

// Root url value of pubnub api without the http/https protocol.
const _origin = "pubsub.pubnub.com"

// The time after which the Subscribe/Presence request will timeout.
// In seconds.
const _subscribeTimeout = 30 //sec

// The time after which the Publish/HereNow/DetailedHitsory/Unsubscribe/
// UnsibscribePresence/Time  request will timeout.
// In seconds.
const _nonSubscribeTimeout = 15 //sec

// On Subscribe/Presence timeout, the number of times the reconnect attempts are made.
const _maxRetries = 50 //times

// The delay in the reconnect attempts on timeout.
// In seconds.
const _retryInterval = 10 //sec

// The HTTP transport Dial timeout.
// In seconds.
const _connectTimeout = 10 //sec

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
// ResetTimeToken: In case of a new request or an error this variable is set to true so that the 
// timeToken will be set to 0 in the next request.
// PresenceChannel: All the presence responses will be routed to this channel. This is the first 
// channel which is sent for the presence routine.
// SubscribeChannel: All the subscribe responses will be routed to this channel. This is the first 
// channel which is sent for the subscribe routine.
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
    ResetTimeToken           bool  
    PresenceChannel          chan []byte
    SubscribeChannel         chan []byte 
    NewSubscribedChannels    string
}

// VersionInfo returns the version of the this code along with the build date. 
func VersionInfo() string{
    return "Version: 3.4; Build Date: May 8, 2013;"
}

// PubnubInit initializes pubnub struct with the user provided values.
// And then initiates the origin by appending the protocol based upon the sslOn argument.
// Then it uses the customuuid or generates the uuid.
// 
// It accepts the following parameters:
// publishKey is the user specific Publish Key. Mansatory.
// subscribeKey is the user specific Subscribe Key. Mandatory.
// secretKey is the user specific Secret Key. Optional, accepts empty string if not used.
// cipherKey stores the user specific Cipher Key. Optional, accepts empty string if not used. 
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
        NewSubscribedChannels: "",
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

// Abort is the struct Pubnub's instance method that closes the open connections for both subscribe 
// and non-subscribe requests.
func (pub *Pubnub) Abort() {
    pub.SubscribedChannels = ""
    if(_conn != nil) {
        _conn.Close()
    }
    if(_subscribeConn!= nil) {
        _subscribeConn.Close()
    }
}

// GetTime is the struct Pubnub's instance method that creates a time request and sends back the 
// response to the channel.
// Closes the channel when the response is sent.
//
// It accepts the following parameters:
// Channel on which to send the response.
func (pub *Pubnub) GetTime(c chan []byte) {
    timeUrl := ""
    timeUrl += "/time"
    timeUrl += "/0"

    value, err := pub.HttpRequest(timeUrl, false)

    if err != nil {
        c <- value
    } else {
         c <- []byte(fmt.Sprintf("%s", value))
    }
    close(c)
}

// SendPublishRequest is the struct Pubnub's instance method that posts a publish request and 
// sends back the response to the channel.
//
// It accepts the following parameters:
// publishUrlString: The url to which the message is to be appended.
// jsonBytes: message to be sent.
// c: Channel on which to send the response.
func (pub *Pubnub) SendPublishRequest(publishUrlString string, jsonBytes []byte, c chan []byte) {
    var publishUrl *url.URL
    publishUrl, urlErr := url.Parse(publishUrlString)
    if urlErr != nil {
        c <- []byte(fmt.Sprintf("%s", urlErr))
    } else {
        publishUrl.Path += string(jsonBytes)
        value, err := pub.HttpRequest(publishUrl.String(), false)
    
        if err != nil {
            c <- value
        } else {
            c <- []byte(fmt.Sprintf("%s", value))
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
// Returns false all the channels is acceptable.
func InvalidChannel(channel string, c chan []byte) bool{
    if (strings.TrimSpace(channel) == "") {
        return true
    } else {
        channelArray := strings.Split(channel, ",")
    
        for i:=0; i < len(channelArray); i++ {
            if (strings.TrimSpace(channelArray[i]) == "") {    
                c <- []byte(fmt.Sprintf("Invalid Channel: %s", channel))
                close(c)    
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
// c: Channel on which to send the response back.
func (pub *Pubnub) Publish(channel string, message interface{}, c chan []byte) {
    if(InvalidChannel(channel, c)){
        return 
    }

    if(InvalidMessage(message)){
        c <- []byte(fmt.Sprintf("Invalid Message"))
        close(c)
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
    //fmt.Println("mess:", string(message))

    jsonSerialized, err := json.Marshal(message)
    if err != nil {
        c <- []byte(fmt.Sprintf("error in serializing: %s", err))
    } else {
        if pub.CipherKey != "" {
            //Encrypt and Serialize
            jsonEncBytes, errEnc := json.Marshal(EncryptString(pub.CipherKey, fmt.Sprintf("%s", jsonSerialized)))
            if errEnc != nil {
                c <- []byte(fmt.Sprintf("error in serializing: %s", errEnc))        
              } else {
                  pub.SendPublishRequest(publishUrlBuffer.String(), jsonEncBytes, c)
              }
        } else {
            pub.SendPublishRequest(publishUrlBuffer.String(), jsonSerialized, c)
        }
    }
    close(c)
}

// SendResponseToChannel is the struct Pubnub's instance method that sends a reponse on the channel 
// provided as an argument or to the subscribe / presence channel is the argument is nil. 
//
// Constructs the response based on the action (1-8). In case the action is 5 sends the response 
// as in the parameter response. 
//
// It accepts the following parameters:
// c: Channel on which to send the response back. Can be nil. If nil, assumes that if the channel name 
// is suffixed with "-pnpres" it is a presence channel else subscribe channel and send the response to the 
// respective channel.
//
// PresenceChannel: All the presence responses will be routed to this channel. 
// This is the first channel which is sent for the presence routine. The routine
// sends reponse only to one (first one that is used to init presence) go channel.
// Calling this routine again with a new go channel will not send teh response back to 
// the new channel.
//
// SubscribeChannel: All the subscribe responses will be routed to this channel. 
// This is the first channel which is sent for the subscribe routine. The routine
// sends reponse only to one (first one that is used to init subscribe) go channel.
// Calling this routine again with a new go channel will not send teh response back to 
// the new channel.
// 
// channels: Pubnub Channels to send a response to. Comma separated string for multiple channels.
// action: (1-8) 
// response: can be nil, is used only in the case action is '5'.  
func (pub *Pubnub) SendResponseToChannel(c chan []byte, channels string, action int, response []byte){
    message := ""
    intResponse := ""
    sendReponseAsIs := false
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
            message = "disconnected due to internet connection issues, trying to reconnect"
            intResponse = "1"
        case 8:
            message = "aborted due to max retry limit"
            intResponse = "1"
    }
    
    channelArray := strings.Split(channels, ",")
    for i := 0; i < len(channelArray); i++ {
        presence := "Subscription to channel "
        channel := channelArray[i]
        
        if(channel == ""){
            continue
        }

        var responseChannel = c

        if (strings.Contains(channel, "-pnpres")) {
            channel = strings.Replace(channel, "-pnpres", "", -1)
            presence = "Presence notifications for channel "
            if (responseChannel == nil){
                responseChannel = pub.PresenceChannel
            }    
        } else {
            if (responseChannel == nil){
                responseChannel = pub.SubscribeChannel
            }
        }
        
        var value string
        
        if(sendReponseAsIs){
            value = strings.Replace(string(response), "-pnpres", "", -1)
        } else {
            value = fmt.Sprintf("[%s, \"%s'%s' %s\", \"%s\"]", intResponse, presence, channel, message, channel)
        }
         
        responseChannel <- []byte(value)
    }
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
//
// Returns:
// subChannels: the Pubnub subscribed channels as a comma separated string.  
// newSubChannels: the new Pubnub subscribed channels as a comma separated string.
// b: The return parameter channelsModified is set to true if new channels are added.
func (pub *Pubnub) GetSubscribedChannels(channels string, c chan []byte, isPresenceSubscribe bool) (subChannels string, newSubChannels string, b bool) {
    channelArray := strings.Split(channels, ",")
    subscribedChannels := pub.SubscribedChannels
    newSubscribedChannels := ""
    channelsModified := false
    alreadySubscribedChannels := ""
        
    for i := 0; i < len(channelArray); i++ {
        channelToSub := strings.TrimSpace(channelArray[i])
        if(isPresenceSubscribe){
            channelToSub += "-pnpres"
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
        pub.SendResponseToChannel(c, alreadySubscribedChannels, 1, nil)
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
//
// Returns:
// b: Bool variable true incase the connection is lost.
func (pub *Pubnub) CheckForTimeoutAndRetries(err error) (bool){
    if (_retryCount >= 0) {
        //closedNetworkError :=strings.Contains(err.Error(), "closed network connection")
        errorInitConn :=strings.Contains(err.Error(), "Error in initializating connection")
        if  (errorInitConn || (strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "no such host"))){
            pub.SendResponseToChannel(nil, pub.SubscribedChannels, 7, nil)
            SleepForAWhile(true)
        }
    }
    
    if(_retryCount >= _maxRetries){
        pub.SendResponseToChannel(nil, pub.SubscribedChannels, 8, nil)
        pub.SubscribedChannels = ""
        _retryCount = 0
        return true
    }
        
    return false
}

// StartSubscribeLoop starts a continuous loop that handles the reponse from pubnub 
// subscribe/presence subscriptions.
//
// It creates subscribe request url and posts it. If the resetTimeToken flag is true 
// it sends 0 to init the subscription. 
// Else sends the last timetoken.
// When the response is received it: 
// Checks For Timeout And Retries: breaks the loop if true.
// If sent timetoken is 0 and the data is empty the connected response is sent back to the channel.
// If no error is received the response is sent to the presence or subscribe pubnub channels. 
// if the channel name is suffixed with "-pnpres" it is a presence channel else subscribe channel 
// and send the response the the respective channel.
//
// It accepts the following parameters:
// c: Channel on which to send the response back. Can be nil. If nil assumes that if the channel 
// name is suffixed with "-pnpres" it is a presence channel else subscribe channel and send the 
// response the the respective channel.
//
// TODO: refactor, remove c.
func (pub *Pubnub) StartSubscribeLoop(c chan []byte) {
    for {
          if len(pub.SubscribedChannels) > 0 {
              var subscribeUrlBuffer bytes.Buffer
              subscribeUrlBuffer.WriteString("/subscribe")
              subscribeUrlBuffer.WriteString("/")
              subscribeUrlBuffer.WriteString(pub.SubscribeKey)
              subscribeUrlBuffer.WriteString("/")
              subscribeUrlBuffer.WriteString(pub.SubscribedChannels)
              subscribeUrlBuffer.WriteString("/0")
              
            sentTimeToken := pub.TimeToken
            
            if pub.ResetTimeToken {
                subscribeUrlBuffer.WriteString("/0")
                sentTimeToken = "0"
                pub.ResetTimeToken = false
            }else{
                subscribeUrlBuffer.WriteString("/")
                if(strings.TrimSpace(pub.TimeToken) == ""){
                    pub.TimeToken = "0"
                }    
                subscribeUrlBuffer.WriteString(pub.TimeToken)
                //fmt.Println("tt:", pub.TimeToken)
            }
                
            if pub.Uuid != "" {
                subscribeUrlBuffer.WriteString("?uuid=")
                subscribeUrlBuffer.WriteString(pub.Uuid)
            }
            //fmt.Println(fmt.Sprintf("Url: %s", subscribeUrlBuffer.String()))
            value, err := pub.HttpRequest(subscribeUrlBuffer.String(), true)
            //fmt.Println(fmt.Sprintf("Value: %s", value))
            
            if err != nil {
                c <- value
                if(pub.CheckForTimeoutAndRetries(err)){
                    break
                }
            } else if string(value) != "" {
                if string(value) == "[]" {
                    SleepForAWhile(false)
                    continue
                }      
                
                data, returnTimeToken, channelName, err := ParseJson(value, pub.CipherKey)
                
                pub.TimeToken = returnTimeToken
                if (data == "[]") {
                    if(sentTimeToken == "0"){
                        pub.SendResponseToChannel(nil, pub.NewSubscribedChannels, 2, nil)
                        pub.NewSubscribedChannels = ""
                    }
                    continue
                }
                
                if err != nil {
                    pub.SendResponseToChannel(nil, channelName, 5, []byte(fmt.Sprintf("Error: %s", err)))
                    if(pub.CheckForTimeoutAndRetries(err)){
                        break
                    }
                } else {
                    if (strings.Contains(channelName, "-pnpres")) {
                        pub.SendResponseToChannel(pub.PresenceChannel, channelName, 5, value)
                    } else {
                        //in case of single subscribe request the channelname will be empty
                        if (channelName == ""){                        
                            channelName = pub.GetSubscribedChannelName()
                        }
                        
                        if(channelName != "") {
                            var buffer bytes.Buffer
                            buffer.WriteString("[")
                            buffer.WriteString(data)
                            buffer.WriteString(",\"")
                            buffer.WriteString(fmt.Sprintf("%s",pub.TimeToken))
                            buffer.WriteString("\",\"")
                            buffer.WriteString(channelName)
                            buffer.WriteString("\"]")
                            
                            pub.SendResponseToChannel(pub.SubscribeChannel, channelName, 5, buffer.Bytes())    
                        }
                    }
                }
            }
        }else {
            break;
        }
    }
    //fmt.Println("Closing Subscribe channel")
}

// GetSubscribedChannels is the struct Pubnub's instance method. 
// In case of single subscribe request the channelname will be empty.
// This methos iterates through the pubnub SubscribedChannels to find the name of the channel.
func (pub *Pubnub) GetSubscribedChannelName() (string){
    channelArray := strings.Split(pub.SubscribedChannels, ",")
    for i := 0; i < len(channelArray); i++ {
        if (strings.Contains(channelArray[i], "-pnpres")) {
            continue
        }else{
            return channelArray[i]
        }            
    }
    return ""
}

// CloseExistingConnection: Closes the open subscribe/presence connection.
func CloseExistingConnection(){
    if(_subscribeConn != nil){
        //fmt.Println("Closing connection")
        _subscribeConn.Close()
    }    
}

// Subscribe is the struct Pubnub's instance method which checks for the InvalidChannels 
// and returns if true.
// Initaiates the presence and subscribe response channels.
// PresenceChannel: All the presence responses will be routed to this channel. 
// This is the first channel which is sent for the presence routine. The routine
// sends reponse only to one (first one that is used to init presence) go channel.
// Calling this routine again with a new go channel will not send teh response back to 
// the new channel.
//
// SubscribeChannel: All the subscribe responses will be routed to this channel. 
// This is the first channel which is sent for the subscribe routine. The routine
// sends reponse only to one (first one that is used to init subscribe) go channel.
// Calling this routine again with a new go channel will not send teh response back to 
// the new channel.
//
// If there is no existing subscribe/presence loop running then it starts a 
// new loop with the new pubnub channels.
// Else closes the earlier connection.
//
// It accepts the following parameters:
// channels: comma separated pubnub channel list.
// c: Channel on which to send the response back.
// isPresenceSubscribe: tells the method that presence subscription is requested.
func (pub *Pubnub) Subscribe(channels string, c chan []byte, isPresenceSubscribe bool) {
    if(InvalidChannel(channels, c)){
        return 
    }

    pub.ResetTimeToken = true
    
    if isPresenceSubscribe {
        if(pub.PresenceChannel == nil){
            pub.PresenceChannel = c
        }
    } else {
        if(pub.SubscribeChannel == nil){
            pub.SubscribeChannel = c
        }
    }
    
    subscribedChannels, newSubscribedChannels, channelsModified := pub.GetSubscribedChannels(channels, c, isPresenceSubscribe)
    
    pub.NewSubscribedChannels = newSubscribedChannels
    
    if(pub.SubscribedChannels == ""){
        pub.SubscribedChannels = subscribedChannels
        pub.StartSubscribeLoop(c)
    }else if (channelsModified){
        CloseExistingConnection()
        pub.SubscribedChannels = subscribedChannels
    }
}    

// SleepForAWhile pauses the subscribe/presence loop for the _retryInterval. 
func SleepForAWhile(retry bool){
    if(retry) {
        _retryCount++
        fmt.Println("Retry count: ", _retryCount)
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
    for i, u := range channels {
        if channel == u {
            found = true
            pub.SendResponseToChannel(c, u, 3, nil)
        } else {
            if len(newChannels)>0 {
                newChannels += ","
            }          
            newChannels += u            
        }
        i++
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
// c: Channel on which to send the response back.
func (pub *Pubnub) Unsubscribe(channels string, c chan []byte) {
    channelArray := strings.Split(channels, ",")
    unsubscribeChannels := ""
    channelRemoved := false
    
    for i := 0; i < len(channelArray); i++ {
        if i>0 {
            unsubscribeChannels += ","
        }
        channelToUnsub := strings.TrimSpace(channelArray[i]);
        unsubscribeChannels += channelToUnsub
        removed := pub.RemoveFromSubscribeList(c, channelToUnsub)
        if !removed {
            pub.SendResponseToChannel(c, channelToUnsub, 4, nil)
        } else {
            channelRemoved = true
        }
    }
    
    if(channelRemoved) {
        CloseExistingConnection()
        pub.ResetTimeToken = true
    }
    close(c)
}

// PresenceUnsubscribe is the struct Pubnub's instance method which unsubscribes a pubnub 
// presence channel(s) from the subscribe loop. 
//
// If all the pubnub channels are not removed the method StartSubscribeLoop will take care 
// of it by starting a new loop.
// When the pubnub channel(s) are removed it creates and posts a leave request. 
// Closes the channel c when the processing is complete. 
// 
// It accepts the following parameters:
// channels: the pubnub channel(s) in a comma separated string.
// c: Channel on which to send the response back.
func (pub *Pubnub) PresenceUnsubscribe(channels string, c chan []byte) {
    channelArray := strings.Split(channels, ",")
    presenceChannels := ""
    channelRemoved := false
    
    for i := 0; i < len(channelArray); i++ {
        if i>0 {
            presenceChannels += ","
        }
        channelToUnsub := strings.TrimSpace(channelArray[i]) + "-pnpres"
        presenceChannels += channelToUnsub
        removed := pub.RemoveFromSubscribeList(c, channelToUnsub) 
        if !removed {
            pub.SendResponseToChannel(c, channelToUnsub, 4, nil)
        }else {
            channelRemoved = true
        }
    }
    
    if(channelRemoved) {
        CloseExistingConnection() 
        pub.ResetTimeToken = true
        
        var subscribeUrlBuffer bytes.Buffer
        subscribeUrlBuffer.WriteString("/v2/presence")
        subscribeUrlBuffer.WriteString("/sub-key/")
        subscribeUrlBuffer.WriteString(pub.SubscribeKey)
        subscribeUrlBuffer.WriteString("/channel/")
        subscribeUrlBuffer.WriteString(presenceChannels)
        subscribeUrlBuffer.WriteString("/leave?uuid=")
        subscribeUrlBuffer.WriteString(pub.Uuid)
        
        value, err := pub.HttpRequest(subscribeUrlBuffer.String(), false)
        c <- value
        if err != nil {
            c <- value
        }
    }    
    close(c)
}

// History is the struct Pubnub's instance method which creates and post the History request 
// for a single pubnub channel.
//
// It parses the response to get the data and return it to the channel.
// Closes the channel c when the processing is complete. 
// 
// It accepts the following parameters:
// channel: a single value of the pubnub channel.
// limit: number of history messages to return.
// start: start time from where to begin the history messages.
// end: end time till where to get the history messages.
// reverse: to fetch the messages in ascending order
// c: channel on which to send the response back.
func (pub *Pubnub) History(channel string, limit int, start int64, end int64, reverse bool, c chan []byte) {
    if(InvalidChannel(channel, c)){
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
        
    value, err := pub.HttpRequest(historyUrlBuffer.String(), false)

    if err != nil {
        c <- value
    } else {
        data, returnOne, returnTwo, err := ParseJson(value, pub.CipherKey)
        if(err != nil){
            c <- value        
        } else {
            var buffer bytes.Buffer
            buffer.WriteString("[")
            buffer.WriteString(data)
            buffer.WriteString(",\"" + returnOne + "\",\"" + returnTwo + "\"]")
               
            c <- []byte(fmt.Sprintf("%s", buffer.Bytes()))
        }
    }
    close(c)
}

// HereNow is the struct Pubnub's instance method which creates and posts the herenow 
// request to get the connected users details.  
//
// Closes the channel c when the processing is complete. 
// 
// It accepts the following parameters:
// channel: a single value of the pubnub channel. 
// c: Channel on which to send the response back.
func (pub *Pubnub) HereNow(channel string, c chan []byte) {
    if(InvalidChannel(channel, c)){
        return 
    }

    var hereNowUrl bytes.Buffer
    hereNowUrl.WriteString("/v2/presence")
    hereNowUrl.WriteString("/sub-key/")
    hereNowUrl.WriteString(pub.SubscribeKey)
    hereNowUrl.WriteString("/channel/")
    hereNowUrl.WriteString(channel)

    value, err := pub.HttpRequest(hereNowUrl.String(), false)

    if err != nil {
        c <- value
    } else {
         c <- []byte(fmt.Sprintf("%s", value))
    }
    close(c)
}

// GetData parses the interface data and decrypts the messages if the cipher key is provided.
// It also unescapes the data and recreates the json response if required to return to the channel.  
//
// It accepts the following parameters:
// interface: the interface to parse.
// cipherKey: the key to decrypt the messages (can be empty).
//
// returns the decrypted and/or unescaped data json data as string.
//
// TODO: refactor
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
            doMarshal := true
            for i, u := range vv {
                if (reflect.TypeOf(u).Kind() == reflect.String){
                    var intf interface{} 
                    if(cipherKey != ""){
                        decrypted, errDecryption := DecryptString(cipherKey, u.(string))
                        if(errDecryption != nil){
                            intf = u
                        }else{
                            decryptedData, _, _, errParseJson := ParseJson([]byte(decrypted), "")
                            if(decryptedData == ""){
                                intf = decrypted
                                doMarshal = false
                            } else if(errParseJson == nil){
                                intf = decryptedData
                            }else {
                                intf = decrypted
                                doMarshal = false
                            }                        
                        }
                    }else{
                        intf = u
                    }
                    
                    unescapeVal, unescapeErr := url.QueryUnescape(intf.(string))
                    if(unescapeErr != nil){
                        vv[i] = intf.(string)
                    }else{
                        vv[i] = unescapeVal
                    }
                }    
            }
            length := len(vv) 
            if(length > 0){
                jsonData, err := json.Marshal(vv)
                if((err == nil) && (doMarshal)){
                    return string(jsonData)
                }else{  
                    return fmt.Sprintf("%s", vv)
                }
            } 
    } 
    return fmt.Sprintf("%s", rawData)    
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
// data: as string
// Timetoken/from time in case of detailed history as string
// pubnub channelname/timetoken/to time in case of detailed history (value 2).
// error if any
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
                    //returnOne = vv[1].(string)
               }
               if(length > 2){
                   returnTwo = ParseInterfaceData(vv[2])
                   //returnTwo = vv[2].(string)
               }
        }
    } else {
        //fmt.Println("Not a valid json, err:", err)
    }
    return returnData, returnOne, returnTwo, err
}

// ParseInterfaceData formats the data to string as per the type of the data.
//
// It accepts the following parameters:
// myInterface: the interface data to parse and convert to string.
//
// returns: the data in string format
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
// error if any.
func (pub *Pubnub) HttpRequest(requestUrl string, isSubscribe bool) ([]byte, error) {
    contents, err := Connect(pub.Origin + requestUrl, isSubscribe)
    
    if err != nil {
        if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
            return []byte(fmt.Sprintf("%s: Timeout", time.Now().String())), nil
        } else if (strings.Contains(fmt.Sprintf("%s", err.Error()), "closed network connection")) {
             return []byte(fmt.Sprintf("%s: Connection aborted", time.Now().String())), nil
        } else {
            return []byte(fmt.Sprintf("Network Error: %s", err.Error())), err
        }
    } else {
        if ((_retryCount > 0) && (isSubscribe)){
            pub.SendResponseToChannel(nil, pub.SubscribedChannels, 6, nil)
        }
        _retryCount = 0
    }
    
    return contents, err
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
                    deadline := time.Now().Add(_subscribeTimeout * time.Second)
                    c.SetDeadline(deadline)            
                    _subscribeConn = c
                } else {
                    deadline := time.Now().Add(_nonSubscribeTimeout * time.Second)
                    c.SetDeadline(deadline)
                    _conn = c
                }
            } else {
                err = fmt.Errorf("Error in initializating connection")
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
// error is any   
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
// error if any.  
func Connect (requestUrl string, isSubscribe bool) ([]byte, error) {
    var contents []byte
    httpClient, err := CreateHttpClient(isSubscribe)
    
    if(err == nil) {
        req, err := http.NewRequest("GET", requestUrl, nil) 
         
        if(err == nil) {
            response, err := httpClient.Do(req)  
            //response, err := httpClient.Get(url)
             if (err == nil) {
                defer response.Body.Close()
                bodyContents, e := ioutil.ReadAll(response.Body)
                if(e == nil){
                    contents = bodyContents
                    return contents, nil
                } else {
                    return nil, e
                }
            }else {
                return nil, err
            }
        }else {
            return nil, err
        }
    }
   
    return nil, err
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
// returns the base64 encoded encrypted string
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
// returns the unencoded encrypted string
// error if any
func DecryptString(cipherKey string, message string) (retVal string, err error) { 
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
// returns the cipher block
// error if any
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
// returns the 256 bit hex of the cipher key
func EncryptCipherKey(cipherKey string) []byte {
    hash := sha256.New()
    hash.Write([]byte(cipherKey))

    sha256String := hash.Sum(nil)[:16]
    return []byte(hex.EncodeToString(sha256String))
}