package account

import (
	"WeddingBackEnd/database/repository"
	"WeddingBackEnd/model"
	"WeddingBackEnd/ultilities"
	"fmt"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

type AddAccount struct {
	Scopes      string        `bson:"scopes" json:"scopes"`
	Email       string        `bson:"email" json:"email" valid:"required,email"`
	IdentityID  string        `bson:"identityID" json:"identityID" valid:"required"`
	PhoneNumber string        `bson:"phoneNumber" json:"phoneNumber"`
	Password    string        `bson:"password" json:"password" valid:"required"`
	UserID      bson.ObjectId `bson:"userID" json:"userID"`
	CreatedAt   time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time     `bson:"updatedAt" json:"updatedAt"`
	CreatedBy   string        `bson:"createdBy"json:"createdBy"`
}

func (c *AddAccount) ValidAccount() error {
	c.Email = strings.TrimSpace(c.Email)
	c.Password = strings.TrimSpace(c.Password)
	if len(c.UserID) == 0 {
		return fmt.Errorf("invalid UserID")
	}
	_, err := govalidator.ValidateStruct(c)
	return err
}

type AddAccountHandler struct {
	AccountRepository repository.AccountRepository
	UserRepository    repository.UserRepository
}

func (h *AddAccountHandler) AccountHandle(c *AddAccount) error {
	if err := c.ValidAccount(); err != nil {
		return err
	}
	u := model.User{
		ID:          bson.NewObjectId(),
		IdentityID:  c.IdentityID,
		PhoneNumber: c.PhoneNumber,
		CreatedTime: ultilities.TimeInUTC(time.Now()),
		UpdatedTime: ultilities.TimeInUTC(time.Now()),
	}
	err := h.UserRepository.Save(u)
	if err != nil {
		return err
	}
	password, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		h.UserRepository.RemoveByID(u.ID.Hex())
		return err
	}

	account := model.Account{
		ID:        bson.NewObjectId(),
		Email:     c.Email,
		Password:  string(password),
		Scopes:    c.Scopes,
		CreatedAt: ultilities.TimeInUTC(time.Now()),
		UpdatedAt: ultilities.TimeInUTC(time.Now()),
		UserID:    bson.ObjectId(u.ID.Hex()),
		CreatedBy: c.CreatedBy,
	}
	err = h.AccountRepository.Save(account)
	if err != nil {
		h.UserRepository.RemoveByID(u.ID.Hex())
		return err
	}
	return nil
}
