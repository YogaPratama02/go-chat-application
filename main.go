package main

import (
	wsdirectmessage "chat-application-api/internal/app/ws/ws_direct_message"
	wsroomchat "chat-application-api/internal/app/ws/ws_room_chat"
	"chat-application-api/internal/pkg/database"
	"chat-application-api/router"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error read env file with err: %s", err)
	}

	db := database.ConnectDB()
	defer db.Close()

	router := router.NewRoutes(db)
	hub := wsroomchat.NewHub()
	hubDM := wsdirectmessage.NewHub()
	go hub.Run()
	go hubDM.Run()
	router.LoadHandlers(hub, hubDM)
	router.Echo.Logger.Fatal(router.Echo.Start(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))))
}
