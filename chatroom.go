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
	path    string
	clients ClientList
}

func (cr *Chatroom) addClient(client *Client) {

	cr.clients[client] = true

}

func (cr *Chatroom) removeClient(client *Client) {

	if _, ok := cr.clients[client]; ok {

		client.connection.Close()

		delete(cr.clients, client)
	}
}

func (cr *Chatroom) handleConnections(w http.ResponseWriter, r *http.Request) {
	log.Println("new client entering the chatroom")

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
	log.Println("new chatroom created")
	return &Chatroom{
		name:    n,
		path:    p,
		clients: make(ClientList),
	}
}
