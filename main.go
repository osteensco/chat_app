package main

import (
	"log"
	"net/http"
)

func main() {

	clients := NewChatroom()

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/ws", clients.handleConnections)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

}

// TODO

// - set up process to create a chat room
// 		- receive room names created from front end
// 		- create a new chat room and http.Handle with room name as the path and serve chatroom.html
// 		- notify frontend of room creation
// 			- this should notify all clients on home page (home page websocket)
// 		- frontend should update list of rooms

// - set up a different websocket connection for the home page
// 		- this will be used to maintain an real time list of chatrooms available on the server

// - utilize redis and cockroachDB for persistent storage of chatrooms, and chatroom data
