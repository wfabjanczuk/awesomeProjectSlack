package handlers

import (
	"github.com/wfabjanczuk/awesomeProjectSlack/internal/connections"
	"github.com/wfabjanczuk/awesomeProjectSlack/internal/responses"
	"log"
)

func sendErrorMessage(clientConnection *connections.ClientConnection, message string) {
	sendWSResponse(clientConnection, responses.WSResponse{
		Message: message,
		Status:  responses.STATUS_ERROR,
	})
}

func sendSuccessMessage(clientConnection *connections.ClientConnection, message string) {
	sendWSResponse(clientConnection, responses.WSResponse{
		Message: message,
		Status:  responses.STATUS_OK,
	})
}

func sendWSResponse(clientConnection *connections.ClientConnection, response responses.WSResponse) {
	err := clientConnection.GetWSConnection().WriteJSON(response)
	if err != nil {
		log.Println("Websocket error: ", err)
		channels[clientConnection.GetChannel()].Delete(clientConnection)
		_ = clientConnection.GetWSConnection().Close()
	}
}
