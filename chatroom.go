package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type NewRoomParse struct {
	Chatroom SubmittedRoom
}

type SubmittedRoom struct {
	Name string
	Path string
}

func NewSubmittedRoom(payload []byte) *SubmittedRoom {

	var rm NewRoomParse
	err := json.Unmarshal(payload, &rm)
	if err != nil {
		log.Println("Error parsing JSON: ", err)
		return nil
	}
	log.Printf("`%v` room received from client with path %v", rm.Chatroom.Name, rm.Chatroom.Path)

	return &rm.Chatroom

}

func createNewChatroomFromMessage(c *Client, payload []byte) {
	newroom := NewSubmittedRoom(payload)
	roompath := newroom.Path
	roomstruct := NewChatroom(newroom.Name, roompath)
	AllRooms[roompath] = roomstruct
	pushToChannel(payload, c.Chatroom.clients)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type ClientList map[*Client]bool

type RoomList map[string]*Chatroom

type Chatroom struct {
	name    string
	Path    string
	clients ClientList
	Channel chan []byte
}

func (cr *Chatroom) registerClient(client *Client) {

	cr.clients[client] = true
	log.Printf("registered client with chatroom %v", cr.Path)

}

func (cr *Chatroom) removeClient(client *Client) {

	if _, ok := cr.clients[client]; ok {
		log.Printf("removing client from %v", cr.Path)
		client.connection.Close()

		delete(cr.clients, client)
	}

}

func (cr *Chatroom) handleConnections(w http.ResponseWriter, r *http.Request) {

	log.Printf("new client entering chatroom with path %v", cr.Path)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Error upgrading connection:", err)
		return
	}

	client := NewClient(conn, cr)

	cr.registerClient(client)

	client.handleMessages()

}

func (cr *Chatroom) startRemovalTimer() {

	startTime := time.Now()
	timer := time.NewTimer(1 * time.Minute)

	for {
		<-timer.C
		if len(cr.clients) != 0 {
			log.Printf("removal timer stopped for chatroom %v with path %v, %v users currently in chatroom", cr.name, cr.Path, len(cr.clients))
			timer.Stop()
			return
		} else {
			if time.Since(startTime) >= 10*time.Minute {

				timer.Stop()
				log.Printf("removing chatroom %v with path %v", cr.name, cr.Path)
				delete(AllRooms, cr.Path)
				go sendRemoveFromLobbyRequest(cr)
				go sendRemoveMessagesRequest(cr)
			}
		}

	}

}

func NewChatroom(n string, p string) *Chatroom {

	log.Printf("new chatroom created with path %v", p)
	return &Chatroom{
		name:    n,
		Path:    p,
		clients: make(ClientList),
		Channel: make(chan []byte),
	}

}
