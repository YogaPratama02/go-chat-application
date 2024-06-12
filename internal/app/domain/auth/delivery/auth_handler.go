package delivery

import (
	"chat-application-api/internal/app/domain/auth/usecase"
	"chat-application-api/internal/app/dto"
	"chat-application-api/internal/pkg/response"
	"chat-application-api/internal/pkg/validation"

	"github.com/labstack/echo"

	_ "github.com/lib/pq"
)

type AuthHandler struct {
	AuthUsecase usecase.AuthUsecase
}

func NewAuthHandler(e *echo.Echo, authUsecase usecase.AuthUsecase) {
	handler := &AuthHandler{
		AuthUsecase: authUsecase,
	}

	route := e.Group("/api/v1")
	route.POST(`/register`, handler.RegisterHandler)
	route.POST(`/login`, handler.LoginHandler)
}

func (h *AuthHandler) RegisterHandler(c echo.Context) error {
	var pl dto.RegisterRequest

	if err := c.Bind(&pl); err != nil {
		response.NewHandlerResponse(err.Error(), nil).BadRequest(c)
		return nil
	}

	if err := validation.DoValidation(&pl); err != nil {
		response.NewHandlerResponse(err, nil).BadRequest(c)
		return nil
	}

	err := h.AuthUsecase.URegister(c, &pl)
	if err != nil {
		response.NewHandlerResponse(err.Error(), nil).Failed(c)
		return nil
	}

	response.NewHandlerResponse("Successfully register", nil).SuccessCreate(c)
	return nil
}

func (h *AuthHandler) LoginHandler(c echo.Context) error {
	var pl dto.LoginRequest

	if err := c.Bind(&pl); err != nil {
		response.NewHandlerResponse(err.Error(), nil).BadRequest(c)
		return nil
	}

	if err := validation.DoValidation(&pl); err != nil {
		response.NewHandlerResponse(err, nil).BadRequest(c)
		return nil
	}

	token, err := h.AuthUsecase.ULogin(c, &pl)
	if err != nil {
		response.NewHandlerResponse(err.Error(), nil).Failed(c)
		return nil
	}

	response.NewHandlerResponse("Succesfully login", token).Success(c)
	return nil
}
