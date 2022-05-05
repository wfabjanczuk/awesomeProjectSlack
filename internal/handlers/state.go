package handlers

import (
	"github.com/gorilla/websocket"
	"github.com/wfabjanczuk/awesomeProjectSlack/internal/connections"
	"github.com/wfabjanczuk/awesomeProjectSlack/internal/requests"
	"net/http"
)

type ClientConnectionSet map[*connections.ClientConnection]struct{}

func NewClientConnectionSet() ClientConnectionSet {
	return make(ClientConnectionSet)
}

func (s ClientConnectionSet) Add(clientConnection *connections.ClientConnection) {
	s[clientConnection] = struct{}{}
}

func (s ClientConnectionSet) Delete(clientConnection *connections.ClientConnection) {
	delete(s, clientConnection)
}

const PublicChannel = "public"

var requestQueue = make(chan *requests.WSRequest)
var channels = make(map[string]ClientConnectionSet)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
