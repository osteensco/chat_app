package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	connection *websocket.Conn
	chatroom   *Chatroom
	channel    chan []byte
}

type NewRoomParse struct {
	Chatroom SubmittedRoom
}

type SubmittedRoom struct {
	name string
	path string
}

func NewSubmittedRoom(payload []byte) *SubmittedRoom {
	var rm NewRoomParse
	err := json.Unmarshal(payload, &rm)
	if err != nil {
		log.Println("Error parsing JSON: ", err)
		return nil
	}

	return &rm.Chatroom
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

		newroom := NewSubmittedRoom(payload)

		if newroom == nil {
			for client := range c.chatroom.clients {
				client.channel <- payload
			}
		} else {

			// break this out into a separate function

			roomstruct := NewChatroom(newroom.name, newroom.path)
			roomname := roomstruct.name

			//TODO
			// send to redis
			AllRooms[roomname] = roomstruct

			// determine correct way to do this
			http.Handle(fmt.Sprintf("/%v", roomstruct.path), http.FileServer(http.Dir("./static/chatroom.html")))
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
