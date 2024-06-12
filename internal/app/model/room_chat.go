package model

import "time"

type ChatRoom struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"user_id"`
	RoomName  string    `json:"room_name"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
