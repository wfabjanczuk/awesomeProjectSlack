package handlers

import (
	"fmt"
	"github.com/wfabjanczuk/awesomeProjectSlack/internal/connections"
	"github.com/wfabjanczuk/awesomeProjectSlack/internal/requests"
	"github.com/wfabjanczuk/awesomeProjectSlack/internal/responses"
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
	go ListenForWS(clientConnection)
}

func ListenForWS(clientConnection *connections.ClientConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload requests.WSPayload

	for {
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

func ListenToRequestQueue() {
	for {
		wsRequest := <-requestQueue
		channelName := wsRequest.ClientConnection.GetChannel()
		payload := wsRequest.Payload

		switch payload.Action {
		case "broadcast":
			broadcastToChannel(channelName, payload.Message)
		case "create":
			createChannel(wsRequest.Payload.Message, wsRequest.ClientConnection)
		case "enter":
			enterChannel(wsRequest.Payload.Message, wsRequest.ClientConnection)
		}
	}
}

func broadcastToChannel(channelName, message string) {
	for clientConnection := range channels[channelName] {
		sendSuccessMessage(clientConnection, message)
	}
}

func createChannel(channelName string, clientConnection *connections.ClientConnection) {
	if _, ok := channels[channelName]; ok {
		errorMessage := fmt.Sprintf("Channel with name: \"%s\" already exists!", channelName)
		sendErrorMessage(clientConnection, errorMessage)

		return
	}

	channels[channelName] = NewClientConnectionSet()

	successMessage := fmt.Sprintf("Channel with name: \"%s\" successfully created.", channelName)
	sendSuccessMessage(clientConnection, successMessage)
}

func enterChannel(channelName string, clientConnection *connections.ClientConnection) {
	if _, ok := channels[channelName]; !ok {
		errorMessage := fmt.Sprintf("Channel with name: \"%s\" does not exist!", channelName)
		sendErrorMessage(clientConnection, errorMessage)

		return
	}

	channels[clientConnection.GetChannel()].Delete(clientConnection)
	clientConnection.SetChannel(channelName)
	channels[channelName].Add(clientConnection)

	successMessage := fmt.Sprintf("Successfully entered channel with name: \"%s\"", channelName)
	sendSuccessMessage(clientConnection, successMessage)
}

// TODO: func leaveChannel

func sendErrorMessage(clientConnection *connections.ClientConnection, message string) {
	sendWsResponse(clientConnection, responses.WSResponse{
		Message: message,
		Status:  responses.STATUS_ERROR,
	})
}

func sendSuccessMessage(clientConnection *connections.ClientConnection, message string) {
	sendWsResponse(clientConnection, responses.WSResponse{
		Message: message,
		Status:  responses.STATUS_OK,
	})
}

func sendWsResponse(clientConnection *connections.ClientConnection, response responses.WSResponse) {
	err := clientConnection.GetWSConnection().WriteJSON(response)
	if err != nil {
		log.Println("Websocket error: ", err)
		channels[clientConnection.GetChannel()].Delete(clientConnection)
		_ = clientConnection.GetWSConnection().Close()
	}
}
