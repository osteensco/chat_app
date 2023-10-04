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

//	 - ${userName} has left chatroom message
//	 - remove chatrooms and messages from DB after 10 min of inactivity
//		 - use go routine chatroom method to start a timer when clientList becomes empty
//	 - add cockroachDB functions

// bugs
//	- 404 error on cached chatrooms
//	- messages aren't displaying
