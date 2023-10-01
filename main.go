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

	err := http.ListenAndServe(":8080", nil)
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

//	 - Data Stores API:
//		 - changeName messages need to POST to DB
//		 - on app load register existing chatrooms from DB
//		 - add logic to delete chatrooms when empty for a period of time
//			 - as well as other cleanup functions for disconnections, etc
//		 - add cockroachDB functions
