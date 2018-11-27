package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chat/src/server"
)

var addr = flag.String("addr", ":8000", "http service address")

func main() {
	flag.Parse()
	fmt.Printf("Starting chat server on port%v. Waiting for connections ... \n\n", *addr)
	chatServer := server.CreateChatServer()
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.WSHandler(chatServer, w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("error starting server: ", err)
	}
}
