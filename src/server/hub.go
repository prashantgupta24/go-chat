package server

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func initHub() {

}

func listen(Messages chan []byte, conn *websocket.Conn) {

	for message := range Messages {
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
