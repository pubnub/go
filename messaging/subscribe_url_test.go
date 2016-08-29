package messaging

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	channelsSingleChannel.Add("blah", successChannel, errorChannel)

	channelsThreeChannels.Add("qwer", successChannel, errorChannel)
	channelsThreeChannels.Add("asdf", successChannel, errorChannel)
	channelsThreeChannels.Add("zxcv", successChannel, errorChannel)

	channelsSingleCG.Add("qwer", successChannel, errorChannel)

	channelsThreeCG.Add("qwer", successChannel, errorChannel)
	channelsThreeCG.Add("asdf", successChannel, errorChannel)
	channelsThreeCG.Add("zxcv", successChannel, errorChannel)

	channelsChannelAndGroupC.Add("asdf", successChannel, errorChannel)
	channelsChannelAndGroupG.Add("qwer", successChannel, errorChannel)
}

func TestCreateSubscribeURLWithoutTimetoken(t *testing.T) {
	initLogging()
	url, timetoken := pubnubSingleChannel.createSubscribeURL("0", "")

	assert.Equal(t, url,
		fmt.Sprintf("/v2/subscribe/my_key/blah/0?uuid=%s&tt=0&%s",
			pubnubSingleChannel.GetUUID(), sdkIdentificationParam), "should be equal")

	assert.Equal(t, "0", timetoken)
}

func TestCreateSubscribeURLWithTimetoken(t *testing.T) {
	initLogging()
	url, timetoken := pubnubSingleChannel.createSubscribeURL("123456", "")

	assert.Equal(t, url,
		fmt.Sprintf("/v2/subscribe/my_key/blah/0?uuid=%s&tt=0&%s",
			pubnubSingleChannel.GetUUID(), sdkIdentificationParam), "should be equal")

	assert.Equal(t, "123456", timetoken)
}

func TestCreateSubscribeURLMultipleChannels(t *testing.T) {
	initLogging()
	url, _ := pubnubThreeChannels.createSubscribeURL("0", "")

	assert.Contains(t, url, fmt.Sprintf("/v2/subscribe/my_key/"), "should be equal")
	assert.Contains(t, url, fmt.Sprintf("/0?uuid=%s&tt=0&%s",
		pubnubThreeChannels.GetUUID(), sdkIdentificationParam), "should be equal")

	assert.Contains(t, url, "asdf", "should be equal")
	assert.Contains(t, url, "qwer", "should be equal")
	assert.Contains(t, url, "zxcv", "should be equal")
}

func TestCreateSubscribeURLSingleCG(t *testing.T) {
	initLogging()
	url, _ := pubnubSingleCG.createSubscribeURL("0", "")

	assert.Equal(t, url,
		fmt.Sprintf("/v2/subscribe/my_key/,/0?channel-group=qwer&uuid=%s&tt=0&%s",
			pubnubSingleCG.GetUUID(), sdkIdentificationParam), "should be equal")
}

func TestCreateSubscribeURLMultipleCG(t *testing.T) {
	initLogging()
	url, _ := pubnubThreeCG.createSubscribeURL("0", "")

	assert.Contains(t, url, "/v2/subscribe/my_key/,/0?channel-group=",
		"should be equal")

	assert.Contains(t, url, fmt.Sprintf("&uuid=%s&tt=0&%s",
		pubnubThreeCG.GetUUID(), sdkIdentificationParam), "should be equal")

	assert.Contains(t, url, "asdf", "should be equal")
	assert.Contains(t, url, "qwer", "should be equal")
	assert.Contains(t, url, "zxcv", "should be equal")
}

func TestCreateSubscribeURLChannelAndCG(t *testing.T) {
	initLogging()
	url, _ := pubnubChannelAndGroup.createSubscribeURL("0", "")

	assert.Equal(t, url,
		fmt.Sprintf("/v2/subscribe/my_key/asdf/0?channel-group=qwer&uuid=%s&tt=0&%s",
			pubnubChannelAndGroup.GetUUID(), sdkIdentificationParam), "should be equal")
}
