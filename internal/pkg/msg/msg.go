package msg

import "errors"

const (
	MsgHeaderTokenNotFound = "Header `token` not found"
	MsgUnauthorized        = "You must login before"

	MsgHeaderTokenUnauthorized = "Unauthorized token"
)

var (
	ErrMissingHeaderData = errors.New("missing header data")
	ErrInvalidToken      = errors.New("invalid token")

	ErrEmailExists                     = errors.New("email is already exists")
	ErrIncorrectEmailOrPassword        = errors.New("incorrect Email or Password")
	ErrUserAlreadyExistsInThisRoomChat = errors.New("user already exists in this room chat")
	ErrUserHaveNotJoinThisRoom         = errors.New("user have not been join in this room chat")
	ErrUserDoesNotExists               = errors.New("user does not exist")
)
