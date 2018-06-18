package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	pubnub "github.com/pubnub/go"
)

func main() {
	config := pubnub.NewConfig()
	config.PublishKey = "pub-c-1bd448ed-05ba-4dbc-81a5-7d6ff5c6e2bb"
	config.SubscribeKey = "sub-c-b9ab9508-43cf-11e8-9967-869954283fb4"

	pn := pubnub.NewPubNub(config)

	fmt.Println("vim-go")

	// Add listeners
	ln1 := pubnub.NewListener()
	ln2 := pubnub.NewListener()
	ln3 := pubnub.NewListener()

	// TODO: listen on 1st listener

	pn.AddListener(ln1)

	pn.Subscribe().
		Channels([]string{"blah"}).
		Execute()

	pn.AddListener(ln2)
	pn.AddListener(ln3)

	time.Sleep(1 * time.Second)

	pn.Destroy()
	time.Sleep(1 * time.Second)

	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
}
