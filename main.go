package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	redisAddr := os.Getenv("REDISADDRESS")
	redisPass := os.Getenv("REDISPASSWORD")
	log.Println(redisAddr)
	log.Println(redisPass)
	cacheClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       0,
	})
	log.Println(cacheClient)
	// ctx := context.Background()

	// state, err2 := cacheClient.Ping(ctx).Result()
	// if err2 != nil {
	// 	log.Printf("error! %v", err2)
	// }
	// log.Println(state)

	index := NewChatroom("index", "home page")
	AllRooms["/ws_roombuilder"] = index

	http.HandleFunc("/ws_roombuilder", func(w http.ResponseWriter, r *http.Request) {
		websocketHandler(r.URL.Path, w, r)
	})

	http.HandleFunc("/ws_chatroom/", func(w http.ResponseWriter, r *http.Request) {
		roomPath := r.URL.Path
		roomPath = roomPath[len("/ws_chatroom/chatroom/"):]
		websocketHandler(roomPath, w, r)
	})

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/chatroom/", chatroomPathHandler)

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
