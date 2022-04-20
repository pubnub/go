package e2e

import (
	"context"
	"fmt"
	"testing"

	pubnub "github.com/pubnub/go/v7"
	"github.com/pubnub/go/v7/tests/stubs"
	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	assert := assert.New(t)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/time/0",
		Query:              "",
		ResponseBody:       `[15078947309567840]`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(configCopy())
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.Time().Execute()

	assert.Nil(err)

	assert.True(int64(15059085932399340) < res.Timetoken)
}

func TestTimeContext(t *testing.T) {
	assert := assert.New(t)

	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/time/0",
		Query:              "",
		ResponseBody:       `[15078947309567840]`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 200,
	})

	pn := pubnub.NewPubNub(configCopy())
	pn.SetClient(interceptor.GetClient())

	res, _, err := pn.TimeWithContext(backgroundContext).Execute()

	assert.Nil(err)

	assert.True(int64(15059085932399340) < res.Timetoken)
}

// Message represents the serializable format we use to send a protobuf message over pubnub
type Message struct {
	Type string `json:"type" mapstructure:"type"`
	Data []byte `json:"data" mapstructure:"data"`
}

type Service struct {
	pubnubClient *pubnub.PubNub
}

// CreateMessage will create a serializable message from any proto message
func CreateMessage() (*Message, error) {
	bytes := []byte("string to byte array or slice")
	return &Message{Type: string("This"), Data: bytes}, nil
}

// SendMessageToChannel will send the given proto message to the specified pubnub channel
func (s *Service) sendMessageToChannel(ch string) error {
	message, err := CreateMessage()
	if err != nil {
		return fmt.Errorf("error marshalling message: %w", err)
	}
	r, _, err := s.pubnubClient.Publish().Channel(ch).Message(message).Serialize(true).UsePost(true).Execute()
	if err != nil {
		return fmt.Errorf("error publishing message: %w", err)
	}
	println(r.Timestamp)
	return nil
}

// SendMessageToChannel will send the given proto message to the specified pubnub channel
func (s *Service) sendMessageToChannel2(ctx context.Context, ch string) error {
	message, err := CreateMessage()
	if err != nil {
		return fmt.Errorf("error marshalling message: %w", err)
	}
	r, _, err := s.pubnubClient.PublishWithContext(ctx).Channel(ch).Message(message).Serialize(true).UsePost(true).Execute()
	if err != nil {
		return fmt.Errorf("error publishing message: %w", err)
	}

	println(r.Timestamp)
	return nil
}

func TestWithoutContext(t *testing.T) {

	pn := pubnub.NewPubNub(configCopy())

	service := Service{
		pubnubClient: pn,
	}

	err := service.sendMessageToChannel("ch1")

	println(err)
}

func TestWithContext(t *testing.T) {

	pn := pubnub.NewPubNub(configCopy())

	service := Service{
		pubnubClient: pn,
	}

	err := service.sendMessageToChannel2(context.Background(), "ch1")

	println(err)
}
