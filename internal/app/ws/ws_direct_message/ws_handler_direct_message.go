package wsdirectmessage

import (
	_authUsecase "chat-application-api/internal/app/domain/auth/usecase"
	"chat-application-api/internal/app/domain/direct_message/usecase"
	"chat-application-api/internal/app/dto"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

type Handler struct {
	hub                  *Hub
	DirectMessageUsecase usecase.DirectMessageUsecase
	AuthUsecase          _authUsecase.AuthUsecase
}

func NewWSDMHandler(h *Hub, e *echo.Echo, directMessageUsecase usecase.DirectMessageUsecase, authUsecase _authUsecase.AuthUsecase) {
	handler := &Handler{
		hub:                  h,
		DirectMessageUsecase: directMessageUsecase,
		AuthUsecase:          authUsecase,
	}

	route := e.Group("/ws")
	route.GET("/direct-message/:user_id", handler.HandleConnectionWSDM)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) HandleConnectionWSDM(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		log.Println("Error convert param room user id:", err)
		return err
	}

	receiverId, err := strconv.Atoi(c.QueryParam("receiver_id"))
	if err != nil {
		log.Println("Error convert param room receiver id:", err)
		return err
	}

	dataUser, err := h.AuthUsecase.UCheckUserIfIsExists(c, &dto.CheckUserIfIsExists{
		UserId: int64(userId),
	})
	if err != nil {
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
		SenderId:   int64(userId),
		SenderName: dataUser.Name,
		ReceiverId: int64(receiverId),
	}
	client.Hub.Register <- client

	client.DirectMessageUseCase = h.DirectMessageUsecase
	go client.WriteMessage()
	go client.ReadMessage()

	return nil
}
