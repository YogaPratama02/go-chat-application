package wsroomchat

import (
	"bytes"
	"chat-application-api/internal/app/domain/room_chat/usecase"
	"chat-application-api/internal/app/dto"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

type Message struct {
	SenderId   int64  `json:"sender_id"`
	SenderName string `json:"sender_name"`
	RoomId     int64  `json:"roomId"`
	Content    string `json:"content"`
	Data       []byte `json:"data"`
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 5 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub *Hub

	// The websocket connection.
	Conn *websocket.Conn

	Send            chan Message
	RoomId          int64  `json:"roomId"`
	SenderId        int64  `json:"user_id"`
	SenderName      string `json:"sender_name"`
	EchoContext     echo.Context
	RoomChatUsecase usecase.RoomChatUsecase
}

func (c *Client) ReadMessage() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			log.Printf("error: %v", err)
			break
		}
		msg.SenderId = c.SenderId
		msg.SenderName = c.SenderName
		msg.RoomId = c.RoomId

		c.Hub.Broadcast <- msg

		_, err = c.RoomChatUsecase.UCreateRoomChatMessage(&dto.CreateRoomChatMessageRequest{
			RoomChatId: c.RoomId,
			UserId:     c.SenderId,
			Message:    msg.Content,
		})
		if err != nil {
			log.Printf("error create room chat message: %v", err)
			break
		}
	}
}

func (c *Client) WriteMessage() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			reqBodyBytes := new(bytes.Buffer)

			if err := json.NewEncoder(reqBodyBytes).Encode(message); err != nil {
				return
			}
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				fmt.Println("hub close the channel")
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.Conn.WriteMessage(websocket.TextMessage, reqBodyBytes.Bytes())
			if err != nil {
				return
			}
			c.Conn.SetWriteDeadline(time.Time{})

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
			c.Conn.SetWriteDeadline(time.Time{})
		}
	}
}
