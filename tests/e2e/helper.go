package e2e

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"time"

	pubnub "github.com/zhashkevych/go"
)

const (
	SPECIAL_CHARACTERS = "-.,_~:/?#[]@!$&'()*+;=`|"
	SPECIAL_CHANNEL    = "-._~:/?#[]@!$&'()*+;=`|"
)

var pamConfig *pubnub.Config
var config *pubnub.Config

var (
	serverErrorTemplate     = "pubnub/server: Server respond with error code %d"
	validationErrorTemplate = "pubnub/validation: %s"
	connectionErrorTemplate = "pubnub/connection: %s"
)

func init() {
	rand.Seed(time.Now().UnixNano())

	config = pubnub.NewConfig()
	config.PublishKey = "pub-c-afeb2ec5-45e9-449f-9a8d-c4940a9c7836"
	config.SubscribeKey = "sub-c-e41d50d4-43ce-11e8-a433-9e6b275e7b64"

	pamConfig = pubnub.NewConfig()
	pamConfig.PublishKey = "pub-c-7e5c6521-91d0-4e60-9656-4bed419a769b"
	pamConfig.SubscribeKey = "sub-c-b9ab9508-43cf-11e8-9967-869954283fb4"
	pamConfig.SecretKey = "sec-c-MjRhODgwMTgtY2RmMS00ZWNmLTgzNTUtYjI3MzZhOThlNTY0"
}

func configCopy() *pubnub.Config {
	cfg := new(pubnub.Config)
	*cfg = *config
	return cfg
}

func pamConfigCopy() *pubnub.Config {
	config := new(pubnub.Config)
	*config = *pamConfig
	return config
}

func randomized(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, rand.Intn(10000000))
}

type fakeTransport struct {
	Status     string
	StatusCode int
	Body       io.ReadCloser
}

func (t fakeTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{
		Status:     t.Status,
		StatusCode: t.StatusCode,
		Body:       t.Body,
	}, nil
}

func (t fakeTransport) Dial(string, string) (net.Conn, error) {
	return nil, errors.New("ooops!")
}

func heyIterator(count int) <-chan string {
	channel := make(chan string)

	init := "hey-"

	go func() {
		for i := 1; i <= i; i++ {
			channel <- fmt.Sprintf("%s%d", init, i)
		}
	}()

	return channel
}
