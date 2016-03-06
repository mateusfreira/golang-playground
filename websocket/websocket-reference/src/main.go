package main

import (
	"log"
	"net/http"

	"./websocket"
)

func main() {
	log.SetFlags(log.Lshortfile)

	// websocket server
	server := websocket.NewServer("/")
	go server.Listen()

	// static files
	//http.Handle("/", http.FileServer(http.Dir("webroot")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
