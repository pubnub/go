package e2e

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"time"

	pubnub "github.com/pubnub/go/v6"
	"github.com/stretchr/testify/assert"
)

var enableDebuggingInTests = false

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

func seedRand() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

func init() {
	seedRand()
	config = pubnub.NewConfig()
	config.PublishKey = "pub-c-3ed95c83-12e6-4cda-9d69-c47ba2abb57e"
	config.SubscribeKey = "sub-c-26a73b0a-c3f2-11e9-8b24-569e8a5c3af3"

	pamConfig = pubnub.NewConfig()
	pamConfig.PublishKey = "pub-c-cdea0ef1-c571-4b72-b43f-ff1dc8aa4c5d"
	pamConfig.SubscribeKey = "sub-c-4757f09c-c3f2-11e9-9d00-8a58a5558306"
	pamConfig.SecretKey = "sec-c-YTYxNzVjYzctNDY2MS00N2NmLTg2NjYtNGRlNWY1NjMxMDBm"

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

func logInTest(format string, a ...interface{}) (n int, err error) {
	if enableDebuggingInTests {
		return fmt.Printf(format, a...)
	}
	return 0, nil
}

func checkFor(assert *assert.Assertions, maxTime, intervalTime time.Duration, fun func() error) {
	maxTimeout := time.NewTimer(maxTime)
	interval := time.NewTicker(intervalTime)
	lastErr := fun()
	if lastErr == nil {
		return
	}
ForLoop:
	for {
		select {
		case <-interval.C:
			lastErr := fun()
			if lastErr != nil {
				logInTest("Error: %s. Checking in next %s\n", lastErr, intervalTime)
				continue
			} else {
				break ForLoop
			}
		case <-maxTimeout.C:
			assert.Fail(lastErr.Error())
			break ForLoop
		}
	}
}

func heyIterator(count int) <-chan string {
	channel := make(chan string)

	init := "hey-"

	go func() {
		for i := 1; i <= count; i++ {
			channel <- fmt.Sprintf("%s%d", init, i)
		}
	}()

	return channel
}
