package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

const (
	socketBufferSize = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize:socketBufferSize, WriteBufferSize:messageBufferSize}

type room struct {
	forward chan []byte
	join    chan *client
	leave   chan *client
	clients map[*client]bool
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case c := <-r.join:
			r.clients[c] = true
		case c := <-r.leave:
			delete(r.clients, c)
			close(c.send)
		case msg := <-r.forward:
			for c := range r.clients {
				c.send <- msg
			}
		}
	}
}

func (r *room) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}

	r.join <- client
	defer func() {
		r.leave <- client
	}()
	go client.write()
	client.read()

}
