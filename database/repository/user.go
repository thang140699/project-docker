package repository

import (
	"WeddingBackEnd/model"
	"time"
)

type UserRepository interface {
	All() ([]model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindByID(id string) (*model.User, error)
	FindByIdentifyID(identifyID string) (*model.User, error)
	FindByPhoneNumber(phoneNumber string) (*model.User, error)

	Save(user model.User) error

	UpdateByID(id string, user model.User) error
	UpdateByIdentifyID(identifyID string, user model.User) error
	UpdateByPhoneNumber(phoneNumber string, user model.User) error

	RemoveByID(id string) error
	RemoveByIdentifyID(identifyID string) error
	RemoveByPhoneNumber(phoneNumber string) error
}

type UserAccessIDRepository interface {
	Set(username, id, exp string) error
	Exists(username, id string) (bool, error)
	Get(username, accessID string) (time.Time, error)
	All(username string) (map[string]string, error)
	Delete(username, id string) error
}
