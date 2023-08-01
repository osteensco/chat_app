package main

import (
	"log"
	"net/http"
)

var AllRooms RoomList = make(RoomList)

func chatroomPathHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("checking if room exists...")
	roomPath := r.URL.Path[len("/chatroom/"):]
	room, ok := AllRooms[roomPath]
	if !ok {
		log.Printf("404 error: room path %v not found", roomPath)
		http.NotFound(w, r)
		return
	} else {
		log.Printf("registering path for room: %v", room.name)
	}
	http.ServeFile(w, r, "./static/chatroom.html")
}

func main() {

	index := NewChatroom("index", "")

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/ws", index.handleConnections)

	http.HandleFunc("/chatroom/", chatroomPathHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

}

// TODO

// - figure out how to handle specific url paths on chatroom creation, break out process into its own function

// - utilize redis and cockroachDB for persistent storage of chatrooms and chatroom data
