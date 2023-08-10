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
	AllRooms["/ws_roombuilder"] = index

	setHandlers()
	initAPI()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

}

// TODO

// bugs

// add features
//	 - chatroom messages:
//		 - ${name of client} has diconnected
//			 - storing names and client cookies required
//		 - ${old display name} has changed their name to ${new name}
//			 - storing names and client cookies required

//	 - utilize redis and cockroachDB for persistent storage of state, defining chatroom lifecycle, and chatroom/chat data
//		 - currently capturing all available chatrooms
//			 - should store as basic key value pair in Redis
//		 - capture message history for each chatroom
//			 - may need to be hash data type in Redis or JSON module if possible with go-redis
//		 - add a client id and implement cookies for logging and chatroom data
//			 - clientList may need to be converted to master client list or just for each chatroom (hash)
//		 - add login and profile page
