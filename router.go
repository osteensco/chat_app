package main

import (
	"log"
	"net/http"
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

func usersEP(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
		}
	}()

	roompath := r.URL.Query().Get("roompath")
	if roompath == "" {
		log.Panicf("roompath query parameter not provided! Request URL provided was %v", r.URL)
	}

	displayname := r.URL.Query().Get("displayname")

	switch r.Method {

	case "GET":
		// used when a new client enters a room
		func(w http.ResponseWriter, r *http.Request) {
			if displayname != "" {
				log.Printf("GET %v FROM usersEP", displayname)
			} else {
				log.Printf("GET %v FROM usersEP", displayname)
			}
		}(w, r)

	case "POST":
		// used when a new client enters a room
		func(w http.ResponseWriter, r *http.Request) {
			if displayname != "" {
				log.Printf("POST %v FROM usersEP", displayname)
			} else {
				log.Printf("POST %v FROM usersEP", displayname)
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

func initAPI() {

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
	http.HandleFunc("/api/users", usersEP)
	// users: {
	// 	chatroom path: xxxxxxxxxxxxx,
	// 	display name: xxxxxx
	// }
}
