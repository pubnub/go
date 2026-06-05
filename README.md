# PubNub Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/pubnub/go/v9.svg)](https://pkg.go.dev/github.com/pubnub/go/v9)
[![codecov.io](https://codecov.io/github/pubnub/go/coverage.svg)](https://codecov.io/github/pubnub/go)
[![Go Report Card](https://goreportcard.com/badge/github.com/pubnub/go/v9)](https://goreportcard.com/report/github.com/pubnub/go/v9)

This is the official PubNub Go SDK repository.

PubNub takes care of the infrastructure and APIs needed for the realtime communication layer of your application. Work on your app's logic and let PubNub handle sending and receiving data across the world in less than 100ms.

## Requirements

* Go 1.25 or later

## Get keys

You will need the publish and subscribe keys to authenticate your app. Get your keys from the [Admin Portal](https://dashboard.pubnub.com/).

## Configure PubNub

1. Integrate PubNub into your project:

    ```bash
    go get github.com/pubnub/go/v9
    ```

2. Create a new file and add the following code:

    ```go
    package main

    import pubnub "github.com/pubnub/go/v9"

    func main() {
        config := pubnub.NewConfigWithUserId("userId")
        config.SubscribeKey = "mySubscribeKey"
        config.PublishKey = "myPublishKey"

        pn := pubnub.NewPubNub(config)
    }
    ```

## Add event listeners

```go
import (
    "fmt"

    pubnub "github.com/pubnub/go/v9"
)

listener := pubnub.NewListener()
doneConnect := make(chan bool)
donePublish := make(chan bool)

msg := map[string]interface{}{
    "msg": "Hello world",
}
go func() {
    for {
        select {
        case status := <-listener.Status:
            switch status.Category {
            case pubnub.PNDisconnectedCategory:
                // This event happens when all channels/groups have been unsubscribed (graceful disconnect)
            case pubnub.PNDisconnectedUnexpectedlyCategory:
                // This event happens when radio / connectivity is lost
            case pubnub.PNConnectedCategory:
                // Connect event. You can do stuff like publish, and know you'll get it.
                // Or just use the connected event to confirm you are subscribed for
                // UI / internal notifications, etc
                doneConnect <- true
            case pubnub.PNReconnectedCategory:
                // Happens as part of our regular operation. This event happens when
                // radio / connectivity is lost, then regained.
            }
        case message := <-listener.Message:
            // Handle new message stored in message.message
            if message.Channel != "" {
                // Message has been received on channel group stored in
                // message.Channel
            } else {
                // Message has been received on channel stored in
                // message.Subscription
            }
            if msg, ok := message.Message.(map[string]interface{}); ok {
                fmt.Println(msg["msg"])
            }
            /*
                log the following items with your favorite logger
                    - message.Message
                    - message.Subscription
                    - message.Timetoken
            */

            donePublish <- true
        case <-listener.Presence:
            // handle presence
        }
    }
}()

pn.AddListener(listener)
```

## Publish and subscribe

In this code, publishing a message is triggered when the subscribe call is finished successfully. The `Publish()` method uses the `msg` variable that is used in the following code.

```go
msg := map[string]interface{}{
        "msg": "Hello world!"
}

pn.Subscribe().
    Channels([]string{"hello_world"}).
    Execute()

<-doneConnect

response, status, err := pn.Publish().
    Channel("hello_world").Message(msg).Execute()

if err != nil {
     // Request processing failed.
     // Handle message publish error
}
```

## Upgrading from v8

v9 requires Go 1.25+ and uses a new module import path:

```bash
go get github.com/pubnub/go/v9@latest
```

```go
import pubnub "github.com/pubnub/go/v9"
```

## Documentation

[PubNub Go SDK documentation](https://www.pubnub.com/docs/sdks/go)

## Support

If you **need help** or have a **general question**, contact <support@pubnub.com>.
