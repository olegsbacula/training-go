package main

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)



type Server struct {
	conns map[*websocket.Conn]bool
}


func NewServer() *Server {
	return &Server{conns: make(map[*websocket.Conn]bool)}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("new socket connection:", ws.RemoteAddr())
	s.conns[ws] = true 
	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	fmt.Println("new socket connection:", ws.RemoteAddr())

	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			return 
		}
		msg := buf[:n]
		s.broadcast(msg)
	}
}

func (s *Server) broadcast (b []byte){
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("write error: ", err)
			}
		}(ws)
	}
}




func main_socket() { // do not forget to add origin if using postman headers -> origin -> http://localhost:8080	
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.ListenAndServe(":8080",nil)
}
