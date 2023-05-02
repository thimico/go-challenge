package chat

import (
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
	"github.com/streadway/amqp"
	"golang.org/x/net/websocket"
	jobsity "my-chat-jobsity-challenge"
	websocket2 "my-chat-jobsity-challenge/pkg/api/chat/platform/websocket"
	"time"
)

// Service represents chat application interface
type Service interface {
	HandleCommand(c echo.Context, conn *websocket.Conn, roomName string, message string) error
	JoinRoom(c echo.Context, conn *websocket.Conn, roomName string, user *jobsity.AuthUser) error
	LeaveRoom(c echo.Context, roomName string, conn *websocket.Conn, user *jobsity.AuthUser) error
	SendMessage(c echo.Context, roomName string, message []byte, user *jobsity.AuthUser) error
	SendMessageOne(c echo.Context, roomName string, user *jobsity.AuthUser, message string) error
	GetUsersInRoom(c echo.Context, roomName string) ([]string, error)
	FetchStockQuote(c echo.Context, stockCode string) (float64, error)
	CreateRoom(c echo.Context, roomName string, user *jobsity.AuthUser) error
}

// New creates new password application service
func New(rooms []string /*map[string]*jobsity.Room,*/, db *pg.DB, rabbit *amqp.Connection, ws RWS) Chat {
	return Chat{
		rooms:  rooms,
		db:     db,
		rabbit: rabbit,
		ws:     ws,
	}
}

// Initialize initalizes password application service with defaults
func Initialize(rooms []string /*map[string]*jobsity.Room,*/, db *pg.DB, rabbit *amqp.Connection) Chat {
	return New(rooms, db, rabbit, &websocket2.Room{})
}

type client struct {
	conn     *websocket.Conn
	username string
	lastMsg  time.Time
}

// Chat represents chat application service
type Chat struct {
	clients map[string][]*client
	Rooms   map[string]*jobsity.Room
	rooms   []string
	db      *pg.DB
	rabbit  *amqp.Connection
	ws      RWS
}

// RWS represents room websocket interface
type RWS interface {
	Run(*jobsity.Room)
	AddClient(*websocket.Conn, *jobsity.Room)
	RemoveClient(*websocket.Conn, *jobsity.Room)
	BroadcastMessage([]byte, *websocket.Conn, *jobsity.Room)
	//BroadcastMessageInRoom([]byte, *websocket.Conn, string)
	//handleStockCommand([]byte, *websocket.Conn, string)
	//handleStockCommandInRoom([]byte, *websocket.Conn, *jobsity.Room)
}
