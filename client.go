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
	log.Printf("new client connected to chatroom with path %v", cr.Path)
	return &Client{
		connection: conn,
		Chatroom:   cr,
	}
}

func pushToChannel(payload []byte, clients ClientList) {
	for c := range clients {
		c.Chatroom.Channel <- payload
	}
}

type Client struct {
	connection *websocket.Conn
	Chatroom   *Chatroom
}

func (c *Client) readMessages() {
	defer c.Chatroom.removeClient(c)

	for {
		messageType, payload, err := c.connection.ReadMessage()

		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("message received {MessageType: %v, Payload: '%v'}", messageType, string(payload))

		if !json.Valid(payload) {
			pushToChannel(payload, c.Chatroom.clients)
		} else {
			//TODO
			// send to redis
			// break this out into a separate function
			newroom := NewSubmittedRoom(payload)
			roomstruct := NewChatroom(newroom.Name, newroom.Path)
			roompath := roomstruct.Path
			AllRooms[roompath] = roomstruct
			pushToChannel(payload, c.Chatroom.clients)
		}

	}

}

func (c *Client) writeMessages() {
	defer c.Chatroom.removeClient(c)

	for message := range c.Chatroom.Channel {
		if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("errror in room %v, error writing message - %v", c.Chatroom.Path, err)
			break
		}
		log.Printf("message sent in chatroom: %v", c.Chatroom.Path)
	}
}

func (c *Client) handleMessages() {
	go c.readMessages()
	go c.writeMessages()
}
