package model

import "time"

type DirectMessageDetail struct {
	Id              int64     `json:"id"`
	DirectMessageId int64     `json:"direct_message_id"`
	UserId          int64     `json:"user_id"`
	Content         string    `json:"content"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
}
