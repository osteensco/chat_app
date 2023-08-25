package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
)

func lobbyEP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("lobbyEP %v", r.Method)
		}(w, r)
	case "POST":
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("lobbyEP %v", r.Method)
		}(w, r)
	case "DELETE":
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("lobbyEP %v", r.Method)
		}(w, r)
	}

}

func chatroomsEP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("chatroomsEP %v", r.Method)
		}(w, r)
	case "POST":
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("chatroomsEP %v", r.Method)
		}(w, r)
	case "PUT":
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("chatroomsEP %v", r.Method)
		}(w, r)
	case "DELETE":
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("chatroomsEP %v", r.Method)
		}(w, r)
	}

}

func usersEP(w http.ResponseWriter, r *http.Request, ctx context.Context, redisClient *redis.Client) {

	defer func() {
		if rec := recover(); rec != nil {
			log.Println("Recovered from panic:", rec)
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{
				"ok":     false,
				"error":  rec,
				"status": http.StatusBadRequest,
			}
			json.NewEncoder(w).Encode(response)
		}
	}()

	roompath := r.URL.Query().Get("roompath")
	if roompath == "" {
		log.Panicf("roompath query parameter not provided! Request URL provided was %v", r.URL)
	}

	displayname := r.URL.Query().Get("displayname")
	if roompath == "" {
		log.Panicf("displayname query parameter not provided! Request URL provided was %v", r.URL)
	}

	switch r.Method {

	case "GET":
		// used when a new client enters a room or a clients displayname is changed
		func(w http.ResponseWriter, r *http.Request) {

			log.Printf("GET %v FROM usersEP", displayname)

			displayNameExists, err := isUserInChatroomRedis(ctx, redisClient, displayname, roompath)

			if err != nil {
				log.Panicf("Error querying Redis with displayname %v %v", displayname, http.StatusInternalServerError)
			} else if !displayNameExists {
				log.Printf("%v does not exist in chatroom %v", displayname, roompath)
				w.WriteHeader(http.StatusBadRequest)
				response := map[string]interface{}{
					"ok":     false,
					"error":  "Display name is missing",
					"status": http.StatusBadRequest,
				}
				json.NewEncoder(w).Encode(response)
			} else {
				w.WriteHeader(http.StatusOK)
			}

		}(w, r)

	case "POST":
		// used when a new client enters a room
		func(w http.ResponseWriter, r *http.Request) {

			err := addUserToChatroomRedis(ctx, redisClient, displayname, roompath)
			if err != nil {
				log.Panicf("Error adding user to chatroom %v", http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
			}

		}(w, r)

	case "PUT":
		// used when a new client changes their display name
		func(w http.ResponseWriter, r *http.Request) {
			if displayname != "" {
				log.Printf("PUT %v FROM usersEP", displayname)
			} else {
				log.Printf("PUT %v FROM usersEP", displayname)
			}
		}(w, r)

	case "DELETE":
		// used when a client leaves a room or a room is removed from the server
		func(w http.ResponseWriter, r *http.Request) {
			if displayname != "" {
				log.Printf("DELETE %v FROM usersEP", displayname)
			} else {
				log.Printf("DELETE %v FROM usersEP", displayname)
			}
		}(w, r)

	}

}

func initAPI(ctx context.Context, redisClient *redis.Client) {

	http.HandleFunc("/api/lobby", lobbyEP)
	// lobby: {
	// 	room name: xxxxxxx,
	// 	room path: xxxxxxxxxxxxxxxxx
	// }
	http.HandleFunc("/api/chatrooms", chatroomsEP)
	// chatrooms: {
	// 	chatroom path:
	// 	messages: {
	// 		timestamp: xxxxxxx,
	// 		display name: xxxxxxx,
	// 		message text: xxxxxx xxxxx xxxxxx xxxxx
	// 	}
	// }
	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		usersEP(w, r, ctx, redisClient)
	})
	// users: {
	// 	chatroom path: xxxxxxxxxxxxx,
	// 	display name: xxxxxx
	// }
}
