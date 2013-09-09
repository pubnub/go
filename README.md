#PubNub 3.4.1 client for Go 1.0.3, 1.1

###Features
* Supports multiplexing, UUID, SSL, Encryption, Proxy, and godoc
* This version is not backward compatible. The major change is in the func calls. A new parameter "error callback" is added to the major functions of the pubnub class.
* The client now supports:
* Error Callback: All the error messages are routed to this callback channel
* Resume on reconnect
* You can now "Subscribe with timetoken"
* An example of Disconnect/Retry has been added in the example 
* Multiple messages received in a single response from the server will now be split into individual messages
* Non 200 response will now be bubbled to the client

###Quick Start Video
We've put together a quick HOWTO video here https://vimeo.com/66431136

###Build Instructions Summary
* Built using Eclipse IDE (juno) 
* Install golang plugin for Eclipse
* Use the command go get github.com/pubnub/go to download and install the package
* Run the project
* Look for the application in the "Console" of the Eclipse IDE

In addition to Eclipse, this has also been tested with Go 1.0.3 on Linux using IntelliJ IDEA 12.

###Demo Console App
We've included a demo console app which documents all the functionality of the client, for example:

* Subscribe
* Subscribe with timetoken
* Publish
* Presence
* Detailed History
* Here_Now
* Unsubscribe
* Presence-Unsubscribe
* Time
* Disconnect/Retry
* Exit

###Quick Implementation Examples
* Init
```
        pubInstance := pubnubMessaging.PubnubInit(<YOUR PUBLISH KEY>, <YOUR SUBSCRIBE KEY>, <SECRET KEY>, <CIPHER>, <SSL ON/OFF>, <UUID>)
```

* Publish

```
        //Init pubnub instance

        var errorChannel = make(chan []byte)
        var callbackChannel = make(chan []byte)
        go pubInstance.Publish(<pubnub channel>, <message to publish>, callbackChannel, errorChannel)
        go ParseResponse(callbackChannel)
        go ParseErrorResponse(errorChannel) 
        // please goto the end of this file see the implementations of ParseResponse and ParseErrorResponse
```

* Subsribe

```
        //Init pubnub instance

        var errorChannel = make(chan []byte)
        var subscribeChannel = make(chan []byte)
        go pubInstance.Subscribe(<pubnub channels, multiple channels can be separated by comma>, <timetoken, should be an empty string in this case>, subscribeChannel, <this field is FALSE for subscribe requests>, errorChannel)
        go ParseResponse(subscribeChannel)  
        go ParseErrorResponse(errorChannel)  
        // please goto the end of this file see the implementations of ParseResponse and ParseErrorResponse
```

* Subscribe with timetoken

```
        //Init pubnub instance

        var errorChannel = make(chan []byte)
        var subscribeChannel = make(chan []byte)
        go pubInstance.Subscribe(<pubnub channel, multiple channels can be separated by comma>, <timetoken to init the request with>, subscribeChannel, <this field is FALSE for subscribe requests>, errorChannel)
        go ParseResponse(subscribeChannel)  
        go ParseErrorResponse(errorChannel)  
        // please goto the end of this file see the implementations of ParseResponse and ParseErrorResponse
```

* Presence
```
        //Init pubnub instance

        var errorChannel = make(chan []byte)
        var presenceChannel = make(chan []byte)
        go pubInstance.Subscribe(<pubnub channel, multiple channels can be separated by comma>, <timetoken, should be an empty string in this case>, presenceChannel, <this field is TRUE for subscribe requests>, errorChannel)
        go ParseResponse(subscribeChannel)  
        go ParseErrorResponse(errorChannel)  
        // please goto the end of this file see the implementations of ParseResponse and ParseErrorResponse
```

* Detailed History
```
        //Init pubnub instance

        var errorChannel = make(chan []byte)
        var channelCallback = make(chan []byte)
        go pubInstance.History(<pubnub channel>, <no of items to fetch>, <start time>, <end time>, false, channelCallback, errorChannel)
        //example: go _pub.History(<pubnub channel>, 100, 0, 0, false, channelCallback, errorChannel)
        go ParseResponse(channel)
        go ParseErrorResponse(errorChannel)  
        // please goto the end of this file see the implementations of ParseResponse and ParseErrorResponse
```

* Here_Now
```
        //Init pubnub instance

        var errorChannel = make(chan []byte)
        var channelCallback = make(chan []byte)
        go pubInstance.HereNow(<pubnub channel>, channelCallback, errorChannel)
        go ParseResponse(channelCallback)
        go ParseErrorResponse(errorChannel) 
        // please goto the end of this file see the implementations of ParseResponse and ParseErrorResponse
```

* Unsubscribe
```
        //Init pubnub instance

        var errorChannel = make(chan []byte)
        var channelCallback = make(chan []byte)
        go pubInstance.Unsubscribe(<pubnub channels, multiple channels can be separated by comma>, channelCallback, errorChannel)
        go ParseUnsubResponse(channelCallback)
        go ParseErrorResponse(errorChannel) 
        // please goto the end of this file see the implementations of ParseResponse and ParseErrorResponse
```

* Presence-Unsubscribe
```
        //Init pubnub instance

        var errorChannel = make(chan []byte)
        var channelCallback = make(chan []byte)
        go pubInstance.PresenceUnsubscribe(<pubnub channels, multiple channels can be separated by comma>, channelCallback, errorChannel)
        go ParseUnsubResponse(channelCallback)
        go ParseErrorResponse(errorChannel) 
        // please goto the end of this file see the implementations of ParseResponse and ParseErrorResponse
```

* Time

```
        //Init pubnub instance

        var errorChannel = make(chan []byte)
        var channelCallback = make(chan []byte)
        go pubInstance.GetTime(channelCallback, errorChannel)
        go ParseResponse(channelCallback)
        go ParseErrorResponse(errorChannel) 
        // please goto the end of this file see the implementations of ParseResponse and ParseErrorResponse
```

* Disconnect/Retry
```
        //Init pubnub instance

        pubInstance.CloseExistingConnection() 
```

* Exit
```
        //Init pubnub instance

        pubInstance.Abort()  
```

* ParseResponse
```
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
```

* ParseErrorResponse
```
func ParseErrorResponse(channel chan []byte){
    for {
        value, ok := <-channel
        if !ok {  
            fmt.Println("")            
            break
        }
        if string(value) != "[]"{
            if(_displayError){
                fmt.Println(fmt.Sprintf("Error Callback: %s", value))
                fmt.Println("")
            }
        }
    }
}
```
