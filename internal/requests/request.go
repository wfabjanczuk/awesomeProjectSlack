package requests

import (
	"github.com/wfabjanczuk/awesomeProjectSlack/internal/connections"
)

type WSPayload struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

type WSRequest struct {
	Payload          WSPayload
	ClientConnection *connections.ClientConnection
}

func NewWSRequest(payload WSPayload, connection *connections.ClientConnection) *WSRequest {
	return &WSRequest{
		Payload:          payload,
		ClientConnection: connection,
	}
}
