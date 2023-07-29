package main

import (
	"fmt"
	"io"

	"github.com/aofei/air"
	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	// creates a server and opens websocket connections among clients
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Printf("incoming connection from client: %v", ws.RemoteAddr())

	s.conns[ws] = true

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read error: ", err)
			continue
		}
		msg := buf[:n]
		fmt.Println(string(msg))

	}
}

func indexHandler(req *air.Request, res *air.Response) error {
	println("yes")
	return res.Render(nil, "/static/index.html")
}

func main() {

	// server := NewServer()

	app := air.New()

	app.GET("/", indexHandler)

	app.Serve()
}
