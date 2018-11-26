package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//WSHandler handles web socket connections
func WSHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Unable to start websockets : %v ", err)
	}
	connection := &Connection{
		name: "",
		conn: conn,
	}
	hub.register <- connection
	go hub.read(conn)
	go hub.write()

	// for {
	// 	// Read in a new message as JSON and map it to a Message object
	// 	//err := conn.ReadJSON(&msg)
	// 	messageType, msg, err := conn.ReadMessage()
	// 	if err != nil {
	// 		log.Printf("error: %v", err)
	// 		break
	// 	}
	// 	fmt.Printf("Message %v received with %v type", msg, messageType)
	// 	// // Send the newly received message to the broadcast channel
	// 	// broadcast <- msg
	// }
}
