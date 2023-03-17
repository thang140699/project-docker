package api

import (
	"WeddingBackEnd/database/repository"
	"WeddingBackEnd/model"
	"WeddingBackEnd/service/auth"
	facbookUltilities "WeddingBackEnd/ultilities/facebook"
	ProviderJWT "WeddingBackEnd/ultilities/provider/jwt"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
)

type FacebookHandler struct {
	UserRepository    repository.UserRepository
	AccountRepository repository.AccountRepository
}
type DataFaceBookResponse struct {
	State string `bson:"state" json:"state"`
	Code  string `bson:"code" json:"code"`
}

func (h *FacebookHandler) FacebookLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	body, _ := io.ReadAll(r.Body) // check for errors
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	accessToken := keyVal["accessToken"]
	UserDetails, err := facbookUltilities.GetUserInfoFromFacebook(accessToken)
	if err != nil {
		fmt.Println(err)
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: "Cannot get user from facebook",
		})
		return
	}
	authToken, err := h.SignInUser(UserDetails)
	if err != nil {
		fmt.Println(err)
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: "Cannot signin from facebook",
		})
		return
	}

	cookie := &http.Cookie{Name: "Authorization", Value: "Bearer " + authToken, Path: "/"}
	http.SetCookie(w, cookie)
	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "Login successfully",
		Code:    http.StatusOK,
		Data: map[string]interface{}{
			"token": authToken,
		},
	})
}

func (h *FacebookHandler) SignInUser(data model.User) (string, error) {

	if reflect.DeepEqual(data, model.User{}) {
		return "", errors.New("user details Can't be empty")
	}

	if data.Email == "" {
		return "", errors.New("last Name can't be empty")
	}

	if data.LastName == "" || data.FirstName == "" {
		return "", errors.New("password can't be empty")
	}
	handler := &auth.SignUpAccountHandler{
		AccountRepository: h.AccountRepository,
		UserRepository:    h.UserRepository,
	}
	acc := &auth.SignUpAccount{
		IdentityID:  data.IdentityID,
		FirstName:   data.FirstName,
		LastName:    data.LastName,
		Birthday:    data.Birthday,
		Gender:      data.Birthday,
		PhoneNumber: data.PhoneNumber,
		Locations:   data.Locations,
		Email:       data.Email,
		Password:    "",
		Scopes:      "facebook",
		UserID:      string(data.ID),
	}
	_, err := handler.UserRepository.FindByEmail(data.Email)
	if err != nil {
		err1, err2 := handler.Handle(acc)
		if err1 != nil || err2 != nil {
			return "", errors.New("error occurred registation" + err1.Error() + err2.Error())
		}
	}
	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        bson.NewObjectId().Hex(),
			Audience:  data.ID.String(),
			Subject:   data.Email,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
		},
		Scopes: "facebook",
	}
	tokenString, _ := ProviderJWT.CreateToken(claims, "JWTSecret")
	if tokenString == "" {
		return "", errors.New("unable generate Auth token")
	}
	return tokenString, nil
}
func GetUserDetails(response http.ResponseWriter, request *http.Request) {

}
