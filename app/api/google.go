package api

import (
	"WeddingBackEnd/database/repository"
	model "WeddingBackEnd/model"
	"WeddingBackEnd/service/auth"
	serviceAuth "WeddingBackEnd/service/auth"
	googleUltis "WeddingBackEnd/ultilities/google"
	ProviderJWT "WeddingBackEnd/ultilities/provider/jwt"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
)

type GoogleHandler struct {
	UserRepository    repository.UserRepository
	AccountRepository repository.AccountRepository
}

type ObjectGoogle struct {
	AccessToken string `bson:"access_token" json:"access_token"`
	Scope       string `bson:"scope" json:"scope"`
}

type DataGoogleResponse struct {
	State string `bson:"state" json:"state"`
	Code  string `bson:"code" json:"code"`
}

func (h *GoogleHandler) HandleGoogleCallBack(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var objGoogle ObjectGoogle

	err := json.NewDecoder(r.Body).Decode(&objGoogle)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: "error",
			Code:    http.StatusBadRequest,
		})
		return
	}

	UserGoogleDetails, err := googleUltis.GetUserInfoFromGoogle(objGoogle.AccessToken)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: err.Error(),
		})
		return
	}

	authToken, err := h.SignInUser(UserGoogleDetails)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: err.Error(),
		})
		return
	}

	cookie := &http.Cookie{Name: "Authorization", Value: "Bearer " + authToken, Path: "/"}
	http.SetCookie(w, cookie)

	WriteJSON(w, http.StatusBadRequest, ResponseBody{
		Message: "success",
		Data: map[string]interface{}{
			"token": authToken,
		},
	})

}

func (h *GoogleHandler) SignInUser(data model.User) (string, error) {

	if reflect.DeepEqual(data, model.User{}) {
		return "", errors.New("user details Can't be empty")
	}

	if data.Email == "" {
		return "", errors.New("lastName can't be empty")
	}

	if data.LastName == "" || data.FirstName == "" {
		return "", errors.New("password can't be empty")
	}

	handler := &serviceAuth.SignUpAccountHandler{
		AccountRepository: h.AccountRepository,
		UserRepository:    h.UserRepository,
	}

	account := &auth.SignUpAccount{
		IdentityID:  data.IdentityID,
		FirstName:   data.FirstName,
		LastName:    data.LastName,
		Birthday:    data.Birthday,
		Gender:      data.Gender,
		PhoneNumber: data.PhoneNumber,
		Locations:   data.Locations,
		Email:       data.Email,
		Password:    "",
		Scopes:      "google",
		UserID:      string(data.ID),
	}

	_, err2 := handler.Handle(account)
	if err2 != nil {
		return "", err2
	}

	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        bson.NewObjectId().Hex(),
			Audience:  string(data.ID),
			Subject:   data.Email,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
		},
		Scopes: "google",
	}
	tokenString, _ := ProviderJWT.CreateToken(claims, "abc")
	if tokenString == "" {
		return "", errors.New("unable generate Auth token")
	}
	return tokenString, nil
}
