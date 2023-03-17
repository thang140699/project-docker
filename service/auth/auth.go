package auth

import (
	"WeddingBackEnd/database/repository"
	"WeddingBackEnd/model"
	"WeddingBackEnd/ultilities"
	"errors"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-passwd/validator"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

type SignUpAccount struct {
	IdentityID  string           `bson:"identityID" json:"identityID"`
	FirstName   string           `bson:"firstName" json:"firstName"`
	LastName    string           `bson:"lastName" json:"lastName"`
	Birthday    int              `bson:"birthday" json:"birthday"`
	Gender      int              `bson:"gender" json:"gender"`
	PhoneNumber string           `bson:"phoneNumber" json:"phoneNumber"`
	Locations   []model.Location `bson:"locations" json:"locations"`
	Email       string           `bson:"email" json:"email" valid:"required,email"`
	Password    string           `bson:"password" json:"password" valid:"required"`
	Scopes      string           `bson:"scopes" json:"scopes"`
	UserID      string           `bson:"userID" json:"userID"`
}

func (r *SignUpAccount) TrimSignUp() {
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
}

func (r *SignUpAccount) ValidSignUp() error {
	r.TrimSignUp()
	if _, err := govalidator.ValidateStruct(r); err != nil {
		return err
	}
	if r.Email == "" || r.Password == "" || r.FirstName == "" || r.LastName == "" || len(r.Locations) == 0 || r.PhoneNumber == "" || string(rune(r.Gender)) == "" || r.IdentityID == "" || string(rune(r.Birthday)) == "" {
		return errors.New("there is still 1 empty box")
	}
	if ok := govalidator.IsEmail(r.Email); !ok {
		return errors.New("email invalid")
	}
	passwordValidator := validator.Regex(ultilities.REGPASS, nil)
	if passwordValidator(r.Password) != nil {
		return errors.New("password is not strong enough")
	}
	phoneNumberValidator := validator.Regex(ultilities.RegPhoneNumber, nil)
	if phoneNumberValidator(r.PhoneNumber) != nil {
		return errors.New("incorrect phone number")
	}
	nameValidator := validator.Regex(ultilities.ReggName, nil)
	if nameValidator(r.FirstName) != nil || nameValidator(r.LastName) != nil {
		return errors.New("incorrect Name")
	}
	if len(r.Scopes) == 0 {
		r.Scopes = "user"
	}
	return nil
}

type SignUpAccountHandler struct {
	AccountRepository repository.AccountRepository
	UserRepository    repository.UserRepository
}

func (h *SignUpAccountHandler) Handle(c *SignUpAccount) (error, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		return err, nil
	}
	user := model.User{
		ID:          bson.NewObjectId(),
		IdentityID:  c.IdentityID,
		FirstName:   c.FirstName,
		LastName:    c.LastName,
		Birthday:    c.Birthday,
		Gender:      c.Gender,
		PhoneNumber: c.PhoneNumber,
		Email:       c.Email,
		Locations:   c.Locations,
		CreatedTime: ultilities.TimeInUTC(time.Now()),
		UpdatedTime: ultilities.TimeInUTC(time.Now()),
	}
	account := model.Account{
		ID:        bson.NewObjectId(),
		Email:     c.Email,
		Password:  string(password),
		Scopes:    c.Scopes,
		UserID:    user.ID,
		CreatedAt: ultilities.TimeInUTC(time.Now()),
		UpdatedAt: ultilities.TimeInUTC(time.Now()),
	}

	err = h.UserRepository.Save(user)
	if err != nil {
		return err, nil
	}
	return nil, h.AccountRepository.Save(account)
}
