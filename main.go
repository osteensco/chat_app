package main

import (
	"log"
	"net/http"
)

func main() {

	index := NewChatroom("index", "home page")
	AllRooms["/ws"] = index

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocketHandler(r.URL.Path, w, r)
	})

	http.HandleFunc("/ws/chatroom/", func(w http.ResponseWriter, r *http.Request) {
		roomPath := r.URL.Path
		roomPath = roomPath[13:]
		websocketHandler(roomPath, w, r)
	})

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/chatroom/", chatroomPathHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

}

// TODO

// bugs

// add features
//	 - utilize redis and cockroachDB for persistent storage of state, defining chatroom lifecycle, and chatroom/chat data
//	 - add a client id and home screen name for logging
