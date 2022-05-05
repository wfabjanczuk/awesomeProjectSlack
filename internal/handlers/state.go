package handlers

import (
	"github.com/gorilla/websocket"
	"github.com/wfabjanczuk/awesomeProjectSlack/internal/connections"
	"github.com/wfabjanczuk/awesomeProjectSlack/internal/requests"
	"net/http"
)

const PublicChannel = "public"

var requestQueue = make(chan *requests.WSRequest)
var channels = make(map[string][]*connections.ClientConnection)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
