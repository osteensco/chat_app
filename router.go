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
	http.HandleFunc("/api/chatrooms", chatroomsEP)
	http.HandleFunc("/api/users", usersEP)

}
