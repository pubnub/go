package messaging

import (
	"encoding/json"
	"fmt"
)

type ResponseType int
type ErrorType int

const (
	ChannelResponse ResponseType = 1 << iota
	ChannelGroupResponse
	WildcardResponse
)

type SuccessResponse struct {
	Data      []byte
	Channel   string
	Source    string
	Timetoken string
	Presence  bool
	Type      ResponseType
}

func (r SuccessResponse) Bytes() []byte {
	// TODO: add cases for Wildcard responses
	switch r.Type {
	case ChannelGroupResponse:
		return []byte(fmt.Sprintf(
			"[[%s], \"%s\", \"%s\", \"%s\"]", r.Data, r.Timetoken,
			removePnpres(r.Channel), removePnpres(r.Source)))
	case ChannelResponse:
		fallthrough
	default:
		return []byte(fmt.Sprintf(
			"[[%s], \"%s\", \"%s\"]", r.Data, r.Timetoken, removePnpres(r.Channel)))
	}
}

type ServerSideErrorData struct {
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
	Error   bool        `json:"error"`
	Service string      `json:"service"`
	Status  int         `json:"status"`
}

type ErrorResponse interface {
	StringForSource(string) string
	BytesForSource(string) []byte
}

type ServerSideErrorResponse struct {
	ErrorResponse

	Data ServerSideErrorData
}

func (e ServerSideErrorResponse) StringForSource(source string) string {
	if val, err := json.Marshal(e.Data.Payload); err != nil || string(val) == "null" {
		return fmt.Sprintf("%s\n", e.Data.Message)
	} else {
		return fmt.Sprintf("%s(%d): %s\n", e.Data.Service, e.Data.Status, val)
	}
}

func (e ServerSideErrorResponse) BytesForSource(source string) []byte {
	return []byte(e.StringForSource(source))
}

func NewPlainServerSideErrorResponse(response interface{}, status int) *ServerSideErrorResponse {
	if responseString, err := json.Marshal(response); err != nil {
		return &ServerSideErrorResponse{
			Data: ServerSideErrorData{
				Message: "Error while marshalling error message",
				Status:  status,
			},
		}
	} else {
		return &ServerSideErrorResponse{
			Data: ServerSideErrorData{
				Message: string(responseString),
			},
		}
	}
}

type clientSideErrorResponse struct {
	ErrorResponse

	Message string
	Reason  ResponseStatus
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
