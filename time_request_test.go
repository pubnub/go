package pubnub

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestTimeRequestHTTP2(t *testing.T) {
	assert := assert.New(t)

	config := NewConfig()
	config.Origin = "ssp.pubnub.com"
	config.UseHTTP2 = true
	pn := NewPubNub(config)

	_, s, err := pn.Time().Execute()

	assert.Nil(err)
	assert.Equal(200, s.StatusCode)
}

func TestNewTimeResponseUnmarshalling(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`s`)

	_, _, err := newTimeResponse(jsonBytes, fakeResponseState)
	assert.Equal("pubnub/parsing: Error unmarshalling response: {s}", err.Error())

	opts := &timeOpts{}
	a, err := opts.buildBody()
	assert.Nil(err)
	assert.Equal(a, []byte{})
}

func TestNewTimeBuilder(t *testing.T) {
	assert := assert.New(t)

	o := newTimeBuilder(pubnub)
	_, err := o.opts.buildBody()
	assert.Nil(err)
}

func TestNewTimeBuilderContext(t *testing.T) {
	assert := assert.New(t)

	o := newTimeBuilderWithContext(pubnub, context.Background())
	_, err := o.opts.buildBody()
	assert.Nil(err)
}
