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
//		 - ${name of client} has diconnected
//			 - storing names and client cookies required

//	 - Data Stores API:
//		 - lobbyEP
//		 - refactor frontend to better utilize API endpoint methods
//		 - add cockroachDB functions
