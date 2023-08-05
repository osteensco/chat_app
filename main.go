package main

import (
	"log"
	"net/http"
)

func main() {

	index := NewChatroom("index", "home page")
	AllRooms["/ws_roombuilder"] = index

	http.HandleFunc("/ws_roombuilder", func(w http.ResponseWriter, r *http.Request) {
		websocketHandler(r.URL.Path, w, r)
	})

	http.HandleFunc("/ws_chatroom/", func(w http.ResponseWriter, r *http.Request) {
		roomPath := r.URL.Path
		roomPath = roomPath[len("/ws_chatroom/chatroom/"):]
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
//	 - chatroom messages:
//		 - ${name of client} has diconnected
//		 - ${old display name} has changed their name to ${new name}
//		 - Randomize Anonymous with a series of numbers e.g. Anonymous239523
//	 - utilize redis and cockroachDB for persistent storage of state, defining chatroom lifecycle, and chatroom/chat data
//	 - add a client id and implement cookies for logging
