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
* Subscribe
* Publish
* Presence
* Detailed History
* Here_Now
* Unsubscribe
* Presence-Unsubscribe
* Time
* Exit

###Features
* Runs on the console similar to the C#, mac and linux example
* Supports multiplexing
* Custom UUID is working
* SSL is working
* Cipher is working
* Subscribe and presence run in the background
* Proxy is working for Basic Authentication with and without SSL
* Reconnect on internet disruption is working
* Timeouts can be set by changing the default values of the constant
* Retry on disconnect limit and interval can be set 
* Naming convention consistency
* Tested with Go 1.1
