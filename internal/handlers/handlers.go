package handlers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/wfabjanczuk/awesomeProjectSlack/internal/requests"
	"github.com/wfabjanczuk/awesomeProjectSlack/internal/responses"
	"log"
	"net/http"
)

func WSEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client connected to the endpoint")

	wsResponse := responses.WSResponse{
		Message: `Connected to the server`,
	}

	conn := &WSConnection{
		connection: ws,
		channel:    PublicChannel,
	}
	channels[PublicChannel] = append(channels[PublicChannel], conn)

	err = ws.WriteJSON(wsResponse)
	if err != nil {
		log.Println("Failed to write JSON:", err)
	}

	go ListenForWS(conn)
}

func ListenForWS(conn *WSConnection) {
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

	channels[channelName] = []*WSConnection{}

	successMessage := fmt.Sprintf("Channel with name: \"%s\" successfully created.", channelName)
	sendSuccessMessage(wsConnection, successMessage)
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
		log.Println("Websocket err", err)
		_ = wsConnection.Close()
		// remove closed connection from state (especially from channels map)
	}
}
