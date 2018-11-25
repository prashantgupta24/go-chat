package server

import (
	"bytes"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients  []*websocket.Conn
	messages chan []byte
	ConnChan chan *websocket.Conn
}

//CreateHub creates a chat hub instance
func CreateHub() *Hub {
	return &Hub{
		messages: make(chan []byte),
		ConnChan: make(chan *websocket.Conn),
	}
}

func (h *Hub) register() {
	for conn := range h.ConnChan {
		h.clients = append(h.clients, conn)
	}
}

func (h *Hub) read(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatalf("unable to read message from web socket: %v", err)
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		h.messages <- message
	}
}

func (h *Hub) write() {

	for message := range h.messages {
		for _, conn := range h.clients {
			writer, err := conn.NextWriter(websocket.TextMessage)

			if err != nil {
				log.Fatalf("Unable to create writer : %v", err)
				return
			}
			fmt.Println("Message received : ", string(message))
			n, err := writer.Write(message)
			if err != nil {
				log.Fatalf("Error writing to websocket : %v", err)
			}
			fmt.Printf("Writing %v bytes!\n", n)
			// Add queued chat messages to the current websocket message.
			// num := len(Messages)
			// for i := 0; i < num; i++ {
			// 	writer.Write(newline)
			// 	writer.Write(<-Messages)
			// }
			if err := writer.Close(); err != nil {
				return
			}
		}

	}

}
