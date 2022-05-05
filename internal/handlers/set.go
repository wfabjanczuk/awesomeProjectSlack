package handlers

import "github.com/wfabjanczuk/awesomeProjectSlack/internal/connections"

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
