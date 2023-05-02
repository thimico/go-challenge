package transport

import (
	"github.com/labstack/echo"
	"golang.org/x/net/websocket"
	"log"
	"my-chat-jobsity-challenge/pkg/api/chat"
	"strings"
)

// HTTP represents chat http service
type HTTP struct {
	svc chat.Service
}

// NewHTTP creates new user http service
func NewHTTP(svc chat.Service, r *echo.Group, mw echo.MiddlewareFunc) {
	h := HTTP{svc}
	ur := r.Group("/chat")

	// swagger:route GET /me auth meReq
	// Gets user's info from session.
	// responses:
	//  200: userResp
	//  500: err
	ur.GET("/ws", h.handleWebSocket, mw)

}

func (h *HTTP) handleWebSocket(c echo.Context) error {
	// Upgrade the HTTP request to a WebSocket connection
	wsHandler := websocket.Handler(func(ws *websocket.Conn) {
		// Read the initial message from the WebSocket
		var msg string
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			log.Println("Error receiving initial message:", err)
			return
		}

		// Parse the message as a join command
		if strings.HasPrefix(msg, "/join ") {
			// Extract the room name from the command
			room := strings.TrimPrefix(msg, "/join ")

			// Join the room
			err := h.svc.JoinRoom(c, ws, room, nil)
			if err != nil {
				log.Println("Error joining room:", err)
				return
			}

			// Send a welcome message to the client
			welcomeMsg := "Welcome to the " + room + " chat room!"
			err = websocket.Message.Send(ws, welcomeMsg)
			if err != nil {
				log.Println("Error sending welcome message:", err)
				return
			}

			// Handle WebSocket events
			for {
				// Read a message from the WebSocket
				var msg string
				err = websocket.Message.Receive(ws, &msg)
				if err != nil {
					log.Println("Error receiving message:", err)
					break
				}

				// Handle the message
				h.svc.HandleCommand(c, ws, room, msg)
			}

			// Leave the room when the WebSocket connection is closed
			h.svc.LeaveRoom(c, room, ws, nil)
		}
	})
	wsHandler.ServeHTTP(c.Response().Writer, c.Request())

	return nil
}
