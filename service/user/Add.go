package user

import (
	"WeddingBackEnd/database/repository"
	"WeddingBackEnd/model"
	"WeddingBackEnd/ultilities"
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"gopkg.in/mgo.v2/bson"
)

type AddUser struct {
	IdentityID  string           `bson:"identityID" json:"identityID"`
	FirstName   string           `bson:"firstName" json:"firstName"`
	LastName    string           `bson:"lastName" json:"lastName"`
	Birthday    int              `bson:"birthday" json:"birthday"`
	PhoneNumber string           `bson:"phoneNumber" json:"phoneNumber"`
	Email       string           `bson:"email" json:"email"`
	Gender      int              `bson:"gender" json:"gender"`
	Locations   []model.Location `bson:"locations" json:"locations"`
	CreatedTime time.Time        `bson:"createdTime" json:"createdTime"`
}

func (e *AddUser) Valid() error {

	_, err := govalidator.ValidateStruct(e)
	if err != nil {
		return err
	}
	if len(e.IdentityID) == 0 {
		return fmt.Errorf("invalid IdentityID")
	}
	if len(e.FirstName) == 0 {
		return fmt.Errorf("invalid FirstName")
	}
	if len(e.LastName) == 0 {
		return fmt.Errorf("invalid LastName")
	}
	if len(e.PhoneNumber) == 0 {
		return fmt.Errorf("invalid PhoneNumber")
	}
	if len(e.Email) == 0 {
		return fmt.Errorf("invalid email")
	}
	if string(rune(e.Gender)) == "" {
		return fmt.Errorf("invalid Gender")
	}
	return nil
}

type AddUserHandler struct {
	UserRepository repository.UserRepository
}

func (h *AddUserHandler) ADD(c *AddUser) (string, error) {
	if err := c.Valid(); err != nil {
		return "", err
	}
	au := model.User{
		ID:          bson.NewObjectId(),
		IdentityID:  c.IdentityID,
		FirstName:   c.FirstName,
		LastName:    c.LastName,
		Birthday:    c.Birthday,
		PhoneNumber: c.PhoneNumber,
		Email:       c.Email,
		Gender:      c.Gender,
		Locations:   []model.Location{},
		CreatedTime: ultilities.TimeInUTC(time.Now()),
	}
	return string(au.ID), h.UserRepository.Save(au)
}
