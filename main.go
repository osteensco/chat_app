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

//	 - refactor cache layer:
//		 - Redis keys should be set to expire after a certain amount of time.
//			- lobby key should never expire
//			- users set 10 min expiration
//			- messages list 10 min expiration

//		 - Define what data should be fast retrieval
//			 - Everything persisted needs fast retrieval
//				 - Lobby
//				 - Users
//				 - Messages

//		 - Check Redis first for data and fallback to cockroachDB
//		 - CockroachDB is considered the source of truth, this should be reflected in the code

//		*** GET request -> check redis -> if error or nothing returned -> check cockroachDB -> write to redis or restart decrepancy timer if already in redis
//		****** All other requests -> cockroachDB function -> if no error -> redis function
//		********* This should ensure cockroachDB is maintained as source of truth

// bugs

//	 - window unloads before both removeDisplayNameFromRoom and roomExitMessage functions fire
