package repository

import (
	"chat-application-api/internal/app/dto"
	"chat-application-api/internal/app/model"
	"database/sql"
	"log"

	"github.com/labstack/echo"
)

type AuthRepository interface {
	RRegister(c echo.Context, request *model.User) error
	RLogin(c echo.Context, request *model.User) error
	RCheckUserIfIsExists(c echo.Context, request *model.User) error
	RGetUserList(c echo.Context, request *model.User) ([]*dto.GetUserListResponse, error)
	RGetSenderName(c echo.Context, request *model.User) (string, error)
	RGetReceiverName(c echo.Context, request *model.User) (string, error)
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) RRegister(c echo.Context, request *model.User) error {
	sqlStatement := `INSERT INTO users (name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	_, err := r.db.ExecContext(c.Request().Context(), sqlStatement, &request.Name, &request.Email, &request.Password, &request.CreatedAt, &request.UpdatedAt)
	if err != nil {
		log.Printf("Error register to database with err: %s", err)
		return err
	}

	return nil
}

func (r *authRepository) RLogin(c echo.Context, request *model.User) error {
	sqlStatement := `SELECT id, name, email, password FROM users WHERE email = $1`

	err := r.db.QueryRow(sqlStatement, request.Email).Scan(&request.Id, &request.Name, &request.Email, &request.Password)
	if err != nil {
		log.Printf("Error login to database with err: %s", err)
		return err
	}

	return nil
}

func (r *authRepository) RCheckUserIfIsExists(c echo.Context, request *model.User) error {
	sqlStatement := `SELECT id, name FROM users WHERE id = $1`

	err := r.db.QueryRow(sqlStatement, request.Id).Scan(&request.Id, &request.Name)
	if err != nil {
		log.Printf("Error check user if is exists or not to database with err: %s", err)
		return err
	}

	return nil
}

func (r *authRepository) RGetUserList(c echo.Context, request *model.User) ([]*dto.GetUserListResponse, error) {
	var resp []*dto.GetUserListResponse
	sqlStatement := `SELECT id, name FROM users WHERE id != $1`

	rows, err := r.db.QueryContext(c.Request().Context(), sqlStatement, request.Id)
	if err != nil {
		log.Printf("Error get user list user to database with err: %s", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		data := &dto.GetUserListResponse{}
		err = rows.Scan(&data.Id, &data.Name)
		if err != nil {
			log.Printf("Error mapping get user list with err: %s", err)
			return nil, err
		}

		resp = append(resp, data)
	}

	return resp, nil
}

func (r *authRepository) RGetSenderName(c echo.Context, request *model.User) (string, error) {
	var senderName string

	sqlStatement := `SELECT name FROM users WHERE id = $1`

	err := r.db.QueryRowContext(c.Request().Context(), sqlStatement, request.Id).Scan(&senderName)
	if err != nil {
		log.Printf("Error get sender name to database with err: %s", err)
		return senderName, err
	}

	return senderName, nil
}

func (r *authRepository) RGetReceiverName(c echo.Context, request *model.User) (string, error) {
	var receiverName string

	sqlStatement := `SELECT name FROM users WHERE id = $1`

	err := r.db.QueryRowContext(c.Request().Context(), sqlStatement, request.Id).Scan(&receiverName)
	if err != nil {
		log.Printf("Error get receiver name to database with err: %s", err)
		return receiverName, err
	}

	return receiverName, nil
}
