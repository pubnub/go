package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignSha256(t *testing.T) {
	assert := assert.New(t)

	signInput := "sub-c-7ba2ac4c-4836-11e6-85a4-0619f8945a4f\npub-c-98863562-19a6-4760-bf0b-d537d1f5c582\ngrant\nchannel=asyncio-pam-FI2FCS0A&pnsdk=PubNub-Python-Asyncio%252F4.0.0&r=1&timestamp=1468409553&uuid=a4dbf92e-e5cb-428f-b6e6-35cce03500a2&w=1"

	res := GetHmacSha256("my_key", signInput)

	assert.Equal("Dq92jnwRTCikdeP2nUs1__gyJthF8NChwbs5aYy2r_I=", res)
}
