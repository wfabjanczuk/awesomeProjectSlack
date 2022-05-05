package handlers

import (
	"fmt"
	"github.com/wfabjanczuk/awesomeProjectSlack/internal/connections"
	"github.com/wfabjanczuk/awesomeProjectSlack/internal/requests"
	"log"
	"net/http"
)

func InitConnection(w http.ResponseWriter, r *http.Request) {
	wsConnection, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Client connected to the endpoint")

	clientConnection := connections.NewClientConnection(wsConnection, PublicChannel)
	channels[PublicChannel] = NewClientConnectionSet()
	channels[PublicChannel].Add(clientConnection)

	sendSuccessMessage(clientConnection, "Connected to the server")
	go listenOnClientConnection(clientConnection)
}

func listenOnClientConnection(clientConnection *connections.ClientConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	for {
		var payload requests.WSPayload

		err := clientConnection.GetWSConnection().ReadJSON(&payload)
		if err != nil {
			log.Println(err)
			sendErrorMessage(clientConnection, "Invalid JSON payload!")

			break
		} else {
			fmt.Printf("New message from client, action: \"%s\", message: \"%s\"\n", payload.Action, payload.Message)

			requestQueue <- requests.NewWSRequest(payload, clientConnection)
		}
	}
}

func ListenOnRequestQueue() {
	for {
		wsRequest := <-requestQueue
		payload := wsRequest.Payload

		switch payload.Action {
		case "broadcast":
			broadcastToChannel(wsRequest.ClientConnection.GetChannel(), payload.Message)
		case "create":
			createChannel(wsRequest.Payload.Message, wsRequest.ClientConnection)
		case "enter":
			enterChannel(wsRequest.Payload.Message, wsRequest.ClientConnection)
		case "leave":
			leaveChannel(wsRequest.ClientConnection)
		}
	}
}
