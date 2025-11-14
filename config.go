package pubnub

import (
	"fmt"
	"log"
	"sync"

	"github.com/pubnub/go/v8/crypto"
)

const (
	presenceTimeout = 0
)

type UserId string

// Config instance is storage for user-provided information which describe further
// PubNub client behaviour. Configuration instance contain additional set of
// properties which allow to perform precise PubNub client configuration.
type Config struct {
	sync.RWMutex
	PublishKey   string // PublishKey you can get it from admin panel (only required if publishing).
	SubscribeKey string // SubscribeKey you can get it from admin panel.
	SecretKey    string // SecretKey (only required for modifying/revealing access permissions).
	AuthKey      string // AuthKey If Access Manager is utilized, client will use this AuthKey in all restricted requests.
	Origin       string // Custom Origin if needed

	// UUID to be used as a device identifier.
	//
	//Deprecated: please use SetUserId/GetUserId
	UUID string
	//DEPRECATED: please use CryptoModule
	CipherKey                    string             // If CipherKey is passed, all communications to/from PubNub will be encrypted.
	Secure                       bool               // True to use TLS
	ConnectTimeout               int                // net.Dialer.Timeout
	NonSubscribeRequestTimeout   int                // http.Client.Timeout for non-subscribe requests
	SubscribeRequestTimeout      int                // http.Client.Timeout for subscribe requests only
	FileUploadRequestTimeout     int                // http.Client.Timeout File Upload Request only
	HeartbeatInterval            int                // The frequency of the pings to the server to state that the client is active
	PresenceTimeout              int                // The time after which the server will send a timeout for the client
	MaximumReconnectionRetries   int                // The config sets how many times to retry to reconnect before giving up.
	MaximumLatencyDataAge        int                // Max time to store the latency data for telemetry
	FilterExpression             string             // Feature to subscribe with a custom filter expression.
	PNReconnectionPolicy         ReconnectionPolicy // Reconnection policy selection
	Log                          *log.Logger        // Deprecated: Logger instance. Use Loggers instead for enhanced logging.
	Loggers                      []PNLogger         // Custom loggers for enhanced logging. If empty, no logging will occur unless the deprecated Log field is set.
	SuppressLeaveEvents          bool               // When true the SDK doesn't send out the leave requests.
	DisablePNOtherProcessing     bool               // PNOther processing looks for pn_other in the JSON on the recevied message
	UseHTTP2                     bool               // HTTP2 Flag
	MessageQueueOverflowCount    int                // When the limit is exceeded by the number of messages received in a single subscribe request, a status event PNRequestMessageCountExceededCategory is fired.
	MaxIdleConnsPerHost          int                // Used to set the value of HTTP Transport's MaxIdleConnsPerHost.
	MaxWorkers                   int                // Number of max workers for Publish and Grant requests
	UsePAMV3                     bool               // Use PAM version 2, Objects requets would still use PAM v3
	StoreTokensOnGrant           bool               // Will store grant v3 tokens in token manager for further use.
	FileMessagePublishRetryLimit int                // The number of tries made in case of Publish File Message failure.
	//DEPRECATED: please use CryptoModule
	UseRandomInitializationVector bool                // When true the IV will be random for all requests and not just file upload. When false the IV will be hardcoded for all requests except File Upload
	CryptoModule                  crypto.CryptoModule // A cryptography module used for encryption and decryption

	validationWarnings []string // Internal field to store validation warnings during config setup
}

// NewDemoConfig initiates the config with demo keys, for tests only.
func NewDemoConfig() *Config {
	demoConfig := NewConfigWithUserId(UserId(GenerateUUID()))

	demoConfig.PublishKey = "demo"
	demoConfig.SubscribeKey = "demo"
	demoConfig.SecretKey = "demo"

	return demoConfig

}

func NewConfigWithUserId(userId UserId) *Config {
	c := Config{
		UUID:                          string(userId),
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

	return &c
}

// Deprecated: Please use NewConfigWithUserId
func NewConfig(uuid string) *Config {
	return NewConfigWithUserId(UserId(uuid))
}

func (c *Config) checkMinTimeout(timeout int) int {
	if timeout < minTimeout {
		warning := fmt.Sprintf("PresenceTimeout value %d is less than the min recommended value of %d, adjusting to %d", timeout, minTimeout, minTimeout)

		// Store warning to be logged when PubNub instance is created
		c.validationWarnings = append(c.validationWarnings, warning)

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

// SetUserId sets userId
func (c *Config) SetUserId(userId UserId) *Config {
	c.UUID = string(userId)
	return c
}

// GetUserId gets value of userId
func (c *Config) GetUserId() UserId {
	return UserId(c.UUID)
}

// GetLogString returns a formatted string representation of the Config for logging purposes.
// Sensitive fields (SecretKey, CipherKey) are masked with ***.
func (c *Config) GetLogString() string {
	c.RLock()
	defer c.RUnlock()

	maskIfNotEmpty := func(value string) string {
		if value != "" {
			return "***"
		}
		return ""
	}

	cryptoModuleStr := "<nil>"
	if c.CryptoModule != nil {
		cryptoModuleStr = "<configured>"
	}

	loggersStr := fmt.Sprintf("%d logger(s)", len(c.Loggers))

	return fmt.Sprintf(`Config{
  PublishKey: %s
  SubscribeKey: %s
  SecretKey: %s
  AuthKey: %s
  Origin: %s
  UUID: %s
  CipherKey: %s
  Secure: %t
  ConnectTimeout: %d
  NonSubscribeRequestTimeout: %d
  SubscribeRequestTimeout: %d
  FileUploadRequestTimeout: %d
  HeartbeatInterval: %d
  PresenceTimeout: %d
  MaximumReconnectionRetries: %d
  MaximumLatencyDataAge: %d
  FilterExpression: %s
  PNReconnectionPolicy: %s
  SuppressLeaveEvents: %t
  DisablePNOtherProcessing: %t
  UseHTTP2: %t
  MessageQueueOverflowCount: %d
  MaxIdleConnsPerHost: %d
  MaxWorkers: %d
  UsePAMV3: %t
  StoreTokensOnGrant: %t
  FileMessagePublishRetryLimit: %d
  UseRandomInitializationVector: %t
  CryptoModule: %s
  Loggers: %s
}`,
		c.PublishKey,
		c.SubscribeKey,
		maskIfNotEmpty(c.SecretKey),
		maskIfNotEmpty(c.AuthKey),
		c.Origin,
		c.UUID,
		maskIfNotEmpty(c.CipherKey),
		c.Secure,
		c.ConnectTimeout,
		c.NonSubscribeRequestTimeout,
		c.SubscribeRequestTimeout,
		c.FileUploadRequestTimeout,
		c.HeartbeatInterval,
		c.PresenceTimeout,
		c.MaximumReconnectionRetries,
		c.MaximumLatencyDataAge,
		c.FilterExpression,
		c.PNReconnectionPolicy,
		c.SuppressLeaveEvents,
		c.DisablePNOtherProcessing,
		c.UseHTTP2,
		c.MessageQueueOverflowCount,
		c.MaxIdleConnsPerHost,
		c.MaxWorkers,
		c.UsePAMV3,
		c.StoreTokensOnGrant,
		c.FileMessagePublishRetryLimit,
		c.UseRandomInitializationVector,
		cryptoModuleStr,
		loggersStr,
	)
}
