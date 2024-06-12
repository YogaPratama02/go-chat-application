package router

import (
	_authUsecase "chat-application-api/internal/app/domain/auth/usecase"
	_directMessageUsecase "chat-application-api/internal/app/domain/direct_message/usecase"
	_roomChatUsecase "chat-application-api/internal/app/domain/room_chat/usecase"
)

type usecases struct {
	AuthUsecase          _authUsecase.AuthUsecase
	RoomChatUsecase      _roomChatUsecase.RoomChatUsecase
	DirectMessageUsecase _directMessageUsecase.DirectMessageUsecase
}

func newUsecases(repo repositories) usecases {
	authUsecase := _authUsecase.NewAuthUsecase(repo.AuthRepository)
	roomChatUsecase := _roomChatUsecase.NewChatRoomUsecase(repo.RoomChatRepository)
	directMessageUsecase := _directMessageUsecase.NewDirectMessageUsecase(repo.DirectMessageRepository)

	return usecases{
		AuthUsecase:          authUsecase,
		RoomChatUsecase:      roomChatUsecase,
		DirectMessageUsecase: directMessageUsecase,
	}
}
