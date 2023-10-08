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
	defer databaseClient.Close()

	index := NewChatroom("index", "home page")
	AllRooms["/ws_lobby"] = index

	setHandlers()
	initAPI(ctx, cacheClient, databaseClient)

	chatrooms, err := getAllChatroomsRedis(ctx, cacheClient, "lobby")
	if err != nil {
		log.Printf("Unable to get chatrooms from cache: %v", err)
	}

	for key, val := range chatrooms {
		val = val[1 : len(val)-1]
		room := NewChatroom(key, val)
		AllRooms[val] = room
		log.Printf("%v: %v", key, val)
	}

	go monitorRoomActivity(&AllRooms)

	err = http.ListenAndServe(":"+Port, nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

}

// TODO

//	 - add cockroachDB functions

// bugs
