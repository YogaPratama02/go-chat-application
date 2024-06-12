package repository

import (
	"chat-application-api/internal/app/dto"
	"chat-application-api/internal/app/model"
	"database/sql"
	"log"

	"github.com/labstack/echo"
)

type RoomChatRepository interface {
	RCreateRoomChat(c echo.Context, request *model.ChatRoom) (int64, error)
	RCreateRoomChatUser(c echo.Context, request *model.RoomChatUser) error
	RCheckUserIsExistsInRoomChat(request *model.RoomChatUser) (bool, error)
	RListRoomChat(c echo.Context, request *model.RoomChatUser) ([]*dto.ListChatRoomResponse, error)
	RCreateRoomChatMessage(request *model.RoomChatMessage) (int64, error)
	RListRoomChatMessage(request *model.RoomChatMessage) ([]*dto.ListRoomChatMessageResponse, error)
	RLeaveRoomChat(c echo.Context, request *model.RoomChatMessage) error
	RGetRoomNameChat(c echo.Context, request *model.RoomChatMessage) (string, error)
}

type roomChatRepository struct {
	db *sql.DB
}

func NewChatRoomRepository(db *sql.DB) RoomChatRepository {
	return &roomChatRepository{
		db: db,
	}
}

func (r *roomChatRepository) RCreateRoomChat(c echo.Context, request *model.ChatRoom) (int64, error) {
	sqlStatement := `INSERT INTO room_chats (room_name, user_id, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var lastInserId int64

	err := r.db.QueryRowContext(c.Request().Context(), sqlStatement, &request.RoomName, &request.UserId, &request.Status, &request.CreatedAt, &request.UpdatedAt).Scan(&lastInserId)
	if err != nil {
		log.Printf("Error create room chat to database with err: %s", err)
		return lastInserId, err
	}

	return lastInserId, nil
}

func (r *roomChatRepository) RCreateRoomChatUser(c echo.Context, request *model.RoomChatUser) error {
	sqlStatement := `INSERT INTO room_chat_user (user_id, room_chat_id, created_at) VALUES ($1, $2, $3)`

	_, err := r.db.ExecContext(c.Request().Context(), sqlStatement, &request.UserId, &request.RoomChatId, &request.CreatedAt)
	if err != nil {
		log.Printf("Error create room chat user to database with err: %s", err)
		return err
	}

	return nil
}

func (r *roomChatRepository) RCheckUserIsExistsInRoomChat(request *model.RoomChatUser) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM room_chat_user WHERE user_id = $1 and room_chat_id = $2)`
	err := r.db.QueryRow(query, request.UserId, request.RoomChatId).Scan(&exists)
	if err != nil {
		log.Printf("Error check user in room chat user table to database with err: %s", err)
		return exists, err
	}

	return exists, nil
}

func (r *roomChatRepository) RListRoomChat(c echo.Context, request *model.RoomChatUser) ([]*dto.ListChatRoomResponse, error) {
	var result []*dto.ListChatRoomResponse
	query := `select b.id, b.room_name FROM room_chat_user a 
		LEFT JOIN room_chats b ON a.room_chat_id = b.id
		WHERE a.user_id = $1`

	rows, err := r.db.QueryContext(c.Request().Context(), query, request.UserId)
	if err != nil {
		log.Printf("Error get room chat to database with err: %s", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		data := &dto.ListChatRoomResponse{}
		err = rows.Scan(&data.RoomChatId, &data.RoomName)
		if err != nil {
			log.Printf("Error mapping list room chat to database with err: %s", err)
			return nil, err
		}

		result = append(result, data)
	}

	return result, nil
}

func (r *roomChatRepository) RCreateRoomChatMessage(request *model.RoomChatMessage) (int64, error) {
	sqlStatement := `INSERT INTO room_chat_messages (user_id, room_chat_id, message, created_at) VALUES ($1, $2, $3, $4) returning id`
	var lastInsertId int64

	err := r.db.QueryRow(sqlStatement, &request.UserId, &request.RoomChatId, &request.Message, &request.CreatedAt).Scan(&lastInsertId)
	if err != nil {
		log.Printf("Error create room chat message to database with err: %s", err)
		return lastInsertId, err
	}

	return lastInsertId, nil
}

func (r *roomChatRepository) RListRoomChatMessage(request *model.RoomChatMessage) ([]*dto.ListRoomChatMessageResponse, error) {
	var result []*dto.ListRoomChatMessageResponse
	query := `select a.user_id, a.id, a.room_chat_id, b.name, a.message, a.created_at FROM room_chat_messages a
		LEFT JOIN users b ON a.user_id = b.id
		WHERE a.room_chat_id = $1
		ORDER BY a.created_at ASC`

	rows, err := r.db.Query(query, request.RoomChatId)
	if err != nil {
		log.Printf("Error get room chat message to database with err: %s", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		data := &dto.ListRoomChatMessageResponse{}
		err = rows.Scan(&data.UserId, &data.RoomChatMessageId, &data.RoomChatId, &data.UserName, &data.Message, &data.CreatedAt)
		if err != nil {
			log.Printf("Error mapping list room chat message to database with err: %s", err)
			return nil, err
		}

		if data.UserId == request.UserId {
			data.IsUser = true
		}

		result = append(result, data)
	}

	return result, nil
}

func (r *roomChatRepository) RLeaveRoomChat(c echo.Context, request *model.RoomChatMessage) error {
	sqlStatement := `DELETE FROM room_chat_user WHERE user_id = $1 AND room_chat_id = $2`

	_, err := r.db.ExecContext(c.Request().Context(), sqlStatement, &request.UserId, &request.RoomChatId)
	if err != nil {
		log.Printf("Error delete room chat user to database with err: %s", err)
		return err
	}

	return nil
}

func (r *roomChatRepository) RGetRoomNameChat(c echo.Context, request *model.RoomChatMessage) (string, error) {
	var roomName string
	sqlStatement := `SELECT room_name FROM room_chats WHERE id = $1`

	err := r.db.QueryRowContext(c.Request().Context(), sqlStatement, &request.RoomChatId).Scan(&roomName)
	if err != nil {
		log.Printf("Error get room chat name to database with err: %s", err)
		return roomName, err
	}

	return roomName, nil
}
