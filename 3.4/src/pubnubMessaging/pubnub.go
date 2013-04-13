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
)

const _origin = "pubsub.pubnub.com"
const _subscribeTimeout = 310 //sec
const _nonSubscribeTimeout = 15 //sec
const _maxRetries = 50 //sec
const _retryInterval = 10 //sec

var _conn net.Conn
var _subscribeConn net.Conn
var _subscribeTransport http.RoundTripper
var _transport http.RoundTripper

type Pubnub struct {
    origin                string
    publishKey           string
    subscribeKey         string
    secretKey            string
    cipherKey            string
    ssl                   bool
    uuid                  string
    subscribedChannels     string 
    timeToken            string
    resetTimeToken         bool  
    presenceChannel     chan []byte
    subscribeChannel     chan []byte 
    newSubscribedChannels     string
}

//Init pubnub struct
func PubnubInit(publishKey string, subscribeKey string, secretKey string, cipherKey string, sslOn bool, customUuid string) *Pubnub {
    newPubnub := &Pubnub{
        origin:                _origin,
        publishKey:           publishKey,
        subscribeKey:         subscribeKey,
        secretKey:            secretKey,
        cipherKey:            cipherKey,
        ssl:                   sslOn,
        uuid:                  "",
        subscribedChannels: "",
        resetTimeToken:        true,
        timeToken:            "0",
        newSubscribedChannels: "",
    }

    if newPubnub.ssl {
        newPubnub.origin = "https://" + newPubnub.origin
    } else {
        newPubnub.origin = "http://" + newPubnub.origin
    }

    if customUuid == "" {
        uuid, err := GenUuid()
        if err == nil {
            newPubnub.uuid = uuid
        } else {
            fmt.Println(err)
        }
    } else {
        newPubnub.uuid = customUuid
    }

    return newPubnub
}

func (pub *Pubnub) Abort() {
    pub.subscribedChannels = ""
    _conn.Close()
    _subscribeConn.Close()
}

func (pub *Pubnub) GetTime(c chan []byte) {
    url := ""
    url += "/time"
    url += "/0"

    value, err := pub.HttpRequest(url, false)

    if err != nil {
        c <- value
    } else {
         c <- []byte(fmt.Sprintf("%s", value))
    }
    close(c)
}

func (pub *Pubnub) Publish(channel string, message string, c chan []byte) {
    signature := ""
    if pub.secretKey != "" {
        signature = GetHmacSha256(pub.secretKey, fmt.Sprintf("%s/%s/%s/%s/%s", pub.publishKey, pub.subscribeKey, pub.secretKey, channel, message))
    } else {
        signature = "0"
    }
    url := ""
    url += "/publish"
    url += "/" + pub.publishKey
    url += "/" + pub.subscribeKey
    url += "/" + signature
    url += "/" + channel
    url += "/0"

    //Now only for string, need add encrypt for other types
    // use "/{\"msg\":\"%s\"}" for sending hash 
    if pub.cipherKey != "" {
        url += fmt.Sprintf("/\"%s\"", EncryptString(pub.cipherKey, fmt.Sprintf("\"%s\"", message)))
    } else {
        url += fmt.Sprintf("/\"%s\"", message)
    }

    value, err := pub.HttpRequest(url, false)

    if err != nil {
        c <- value
    } else {
         c <- []byte(fmt.Sprintf("%s", value))
    }
    close(c)
}

func (pub *Pubnub) SendResponseToChannel(c chan []byte, channels string, action int){
    message := ""
    intResponse := ""
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
    }
    
    channelArray := strings.Split(channels, ",")
    for i := 0; i < len(channelArray); i++ {
        presence := ""
        channel := channelArray[i]
        
        if(channel == ""){
            continue
        }

        var responseChannel = c
        
        if (strings.Contains(channel, "-pnpres")) {
            channel = strings.Replace(channel, "-pnpres", "", -1)
            presence = "Presence notifications for "
            if (responseChannel == nil){
                responseChannel = pub.presenceChannel
            }    
        } else {
            if (responseChannel == nil){
                responseChannel = pub.subscribeChannel
            }
        }
        
         value := fmt.Sprintf("[%s, \"%s%s %s\", \"%s\"]", intResponse, presence, channel, message, channel)
         
        responseChannel <- []byte(value)
    }
}

func (pub *Pubnub) GetSubscribedChannels(channels string, c chan []byte, isPresenceSubscribe bool) (subChannels string, newSubChannels string, b bool) {
    channelArray := strings.Split(channels, ",")
    subscribedChannels := pub.subscribedChannels
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
        pub.SendResponseToChannel(c, alreadySubscribedChannels, 1)
    }    

    return subscribedChannels, newSubscribedChannels, channelsModified
}

func (pub *Pubnub) StartSubscribeLoop(c chan []byte) {
    for {
          if len(pub.subscribedChannels) > 0 {
            url := ""
            url += "/subscribe"
            url += "/" + pub.subscribeKey
            url += "/" + pub.subscribedChannels
            url += "/0"
            
            sentTimeToken := pub.timeToken
            
            if pub.resetTimeToken {
                url += "/0"
                sentTimeToken = "0"
                pub.resetTimeToken = false
            }else{
                url += "/" + pub.timeToken
               }
                
            if pub.uuid != "" {
                url += "?uuid=" + pub.uuid
            }
            //fmt.Println(fmt.Sprintf("Url: %s", url))
            value, err := pub.HttpRequest(url, true)
            //fmt.Println(fmt.Sprintf("Value: %s", value))
            
            if err != nil {
                c <- value
                if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "no such host") {
                    SleepForAWhile()
                }
            } else if string(value) != "" {
                if string(value) == "[]" {
                    SleepForAWhile()
                    continue
                }      
                
                   data, returnTimeToken, channelName , err := ParseJson(value)
                
                pub.timeToken = returnTimeToken
                if (data == "[]") {
                    if(sentTimeToken == "0"){
                        pub.SendResponseToChannel(nil, pub.newSubscribedChannels, 2)
                        pub.newSubscribedChannels = ""
                    }
                    continue
                }
                           
                if err != nil {  
                    fmt.Println(fmt.Sprintf("Error: %s", err))
                }

                if (strings.Contains(channelName, "-pnpres")) {
                       pub.presenceChannel <- []byte(strings.Replace(string(value), "-pnpres", "", -1))
                   } else {
                       pub.subscribeChannel <- []byte(strings.Replace(string(value), "-pnpres", "", -1))
                   }
                
                /*if pub.cipherKey != "" {
                    c <- []byte(DecryptString(pub.cipherKey, fmt.Sprintf("%s", value)))
                } else {
                       //c <- []byte(fmt.Sprintf("%s %s %s %s ", value, data, timeToken, returnedChannels))
                    c <- []byte(fmt.Sprintf("%s", value))
                  }*/
            }
        }else {
            break;
        }
    }
    fmt.Println("Closing Subscribe channel")
}

func CloseExistingConnection(){
    if(_subscribeConn != nil){
        fmt.Println("Closing connection")
        _subscribeConn.Close()
    }    
}

func (pub *Pubnub) Subscribe(channels string, c chan []byte, isPresenceSubscribe bool) {
    
    pub.resetTimeToken = true
    
    if isPresenceSubscribe {
        if(pub.presenceChannel == nil){
            pub.presenceChannel = c
        }
    } else {
        if(pub.subscribeChannel == nil){
            pub.subscribeChannel = c
        }
    }
    
    subscribedChannels, newSubscribedChannels, channelsModified := pub.GetSubscribedChannels(channels, c, isPresenceSubscribe)
    pub.newSubscribedChannels = newSubscribedChannels
    
    if(pub.subscribedChannels == ""){
        pub.subscribedChannels = subscribedChannels
        pub.StartSubscribeLoop(c)
    }else if (channelsModified){
        CloseExistingConnection()
        pub.subscribedChannels = subscribedChannels
    }
}    

func SleepForAWhile(){
    //TODO: change to reconnect val
    time.Sleep(_retryInterval * time.Second)
}

func (pub *Pubnub) NotDuplicate(channel string) (b bool){
    var channels = strings.Split(pub.subscribedChannels, ",")
    for i, u := range channels {
        if channel == u {
            return false
        } 
        i++
    }
    return true 
}

func (pub *Pubnub) RemoveFromSubscribeList(c chan []byte, channel string) (b bool){
    var channels = strings.Split(pub.subscribedChannels, ",")
    newChannels := ""
    found := false
    for i, u := range channels {
        if channel == u {
            found = true
            pub.SendResponseToChannel(c, u, 3)
        } else {
            if len(newChannels)>0 {
                newChannels += ","
            }          
            newChannels += u            
        }
        i++
    }
    if found {
        pub.subscribedChannels = newChannels
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
            pub.SendResponseToChannel(c, channelToUnsub, 4)
        } else {
            channelRemoved = true
        }
    }
    
    if(channelRemoved) {
        CloseExistingConnection()
        pub.resetTimeToken = true
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
            pub.SendResponseToChannel(c, channelToUnsub, 4)
        }else {
            channelRemoved = true
        }
    }
    
    if(channelRemoved) {
        CloseExistingConnection() 
        pub.resetTimeToken = true
        
        url := ""
        url += "/v2/presence"
        url += "/sub-key/" + pub.subscribeKey
        url += "/channel/" + presenceChannels
        url += "/leave?uuid=" + pub.uuid
            
        value, err := pub.HttpRequest(url, false)
        c <- value
        if err != nil {
            c <- value
        }
    }    
    close(c)
}

func (pub *Pubnub) History(channel string, limit int, c chan []byte) {
    url := ""
    url += "/history"
    url += "/" + pub.subscribeKey
    url += "/" + channel
    url += "/0"
    url += "/" + fmt.Sprintf("%d", limit)

    value, err := pub.HttpRequest(url, false)

    if err != nil {
        c <- value
    } else {
         c <- []byte(fmt.Sprintf("%s", value))
    }
    close(c)
}

func (pub *Pubnub) HereNow(channel string, c chan []byte) {
    url := ""
    url += "/v2/presence"
    url += "/sub-key/" + pub.subscribeKey
    url += "/channel/" + channel

    value, err := pub.HttpRequest(url, false)

    if err != nil {
        c <- value
    } else {
         c <- []byte(fmt.Sprintf("%s", value))
    }
    close(c)
}

func ParseJson (contents []byte) (data string, timeToken string, channels string, err error){
    var s interface{}
    returnData := ""
    returnTimeToken := ""
    returnChannels := ""
    
    if err := json.Unmarshal(contents, &s); err == nil {
        v := s.(interface{})
        
        switch vv := v.(type) {
           case []interface{}:
               length := len(vv)
               if(length > 0){
                   returnData = fmt.Sprintf("%s", vv[0])
               }
               if(length > 1){
                   returnTimeToken = fmt.Sprintf("%s", vv[1])
               }
               if(length > 2){
                   returnChannels = fmt.Sprintf("%s", vv[2])
               }
        }
    } else {
        //fmt.Println("Not a valid json, err:", err)
    }
    return returnData, returnTimeToken, returnChannels, err
}

func (pub *Pubnub) HttpRequest(url string, subscribe bool) ([]byte, error) {
    contents, err := Connect(pub.origin+url, subscribe)
    
    if err != nil {
        if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
            return []byte(fmt.Sprintf("%s: Reconnecting from timeout", time.Now().String())), nil
        } else if (strings.Contains(fmt.Sprintf("%s", err.Error()), "closed network connection")) {
             return []byte(fmt.Sprintf("%s: Connection aborted", time.Now().String())), nil
        } else {
            return []byte(fmt.Sprintf("Network Error: %s", err.Error())), err
        }
    }
    
    return contents, err
}

func SetOrGetTransport(isSubscribe bool) (http.RoundTripper){
    transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        Dial: func(netw, addr string) (net.Conn, error) {
            //connect timeout
            deadline := time.Now().Add(_subscribeTimeout * time.Second)
            c, err := net.DialTimeout(netw, addr, time.Second)
            if(isSubscribe){
                //rw timeout
                c.SetDeadline(deadline)            
                _subscribeConn = c
            } else {
                //rw timeout
                //deadline := time.Now().Add(_nonSubscribeTimeout * time.Second)
                c.SetDeadline(deadline)
                _conn = c
            }
            
            if err != nil {
                return nil, err
            }
            
            return c, nil
    }}
    return transport
}

func CreateHttpClient (isSubscribe bool) (*http.Client) {
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

    httpClient := &http.Client{Transport: transport, CheckRedirect: nil}
    return httpClient
}

func Connect (url string, isSubscribe bool) ([]byte, error) {
    var contents []byte
    httpClient := CreateHttpClient(isSubscribe)
    response, err := httpClient.Get(url)
    
    if (err == nil) {
        defer response.Body.Close()
        bodyContents, e := ioutil.ReadAll(response.Body)
        err = e
        contents = bodyContents
    }
    return contents, err
}
