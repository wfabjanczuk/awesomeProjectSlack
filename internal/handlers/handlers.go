package handlers

import (
	"fmt"
	"github.com/wfabjanczuk/awesomeProjectSlack/internal/models"
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

	conn := WSConnection{
		connection: ws,
		channel:    *publicChannel,
	}
	channels[*publicChannel] = append(channels[*publicChannel], conn)

	err = ws.WriteJSON(wsResponse)
	if err != nil {
		log.Println("Failed to write JSON:", err)
	}

	go ListenForWS(&conn)
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
			break
		} else {
			fmt.Println(payload)

			requestQueue <- requests.NewWSRequest(payload, conn)
		}
	}
}

func ListenToRequestQueue() {
	for {
		wsRequest := <-requestQueue
		channel := wsRequest.Connection.GetChannel()
		payload := wsRequest.Payload

		switch payload.Action {
		case "broadcast":
			broadcastToChannel(channel, responses.WSResponse{
				Message: payload.Message,
				Status:  responses.WS_STATUS_OK,
			})
		}
	}
}

func broadcastToChannel(channel models.Channel, response responses.WSResponse) {
	for _, clientConnection := range channels[channel] {
		err := clientConnection.GetWSConnection().WriteJSON(response)
		if err != nil {
			log.Println("Websocket err", err)
			_ = clientConnection.GetWSConnection().Close()
		}
	}
}
