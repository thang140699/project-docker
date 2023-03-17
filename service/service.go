package service

import (
	"WeddingBackEnd/database/repository"
	"WeddingBackEnd/model"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

type ServiceHandler struct {
	UserRepository repository.UserRepository
}

func (handle *ServiceHandler) VerifyTokenUser(tokenStr string, verifyKey string) (model.User, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(verifyKey), nil
	})
	if token == nil {
		return model.User{}, err
	}
	if len(claims["Email"].(string)) == 0 {
		return model.User{}, errors.New("undefined Email")
	}
	email := claims["Email"]
	User, err := handle.UserRepository.FindByEmail(email.(string))
	if err != nil {
		return model.User{}, errors.New("cannot find User by email")
	}
	return *User, nil
}
