package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chat/src/server"
)

func main() {
	fmt.Println("starting")
	http.Handle("/", http.FileServer(http.Dir("./static")))
	//http.HandleFunc("/", server.HTTPHandler)
	http.HandleFunc("/ws", server.WSHandler)

	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		log.Fatal("error starting server: ", err)
	}
	log.Println("http server started on :8000")
}
