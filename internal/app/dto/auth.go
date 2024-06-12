package dto

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type CheckUserIfIsExists struct {
	UserId int64 `json:"id" validate:"required"`
}

type GetUserListRequest struct {
	UserId int64 `json:"user_id"`
}

type GetUserListResponse struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type GetSenderReceiverNameRequest struct {
	SenderId   int64 `json:"sender_d"`
	ReceiverId int64 `json:"receiver_id"`
}

type GetSenderReceiverNameResponse struct {
	SenderNane   string `json:"sender_name"`
	ReceiverName string `json:"receiver_name"`
}
