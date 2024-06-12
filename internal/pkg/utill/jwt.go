package utill

import (
	"chat-application-api/internal/app/model"
	"chat-application-api/internal/pkg/msg"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

var SecretKey = []byte(os.Getenv("SECRET_KEY"))

type JWTClaim struct {
	jwt.StandardClaims
}

type ResponseTokenJwt struct {
	Id           int       `json:"id"`
	Token        string    `json:"token"`
	TokenExpired time.Time `json:"token_expired"`
}

func GenerateJWT(user *model.User) (response *ResponseTokenJwt, err error) {
	tokenJwt := jwt.New(jwt.SigningMethodHS256)
	expiredToken := time.Now().Add(24 * time.Hour)

	claims := tokenJwt.Claims.(jwt.MapClaims)
	claims["id"] = user.Id
	claims["exp"] = expiredToken.Unix()

	token, err := tokenJwt.SignedString(SecretKey)
	if err != nil {
		err = errors.Wrap(err, "failed generate jwt token")
		return
	}

	response = &ResponseTokenJwt{
		Id:           user.Id,
		Token:        token,
		TokenExpired: expiredToken,
	}

	return
}

func ClaimsTokenJwt(token string) (response ResponseTokenJwt, err error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, errors.Wrapf(err, "Unexpected signing method: %v", token.Header["alg"])
		}
		return SecretKey, nil
	})
	if err != nil {
		return
	}

	if jwtToken == nil {
		err = errors.WithStack(msg.ErrInvalidToken)
		return
	}

	if !jwtToken.Valid {
		err = errors.WithStack(msg.ErrInvalidToken)
		return
	}
	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		response = ResponseTokenJwt{
			Id: int(claims["id"].(float64)),
		}
	}

	return
}
