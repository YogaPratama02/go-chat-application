package model

import "time"

type DirectMessage struct {
	Id         int64     `json:"id"`
	SenderId   int64     `json:"sender_id"`
	ReceiverId int64     `json:"receiver_Id"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}
