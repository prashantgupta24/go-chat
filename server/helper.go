package server

import (
	"errors"
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
func WSHandler(server Server, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Unable to start websockets : %v ", err)
	}

	connectionStruct := &ConnectionStruct{
		conn: conn,
	}
	go server.Read(connectionStruct)
	go server.Write()
}

//GetConn function
func (c *ConnectionStruct) GetConn() *websocket.Conn {
	return c.conn
}

//Read function
func (c *ConnectionStruct) Read() (*MessageJSON, error) {
	var i *MessageJSON
	err := c.conn.ReadJSON(&i)
	if i == nil {
		return nil, errors.New("error while parsing")
	}
	return i, err
}
