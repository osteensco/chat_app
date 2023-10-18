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

	cacheClient := connectRedis(ctx, 1)
	log.Println(cacheClient)

	databaseClient := connectCockrochDB(ctx)
	defer databaseClient.Close()

	index := NewChatroom("index", "home page")
	AllRooms["/ws_lobby"] = index

	setHandlers()
	initAPI(ctx, cacheClient, databaseClient)

	chatrooms, err := getAllChatroomsCRDB(ctx, databaseClient, "lobby")
	if err != nil {
		log.Printf("Unable to get chatrooms from database: %v", err)
	} else {
		log.Println("Building out RoomList map...")
	}

	for key, val := range chatrooms {
		room := NewChatroom(key, val)
		AllRooms[val] = room
		log.Printf("Added %v: %v to RoomList map", key, val)
	}

	go monitorRoomActivity(&AllRooms)

	err = http.ListenAndServe(":"+Port, nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

}

// TODO
//	 - app should be resilient to cache failure
