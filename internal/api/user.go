package api

import (
	"encoding/json"
	"net/http"

	"github.com/ferdiebergado/fullstackgo/internal/model"
	"github.com/ferdiebergado/fullstackgo/internal/service"
)

type UserHandler interface {
	HandleCreateUser(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	service service.UserService
}

var _ UserHandler = (*userHandler)(nil)

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{
		service: userService,
	}
}

func (h *userHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var params model.UserCreateParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newUser, err := h.service.CreateUser(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userJson, err := json.Marshal(newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(userJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
