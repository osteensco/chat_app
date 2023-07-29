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

type Chatroom struct {
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

func (c Chatroom) handleConnections(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Error upgrading connection:", err)
		return
	}
	conn.Close()

	//print messages
}

func NewChatroom() *Chatroom {
	return &Chatroom{
		clients: make(ClientList),
	}
}

type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

type Client struct {
	connection *websocket.Conn
	chatroom   *Chatroom
}

type ClientList map[*Client]bool

func NewClient(conn *websocket.Conn, cr *Chatroom) *Client {
	return &Client{
		connection: conn,
		chatroom:   cr,
	}
}

func handleMessages() {
	// send messages to server
	// broadcast messages to each client
}

func main() {

	clients := NewChatroom()

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/ws", clients.handleConnections)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

}
