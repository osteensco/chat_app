package main

import (
	"log"
	"net/http"
)

var AllRooms RoomList = make(RoomList)

func chatroomPathHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("checking if room exists...")
	roomPath := r.URL.Path[len("/chatroom/"):]
	if roomPath[len(roomPath)-3:] == "/ws" { //unnecessary?
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

func websocketHandler(roompath string, w http.ResponseWriter, r *http.Request) {
	log.Printf("websocket handle for path %v", roompath)
	AllRooms[roompath].handleConnections(w, r)
}