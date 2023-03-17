package api

import (
	"WeddingBackEnd/database/repository"
	model "WeddingBackEnd/model"
	service "WeddingBackEnd/service"
	serviceAuth "WeddingBackEnd/service/auth"
	"WeddingBackEnd/ultilities"
	ProviderJWT "WeddingBackEnd/ultilities/provider/jwt"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-passwd/validator"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

type MyClaims struct {
	jwt.StandardClaims
	Scopes    string   `json:"scopes"`
	DeviceIDs []string `json:"deviceIDs"`
}

type VerifyGmailCode struct {
	Code string `json:"verifycode"`
}

type AuthSignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type AuthSignUp struct {
	IdentityID  string   `bson:"identityID" json:"identityID"`
	FirstName   string   `bson:"firstName" json:"firstName"`
	LastName    string   `bson:"lastName" json:"lastName"`
	Birthday    int      `bson:"birthday" json:"birthday"`
	Gender      int      `bson:"gender" json:"gender"`
	PhoneNumber string   `bson:"phoneNumber" json:"phoneNumber"`
	Locations   []string `bson:"locations" json:"locations"`
	Email       string   `bson:"email" json:"email" valid:"required,email"`
	Password    string   `bson:"password" json:"password" valid:"required"`
	Scopes      string   `bson:"scopes" json:"scopes"`
	UserID      string   `bson:"userID" json:"userID"`
}
type ResetPassword struct {
	Contact         string
	Password        string
	ConfirmPassword string
}
type ChangePassword struct {
	Contact         string
	Password        string
	NewPassword     string
	ConfirmPassword string
}
type AuthHandler struct {
	JwtSecret          string
	AccountRepository  repository.AccountRepository
	AccessIDRepository repository.UserAccessIDRepository
	UserRepository     repository.UserRepository
}

const verifyKey = "Si3wpzWhAaoP9FKQUIMuguAUKm7eYDtfDoBpQJD5iWRxzD5pV8AL866oUUua01EF"

func (h *AuthHandler) VerifyToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	token, err := ultilities.GetQuery(r, "headers")
	if err != true {
		WriteJSON(w, http.StatusConflict, ResponseBody{
			Message: "cannot get token",
		})
		return
	}
	handler := &service.ServiceHandler{
		UserRepository: h.UserRepository,
	}
	user, errToken := handler.VerifyTokenUser(token, verifyKey)
	if errToken != nil {
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: errToken.Error(),
		})
		return
	}
	WriteJSON(w, 200, user)
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var data = new(AuthSignIn)
	if err := BindJSON(r, &data); err != nil {
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: err.Error(),
		})
		return
	}
	if err := data.ValidSignIn(); err != nil {
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: "Cannot Validate Signin",
			Code:    400,
		})
	}
	account, err := h.AccountRepository.FindByEmail(data.Email)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, ResponseBody{
			Message: "account is not existed",
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(data.Password))
	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, ResponseBody{
			Message: "email or password is incorrect",
			Code:    http.StatusUnauthorized,
		})
		return
	}
	// if mail, err, code := checkAccount(data, h, r) {

	// }
	claims := h.makeClaims(*account, time.Now().Add(time.Hour*12).Unix())
	tokenString, err := ProviderJWT.CreateToken(claims, h.JwtSecret)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ResponseBody{
			Message: "can not create token",
			Code:    http.StatusInternalServerError,
		})
		return
	}
	if err := h.AccessIDRepository.Set(account.Email, claims.Id, fmt.Sprintf("%d", claims.ExpiresAt)); err != nil {
		fmt.Println(err)
		WriteJSON(w, http.StatusInternalServerError, ResponseBody{
			Message: "Unable to create token",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "Successfully login",
		Data: map[string]interface{}{
			"token": tokenString,
		},
	})
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var data = new(serviceAuth.SignUpAccount)
	if err := BindJSON(r, data); err != nil {
		fmt.Println(err)
		WriteJSON(w, http.StatusBadRequest, err)
		return
	}
	if err := data.ValidSignUp(); err != nil {
		fmt.Println(err)
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: err.Error(),
			Code:    400,
		})
	}
	handler := &serviceAuth.SignUpAccountHandler{
		AccountRepository: h.AccountRepository,
		UserRepository:    h.UserRepository,
	}
	err1, err2 := handler.Handle(data)
	if err1 != nil || err2 != nil {
		WriteJSON(w, http.StatusNotFound, ResponseBody{
			Message: "Cannot register",
		})
		return
	}
	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "Sign Up Successfully",
		Data:    data,
	})
}

// find account before reset password
func (h *AuthHandler) FindEmail(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var email string
	if err := govalidator.IsEmail(email); !err {
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: "Email is not existed",
		})
		return
	}

	_, err := h.AccountRepository.FindByEmail(email)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, ResponseBody{
			Message: "account is not existed",
		})
		return
	}
	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "Email exists",
	})
}

func (h *AuthHandler) ForgetPassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var data *ResetPassword
	if data.Contact == "" {
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: "Contact is empty",
		})
		return
	}
	passwordValidator := validator.Regex(ultilities.REGPASS, nil)
	if passwordValidator(data.Password) != nil {
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: "Password is not strong enough",
			Code:    http.StatusBadRequest,
		})
		return
	}
	if data.Password != data.ConfirmPassword {
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: "Authentication password is not correct",
			Code:    http.StatusBadRequest,
		})
		return
	}
	password, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, err)
		return
	}
	data.Password = string(password)
	if govalidator.IsEmail(data.Contact) {
		err := forgetPasswordByEmail(h, *data)
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, ResponseBody{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}
		WriteJSON(w, http.StatusOK, ResponseBody{
			Message: "Successfully",
			Code:    http.StatusOK,
		})
		return
	}
	err = forgetPasswordByPhoneNumber(h, *data)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "Successfully",
		Code:    http.StatusOK,
	})

}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request, _ httprouter.Router) {
	var data ChangePassword
	passwordValidator := validator.Regex(ultilities.REGPASS, nil)
	if passwordValidator(data.Password) != nil {
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: "Password is not strong enough",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if data.Password != data.ConfirmPassword {
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: "Authentication password is not correct",
			Code:    http.StatusBadRequest,
		})
		return
	}
	password, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, err)
		return
	}
	data.Password = string(password)
	findAccount, _ := h.AccountRepository.FindByEmail(data.Contact)
	err = bcrypt.CompareHashAndPassword([]byte(findAccount.Password), []byte(data.Password))
	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, ResponseBody{
			Message: "email or password is incorrect",
			Code:    http.StatusUnauthorized,
		})
		return
	}
	account := model.Account{
		ID:        findAccount.ID,
		Email:     findAccount.Email,
		Password:  data.NewPassword,
		UserID:    findAccount.UserID,
		CreatedAt: findAccount.CreatedAt,
		UpdatedAt: time.Now(),
	}
	err = h.AccountRepository.UpdateByEmail(data.Contact, account)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ResponseBody{
			Message: err.Error(),
		})
		return
	}
	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "Change password succesfully",
	})
}

func forgetPasswordByEmail(h *AuthHandler, data ResetPassword) error {
	findAccount, err := h.AccountRepository.FindByEmail(data.Contact)
	if err != nil {
		return err
	}
	account := model.Account{
		ID:        findAccount.ID,
		Email:     findAccount.Email,
		Password:  data.Password,
		UserID:    findAccount.UserID,
		CreatedAt: findAccount.CreatedAt,
		UpdatedAt: time.Now(),
	}
	err = h.AccountRepository.UpdateByEmail(account.Email, account)
	if err != nil {
		return err
	}
	return nil
}
func forgetPasswordByPhoneNumber(h *AuthHandler, data ResetPassword) error {
	findAccount, err := h.AccountRepository.FindByPhoneNumber(data.Contact)
	if err != nil {
		return err
	}
	account := model.Account{
		ID:        findAccount.ID,
		Email:     findAccount.Email,
		Password:  data.Password,
		UserID:    findAccount.UserID,
		CreatedAt: findAccount.CreatedAt,
		UpdatedAt: time.Now(),
	}
	err = h.AccountRepository.UpdateByPhoneNumber(data.Contact, account)
	if err != nil {
		return err
	}
	return nil
}
func (r *AuthSignIn) TrimSignIn() {
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)
}

func (r *AuthSignIn) ValidSignIn() error {
	r.TrimSignIn()
	if r.Email == "" || r.Password == "" {
		return errors.New("there is still 1 empty box")
	}
	return nil
}

func (h *AuthHandler) makeClaims(account model.Account, tme int64) MyClaims {
	return MyClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        bson.NewObjectId().Hex(),
			Audience:  account.UserID.String(),
			Subject:   account.Email,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: tme,
		},
		Scopes: account.Scopes,
	}
}
