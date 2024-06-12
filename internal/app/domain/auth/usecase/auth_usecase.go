package usecase

import (
	"chat-application-api/internal/app/domain/auth/repository"
	"chat-application-api/internal/app/dto"
	"chat-application-api/internal/app/model"
	"chat-application-api/internal/pkg/msg"
	"chat-application-api/internal/pkg/utill"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	URegister(c echo.Context, pl *dto.RegisterRequest) error
	ULogin(c echo.Context, pl *dto.LoginRequest) (*dto.LoginResponse, error)
	UCheckUserIfIsExists(c echo.Context, pl *dto.CheckUserIfIsExists) (*model.User, error)
	UGetUserList(c echo.Context, pl *dto.GetUserListRequest) ([]*dto.GetUserListResponse, error)
	UGetSenderReceiverName(c echo.Context, pl *dto.GetSenderReceiverNameRequest) (*dto.GetSenderReceiverNameResponse, error)
	UGetSenderName(c echo.Context, pl *dto.GetSenderReceiverNameRequest) (*dto.GetSenderReceiverNameResponse, error)
}

type authUsecase struct {
	authRepository repository.AuthRepository
}

func NewAuthUsecase(repository repository.AuthRepository) AuthUsecase {
	return &authUsecase{repository}
}

func (s *authUsecase) URegister(c echo.Context, pl *dto.RegisterRequest) error {
	password, err := bcrypt.GenerateFromPassword([]byte(pl.Password), 4)
	if err != nil {
		return err
	}
	pl.Password = string(password)

	registerCreate := model.User{
		Name:      pl.Name,
		Email:     pl.Email,
		Password:  pl.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = s.authRepository.RRegister(c, &registerCreate)
	if err != nil {
		pqErr := err.(*pq.Error)
		if pqErr.Code == "23505" {
			err = msg.ErrEmailExists
		}
		return err
	}

	return nil
}

func (s *authUsecase) ULogin(c echo.Context, pl *dto.LoginRequest) (*dto.LoginResponse, error) {
	data := &model.User{
		Email: pl.Email,
	}

	if err := s.authRepository.RLogin(c, data); err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(pl.Password)); err != nil {
		log.Printf("Email or Password Incorrect with err: %s\n", err)
		err = msg.ErrIncorrectEmailOrPassword
		return nil, err
	}

	tokenResponse, err := utill.GenerateJWT(data)
	if err != nil {
		log.Printf("Can't generate JTW with err: %s\n", err)
		return nil, err
	}

	userIdStr := strconv.Itoa(data.Id)
	cookie := &http.Cookie{
		Name:    "userId",
		Value:   userIdStr,
		Expires: time.Now().Add(24 * time.Hour),
		// Secure:  true,
		// Path:    "/",
		// Local
		// SameSite: 2,
		// HttpOnly: true,
	}

	c.SetCookie(cookie)

	resp := &dto.LoginResponse{
		Token: tokenResponse.Token,
	}

	return resp, nil
}

func (s *authUsecase) UCheckUserIfIsExists(c echo.Context, pl *dto.CheckUserIfIsExists) (*model.User, error) {
	dataCheckUser := &model.User{
		Id: int(pl.UserId),
	}

	err := s.authRepository.RCheckUserIfIsExists(c, dataCheckUser)
	if err != nil {
		return dataCheckUser, err
	}

	return dataCheckUser, nil
}

func (s *authUsecase) UGetUserList(c echo.Context, pl *dto.GetUserListRequest) ([]*dto.GetUserListResponse, error) {
	result, err := s.authRepository.RGetUserList(c, &model.User{
		Id: int(pl.UserId),
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *authUsecase) UGetSenderReceiverName(c echo.Context, pl *dto.GetSenderReceiverNameRequest) (*dto.GetSenderReceiverNameResponse, error) {
	var result dto.GetSenderReceiverNameResponse
	senderName, err := s.authRepository.RGetSenderName(c, &model.User{
		Id: int(pl.SenderId),
	})
	if err != nil {
		return nil, err
	}

	receiverName, err := s.authRepository.RGetReceiverName(c, &model.User{
		Id: int(pl.ReceiverId),
	})
	if err != nil {
		return nil, err
	}

	result.SenderNane = senderName
	result.ReceiverName = receiverName

	return &result, nil
}

func (s *authUsecase) UGetSenderName(c echo.Context, pl *dto.GetSenderReceiverNameRequest) (*dto.GetSenderReceiverNameResponse, error) {
	var result dto.GetSenderReceiverNameResponse
	senderName, err := s.authRepository.RGetSenderName(c, &model.User{
		Id: int(pl.SenderId),
	})
	if err != nil {
		return nil, err
	}

	result.SenderNane = senderName

	return &result, nil
}
