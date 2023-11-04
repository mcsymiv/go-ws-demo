package handlers

import (
	"log"
	"net/http"
	"sort"

	"github.com/gorilla/websocket"
)

// gorilla websocket connection wrapper
type WsConnection struct {
	*websocket.Conn
}

// defines the response send back from websocket
type WsResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"messageType"`
	ConnectedUsers []string `json:"connectedUsers"`
}

type WsPayload struct {
	Action     string       `json:"action"`
	Username   string       `json:"username"`
	Message    string       `json:"message"`
	Connection WsConnection `json:"-"`
}

var wsChannel = make(chan WsPayload)

var clients = make(map[WsConnection]string)

func ListenWsConnection(conn *WsConnection) {
	defer func() {
		log.Println("Defer ListenWsConnection, recovering")
		if r := recover(); r != nil {
			log.Println("Unable to recover in listen", r)
		}
	}()

	var payload WsPayload
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
      // do nothing
		} else {
			payload.Connection = *conn
			wsChannel <- payload
		}
	}
}

// ListenWsChannel
// handles client actions
func ListenWsChannel() {
	var response WsResponse

	for {
		event := <-wsChannel
		switch event.Action {
    case "login", "connected":
			// get list of all users and send it back via broadcast
			clients[event.Connection] = event.Username
			users := getUsersFromClientsAsList()
      response.Action = "connectedUsers"
      response.ConnectedUsers = users
		  BroadcastToAll(response)
    case "left":
      response.Action = "connectedUsers"
      delete(clients, event.Connection)
			users := getUsersFromClientsAsList()
      response.ConnectedUsers = users
		  BroadcastToAll(response)
		}
	}
}

func getUsersFromClientsAsList() []string {
	var users []string
	for _, user := range clients {
    if user != "" {
		  users = append(users, user)
    }
	}

	sort.Strings(users)
	return users
}

func BroadcastToAll(res WsResponse) {
	for client := range clients {
		err := client.WriteJSON(res)
		if err != nil {
			log.Println("Unable to write to client", err)
		}
	}
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := updateConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Unable to upgrade ws connection", err)
	}
	log.Println("Connected to ws")

	var res WsResponse
	res.Message = "Connected to server"

	connection := WsConnection{
		Conn: ws,
	}

	err = ws.WriteJSON(res)
	if err != nil {
		log.Println("Unable to write json", err)
	}

	go ListenWsConnection(&connection)
}
