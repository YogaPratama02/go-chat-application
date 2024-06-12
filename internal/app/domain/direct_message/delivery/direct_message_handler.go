package delivery

import (
	_authUsecase "chat-application-api/internal/app/domain/auth/usecase"
	"chat-application-api/internal/app/domain/direct_message/usecase"
	"chat-application-api/internal/app/dto"
	"chat-application-api/internal/pkg/middleware"
	"chat-application-api/internal/pkg/response"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type DirectMessageHandler struct {
	DirectMessageUsecase usecase.DirectMessageUsecase
	AuthUsecase          _authUsecase.AuthUsecase
}

func NewDirectMessageHandler(e *echo.Echo, directMessageUsecase usecase.DirectMessageUsecase, authUseacse _authUsecase.AuthUsecase) {
	handler := &DirectMessageHandler{
		DirectMessageUsecase: directMessageUsecase,
		AuthUsecase:          authUseacse,
	}

	route := e.Group("/api/v1/direct-message")
	route.Use(middleware.ValidateToken)

	// html
	html := e.Group("/page")
	html.GET(`/direct-message/:user_id`, handler.DirectMessagePage)
}

func (h *DirectMessageHandler) DirectMessagePage(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		log.Println("Error get param room chat id:", err)

		response.NewHandlerResponse(err, nil).BadRequest(c)
		return err
	}

	receiverId, err := strconv.Atoi(c.QueryParam("receiver_id"))
	if err != nil {
		log.Println("Error get param receiver id id:", err)

		response.NewHandlerResponse(err, nil).BadRequest(c)
		return err
	}

	userName, err := h.AuthUsecase.UGetSenderReceiverName(c, &dto.GetSenderReceiverNameRequest{
		SenderId:   int64(userId),
		ReceiverId: int64(receiverId),
	})
	if err != nil {
		response.NewHandlerResponse(err, nil).BadRequest(c)
		return err
	}

	return c.Render(http.StatusOK, "direct_message.html", map[string]interface{}{
		"userId":       userId,
		"senderName":   userName.SenderNane,
		"receiverName": userName.ReceiverName,
		"receiverId":   receiverId,
	})
}
