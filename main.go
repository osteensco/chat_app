package main

import (
	"log"
	"net/http"
)

// type Server struct {
// 	conns map[*websocket.Conn]bool
// }

// func NewServer() *Server {
// 	// creates a server and opens websocket connections among clients
// 	return &Server{
// 		conns: make(map[*websocket.Conn]bool),
// 	}
// }

func main() {

	http.Handle("/", http.FileServer(http.Dir("./static")))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
