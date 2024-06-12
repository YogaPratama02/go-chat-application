package repository

import (
	"chat-application-api/internal/app/model"
	"database/sql"
	"log"
)

type DirectMessageRepository interface {
	RCheckWhetherSenderEverChatReceiverOrNot(request *model.DirectMessage) (bool, error)
	RCreateDirectMessage(request *model.DirectMessage) (int64, error)
	RGetIdDirectMessageBetweenSenderAndReceiver(request *model.DirectMessage) (int64, error)
	RCreateDirectMessageDetail(request *model.DirectMessageDetail) (int64, error)
}

type directMessageRepository struct {
	db *sql.DB
}

func NewDirectMessageRepository(db *sql.DB) DirectMessageRepository {
	return &directMessageRepository{
		db: db,
	}
}

func (r *directMessageRepository) RCheckWhetherSenderEverChatReceiverOrNot(request *model.DirectMessage) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM direct_messages WHERE sender_id = $1 and receiver_id = $2 OR sender_id = $2 and receiver_id = $1)`
	err := r.db.QueryRow(query, request.SenderId, request.ReceiverId).Scan(&exists)
	if err != nil {
		log.Printf("Error check whether sender is ever chat receiver or not to database with err: %s", err)
		return exists, err
	}

	return exists, nil
}

func (r *directMessageRepository) RCreateDirectMessage(request *model.DirectMessage) (int64, error) {
	sqlStatement := `INSERT INTO direct_messages (sender_id, receiver_id, created_at) VALUES ($1, $2, $3) RETURNING id`
	var lastInsertId int64

	err := r.db.QueryRow(sqlStatement, &request.SenderId, &request.ReceiverId, &request.CreatedAt).Scan(&lastInsertId)
	if err != nil {
		log.Printf("Error create direct message to database with err: %s", err)
		return lastInsertId, err
	}

	return lastInsertId, nil
}

func (r *directMessageRepository) RGetIdDirectMessageBetweenSenderAndReceiver(request *model.DirectMessage) (int64, error) {
	var id int64
	query := `SELECT id FROM direct_messages WHERE sender_id = $1 and receiver_id = $2 OR sender_id = $2 and receiver_id = $1`
	err := r.db.QueryRow(query, request.SenderId, request.ReceiverId).Scan(&id)
	if err != nil {
		log.Printf("Error get id direct message between sender and receiver to database with err: %s", err)
		return id, err
	}

	return id, nil
}

func (r *directMessageRepository) RCreateDirectMessageDetail(request *model.DirectMessageDetail) (int64, error) {
	sqlStatement := `INSERT INTO direct_message_details (direct_message_id, user_id, content, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	var lastInsertId int64

	err := r.db.QueryRow(sqlStatement, &request.DirectMessageId, &request.UserId, &request.Content, &request.CreatedAt).Scan(&lastInsertId)
	if err != nil {
		log.Printf("Error create direct message detail to database with err: %s", err)
		return lastInsertId, err
	}

	return lastInsertId, nil
}
