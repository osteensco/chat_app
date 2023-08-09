package main

import (
	"log"
	"net/http"
)

func initAPI() {
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		route := r.URL.Path
		log.Println(route)
	})
}
