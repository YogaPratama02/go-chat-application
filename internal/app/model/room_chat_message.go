package model

import "time"

type RoomChatMessage struct {
	Id         int64     `json:"id"`
	UserId     int64     `json:"user_id"`
	RoomChatId int64     `json:"room_chat_id"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}
