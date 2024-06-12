package usecase

import (
	"chat-application-api/internal/app/domain/room_chat/repository"
	"chat-application-api/internal/app/dto"
	"chat-application-api/internal/app/model"
	"chat-application-api/internal/pkg/msg"
	"time"

	"github.com/labstack/echo"
)

type RoomChatUsecase interface {
	UCreateChatRoom(c echo.Context, pl *dto.CreateChatRoomRequest) error
	UJoinUserRoomChat(c echo.Context, pl *dto.AddUserRoomChatRequest) error
	UListRoomChat(c echo.Context, pl *dto.ListChatRoomRequest) ([]*dto.ListChatRoomResponse, error)
	UCreateRoomChatMessage(pl *dto.CreateRoomChatMessageRequest) (int64, error)
	UListRoomChatMessage(c echo.Context, request *dto.ListRoomChatMessageRequest) ([]*dto.ListRoomChatMessageResponse, error)
	ULeaveRoomChat(c echo.Context, pl *dto.LeaveRoomChatRequest) error
	UCheckUserIsExistsInRoomChat(c echo.Context, pl *dto.CreateRoomChatMessageRequest) (bool, error)
	UGetRoomNameChat(c echo.Context, pl *dto.GetRoomChatNameRequest) (string, error)
}

type roomChatUsecase struct {
	roomChatRepository repository.RoomChatRepository
}

func NewChatRoomUsecase(repository repository.RoomChatRepository) RoomChatUsecase {
	return &roomChatUsecase{repository}
}

func (s *roomChatUsecase) UCreateChatRoom(c echo.Context, pl *dto.CreateChatRoomRequest) error {
	dataRoomChat := &model.ChatRoom{
		RoomName:  pl.RoomName,
		UserId:    pl.UserId,
		Status:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	chatRoomId, err := s.roomChatRepository.RCreateRoomChat(c, dataRoomChat)
	if err != nil {
		return err
	}

	dataRoomChatUser := &model.RoomChatUser{
		UserId:     pl.UserId,
		RoomChatId: chatRoomId,
		CreatedAt:  time.Now(),
	}

	err = s.roomChatRepository.RCreateRoomChatUser(c, dataRoomChatUser)
	if err != nil {
		return err
	}

	return nil
}

func (s *roomChatUsecase) UJoinUserRoomChat(c echo.Context, pl *dto.AddUserRoomChatRequest) error {
	dataRoomChatUser := &model.RoomChatUser{
		UserId:     pl.UserId,
		RoomChatId: pl.ChatRoomId,
		CreatedAt:  time.Now(),
	}

	// Check the user whether registered in this room chat or not
	isExists, err := s.roomChatRepository.RCheckUserIsExistsInRoomChat(dataRoomChatUser)
	if err != nil {
		return err
	}

	if isExists {
		err = msg.ErrUserAlreadyExistsInThisRoomChat
		return err
	}

	err = s.roomChatRepository.RCreateRoomChatUser(c, dataRoomChatUser)
	if err != nil {
		return err
	}

	return nil
}

func (s *roomChatUsecase) UListRoomChat(c echo.Context, pl *dto.ListChatRoomRequest) ([]*dto.ListChatRoomResponse, error) {
	data := &model.RoomChatUser{
		UserId: pl.UserId,
	}

	result, err := s.roomChatRepository.RListRoomChat(c, data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *roomChatUsecase) UCreateRoomChatMessage(pl *dto.CreateRoomChatMessageRequest) (int64, error) {
	var lastInsertId int64

	// Check the user whether registered in this room chat or not
	isExists, err := s.roomChatRepository.RCheckUserIsExistsInRoomChat(&model.RoomChatUser{
		UserId:     pl.UserId,
		RoomChatId: pl.RoomChatId,
	})
	if err != nil {
		return lastInsertId, err
	}

	if !isExists {
		err = msg.ErrUserHaveNotJoinThisRoom
		return lastInsertId, err
	}

	data := &model.RoomChatMessage{
		UserId:     pl.UserId,
		RoomChatId: pl.RoomChatId,
		Message:    pl.Message,
		CreatedAt:  time.Now(),
	}

	lastInsertId, err = s.roomChatRepository.RCreateRoomChatMessage(data)
	if err != nil {
		return lastInsertId, err
	}

	return lastInsertId, nil
}

func (s *roomChatUsecase) UListRoomChatMessage(c echo.Context, pl *dto.ListRoomChatMessageRequest) ([]*dto.ListRoomChatMessageResponse, error) {
	data := &model.RoomChatMessage{
		UserId:     pl.UserId,
		RoomChatId: pl.RoomChatId,
	}

	result, err := s.roomChatRepository.RListRoomChatMessage(data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *roomChatUsecase) ULeaveRoomChat(c echo.Context, pl *dto.LeaveRoomChatRequest) error {
	if err := s.roomChatRepository.RLeaveRoomChat(c, &model.RoomChatMessage{
		UserId:     pl.UserId,
		RoomChatId: pl.RoomChatId,
	}); err != nil {
		return err
	}

	return nil
}

func (s *roomChatUsecase) UCheckUserIsExistsInRoomChat(c echo.Context, pl *dto.CreateRoomChatMessageRequest) (bool, error) {
	isExists, err := s.roomChatRepository.RCheckUserIsExistsInRoomChat(&model.RoomChatUser{
		UserId:     pl.UserId,
		RoomChatId: pl.RoomChatId,
	})
	if err != nil {
		return isExists, err
	}

	return isExists, nil
}

func (s *roomChatUsecase) UGetRoomNameChat(c echo.Context, pl *dto.GetRoomChatNameRequest) (string, error) {
	roomName, err := s.roomChatRepository.RGetRoomNameChat(c, &model.RoomChatMessage{
		RoomChatId: pl.RoomChatId,
	})
	if err != nil {
		return roomName, err
	}

	return roomName, nil
}
