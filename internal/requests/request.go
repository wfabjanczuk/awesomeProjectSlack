package requests

import (
	"github.com/gorilla/websocket"
)

type ClientConnection interface {
	GetWSConnection() *websocket.Conn
	GetChannel() string
	SetChannel(channel string)
}

type WSPayload struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

type WSRequest struct {
	Payload          WSPayload
	ClientConnection ClientConnection
}

func NewWSRequest(payload WSPayload, connection ClientConnection) *WSRequest {
	return &WSRequest{
		Payload:          payload,
		ClientConnection: connection,
	}
}
