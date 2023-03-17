package user

import "WeddingBackEnd/database/repository"

type DeleteUserHandler struct {
	UserRepository repository.UserRepository
}

func (h *DeleteUserHandler) Handle(Id string) error {
	_, err := h.UserRepository.FindByID(Id)
	if err != nil {
		return err
	}
	err = h.UserRepository.RemoveByID(Id)
	if err != nil {
		return err
	}
	return nil

}
