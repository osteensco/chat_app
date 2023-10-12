package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

	// roomname := strings.ReplaceAll(r.URL.Query().Get("roomname"), " ", "%20")
	roomname := r.URL.Query().Get("roomname")
	if roomname == "" {
		log.Panicf("roompath query parameter not provided! Request URL provided was %v", r.URL)
	}

	key := "lobby"

	log.Printf("%v room %v with path %v at lobbyEP", r.Method, roomname, roompath)

	switch r.Method {
	case "GET":
		// used on new client entering lobby
		func(w http.ResponseWriter, r *http.Request) {

			allChatrooms, err := getAllChatroomsRedis(ctx, redisClient, key)
			if err != nil || len(allChatrooms) == 0 {
				allChatrooms, err = getAllChatroomsCRDB(ctx, CRDBClient, key)
			}
			if err != nil {
				log.Panicf("Error getting all chatromms in %v from Redis %v", roompath, http.StatusInternalServerError)
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

			jsonroompath, err := json.Marshal(roompath)
			if err != nil {
				log.Panicf("error marchaling JSON: %v", err)
			}

			room := map[string]interface{}{roomname: fmt.Sprintf("%v", string(jsonroompath))}
			err = addChatroomToLobbyRedis(ctx, redisClient, key, room)
			if err == nil {
				addChatroomToLobbyCRDB(ctx, CRDBClient, roomname, roompath)
			}
			if err != nil {
				log.Panicf("Error adding chatroom to %v in Redis %v", key, http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
			}

		}(w, r)

	case "DELETE":
		// used when chatroom is removed
		func(w http.ResponseWriter, r *http.Request) {

			err := removeChatroomFromLobbyRedis(ctx, redisClient, key, roomname)
			if err == nil {
				removeChatroomFromLobbyCRDB(ctx, CRDBClient, roomname)
			}
			if err != nil {
				log.Panicf("Error removing chatroom from %v in Redis %v", key, http.StatusInternalServerError)
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

	switch r.Method {
	case "GET":
		//used to get chat history on user entering room
		func(w http.ResponseWriter, r *http.Request) {

			chatMessages, err := getMessageHistoryRedis(ctx, redisClient, roompath)

			if err != nil || len(chatMessages) == 0 {
				chatMessages, err = getMessageHistoryCRDB(ctx, CRDBClient, roompath)

				if err != nil {
					log.Panicf("Error querying Redis with roompath %v %v", roompath, http.StatusInternalServerError)
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
			length, err := getMessageHistoryLengthRedis(ctx, redisClient, roompath)

			if err != nil || length == 0 {

				var messagesArray []string
				messagesArray, err = getMessageHistoryCRDB(ctx, CRDBClient, roompath)

				if err != nil {
					log.Panicf("Error message count from chatroom history %v", http.StatusInternalServerError)
				} else {
					length = int64(len(messagesArray))
				}

			}

			if length == 10 {
				err = removeMessageFromHistoryRedis(ctx, redisClient, roompath)
				if err == nil {
					err = removeMessageFromHistoryCRDB(ctx, CRDBClient, roompath)
				}
				if err != nil {
					log.Panicf("Error removing message from chatroom history %v", http.StatusInternalServerError)
				}
			}

			err = addMessageToHistoryRedis(ctx, redisClient, roompath, chatMessage)
			if err == nil {
				err = addMessageToHistoryCRDB(ctx, CRDBClient, roompath, chatMessage)
			}
			if err != nil {
				log.Panicf("Error adding message to chatroom history %v", http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
			}

		}(w, r)

	case "DELETE":
		// used on chatroom removal
		func(w http.ResponseWriter, r *http.Request) {

			err := deleteKeyRedis(ctx, redisClient, roompath)
			if err == nil {
				deleteKeyCRDB(ctx, CRDBClient, roompath)
			}
			if err != nil {
				log.Panicf("Error removing chatroom history for room %v, %v", roompath, http.StatusInternalServerError)
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

	switch r.Method {

	case "GET":
		// used when a new client is being assigned an anonymous display name
		func(w http.ResponseWriter, r *http.Request) {

			displayNameExists, err := isUserInChatroomRedis(ctx, redisClient, displayname, roompath)

			if err != nil {
				displayNameExists, err = isUserInChatroomCRDB(ctx, CRDBClient, displayname, roompath)
				if err != nil {
					log.Panicf("Error querying Redis with displayname %v %v", displayname, http.StatusInternalServerError)
				}
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

			err := addUserToChatroomRedis(ctx, redisClient, displayname, roompath)

			if err == nil {
				err = addUserToChatroomCRDB(ctx, CRDBClient, displayname, roompath)
			}

			if err != nil {
				log.Panicf("Error adding user to chatroom %v, %v", http.StatusInternalServerError, err)
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

			err := changeUserNameRedis(ctx, redisClient, displayname, newname, roompath)

			if err == nil {
				changeUserNameCRDB(ctx, CRDBClient, displayname, newname, roompath)
			}

			if err != nil {
				log.Panicf("Error changing username %v to %v in chatroom %v", displayname, newname, http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
			}

		}(w, r)

	case "DELETE":
		// used when a client leaves a room or a room is removed from the server
		func(w http.ResponseWriter, r *http.Request) {

			err := removeUserFromChatroomRedis(ctx, redisClient, displayname, roompath)

			if err == nil {
				err = removeUserFromChatroomCRDB(ctx, CRDBClient, displayname, roompath)
			}

			if err != nil {
				log.Panicf("Error removing user from chatroom %v", http.StatusInternalServerError)
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
