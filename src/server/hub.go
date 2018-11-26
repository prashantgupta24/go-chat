package server

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
)

//Hub is the core component of the chat server
type Hub struct {
	clients    map[*websocket.Conn]string
	messages   chan []byte
	register   chan *Connection
	unregister chan *websocket.Conn
	clientNum  int
}

//Connection is used for handling the various web socket connections
type Connection struct {
	name string
	conn *websocket.Conn
}

//CreateHub creates a chat hub instance
func CreateHub() *Hub {
	hub := &Hub{
		clients:    make(map[*websocket.Conn]string),
		messages:   make(chan []byte),
		register:   make(chan *Connection),
		unregister: make(chan *websocket.Conn),
	}
	go startRegisterChannel(hub)
	go startUnregisterChannel(hub)
	return hub
}

func startRegisterChannel(hub *Hub) {
	for connection := range hub.register {
		//hub.clients[connection.conn] = connection.name
		hub.clientNum = hub.clientNum + 1 //This will be removed when we send name along with connection
		hub.clients[connection.conn] = "Client-" + strconv.Itoa(hub.clientNum)
		fmt.Printf("Started new web socket connection! Total connections : %v \n\n", len(hub.clients))
	}
}

func startUnregisterChannel(hub *Hub) {
	for conn := range hub.unregister {
		if _, ok := hub.clients[conn]; ok {
			conn.Close()
			fmt.Printf("Closing connection with %v, connections remaining : %v \n\n", hub.clients[conn], len(hub.clients)-1)
			delete(hub.clients, conn)
		}
	}
}

func (h *Hub) read(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("unable to read message from %v with error : %v \n\n", h.clients[conn], err)
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		h.messages <- message
	}
	h.unregister <- conn
}

func (h *Hub) write() {

	for message := range h.messages {
		for conn := range h.clients {
			writer, err := conn.NextWriter(websocket.TextMessage)

			if err != nil {
				log.Printf("Unable to create writer for client: %v Error: %v \n\n", h.clients[conn], err)
				return
			}
			//fmt.Println("Message received : ", string(message))
			b := []byte(h.clients[conn])
			b = append(b, []byte(" : ")...)
			b = append(b, message...)
			_, err = writer.Write(b)
			if err != nil {
				log.Fatalf("Error writing to websocket : %v", err)
			}
			//fmt.Printf("Writing %v bytes!\n", n)
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
