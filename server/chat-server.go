package server

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
)

const (
	REGISTER   = "reg"
	UNREGISTER = "unreg"
	MESSAGE    = "msg"
)

//CreateChatServer creates a chat server instance
func CreateChatServer() Server {
	myServer := &MyChatServer{
		clients:  make(map[*websocket.Conn]string),
		messages: make(chan *MessageJSON),
		//register:   make(chan *Connection),
		unregister: make(chan *websocket.Conn),
	}
	//go startRegisterChannel(myServer)
	go startUnregisterChannel(myServer)
	return myServer
}

// func startRegisterChannel(myServer *MyChatServer) {
// 	for connection := range myServer.register {
// 		myServer.clients[connection.conn] = connection.name
// 		fmt.Printf("Started new web socket connection %v! Total connections : %v \n\n", connection.name, len(myServer.clients))
// 	}
// }

func startUnregisterChannel(myServer *MyChatServer) {
	for conn := range myServer.unregister {
		if _, ok := myServer.GetClients()[conn]; conn != nil && ok {
			conn.Close()
			myServer.messages <- &MessageJSON{
				Sender:  myServer.clients[conn],
				Message: strconv.Itoa(len(myServer.clients) - 1),
				MsgType: UNREGISTER,
			}
			fmt.Printf("Closing connection with %v, connections remaining : %v \n\n", myServer.clients[conn], len(myServer.clients)-1)
			delete(myServer.clients, conn)
		}
	}
}

func (s *MyChatServer) Read(conn Connection) {
	wbConnection := conn.GetConn()
	for {
		//var messageJSON MessageJSON
		messageJSON, err := conn.Read()
		if err != nil {
			log.Printf("error while parsing json message: %v", err)
			break
		}
		if messageJSON.MsgType == REGISTER {
			s.clients[wbConnection] = messageJSON.Sender
			fmt.Printf("Started new web socket connection %v! Total connections : %v \n\n", messageJSON.Message, len(s.clients))
			s.messages <- &MessageJSON{
				Sender:  messageJSON.Sender,
				Message: strconv.Itoa(len(s.clients)),
				MsgType: REGISTER,
			}
		} else {
			fmt.Printf("Reading %v from %v \n", messageJSON.Message, messageJSON.Sender)
			messageToSend := bytes.TrimSpace(bytes.Replace([]byte(messageJSON.Message), newline, space, -1))
			s.messages <- &MessageJSON{
				Sender:  messageJSON.Sender,
				Message: string(messageToSend),
				MsgType: MESSAGE,
			}
		}
	}
	s.unregister <- wbConnection
}

func (s *MyChatServer) Write() {

	for message := range s.messages {
		for conn := range s.clients {
			err := conn.WriteJSON(message)
			if err != nil {
				log.Printf("Unable to write for client: %v Error: %v \n\n", s.clients[conn], err)
				return
			}
		}

	}

}

//GetMessagesChan retrieves the messages chan
func (s *MyChatServer) GetMessagesChan() chan *MessageJSON {
	return s.messages
}

//GetClients retrieves the clients
func (s *MyChatServer) GetClients() map[*websocket.Conn]string {
	return s.clients
}

//Register function adds the connection to the chat server
// func (s *MyChatServer) Register(connection *Connection) {
// 	s.register <- connection
// }
