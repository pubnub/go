package messaging

import (
	"encoding/json"
	"fmt"
)

type responseType int
type errorType int

const (
	channelResponse responseType = 1 << iota
	channelGroupResponse
	wildcardResponse
)

type successResponse struct {
	Data      []byte
	Channel   string
	Source    string
	Timetoken string
	Presence  bool
	Type      responseType
}

func (r successResponse) Bytes() []byte {
	switch r.Type {
	case wildcardResponse:
		fallthrough
	case channelGroupResponse:
		return []byte(fmt.Sprintf(
			"[[%s], \"%s\", \"%s\", \"%s\"]", r.Data, r.Timetoken,
			removePnpres(r.Channel), removePnpres(r.Source)))
	case channelResponse:
		fallthrough
	default:
		return []byte(fmt.Sprintf(
			"[[%s], \"%s\", \"%s\"]", r.Data, r.Timetoken, removePnpres(r.Channel)))
	}
}

type serverSideErrorData struct {
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
	Error   bool        `json:"error"`
	Service string      `json:"service"`
	Status  int         `json:"status"`
}

type errorResponse interface {
	StringForSource(string) string
	BytesForSource(string) []byte
}

type serverSideErrorResponse struct {
	errorResponse

	Data serverSideErrorData
}

func (e serverSideErrorResponse) StringForSource(source string) string {
	if val, err := json.Marshal(e.Data.Payload); err != nil || string(val) == "null" {
		return fmt.Sprintf("%s\n", e.Data.Message)
	} else {
		return fmt.Sprintf("%s(%d): %s\n", e.Data.Service, e.Data.Status, val)
	}
}

func (e serverSideErrorResponse) BytesForSource(source string) []byte {
	return []byte(e.StringForSource(source))
}

func newPlainServerSideErrorResponse(response interface{}, status int) *serverSideErrorResponse {
	if responseString, err := json.Marshal(response); err != nil {
		return &serverSideErrorResponse{
			Data: serverSideErrorData{
				Message: "Error while marshalling error message",
				Status:  status,
			},
		}
	} else {
		return &serverSideErrorResponse{
			Data: serverSideErrorData{
				Message: string(responseString),
			},
		}
	}
}

type clientSideErrorResponse struct {
	errorResponse

	Message string
	Reason  responseStatus
}

func newClientSideErrorResponse(msg string) *clientSideErrorResponse {
	return &clientSideErrorResponse{
		Message: msg,
	}
}

func (e clientSideErrorResponse) StringForSource(source string) string {
	// TODO: handle all reasons
	switch e.Reason {
	case responseAlreadySubscribed:
		return fmt.Sprintf("[0, \"%s channel '%s' %s\", \"%s\"]",
			stringPresenceOrSubscribe(source),
			source,
			stringResponseReason(e.Reason),
			source)
	case responseAsIsError:
		fallthrough
	default:
		return fmt.Sprintf("[0, \"%s\", \"%s\"]", e.Message, source)
	}
}

func (e clientSideErrorResponse) BytesForSource(source string) []byte {
	return []byte(e.StringForSource(source))
}

func (e clientSideErrorResponse) Bytes(source string) []byte {
	return []byte(e.StringForSource(source))
}
