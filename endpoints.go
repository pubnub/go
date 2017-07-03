package pubnub

import (
	"bytes"
	"net/http"
	"net/url"
)

type endpointOpts interface {
	config() Config
	client() *http.Client
	context() Context
	validate() error

	buildPath() (string, error)
	buildQuery() (*url.Values, error)
	// or bytes[]?
	buildBody() ([]byte, error)

	httpMethod() string
}

func defaultQuery(uuid string) *url.Values {
	v := &url.Values{}

	v.Set("pnsdk", "PubNub-Go/4.0.0")
	v.Set("uuid", uuid)

	return v
}

func buildUrl(o endpointOpts) (*url.URL, error) {
	var buffer bytes.Buffer

	path, err := o.buildPath()
	if err != nil {
		return &url.URL{}, err
	}

	query, err := o.buildQuery()
	if err != nil {
		return &url.URL{}, err
	}

	if o.config().Secure == true {
		buffer.WriteString("https")
	} else {
		buffer.WriteString("http")
	}

	retUrl := &url.URL{
		Path:     path,
		Scheme:   "https",
		Host:     o.config().Origin,
		RawQuery: query.Encode(),
	}

	buffer.WriteString("://")

	buffer.WriteString(o.config().Origin)
	buffer.WriteString(path)
	buffer.WriteString("?")

	buffer.WriteString(query.Encode())

	// return buffer.String(), nil
	return retUrl, nil
}
