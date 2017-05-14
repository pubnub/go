package pubnub

import (
	"bytes"
	"net/url"
)

type Endpoint interface {
	buildPath() string
	buildQuery() *url.Values
	// or bytes[]?
	buildBody() string
	PubNub() *PubNub
}

type TransactionalEndpoint interface {
	Sync() (interface{}, error)
	Async()
}

func defaultQuery() *url.Values {
	v := &url.Values{}

	v.Set("pnsdk", "4")
	v.Set("uuid", "TODO-setup-uuid")

	return v
}

func buildUrl(e Endpoint) string {
	var buffer bytes.Buffer

	if e.PubNub().PNConfig.Secure == true {
		buffer.WriteString("https")
	} else {
		buffer.WriteString("http")
	}

	buffer.WriteString("://")

	buffer.WriteString(e.PubNub().PNConfig.Origin)
	buffer.WriteString(e.buildPath())
	buffer.WriteString("?")

	buffer.WriteString(e.buildQuery().Encode())

	return buffer.String()
}
