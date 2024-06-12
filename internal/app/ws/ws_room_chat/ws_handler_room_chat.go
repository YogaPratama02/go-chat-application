package wsroomchat

import (
	_authUsecase "chat-application-api/internal/app/domain/auth/usecase"
	"chat-application-api/internal/app/domain/room_chat/usecase"
	"chat-application-api/internal/app/dto"
	"chat-application-api/internal/pkg/msg"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

type Handler struct {
	hub             *Hub
	RoomChatUsecase usecase.RoomChatUsecase
	AuthUsecase     _authUsecase.AuthUsecase
}

func NewWSHandler(h *Hub, e *echo.Echo, roomChatUsecase usecase.RoomChatUsecase, authUsecase _authUsecase.AuthUsecase) {
	handler := &Handler{
		hub:             h,
		RoomChatUsecase: roomChatUsecase,
		AuthUsecase:     authUsecase,
	}

	route := e.Group("/ws")
	route.GET("/room-chat/:user_id", handler.HandleConnectionWS)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) HandleConnectionWS(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		log.Println("Error convert param room user id:", err)
		return err
	}

	roomChatId, err := strconv.Atoi(c.QueryParam("room_chat_id"))
	if err != nil {
		log.Println("Error convert room chat id to integer:", err)
		return err
	}

	dataUser, err := h.AuthUsecase.UCheckUserIfIsExists(c, &dto.CheckUserIfIsExists{
		UserId: int64(userId),
	})
	if err != nil {
		return err
	}

	isExists, err := h.RoomChatUsecase.UCheckUserIsExistsInRoomChat(c, &dto.CreateRoomChatMessageRequest{
		UserId:     int64(userId),
		RoomChatId: int64(roomChatId),
	})
	if err != nil {
		return err
	}

	if !isExists {
		err = msg.ErrUserHaveNotJoinThisRoom
		log.Printf("Error user have not join in this room chat with err: %s", err)
		return err
	}

	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		log.Println("Error upgrading to websocket:", err)
		return nil
	}

	client := &Client{
		Hub:        h.hub,
		Conn:       conn,
		Send:       make(chan Message),
		RoomId:     int64(roomChatId),
		SenderId:   int64(userId),
		SenderName: dataUser.Name,
	}
	client.Hub.Register <- client
	client.RoomChatUsecase = h.RoomChatUsecase

	go client.WriteMessage()
	go client.ReadMessage()

	return nil
}
