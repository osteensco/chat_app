package main

import (
	"log"
	"net/http"
)

func initAPI() {
	http.HandleFunc("/api/roombuilder", func(w http.ResponseWriter, r *http.Request) {
		// r.Method to determine type of request i.e. GET, POST, PUT, DELETE
		// probably need all types of requests handled
		// use a switch case
		route := r.URL.Path
		log.Println(route)
	})

	http.HandleFunc("/api/chatrooms", func(w http.ResponseWriter, r *http.Request) {
		//GET
		//GET {ID}
		route := r.URL.Path
		log.Println(route)
	})

	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		//GET
		//GET {ID}
		route := r.URL.Path
		log.Println(route)
	})

}
