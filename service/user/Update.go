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

type UpdateUser struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	IdentityID  string        `bson:"identityID" json:"identityID"`
	FirstName   string        `bson:"firstName" json:"firstName"`
	LastName    string        `bson:"lastName" json:"lastName"`
	Birthday    int           `bson:"birthday" json:"birthday"`
	PhoneNumber string        `bson:"phoneNumber" json:"phoneNumber"`
	Email       string        `bson:"email" json:"email"`
	Gender      int           `bson:"gender" json:"gender"`
	UpdatedTime time.Time     `bson:"updatedTime" json:"updatedTime"`
}
type UpdateUserHandler struct {
	UserRepository repository.UserRepository
}

func (e *UpdateUser) Validate() error {

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
	err2 := govalidator.IsNull(string(e.Gender))
	if err2 {
		return fmt.Errorf("invalid Gender")
	}
	return nil
}
func (handle *UpdateUserHandler) Update(id string, e *UpdateUser) error {
	if err := e.Validate(); err != nil {
		return err
	}
	user := model.User{
		ID:          e.ID,
		IdentityID:  e.IdentityID,
		FirstName:   e.FirstName,
		LastName:    e.LastName,
		PhoneNumber: e.PhoneNumber,
		Birthday:    e.Birthday,
		Email:       e.Email,
		Gender:      e.Gender,
		UpdatedTime: ultilities.TimeInUTC(time.Now()),
	}
	return handle.UserRepository.UpdateByID(id, user)
}
