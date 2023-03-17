package account

import (
	"WeddingBackEnd/database/repository"
	"WeddingBackEnd/ultilities"
	"time"
)

type UpdateAccount struct {
	UserID      string `bson:"userID" json:"userID"`
	IdentityID  string `bson:"identityID" json:"identityID"`
	PhoneNumber string `bson:"phoneNumber" json:"phoneNumber"`
}

func (c *UpdateAccount) validUpdateAccount() error {
	return nil
}

type UpdateAccountHandler struct {
	AccountRepository repository.AccountRepository
	UserRepository    repository.UserRepository
}

func (h *UpdateAccountHandler) UpdateHandle(a *UpdateAccount) error {
	err := a.validUpdateAccount()
	if err != nil {
		return err
	}
	user, err := h.UserRepository.FindByID(a.UserID)
	if err != nil {
		return nil
	}
	if a.IdentityID != "" {
		user.IdentityID = a.IdentityID
	}
	if a.PhoneNumber != "" {
		user.PhoneNumber = a.PhoneNumber
	}
	user.UpdatedTime = ultilities.TimeInUTC(time.Now())
	return h.UserRepository.UpdateByID(a.UserID, *user)
}
