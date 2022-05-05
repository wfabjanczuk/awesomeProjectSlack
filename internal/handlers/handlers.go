package handlers

import (
	"fmt"
	"github.com/gorilla/websocket"
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

	conn := connections.NewClientConnection(wsConnection, PublicChannel)
	channels[PublicChannel] = append(channels[PublicChannel], conn)

	sendSuccessMessage(wsConnection, "Connected to the server")
	go ListenForWS(conn)
}

func ListenForWS(conn *connections.ClientConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload requests.WSPayload

	for {
		err := conn.GetWSConnection().ReadJSON(&payload)
		if err != nil {
			log.Println(err)
			sendErrorMessage(conn.GetWSConnection(), "Invalid JSON payload!")

			break
		} else {
			fmt.Printf("New message from client, action: \"%s\", message: \"%s\"\n", payload.Action, payload.Message)

			requestQueue <- requests.NewWSRequest(payload, conn)
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
			createChannel(wsRequest.Payload.Message, wsRequest.ClientConnection.GetWSConnection())
		case "enter":
			enterChannel(wsRequest.Payload.Message, wsRequest.ClientConnection)
		}
	}
}

func broadcastToChannel(channelName, message string) {
	for _, clientConnection := range channels[channelName] {
		sendSuccessMessage(clientConnection.GetWSConnection(), message)
	}
}

func createChannel(channelName string, wsConnection *websocket.Conn) {
	if _, ok := channels[channelName]; ok {
		errorMessage := fmt.Sprintf("Channel with name: \"%s\" already exists!", channelName)
		sendErrorMessage(wsConnection, errorMessage)

		return
	}

	channels[channelName] = []*connections.ClientConnection{}

	successMessage := fmt.Sprintf("Channel with name: \"%s\" successfully created.", channelName)
	sendSuccessMessage(wsConnection, successMessage)
}

func enterChannel(channelName string, clientConnection *connections.ClientConnection) {
	if _, ok := channels[channelName]; !ok {
		errorMessage := fmt.Sprintf("Channel with name: \"%s\" does not exist!", channelName)
		sendErrorMessage(clientConnection.GetWSConnection(), errorMessage)

		return
	}

	clientConnection.SetChannel(channelName)
	channels[channelName] = append(channels[channelName], clientConnection)

	successMessage := fmt.Sprintf("Successfully entered channel with name: \"%s\"", channelName)
	sendSuccessMessage(clientConnection.GetWSConnection(), successMessage)
}

// TODO: func enterChannel, func leaveChannel

func sendErrorMessage(wsConnection *websocket.Conn, message string) {
	sendWsResponse(wsConnection, responses.WSResponse{
		Message: message,
		Status:  responses.STATUS_ERROR,
	})
}

func sendSuccessMessage(wsConnection *websocket.Conn, message string) {
	sendWsResponse(wsConnection, responses.WSResponse{
		Message: message,
		Status:  responses.STATUS_OK,
	})
}

func sendWsResponse(wsConnection *websocket.Conn, response responses.WSResponse) {
	err := wsConnection.WriteJSON(response)
	if err != nil {
		log.Println("Websocket error: ", err)
		_ = wsConnection.Close()
		// TODO: remove closed connection from state (especially from channels map)
	}
}
