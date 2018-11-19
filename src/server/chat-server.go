package server

import (
	"bytes"
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
func WSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Unable to start websockets : %v ", err)
	}
	//defer conn.Close()

	Messages := make(chan []byte)
	//writer, err := conn.NextWriter(websocket.TextMessage)

	// if err != nil {
	// 	log.Fatalf("Unable to create writer : %v", err)
	// 	return
	// }
	go listen(Messages, conn)

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Fatalf("unable to read message from web socket: %v", err)
				break
			}
			message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
			//writer.Write(message)
			Messages <- message
		}
	}()

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

func HTTPHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "/home.html")
}
