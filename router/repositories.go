package router

import (
	_authRepository "chat-application-api/internal/app/domain/auth/repository"
	_directMessageRepository "chat-application-api/internal/app/domain/direct_message/repository"
	_roomChatRepository "chat-application-api/internal/app/domain/room_chat/repository"
	"database/sql"
)

type repositories struct {
	AuthRepository          _authRepository.AuthRepository
	RoomChatRepository      _roomChatRepository.RoomChatRepository
	DirectMessageRepository _directMessageRepository.DirectMessageRepository
}

func newRepositories(db *sql.DB) repositories {
	authRepository := _authRepository.NewAuthRepository(db)
	roomChatRepository := _roomChatRepository.NewChatRoomRepository(db)
	directMessageRepository := _directMessageRepository.NewDirectMessageRepository(db)

	return repositories{
		AuthRepository:          authRepository,
		RoomChatRepository:      roomChatRepository,
		DirectMessageRepository: directMessageRepository,
	}
}
