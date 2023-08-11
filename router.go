package main

import (
	"log"
	"net/http"
)

func lobbyEP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method)
		}(w, r)
	case "POST":
		func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method)
		}(w, r)
	case "DELETE":
		func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method)
		}(w, r)
	}

}

func chatroomsEP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method)
		}(w, r)
	case "POST":
		func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method)
		}(w, r)
	case "PUT":
		func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method)
		}(w, r)
	case "DELETE":
		func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method)
		}(w, r)
	}

}

func usersEP(w http.ResponseWriter, r *http.Request) {

	// TODO
	// Need parameter for room path

	switch r.Method {
	case "GET":
		func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method)
		}(w, r)
	case "POST":
		func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method)
		}(w, r)
	case "PUT":
		func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method)
		}(w, r)
	case "DELETE":
		func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method)
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
