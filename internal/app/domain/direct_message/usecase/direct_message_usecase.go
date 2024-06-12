package usecase

import (
	"chat-application-api/internal/app/domain/direct_message/repository"
	"chat-application-api/internal/app/dto"
	"chat-application-api/internal/app/model"
	"log"
	"time"
)

type DirectMessageUsecase interface {
	UCreateDirectMessage(pl *dto.CreateDirectMessageRequest) error
}

type directMessageUsecase struct {
	directMessageRepository repository.DirectMessageRepository
}

func NewDirectMessageUsecase(repository repository.DirectMessageRepository) DirectMessageUsecase {
	return &directMessageUsecase{repository}
}

func (s *directMessageUsecase) UCreateDirectMessage(pl *dto.CreateDirectMessageRequest) error {
	// Check whether sender is ever chat receiver or not
	now := time.Now()
	var idDirectMessage int64
	dataDirectMessage := &model.DirectMessage{
		SenderId:   pl.SenderId,
		ReceiverId: pl.ReceiverId,
		CreatedAt:  now,
	}
	isEver, err := s.directMessageRepository.RCheckWhetherSenderEverChatReceiverOrNot(dataDirectMessage)
	if err != nil {
		return err
	}

	switch isEver {
	case true:
		idDirectMessage, err = s.directMessageRepository.RGetIdDirectMessageBetweenSenderAndReceiver(dataDirectMessage)
		if err != nil {
			return err
		}
	default:
		idDirectMessage, err = s.directMessageRepository.RCreateDirectMessage(dataDirectMessage)
		if err != nil {
			return err
		}
	}

	// Create direct message content
	lastIdDirectMessageContent, err := s.directMessageRepository.RCreateDirectMessageDetail(&model.DirectMessageDetail{
		Content:         pl.Content,
		UserId:          pl.SenderId,
		DirectMessageId: idDirectMessage,
		CreatedAt:       now,
	})
	if err != nil {
		return err
	}

	log.Println("SUCCESSFULLY CREATE DM WITH ID DM DETAIL", lastIdDirectMessageContent)
	return nil
}
