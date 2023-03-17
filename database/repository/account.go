package repository

import "WeddingBackEnd/model"

type AccountRepository interface {
	All() ([]model.Account, error)
	FindByID(Id string) (*model.Account, error)
	FindByEmail(email string) (*model.Account, error)
	FindByPhoneNumber(phoneNumber string) (*model.Account, error)
	Save(account model.Account) error

	UpdateByEmail(email string, account model.Account) error
	UpdateByPhoneNumber(phoneNumber string, account model.Account) error
	RemoveByID(id string) error
	RemoveByEmail(email string) error
	RemoveByUserID(userID string) error
}
