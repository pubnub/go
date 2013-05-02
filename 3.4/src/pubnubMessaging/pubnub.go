// Package pubnubMessaging provides the implemetation to connect to pubnub api
// TODO change string concat to buffer
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
)

const _origin = "pubsub.pubnub.com"
const _subscribeTimeout = 30 //sec
const _nonSubscribeTimeout = 15 //sec
const _maxRetries = 5 //times
const _retryInterval = 10 //sec
const _connectTimeout = 10 //sec

var _conn net.Conn
var _subscribeConn net.Conn
var _subscribeTransport http.RoundTripper
var _transport http.RoundTripper

var _retryCount = 0

var _proxyServer string
var _proxyPort int
var _proxyUser string
var _proxyPassword string
    
var _proxyServerEnabled = false

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

//Init pubnub struct
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

func SetProxy(proxyServer string, proxyPort int, proxyUser string, proxyPassword string){
    _proxyServer = proxyServer
    _proxyPort = proxyPort
    _proxyUser = proxyUser
    _proxyPassword = proxyPassword
    _proxyServerEnabled = true
}

func (pub *Pubnub) Abort() {
    pub.SubscribedChannels = ""
    if(_conn != nil) {
        _conn.Close()
    }
    if(_subscribeConn!= nil) {
        _subscribeConn.Close()
    }
}

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
    publishUrl := ""
    publishUrl += "/publish"
    publishUrl += "/" + pub.PublishKey
    publishUrl += "/" + pub.SubscribeKey
    publishUrl += "/" + signature
    publishUrl += "/" + channel
    publishUrl += "/0/"
    
    //fmt.Println("mess:", string(message))

    jsonSerialized, err := json.Marshal(message)
    if err != nil {
        c <- []byte(fmt.Sprintf("error in serializing: %s", err))
    } else {
        if pub.CipherKey != "" {
            jsonEncBytes, errEnc := json.Marshal(EncryptString(pub.CipherKey, fmt.Sprintf("%s", jsonSerialized)))
            if errEnc != nil {
                c <- []byte(fmt.Sprintf("error in serializing: %s", errEnc))        
              } else {
                  pub.SendPublishRequest(publishUrl, jsonEncBytes, c)
              }
        } else {
            pub.SendPublishRequest(publishUrl, jsonSerialized, c)
        }
    }
    close(c)
}

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

func (pub *Pubnub) CheckForTimeoutAndRetries(err error) (bool){
    //if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "no such host") {
    if (_retryCount == 0) {
        if !strings.Contains(err.Error(), "closed network connection") {
            pub.SendResponseToChannel(nil, pub.SubscribedChannels, 7, nil)
        }
    }
    
    SleepForAWhile(true)
    
    if(_retryCount >= _maxRetries){
        pub.SendResponseToChannel(nil, pub.SubscribedChannels, 8, nil)
        pub.SubscribedChannels = ""
        _retryCount = 0
        return true
    }
        
    //}
    return false
}

//TODO refactor
func (pub *Pubnub) StartSubscribeLoop(c chan []byte) {
    for {
          if len(pub.SubscribedChannels) > 0 {
            subscribeUrl := ""
            subscribeUrl += "/subscribe"
            subscribeUrl += "/" + pub.SubscribeKey
            subscribeUrl += "/" + pub.SubscribedChannels
            subscribeUrl += "/0"
            
            sentTimeToken := pub.TimeToken
            
            if pub.ResetTimeToken {
                subscribeUrl += "/0"
                sentTimeToken = "0"
                pub.ResetTimeToken = false
            }else{
                subscribeUrl += "/" + pub.TimeToken
               }
                
            if pub.Uuid != "" {
                subscribeUrl += "?uuid=" + pub.Uuid
            }
            //fmt.Println(fmt.Sprintf("Url: %s", url))
            value, err := pub.HttpRequest(subscribeUrl, true)
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
                            buffer.WriteString(",\"" + fmt.Sprintf("%s",pub.TimeToken) + "\",\"" + channelName + "\"]")
                            
                            pub.SendResponseToChannel(pub.SubscribeChannel, channelName, 5, buffer.Bytes())    
                        }
                    }
                }
            }
        }else {
            break;
        }
    }
    fmt.Println("Closing Subscribe channel")
}

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

func CloseExistingConnection(){
    if(_subscribeConn != nil){
        fmt.Println("Closing connection")
        _subscribeConn.Close()
    }    
}

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

func SleepForAWhile(retry bool){
    if(retry) {
        _retryCount++
        fmt.Println("Retry count: ", _retryCount)
    }
    time.Sleep(_retryInterval * time.Second)
}

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
        
        subscribeUrl := ""
        subscribeUrl += "/v2/presence"
        subscribeUrl += "/sub-key/" + pub.SubscribeKey
        subscribeUrl += "/channel/" + presenceChannels
        subscribeUrl += "/leave?uuid=" + pub.Uuid
            
        value, err := pub.HttpRequest(subscribeUrl, false)
        c <- value
        if err != nil {
            c <- value
        }
    }    
    close(c)
}

func (pub *Pubnub) History(channel string, limit int, start int64, end int64, reverse bool, c chan []byte) {
    if(InvalidChannel(channel, c)){
        return 
    }

    if(limit < 0){
        limit = 100
    }
    
    parameters := "&reverse=" + fmt.Sprintf("%t", reverse)
    if(start > 0){
        parameters += "&start=" + fmt.Sprintf("%d", start)
    }
    if(end > 0){
        parameters += "&end=" + fmt.Sprintf("%d", end)
    }

    historyUrl := ""
    historyUrl += "/v2/history"
    historyUrl += "/sub-key/" + pub.SubscribeKey
    historyUrl += "/channel/" + channel
    historyUrl += "?count=" + fmt.Sprintf("%d", limit)
    historyUrl += parameters
    
    value, err := pub.HttpRequest(historyUrl, false)

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
            buffer.WriteString(",\"" + fmt.Sprintf("%s",returnOne) + "\",\"" + returnTwo + "\"]")
               
            c <- []byte(fmt.Sprintf("%s", buffer.Bytes()))
        }
    }
    close(c)
}

func (pub *Pubnub) HereNow(channel string, c chan []byte) {
    if(InvalidChannel(channel, c)){
        return 
    }

    hereNowUrl := ""
    hereNowUrl += "/v2/presence"
    hereNowUrl += "/sub-key/" + pub.SubscribeKey
    hereNowUrl += "/channel/" + channel

    value, err := pub.HttpRequest(hereNowUrl, false)

    if err != nil {
        c <- value
    } else {
         c <- []byte(fmt.Sprintf("%s", value))
    }
    close(c)
}

//TODO: refactor
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

func UnescapeContents(contents []byte) ([]byte){
    if(contents != nil){
        stringContents := string(contents)
        stringContents, err := url.QueryUnescape(stringContents)
        if(err == nil){
            contents = []byte(stringContents)
            return contents
        } 
    }
    return contents
}

func ParseJson (contents []byte, cipherKey string) (string, string, string, error){
    //contents = UnescapeContents(contents)
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
                   returnOne = fmt.Sprintf("%s", vv[1])
               }
               if(length > 2){
                   returnTwo = fmt.Sprintf("%s", vv[2])
               }
        }
    } else {
        //fmt.Println("Not a valid json, err:", err)
    }
    return returnData, returnOne, returnTwo, err
}

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
