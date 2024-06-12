package router

import (
	_authDelivery "chat-application-api/internal/app/domain/auth/delivery"
	_directMessageDelivery "chat-application-api/internal/app/domain/direct_message/delivery"
	_roomChatDelivery "chat-application-api/internal/app/domain/room_chat/delivery"
	wsdirectmessage "chat-application-api/internal/app/ws/ws_direct_message"
	wsroomchat "chat-application-api/internal/app/ws/ws_room_chat"
	"html/template"
	"io"

	// _wsHandler "chat-application-api/internal/app/ws"
	"chat-application-api/internal/pkg/middleware"
	"database/sql"

	"github.com/labstack/echo"
)

type Router struct {
	Echo         *echo.Echo
	Repositories repositories
	Usecases     usecases
}

type TemplateRegistry struct {
	templates *template.Template
}

// Implement e.Renderer interface
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewRoutes(db *sql.DB) Router {
	e := echo.New()
	e.Debug = false

	// CORS
	middleware.MiddlewareCors(e)

	// Logger
	middleware.LoggerConfig(e)

	repos := newRepositories(db)

	e.Renderer = &TemplateRegistry{
		templates: template.Must(template.ParseGlob("public/*.html")),
	}

	return Router{
		Echo:         e,
		Repositories: repos,
		Usecases:     newUsecases(repos),
	}
}

func (r *Router) LoadHandlers(hub *wsroomchat.Hub, hubDM *wsdirectmessage.Hub) {
	_authDelivery.NewAuthHandler(r.Echo, r.Usecases.AuthUsecase)
	_roomChatDelivery.NewChatRoomHandler(r.Echo, r.Usecases.RoomChatUsecase, r.Usecases.AuthUsecase)
	_directMessageDelivery.NewDirectMessageHandler(r.Echo, r.Usecases.DirectMessageUsecase, r.Usecases.AuthUsecase)
	wsroomchat.NewWSHandler(hub, r.Echo, r.Usecases.RoomChatUsecase, r.Usecases.AuthUsecase)
	wsdirectmessage.NewWSDMHandler(hubDM, r.Echo, r.Usecases.DirectMessageUsecase, r.Usecases.AuthUsecase)
}
