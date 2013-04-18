package main 

import (
    "bufio"
    "os"
    "fmt"
    "time"
    "pubnubMessaging"
    "strings"
    "strconv"
)

var _connectChannels = ""
var _ssl bool
var _cipher = ""
var _uuid = ""
var _pub *pubnubMessaging.Pubnub

func main() {
    b := Init()
    if b {
        ch := make(chan int)
        ReadLoop(ch)
    }
    fmt.Println("Exit")
}

func Init() (b bool){
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

func AskUser() (string){
	var proxyUser string
	
    fmt.Println("Enter proxy username (optional)")
    fmt.Scanln(&proxyUser)
    
    return proxyUser
}

func AskPassword() (string){
	var proxyPassword string
	
    fmt.Println("Enter proxy password (optional)")
    fmt.Scanln(&proxyPassword)
    
    return proxyPassword
}


func ReadLoop(ch chan int){
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
                go SubscribeRoutine()
            case "2":
                fmt.Println("Please enter the message")
                message, _ , err := reader.ReadLine()
                if err != nil {
                    fmt.Println(err)
                }else{
                    go PublishRoutine(string(message))
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
                go UnsubscribeRoutine()
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
    close(ch)
}

func ParseResponseSubscribe(channel chan []byte){
    for {
        value, ok := <-channel
        if !ok {  
            fmt.Println("")            
            break
        }
        if string(value) != "[]"{
            fmt.Println(fmt.Sprintf("Subscribe: %s", value))
            //fmt.Println(fmt.Sprintf("%s", value))
            fmt.Println("")
        }
    }
}

func ParseResponsePresence(channel chan []byte){
    for {
        value, ok := <-channel
        if !ok {  
            break
        }
        if string(value) != "[]"{
            fmt.Println(fmt.Sprintf("Presence: %s ", value))
            //fmt.Println(fmt.Sprintf("%s", value))
            fmt.Println("");
        }
    }
}

func ParseResponse(channel chan []byte){
    for {
        value, ok := <-channel
        if !ok {
            break
        }
        if string(value) != "[]"{
            fmt.Println(fmt.Sprintf("Response: %s ", value))
            //fmt.Println(fmt.Sprintf("%s", value))
            fmt.Println("");
        }
    }
}

func SubscribeRoutine(){
    var subscribeChannel = make(chan []byte)
    go _pub.Subscribe(_connectChannels, subscribeChannel, false)
    ParseResponseSubscribe(subscribeChannel)
}

func PublishRoutine(message string){
    channelArray := strings.Split(_connectChannels, ",");
    
    for i:=0; i < len(channelArray); i++ {
        ch := strings.TrimSpace(channelArray[i])
        fmt.Println("Publish to channel: ",ch)
        channel := make(chan []byte)
        go _pub.Publish(ch, message, channel)
        ParseResponse(channel)
    }
}

func PresenceRoutine(){
    var presenceChannel = make(chan []byte)
    go _pub.Subscribe(_connectChannels, presenceChannel, true)
    ParseResponsePresence(presenceChannel)
}

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

func UnsubscribeRoutine(){
    channel := make(chan []byte)
    
    go _pub.Unsubscribe(_connectChannels, channel)
    ParseResponse(channel)
}

func UnsubscribePresenceRoutine(){
    channel := make(chan []byte)
    
    go _pub.PresenceUnsubscribe(_connectChannels, channel)
    ParseResponse(channel)
}

func TimeRoutine(){
    channel := make(chan []byte)
    go _pub.GetTime(channel)
    ParseResponse(channel)
}