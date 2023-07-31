package main

import (
	"log"
	"net/http"
)

var AllRooms RoomList = make(RoomList)

func main() {

	index := NewChatroom("index", "")

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/ws", index.handleConnections)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

}

// TODO

// - figure out how to handle specific url paths on chatroom creation, break out process into its own function

// - utilize redis and cockroachDB for persistent storage of chatrooms and chatroom data
