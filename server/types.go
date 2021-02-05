package server

import "github.com/gorilla/websocket"

//Server is the main chat server interface
type Server interface {
	Read(Connection)
	Write()
	Messages() chan *MessageJSON
}

//MyChatServer is an implementation of the chat server interface
type MyChatServer struct {
	clients    map[*websocket.Conn]string
	messages   chan *MessageJSON
	unregister chan *websocket.Conn
	clientNum  int
}

//Connection is the interface for handling web socket connections
type Connection interface {
	Read() (*MessageJSON, error)
	Conn() *websocket.Conn
}

//ConnectionStruct is the implementation for handling the various web socket connections
type ConnectionStruct struct {
	name string
	conn *websocket.Conn
}

//MessageJSON is the core component of all messages being passed around in the chat server
type MessageJSON struct {
	MsgType messageType `json:"type"`
	Sender  string      `json:"sender"`
	Message string      `json:"message"`
}
