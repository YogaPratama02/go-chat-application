package dto

type CreateDirectMessageRequest struct {
	SenderId   int64  `json:"sender_id"`
	Content    string `json:"content"`
	ReceiverId int64  `json:"receiver_id"`
}
