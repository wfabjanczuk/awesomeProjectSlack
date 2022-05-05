package handlers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type WSConnection struct {
	*websocket.Conn
}

type WSPayload struct {
	Action   string       `json:"action"`
	Username string       `json:"username"`
	Message  string       `json:"message"`
	Conn     WSConnection `json:"-"`
}

type WSResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

var inputChan = make(chan WSPayload)
var clients = make(map[WSConnection]string)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func WSEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client connected to the endpoint")

	response := WSResponse{
		Message: `Connected to the server`,
	}

	conn := WSConnection{
		Conn: ws,
	}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println("Failed to write JSON:", err)
	}

	go ListenForWs(&conn)
}

func ListenForWs(conn *WSConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload WSPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			log.Println(err)
			break
		} else {
			fmt.Println(payload)

			payload.Conn = *conn
			inputChan <- payload
		}
	}
}

func ListenToWsChannel() {
	var response WSResponse

	for {
		event := <-inputChan

		switch event.Action {
		case "broadcast":
			response.Action = "broadcast"
			response.Message = event.Message
			broadcastToAll(response)
		}
	}
}

func broadcastToAll(response WSResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("Websocket err", err)
			_ = client.Close()
			delete(clients, client)
		}
	}
}
