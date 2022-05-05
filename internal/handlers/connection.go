package handlers

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type WSConnection struct {
	connection *websocket.Conn
	channel    string
}

func (w *WSConnection) GetWSConnection() *websocket.Conn {
	return w.connection
}

func (w *WSConnection) GetChannel() string {
	return w.channel
}

func (w *WSConnection) SetChannel(channelName string) {
	w.channel = channelName
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
