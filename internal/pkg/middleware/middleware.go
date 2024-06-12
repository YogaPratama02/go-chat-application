package middleware

import (
	"chat-application-api/internal/pkg/msg"
	"chat-application-api/internal/pkg/response"
	"chat-application-api/internal/pkg/utill"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func MiddlewareCors(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:          middleware.DefaultSkipper,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))
}

func LoggerConfig(e *echo.Echo) {
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
}

func ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var (
			request  = ctx.Request()
			reqToken string
		)
		headerDataToken := request.Header.Get("Authorization")
		if !strings.Contains(headerDataToken, "Bearer") {
			response.NewHandlerResponse(msg.MsgUnauthorized, nil).Failed(ctx)
			return nil
		}

		splitToken := strings.Split(headerDataToken, "Bearer ")
		if len(splitToken) > 1 {
			reqToken = splitToken[1]
		}

		jwtResponse, err := utill.ClaimsTokenJwt(reqToken)
		if err != nil {
			response.NewHandlerResponse(msg.MsgUnauthorized, nil).Failed(ctx)
			return nil
		}

		// Set data jwt response to ...
		ctx.Set("token-data", jwtResponse)

		return next(ctx)
	}
}
