package pubnub

import (
	"bytes"
	"net/http"
	"net/url"
)

type endpointOpts interface {
	config() Config
	client() *http.Client
	validate() error

	buildPath() string
	buildQuery() *url.Values
	// or bytes[]?
	buildBody() string
}

func defaultQuery() *url.Values {
	v := &url.Values{}

	v.Set("pnsdk", "4")
	v.Set("uuid", "TODO-setup-uuid")

	return v
}

func buildUrl(o endpointOpts) string {
	var buffer bytes.Buffer

	if o.config().Secure == true {
		buffer.WriteString("https")
	} else {
		buffer.WriteString("http")
	}

	buffer.WriteString("://")

	buffer.WriteString(o.config().Origin)
	buffer.WriteString(o.buildPath())
	buffer.WriteString("?")

	buffer.WriteString(o.buildQuery().Encode())

	return buffer.String()
}
