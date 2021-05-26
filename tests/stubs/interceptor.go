package stubs

import (
	"bytes"
	//"fmt"
	"io/ioutil"
	//"log"
	"net/http"
	"net/url"
	"time"

	"github.com/pubnub/go/v5/tests/helpers"
)

type Interceptor struct {
	Transport *interceptTransport
}

func NewInterceptor() *Interceptor {
	return &Interceptor{
		Transport: &interceptTransport{},
	}
}

func (i *Interceptor) AddStub(stub *Stub) {
	i.Transport.AddStub(stub)
}

func (i *Interceptor) GetClient() *http.Client {
	return &http.Client{
		Transport: i.Transport,
	}
}

type Stub struct {
	Method             string
	Path               string
	Query              string
	ResponseBody       string
	ResponseStatusCode int
	MixedPathPositions []int
	IgnoreQueryKeys    []string
	MixedQueryKeys     []string
	Hang               bool
}

func (s *Stub) Match(req *http.Request) bool {
	if s.Hang {
		time.Sleep(1000 * time.Second)
	}

	if s.Method != req.Method {
		//log.Printf("Methods are not equal: %s != %s\n", s.Method, req.Method)
		return false
	}

	parsedUrl, _ := req.URL.Parse(req.URL.String())
	if !helpers.PathsEqual(s.Path, parsedUrl.EscapedPath(), s.MixedPathPositions) {
		return false
	}

	expectedQuery, _ := url.ParseQuery(s.Query)
	actualQuery := req.URL.Query()

	//fmt.Println("ex", expectedQuery, "\nact", actualQuery, "\nignore",
	//	s.IgnoreQueryKeys)

	if !helpers.QueriesEqual(&expectedQuery,
		&actualQuery,
		s.IgnoreQueryKeys,
		s.MixedQueryKeys) {
		//fmt.Println("NOT EQUAL")
		return false

	}

	return true
}

type interceptTransport struct {
	Stubs []*Stub
}

func (i *interceptTransport) RoundTrip(req *http.Request) (*http.Response,
	error) {

	//fmt.Println(req.URL)

	for _, v := range i.Stubs {
		if v.Match(req) {
			var statusString string

			switch v.ResponseStatusCode {
			case 200:
				statusString = "200 OK"
			case 403:
				statusString = "403 Forbidden"
			default:
				statusString = ""
			}

			return &http.Response{
				Status:           statusString,
				StatusCode:       v.ResponseStatusCode,
				Proto:            "HTTP/1.0",
				ProtoMajor:       1,
				ProtoMinor:       0,
				Request:          req,
				Header:           http.Header{"Content-Length": {"256"}},
				TransferEncoding: nil,
				Close:            true,
				Body:             ioutil.NopCloser(bytes.NewBufferString(v.ResponseBody)),
				ContentLength:    256,
			}, nil
		}
	}

	// Nothing was found
	return &http.Response{
		Status:           "530 No stub matched",
		StatusCode:       530,
		Proto:            "HTTP/1.0",
		ProtoMajor:       1,
		ProtoMinor:       0,
		Request:          req,
		TransferEncoding: nil,
		Body:             ioutil.NopCloser(bytes.NewBufferString("No Stub Matched")),
		Close:            true,
		ContentLength:    256,
	}, nil
}

func (i *interceptTransport) AddStub(stub *Stub) {
	i.Stubs = append(i.Stubs, stub)
}
