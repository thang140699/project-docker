package api

import (
	"WeddingBackEnd/database/repository"
	"WeddingBackEnd/model"
	serviceUser "WeddingBackEnd/service/user"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type UserHandler struct {
	UserRepository repository.UserRepository
}

func getByPage(start string, end string, limit string, page string, arr []model.User) ([]model.User, error) {
	st, errST := strconv.Atoi(start)
	ed, errED := strconv.Atoi(end)
	lm, errLM := strconv.Atoi(limit)
	pg, errPG := strconv.Atoi(page)
	if errST != nil {
		st = 0
	}
	if errED != nil {
		ed = len(arr)
	}
	if errLM != nil {
		lm = len(arr)
	}
	if errPG != nil {
		pg = 1
	}
	divideResult := len(arr) / lm
	surplus := len(arr) % lm
	if surplus == 0 && pg > divideResult || surplus != 0 && pg > divideResult+1 {
		return nil, errors.New("don't have record")
	}
	if surplus != 0 && pg == divideResult+1 {
		return arr[lm*(pg-1) : ed], nil
	}
	return arr[st:ed], nil

}

func (h *UserHandler) Add(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var newUser = new(serviceUser.AddUser)
	if err := BindJSON(r, newUser); err != nil {
		fmt.Println(err)
		return
	}
	handler := &serviceUser.AddUserHandler{
		UserRepository: h.UserRepository,
	}
	_, err := handler.ADD(newUser)
	if err != nil {
		WriteJSON(w, HTTP_ERROR_CODE_ADD_FAILED, ResponseBody{
			Message: err.Error(),
			Code:    HTTP_ERROR_CODE_ADD_FAILED,
		})
		return
	}
	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "Add user succesfuly",
		Code:    http.StatusOK,
	})

}
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	query := r.URL.Query()
	start := query.Get("start")
	end := query.Get("end")
	limit := query.Get("limit")
	page := query.Get("page")
	Users, err := h.UserRepository.All()
	if err != nil {
		WriteJSON(w, HTTP_ERROR_CODE_READ_FAILED, ResponseBody{
			Message: "unable to get User from server",
		})
	}
	results, err := getByPage(start, end, limit, page, Users)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, ResponseBody{
			Message: err.Error(),
			Code:    http.StatusNotFound,
		})
		return
	}
	WriteJSON(w, http.StatusOK, results)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	User, err := h.UserRepository.FindByID(id)
	if err != nil {
		WriteJSON(w, HTTP_ERROR_CODE_READ_FAILED, ResponseBody{
			Message: "unable to get User by id : " + id,
		})
	}
	WriteJSON(w, http.StatusOK, User)
}

func (h *UserHandler) UpdateByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	var updateUsers model.User
	if err := BindJSON(r, updateUsers); err != nil {
		fmt.Println(err)
		return
	}
	if err := h.UserRepository.UpdateByID(id, updateUsers); err != nil {
		WriteJSON(w, HTTP_ERROR_CODE_UPDATE_FAILED, ResponseBody{
			Message: err.Error(),
			Code:    HTTP_ERROR_CODE_UPDATE_FAILED,
		})
		return
	}
	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "Update succesfully",
	})
}
func (h *UserHandler) RemoveByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if err := h.UserRepository.RemoveByID(id); err != nil {
		WriteJSON(w, HTTP_ERROR_CODE_DELETE_FAILED, ResponseBody{
			Message: err.Error(),
			Code:    HTTP_ERROR_CODE_DELETE_FAILED,
		})
		return
	}
	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "Delete succesfully",
	})
}
func (h *UserHandler) RemoveByPhoneNumber(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if err := h.UserRepository.RemoveByPhoneNumber(id); err != nil {
		WriteJSON(w, HTTP_ERROR_CODE_DELETE_FAILED, ResponseBody{
			Message: err.Error(),
			Code:    HTTP_ERROR_CODE_DELETE_FAILED,
		})
		return
	}
	WriteJSON(w, http.StatusOK, ResponseBody{
		Message: "Delete succesfully",
	})
}
