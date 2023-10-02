package main

import (
	"context"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	ctx := context.Background()

	cacheClient := connectRedis(ctx)
	log.Println(cacheClient)

	databaseClient := connectCockrochDB(ctx)
	defer databaseClient.Close(ctx)

	index := NewChatroom("index", "home page")
	AllRooms["/ws_lobby"] = index

	setHandlers()
	initAPI(ctx, cacheClient)

	chatrooms, err := getAllChatroomsRedis(ctx, cacheClient, "lobby")
	if err != nil {
		log.Printf("Unable to get chatrooms from cache: %v", err)
	}

	for key, val := range chatrooms {
		room := NewChatroom(key, val)
		AllRooms[val] = room
	}

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

}

// TODO

// bugs

// add features
//	 - chatroom messages:
//		 -  ${userName} has entered the room
//		 - ${userName} has left the room
//			 - cookies necessary?
//		* on client disconnect message, parse payload and grab displayname in backend
//			-need to send displayname on disconnect from FE
//		* read message on FE for other clients
//		* store in DB

//	 - Data Stores API:

//		 - add logic to delete chatrooms when empty for a period of time
//			 - as well as other cleanup functions for disconnections, messages when room doesn't exist, etc
//		 - add cockroachDB functions
