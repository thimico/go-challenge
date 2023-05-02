package jobsity

import (
	"github.com/streadway/amqp"
	"golang.org/x/net/websocket"
)

type Room struct {
	Name      string
	Clients   map[*websocket.Conn]bool
	Broadcast chan []byte
	Quit      chan bool
	rabbit    *amqp.Connection
}

func NewRoom(name string, rabbit *amqp.Connection) *Room {
	return &Room{
		Name:      name,
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan []byte),
		Quit:      make(chan bool),
		rabbit:    rabbit,
	}
}
