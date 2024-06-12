package dto

import "time"

type CreateChatRoomRequest struct {
	UserId   int64  `json:"user_id" validate:"required"`
	RoomName string `json:"room_name" validate:"required"`
}

type AddUserRoomChatRequest struct {
	ChatRoomId int64 `json:"chat_room_id" validate:"required"`
	UserId     int64 `json:"user_id" validate:"required"`
}

type ListChatRoomRequest struct {
	UserId int64 `json:"user_id" validate:"required"`
}

type ListChatRoomResponse struct {
	RoomChatId int64  `json:"room_chat_id"`
	RoomName   string `json:"room_name"`
}

type CreateRoomChatMessageRequest struct {
	RoomChatId int64  `json:"room_chat_id" validate:"required"`
	UserId     int64  `json:"user_id" validate:"required"`
	Message    string `json:"message" validate:"required"`
}

type ListRoomChatMessageRequest struct {
	RoomChatId int64 `json:"room_chat_id" validate:"required"`
	UserId     int64 `json:"user_id" validate:"required"`
}

type ListRoomChatMessageResponse struct {
	RoomChatMessageId int64     `json:"room_chat_message_id"`
	RoomChatId        int64     `json:"room_chat_id"`
	UserId            int64     `json:"user_id"`
	UserName          string    `json:"user_name"`
	Message           string    `json:"message"`
	IsUser            bool      `json:"is_user"`
	CreatedAt         time.Time `json:"created_at"`
}

type LeaveRoomChatRequest struct {
	RoomChatId int64 `json:"room_chat_id" validate:"required"`
	UserId     int64 `json:"user_id" validate:"required"`
}

type GetRoomChatNameRequest struct {
	RoomChatId int64 `json:"room_chat_id" validate:"required"`
}
