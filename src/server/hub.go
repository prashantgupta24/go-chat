package server

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/gorilla/websocket"
)

func initHub() {

}

func listen(Messages chan []byte, w *io.WriteCloser, conn *websocket.Conn) {
	writer, err := conn.NextWriter(websocket.TextMessage)

	if err != nil {
		log.Fatalf("Unable to create writer : %v", err)
		return
	}
	//writer := *w

	for message := range Messages {
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		fmt.Println("Message received : ", message)
		writer.Write(message)
	}
}
