package main

import (
	"encoding/json"
	"log"

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

func NewClient(conn *websocket.Conn, cr *Chatroom) *Client {
	log.Println("new client connected to chat room")
	return &Client{
		connection: conn,
		chatroom:   cr,
		channel:    make(chan []byte),
	}
}

func pushToChannel(payload []byte, clients ClientList) {
	for c := range clients {
		c.channel <- payload
	}
}

type Client struct {
	connection *websocket.Conn
	chatroom   *Chatroom
	channel    chan []byte
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

		if !json.Valid(payload) {
			pushToChannel(payload, c.chatroom.clients)
		} else {
			//TODO
			// send to redis
			// break this out into a separate function
			newroom := NewSubmittedRoom(payload)
			roomstruct := NewChatroom(newroom.Name, newroom.Path)
			roompath := roomstruct.path
			AllRooms[roompath] = roomstruct
			pushToChannel(payload, c.chatroom.clients)
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
