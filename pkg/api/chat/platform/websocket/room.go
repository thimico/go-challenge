package websocket

import (
	"fmt"
	"golang.org/x/net/websocket"
	jobsity "my-chat-jobsity-challenge"
)

// Room represents the client for room websocket
type Room struct{}

func (r *Room) Run(room *jobsity.Room) {
	for {
		select {
		case message := <-room.Broadcast:
			for client := range room.Clients {
				err := websocket.Message.Send(client, string(message))
				if err != nil {
					fmt.Printf("error sending message to client: %v", err)
					r.RemoveClient(client, room)
				}
			}
		case <-room.Quit:
			return
		}
	}
}

func (r *Room) AddClient(conn *websocket.Conn, room *jobsity.Room) {
	room.Clients[conn] = true
}

func (r *Room) RemoveClient(conn *websocket.Conn, room *jobsity.Room) {
	delete(room.Clients, conn)
}

func (r *Room) BroadcastMessage(message []byte, sender *websocket.Conn, room *jobsity.Room) {
	for client := range room.Clients {
		if client != sender {
			err := websocket.Message.Send(client, string(message))
			if err != nil {
				fmt.Printf("error sending message to client: %v", err)
				r.RemoveClient(client, room)
			}
		}
	}
}
