package main

import (
	"log"
	"net/http"
)

var AllRooms RoomList = make(RoomList)

func chatroomPathHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("checking if room exists...")
	roomPath := r.URL.Path[len("/chatroom/"):]
	if roomPath[len(roomPath)-3:] == "/ws" {
		roomPath = roomPath[:len(roomPath)-3]
	}
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
	AllRooms["/ws"] = index

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("websocket handle for path %v", r.URL.Path)
		AllRooms[r.URL.Path].handleConnections(w, r)
	})

	http.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		roomPath := r.URL.Path
		roomPath = roomPath[13:]
		log.Printf("websocket handle for path %v", roomPath)
		AllRooms[roomPath].handleConnections(w, r)
	})

	http.HandleFunc("/chatroom/", chatroomPathHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

}

// TODO

// - figure out how to handle a basic get request vs a /ws request

// - utilize redis and cockroachDB for persistent storage of chatrooms and chatroom data
