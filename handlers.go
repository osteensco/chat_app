package main

import (
	"log"
	"net/http"
)

var AllRooms RoomList = make(RoomList)

func chatroomPathHandler(w http.ResponseWriter, r *http.Request) {

	roomPath := r.URL.Path[len("/chatroom/"):]
	log.Printf("requested chatroom for %v", roomPath)

	room, ok := AllRooms[roomPath]
	log.Println(AllRooms)

	if !ok {

		log.Printf("404 error: room path %v not found", roomPath)
		http.NotFound(w, r)
		return
	} else {
		log.Printf("serving %v chatroom for %v", room.name, roomPath)
		http.ServeFile(w, r, "./static/chatroom.html")
	}

}

func websocketHandler(roompath string, w http.ResponseWriter, r *http.Request) {

	log.Printf("websocket handle for path %v", roompath)
	room, ok := AllRooms[roompath]
	if !ok {
		log.Printf("websocket connection failed for room %v, room not found in AllRooms map", roompath)
	} else {
		room.handleConnections(w, r)
	}

}

func setHandlers() {

	http.HandleFunc("/ws_lobby", func(w http.ResponseWriter, r *http.Request) {
		websocketHandler(r.URL.Path, w, r)
	})

	http.HandleFunc("/ws_chatroom/", func(w http.ResponseWriter, r *http.Request) {
		roomPath := r.URL.Path
		roomPath = roomPath[len("/ws_chatroom/chatroom/"):]
		websocketHandler(roomPath, w, r)
	})

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/chatroom/", chatroomPathHandler)

}
