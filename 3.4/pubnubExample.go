// Package main provides the example implemetation to connect to pubnub api.
// Runs on the console. 
package main 

import (
    "bufio"
    "os"
    "fmt"
    "time"
    "github.com/pubnub/go/3.4/pubnubMessaging"
    "strings"
    "strconv"
    "unicode/utf16"
    "unicode/utf8"
    "encoding/binary"
)

// _connectChannels: the conected pubnub channels, multiple channels are stored separated by comma.
var _connectChannels = ""

// _ssl: true if the ssl is enabled else false.
var _ssl bool

// _cipher: stores the cipher key set by the user.
var _cipher = ""

// _uuid stores the custom uuid set by the user.
var _uuid = ""

// _pub instance of the Pubnub struct from the pubnubMessaging package.
var _pub *pubnubMessaging.Pubnub

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
func Init() (b bool){
    fmt.Println("PubNub Api for go;", pubnubMessaging.VersionInfo())
    fmt.Println("Please enter the channel name(s). Enter multiple channels separated by comma.")
    reader := bufio.NewReader(os.Stdin)
    
    line, _ , err := reader.ReadLine()
    if err != nil {
        fmt.Println(err)
    }else{
        _connectChannels = string(line)
        if strings.TrimSpace(_connectChannels) != "" { 
            fmt.Println("Channel: ", _connectChannels)
            fmt.Println("Enable SSL. Enter y for Yes, n for No.")
            var enableSsl string
            fmt.Scanln(&enableSsl)
            
            if enableSsl == "y" || enableSsl == "Y" {
                _ssl = true
                fmt.Println("SSL enabled")    
            }else{
                _ssl = false
                fmt.Println("SSL disabled")
            }
            
            fmt.Println("Please enter a CIPHER key, leave blank if you don't want to use this.")
            fmt.Scanln(&_cipher)
            fmt.Println("Cipher: ", _cipher)
            
            fmt.Println("Please enter a Custom UUID, leave blank for default.")
            fmt.Scanln(&_uuid)
            fmt.Println("UUID: ", _uuid)
            
            pubInstance := pubnubMessaging.PubnubInit("demo", "demo", "", _cipher, _ssl, _uuid)
            
            _pub = pubInstance
            
            SetupProxy()
            
            return true
        }else{
            fmt.Println("Channel cannot be empty.")
        }    
    }
    return false
}

// SetupProxy asks the user the Proxy details and calls the SetProxy of the pubnubMessaging 
// package with the details. 
func SetupProxy(){
    fmt.Println("Using Proxy? Enter y to setup.")
    var enableProxy string
    fmt.Scanln(&enableProxy)
            
    if enableProxy == "y" || enableProxy == "Y" {
        proxyServer := AskServer();
        proxyPort := AskPort();
        proxyUser := AskUser();
        proxyPassword := AskPassword();
    
        pubnubMessaging.SetProxy(proxyServer, proxyPort, proxyUser, proxyPassword)
    
        fmt.Println("Proxy sever set")    
    }else{
        fmt.Println("Proxy not used")
    }
}

// AskServer asks the user to enter the proxy server name or IP. 
// It validates the input and returns the value if validated.
func AskServer() (string){
    var proxyServer string
    
    fmt.Println("Enter proxy servername or IP.")
    fmt.Scanln(&proxyServer)
    
    if(strings.TrimSpace(proxyServer) == ""){
        fmt.Println("Proxy servername or IP is empty.")
        AskServer()
    }
    return proxyServer
}

// AskPort asks the user to enter the proxy port number. 
// It validates the input and returns the value if validated.
func AskPort() (int){
    var proxyPort string
    
    fmt.Println("Enter proxy port.")
    fmt.Scanln(&proxyPort)
    
    port, err := strconv.Atoi(proxyPort)
    if (err != nil) || ((port <= 0) || (port > 65536)) {
        fmt.Println("Proxy port is invalid.")
        AskPort()
    }
    return port
}

// AskUser asks the user to enter the proxy username. 
// returns the value, can be empty.
func AskUser() (string){
    var proxyUser string
    
    fmt.Println("Enter proxy username (optional)")
    fmt.Scanln(&proxyUser)
    
    return proxyUser
}

// AskPassword asks the user to enter the proxy password. 
// returns the value, can be empty.
func AskPassword() (string){
    var proxyPassword string
    
    fmt.Println("Enter proxy password (optional)")
    fmt.Scanln(&proxyPassword)
    
    return proxyPassword
}

// AskChannel asks the user to channel name.
// If the channel(s) are not provided the channel(s) provided by the user
// at the beginning will be used.
// returns the read channel(s), or error
func AskChannel() (string, error){
    fmt.Println("Please enter the channel name. Leave empty to use the channel(s) provided at the beginning.")
    reader := bufio.NewReader(os.Stdin)
    channels, _ , errReadingChannel := reader.ReadLine()
    if(errReadingChannel != nil){
        return "", errReadingChannel
    } else {
        if(strings.TrimSpace(string(channels)) == ""){
            fmt.Println("Using channel(s): ", _connectChannels)
            return _connectChannels, nil   
        }  
    }
    return string(channels), nil
}

// UTF16BytesToString converts UTF-16 encoded bytes, in big or little endian byte order,
// to a UTF-8 encoded string.
func UTF16BytesToString(b []byte, o binary.ByteOrder) string {
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
func ReadLoop(){
    fmt.Println("")
    fmt.Println("ENTER 1 FOR Subscribe")
    fmt.Println("ENTER 2 FOR Publish")
    fmt.Println("ENTER 3 FOR Presence")
    fmt.Println("ENTER 4 FOR Detailed History")
    fmt.Println("ENTER 5 FOR Here_Now")
    fmt.Println("ENTER 6 FOR Unsubscribe")
    fmt.Println("ENTER 7 FOR Presence-Unsubscribe")
    fmt.Println("ENTER 8 FOR Time")
    fmt.Println("ENTER 9 FOR Exit")
    fmt.Println("")
    reader := bufio.NewReader(os.Stdin)
    
    for{
        var action string
        fmt.Scanln(&action)
        breakOut := false
        switch action {
            case "1":
                fmt.Println("Running Subscribe")
                channels, errReadingChannel := AskChannel()
                if(errReadingChannel != nil){
                    fmt.Println("errReadingChannel: ", errReadingChannel)
                } else {                
                    go SubscribeRoutine(channels)
                }    
            case "2":
                channels, errReadingChannel := AskChannel()
                if(errReadingChannel != nil){
                    fmt.Println("errReadingChannel: ", errReadingChannel)
                } else {
                    fmt.Println("Please enter the message")
                    message, _ , err := reader.ReadLine()
                    if err != nil {
                        fmt.Println(err)
                    }else{
                        go PublishRoutine(channels, string(message))
                    }
                }
            case "3":
                fmt.Println("Running Presence")
                go PresenceRoutine()    
            case "4":
                fmt.Println("Running detailed history")
                go DetailedHistoryRoutine()
            case "5":
                fmt.Println("Running here now")
                go HereNowRoutine()            
            case "6":
                fmt.Println("Running Unsubscribe")
                channels, errReadingChannel := AskChannel()
                if(errReadingChannel != nil){
                    fmt.Println("errReadingChannel: ", errReadingChannel)
                } else {                
                    go UnsubscribeRoutine(channels)
                }    
            case "7":
                fmt.Println("Running Unsubscribe Presence")
                go UnsubscribePresenceRoutine()
            case "8":
                fmt.Println("Running Time")
                go TimeRoutine()
            case "9":
                fmt.Println("Exiting") 
                _pub.Abort()   
                breakOut = true
            default: 
                fmt.Println("Invalid choice!")            
        }
        if breakOut {
            break
        }else{
            time.Sleep(1000 * time.Millisecond)
        }
    }
}

// ParseResponseSubscribe parses the response of the Subscribed pubnub channel.
// It prints the response as-is in the console.
func ParseResponseSubscribe(channel chan []byte){
    for {
        value, ok := <-channel
        if !ok {  
            fmt.Println("")            
            break
        }
        if string(value) != "[]"{
            fmt.Println(fmt.Sprintf("Subscribe: %s", value))
            fmt.Println("")
        }
    }
}

// ParseResponsePresence parses the response of the presence subscription pubnub channel.
// It prints the response as-is in the console.
func ParseResponsePresence(channel chan []byte){
    for {
        value, ok := <-channel
        if !ok {  
            break
        }
        if string(value) != "[]"{
            fmt.Println(fmt.Sprintf("Presence: %s ", value))
            fmt.Println("");
        }
    }
}

// ParseResponse parses the response of all the other activities apart 
// from subscribe and presence on the pubnub channel.
// It prints the response as-is in the console.
func ParseResponse(channel chan []byte){
    for {
        value, ok := <-channel
        if !ok {
            break
        }
        if string(value) != "[]"{
            fmt.Println(fmt.Sprintf("Response: %s ", value))
            fmt.Println("");
        }
    }
}

// SubscribeRoutine calls the Subscribe routine of the pubnubMessaging package
// as a parallel process. 
func SubscribeRoutine(channels string){
    var subscribeChannel = make(chan []byte)
    go _pub.Subscribe(channels, subscribeChannel, false)
    ParseResponseSubscribe(subscribeChannel)    
}

// PublishRoutine asks the user the message to send to the pubnub channel(s) and 
// calls the Publish routine of the pubnubMessaging package as a parallel 
// process. If we have multiple pubnub channels then this method will spilt the 
// channel by comma and send the message on all the pubnub channels.
func PublishRoutine(channels string, message string){
    channelArray := strings.Split(channels, ",");
    
    for i:=0; i < len(channelArray); i++ {
        ch := strings.TrimSpace(channelArray[i])
        fmt.Println("Publish to channel: ",ch)
        channel := make(chan []byte)
        go _pub.Publish(ch, message, channel)
        ParseResponse(channel)
    }
}

// PresenceRoutine calls the Subscribe routine of the pubnubMessaging package,
// by setting the last argument as true, as a parallel process. 
func PresenceRoutine(){
    var presenceChannel = make(chan []byte)
    go _pub.Subscribe(_connectChannels, presenceChannel, true)
    ParseResponsePresence(presenceChannel)
}

// DetailedHistoryRoutine calls the History routine of the pubnubMessaging package as a parallel 
// process. If we have multiple pubnub channels then this method will spilt the _connectChannels 
// by comma and send the message on all the pubnub channels.
func DetailedHistoryRoutine(){
    channelArray := strings.Split(_connectChannels, ",");
    for i:=0; i < len(channelArray); i++ {
        ch := strings.TrimSpace(channelArray[i])
        fmt.Println("DetailedHistory for channel: ", ch)
        
        channel := make(chan []byte)
        
        //go _pub.History(ch, 100, 13662867154115803, 13662867243518473, false, channel)
        go _pub.History(ch, 100, 0, 0, false, channel)
        ParseResponse(channel)
    }
}

// HereNowRoutine calls the HereNow routine of the pubnubMessaging package as a parallel 
// process. If we have multiple pubnub channels then this method will spilt the _connectChannels 
// by comma and send the message on all the pubnub channels.
func HereNowRoutine(){
    channelArray := strings.Split(_connectChannels, ",");
    for i:=0; i < len(channelArray); i++ {    
        channel := make(chan []byte)
        ch := strings.TrimSpace(channelArray[i])
        fmt.Println("HereNow for channel: ", ch)
        
        go _pub.HereNow(ch, channel)
        ParseResponse(channel)
    }
}
	
// UnsubscribeRoutine calls the Unsubscribe routine of the pubnubMessaging package as a parallel 
// process. All the channels in the _connectChannels string will be unsubscribed.
func UnsubscribeRoutine(channels string){
    channel := make(chan []byte)
    
    go _pub.Unsubscribe(channels, channel)
    ParseResponse(channel)
}

// UnsubscribePresenceRoutine calls the PresenceUnsubscribe routine of the pubnubMessaging package as a parallel 
// process. All the channels in the _connectChannels string will be unsubscribed.
func UnsubscribePresenceRoutine(){
    channel := make(chan []byte)
    
    go _pub.PresenceUnsubscribe(_connectChannels, channel)
    ParseResponse(channel)
}

// TimeRoutine calls the GetTime routine of the pubnubMessaging package as a parallel 
// process. 
func TimeRoutine(){
    channel := make(chan []byte)
    go _pub.GetTime(channel)
    ParseResponse(channel)
}