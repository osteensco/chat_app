package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

var expirationSeconds = 300

func lobbyEP(w http.ResponseWriter, r *http.Request, ctx context.Context, redisClient *redis.Client, CRDBClient *pgxpool.Pool) {

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

	roomname := r.URL.Query().Get("roomname")

	if roomname == "" {
		log.Panicf("roompath query parameter not provided! Request URL provided was %v", r.URL)
	}

	log.Printf("%v room %v with path %v at lobbyEP", r.Method, roomname, roompath)

	key := "lobby"

	switch r.Method {
	case "GET":
		// used on new client entering lobby
		func(w http.ResponseWriter, r *http.Request) {

			var allChatrooms map[string]string
			var err error

			if redisKeyExists(ctx, redisClient, key) {
				allChatrooms, err = getAllChatroomsRedis(ctx, redisClient, key)
				resetRedisKeyExpiration(ctx, redisClient, key)
			} else {
				allChatrooms, err = getAllChatroomsCRDB(ctx, CRDBClient, key)
				if err == nil {
					go func() {
						log.Println("Adding lobby to cache (Redis)")
						for chatroomName, chatroomPath := range allChatrooms {
							err := addChatroomToLobbyRedis(ctx, redisClient, key, chatroomName, chatroomPath)
							if err != nil {
								log.Println("Error adding chatroom to lobby in Redis:", err)
							}
						}
						resetRedisKeyExpiration(ctx, redisClient, key)
					}()
				}
			}

			if err != nil {
				log.Panicf("Error getting all chatrooms in %v:  %v", key, err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {

				responsePayload, err := json.Marshal(allChatrooms)
				if err != nil {
					log.Panicf("Error marshaling chatMessages to JSON: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				_, err = w.Write(responsePayload)
				if err != nil {
					log.Panicf("Error writing JSON response: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			}

		}(w, r)

	case "POST":
		// used when new chatroom is created
		func(w http.ResponseWriter, r *http.Request) {

			err := addChatroomToLobbyCRDB(ctx, CRDBClient, roomname, roompath)
			if err == nil {
				err = addChatroomToLobbyRedis(ctx, redisClient, key, roomname, roompath)
			}
			if err != nil {
				log.Panicf("Error adding chatroom to %v: %v", key, err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)

		}(w, r)

	case "DELETE":
		// used when chatroom is removed
		func(w http.ResponseWriter, r *http.Request) {

			err := removeChatroomFromLobbyCRDB(ctx, CRDBClient, roomname)
			if err == nil {
				err = removeChatroomFromLobbyRedis(ctx, redisClient, key, roomname)
			}
			if err != nil {
				log.Panicf("Error removing chatroom from %v: %v", key, err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				w.WriteHeader(http.StatusOK)
			}

		}(w, r)
	}

}

func messagesEP(w http.ResponseWriter, r *http.Request, ctx context.Context, redisClient *redis.Client, CRDBClient *pgxpool.Pool) {

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

	log.Printf("%v message(s) for room %v at messagesEP", r.Method, roompath)

	var key = "messages_" + roompath

	switch r.Method {
	case "GET":
		//used to get chat history on user entering room
		func(w http.ResponseWriter, r *http.Request) {

			var chatMessages []string
			var err error

			if redisKeyExists(ctx, redisClient, key) {
				chatMessages, err = getMessageHistoryRedis(ctx, redisClient, roompath)
				resetRedisKeyExpiration(ctx, redisClient, key)
			} else {
				chatMessages, err = getMessageHistoryCRDB(ctx, CRDBClient, roompath)
				log.Printf("chatMessages: %v", chatMessages)
				log.Printf("error: %v", err)
				if err == nil {
					go func() {
						log.Printf("Adding messages in room %v to cache (Redis)", roompath)
						for _, message := range chatMessages {
							err = addMessageToHistoryRedis(ctx, redisClient, roompath, message)
							if err != nil {
								log.Panicf("Error adding message to chatroom history: %v", err)
								http.Error(w, "Internal Server Error", http.StatusInternalServerError)
								return
							}
						}
						resetRedisKeyExpiration(ctx, redisClient, key)
					}()
				}
			}

			if err != nil {
				log.Panicf("Error querying messages for roompath %v, %v", roompath, err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			responsePayload, err := json.Marshal(chatMessages)
			if err != nil {
				log.Panicf("Error marshaling chatMessages to JSON: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			_, err = w.Write(responsePayload)
			if err != nil {
				log.Panicf("Error writing JSON response: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

		}(w, r)

	case "POST":
		// used to add to chat history

		chatMessage := r.URL.Query().Get("chatmessage")
		if chatMessage == "" {
			log.Panicf("chatmessage query parameter not provided! Request URL provided was %v", r.URL)
		}

		func(w http.ResponseWriter, r *http.Request) {

			// check length, remove earliest message if applicable

			// Update Database
			var length int8
			messagesArray, err := getMessageHistoryCRDB(ctx, CRDBClient, roompath)
			if err != nil {
				log.Panicf("Error message count from %v chatroom history: %v", roompath, err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				length = int8(len(messagesArray))
				log.Printf("message history length for %v: %v", roompath, length)
				if length == 10 {
					err = removeMessageFromHistoryCRDB(ctx, CRDBClient, roompath)
					if err != nil {
						log.Panicf("Error removing message from %v chatroom history: %v", roompath, err)
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
						return
					}
				}
				err = addMessageToHistoryCRDB(ctx, CRDBClient, roompath, chatMessage)
				if err != nil {
					log.Panicf("Error adding message to %v chatroom history: %v", roompath, err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			}

			// Update Cache if applicable
			if redisKeyExists(ctx, redisClient, key) {
				length, err = getMessageHistoryLengthRedis(ctx, redisClient, roompath)
				if err != nil {
					log.Panicf("Error getting message history length from chatroom %v", roompath)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				} else {
					if length == 10 {
						err = removeMessageFromHistoryRedis(ctx, redisClient, roompath)
						if err != nil {
							log.Panicf("Error removing message from %v chatroom history: %v", roompath, err)
							http.Error(w, "Internal Server Error", http.StatusInternalServerError)
							return
						}
					}
					err = addMessageToHistoryRedis(ctx, redisClient, roompath, chatMessage)
					if err != nil {
						log.Panicf("Error adding message to %v chatroom history: %v", roompath, err)
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
						return
					}
				}
			}
			w.WriteHeader(http.StatusOK)

		}(w, r)

	case "DELETE":
		// used on chatroom removal
		func(w http.ResponseWriter, r *http.Request) {

			err := deleteKeyCRDB(ctx, CRDBClient, roompath)
			if err == nil {
				deleteKeyRedis(ctx, redisClient, roompath)
			}
			if err != nil {
				log.Panicf("Error removing chatroom history for room %v: %v", roompath, err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				w.WriteHeader(http.StatusOK)
			}

		}(w, r)
	}

}

func usersEP(w http.ResponseWriter, r *http.Request, ctx context.Context, redisClient *redis.Client, CRDBClient *pgxpool.Pool) {

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
	if displayname == "" {
		log.Panicf("displayname query parameter not provided! Request URL provided was %v", r.URL)
	}

	log.Printf("%v %v from %v at usersEP", r.Method, displayname, roompath)

	var key = "users_" + roompath

	switch r.Method {

	case "GET":
		// used when a new client is being assigned an anonymous display name
		func(w http.ResponseWriter, r *http.Request) {

			var displayNameExists bool
			var err error

			if redisKeyExists(ctx, redisClient, key) {
				displayNameExists, err = isUserInChatroomRedis(ctx, redisClient, displayname, roompath)
				resetRedisKeyExpiration(ctx, redisClient, key)
			} else {
				displayNameExists, err = isUserInChatroomCRDB(ctx, CRDBClient, displayname, roompath)
				if err == nil {
					go func() {
						log.Printf("Adding users in room %v to cache (Redis)", roompath)
						var allUsers []string
						allUsers, err = getAllUsersInChatroomCRDB(ctx, CRDBClient, roompath)
						if err != nil {
							log.Panicf("Error getting all users in chatroom from cockroachDB: %v", err)
							http.Error(w, "Internal Server Error", http.StatusInternalServerError)
							return
						}
						for _, user := range allUsers {
							err = addUserToChatroomRedis(ctx, redisClient, user, roompath)
							if err != nil {
								log.Panicf("Error adding user to chatroom %v: %v", roompath, err)
								http.Error(w, "Internal Server Error", http.StatusInternalServerError)
								return
							}
						}
						resetRedisKeyExpiration(ctx, redisClient, key)
					}()
				}
			}

			if err != nil {
				log.Panicf("Error querying Redis with displayname %v: %v", displayname, err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !displayNameExists {
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

			err := addUserToChatroomCRDB(ctx, CRDBClient, displayname, roompath)

			if err == nil && redisKeyExists(ctx, redisClient, key) {
				err = addUserToChatroomRedis(ctx, redisClient, displayname, roompath)
			}

			if err != nil {
				log.Panicf("Error adding user %v to chatroom %v: %v", displayname, roompath, err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				w.WriteHeader(http.StatusOK)
			}

		}(w, r)

	case "PUT":
		// used when a client changes their name
		func(w http.ResponseWriter, r *http.Request) {

			newname := r.URL.Query().Get("newname")
			if newname == "" {
				log.Panicf("newname query parameter not provided! Request URL provided was %v", r.URL)
			}

			log.Printf("changing name from %v to %v in %v", displayname, newname, roompath)

			err := changeUserNameCRDB(ctx, CRDBClient, displayname, newname, roompath)

			if err == nil {
				room, ok := AllRooms[roompath]
				if !ok {
					log.Panicf("Roompath %v not found in AllRooms map!", roompath)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				err = room.UpdateClientName(displayname, newname)
				if redisKeyExists(ctx, redisClient, key) {
					err = changeUserNameRedis(ctx, redisClient, displayname, newname, roompath)
				}
			}

			if err != nil {
				log.Panicf("Error changing username %v to %v in chatroom %v: %v", displayname, newname, roompath, err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				w.WriteHeader(http.StatusOK)
			}

		}(w, r)

	case "DELETE":
		// used when a client leaves a room or a room is removed from the server
		func(w http.ResponseWriter, r *http.Request) {

			err := removeUserFromChatroomCRDB(ctx, CRDBClient, displayname, roompath)

			if err == nil {
				err = removeUserFromChatroomRedis(ctx, redisClient, displayname, roompath)
			}

			if err != nil {
				log.Panicf("Error removing user %v from chatroom %v: %v", displayname, roompath, err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				w.WriteHeader(http.StatusOK)
			}

		}(w, r)

	}

}

func initAPI(ctx context.Context, redisClient *redis.Client, CRDBClient *pgxpool.Pool) {

	http.HandleFunc("/api/lobby", func(w http.ResponseWriter, r *http.Request) {
		lobbyEP(w, r, ctx, redisClient, CRDBClient)
	})
	// Redis HASH
	// lobby: {
	// 	room name: xxxxxxx,
	// 	room path: xxxxxxxxxxxxxxxxx
	// }
	http.HandleFunc("/api/messages", func(w http.ResponseWriter, r *http.Request) {
		messagesEP(w, r, ctx, redisClient, CRDBClient)
	})
	// Redis LIST
	// chatrooms: {
	// 	chatroom path:
	// 	messages: {
	// 		timestamp: xxxxxxx,
	// 		display name: xxxxxxx,
	// 		message text: xxxxxx xxxxx xxxxxx xxxxx
	// 	}
	// }
	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		usersEP(w, r, ctx, redisClient, CRDBClient)
	})
	// Redis SET
	// users: {
	// 	chatroom path: xxxxxxxxxxxxx,
	// 	display name: xxxxxx
	// }
}
