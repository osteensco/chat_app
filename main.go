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

func (cr *Chatroom) handleConnections(w http.ResponseWriter, r *http.Request) {
	log.Println("new chatroom created")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Error upgrading connection:", err)
		return
	}

	client := NewClient(conn, cr)

	cr.addClient(client)

	client.handleMessages()
}

func NewChatroom() *Chatroom {
	return &Chatroom{
		clients: make(ClientList),
	}
}

type Client struct {
	connection *websocket.Conn
	chatroom   *Chatroom
	channel    chan []byte
}

func NewClient(conn *websocket.Conn, cr *Chatroom) *Client {
	log.Println("new client connected to chat room")
	return &Client{
		connection: conn,
		chatroom:   cr,
		channel:    make(chan []byte),
	}
}

func (c *Client) readMessages() {
	defer c.chatroom.removeClient(c)

	for {
		messageType, payload, err := c.connection.ReadMessage()

		if err != nil {
			log.Println(err)
			break
		}
		log.Printf("Message received {MessageType: %v, Payload: %v", messageType, string(payload))

		for client := range c.chatroom.clients {
			client.channel <- payload
		}
	}

}

func (c *Client) writeMessages() {
	defer c.chatroom.removeClient(c)

	for message := range c.channel {
		if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println(err)
			return
		}
		log.Println("message sent")
	}
}

func (c *Client) handleMessages() {
	go c.readMessages()
	go c.writeMessages()
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

// TODO
// - rename chatroom to better reflect its the actions of a client
// - consolidate client actions
// - create chat room creation workflows
// - randomize chat room and client names by default
// - pass more information about clients and chat rooms to backend
// - add in a data store?
