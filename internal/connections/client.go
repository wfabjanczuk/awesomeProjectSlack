package connections

import (
	"github.com/gorilla/websocket"
)

type ClientConnection struct {
	wsConnection *websocket.Conn
	channel      string
}

func NewClientConnection(wsConnection *websocket.Conn, channel string) *ClientConnection {
	return &ClientConnection{
		wsConnection: wsConnection,
		channel:      channel,
	}
}

func (c *ClientConnection) GetWSConnection() *websocket.Conn {
	return c.wsConnection
}

func (c *ClientConnection) GetChannel() string {
	return c.channel
}

func (c *ClientConnection) SetChannel(channelName string) {
	c.channel = channelName
}
