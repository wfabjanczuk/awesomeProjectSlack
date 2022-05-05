package handlers

import (
	"github.com/gorilla/websocket"
	"github.com/wfabjanczuk/awesomeProjectSlack/internal/models"
	"net/http"
)

type WSConnection struct {
	connection *websocket.Conn
	channel    models.Channel
}

func (w *WSConnection) GetWSConnection() *websocket.Conn {
	return w.connection
}

func (w *WSConnection) GetChannel() models.Channel {
	return w.channel
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
