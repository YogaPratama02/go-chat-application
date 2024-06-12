package delivery

import (
	"chat-application-api/internal/app/domain/room_chat/usecase"
	"chat-application-api/internal/app/dto"
	"chat-application-api/internal/pkg/middleware"
	"chat-application-api/internal/pkg/response"
	"chat-application-api/internal/pkg/utill"
	"chat-application-api/internal/pkg/validation"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_authUsecase "chat-application-api/internal/app/domain/auth/usecase"

	"github.com/labstack/echo"
)

type ChatRoomHandler struct {
	RoomChatUsecase usecase.RoomChatUsecase
	AuthUsecase     _authUsecase.AuthUsecase
}

func NewChatRoomHandler(e *echo.Echo, roomChatUsecase usecase.RoomChatUsecase, authUseacse _authUsecase.AuthUsecase) {
	handler := &ChatRoomHandler{
		RoomChatUsecase: roomChatUsecase,
		AuthUsecase:     authUseacse,
	}

	route := e.Group("/api/v1/chat")
	route.Use(middleware.ValidateToken)
	route.POST(`/room-chat`, handler.CreateChatRoom)
	route.POST(`/room-chat/join`, handler.JoinUserRoomChat)
	route.GET(`/room-chat`, handler.ListRoomChat)
	route.DELETE(`/room-chat/:room_chat_id`, handler.LeaveRoomChat)
	route.GET(`/room-chat-message/:room_chat_id`, handler.ListRoomChatMessage)

	// html
	html := e.Group("/page")
	html.GET(`/room-chat/:user_id`, handler.RoomChatPage)
}

func (h *ChatRoomHandler) RoomChatPage(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		log.Println("Error get param user id:", err)

		response.NewHandlerResponse(err, nil).BadRequest(c)
		return err
	}

	roomChatId, err := strconv.Atoi(c.QueryParam("room_chat_id"))
	if err != nil {
		log.Println("Error get param room chat id id:", err)

		response.NewHandlerResponse(err, nil).BadRequest(c)
		return nil
	}

	roomName, err := h.RoomChatUsecase.UGetRoomNameChat(c, &dto.GetRoomChatNameRequest{
		RoomChatId: int64(roomChatId),
	})
	if err != nil {
		response.NewHandlerResponse(err.Error(), nil).Failed(c)
		return err
	}

	userName, err := h.AuthUsecase.UGetSenderName(c, &dto.GetSenderReceiverNameRequest{
		SenderId: int64(userId),
	})
	if err != nil {
		response.NewHandlerResponse(err, nil).BadRequest(c)
		return err
	}

	return c.Render(http.StatusOK, "room_chat.html", map[string]interface{}{
		"roomChatId": roomChatId,
		"roomName":   roomName,
		"userId":     userId,
		"senderName": userName.SenderNane,
	})
}

func (h *ChatRoomHandler) CreateChatRoom(c echo.Context) error {
	var (
		pl          dto.CreateChatRoomRequest
		jwtResponse utill.ResponseTokenJwt = c.Get("token-data").(utill.ResponseTokenJwt)
	)

	if err := c.Bind(&pl); err != nil {
		response.NewHandlerResponse(err.Error(), nil).BadRequest(c)
		return nil
	}

	pl.UserId = int64(jwtResponse.Id)

	if err := validation.DoValidation(&pl); err != nil {
		response.NewHandlerResponse(err, nil).BadRequest(c)
		return nil
	}

	err := h.RoomChatUsecase.UCreateChatRoom(c, &pl)
	if err != nil {
		response.NewHandlerResponse(err.Error(), nil).Failed(c)
		return nil
	}

	response.NewHandlerResponse("Successfully create room chat", nil).SuccessCreate(c)
	return nil
}

func (h *ChatRoomHandler) JoinUserRoomChat(c echo.Context) error {
	var (
		pl dto.AddUserRoomChatRequest
	)

	if err := c.Bind(&pl); err != nil {
		response.NewHandlerResponse(err.Error(), nil).BadRequest(c)
		return nil
	}

	if err := validation.DoValidation(&pl); err != nil {
		response.NewHandlerResponse(err, nil).BadRequest(c)
		return nil
	}

	err := h.RoomChatUsecase.UJoinUserRoomChat(c, &pl)
	if err != nil {
		response.NewHandlerResponse(err.Error(), nil).Failed(c)
		return nil
	}

	response.NewHandlerResponse("Successfully add user in chat room", nil).SuccessCreate(c)
	return nil
}

func (h *ChatRoomHandler) ListRoomChat(c echo.Context) error {
	var (
		pl          dto.ListChatRoomRequest
		jwtResponse utill.ResponseTokenJwt = c.Get("token-data").(utill.ResponseTokenJwt)
	)

	userId, err := c.Cookie("userId")
	if err != nil {
		log.Println("Error get cookie user id:", err)
		return err
	}

	fmt.Println(userId)

	pl.UserId = int64(jwtResponse.Id)

	result, err := h.RoomChatUsecase.UListRoomChat(c, &pl)
	if err != nil {
		response.NewHandlerResponse(err.Error(), nil).Failed(c)
		return nil
	}

	response.NewHandlerResponse("Successfully get list room chat user", result).Success(c)

	return nil
}

func (h *ChatRoomHandler) ListRoomChatMessage(c echo.Context) error {
	var (
		pl          dto.ListRoomChatMessageRequest
		jwtResponse utill.ResponseTokenJwt = c.Get("token-data").(utill.ResponseTokenJwt)
	)

	roomChatId, err := strconv.Atoi(c.Param("room_chat_id"))
	if err != nil {
		response.NewHandlerResponse(err, nil).BadRequest(c)
		return nil
	}

	pl.RoomChatId = int64(roomChatId)
	pl.UserId = int64(jwtResponse.Id)

	result, err := h.RoomChatUsecase.UListRoomChatMessage(c, &pl)
	if err != nil {
		response.NewHandlerResponse(err.Error(), nil).Failed(c)
		return nil
	}

	response.NewHandlerResponse("Successfully get list room chat message", result).Success(c)

	return nil
}

func (h *ChatRoomHandler) LeaveRoomChat(c echo.Context) error {
	var (
		pl          dto.LeaveRoomChatRequest
		jwtResponse utill.ResponseTokenJwt = c.Get("token-data").(utill.ResponseTokenJwt)
	)

	roomChatId, err := strconv.Atoi(c.Param("room_chat_id"))
	if err != nil {
		response.NewHandlerResponse(err, nil).BadRequest(c)
		return nil
	}

	pl.RoomChatId = int64(roomChatId)
	pl.UserId = int64(jwtResponse.Id)

	err = h.RoomChatUsecase.ULeaveRoomChat(c, &pl)
	if err != nil {
		response.NewHandlerResponse(err.Error(), nil).Failed(c)
		return nil
	}

	response.NewHandlerResponse("Successfully leave from this room chat", nil).Success(c)

	return nil
}
