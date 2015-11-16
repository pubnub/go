package messaging

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// AUDIT CHANNELS
func TestPamChGenerateParamsStringAudit(t *testing.T) {
	params := pubnub.pamGenerateParamsForChannel("audit", "testc", true,
		false, 4, "blah")

	assert.Equal(t, params,
		fmt.Sprintf("auth=blah&channel=testc&%s&timestamp=%s&uuid=%s",
			sdkIdentificationParam, timestamp(), pubnub.GetUUID()), "should be equal")
}

func TestPamChGenerateParamsStringAuditNoAuth(t *testing.T) {
	params := pubnub.pamGenerateParamsForChannel("audit", "testc", true, false,
		4, "")

	assert.Equal(t, params, fmt.Sprintf("channel=testc&%s&timestamp=%s&uuid=%s",
		sdkIdentificationParam, timestamp(), pubnub.GetUUID()), "should be equal")
}

func TestPamChGenerateParamsStringAuditNoChannel(t *testing.T) {
	params := pubnub.pamGenerateParamsForChannel("audit", "", true, false, 4,
		"blah")

	assert.Equal(t, params, fmt.Sprintf("auth=blah&%s&timestamp=%s&uuid=%s",
		sdkIdentificationParam, timestamp(), pubnub.GetUUID()), "should be equal")
}

func TestPamChGenerateParamsStringAuditNoAuthNoChannel(t *testing.T) {
	params := pubnub.pamGenerateParamsForChannel("audit", "", true, false, 4, "")

	assert.Equal(t, params, fmt.Sprintf("%s&timestamp=%s&uuid=%s",
		sdkIdentificationParam, timestamp(), pubnub.GetUUID()), "should be equal")
}

// GRANT CHANNELS
func TestPamChGenerateParamsStringGrant(t *testing.T) {
	params := pubnub.pamGenerateParamsForChannel("grant", "testc", true, false,
		4, "blah")

	assert.Equal(t, params,
		fmt.Sprintf("auth=blah&channel=testc&%s&r=1&timestamp=%s&ttl=4&uuid=%s&w=0",
			sdkIdentificationParam, timestamp(), pubnub.GetUUID()), "should be equal")
}

// no ttl
func TestPamChGenerateParamsStringGrantNoTTL(t *testing.T) {
	params := pubnub.pamGenerateParamsForChannel("grant", "testc", true, false,
		-1, "blah")

	assert.Equal(t, params,
		fmt.Sprintf("auth=blah&channel=testc&%s&r=1&timestamp=%s&uuid=%s&w=0",
			sdkIdentificationParam, timestamp(), pubnub.GetUUID()), "should be equal")
}

func TestPamChGenerateParamsStringGrantNoTTLZero(t *testing.T) {
	params := pubnub.pamGenerateParamsForChannel("grant", "testc", true, false,
		0, "blah")

	assert.Equal(t, params,
		fmt.Sprintf("auth=blah&channel=testc&%s&r=1&timestamp=%s&ttl=0&uuid=%s&w=0",
			sdkIdentificationParam, timestamp(), pubnub.GetUUID()), "should be equal")
}

// no auth
func TestPamChGenerateParamsStringGrantNoAuth(t *testing.T) {
	params := pubnub.pamGenerateParamsForChannel("grant", "testc", true, false,
		4, "")

	assert.Equal(t, params,
		fmt.Sprintf("channel=testc&%s&r=1&timestamp=%s&ttl=4&uuid=%s&w=0",
			sdkIdentificationParam, timestamp(), pubnub.GetUUID()), "should be equal")
}

// no channel
func TestPamChGenerateParamsStringGrantNoChannel(t *testing.T) {
	params := pubnub.pamGenerateParamsForChannel("grant", "", true, false, 4,
		"blah")

	assert.Equal(t, params,
		fmt.Sprintf("auth=blah&%s&r=1&timestamp=%s&ttl=4&uuid=%s&w=0",
			sdkIdentificationParam, timestamp(), pubnub.GetUUID()), "should be equal")
}

// no auth&channel
func TestPamChGenerateParamsStringGrantNoAuthNoChannel(t *testing.T) {
	params := pubnub.pamGenerateParamsForChannel("grant", "", true, true, 4, "")

	assert.Equal(t, params, fmt.Sprintf("%s&r=1&timestamp=%s&ttl=4&uuid=%s&w=1",
		sdkIdentificationParam, timestamp(), pubnub.GetUUID()), "should be equal")
}

// AUDIT CHANNEL GROUPS
func TestPamCgGenerateParamsStringAudit(t *testing.T) {
	params := pubnub.pamGenerateParamsForChannelGroup("audit", "testc", true,
		false, 4, "blah")

	assert.Equal(t, params,
		fmt.Sprintf("auth=blah&channel-group=testc&%s&timestamp=%s&uuid=%s",
			sdkIdentificationParam, timestamp(), pubnub.GetUUID()), "should be equal")
}

func TestPamCgGenerateParamsStringAuditNoAuth(t *testing.T) {
	params := pubnub.pamGenerateParamsForChannelGroup("audit", "testc", true, false,
		4, "")

	assert.Equal(t, params, fmt.Sprintf("channel-group=testc&%s&timestamp=%s&uuid=%s",
		sdkIdentificationParam, timestamp(), pubnub.GetUUID()), "should be equal")
}

func TestPamCgGenerateParamsStringAuditNoChannelGroup(t *testing.T) {
	params := pubnub.pamGenerateParamsForChannelGroup("audit", "", true, false, 4,
		"blah")

	assert.Equal(t, params, fmt.Sprintf("auth=blah&%s&timestamp=%s&uuid=%s",
		sdkIdentificationParam, timestamp(), pubnub.GetUUID()), "should be equal")
}

func TestPamCgGenerateParamsStringAuditNoAuthNoChannelGroup(t *testing.T) {
	params := pubnub.pamGenerateParamsForChannelGroup("audit", "", true, false, 4, "")

	assert.Equal(t, params, fmt.Sprintf("%s&timestamp=%s&uuid=%s",
		sdkIdentificationParam, timestamp(), pubnub.GetUUID()), "should be equal")
}

// GRANT CHANNEL GROUPS
func TestPamCgGenerateParamsStringGrant(t *testing.T) {
	params := pubnub.pamGenerateParamsForChannelGroup("grant", "testc", true, false,
		4, "blah")

	assert.Equal(t, params,
		fmt.Sprintf("auth=blah&channel-group=testc&m=0&%s&r=1&timestamp=%s&ttl=4&uuid=%s",
			sdkIdentificationParam, timestamp(), pubnub.GetUUID()), "should be equal")
}

// no ttl
func TestPamCgGenerateParamsStringGrantNoTTL(t *testing.T) {
	params := pubnub.pamGenerateParamsForChannelGroup("grant", "testc", true, false,
		-1, "blah")

	assert.Equal(t, params,
		fmt.Sprintf("auth=blah&channel-group=testc&m=0&%s&r=1&timestamp=%s&uuid=%s",
			sdkIdentificationParam, timestamp(), pubnub.GetUUID()), "should be equal")
}

func TestPamCgGenerateParamsStringGrantNoTTLZero(t *testing.T) {
	params := pubnub.pamGenerateParamsForChannelGroup("grant", "testc", true, false,
		0, "blah")

	assert.Equal(t, params,
		fmt.Sprintf("auth=blah&channel-group=testc&m=0&%s&r=1&timestamp=%s&ttl=0&uuid=%s",
			sdkIdentificationParam, timestamp(), pubnub.GetUUID()), "should be equal")
}

// no auth
func TestPamCgGenerateParamsStringGrantNoAuth(t *testing.T) {
	params := pubnub.pamGenerateParamsForChannelGroup("grant", "testc", true, false,
		4, "")

	assert.Equal(t, params,
		fmt.Sprintf("channel-group=testc&m=0&%s&r=1&timestamp=%s&ttl=4&uuid=%s",
			sdkIdentificationParam, timestamp(), pubnub.GetUUID()), "should be equal")
}

// no channel group
func TestPamCgGenerateParamsStringGrantNoChannelGroup(t *testing.T) {
	params := pubnub.pamGenerateParamsForChannelGroup("grant", "", true, false, 4,
		"blah")

	assert.Equal(t, params,
		fmt.Sprintf("auth=blah&m=0&%s&r=1&timestamp=%s&ttl=4&uuid=%s",
			sdkIdentificationParam, timestamp(), pubnub.GetUUID()), "should be equal")
}

// no auth&channel group
func TestPamCgGenerateParamsStringGrantNoAuthNoChannelGroup(t *testing.T) {
	params := pubnub.pamGenerateParamsForChannelGroup("grant", "", true, true, 4, "")

	assert.Equal(t, params, fmt.Sprintf("m=1&%s&r=1&timestamp=%s&ttl=4&uuid=%s",
		sdkIdentificationParam, timestamp(), pubnub.GetUUID()), "should be equal")
}

func timestamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}
