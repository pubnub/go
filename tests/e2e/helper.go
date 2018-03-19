package e2e

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"time"

	pubnub "github.com/pubnub/go"
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
	config.PublishKey = "pub-c-071e1a3f-607f-4351-bdd1-73a8eb21ba7c"
	config.SubscribeKey = "sub-c-5c4fdcc6-c040-11e5-a316-0619f8945a4f"

	pamConfig = pubnub.NewConfig()
	pamConfig.PublishKey = "pub-c-1bd448ed-05ba-4dbc-81a5-7d6ff5c6e2bb"
	pamConfig.SubscribeKey = "sub-c-90c51098-c040-11e5-a316-0619f8945a4f"
	pamConfig.SecretKey = "sec-c-ZDA1ZTdlNzAtYzU4Zi00MmEwLTljZmItM2ZhMDExZTE2ZmQ5"
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
