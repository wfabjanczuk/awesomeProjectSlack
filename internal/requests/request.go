package requests

import (
	"github.com/gorilla/websocket"
	"github.com/wfabjanczuk/awesomeProjectSlack/internal/models"
)

type ClientConnection interface {
	GetWSConnection() *websocket.Conn
	GetChannel() models.Channel
}

type WSPayload struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

type WSRequest struct {
	Payload    WSPayload
	Connection ClientConnection
}

func NewWSRequest(payload WSPayload, connection ClientConnection) *WSRequest {
	return &WSRequest{
		Payload:    payload,
		Connection: connection,
	}
}
