package pubnub

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestIAmHere(t *testing.T) {
	pubnub := NewPubNubDemo()
	pubnub.Config.Log = log.Default()

	resChan, errChan := pubnub.iAmHere(pubnub.ctx,
		[]string{"ch1"},
		[]string{})

	select {
	case info := <-errChan:
		fmt.Println("error", info)
	case <-resChan:
		fmt.Println("ok")
	}
}

func TestIAmAway(t *testing.T) {
	pubnub := NewPubNubDemo()
	pubnub.Config.Log = log.Default()

	resChan, errChan := pubnub.iAmAway(pubnub.ctx,
		[]string{"ch1"},
		[]string{})

	select {
	case info := <-errChan:
		fmt.Println("error", info)
	case <-resChan:
		fmt.Println("ok")
	}
}

func TestSetPresenceState(t *testing.T) {
	pubnub := NewPubNubDemo()
	pubnub.Config.Log = log.Default()

	resChan, errChan := pubnub.setPresenceState(pubnub.ctx,
		[]string{"ch1"},
		[]string{},
		map[string]interface{}{
			"t": "bla",
		},
	)

	select {
	case info := <-errChan:
		fmt.Println("error", info)
	case res := <-resChan:
		fmt.Println("result", res)
	}
}

func TestHandshake(t *testing.T) {
	pubnub := NewPubNubDemo()
	pubnub.Config.Log = log.Default()

	resChan, errChan := pubnub.handshake(pubnub.ctx,
		[]string{"ch1"},
		[]string{},
	)

	select {
	case info := <-errChan:
		fmt.Println("error", info)
	case res := <-resChan:
		fmt.Println("result", res)
	}
}

func TestReceiveMessages(t *testing.T) {
	pubnub := NewPubNubDemo()
	pubnub.Config.Log = log.Default()

	newCtx, _ := contextWithCancel(pubnub.ctx)
	newCtx2, cancel2 := context.WithTimeout(newCtx, time.Duration(2)*time.Second)

	resChan, errChan := pubnub.receiveMessages(newCtx2,
		[]string{"ch1"},
		[]string{},
		42,
		"42",
	)

	select {
	case info := <-errChan:
		fmt.Println("error", info)
	case res := <-resChan:
		fmt.Println("result", res)
	case <-time.NewTicker(time.Duration(1) * time.Second).C:
		fmt.Println("Done waiting. Cancel")
		cancel2()
	}

	select {
	case <-time.NewTicker(time.Duration(1) * time.Second).C:
		println("Line")
	}
}
