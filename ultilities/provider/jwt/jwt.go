package jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTService struct {
	JwtKey string
}

var expirationTime = time.Now().Add(time.Minute * 10).Unix()

type Claims struct {
	Data interface{} `json:"data"`
	jwt.StandardClaims
}

func (JS JWTService) IssueJWT(data interface{}) (string, error) {
	claims := Claims{
		Data: data,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expirationTime,
		},
	}
	// token, := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return CreateToken(claims, JS.JwtKey)
}

func (JS JWTService) ValidateJWT(tokenString string, data interface{}) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(JS.JwtKey), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return errors.New("invalid signature")
		}
	}
	if !token.Valid {
		return errors.New("token not valid")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return JS.parseDataFromMap(claims, data)
	}

	return fmt.Errorf("invalid JWTService")
}

func (JS JWTService) parseDataFromMap(m jwt.MapClaims, out interface{}) error {
	data, ok := m["data"]
	if !ok || data == nil {
		return fmt.Errorf("invalid JWTService Claims: Data")
	}

	buffers, e := json.Marshal(data)
	if e != nil {
		return e
	}

	return json.Unmarshal(buffers, out)
}
func CreateToken(claims jwt.Claims, JwtSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(JwtSecret))
}
func NewJWT(Key string) *JWTService {
	service := new(JWTService)
	service.JwtKey = Key
	return service
}
