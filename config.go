package pubnub

import (
	"fmt"
	"log"
	"sync"

	"github.com/pubnub/go/v5/utils"
)

const (
	presenceTimeout = 0
)

// Config instance is storage for user-provided information which describe further
// PubNub client behaviour. Configuration instance contain additional set of
// properties which allow to perform precise PubNub client configuration.
type Config struct {
	sync.RWMutex
	PublishKey                    string             // PublishKey you can get it from admin panel (only required if publishing).
	SubscribeKey                  string             // SubscribeKey you can get it from admin panel.
	SecretKey                     string             // SecretKey (only required for modifying/revealing access permissions).
	AuthKey                       string             // AuthKey If Access Manager is utilized, client will use this AuthKey in all restricted requests.
	Origin                        string             // Custom Origin if needed
	UUID                          string             // UUID to be used as a device identifier, a default uuid is generated if not passed.
	CipherKey                     string             // If CipherKey is passed, all communications to/from PubNub will be encrypted.
	Secure                        bool               // True to use TLS
	ConnectTimeout                int                // net.Dialer.Timeout
	NonSubscribeRequestTimeout    int                // http.Client.Timeout for non-subscribe requests
	SubscribeRequestTimeout       int                // http.Client.Timeout for subscribe requests only
	FileUploadRequestTimeout      int                // http.Client.Timeout File Upload Request only
	HeartbeatInterval             int                // The frequency of the pings to the server to state that the client is active
	PresenceTimeout               int                // The time after which the server will send a timeout for the client
	MaximumReconnectionRetries    int                // The config sets how many times to retry to reconnect before giving up.
	MaximumLatencyDataAge         int                // Max time to store the latency data for telemetry
	FilterExpression              string             // Feature to subscribe with a custom filter expression.
	PNReconnectionPolicy          ReconnectionPolicy // Reconnection policy selection
	Log                           *log.Logger        // Logger instance
	SuppressLeaveEvents           bool               // When true the SDK doesn't send out the leave requests.
	DisablePNOtherProcessing      bool               // PNOther processing looks for pn_other in the JSON on the recevied message
	UseHTTP2                      bool               // HTTP2 Flag
	MessageQueueOverflowCount     int                // When the limit is exceeded by the number of messages received in a single subscribe request, a status event PNRequestMessageCountExceededCategory is fired.
	MaxIdleConnsPerHost           int                // Used to set the value of HTTP Transport's MaxIdleConnsPerHost.
	MaxWorkers                    int                // Number of max workers for Publish and Grant requests
	UsePAMV3                      bool               // Use PAM version 2, Objects requets would still use PAM v3
	StoreTokensOnGrant            bool               // Will store grant v3 tokens in token manager for further use.
	FileMessagePublishRetryLimit  int                // The number of tries made in case of Publish File Message failure.
	UseRandomInitializationVector bool               // When true the IV will be random for all requests and not just file upload. When false the IV will be hardcoded for all requests except File Upload
}

// NewDemoConfig initiates the config with demo keys, for tests only.
func NewDemoConfig() *Config {
	demoConfig := NewConfig()

	demoConfig.PublishKey = "demo"
	demoConfig.SubscribeKey = "demo"
	demoConfig.SecretKey = "demo"

	return demoConfig

}

// NewConfig initiates the config with default values.
func NewConfig() *Config {
	c := Config{
		Origin:                        "ps.pndsn.com",
		Secure:                        true,
		ConnectTimeout:                10,
		NonSubscribeRequestTimeout:    10,
		SubscribeRequestTimeout:       310,
		FileUploadRequestTimeout:      60,
		MaximumLatencyDataAge:         60,
		MaximumReconnectionRetries:    50,
		SuppressLeaveEvents:           false,
		DisablePNOtherProcessing:      false,
		PNReconnectionPolicy:          PNNonePolicy,
		MessageQueueOverflowCount:     100,
		MaxIdleConnsPerHost:           30,
		MaxWorkers:                    20,
		UsePAMV3:                      true,
		StoreTokensOnGrant:            true,
		FileMessagePublishRetryLimit:  5,
		UseRandomInitializationVector: true,
	}

	c.UUID = fmt.Sprintf("pn-%s", utils.UUID())

	return &c
}

func (c *Config) checkMinTimeout(timeout int) int {
	if timeout < minTimeout {
		if c.Log != nil {
			c.Log.Println(fmt.Sprintf("PresenceTimeout value less than the min recommended value of %[1]d, setting value to %[1]d %d", minTimeout, timeout))
		}
		timeout = minTimeout
	}
	return timeout
}

// SetPresenceTimeoutWithCustomInterval sets the presence timeout and interval.
// timeout: How long the server will consider the client alive for presence.
// interval: How often the client will announce itself to server.
func (c *Config) SetPresenceTimeoutWithCustomInterval(
	timeout, interval int) *Config {
	timeout = c.checkMinTimeout(timeout)
	c.Lock()
	c.PresenceTimeout = timeout
	c.HeartbeatInterval = interval
	c.Unlock()
	return c
}

var minTimeout = 20

// SetPresenceTimeout sets the presence timeout and automatically calulates the preferred timeout value.
// timeout: How long the server will consider the client alive for presence.
func (c *Config) SetPresenceTimeout(timeout int) *Config {
	timeout = c.checkMinTimeout(timeout)
	return c.SetPresenceTimeoutWithCustomInterval(timeout, (timeout/2)-1)
}
