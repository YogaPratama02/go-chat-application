package model

import "time"

type RoomChatUser struct {
	Id         int64     `json:"id"`
	UserId     int64     `json:"user_id"`
	RoomChatId int64     `json:"room_chat_id"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}
