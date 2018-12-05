package server

import "github.com/gorilla/websocket"

//Server is the main chat server interface
type Server interface {
	Read(*websocket.Conn)
	Write()
	Register(*Connection)
}

//MyChatServer is an implementation of the chat server interface
type MyChatServer struct {
	clients    map[*websocket.Conn]string
	messages   chan *MessageJSON
	register   chan *Connection
	unregister chan *websocket.Conn
	clientNum  int
}

//Connection is used for handling the various web socket connections
type Connection struct {
	name string
	conn *websocket.Conn
}

//MessageJSON is the core component of all messages being passed around in the chat server
type MessageJSON struct {
	MsgType string `json:"type"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
}
