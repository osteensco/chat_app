package main

import (
	"log"

	"github.com/gorilla/websocket"
)

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
		log.Printf("Message received {MessageType: %v, Payload: %v}", messageType, string(payload))

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
