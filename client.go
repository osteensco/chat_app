package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func pushToChannel(payload []byte, clients ClientList) {
	for c := range clients {
		c.Chatroom.Channel <- payload
	}
}

type Client struct {
	Name       string
	connection *websocket.Conn
	Chatroom   *Chatroom
}

func (c *Client) readMessages() {

	defer func() {
		log.Println("closing client connection in RM go routine")
		c.Chatroom.removeClient(c)
	}()

	for {
		messageType, payload, err := c.connection.ReadMessage()

		if err != nil {
			log.Println(err)
			pushToChannel([]byte("client disconnect"), c.Chatroom.clients) //push to channel so that writeMessage errors and closes the connection as well
			return
		}
		log.Printf("message received {MessageType: %v, Payload: '%v'}", messageType, string(payload))

		if !json.Valid(payload) {
			pushToChannel(payload, c.Chatroom.clients)
		} else {
			createNewChatroomFromMessage(c, payload)
		}

	}

}

func (c *Client) sendMessage(message []byte) error {
	if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Printf("error in room %v, error writing message - %v", c.Chatroom.Path, err)
		return err
	}
	log.Printf("message sent in chatroom: %v", c.Chatroom.Path)
	return nil
}

func (c *Client) writeMessages() {

	defer func() {

		log.Println("closing client connection in WM go routine")
		msg := c.Name + " has left the chatroom"
		go sendPostUserLeftMessage(c.Chatroom, msg)
		pushToChannel([]byte(msg), c.Chatroom.clients)
		c.Chatroom.removeClient(c)

	}()

	for message := range c.Chatroom.Channel {
		err := c.sendMessage(message)
		if err != nil {
			return
		}
	}

}

func (c *Client) handleMessages() {
	go c.readMessages()
	go c.writeMessages()
}

func NewClient(conn *websocket.Conn, cr *Chatroom, username string) *Client {

	log.Printf("new client connected to chatroom with path %v", cr.Path)
	return &Client{
		Name:       username,
		connection: conn,
		Chatroom:   cr,
	}

}
