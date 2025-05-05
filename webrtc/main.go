package main

import (
	"fmt"
	"log"
	"net/http"
	"webrtc/server"
)

func main() {
	server.AllRooms.Init()

	http.HandleFunc("/create", server.CreateRoomRequestHandler)
	http.HandleFunc("/join", server.JoinRoomRequestHandler)

	log.Println("starting server on port 8080")
	fmt.Println(" ")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
