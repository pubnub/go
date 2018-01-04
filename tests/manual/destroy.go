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
	config.PublishKey = "pub-c-071e1a3f-607f-4351-bdd1-73a8eb21ba7c"
	config.SubscribeKey = "sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f"

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
