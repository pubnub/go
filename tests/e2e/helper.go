package e2e

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	godotenv "github.com/joho/godotenv"
	pubnub "github.com/pubnub/go/v7"
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
	godotenv.Load("../../.env")

	seedRand()
	config = pubnub.NewConfigWithUserId(pubnub.UserId(pubnub.GenerateUUID()))
	config.PublishKey = os.Getenv("PUBLISH_KEY")
	config.SubscribeKey = os.Getenv("SUBSCRIBE_KEY")

	pamConfig = pubnub.NewConfigWithUserId(pubnub.UserId(pubnub.GenerateUUID()))
	pamConfig.PublishKey = os.Getenv("PAM_PUBLISH_KEY")
	pamConfig.SubscribeKey = os.Getenv("PAM_SUBSCRIBE_KEY")
	pamConfig.SecretKey = os.Getenv("PAM_SECRET_KEY")
}

func configCopy() *pubnub.Config {
	cfg := new(pubnub.Config)
	*cfg = *config
	cfg.SetUserId(pubnub.UserId(pubnub.GenerateUUID()))
	return cfg
}

func pamConfigCopy() *pubnub.Config {
	config := new(pubnub.Config)
	*config = *pamConfig
	config.SetUserId(pubnub.UserId(pubnub.GenerateUUID()))
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

func subscribeWithATimeout(t *testing.T, pn *pubnub.PubNub, channel string, duration time.Duration) error {
	listener := pubnub.NewListener()
	pn.AddListener(listener)
	pn.Subscribe().Channels([]string{channel}).Execute()
	timer := time.NewTimer(duration)
	select {
	case s := <-listener.Status:
		timer.Stop()
		if s.Category == pubnub.PNConnectedCategory {
			pn.RemoveListener(listener)
			return nil
		} else {
			errMsg := fmt.Sprintf("didn't receive connected but %s", s.Category)
			t.Error(errMsg)
			return errors.New(errMsg)
		}
	case <-timer.C:
		timer.Stop()
		errMsg := "connected didn't came in desired time"
		t.Error(errMsg)
		return errors.New(errMsg)
	}
}

func checkForAsserted(t *testing.T, maxTime, intervalTime time.Duration, fun func() error) {

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
