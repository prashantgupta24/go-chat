package server

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
)

//CreateChatServer creates a chat server instance
func CreateChatServer() Server {
	myServer := &MyChatServer{
		clients:    make(map[*websocket.Conn]string),
		messages:   make(chan *Message),
		register:   make(chan *Connection),
		unregister: make(chan *websocket.Conn),
	}
	go startRegisterChannel(myServer)
	go startUnregisterChannel(myServer)
	return myServer
}

func startRegisterChannel(myServer *MyChatServer) {
	for connection := range myServer.register {
		//myServer.clients[connection.conn] = connection.name
		myServer.clientNum = myServer.clientNum + 1 //This will be removed when we send name along with connection
		connectionName := "Client-" + strconv.Itoa(myServer.clientNum)
		myServer.clients[connection.conn] = connectionName
		fmt.Printf("Started new web socket connection %v! Total connections : %v \n\n", connectionName, len(myServer.clients))
	}
}

func startUnregisterChannel(myServer *MyChatServer) {
	for conn := range myServer.unregister {
		if _, ok := myServer.clients[conn]; ok {
			conn.Close()
			fmt.Printf("Closing connection with %v, connections remaining : %v \n\n", myServer.clients[conn], len(myServer.clients)-1)
			delete(myServer.clients, conn)
		}
	}
}

func (s *MyChatServer) Read(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("unable to read message with error : %v \n\n", err)
			break
		}
		clientName := s.clients[conn]
		//fmt.Printf("Reading %v from %v \n", message, clientName)
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		s.messages <- &Message{
			data:   &message,
			sender: &clientName,
		}
	}
	s.unregister <- conn
}

func (s *MyChatServer) Write() {

	for message := range s.messages {
		for conn := range s.clients {
			writer, err := conn.NextWriter(websocket.TextMessage)

			if err != nil {
				log.Printf("Unable to create writer for client: %v Error: %v \n\n", s.clients[conn], err)
				return
			}
			//fmt.Println("Message received : ", string(message))
			b := []byte(*message.sender)
			b = append(b, []byte(" : ")...)
			b = append(b, *message.data...)
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

//Register function adds the connection to the chat server
func (s *MyChatServer) Register(connection *Connection) {
	s.register <- connection
}
