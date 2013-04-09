#PubNub clients for Go
## Preview Release -- Beta 1

####Build Instructions
* Built using Eclipse IDE (juno) 
* Install golang plugin for Eclipse
* Copy the src directory in the project
* Run the project
* Look for the application in the "Console" of the Eclipse IDE

This has also been tested with Go 1.0.3 on Linux using IntelliJ IDEA 12.

###Flow
* Runs as a console application
* Asks for the channel name (multiple channels separated by comma can be entered), SSL, Cipher and Custom UUID
####User then chooses the options to
* ENTER 1 FOR Subscribe
* ENTER 2 FOR Publish
* ENTER 3 FOR Presence
* ENTER 4 FOR Detailed History
* ENTER 5 FOR Here_Now
* ENTER 6 FOR Unsubscribe
* ENTER 7 FOR Presence-Unsubscribe
* ENTER 8 FOR Time
* ENTER 9 FOR Exit


###Features
* Runs on the console similar to the C#, mac and linux example
* Supports multiplexing
* Custom UUID is working
* SSL is working
* Cipher is working
* Subscribe and presence run in the background
* Other calls to pubnub service will also be async
* Naming convention consistency

###Known issues:
* When unsubscribed from a channel the web request doesn't abort at once. It will wait for response of the active request, parse the response and then close the connection for further requests 
* Presence and Subscribe messages appear twice in the console. 
* Optimize ParseJson method 
* Channel name is with "-pnpres" in case of presence response messages 
* Multiple request can be initiated for subscribe and presence. If the user choose 3 and then again chooses 3, same response will be displayed multiple times.

###Pending work: 
* Test cases
* Abort Http request 
* Proxy 
* Some additional settings like reconnect wait interval etc. 
* Add info messages like [1,"Connected","test"], [1,"Presence Connected","test"] for subscribe, presence, unsubscribe and unsubscribe-presence
* Notify the user when the channel is already subscribed (when subscribe request is received), or not subscribed (when unsubscribe request is received). Right now these are handled but the messages are not displayed to the user 
* Reconnect / no of retries on reconnect
