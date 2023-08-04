package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

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

func (cr *Chatroom) addClient(client *Client) {
	log.Printf("adding client to chatroom with path %v", cr.Path)
	cr.clients[client] = true

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

	cr.addClient(client)

	client.handleMessages()
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
