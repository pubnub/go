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
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"

	pubnub "github.com/pubnub/go/v8"
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

// eventually polls check at an exponentially-growing interval (starting at initialInterval,
// capped at 2s) until it returns ok=true, or timeout elapses. On success it returns the
// converged value; on timeout it fails the test with a description that includes the last
// observed value, so CI failures stay diagnosable instead of producing a bare "expected X got Y".
//
// Use this for e2e assertions against eventually-consistent PubNub state (presence
// registration, message delivery, channel-group propagation, …) in place of fixed
// time.Sleep calls. Returning the converged value keeps callers free of outer
// state-capture variables.
//
// Example:
//
//	res := eventually(t, 30*time.Second, 250*time.Millisecond,
//	    fmt.Sprintf("channel %s reaches occupancy 3", ch),
//	    func() (*pubnub.HereNowResponse, bool) {
//	        r, _, err := pn.HereNow().Channels([]string{ch}).Execute()
//	        return r, err == nil && r != nil && r.TotalOccupancy >= 3
//	    })
func eventually[T any](
	t *testing.T,
	timeout, initialInterval time.Duration,
	description string,
	check func() (T, bool),
) T {
	t.Helper()
	const maxInterval = 2 * time.Second
	deadline := time.Now().Add(timeout)
	delay := initialInterval
	var last T
	for {
		v, ok := check()
		last = v
		if ok {
			return v
		}
		if time.Now().After(deadline) {
			t.Fatalf("eventually: %s did not converge within %v (last value=%#v)",
				description, timeout, last)
			return last
		}
		time.Sleep(delay)
		if delay < maxInterval {
			if delay *= 2; delay > maxInterval {
				delay = maxInterval
			}
		}
	}
}

// waitForOccupancy polls HereNow until channel reports at least expected total occupants,
// or timeout elapses. Returns the last HereNow response so callers can run further
// assertions on real data. Built on top of eventually.
func waitForOccupancy(t *testing.T, pn *pubnub.PubNub, channel string, expected int, timeout time.Duration) *pubnub.HereNowResponse {
	t.Helper()
	return eventually(t, timeout, 250*time.Millisecond,
		fmt.Sprintf("channel %s reaches occupancy %d", channel, expected),
		func() (*pubnub.HereNowResponse, bool) {
			r, _, err := pn.HereNow().
				Channels([]string{channel}).
				Limit(100).
				Offset(0).
				IncludeUUIDs(true).
				Execute()
			return r, err == nil && r != nil && r.TotalOccupancy >= expected
		},
	)
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
