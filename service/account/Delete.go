package account

import (
	"WeddingBackEnd/database/repository"
)

type DeleteAccountHandler struct {
	AccountRepository repository.AccountRepository
	UserRepository    repository.UserRepository
}

func (h *DeleteAccountHandler) DeleteHandle(id string) error {
	account, err := h.AccountRepository.FindByID(id)
	if err != nil {
		return err
	}

	err = h.UserRepository.RemoveByID(string(account.UserID))
	if err != nil {
		return err
	}
	return h.AccountRepository.RemoveByID(account.ID.Hex())
}
