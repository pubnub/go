package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	pubnub "github.com/pubnub/go/v5"
)

func main() {
	config := pubnub.NewConfig()

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
