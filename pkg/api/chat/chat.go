package chat

import (
	"encoding/csv"
	"fmt"
	"github.com/labstack/echo"
	"golang.org/x/net/websocket"
	jobsity "my-chat-jobsity-challenge"
	"net/http"
	"strconv"
	"strings"
)

func (s *Chat) JoinRoom(c echo.Context, conn *websocket.Conn, roomName string, user jobsity.AuthUser) error {
	// Look up the room in the map
	room, ok := s.Rooms[roomName]
	if !ok {
		return fmt.Errorf("room %s not found", roomName)
	}

	// Add the client to the room
	//room.AddClient(conn, user)
	s.ws.AddClient(conn, room)
	return nil
}

func (s *Chat) LeaveRoom(c echo.Context, roomName string, conn *websocket.Conn, user *jobsity.AuthUser) error {
	if room, ok := s.Rooms[roomName]; ok {
		s.ws.RemoveClient(conn, room)
		s.sendMessageToRoom(roomName, fmt.Sprintf("%s left the room", user.Username), nil)
		return nil
	}
	return fmt.Errorf("room %s not found", roomName)
}

func (s *Chat) CreateRoom(c echo.Context, roomName string, user *jobsity.AuthUser) error {
	// Only admins can create new rooms
	if user.Role != jobsity.AdminRole {
		return fmt.Errorf("not authorized")
	}

	// Check if the room already exists
	if _, ok := s.Rooms[roomName]; ok {
		return fmt.Errorf("room %s already exists", roomName)
	}

	// Create a new room and start it
	room := jobsity.NewRoom(roomName, s.rabbit)
	s.Rooms[roomName] = room
	go s.ws.Run(room)

	return nil
}

func (s *Chat) DeleteRoom(c echo.Context, roomName string, user *jobsity.AuthUser) error {
	// Only admins can delete rooms
	if user.Role != jobsity.AdminRole {
		return fmt.Errorf("not authorized")
	}

	// Check if the room exists
	if room, ok := s.Rooms[roomName]; ok {
		// Stop the room's message broadcasting loop
		room.Quit <- true
		delete(s.Rooms, roomName)
		return nil
	}
	return fmt.Errorf("room %s not found", roomName)
}

func (s *Chat) SendMessageOne(c echo.Context, roomName string, user *jobsity.AuthUser, message string) error {
	if _, ok := s.Rooms[roomName]; ok {
		msg := fmt.Sprintf("%s: %s", user.Username, message)
		s.sendMessageToRoom(roomName, msg, user)
		return nil
	}
	return fmt.Errorf("room %s not found", roomName)
}

func (s *Chat) HandleCommand(c echo.Context, conn *websocket.Conn, roomName string, message string) error {
	if strings.HasPrefix(message, "/join ") {
		// Join the specified room
		room := strings.TrimPrefix(message, "/join ")
		return s.JoinRoom(c, conn, room, jobsity.AuthUser{})
	} else if strings.HasPrefix(message, "/leave") {
		// Leave the current room
		return s.LeaveRoom(c, roomName, conn, nil)
	} else if strings.HasPrefix(message, "/users") {
		// Get the list of users in the current room
		users, err := s.GetUsersInRoom(c, roomName)
		if err != nil {
			return err
		}
		// Construct a message with the list of users
		var msg string
		if len(users) == 0 {
			msg = "There are no users in this room."
		} else {
			msg = "Users in this room: " + strings.Join(users, ", ")
		}
		// Send the message to the client
		err = websocket.Message.Send(conn, msg)
		if err != nil {
			return err
		}
		return nil
	} else if strings.HasPrefix(message, "/stock=") {
		// Handle stock command
		stockCode := strings.TrimPrefix(message, "/stock=")
		go s.handleStockCommand(conn, roomName, nil, stockCode)
		return nil
	} else if strings.HasPrefix(message, "/create ") {
		// Create a new room (admin only)
		//if r.RoleID < jobsity.SuperAdminRole || r.RoleID > jobsity.UserRole {
		//	return jobsity.ErrBadRequest
		//}
		//err := s.RBAC.IsLowerRole(c, jobsity.UserRole)
		//if err != nil {
		//	return err
		//}
		roomName := strings.TrimPrefix(message, "/create ")
		return s.CreateRoom(c, roomName, nil)
	} else {
		// Unrecognized command
		return fmt.Errorf("unrecognized command: %s", message)
	}
}

func (s *Chat) GetUsersInRoom(c echo.Context, roomName string) ([]string, error) {
	// Get the list of clients in the room
	clients, ok := s.clients[roomName]
	if !ok {
		return nil, fmt.Errorf("room %s not found", roomName)
	}

	// Extract the list of usernames from the clients
	var usernames []string
	for _, client := range clients {
		usernames = append(usernames, client.username)
	}

	return usernames, nil
}

func (s *Chat) HandleMessage(ws *websocket.Conn, room string, message string) {
	if strings.HasPrefix(message, "/stock=") {
		// Handle stock command
		stockCode := strings.TrimPrefix(message, "/stock=")
		go s.handleStockCommand(ws, room, nil, stockCode)
	} else {
		// Broadcast the message to all clients in the room
		roomFrom, ok := s.Rooms[room]
		if !ok {
			return
		}

		s.ws.BroadcastMessage([]byte(message), ws, roomFrom)

	}
}

func (s *Chat) SendMessage(roomName string, user *jobsity.AuthUser, message string) error {
	if room, ok := s.Rooms[roomName]; ok {
		// Handle stock command messages
		if strings.HasPrefix(message, "/stock=") {
			stockCode := strings.TrimPrefix(message, "/stock=")
			err := s.handleStockCommand(nil, roomName, user, stockCode)
			if err != nil {
				return err
			}
			return nil
		}

		// Broadcast regular messages
		msg := fmt.Sprintf("%s: %s", user.Username, message)
		s.ws.BroadcastMessage([]byte(msg), nil, room)
		return nil
	}
	return fmt.Errorf("room %s not found", roomName)
}

func (s *Chat) handleStockCommand(ws *websocket.Conn, roomName string, user *jobsity.AuthUser, stockCode string) error {
	if room, ok := s.Rooms[roomName]; ok {
		return s.handleStockCommandInRoom(ws, room, user, stockCode)
	}
	return nil
}

func (s *Chat) handleStockCommandInRoom(ws *websocket.Conn, room *jobsity.Room, user *jobsity.AuthUser, stockCode string) error {
	// Call the stock API to get the stock quote
	stockQuote, err := s.FetchStockQuote(stockCode)
	if err != nil {
		return fmt.Errorf("failed to fetch stock quote for %s: %v", stockCode, err)
	}

	// Send the stock quote to the room
	msg := fmt.Sprintf("%s quote is $%.2f per share", stockCode, stockQuote)
	//&jobsity.AuthUser{Username: "Bot"}
	s.ws.BroadcastMessage([]byte(msg), ws, room)
	return nil
}

func (s *Chat) sendMessageToRoom(roomName string, message string, sender *jobsity.AuthUser) {
	if room, ok := s.Rooms[roomName]; ok {
		if sender != nil {
			fmt.Printf("[%s] %s: %s", roomName, sender.Username, message)
		} else {
			fmt.Printf("[%s] %s", roomName, message)
		}
		room.Broadcast <- []byte(message)
	}
}

func (s *Chat) FetchStockQuote(stockCode string) (float64, error) {
	// Fetch the stock data from the API
	resp, err := http.Get(fmt.Sprintf("https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv", stockCode))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Parse the CSV data
	reader := csv.NewReader(resp.Body)
	record, err := reader.Read()
	if err != nil {
		return 0, err
	}

	// Extract the last closing price
	price, err := strconv.ParseFloat(record[6], 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}
