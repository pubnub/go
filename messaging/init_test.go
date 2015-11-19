package messaging

import (
	"fmt"
	"regexp"
	"time"
)

var pubnub = Pubnub{
	subscribeKey: "demo",
}

var signatureRegexp, _ = regexp.Compile("&signature=.*$")

func timestamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

func truncateSignature(input string) (output string) {
	return signatureRegexp.ReplaceAllString(input, "")
}
