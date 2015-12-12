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
	Error     bool
	Type      ResponseType
	Info      string
}

type ServerSideErrorData struct {
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
	Error   bool        `json:"error"`
	Service string      `json:"service"`
	Status  int         `json:"status"`
}

type ErrorResponse interface {
	error
}

type ServerSideErrorResponse struct {
	ErrorResponse

	Data ServerSideErrorData
}

func (e ServerSideErrorResponse) Error() string {
	if val, err := json.Marshal(e.Data.Payload); err != nil || string(val) == "null" {
		return fmt.Sprintf("%s\n", e.Data.Message)
	} else {
		return fmt.Sprintf("%s(%d): %s\n", e.Data.Service, e.Data.Status, val)
	}
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

type ClientSideErrorResponse struct {
	ErrorResponse

	Message string
	Reason  ResponseStatus
}

func (e ClientSideErrorResponse) Error() string {
	return fmt.Sprintf("Client-Side Error reason: %s %s",
		getResponseReasonString(e.Reason), e.Message)
}
