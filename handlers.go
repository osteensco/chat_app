package main

import (
	"log"
	"net/http"
	"strings"
	"time"
)

const Address = "localhost"
const Port = "8080"

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

func sendRemoveFromLobbyRequest(cr *Chatroom) {

	apiURL := "http://" + Address + ":" + Port + "/api/lobby?roomname=" + cr.name + "&roompath=" + cr.Path
	req, err := http.NewRequest("DELETE", strings.ReplaceAll(apiURL, " ", "%20"), nil)
	if err != nil {
		log.Println("Error creating DELETE request:", err)
		return
	}

	reqclient := &http.Client{}
	resp, err := reqclient.Do(req)
	if err != nil {
		log.Println("Error sending DELETE request:", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("DELETE request to %v failed with status code: %v", apiURL, resp.StatusCode)
		return
	}

	defer resp.Body.Close()

}

func sendRemoveMessagesRequest(cr *Chatroom) {

	apiURL := "http://" + Address + ":" + Port + "/api/messages?roompath=" + cr.Path
	req, err := http.NewRequest("DELETE", strings.ReplaceAll(apiURL, " ", "%20"), nil)
	if err != nil {
		log.Println("Error creating DELETE request:", err)
		return
	}

	reqclient := &http.Client{}
	resp, err := reqclient.Do(req)
	if err != nil {
		log.Println("Error sending DELETE request:", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("DELETE request to %v failed with status code: %v", apiURL, resp.StatusCode)
		return
	}

	defer resp.Body.Close()

}

func sendPostUserLeftMessage(cr *Chatroom, msg string) {
	apiURL := "http://" + Address + ":" + Port + "/api/messages?roompath=" + cr.Path + "&chatmessage=" + msg
	req, err := http.NewRequest("POST", strings.ReplaceAll(apiURL, " ", "%20"), nil)
	if err != nil {
		log.Println("Error creating POST request:", err)
		return
	}

	reqclient := &http.Client{}
	resp, err := reqclient.Do(req)
	if err != nil {
		log.Println("Error sending POST request:", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("POST request to %v failed with status code: %v", apiURL, resp.StatusCode)
		return
	}

	defer resp.Body.Close()
}

func monitorRoomActivity(rooms *RoomList) {

	for {

		for _, chatroom := range *rooms {
			if chatroom.name != "index" && len(chatroom.clients) == 0 {
				log.Printf("Chatroom %v with path %v is not populated, starting timer for removal", chatroom.name, chatroom.Path)
				go chatroom.startRemovalTimer()
			}
		}

		time.Sleep(10 * time.Minute)

	}

}
