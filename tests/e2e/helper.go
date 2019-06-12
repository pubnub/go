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
	config.PublishKey = "pub-c-c1648ded-d156-4a2d-9dbb-a23262945fe2"
	config.SubscribeKey = "sub-c-c14b8948-7dfe-11e9-aee4-2e27e4d79cf8"

	pamConfig = pubnub.NewConfig()
	pamConfig.PublishKey = "pub-c-2dea72e4-e0aa-4c85-9411-d75baf7568b7"
	pamConfig.SubscribeKey = "sub-c-490a8ac8-7e0e-11e9-84e9-eed29b7b36d8"
	pamConfig.SecretKey = "sec-c-MDU3OGY1ZjMtMDUwZS00NTc4LWFhM2ItN2E3NzhmMDVkZmQx"

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
